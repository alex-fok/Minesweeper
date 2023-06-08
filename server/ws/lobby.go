package ws

import (
	"minesweeper/utils"
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

func (l *Lobby) createRoom(c *Client) *Room {
	var id uint
	for {
		id = utils.CreateRoomId()
		if _, ok := l.rooms[id]; !ok {
			break
		}
	}
	if c.room != nil {
		c.room.unregister <- c
	}
	r := newRoom(id, c, l)

	l.register <- r
	return r
}

func (l *Lobby) findRoom(id uint) (*Room, bool) {
	r, ok := l.rooms[id]
	return r, ok
}

func (l *Lobby) run() {
	for {
		select {
		case r := <-l.register:
			l.rooms[r.id] = r
		case r := <-l.unregister:
			r.stop <- true
			delete(l.rooms, r.id)
		}
	}
}
