package main

import (
	"boardhelper"
	"fmt"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

var board [][]boardhelper.Block

func main() {
	router := gin.Default()
	server := socketio.NewServer(nil)
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})

	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		fmt.Println("notice:", msg)
		s.Emit("reply", "have"+msg)
	})
	server.OnEvent("/", "startGame", func(s socketio.Conn) {
		board = boardhelper.GetBoard(50, 16) // 16 * 16 Board

	})
	go server.Serve()
	defer server.Close()

	router.GET("/socket.io/", gin.WrapH(server))

	router.Run("localhost:8080")
}
