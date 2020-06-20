package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"time"
)

type Server struct {
	Network    *Network
	Game       *Game
	Operations chan *Actions
}

type Actions struct {
	MessageType int8
	Client      *Client
}

func newServer() *Server {
	n := newNetwork()

	s := &Server{
		Network:    n,
		Game:       newGame(n),
		Operations: make(chan *Actions),
	}

	go s.Network.start(s.Operations)
	return s
}

func (s *Server) ops() {
	for {
		select {
		case message, ok := <-s.Operations:
			if ok {
				fmt.Println("SERVER OPS:", message)

				switch message.MessageType {
				case 0:
					s.Game.addPlayer(message.Client)
				case 1:
					for k, v := range s.Game.Players {
						if v.Client == message.Client {
							fmt.Println("player to remove found!!!")
							s.Game.removePlayer(k)

							// TODO: send del message to all clients
						}
					}
				}

			}
		}
	}
}

func (s *Server) gameLoop() {
	var step float32
	step = 50
	ticker := time.NewTicker(time.Duration(step) * time.Millisecond)

	for range ticker.C {
		//fmt.Println("do world update")

		// bufx := make([]byte, 0, 100)
		// x := make([]byte, 4)
		// binary.LittleEndian.PutUint32(x, 12342)

		// bufx = append(bufx, x...)
		// fmt.Println("test:", x, bufx)

		// var buf string
		// for _, p := range s.Game.Players {
		// 	// fmt.Println(v.Name)
		// 	p.Client.handleInputs(p, step)

		// 	buf = `[{"x": ` + fmt.Sprintf("%f", p.Position.X) + `, "y": ` + fmt.Sprintf("%f", p.Position.Y) + `, "z": 0}]`
		// }

		buf := make([]byte, 0, 20)
		for _, p := range s.Game.Players {
			s.handleInputs(p.Client, p, step)
			//p.Client.handleInputs(p, step)

			id := make([]byte, 4)
			x := make([]byte, 4)
			y := make([]byte, 4)
			z := make([]byte, 4)

			vx := make([]byte, 4)
			vy := make([]byte, 4)
			vz := make([]byte, 4)

			binary.LittleEndian.PutUint32(id, math.Float32bits(float32(p.ID)))

			binary.LittleEndian.PutUint32(x, math.Float32bits(p.Position.X))
			binary.LittleEndian.PutUint32(y, math.Float32bits(p.Position.Y))
			binary.LittleEndian.PutUint32(z, math.Float32bits(p.Position.Z))

			binary.LittleEndian.PutUint32(vx, math.Float32bits(p.Velocity.X))
			binary.LittleEndian.PutUint32(vy, math.Float32bits(p.Velocity.Y))
			binary.LittleEndian.PutUint32(vz, math.Float32bits(p.Velocity.Z))

			buf = append(buf, id...)
			buf = append(buf, x...)
			buf = append(buf, y...)
			buf = append(buf, z...)

			buf = append(buf, vx...)
			buf = append(buf, vy...)
			buf = append(buf, vz...)
		}

		//fmt.Println(buf)

		s.Network.broadcast <- buf
	}
}

func (s *Server) handleInputs(c *Client, p *Player, dt float32) {
	//func (c *Client) handleInputs(p *Player, dt float32) {
	for len(c.NetworkInput) != 0 {
		message := <-c.NetworkInput

		//fmt.Println("handleInputs", message)

		switch message[0] {
		case 0:
			fmt.Println("add new fkn player here")

		case 1:
			//fmt.Println(message)

			x := 0.0001 * dt

			if message[1] == 1 || message[3] == 1 {
				if message[1] == 1 {
					p.Velocity.Y += x
				} else {
					p.Velocity.Y -= x
				}
			} else {
				p.Velocity.Y = 0
			}

			if message[2] == 1 || message[4] == 1 {
				if message[2] == 1 {
					p.Velocity.X -= x
				} else {
					p.Velocity.X += x
				}
			} else {
				p.Velocity.X = 0
			}

		default:
			fmt.Println("unknown command")
		}

		// if message[0] == 1 {
		// 	p.Velocity.Y += dt * 0.005
		// } else {
		// 	//p.Velocity.Y = 0
		// }

		// if message[1] == 1 {
		// 	p.Velocity.X -= dt * 0.005
		// } else {
		// 	p.Velocity.X = 0
		// }

		// if message[2] == 1 {
		// 	p.Velocity.Y -= dt * 0.005
		// } else {
		// 	p.Velocity.Y = 0
		// }

		// if message[3] == 1 {
		// 	p.Velocity.X += dt * 0.005
		// } else {
		// 	p.Velocity.X = 0
		// }

		// fmt.Println(p.Velocity)

		p.Position.X += dt * p.Velocity.X
		p.Position.Y += dt * p.Velocity.Y
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
