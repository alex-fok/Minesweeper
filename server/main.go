package main

import (
	"log"
	"math/rand"
	"minesweeper/ws"
	"net/http"
	"time"
)

//	func allowCors(w *http.ResponseWriter) {
//		(*w).Header().Set("Access-Control-Allow-Origin", "*")
//	}
func init() {
	rand.Seed(time.Now().UnixNano())
}
func main() {
	var l *ws.Lobby = ws.CreateLobby()
	http.Handle("/", http.FileServer(http.Dir("dist")))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		//	allowCors(&w)
		ws.ServeWs(w, r, l)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
