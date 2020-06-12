package main

import "fmt"

type Point struct {
	X float32
	Y float32
	Z float32
}

func (p *Point) print() {
	fmt.Printf("Point: {x: %f, y: %f, z: %f} \n", p.X, p.Y, p.Z)
}

func (p *Point) move() {
	p.X += 0.1
}
