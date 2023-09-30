package ws

import (
	"encoding/json"
	"log"
	"minesweeper/types"
	"minesweeper/utils"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type Request struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type Message struct {
	Message string `json:"message"`
}

type Socket struct {
	conn   *websocket.Conn
	IsOpen bool
	close  func()
}

type Writer struct {
	update chan *Action
	ping   chan bool
	stop   chan bool
}

type ClientId = types.ClientId

type Client struct {
	id     ClientId
	socket *Socket
	lobby  *Lobby
	room   *Room
	alias  string
	writer *Writer
}

func NewClient(conn *websocket.Conn, lobby *Lobby) *Client {
	socket := &Socket{
		conn:   conn,
		IsOpen: true,
	}
	socket.close = func() {
		if !socket.IsOpen {
			return
		}
		socket.conn.Close()
		socket.IsOpen = false
	}
	return &Client{
		id:     ClientId(utils.CreateClientId()),
		socket: socket,
		lobby:  lobby,
		alias:  "Anonymous",
		writer: &Writer{
			update: make(chan *Action),
			ping:   make(chan bool),
			stop:   make(chan bool),
		},
	}
}

func (c *Client) reconnect(req *Request) {
	type ReconnReq struct {
		UserId string `json:"userId"`
		RoomId string `json:"roomId"`
	}

	var reconnReq ReconnReq
	if err := json.Unmarshal([]byte(req.Content), &reconnReq); err != nil {
		log.Println(err)
		return
	}
	c.id = ClientId(reconnReq.UserId)

	if rId, err := strconv.ParseUint(reconnReq.RoomId, 10, 64); err == nil {
		if r, ok := c.lobby.findRoom(uint(rId)); ok {
			c.room = r
			c.room.reconnect <- c
			return
		}
	}
	c.writer.update <- &Action{
		Name:    "reconnFailed",
		Content: "{}",
	}
}

func (c *Client) createRoom(req *Request) {
	type CreateReq struct {
		Alias     string `json:"alias"`
		RoomType  string `json:"roomType"`
		Pass      string `json:"passcode"`
		Capacity  int    `json:"capacity"`
		Player    uint   `json:"player"`
		Size      uint   `json:"size"`
		Bomb      uint   `json:"bomb"`
		TimeLimit uint   `json:"timeLimit"`
	}
	var createReq CreateReq
	if err := json.Unmarshal([]byte(req.Content), &createReq); err != nil {
		log.Println(err)
		return
	}
	c.alias = createReq.Alias
	if c.room != nil {
		c.room.unregister <- c.id
	}
	// Restrict roomType to equal to either "public" or "private"
	var roomType string
	if createReq.RoomType == "public" {
		roomType = createReq.RoomType
	} else {
		roomType = "private"
	}
	c.room = c.lobby.createRoom(c, &RoomConfig{
		Type:      roomType,
		Pass:      createReq.Pass,
		Capacity:  createReq.Capacity,
		Player:    createReq.Player,
		Size:      createReq.Size,
		Bomb:      createReq.Bomb,
		TimeLimit: createReq.TimeLimit,
	})
	c.room.register <- &RoomLogin{
		client: c,
		pass:   createReq.Pass,
	}
	log.Println("Room", c.room.id, "created by Client", createReq.Alias)
}

func (c *Client) joinRoom(req *Request) {
	type JoinReq struct {
		Alias    string `json:"alias"`
		RoomType string `json:"roomType"`
		Id       uint   `json:"id"`
		Passcode string `jon:"passcode"`
	}
	var joinReq JoinReq
	if err := json.Unmarshal([]byte(req.Content), &joinReq); err != nil {
		log.Println(err)
		return
	}
	if c.room != nil {
		if c.room.id == joinReq.Id {
			return
		} else {
			c.room.unregister <- c.id
		}
	}
	c.alias = joinReq.Alias
	// Find Room. Register user if valid
	r, ok := c.lobby.findRoom(joinReq.Id)

	mismatch := joinReq.RoomType == "public" && r.roomType == "private"

	// Room not found, OR
	// Mismatch room type
	if mismatch || !ok {
		log.Println("Room", joinReq.Id, "not found")
		message, _ := json.Marshal(&Message{
			Message: "Room #" + strconv.FormatUint(uint64(joinReq.Id), 10) + " not found",
		})
		c.writer.update <- &Action{
			Name:    "message",
			Content: string(message),
		}
		return
	}
	c.room = r

	if c.room.roomType == "private" {
		c.room.register <- &RoomLogin{
			client: c,
			pass:   joinReq.Passcode,
		}
	} else {
		c.room.register <- &RoomLogin{
			client: c,
			pass:   "",
		}
		log.Println("Client", joinReq.Alias, "joined Room", r.id)
	}
}

func (c *Client) findPublic(req *Request) {
	// Request should be empty
	type RoomIds struct {
		Rooms []PublicRoomInfo `json:"rooms"`
	}

	rIds, _ := json.Marshal(RoomIds{
		Rooms: c.lobby.getPublicRIds(),
	})

	c.writer.update <- &Action{
		Name:    "publicRoomIds",
		Content: string(rIds),
	}
}

func (c *Client) handleInviteCode(req *Request) {
	type Inivitation struct {
		Id string `json:"id"`
	}
	var invitation Inivitation
	json.Unmarshal([]byte(req.Content), &invitation)
	if r := c.lobby.findInviteCode(invitation.Id); r != nil {
		c.room = r
		if c.room.roomType == "private" {
			c.writer.update <- &Action{
				Name:    "passcode",
				Content: "{}",
			}
		} else {
			c.room.register <- &RoomLogin{
				client: c,
				pass:   "",
			}
		}
	} else {
		c.writer.update <- &Action{
			Name:    "reconnFailed",
			Content: "{}",
		}
	}
}

func (c *Client) confirmPasscode(req *Request) {
	type Login struct {
		Pass string `json:"passcode"`
	}
	var login Login
	json.Unmarshal([]byte(req.Content), &login)
	if c.room != nil {
		c.room.register <- &RoomLogin{
			client: c,
			pass:   login.Pass,
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

func (c *Client) keepAlive() {
	ticker := time.NewTicker(15 * time.Second)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case <-ticker.C:
			if c.socket.IsOpen {
				c.writer.ping <- true
			} else {
				break
			}
		}
	}
}

func (c *Client) writeBuffer() {
	defer func() {
		close(c.writer.update)
		close(c.writer.ping)
		close(c.writer.stop)
	}()

	for {
		select {
		case action, ok := <-c.writer.update:
			if !ok {
				c.socket.conn.WriteMessage(websocket.CloseMessage, []byte{})
				c.socket.close()
			}
			if len(action.Content) != 0 {
				if err := c.socket.conn.WriteJSON(action); err != nil {
					log.Println(err)
					c.socket.close()
				}
			}
		case <-c.writer.ping:
			if err := c.socket.conn.WriteMessage(websocket.PingMessage, []byte("keep alive")); err != nil {
				log.Println(err)
				c.socket.close()
			}
		case <-c.writer.stop:
			break
		}
	}
}

func (c *Client) readBuffer() {
	defer func() {
		c.writer.stop <- true

		if c.room != nil {
			c.room.disconnect <- c.id
		}
	}()

	var req Request
	for {
		// Get Message Type
		if err := c.socket.conn.ReadJSON(&req); err != nil {
			log.Println(err)
			break
		}
		if !c.socket.IsOpen {
			break
		}
		switch req.Name {
		case "reconnect":
			go c.reconnect(&req)
		case "createRoom":
			go c.createRoom(&req)
		case "joinRoom":
			go c.joinRoom(&req)
		case "findPublicRoomIds":
			go c.findPublic(&req)
		case "inviteCode":
			go c.handleInviteCode(&req)
		case "passcode":
			go c.confirmPasscode(&req)
		case "rename":
			go c.rename(&req)
		case "share", "reveal", "rematch", "ready":
			go c.updateRoom(&req)
		default:
			continue
		}
	}
}

func (c *Client) run() {
	go c.readBuffer()
	go c.writeBuffer()
	go c.keepAlive()
	type IdMsg struct {
		Id string `json:"id"`
	}
	id, _ := json.Marshal(&IdMsg{
		Id: string(c.id),
	})
	c.writer.update <- &Action{
		Name:    "userId",
		Content: string(id),
	}
}
