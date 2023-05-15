package ws

import (
	"math/rand"
)

type Lobby struct {
	rooms      map[uint]*Room
	register   chan *Room
	unregister chan *Room
}

var lobby *Lobby = nil

func CreateLobby() *Lobby {
	if lobby != nil {
		return lobby
	}
	lobby = &Lobby{
		rooms:      make(map[uint]*Room),
		register:   make(chan *Room),
		unregister: make(chan *Room),
	}
	go lobby.run()
	return lobby
}

func (l *Lobby) createRoom(c *Client) {
	id := uint(rand.Uint32())
	c.room = newRoom(id)
	go c.room.run()
	l.register <- c.room
	c.room.register <- c
}

func (l *Lobby) findRoom(id uint) *Room {
	return l.rooms[id]
}

func (l *Lobby) run() {
	for {
		select {
		case r := <-l.register:
			l.rooms[r.id] = r
		case r := <-l.unregister:
			delete(l.rooms, r.id)
		}
	}
}
