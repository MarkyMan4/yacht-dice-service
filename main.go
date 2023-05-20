package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/websocket"
)

var logfile *os.File

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>test</h1>")
}

func main() {
	logfile, _ = os.OpenFile("/var/log/yacht_dice/log", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	defer logfile.Close()

	log.SetOutput(logfile)
	log.Println("listening on port 8000")

	server := NewServer()
	http.HandleFunc("/", handler)
	http.Handle("/ws/", websocket.Handler(server.handleWebSocket))
	log.Fatalln(http.ListenAndServe(":8000", nil))
}
