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
		// xx := make([]byte, 0, 30)
		// x1 := make([]byte, 4)
		// x2 := make([]byte, 4)
		// x3 := make([]byte, 4)
		// x4 := make([]byte, 4)
		// x5 := make([]byte, 4)
		// x6 := make([]byte, 4)
		// x7 := make([]byte, 4)
		// x8 := make([]byte, 4)
		// binary.LittleEndian.PutUint32(x1, 1)
		// binary.LittleEndian.PutUint32(x2, 1)
		// binary.LittleEndian.PutUint32(x3, 1)
		// xx = append(xx, x1...)
		// xx = append(xx, x2...)
		// xx = append(xx, x3...)
		// xx = append(xx, x4...)
		// xx = append(xx, x5...)
		// xx = append(xx, x6...)
		// xx = append(xx, x7...)
		// xx = append(xx, x8...)
		// c.NetworkInput <- xx

		buf := make([]byte, 0, 30)
		for _, p := range s.Game.Players {
			if p.Client == nil {
				//p.moveBot(step)
			} else {
				s.handleInputs(p.Client, p, step)
			}
			//p.Client.handleInputs(p, step)

			id := make([]byte, 4)
			x := make([]byte, 4)
			y := make([]byte, 4)
			z := make([]byte, 4)

			vx := make([]byte, 4)
			vy := make([]byte, 4)
			vz := make([]byte, 4)

			rot := make([]byte, 4)
			trot := make([]byte, 4)

			binary.LittleEndian.PutUint32(id, math.Float32bits(float32(p.ID)))

			binary.LittleEndian.PutUint32(x, math.Float32bits(p.Position.X))
			binary.LittleEndian.PutUint32(y, math.Float32bits(p.Position.Y))
			binary.LittleEndian.PutUint32(z, math.Float32bits(p.Position.Z))

			binary.LittleEndian.PutUint32(vx, math.Float32bits(p.Velocity.X))
			binary.LittleEndian.PutUint32(vy, math.Float32bits(p.Velocity.Y))
			binary.LittleEndian.PutUint32(vz, math.Float32bits(p.Velocity.Z))

			binary.LittleEndian.PutUint32(rot, math.Float32bits(p.Rotation))
			binary.LittleEndian.PutUint32(trot, math.Float32bits(p.TurretRotation))

			buf = append(buf, id...)
			buf = append(buf, x...)
			buf = append(buf, y...)
			buf = append(buf, z...)

			buf = append(buf, vx...)
			buf = append(buf, vy...)
			buf = append(buf, vz...)

			buf = append(buf, rot...)
			buf = append(buf, trot...)
		}

		//fmt.Println(buf)

		s.Network.Broadcast <- buf

		tt := time.Since(start) * 1000 //* time.Millisecond
		fmt.Println("Executing time:", tt)
	}
}

func (s *Server) handleInputs(c *network.Client, p *game.Player, dt float32) {
	//func (c *Client) handleInputs(p *Player, dt float32) {
	for len(c.NetworkInput) != 0 {
		message := <-c.NetworkInput

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
