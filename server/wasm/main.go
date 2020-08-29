package main

import (
	"encoding/binary"
	"fmt"
	"syscall/js"
	"time"

	"github.com/slarsson/tanks/game"
)

// updateRate: is the rate (time in ms) which updates are sent to the server
// this value should match the server tickrate (the rate updates are sent from the server to the client)
const updateRate = 50

// prev: is a snapshoot of previously state
// used to detect if server and client game state has deviated from each other
var prev = game.LastState{
	Position:       game.Vector3{X: 0, Y: 0, Z: 0},
	SequenceNumber: 0,
	ShouldUpdate:   true,
	Timestamp:      time.Now(),
}

// localPlayer: Player object for local state
var localPlayer = game.NewLocalPlayer()

// networkPlayers: is a map contaning all players in the game (including localPlayer as seen by server)
var networkPlayers = make(map[int]*game.Player)

// gameMap: the game map
var gameMap = game.NewMap()

// ...
var pmanager = game.NewProjectileManager()

// sequence: is the sequence number, incremented and sent to the server on every tick
var sequence = uint32(0)

// setSelf: sets the ID of the localPlayer.
//
// args[0]: Int  # ID of localPlayer
func setSelf(_ js.Value, args []js.Value) interface{} {
	localPlayer.ID = args[0].Int()
	fmt.Println("WASM: local player has ID =", localPlayer.ID)
	return js.ValueOf(nil)
}

// keypress: maps keypresses to actions.
// Returns TRUE if key exists otherwise FALSE.
//
// args[0]: String  # the JavaScript KeyboardEvent "key" value (w, a, ArrowLeft, ..)
// args[1]: Bool    # true or false, if the key is pressed or released
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

// poll: gets the current state (the message sent to the server).
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
		buf[7] = 1
	}

	// increment sequence number and append to message
	sequence++
	s := make([]byte, 4)
	binary.LittleEndian.PutUint32(s, sequence)
	buf = append(buf, s...)

	updateLocalPlayer()

	// append to js array
	uint8Array := js.Global().Get("Uint8Array")
	dst := uint8Array.New(len(buf))
	js.CopyBytesToJS(dst, buf)
	return dst
}

func getPosition(this js.Value, args []js.Value) interface{} {
	id := args[0].Int()

	if id == localPlayer.ID {
		return []interface{}{
			localPlayer.Position.X,
			localPlayer.Position.Y,
			localPlayer.Position.Z,
			localPlayer.Rotation,
			localPlayer.TurretRotation,
		}
	}

	if p, ok := networkPlayers[id]; ok {
		return []interface{}{
			p.Position.X,
			p.Position.Y,
			p.Position.Z,
			p.Rotation,
			p.TurretRotation,
		}
	}

	return js.ValueOf(nil)
}

func removePlayer(this js.Value, args []js.Value) interface{} {
	key := args[0].Int()
	fmt.Println("WASM: remove id:", key)
	delete(networkPlayers, key)
	return js.ValueOf(nil)
}

// main exports functions to js
func main() {
	js.Global().Set("wasm__poll", js.FuncOf(poll))
	js.Global().Set("wasm__keypress", js.FuncOf(keypress))
	js.Global().Set("wasm__update", js.FuncOf(update))
	js.Global().Set("wasm__get", js.FuncOf(getPosition))
	js.Global().Set("wasm__setSelf", js.FuncOf(setSelf))
	js.Global().Set("wasm__updateProjectiles", js.FuncOf(updateProjectiles))
	js.Global().Set("wasm__removePlayer", js.FuncOf(removePlayer))
	js.Global().Set("wasm__updateCrane", js.FuncOf(updateCrane))

	fmt.Println("WASM: WebAssembly init!")

	c := make(chan struct{}, 0)
	<-c
}
