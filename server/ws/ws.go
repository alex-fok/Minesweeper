package ws

import (
	"encoding/json"
	"log"
	"minesweeper/boardhelper"
	"net/http"

	"github.com/gorilla/websocket"
)

const DEFAULT_SIZE = 26
const DEFAULT_BOMB_COUNT = 100

var board = boardhelper.GetBoard(DEFAULT_BOMB_COUNT, DEFAULT_SIZE)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func getBlockArr(v *boardhelper.Vertex) []boardhelper.BlockInfo {
	return boardhelper.GetRevealables(v, board, DEFAULT_SIZE)
}

func ServeWs(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)

	defer conn.Close()

	if err != nil {
		log.Println(err)
		return
	}
	for {
		var v boardhelper.Vertex
		err := conn.ReadJSON(&v)
		if err != nil {
			log.Println(err)
			return
		}
		data, err := json.Marshal(getBlockArr(&v))

		if err != nil {
			log.Println(err)
			return
		}
		dataArr := []string{string(data)}
		if err := conn.WriteJSON(dataArr); err != nil {
			log.Println(err)
			return
		}
	}
}
