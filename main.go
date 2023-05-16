package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/websocket"
)

type Server struct {
	rooms map[string]map[*websocket.Conn]bool
}

type Room struct {
}

func NewServer() *Server {
	return &Server{
		rooms: make(map[string]map[*websocket.Conn]bool),
	}
}

func (s *Server) handleWebSocket(ws *websocket.Conn) {
	// based on URL, can close socket if not valid
	// client can check ws.readyState to see if they are actually connected
	urlParts := strings.Split(ws.Request().URL.Path, "/")

	if len(urlParts) != 3 || len(strings.Trim(urlParts[2], " ")) == 0 {
		fmt.Println(len(urlParts))
		fmt.Println("invalid URL:", ws.Request().URL.Path)
		ws.Close()
		return
	}

	roomId := urlParts[2]
	fmt.Println("new incoming connection from client:", ws.RemoteAddr())
	fmt.Println("room ID:", roomId)

	// if the room ID exists and there are already 2 people in it, don't allow to join
	if _, ok := s.rooms[roomId]; ok && len(s.rooms[roomId]) >= 2 {
		fmt.Printf("room %s is full\n", roomId)
		ws.Close()
		return
	} else if !ok {
		s.rooms[roomId] = make(map[*websocket.Conn]bool)
		fmt.Println("new room created:", roomId)
	}

	s.rooms[roomId][ws] = true // possibly use a mutex instead here
	s.readLoop(ws, roomId)
}

func (s *Server) readLoop(ws *websocket.Conn, roomId string) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				// remove connection once closed by client
				delete(s.rooms[roomId], ws)

				// if no one left in the room, delete the room
				if len(s.rooms[roomId]) == 0 {
					delete(s.rooms, roomId)
					fmt.Printf("room %s deleted", roomId)
				}

				break
			}

			fmt.Println("read error:", err)
			continue
		}

		msg := buf[:n]
		s.broadcast(msg, roomId)
	}
}

func (s *Server) broadcast(b []byte, roomId string) {
	for ws := range s.rooms[roomId] {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				fmt.Println("write error:", err)
			}
		}(ws)
	}
}

func main() {
	server := NewServer()
	http.Handle("/ws/", websocket.Handler(server.handleWebSocket))
	http.ListenAndServe(":8000", nil)
}
