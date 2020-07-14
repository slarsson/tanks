package main

import (
	"encoding/binary"
	"fmt"
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
		//start := time.Now()

		s.Game.PManager.UpdateAll(step, &s.Game.Players, s.Game.Map, s.Game.Network.Broadcast)
		//fmt.Println("antal:", len(s.Game.PManager.Projectiles))

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
		//s.Network.Broadcast <- *game.TestMessage()

		//fmt.Println("Executing time:", time.Since(start)*1000)
	}
}

func (s *Server) handleInputs(c *network.Client, p *game.Player, dt float32) {
	for len(c.NetworkInput) != 0 {
		message := <-c.NetworkInput

		switch message[0] {
		case 0:
			fmt.Println("do something else..")
		case 1:
			p.Controls.Decode(&message)
			p.SetSequenceNumber(&message)
			p.Move(dt)

			if s.Game.Map.OutOfBounds(p.Position) {
				fmt.Println("yeees")
				p.Position.Set(0, 0, 0)
			}

			p.HandleCollsionWithObjects(&s.Game.Map.Obstacles)
			p.HandleCollsionWithPlayers(&s.Game.Players, dt)
			// if p.Controls.Shoot {
			// 	s.Game.PManager.NewProjectile(p)
			// 	//s.Game.AddProjectile(p)
			// }

			if projectile, ok := p.Shoot(); ok {
				s.Game.PManager.Add(projectile)
			}

		default:
			fmt.Println("unknown command")
		}
	}
}
