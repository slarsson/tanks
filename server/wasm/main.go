package main

import (
	"fmt"
	"math"
	"syscall/js"

	"github.com/slarsson/game/game"
)

const UpdateRate = 50

type keys struct {
	w        bool
	a        bool
	s        bool
	d        bool
	left     bool
	right    bool
	spacebar bool
}

var controls = keys{
	w:        false,
	a:        false,
	s:        false,
	d:        false,
	left:     false,
	right:    false,
	spacebar: false,
}

var self = -1
var players = make(map[int]*game.Player)

func setSelf(this js.Value, args []js.Value) interface{} {
	self = args[0].Int()
	fmt.Println("WASM: self has been set to", self)
	return js.ValueOf(nil)
}

func keypress(this js.Value, args []js.Value) interface{} {
	key := args[0].String()

	if key == "w" {
		controls.w = args[1].Bool()
		return js.ValueOf(true)
	}

	if key == "a" {
		controls.a = args[1].Bool()
		return js.ValueOf(true)
	}

	if key == "s" {
		controls.s = args[1].Bool()
		return js.ValueOf(true)
	}

	if key == "d" {
		controls.d = args[1].Bool()
		return js.ValueOf(true)
	}

	if key == "ArrowLeft" {
		controls.left = args[1].Bool()
		return js.ValueOf(true)
	}

	if key == "ArrowRight" {
		controls.right = args[1].Bool()
		return js.ValueOf(true)
	}

	return js.ValueOf(false)
}

func poll(this js.Value, args []js.Value) interface{} {
	buf := make([]byte, 7)
	buf[0] = 1

	if controls.w {
		buf[1] = 1
	}

	if controls.a {
		buf[2] = 1
	}

	if controls.s {
		buf[3] = 1
	}

	if controls.d {
		buf[4] = 1
	}

	if controls.left {
		buf[5] = 1
	}

	if controls.right {
		buf[6] = 1
	}

	uint8Array := js.Global().Get("Uint8Array")
	dst := uint8Array.New(len(buf))
	js.CopyBytesToJS(dst, buf)
	return dst
}

func update(this js.Value, args []js.Value) interface{} {

	//fmt.Println(args)

	for i := 0; i < len(args); i += 9 {
		key := args[i].Int()
		if p, ok := players[key]; ok {
			if key == self {
				//continue
			}

			p.Position.Set(float32(args[i+1].Float()), float32(args[i+2].Float()), float32(args[i+3].Float()))
			p.Rotation = float32(args[i+7].Float())
			p.TurretRotation = float32(args[i+8].Float())

			if key != self {
				p.Position.X += UpdateRate * p.Velocity.X
				p.Position.Y += UpdateRate * p.Velocity.Y

				for _, v := range players {
					if v.ID == p.ID || v.ID == self {
						continue
					}

					poly1 := game.NewTankHullPolygon()
					poly1.Rotate(p.Rotation, p.Position)

					poly2 := game.NewTankHullPolygon()
					poly2.Rotate(v.Rotation, v.Position)

					ok, mtv := poly1.Collision(poly2)
					if ok {
						fmt.Println("CLIENT HAS CRASHED..")

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
					}
				}
			}

			// players[args[i].Int()].Position.Set(float32(args[i+1].Float()), float32(args[i+2].Float()), float32(args[i+3].Float()))
			// players[args[i].Int()].Rotation = float32(args[i+7].Float())
			// players[args[i].Int()].TurretRotation = float32(args[i+8].Float())

			// if key != self {
			// 	players[args[i].Int()].Position.X += 150 * players[args[i].Int()].Velocity.X
			// 	players[args[i].Int()].Position.Y += 150 * players[args[i].Int()].Velocity.Y
			// 	//continue
			// }

			// // predict next step ?
			// // dt := float32(50)
			// // p := players[args[i].Int()]
			// // p.Position.X += dt * p.Velocity.X
			// // p.Position.Y += dt * p.Velocity.Y

		} else {
			fmt.Println("add new player?")
			players[args[i].Int()] = &game.Player{
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
			}
		}
	}

	return js.ValueOf(nil)
}

func getPosition(this js.Value, args []js.Value) interface{} {
	if p, ok := players[args[0].Int()]; ok {
		return []interface{}{p.Position.X, p.Position.Y, p.Position.Z, p.Rotation, p.TurretRotation}
	}
	return js.ValueOf(nil)
}

func guessPosition(this js.Value, args []js.Value) interface{} {
	dt := float32(args[0].Float())
	for _, p := range players {
		if p.ID == self {
			continue
		}
		p.Position.X += dt * p.Velocity.X
		p.Position.Y += dt * p.Velocity.Y
	}

	return js.ValueOf(nil)
}

func local(this js.Value, args []js.Value) interface{} {
	//fmt.Println("varför är min dator fkn sämst?")

	p, ok := players[args[0].Int()]
	if !ok {
		return js.ValueOf(nil)
	}

	dt := float32(args[1].Float())

	// fmt.Println(p)
	// fmt.Println(dt)

	if controls.w || controls.s {
		fmt.Println("do update shiett..")
		if controls.w {
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

	if controls.a {
		p.Rotation += 0.002 * dt
	}

	if controls.d {
		p.Rotation -= 0.002 * dt
	}

	if controls.left {
		p.TurretRotation += 0.002 * dt
	}

	if controls.right {
		p.TurretRotation -= 0.002 * dt
	}

	p.Position.X += dt * p.Velocity.X
	p.Position.Y += dt * p.Velocity.Y

	return js.ValueOf(nil)
}

func main() {
	js.Global().Set("wasm__poll", js.FuncOf(poll))
	js.Global().Set("wasm__keypress", js.FuncOf(keypress))
	js.Global().Set("wasm__update", js.FuncOf(update))
	js.Global().Set("wasm__get", js.FuncOf(getPosition))
	js.Global().Set("wasm__local", js.FuncOf(local))
	js.Global().Set("wasm__setSelf", js.FuncOf(setSelf))
	js.Global().Set("wasm__guessPosition", js.FuncOf(guessPosition))

	fmt.Println("WebAssembly init!")

	c := make(chan struct{}, 0)
	<-c
}
