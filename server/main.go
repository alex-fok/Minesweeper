package main

import (
	"log"
	"minesweeper/ws"
	"net/http"
)

// func allowCors(w *http.ResponseWriter) {
// 	(*w).Header().Set("Access-Control-Allow-Origin", "*")
// }

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		//	allowCors(&w)
		ws.ServeWs(w, r)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
