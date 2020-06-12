package main

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn         *websocket.Conn
	position     *Point
	NetworkInput chan []byte
}

func (c *Client) swag() {
	fmt.Println("myclient")
}

type Server struct {

	// list of fkn clients
	clients map[*Client]bool

	// okej?
	register chan *Client

	unregister chan *Client

	// channel that recives messages from clients
	broadcast chan []byte
}

func newServer() *Server {
	return &Server{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
}

func (s *Server) run() {
	go s.gameLoop()
	for {
		select {
		case client := <-s.register:
			s.clients[client] = true
		case client := <-s.unregister:
			delete(s.clients, client)
		case message := <-s.broadcast:
			for client := range s.clients {
				client.swag()
			}
			fmt.Printf("message: %s", message)
		}
	}
}

func (s *Server) broadcastState() {
	buf := "["
	for x := range s.clients {
		buf += `{"x": ` + fmt.Sprintf("%f", x.position.X) + `, "y": 0, "z": 0},`
	}
	if len(buf) > 1 {
		buf = buf[:len(buf)-1]
	}

	buf += "]"

	for client := range s.clients {
		// this will cause fkn problems?
		client.conn.WriteMessage(websocket.TextMessage, []byte(buf))
		//client.position.print()
	}
}

func (c *Client) inputs() {
	for len(c.NetworkInput) != 0 {
		message := <-c.NetworkInput
		//x := uint32(message)
		//x := binary.BigEndian.Uint64(message)

		//fmt.Println(x)
		fmt.Printf("message: %s", message)
		c.position.move()
	}
}

func (s *Server) gameLoop() {
	ticker := time.NewTicker(100 * time.Millisecond)
	for range ticker.C {
		//fmt.Println("do update osv", len(s.clients))

		for c := range s.clients {
			//fmt.Println("get inputz")
			//fmt.Println("do shiet")
			c.inputs()
		}

		s.broadcastState()
	}
}
