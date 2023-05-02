package main

import (
	"fmt"
	boardhelper "minesweeper/boardhelper"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

var board [][]boardhelper.Block

func getServer(server *socketio.Server) *socketio.Server {
	fmt.Println("httpServer")
	return server
}
func main() {
	router := gin.Default()
	server := socketio.NewServer(nil)
	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
		return
	})
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})

	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		fmt.Println("notice:", msg)
		s.Emit("reply", "have"+msg)
		return
	})
	server.OnEvent("/", "startGame", func(s socketio.Conn) {
		board = boardhelper.GetBoard(50, 16) // 16 * 16 Board
		fmt.Println(board)
		s.Emit("board", board)
		return
	})
	go server.Serve()
	defer server.Close()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
	}))

	router.GET("/socket.io/", gin.WrapH(getServer(server)))
	router.POST("/socket.io/", gin.WrapH(getServer(server)))
	router.Run("localhost:8080")
}
