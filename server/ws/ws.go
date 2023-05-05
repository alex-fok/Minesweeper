package ws

import (
	"encoding/json"
	"log"
	"minesweeper/boardhelper"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ServeWs(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)

	defer conn.Close()

	if err != nil {
		log.Println(err)
		return
	}
	board := boardhelper.GetBoard(50, 26)
	// for _, line := range board {
	// 	fmt.Println(line)
	// }
	type V struct {
		X, Y int
	}

	type BlockInfo struct {
		X int `json:"x"`
		Y int `json:"y"`
		boardhelper.Block
	}

	for {
		var v V
		err := conn.ReadJSON(&v)
		if err != nil {
			log.Println(err)
			return
		}
		block := BlockInfo{
			X:     v.X,
			Y:     v.Y,
			Block: board[v.Y][v.X],
		}
		data, err := json.Marshal(block)

		if err != nil {
			log.Println(err)
			return
		}
		if err := conn.WriteJSON(string(data)); err != nil {
			log.Println(err)
			return
		}
	}
}
