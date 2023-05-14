package ws

import (
	"encoding/json"
	"log"
	"minesweeper/boardhelper"
)

type Action struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type Room struct {
	id      uint
	clients map[*Client]bool

	handler map[string]func(string)
	turn    struct {
		number uint
		player *Client
	}
	update     chan *Action
	register   chan *Client
	unregister chan *Client
}

func (r *Room) revealBlocks(content string) {
	var v boardhelper.Vertex
	json.Unmarshal([]byte(content), &v)
	data, err := json.Marshal(boardhelper.GetRevealables(&v, board))
	if err != nil {
		log.Println(err)
		return
	}
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

func newRoom(id uint) *Room {
	room := Room{
		id:         id,
		clients:    make(map[*Client]bool),
		update:     make(chan *Action),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

	room.handler = make(map[string]func(string))
	room.handler["reveal"] = room.revealBlocks

	return &room
}

func (r *Room) run() {
	for {
		select {
		case c := <-r.register:
			r.clients[c] = true
		case c := <-r.unregister:
			if _, ok := r.clients[c]; ok {
				delete(r.clients, c)
				close(c.update)
			}
		case action := <-r.update:
			// for c := range r.clients {
			// 	c.update(action)
			// }
			r.handler[action.Name](action.Content)
		}
	}
}
