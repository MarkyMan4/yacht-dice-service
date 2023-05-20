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
	games   map[string]*Game // map from room ID to game
}

type Event struct {
	EventType string `json:"eventType"`

	// event type determines which of these fields will be populated
	Payload struct {
		Name     string `json:"name"`
		Die      int    `json:"die"`
		Category string `json:"category"`
	} `json:"payload"`
}

type Player struct {
	PlayerNum string `json:"playerNum"` // either p1 or p2
	Nickname  string `json:"nickname"`
}

func NewServer() *Server {
	return &Server{
		rooms:   make(map[string]map[*websocket.Conn]bool),
		players: make(map[*websocket.Conn]*Player),
		games:   make(map[string]*Game),
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

		// create the game
		s.games[roomId] = NewGame()
		s.games[roomId].Player1 = s.players[ws]
	} else {
		// else add the player as player 2
		s.players[ws] = &Player{PlayerNum: "p2"}

		// add player 2 to the game
		s.games[roomId].Player2 = s.players[ws]
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
		s.handleEvent(msg, ws, roomId)
	}
}

// determine what to do based on event type, this will interact with a yacht dice game instance
func (s *Server) handleEvent(msg []byte, fromConn *websocket.Conn, roomId string) {
	var e Event
	// TODO handle errors when unmarshaling
	json.Unmarshal(msg, &e)

	// TODO handle event where player puts dice back into play
	switch e.EventType {
	case "name":
		// set the players name
		s.players[fromConn].Nickname = e.Payload.Name

		// if s.players[fromConn].PlayerNum == "p2" broadcast message with game state (i.e. start the game)
		if s.players[fromConn].PlayerNum == "p2" {
			s.broadcastGameToRoom(roomId)
		}
	case "roll":
		s.games[roomId].rollDice()
		s.broadcastGameToRoom(roomId)
	case "keep":
		s.games[roomId].keepDie(e.Payload.Die)
		s.broadcastGameToRoom(roomId)
	case "unkeep":
		s.games[roomId].unkeepDie(e.Payload.Die)
		s.broadcastGameToRoom(roomId)
	case "score":
		s.games[roomId].scoreRoll(e.Payload.Category)
		s.broadcastGameToRoom(roomId)
	}
}

func (s *Server) broadcastGameToRoom(roomId string) {
	gameData, err := json.Marshal(s.games[roomId])
	if err != nil {
		fmt.Println(err)
	}

	s.broadcast(gameData, roomId)
}

func (s *Server) broadcast(b []byte, roomId string) {
	// to include custom information for each player, I could marshal the game object,
	// unmarshal it to a map[string]interface{}, add custom keys then marshal it again
	for ws := range s.rooms[roomId] {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				fmt.Println("write error:", err)
			}
		}(ws)
	}
}
