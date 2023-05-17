package main

import (
	"net/http"

	"golang.org/x/net/websocket"
)

func main() {
	server := NewServer()
	http.Handle("/ws/", websocket.Handler(server.handleWebSocket))
	http.ListenAndServe(":8000", nil)
}
