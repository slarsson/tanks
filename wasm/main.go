package main

import (
	"fmt"
	"syscall/js"
)

type keys struct {
	w     bool
	a     bool
	s     bool
	d     bool
	left  bool
	right bool
}

type Player struct {
	ID             int
	Position       *Vector3
	Velocity       *Vector3
	Rotation       float32
	TurretRotation float32
}

type Vector3 struct {
	X float32
	Y float32
	Z float32
}

func (v *Vector3) setPosition(x float32, y float32, z float32) {
	v.X = x
	v.Y = y
	v.Z = z
}

func (v *Vector3) getPosition() []float32 {
	return []float32{v.X, v.Y, v.Z}
}

var controls = keys{w: false, a: false, s: false, d: false}
var players = make(map[int]*Player)

var localPlayer = 8081

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
		if _, ok := players[args[i].Int()]; ok {

			players[args[i].Int()].Position.setPosition(float32(args[i+1].Float()), float32(args[i+2].Float()), float32(args[i+3].Float()))
			players[args[i].Int()].Rotation = float32(args[i+7].Float())
			players[args[i].Int()].TurretRotation = float32(args[i+8].Float())

			// predict next step ?
			dt := float32(50)
			p := players[args[i].Int()]
			p.Position.X += dt * p.Velocity.X
			p.Position.Y += dt * p.Velocity.Y

		} else {
			fmt.Println("add new player?")
			players[args[i].Int()] = &Player{
				ID: args[i].Int(),
				Position: &Vector3{
					X: float32(args[i+1].Float()),
					Y: float32(args[i+2].Float()),
					Z: float32(args[i+3].Float()),
				},
				Velocity: &Vector3{
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

func getp(this js.Value, args []js.Value) interface{} {

	play := players[args[0].Int()]
	// dt := float32(args[0].Float())
	// play.Position.X += dt * play.Velocity.Y
	// play.Position.Y += dt * play.Velocity.X
	p := play.Position.getPosition()
	return []interface{}{p[0], p[1], p[2], play.Rotation, play.TurretRotation}
}

func printState(this js.Value, args []js.Value) interface{} {
	fmt.Println(controls)
	return js.ValueOf(nil)
}

func test(this js.Value, x []js.Value) interface{} {
	//js.Global().Set("output", js.ValueOf(x[0]))
	fmt.Println("wtf")

	for _, v := range x {
		println(v.Int())
	}

	return js.ValueOf(nil)
}

func local(this js.Value, args []js.Value) interface{} {

	// dt := float32(args[0].Float())
	// p := players[8081]

	// if controls.w {
	// 	p.Position.Y += dt * 0.005
	// }

	// if controls.a {
	// 	p.Position.X -= dt * 0.005
	// }

	// if controls.s {
	// 	p.Position.Y -= dt * 0.005
	// }

	// if controls.d {
	// 	p.Position.X += dt * 0.005
	// }

	return js.ValueOf(nil)
}

func registerCallbacks() {
	js.Global().Set("poll", js.FuncOf(poll))

	js.Global().Set("swag", js.FuncOf(test))

	js.Global().Set("state", js.FuncOf(keypress))
	js.Global().Set("wasmprint", js.FuncOf(printState))

	js.Global().Set("wasmupdate", js.FuncOf(update))
	js.Global().Set("wasmgetpos", js.FuncOf(getp))
	js.Global().Set("wasmglocal", js.FuncOf(local))
}

func main() {
	c := make(chan struct{}, 0)

	registerCallbacks()

	fmt.Println("Hello, WebAssembly!")
	<-c
}
