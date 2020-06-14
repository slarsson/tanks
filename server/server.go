package main

import (
	"fmt"
	"time"
)

type Server struct {
	Network *Network
	Game    *Game
}

func newServer() *Server {
	return &Server{
		Network: newNetwork(),
		Game:    newGame(),
	}
}

func (s *Server) gameLoop() {
	var step float32
	step = 50
	ticker := time.NewTicker(time.Duration(step) * time.Millisecond)

	for range ticker.C {
		//fmt.Println("do world update")

		var buf string
		for _, p := range s.Game.Players {
			// fmt.Println(v.Name)
			p.Client.handleInputs(p, step)

			buf = `[{"x": ` + fmt.Sprintf("%f", p.Position.X) + `, "y": ` + fmt.Sprintf("%f", p.Position.Y) + `, "z": 0}]`
		}

		// x := rand.Float32() * 5
		// y := rand.Float32() * 3
		s.Network.broadcast <- []byte(buf)
	}
}

func (c *Client) handleInputs(p *Player, dt float32) {
	for len(c.NetworkInput) != 0 {
		message := <-c.NetworkInput

		//fmt.Println("ID: ", c.ID, " ->", message)

		if message[0] == byte(49) {
			fmt.Println("do somethings ffs")
			p.Position.Y += dt * 0.005
		}

		if message[1] == byte(49) {
			p.Position.X -= dt * 0.005
		}

		if message[2] == byte(49) {
			p.Position.Y -= dt * 0.005
		}

		if message[3] == byte(49) {
			p.Position.X += dt * 0.005
		}
	}
}

// package main

// import (
// 	"fmt"
// 	"time"

// 	"github.com/gorilla/websocket"
// )

// func (c *Client) swag() {
// 	fmt.Println("myclient")
// }

// type Server struct {

// 	// list of fkn clients
// 	clients map[*Client]bool

// 	// okej?
// 	register chan *Client

// 	unregister chan *Client

// 	// channel that recives messages from clients
// 	broadcast chan []byte
// }

// func newServer() *Server {
// 	return &Server{
// 		clients:    make(map[*Client]bool),
// 		register:   make(chan *Client),
// 		unregister: make(chan *Client),
// 		broadcast:  make(chan []byte),
// 	}
// }

// func (s *Server) reader(client *Client) {
// 	defer func() {
// 		s.unregister <- client // close connection?
// 	}()

// 	for {
// 		_, message, err := client.conn.ReadMessage()
// 		if err != nil {
// 			break
// 		}
// 		client.NetworkInput <- message
// 	}
// }

// func (s *Server) writer(client *Client) {
// 	// TODO: write from this connection osv
// 	// for {
// 	// 	select {
// 	// 	case message, ok := <-client.NetworkOutput:
// 	// 		client.conn.w
// 	// 	}
// 	// }
// }

// func (s *Server) run() {
// 	g := newGame()

// 	go s.gameLoop()
// 	for {
// 		select {
// 		case client := <-s.register:
// 			s.clients[client] = true
// 			g.addPlayer(client)
// 			go s.reader(client)
// 			go s.writer(client)
// 		case client := <-s.unregister:
// 			delete(s.clients, client)
// 		case message := <-s.broadcast:
// 			for client := range s.clients {
// 				client.swag()
// 			}
// 			fmt.Printf("message: %s", message)
// 		}
// 	}
// }

// func (s *Server) broadcastState() {
// 	buf := "["
// 	for x := range s.clients {
// 		buf += `{"x": ` + fmt.Sprintf("%f", x.position.X) + `, "y": ` + fmt.Sprintf("%f", x.position.Y) + `, "z": 0},`
// 	}
// 	if len(buf) > 1 {
// 		buf = buf[:len(buf)-1]
// 	}

// 	buf += "]"

// 	for client := range s.clients {
// 		// this will cause fkn problems?
// 		client.conn.WriteMessage(websocket.TextMessage, []byte(buf))
// 		//client.position.print()
// 	}
// }

// func (c *Client) inputs(dt float32) {
// 	for len(c.NetworkInput) != 0 {
// 		message := <-c.NetworkInput

// 		//fmt.Println("ID: ", c.ID, " ->", message)

// 		if message[0] == byte(49) {
// 			fmt.Println("do somethings ffs")
// 			c.position.Y += dt * 0.005
// 		}

// 		if message[1] == byte(49) {
// 			c.position.X -= dt * 0.005
// 		}

// 		if message[2] == byte(49) {
// 			c.position.Y -= dt * 0.005
// 		}

// 		if message[3] == byte(49) {
// 			c.position.X += dt * 0.005
// 		}
// 	}
// }

// func (s *Server) gameLoop() {
// 	fmt.Println("gameloop started..")
// 	var step float32
// 	step = 50
// 	ticker := time.NewTicker(time.Duration(step) * time.Millisecond)

// 	for range ticker.C {
// 		//fmt.Println("do update osv", len(s.clients))

// 		for c := range s.clients {
// 			//fmt.Println("get inputz")
// 			//fmt.Println("do shiet")
// 			c.inputs(step)
// 		}

// 		s.broadcastState()
// 	}
// }
