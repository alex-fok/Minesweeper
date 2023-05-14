package ws

import (
	"log"
	"minesweeper/boardhelper"
	"net/http"

	"github.com/gorilla/websocket"
)

const DEFAULT_SIZE = 26
const DEFAULT_BOMB_COUNT = 100

type room struct {
	conn []*websocket.Conn
	id   uint
}

var board = boardhelper.GetBoard(DEFAULT_BOMB_COUNT, DEFAULT_SIZE)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ServeWs(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{
		conn:   conn,
		update: make(chan *Action),
	}
	go client.readBuffer()
	go client.writeBuffer()
}
