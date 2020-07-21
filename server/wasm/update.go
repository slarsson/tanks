package main

import (
	"fmt"
	"math/rand"
	"syscall/js"
	"time"

	"github.com/slarsson/game/game"
)

func updateLocalPlayer() {
	if v, ok := networkPlayers[localPlayer.ID]; ok {
		if !v.IsAlive {
			return
		}
	}

	// update localPlayer with current input
	// TODO: sync actions with server?
	localPlayer.Move(updateRate)
	localPlayer.HandleCollsionWithPlayers(&networkPlayers, updateRate)
	localPlayer.HandleCollsionWithObjects(&gameMap.Obstacles)
	if pr, ok := localPlayer.Shoot(); ok {
		var wtf int
		for {
			wtf = rand.Intn(10000) // fejk random?
			_, ok := projectiles[wtf]
			if !ok {
				break
			}
		}
		projectiles[wtf] = pr
	}

	// update last snapshoot of state for current sequence number
	if prev.ShouldUpdate {
		//fmt.Println(time.Now().Sub(prev.Timestamp))
		prev.Position.X = localPlayer.Position.X
		prev.Position.Y = localPlayer.Position.Y
		prev.SequenceNumber = sequence
		prev.ShouldUpdate = false
		prev.Timestamp = time.Now()
	}
}

func update(this js.Value, args []js.Value) interface{} {
	for i := 0; i < len(args)-1; i += 12 {
		key := args[i].Int()
		if p, ok := networkPlayers[key]; ok {

			status := (args[i+11].Int() != 0)
			if p.IsAlive != status {
				p.IsAlive = status
				fmt.Println("WASM: alive status has change (id:", p.ID, ")")
			}

			if !p.IsAlive {
				continue
			}

			p.Position.Set(float32(args[i+1].Float()), float32(args[i+2].Float()), float32(args[i+3].Float()))
			p.Rotation = float32(args[i+7].Float())
			p.TurretRotation = float32(args[i+8].Float())

			// check if localPlayer position has deviated to much from the server
			// reset the localPlayer to the current (should be the latest state) server position
			if key == localPlayer.ID {
				if prev.SequenceNumber == uint32(args[i+10].Int()) {
					if prev.Compare(p.Position) {
						fmt.Println("WASM: correct position with server position")
						localPlayer.SyncState(p)
					}
					prev.ShouldUpdate = true
				} else if (prev.SequenceNumber + 2) < (uint32(args[i+10].Int())) {
					fmt.Println("WASM: missed nummer / error :(")
					//localPlayer.SyncState(p)
					prev.ShouldUpdate = true
				}

				continue
			}

			// TODO: kanske inte så här...
			if args[i+9].Int() == 1 {
				addProjectile4Real(p)
			}
		} else {
			networkPlayer := game.NewLocalPlayer()
			networkPlayer.Position.Set(float32(args[i+1].Float()), float32(args[i+2].Float()), float32(args[i+3].Float()))
			networkPlayer.Position.Set(float32(args[i+4].Float()), float32(args[i+5].Float()), float32(args[i+6].Float()))
			networkPlayers[args[i].Int()] = networkPlayer
		}
	}

	return js.ValueOf(nil)
}
