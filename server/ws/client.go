package ws

import (
	"encoding/json"
	"log"
	"minesweeper/types"
	"minesweeper/utils"
	"strconv"

	"github.com/gorilla/websocket"
)

type ReconnectReq struct {
	UserId string `json:"userId"`
	RoomId string `json:"roomId"`
}

type CreateRequest struct {
	Alias string `json:"alias"`
}

type JoinRequest struct {
	Id    uint   `json:"id"`
	Alias string `json:"alias"`
}

type IdMsg struct {
	Id string `json:"id"`
}

type Request struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type ClientId = types.ClientId

type Client struct {
	id     ClientId
	conn   *websocket.Conn
	isOpen bool
	lobby  *Lobby
	room   *Room
	alias  string
	update chan *Action
	stop   chan bool
}

func NewClient(conn *websocket.Conn, lobby *Lobby) *Client {
	c := &Client{
		id:     ClientId(utils.CreateClientId()),
		conn:   conn,
		isOpen: true,
		lobby:  lobby,
		alias:  "Anonymous",
		update: make(chan *Action),
		stop:   make(chan bool),
	}
	socketCloseHandler := c.conn.CloseHandler()
	c.conn.SetCloseHandler(func(code int, text string) error {
		c.isOpen = false
		return socketCloseHandler(code, text)
	})
	return c
}

func (c *Client) reconnect(req *Request) {
	var reconnectReq ReconnectReq
	if err := json.Unmarshal([]byte(req.Content), &reconnectReq); err != nil {
		log.Println(err)
		return
	}
	c.id = ClientId(reconnectReq.UserId)
	if rId, err := strconv.ParseUint(reconnectReq.RoomId, 10, 64); err == nil {
		if r, ok := c.lobby.findRoom(uint(rId)); ok {
			r.reconnect <- c
		} else {
			log.Println("Room", reconnectReq.RoomId, "not found")
		}
	}
}

func (c *Client) createRoom(req *Request) {
	var createReq CreateRequest
	if err := json.Unmarshal([]byte(req.Content), &createReq); err != nil {
		log.Println(err)
		return
	}
	c.alias = createReq.Alias
	if c.room != nil {
		c.room.unregister <- c
	}
	c.room = c.lobby.createRoom(c)

	log.Println("Room", c.room.id, "created by Client", createReq.Alias)
}

func (c *Client) joinRoom(req *Request) {
	var joinReq JoinRequest
	if err := json.Unmarshal([]byte(req.Content), &joinReq); err != nil {
		log.Println(err)
		return
	}
	if c.room != nil && c.room.id == joinReq.Id {
		if c.room.id == joinReq.Id {
			return
		} else {
			c.room.unregister <- c
		}
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

func (c *Client) reveal(req *Request) {
	if c.room == nil {
		return
	}
	c.room.update <- &RoomUpdate{
		Client: c.id,
		Action: &Action{
			Name:    req.Name,
			Content: req.Content,
		},
	}
}

func (c *Client) writeBuffer() {
	defer func() {
		c.conn.Close()
		close(c.update)
	}()
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
		case <-c.stop:
			break
		}
	}
}

func (c *Client) readBuffer() {
	defer func() {
		if c.room != nil {
			c.room.disconnect <- c
			c.conn.Close()
		}
	}()

	var req Request
	for {
		// Get Message Type
		if err := c.conn.ReadJSON(&req); err != nil {
			log.Println(err)
			c.stop <- true
			return
		}
		switch req.Name {
		case "reconnect":
			go c.reconnect(&req)
		case "createRoom":
			go c.createRoom(&req)
		case "joinRoom":
			go c.joinRoom(&req)
		case "reveal":
			go c.reveal(&req)
		default:
			continue
		}
	}
}

func (c *Client) run() {
	go c.readBuffer()
	go c.writeBuffer()

	id, _ := json.Marshal(&IdMsg{
		Id: string(c.id),
	})
	c.update <- &Action{
		Name:    "userId",
		Content: string(id),
	}
}
