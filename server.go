package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/websocket"
)

type Server struct {
	rooms   map[string]map[*websocket.Conn]bool
	players map[*websocket.Conn]*Player
}

type Event struct {
	EventType string `json:"eventType"`

	// event type determines which of these fields will be populated
	Payload struct {
		Name     string `json:"name"`
		Die      int    `json:"die"`
		Category int    `json:"category"`
	} `json:"payload"`
}

type Player struct {
	PlayerNum string // either p1 or p2
	Nickname  string
}

func NewServer() *Server {
	return &Server{
		rooms:   make(map[string]map[*websocket.Conn]bool),
		players: make(map[*websocket.Conn]*Player),
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

		// add person who created the room as player 1
		s.players[ws] = &Player{PlayerNum: "p1"}
	} else {
		// else add the player as player 2
		s.players[ws] = &Player{PlayerNum: "p2"}
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
		s.handleEvent(msg, ws)
		s.broadcast(msg, roomId)
	}
}

// determine what to do based on event type, this will interact with a yacht dice game instance
func (s *Server) handleEvent(msg []byte, fromConn *websocket.Conn) {
	var e Event
	// TODO handle errors when unmarshaling
	json.Unmarshal(msg, &e)

	switch e.EventType {
	case "name":
		// set the players name
		s.players[fromConn].Nickname = e.Payload.Name
		fmt.Println(s.players[fromConn])
	case "roll":
		fmt.Println("rolled")
	case "keep":
		fmt.Println(e.Payload.Die)
	case "score":
		fmt.Println(e.Payload.Category)
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
