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
	Operations chan *Action
}

type Action struct {
	MessageType int8
	Client      *Client
}

func NewServer() *Server {
	n := NewNetwork()

	s := &Server{
		Network:    n,
		Game:       newGame(n),
		Operations: make(chan *Action),
	}

	go s.Network.Start(s.Operations)
	return s
}

func (s *Server) Manager() {
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

func (s *Server) GameLoop() {
	//var step float32
	step := float32(50)
	ticker := time.NewTicker(time.Duration(step) * time.Millisecond)

	for range ticker.C {

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

		p.Position.X += dt * p.Velocity.X
		p.Position.Y += dt * p.Velocity.Y
	}
}
