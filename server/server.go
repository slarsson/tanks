package main

import (
	"encoding/binary"
	"fmt"
	"regexp"
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

func validateName(input []byte) bool {
	if len(input) > 20 {
		return false
	}

	matched, err := regexp.Match(`^[A-Za-z0-9_-]+$`, input)

	if err != nil || !matched {
		return false
	}
	return true
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
							s.Game.RemovePlayer(k)
						}
					}
				case 99:
					for _, v := range s.Game.Players {
						if v.Client == message.Client {
							if validateName(message.Payload) {
								name := string(message.Payload)
								if s.Game.SetPlayerName(v.ID, name) {
									v.ExitLobby()
									s.Network.Broadcast <- game.PlayerNameMessage(v.ID, name)
								} else {
									message.Client.NetworkOutput <- game.ErrorMessage()
								}
							} else {
								message.Client.NetworkOutput <- game.ErrorMessage()
							}
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
			p.SetSequenceNumber(&message)

			if !p.IsAlive {
				if p.Lobby {
					return
				}

				if p.RespawnTime > game.RespawnTime {
					p.Respawn(s.Game.Map.RandomSpawn())
				} else {
					p.RespawnTime += dt
				}
				return
			}

			p.Controls.Decode(&message)
			p.Move(dt)
			p.HandleCollsionWithObjects(&s.Game.Map.Obstacles, dt)
			p.HandleCollsionWithPlayers(&s.Game.Players, dt)

			if s.Game.Map.OutOfBounds(p.Position) {
				p.OutOfMapTime += dt
				if p.OutOfMapTime > 3000 {
					p.Kill()
				}
			} else {
				p.OutOfMapTime = 0
				if projectile, ok := p.Shoot(); ok {
					s.Game.PManager.Add(projectile)
				}
			}
		default:
			fmt.Println("unknown command")
		}
	}
}

// func (s *Server) Swag() {
// 	ticker := time.NewTicker(time.Duration(523) * time.Millisecond)

// 	for range ticker.C {
// 		//fmt.Println(game.SendPlayerName())

// 		s.Network.Broadcast <- *game.SendPlayerName()
// 	}
// }
