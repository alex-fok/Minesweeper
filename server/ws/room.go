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
	id                uint
	clients           map[ClientId]*Client
	lobby             *Lobby
	board             [][]game.Block
	gameUpdateHandler map[string]func(ClientId, string) []*Action
	gameDriver        game.Driver
	inviteCode        string
	update            chan *RoomUpdate
	register          chan *Client
	unregister        chan *Client
	disconnect        chan *Client
	reconnect         chan *Client
	timeouts          map[ClientId]int64
	stop              chan bool
}

const TIMELIMIT_IN_SEC = 60 * 5 // 5 minutes

func newRoom(id uint, c *Client, l *Lobby) *Room {
	r := &Room{
		id:         id,
		clients:    make(map[ClientId]*Client),
		lobby:      l,
		gameDriver: *game.NewDriver(),
		inviteCode: l.createInviteCode(id),
		update:     make(chan *RoomUpdate),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		disconnect: make(chan *Client),
		reconnect:  make(chan *Client),
		timeouts:   make(map[ClientId]int64),
	}

	// Init game update handler
	r.gameUpdateHandler = make(map[string]func(ClientId, string) []*Action)
	r.gameUpdateHandler["reveal"] = r.gameDriver.Reveal
	r.gameUpdateHandler["rematch"] = r.gameDriver.Rematch
	go r.run()

	return r
}

func (r *Room) registerClient(c *Client) {
	r.clients[c.id] = c

	// Notify room info
	r.notifyRoomInfo(c)

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
		c.writer.update <- &Action{
			Name:    "reconnFailed",
			Content: "{}",
		}
		return
	}
	// Remove from timeout map
	delete(r.timeouts, c.id)

	// Notify room info
	r.notifyRoomInfo(c)

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

	c.writer.update <- &Action{
		Name:    "gameStat",
		Content: string(stat),
	}
}

func (r *Room) notifyRoomInfo(c *Client) {
	type RoomInfo struct {
		Id         uint   `json:"id"`
		InviteCode string `json:"inviteCode"`
	}
	rInfo, _ := json.Marshal(RoomInfo{
		Id:         r.id,
		InviteCode: r.inviteCode,
	})

	c.writer.update <- &Action{
		Name:    "roomInfo",
		Content: string(rInfo),
	}
}

func (r *Room) rename(client ClientId, content string) {
	type Req struct {
		Alias string `json:"alias"`
	}
	var req Req
	json.Unmarshal([]byte(content), &req)
	r.gameDriver.RenamePlayer(client, req.Alias)

	type PlayerAlias struct {
		Client ClientId `json:"client"`
		Alias  string   `json:"alias"`
	}
	alias, _ := json.Marshal(PlayerAlias{
		Client: client,
		Alias:  req.Alias,
	})
	r.broadcast(&Action{
		Name:    "playerAlias",
		Content: string(alias),
	})
}

func (r *Room) handleGameUpdate(client ClientId, action *Action) {
	actions := r.gameUpdateHandler[action.Name](client, action.Content)
	for _, a := range actions {
		r.broadcast(a)
	}
}

func (r *Room) handleShare(client ClientId, action *Action) {
	type Content struct {
		Name    string
		Content string
	}

	var req Content
	innerContent := map[string]interface{}{}

	json.Unmarshal([]byte(action.Content), &req)
	json.Unmarshal([]byte(req.Content), &innerContent)

	innerContent["id"] = client

	contentBytes, _ := json.Marshal(innerContent)
	r.broadcast(&Action{
		Name:    req.Name,
		Content: string(contentBytes),
	})
}

func (r *Room) handleRoomUpdate(update *RoomUpdate) {
	switch update.Action.Name {
	case "rename":
		r.rename(update.Client, update.Action.Content)
	case "reveal", "rematch":
		r.handleGameUpdate(update.Client, update.Action)
	case "share":
		r.handleShare(update.Client, update.Action)
	default:
		return
	}

}

func (r *Room) broadcast(action *Action) {
	for cId := range r.clients {
		if r.clients[cId].socket.IsOpen {
			r.clients[cId].writer.update <- action
		}
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
			r.handleRoomUpdate(update)
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
