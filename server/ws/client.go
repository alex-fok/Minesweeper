package ws

import (
	"encoding/json"
	"log"
	"minesweeper/types"
	"minesweeper/utils"
	"strconv"

	"github.com/gorilla/websocket"
)

type Request struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type Message struct {
	Message string `json:"message"`
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
	return &Client{
		id:     ClientId(utils.CreateClientId()),
		conn:   conn,
		isOpen: true,
		lobby:  lobby,
		alias:  "Anonymous",
		update: make(chan *Action),
		stop:   make(chan bool),
	}
}

func (c *Client) reconnect(req *Request) {
	type ReconnectReq struct {
		UserId string `json:"userId"`
		RoomId string `json:"roomId"`
	}
	var reconnectReq ReconnectReq
	if err := json.Unmarshal([]byte(req.Content), &reconnectReq); err != nil {
		log.Println(err)
		return
	}
	c.id = ClientId(reconnectReq.UserId)

	roomNotFound := false

	if rId, err := strconv.ParseUint(reconnectReq.RoomId, 10, 64); err == nil {
		if r, ok := c.lobby.findRoom(uint(rId)); ok {
			r.reconnect <- c
		} else {
			roomNotFound = true
		}
	} else {
		roomNotFound = true
	}

	if roomNotFound {
		c.update <- &Action{
			Name:    "reconnFailed",
			Content: "{}",
		}
	}
}

func (c *Client) createRoom(req *Request) {
	type CreateReq struct {
		Alias string `json:"alias"`
	}
	var createReq CreateReq
	if err := json.Unmarshal([]byte(req.Content), &createReq); err != nil {
		log.Println(err)
		return
	}
	c.alias = createReq.Alias
	if c.room != nil {
		c.room.unregister <- c
	}

	c.room = c.lobby.createRoom(c)
	c.room.register <- c

	// Update client

	c.update <- &Action{
		Name:    "roomCreated",
		Content: string("{}"),
	}

	log.Println("Room", c.room.id, "created by Client", createReq.Alias)
}

func (c *Client) joinRoom(req *Request) {
	type JoinReq struct {
		Id    uint   `json:"id"`
		Alias string `json:"alias"`
	}
	var joinReq JoinReq
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
		message, _ := json.Marshal(&Message{
			Message: "Room #" + strconv.FormatUint(uint64(joinReq.Id), 10) + " not found",
		})
		c.update <- &Action{
			Name:    "message",
			Content: string(message),
		}
		return
	}
	log.Println("Client", joinReq.Alias, "joined Room", r.id)
	c.room = r
	c.room.register <- c
}

func (c *Client) handleInviteCode(req *Request) {
	type Inivitation struct {
		Id string `json:"id"`
	}

	var invitation Inivitation
	json.Unmarshal([]byte(req.Content), &invitation)
	if r := c.lobby.findInviteCode(invitation.Id); r != nil {
		c.room = r
		c.room.register <- c
	} else {
		c.update <- &Action{
			Name:    "reconnFailed",
			Content: "{}",
		}
	}
}

func (c *Client) rename(req *Request) {
	type RenameReq struct {
		Alias string `json:"alias"`
	}
	var renameReq RenameReq
	if err := json.Unmarshal([]byte(req.Content), &renameReq); err != nil {
		log.Println(err)
		return
	}
	c.alias = renameReq.Alias
	if c.room == nil {
		return
	}
	c.updateRoom(req)
}

func (c *Client) updateRoom(req *Request) {
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
		c.isOpen = false
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
		c.conn.Close()
		c.isOpen = false
		close(c.update)

		if c.room != nil {
			c.room.disconnect <- c
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
		case "inviteCode":
			go c.handleInviteCode(&req)
		case "rename":
			go c.rename(&req)
		case "reveal", "rematch":
			go c.updateRoom(&req)
		default:
			continue
		}
	}
}

func (c *Client) run() {
	go c.readBuffer()
	go c.writeBuffer()
	type IdMsg struct {
		Id string `json:"id"`
	}
	id, _ := json.Marshal(&IdMsg{
		Id: string(c.id),
	})
	c.update <- &Action{
		Name:    "userId",
		Content: string(id),
	}
}
