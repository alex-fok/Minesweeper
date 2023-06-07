package ws

import (
	"encoding/json"
	"log"
	"math/rand"
	"strings"

	"github.com/gorilla/websocket"
)

const idLetters string = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const bits = 6               // Number of bits representing an index; Log2(len(idLetters)) + 1
const mask = 1<<bits - 1     // Masking for 6 bit value
const count = int(63 / bits) // Number of indices per 63 bit

// Generate 16 letter ID
func generateId() string {
	sb := strings.Builder{}
	sb.Grow(16)
	for idIdx, value, remain := 0, rand.Int63(), count; idIdx < 16; {
		if remain == 0 {
			value, remain = rand.Int63(), count
		}
		if letterIdx := int(value & mask); letterIdx < len(idLetters) {
			sb.WriteByte(idLetters[letterIdx])
			idIdx++
		}
		value >>= bits
		remain--
	}
	return sb.String()
}

type CreateRequest struct {
	Alias string `json:"alias"`
}
type JoinRequest struct {
	Id    uint   `json:"id"`
	Alias string `json:"alias"`
}
type Message struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type Client struct {
	id     string
	conn   *websocket.Conn
	lobby  *Lobby
	room   *Room
	alias  string
	update chan *Action
}

func NewClient(conn *websocket.Conn, lobby *Lobby) *Client {
	return &Client{
		id:     generateId(),
		conn:   conn,
		lobby:  lobby,
		alias:  "Anonymous",
		update: make(chan *Action),
	}
}

func (c *Client) createRoom(msg *Message) {
	var createReq CreateRequest
	if err := json.Unmarshal([]byte(msg.Content), &createReq); err != nil {
		log.Println(err)
		return
	}
	c.alias = createReq.Alias
	c.lobby.createRoom(c)
	log.Println("Room", c.room.id, "created by Client", createReq.Alias)
}

func (c *Client) joinRoom(msg *Message) {
	var joinReq JoinRequest
	if err := json.Unmarshal([]byte(msg.Content), &joinReq); err != nil {
		log.Println(err)
		return
	}
	if c.room != nil && c.room.id == joinReq.Id {
		return
	}
	c.alias = joinReq.Alias
	// Find Room. Register user if valid
	r, ok := c.lobby.findRoom(joinReq.Id)
	if !ok {
		log.Println("Room", joinReq.Id, "not found")
		return
	}
	log.Println("Client", joinReq.Alias, "joined Room", r.id)
	c.room = r
	c.room.register <- c
}

func (c *Client) reveal(msg *Message) {
	if c.room == nil || !c.room.isPlayable(c) {
		return
	}
	c.room.update <- &Action{
		Name:    msg.Name,
		Content: msg.Content,
	}
}

func (c *Client) writeBuffer() {
	defer c.conn.Close()
	for {
		select {
		case action, ok := <-c.update:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if len(action.Content) != 0 {
				if err := c.conn.WriteJSON(action); err != nil {
					log.Println(err)
					return
				}
			}
		}
	}
}

func (c *Client) readBuffer() {
	defer func() {
		if c.room != nil {
			c.room.unregister <- c
			c.conn.Close()
			close(c.update)
		}
	}()
	for {
		// Get Message Type
		var msg Message
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			log.Println(err)
			return
		}
		switch msg.Name {
		case "createRoom":
			go c.createRoom(&msg)
		case "joinRoom":
			go c.joinRoom(&msg)
		case "reveal":
			go c.reveal(&msg)
		default:
			continue
		}
	}
}
