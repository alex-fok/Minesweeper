package ws

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type CreateRequest struct {
	Alias string `json:"alias"`
}
type JoinRequest struct {
	Id uint `json:"id"`
}
type Message struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type Client struct {
	conn   *websocket.Conn
	lobby  *Lobby
	room   *Room
	alias  string
	update chan *Action
}

func newClient(conn *websocket.Conn, lobby *Lobby) *Client {
	return &Client{
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
}

func (c *Client) joinRoom(msg *Message) {
	var joinReq JoinRequest
	if err := json.Unmarshal([]byte(msg.Content), &joinReq); err != nil {
		log.Println(err)
		return
	}
	// Verify Room. Register user if valid
	r, ok := c.lobby.findRoom(joinReq.Id)
	if !ok {
		return
	}
	c.room = r
	c.room.register <- c

	isClientCurr := c.room.turn.curr == c

	// Send Client 'roomJoined' message
	joinedMsg, _ := json.Marshal(struct {
		IsPlayerTurn bool `json:"isPlayerTurn"`
	}{
		IsPlayerTurn: isClientCurr,
	})
	c.update <- &Action{
		Name:    "roomJoined",
		Content: string(joinedMsg),
	}

	// Send opponent 'newPlayer' message
	newPlayerMsg, _ := json.Marshal(struct {
	}{})
	var opponent *Client
	if isClientCurr {
		opponent = c.room.turn.next
	} else {
		opponent = c.room.turn.curr
	}
	opponent.update <- &Action{
		Name:    "newPlayer",
		Content: string(newPlayerMsg),
	}
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
