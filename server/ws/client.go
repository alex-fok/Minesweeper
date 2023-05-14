package ws

import (
	"encoding/json"
	"log"
	"math/rand"
	"minesweeper/boardhelper"

	"github.com/gorilla/websocket"
)

type Message struct {
	Name    string `json:"Name"`
	Content string `json:"Content"`
}

type Client struct {
	conn   *websocket.Conn
	room   *Room
	update chan *Action
}

func getBlockArr(m *Message) []boardhelper.BlockInfo {
	var v boardhelper.Vertex
	json.Unmarshal([]byte(m.Content), &v)
	return boardhelper.GetRevealables(&v, board)
}

func (c *Client) createRoom() {
	c.room = newRoom(uint(rand.Uint32()))
	go c.room.run()
	c.room.register <- c
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
		switch messageType := msg.Name; messageType {
		case "newGame":
			c.createRoom()
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

		// data, err := json.Marshal(resObj)

		// if err != nil {
		// 	log.Println(err)
		// 	return
		// }
		// dataArr := []string{string(data)}
		// if err := c.conn.WriteJSON(dataArr); err != nil {
		// 	log.Println(err)
		// 	return
		// }
	}
}
