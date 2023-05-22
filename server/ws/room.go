package ws

import (
	"encoding/json"
	"math/rand"
	"minesweeper/boardhelper"
	"time"
)

type Action struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type Turn struct {
	count uint
	curr  *Client
	next  *Client
}

type Room struct {
	id            uint
	clients       map[*Client]bool
	lobby         *Lobby
	clientCount   uint
	board         [][]boardhelper.Block
	actionHandler map[string]func(string)
	turn          Turn
	update        chan *Action
	register      chan *Client
	unregister    chan *Client
}

const DEFAULT_SIZE = 26
const DEFAULT_BOMB_COUNT = 100

func newRoom(id uint, c *Client, l *Lobby) *Room {
	room := Room{
		id:          id,
		clients:     make(map[*Client]bool),
		lobby:       l,
		clientCount: 0,
		board:       boardhelper.GetBoard(DEFAULT_SIZE, DEFAULT_BOMB_COUNT),
		turn: Turn{
			count: 0,
			curr:  nil,
			next:  nil,
		},
		update:     make(chan *Action),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	// Init handler
	room.actionHandler = make(map[string]func(string))
	room.actionHandler["reveal"] = room.revealBlocks

	// Init turn-related info
	room.turn.count = 1

	room.assignTurn(c)
	room.clients[c] = true
	return &room
}

func (r *Room) assignTurn(c *Client) {
	isEmpty := r.turn.curr == nil && r.turn.next == nil
	if isEmpty {
		rand.Seed(time.Now().UnixNano())
		if rand.Intn(2) == 0 {
			r.turn.curr = c
		} else {
			r.turn.next = c
		}
	} else if r.turn.curr == nil {
		r.turn.curr = c
	} else if r.turn.next == nil { // Not 'else'. Could be more than 3 person in a room
		r.turn.next = c
	}
}

func (r *Room) isPlayable(c *Client) bool {
	return r.turn.curr == c && r.turn.next != nil
}

func (r *Room) registerClient(c *Client) {
	r.clients[c] = true
	r.assignTurn(c)
	r.clientCount++
}

func (r *Room) unregisterClient(c *Client) {
	if _, ok := r.clients[c]; ok {
		delete(r.clients, c)
		if r.turn.curr == c {
			r.turn.curr = nil
		} else if r.turn.next == c {
			r.turn.next = nil
		}
	}
	if r.clientCount--; r.clientCount == 0 {
		r.lobby.unregister <- r
	}
}

func (r *Room) revealBlocks(content string) {
	// Get revealable blocks
	var v boardhelper.Vertex
	json.Unmarshal([]byte(content), &v)
	if r.board[v.Y][v.X].Visited {
		return
	}
	revealables := boardhelper.GetRevealables(&v, r.board)

	// Update visited blocks
	for _, block := range revealables {
		r.board[block.Y][block.X].Visited = true
	}

	// Update turn-related info
	if !(revealables[0].Type == boardhelper.BOMB) {
		r.turn.count++
		r.turn.next, r.turn.curr = r.turn.curr, r.turn.next
		turn, _ := json.Marshal(struct {
			Count uint `json:"count"`
		}{Count: r.turn.count})
		r.broadcast(&Action{
			Name:    "turnPassed",
			Content: string(turn),
		})
	}

	data, _ := json.Marshal(struct {
		Blocks []boardhelper.BlockInfo `json:"blocks"`
	}{
		Blocks: revealables,
	})

	// Broadcast to clients
	r.broadcast(&Action{
		Name:    "reveal",
		Content: string(data),
	})
}

func (r *Room) broadcast(action *Action) {
	for client := range r.clients {
		client.update <- action
	}
}

func (r *Room) run() {
	for {
		select {
		case c := <-r.register:
			r.registerClient(c)
		case c := <-r.unregister:
			r.unregisterClient(c)
		case action := <-r.update:
			r.actionHandler[action.Name](action.Content)
		}
	}
}
