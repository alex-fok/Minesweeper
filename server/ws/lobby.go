package ws

import (
	"math/rand"
	"strconv"
	"time"
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
	var id uint
	for {
		rand.Seed(time.Now().UnixNano())
		id = uint(rand.Uint32() % 10000)
		if _, ok := l.rooms[id]; !ok {
			break
		}
	}

	c.room = newRoom(id, c, l)
	go c.room.run()

	l.register <- c.room
	c.update <- &Action{
		Name:    "room_id",
		Content: strconv.FormatUint(uint64(id), 10),
	}
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
