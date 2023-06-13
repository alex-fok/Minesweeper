package ws

import (
	"encoding/json"
	"log"
	"minesweeper/game"
	"minesweeper/types"
	"time"
)

type Game = game.Game
type ClientMeta = types.ClientMeta

type GameEnded struct {
	Winner ClientId `json:"winner"`
}

type RoomUpdate struct {
	Client ClientId
	Action *Action
}

type Room struct {
	id            uint
	clients       map[ClientId]*Client
	lobby         *Lobby
	board         [][]game.Block
	actionHandler map[string]func(ClientId, string)
	gameDriver    game.Driver
	update        chan *RoomUpdate
	register      chan *Client
	unregister    chan *Client
	disconnect    chan *Client
	reconnect     chan *Client
	timeouts      map[ClientId]int64
	stop          chan bool
}

const DEFAULT_SIZE = 26
const DEFAULT_BOMB_COUNT = 100
const TIMELIMIT_IN_SEC = 5 * 60 // 5 minutes

func newRoom(id uint, c *Client, l *Lobby) *Room {
	r := &Room{
		id:         id,
		clients:    make(map[ClientId]*Client),
		lobby:      l,
		gameDriver: *game.NewDriver(),
		update:     make(chan *RoomUpdate),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		disconnect: make(chan *Client),
		reconnect:  make(chan *Client),
		timeouts:   make(map[ClientId]int64),
	}
	// Init handler
	r.actionHandler = make(map[string]func(ClientId, string))
	r.actionHandler["reveal"] = func(cId ClientId, content string) {
		actions := r.gameDriver.Reveal(cId, content)
		for _, a := range actions {
			r.broadcast(a)
		}
	}
	go r.run()

	// Register Client
	r.registerClient(c)

	// Update client
	c.update <- &Action{
		Name:    "roomCreated",
		Content: "{}",
	}
	return r
}

func (r *Room) registerClient(c *Client) {
	r.clients[c.id] = c

	type RoomIdMsg struct {
		Id uint `json:"id"`
	}
	rId, _ := json.Marshal(RoomIdMsg{Id: r.id})

	c.update <- &Action{
		Name:    "roomId",
		Content: string(rId),
	}
	actions := r.gameDriver.RegisterPlayer(&ClientMeta{
		Id:       c.id,
		Alias:    c.alias,
		IsOnline: true,
	})
	for _, a := range actions {
		r.broadcast(a)
	}
}

func (r *Room) unregisterClient(c *Client) {
	if _, ok := r.clients[c.id]; ok {
		delete(r.clients, c.id)
		r.gameDriver.UnregisterPlayer(c.id)
	}
	if len(r.clients) == 0 {
		r.stop <- true
		r.lobby.unregister <- r
	}
}

func (r *Room) disconnectClient(c *Client) {
	log.Println("Client", c.alias, "disconnected")
	if _, ok := r.timeouts[c.id]; !ok {
		r.timeouts[c.id] = time.Now().Unix()

		if action := r.gameDriver.DisconnectPlayer(c.id); action != nil {
			r.broadcast(action)
		}
	}
}

func (r *Room) reconnectClient(c *Client) {
	// Reconnect user if id is in timeout map
	_, timeoutOk := r.timeouts[c.id]
	_, clientsOk := r.clients[c.id]
	if !timeoutOk || !clientsOk {
		log.Println("Cannot reconnect client")
		c.update <- &Action{
			Name:    "reconnFailed",
			Content: "{}",
		}
		return
	}
	// Remove from timeout map
	delete(r.timeouts, c.id)

	// Notify all if client is player
	if action := r.gameDriver.ReconnectPlayer(c.id); action != nil {
		r.broadcast(action)
	}

	// Reuse alias before disconnection, reassign
	c.alias = r.clients[c.id].alias
	r.clients[c.id] = c
	c.room = r

	log.Println("Client", c.alias, "reconnected")

	// Return Game Stat
	stat, _ := json.Marshal(r.gameDriver.GetGameStat())

	c.update <- &Action{
		Name:    "gameStat",
		Content: string(stat),
	}
}

func (r *Room) checkActivity(now int64) {
	for cId, t := range r.timeouts {
		if now-t > TIMELIMIT_IN_SEC {
			log.Println("Removing client", r.clients[cId].alias, "from room", r.id)
			delete(r.clients, cId)
			delete(r.timeouts, cId)
		}
	}
}

func (r *Room) broadcast(action *Action) {
	for cId := range r.clients {
		if r.clients[cId].isOpen {
			r.clients[cId].update <- action
		}
	}
}

func (r *Room) run() {
	// Setup ticker
	ticker := time.NewTicker(time.Minute)
	doneChecking := make(chan bool)
	defer func() {
		doneChecking <- true
		ticker.Stop()
	}()

	// Check for client activity
	go func() {
		for {
			select {
			case t := <-ticker.C:
				r.checkActivity(t.Unix())
			case <-doneChecking:
				break
			}
		}
	}()

	for {
		select {
		case update := <-r.update:
			r.actionHandler[update.Action.Name](update.Client, update.Action.Content)
		case c := <-r.register:
			r.registerClient(c)
		case c := <-r.unregister:
			r.unregisterClient(c)
		case c := <-r.disconnect:
			r.disconnectClient(c)
		case c := <-r.reconnect:
			r.reconnectClient(c)
		case <-r.stop:
			break
		}
	}
}
