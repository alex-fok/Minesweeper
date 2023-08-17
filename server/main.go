package main

import (
	"log"
	"minesweeper/ws"
	"net/http"
)

//	func allowCors(w *http.ResponseWriter) {
//		(*w).Header().Set("Access-Control-Allow-Origin", "*")
//	}

func main() {
	var l *ws.Lobby = ws.CreateLobby()
	http.Handle("/", http.FileServer(http.Dir("dist")))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		//	allowCors(&w)
		ws.ServeWs(w, r, l)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
