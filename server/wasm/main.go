package main

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"syscall/js"

	"github.com/slarsson/game/game"
)

const updateRate = 50

var prev = game.LastState{
	Position:       game.Vector3{X: 0, Y: 0, Z: 0},
	SequenceNumber: 0,
	ShouldUpdate:   true,
}

var localPlayer = game.NewLocalPlayer()

var networkPlayers = make(map[int]*game.Player)
var gameMap = game.NewMap()
var projectiles = make(map[int]*game.Projectile)
var sequence = uint32(0)

func setSelf(this js.Value, args []js.Value) interface{} {
	localPlayer.ID = args[0].Int()
	fmt.Println("WASM: local player has ID =", localPlayer.ID)
	return js.ValueOf(nil)
}

func keypress(this js.Value, args []js.Value) interface{} {
	key := args[0].String()
	status := args[1].Bool()

	if key == "w" {
		localPlayer.Controls.Forward = status
		return js.ValueOf(true)
	}

	if key == "a" {
		localPlayer.Controls.RotateLeft = status
		return js.ValueOf(true)
	}

	if key == "s" {
		localPlayer.Controls.Backward = status
		return js.ValueOf(true)
	}

	if key == "d" {
		localPlayer.Controls.RotateRight = status
		return js.ValueOf(true)
	}

	if key == "ArrowLeft" {
		localPlayer.Controls.RotateTurretLeft = status
		return js.ValueOf(true)
	}

	if key == "ArrowRight" {
		localPlayer.Controls.RotateTurretRight = status
		return js.ValueOf(true)
	}

	if key == " " {
		localPlayer.Controls.Shoot = status
		return js.ValueOf(true)
	}

	return js.ValueOf(false)
}

func poll(this js.Value, args []js.Value) interface{} {
	buf := make([]byte, 8)

	// set MessageType = 1
	buf[0] = 1

	if localPlayer.Controls.Forward {
		buf[1] = 1
	}

	if localPlayer.Controls.RotateLeft {
		buf[2] = 1
	}

	if localPlayer.Controls.Backward {
		buf[3] = 1
	}

	if localPlayer.Controls.RotateRight {
		buf[4] = 1
	}

	if localPlayer.Controls.RotateTurretLeft {
		buf[5] = 1
	}

	if localPlayer.Controls.RotateTurretRight {
		buf[6] = 1
	}

	if localPlayer.Controls.Shoot {
		addProjectile4Real(localPlayer)
		buf[7] = 1
	}

	// increment sequence number and append to message
	sequence++
	s := make([]byte, 4)
	binary.LittleEndian.PutUint32(s, sequence)
	buf = append(buf, s...)

	// update localPlayer with current input
	// TODO: sync actions with server?
	localPlayer.Move(updateRate)
	localPlayer.HandleCollsionWithPlayers(&networkPlayers, updateRate)
	localPlayer.HandleCollsionWithObjects(&gameMap.Obstacles)

	// update last snapshoot of state for current sequence number
	if prev.ShouldUpdate {
		prev.Position.X = localPlayer.Position.X
		prev.Position.Y = localPlayer.Position.Y
		prev.SequenceNumber = sequence
		prev.ShouldUpdate = false
	}

	// append to js array
	uint8Array := js.Global().Get("Uint8Array")
	dst := uint8Array.New(len(buf))
	js.CopyBytesToJS(dst, buf)
	return dst
}

func update(this js.Value, args []js.Value) interface{} {
	for i := 0; i < len(args)-1; i += 11 {
		key := args[i].Int()
		if p, ok := networkPlayers[key]; ok {

			// check if localPlayer position has deviated to much from the server
			// reset the localPlayer to the current (should be the latest state) server position
			if key == localPlayer.ID {
				if prev.SequenceNumber == uint32(args[i+10].Int()) {
					x := float32(args[i+1].Float())
					y := float32(args[i+2].Float())
					z := float32(args[i+3].Float())

					if prev.Compare(x, y, z) {
						localPlayer.Position.Set(x, y, z)
						localPlayer.Rotation = float32(args[i+7].Float())
						localPlayer.TurretRotation = float32(args[i+8].Float())
					}

					// TODO: this should not have to be done?
					p.Position.Set(x, y, z)
					p.Rotation = float32(args[i+7].Float())
					p.TurretRotation = float32(args[i+8].Float())

					prev.ShouldUpdate = true
				} else if (prev.SequenceNumber + 2) < (uint32(args[i+10].Int())) {
					fmt.Println("missed nummer / error :(")
					prev.ShouldUpdate = true
				}

				continue
			}

			// update networkPlayer position / rotation
			p.Position.Set(float32(args[i+1].Float()), float32(args[i+2].Float()), float32(args[i+3].Float()))
			p.Rotation = float32(args[i+7].Float())
			p.TurretRotation = float32(args[i+8].Float())

			if args[i+9].Int() == 1 {
				//fmt.Println("add to:", p)
				addProjectile4Real(p)
			}

			// TEST: guess next position, should be done the other way around? interpolate from old to current?
			// https://developer.valvesoftware.com/wiki/Source_Multiplayer_Networking
			// if key != self {
			// 	p.Position.X += updateRate * p.Velocity.X
			// 	p.Position.Y += updateRate * p.Velocity.Y

			// 	for _, v := range networkPlayers {
			// 		if v.ID == p.ID || v.ID == self {
			// 			continue
			// 		}

			// 		poly1 := game.NewTankHullPolygon()
			// 		poly1.Rotate(p.Rotation, p.Position)

			// 		poly2 := game.NewTankHullPolygon()
			// 		poly2.Rotate(v.Rotation, v.Position)

			// 		ok, mtv := poly1.Collision(poly2)
			// 		if ok {
			// 			fmt.Println("CLIENT HAS CRASHED..")

			// 			dx := mtv.Vector.X * mtv.Magnitude
			// 			dy := mtv.Vector.Y * mtv.Magnitude

			// 			if (dx < 0 && p.Velocity.X < 0) || (dx > 0 && p.Velocity.X > 0) {
			// 				dx = -dx
			// 			}

			// 			if (dy < 0 && p.Velocity.Y < 0) || (dy > 0 && p.Velocity.Y > 0) {
			// 				dy = -dy
			// 			}

			// 			p.Position.X += dx
			// 			p.Position.Y += dy
			// 		}
			// 	}
			// }

		} else {
			fmt.Println("add new player?")
			networkPlayers[args[i].Int()] = &game.Player{
				ID: args[i].Int(),
				Position: &game.Vector3{
					X: float32(args[i+1].Float()),
					Y: float32(args[i+2].Float()),
					Z: float32(args[i+3].Float()),
				},
				Velocity: &game.Vector3{
					X: float32(0),
					Y: float32(0),
					Z: float32(0),
				},
				Rotation:       0,
				TurretRotation: 0,
				Controls:       game.NewControls(),
			}
		}
	}

	return js.ValueOf(nil)
}

