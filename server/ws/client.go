package ws

import (
	"encoding/json"
	"log"
	"minesweeper/utils"

	"github.com/gorilla/websocket"
)

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
		id:     utils.CreateClientId(),
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
