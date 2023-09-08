package ws

import (
	"minesweeper/utils"
)

type Lobby struct {
	rooms       map[uint]*Room
	inviteCodes map[string]uint
	register    chan *Room
	unregister  chan *Room
}

var lobby *Lobby = nil

func CreateLobby() *Lobby {
	if lobby != nil {
		return lobby
	}
	lobby = &Lobby{
		rooms:       make(map[uint]*Room),
		inviteCodes: make(map[string]uint),
		register:    make(chan *Room),
		unregister:  make(chan *Room),
	}
	go lobby.run()
	return lobby
}

func (l *Lobby) createRoom(c *Client, config *RoomConfig) *Room {
	var id uint
	for {
		id = utils.CreateRoomId()
		if _, ok := l.rooms[id]; !ok {
			break
		}
	}
	r := newRoom(id, c, l, config)

	l.register <- r
	return r
}

func (l *Lobby) findRoom(id uint) (*Room, bool) {
	r, ok := l.rooms[id]
	return r, ok
}

func (l *Lobby) createInviteCode(rId uint) string {
	id := utils.CreateInvitationId()
	l.inviteCodes[id] = rId
	return id
}

func (l *Lobby) findInviteCode(id string) *Room {
	if rId, idOk := l.inviteCodes[id]; idOk {
		if r, ok := l.findRoom(rId); ok {
			return r
		}
		return nil
	}
	return nil
}

func (l *Lobby) getPublicRIds() []uint {
	roomIds := []uint{}
	for _, r := range l.rooms {
		if r.roomType == "public" {
			roomIds = append(roomIds, r.id)
		}
	}
	return roomIds
}

func (l *Lobby) run() {
	for {
		select {
		case r := <-l.register:
			l.rooms[r.id] = r
		case r := <-l.unregister:
			r.stop <- true
			delete(l.rooms, r.id)
			delete(l.inviteCodes, r.inviteCode)
		}
	}
}