func getPosition(this js.Value, args []js.Value) interface{} {
	id := args[0].Int()

	if id == localPlayer.ID {
		return []interface{}{localPlayer.Position.X, localPlayer.Position.Y, localPlayer.Position.Z, localPlayer.Rotation, localPlayer.TurretRotation}
	}

	if p, ok := networkPlayers[id]; ok {
		return []interface{}{p.Position.X, p.Position.Y, p.Position.Z, p.Rotation, p.TurretRotation}
	}

	return js.ValueOf(nil)
}

func guessPosition(this js.Value, args []js.Value) interface{} {
	dt := float32(args[0].Float())
	for _, p := range networkPlayers {
		if p.ID == localPlayer.ID {
			continue
		}
		p.Position.X += dt * p.Velocity.X
		p.Position.Y += dt * p.Velocity.Y
	}

	return js.ValueOf(nil)
}

// Projectiles ...
//
// TODO: testing only..
func addProjectile(this js.Value, args []js.Value) interface{} {
	// var wtf int
	// for {
	// 	wtf = rand.Intn(10000) // fejk random?
	// 	_, ok := projectiles[wtf]
	// 	if !ok {
	// 		break
	// 	}
	// }

	// id := args[0].Int()

	// if id == localPlayer.ID {
	// 	projectiles[wtf] = localPlayer.NewProjectile()
	// } else {
	// 	//projectiles[wtf] = networkPlayers[id].NewProjectile()
	// }

	return js.ValueOf(nil)
}

func addProjectile4Real(p *game.Player) {
	var wtf int
	for {
		wtf = rand.Intn(10000) // fejk random?
		_, ok := projectiles[wtf]
		if !ok {
			break
		}
	}
	projectiles[wtf] = p.NewProjectile()
}

func updateProjectiles(this js.Value, args []js.Value) interface{} {
	buf := make([]interface{}, len(projectiles)*5)
	//buf := make([]float32, len(projectiles)*3)

	dt := float32(args[0].Float())

	wtf := 0
	for i, val := range projectiles {
		fmt.Println("owner:", val.Owner.ID)
		if !val.IsAlive {
			delete(projectiles, i)
			buf[wtf] = i
			wtf++
			buf[wtf] = 0
			wtf++
			buf[wtf] = 0
			wtf++
			buf[wtf] = 0
			wtf++
			buf[wtf] = 0
			wtf++
			continue
		}

		x := val.Position.X
		y := val.Position.Y
		z := val.Position.Z

		val.Move(dt)
		val.CollisionTest()
		val.CollisionTestPlayers(&networkPlayers)
		buf[wtf] = i
		wtf++
		buf[wtf] = x
		wtf++
		buf[wtf] = y
		wtf++
		buf[wtf] = z
		wtf++
		buf[wtf] = 1
		wtf++

		// val.Move(dt)
		// val.CollisionTest()
		// buf[wtf] = i
		// wtf++
		// buf[wtf] = val.Position.X
		// wtf++
		// buf[wtf] = val.Position.Y
		// wtf++
		// buf[wtf] = val.Position.Z
		// wtf++

	}

	return buf
}

// main exports functions to js
func main() {
	js.Global().Set("wasm__poll", js.FuncOf(poll))
	js.Global().Set("wasm__keypress", js.FuncOf(keypress))
	js.Global().Set("wasm__update", js.FuncOf(update))
	js.Global().Set("wasm__get", js.FuncOf(getPosition))
	js.Global().Set("wasm__setSelf", js.FuncOf(setSelf))
	js.Global().Set("wasm__guessPosition", js.FuncOf(guessPosition))
	js.Global().Set("wasm__updateProjectiles", js.FuncOf(updateProjectiles))
	js.Global().Set("wasm__addProjectile", js.FuncOf(addProjectile))

	fmt.Println("WebAssembly init!")

	c := make(chan struct{}, 0)
	<-c
}
