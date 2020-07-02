package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"time"

	"github.com/slarsson/game/game"
	"github.com/slarsson/game/network"
)

type Server struct {
	Network    *network.Network
	Game       *game.Game
	Operations chan *network.Action
}

func NewServer() *Server {
	n := network.NewNetwork()

	s := &Server{
		Network:    n,
		Game:       game.NewGame(n),
		Operations: make(chan *network.Action),
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
					s.Game.AddPlayer(message.Client)
				case 1:
					for k, v := range s.Game.Players {
						if v.Client == message.Client {
							fmt.Println("player to remove found!!!")
							s.Game.RemovePlayer(k)

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

	// c := &Client{
	// 	conn:          nil,
	// 	NetworkInput:  make(chan []byte, 100),
	// 	NetworkOutput: make(chan []byte, 100),
	// }
	// s.Game.addPlayer(c)

	for range ticker.C {
		start := time.Now()

		buf := make([]byte, 0, 30)
		mt := make([]byte, 4)
		binary.LittleEndian.PutUint32(mt, 0)
		buf = append(buf, mt...)
		for _, p := range s.Game.Players {
			if p.Client == nil {
				//p.moveBot(step)
			} else {
				s.handleInputs(p.Client, p, step)
			}
			//p.Client.handleInputs(p, step)

			p.AppendPlayerState(&buf)
		}

		s.Network.Broadcast <- buf

		fmt.Println("Executing time:", time.Since(start)*1000)
	}
}

func (s *Server) handleInputs(c *network.Client, p *game.Player, dt float32) {
	//func (c *Client) handleInputs(p *Player, dt float32) {
	for len(c.NetworkInput) != 0 {

		// swag
		message := <-c.NetworkInput

		// game.DecodeControls(message)
		p.Controls.Decode(&message)
		//p.Controls.Print()
		p.Move(&s.Game.Players, dt)
		//fmt.Println("handleInputs", message)

		switch message[0] {
		case 0:
			fmt.Println("add new fkn player here")

		case 1:
			if message[1] == 1 || message[3] == 1 {
				if message[1] == 1 {
					p.Velocity.X -= float32(math.Sin(float64(p.Rotation))) * 0.0001 * dt
					p.Velocity.Y += float32(math.Cos(float64(p.Rotation))) * 0.0001 * dt
				} else {
					p.Velocity.X += float32(math.Sin(float64(p.Rotation))) * 0.0001 * dt
					p.Velocity.Y -= float32(math.Cos(float64(p.Rotation))) * 0.0001 * dt
				}
			} else {
				p.Velocity.Y = 0
				p.Velocity.X = 0
			}

			if message[2] == 1 {
				p.Rotation += 0.002 * dt
			}

			if message[4] == 1 {
				p.Rotation -= 0.002 * dt
			}

			if message[5] == 1 {
				p.TurretRotation += 0.002 * dt
			}

			if message[6] == 1 {
				p.TurretRotation -= 0.002 * dt
			}

			if message[7] == 1 {
				p.Shoot = true

				fmt.Println("add new projectile for ID:", p.ID)
			} else {
				p.Shoot = false
			}

		default:
			fmt.Println("unknown command")
		}

		p.Position.X += dt * p.Velocity.X
		p.Position.Y += dt * p.Velocity.Y

		for _, v := range s.Game.Map.Obstacles {
			tank := game.NewTankHullPolygon()
			tank.Rotate(p.Rotation, p.Position)
			ok, mtv := tank.Collision(v)
			if ok {
				fmt.Println("CRASH WITH THE FKN WALL")
				dx := mtv.Vector.X * mtv.Magnitude
				dy := mtv.Vector.Y * mtv.Magnitude

				if (dx < 0 && p.Velocity.X < 0) || (dx > 0 && p.Velocity.X > 0) {
					dx = -dx
				}

				if (dy < 0 && p.Velocity.Y < 0) || (dy > 0 && p.Velocity.Y > 0) {
					dy = -dy
				}

				p.Position.X += dx
				p.Position.Y += dy

				p.Velocity.X = 0
				p.Velocity.Y = 0
			}
		}

		for i, v := range s.Game.Players {
			if i == p.ID {
				continue
			}

			poly1 := game.NewTankHullPolygon()
			poly1.Rotate(p.Rotation, p.Position)

			poly2 := game.NewTankHullPolygon()
			poly2.Rotate(v.Rotation, v.Position)

			ok, mtv := poly1.Collision(poly2)
			if !ok {
				continue
			}

			if message[2] == 1 || message[4] == 1 {
				fmt.Println("ROTATION WILL FUCK IT UP")
				if message[2] == 1 {
					p.Rotation -= 0.002 * dt
				}

				if message[4] == 1 {
					p.Rotation += 0.002 * dt
				}

				poly1 = game.NewTankHullPolygon()
				poly1.Rotate(p.Rotation, p.Position)

				poly2 = game.NewTankHullPolygon()
				poly2.Rotate(v.Rotation, v.Position)

				ok, mtv = poly1.Collision(poly2)
				if !ok {
					continue
				}
			}

			dx := mtv.Vector.X * mtv.Magnitude
			dy := mtv.Vector.Y * mtv.Magnitude

			if (dx < 0 && p.Velocity.X < 0) || (dx > 0 && p.Velocity.X > 0) {
				dx = -dx
			}

			if (dy < 0 && p.Velocity.Y < 0) || (dy > 0 && p.Velocity.Y > 0) {
				dy = -dy
			}

			p.Position.X += dx
			p.Position.Y += dy

			p.Velocity.X = 0
			p.Velocity.Y = 0
		}

	}
}
