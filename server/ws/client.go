package ws

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

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
	update chan *Action
}

func newClient(conn *websocket.Conn, lobby *Lobby) *Client {
	return &Client{
		conn:   conn,
		lobby:  lobby,
		update: make(chan *Action),
	}
}

func (c *Client) newGame() {
	c.lobby.createRoom(c)
}

func (c *Client) joinGame(msg *Message) {
	var joinReq JoinRequest
	if err := json.Unmarshal([]byte(msg.Content), &joinReq); err != nil {
		log.Println(err)
		return
	}
	c.lobby.findRoom(joinReq.Id)
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
		case "newGame":
			go c.newGame()
		case "joinGame":
			go c.joinGame(&msg)
		case "reveal":
			go c.reveal(&msg)
		default:
			continue
		}
	}
}
