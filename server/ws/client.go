package ws

import (
	"log"

	"github.com/gorilla/websocket"
)

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

func (c *Client) handleAction(action *Action) string {
	return action.Content
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
			if data := c.handleAction(action); len(data) != 0 {
				if err := c.conn.WriteJSON(data); err != nil {
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
			lobby.createRoom(c)
		case "reveal":
			if c.room == nil {
				continue
			}
			c.room.update <- &Action{
				Name:    msg.Name,
				Content: msg.Content,
			}

		default:
			continue
		}
	}
}
