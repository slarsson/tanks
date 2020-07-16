package main

import (
	"math/rand"
	"syscall/js"

	"github.com/slarsson/game/game"
)

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
		//fmt.Println("owner:", val.Owner.ID)
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
		val.CollisionTest(gameMap)
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
