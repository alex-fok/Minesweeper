package ws

import (
	"encoding/json"
	"log"
	"minesweeper/boardhelper"
	"net/http"

	"github.com/gorilla/websocket"
)

type V struct {
	X, Y int
}
type BlockInfo struct {
	X int `json:"x"`
	Y int `json:"y"`
	boardhelper.Block
}

var board = boardhelper.GetBoard(50, 26)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func getBlockArr(v *V) []BlockInfo {
	slice := []BlockInfo{}
	block := BlockInfo{
		X:     v.X,
		Y:     v.Y,
		Block: board[v.X][v.Y],
	}
	slice = append(slice, block)
	return slice
}

func ServeWs(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)

	defer conn.Close()

	if err != nil {
		log.Println(err)
		return
	}
	// for _, line := range board {
	// 	fmt.Println(line)
	// }

	for {
		var v V
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
