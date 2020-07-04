package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/slarsson/game/network"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func start(s *Server, w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("UPGRADE ERROR:", err)
		return
	}

	client := &network.Client{
		Conn:          c,
		NetworkInput:  make(chan []byte, 100),
		NetworkOutput: make(chan []byte, 100),
	}

	s.Network.Register <- client
}

func main() {
	server := NewServer()
	go server.GameLoop()
	go server.Manager()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		start(server, w, r)
	})

	log.Fatal(http.ListenAndServe("0.0.0.0:1337", nil))
}
