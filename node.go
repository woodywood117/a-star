package main

import (
	"github.com/faiface/pixel"
	"math"
)

type Node struct {
	x, y, g, h float64
	last *Node
	traversable bool
}

func NewNode(x, y float64, traversable bool) *Node {
	n := &Node{
		x: x,
		y: y,
		g: -1,
		h: 0,
		last: nil,
		traversable: traversable,
	}
	return n
}

func (n *Node) Cost() float64 {
	if n == nil {
		return math.MaxFloat64
	}
	return n.g + n.h
}

func (n *Node) Hscore(dst_x, dst_y float64) {
	n.h = n.dist(dst_x, dst_y)
}

func (n *Node) dist(x, y float64) float64 {
	first := math.Pow(float64(x-n.x), 2)
	second := math.Pow(float64(y-n.y), 2)
	return math.Sqrt(first + second)
}

func (n *Node) Draw (win pixel.Target, tile *pixel.Sprite) {
	tile.Draw(win, pixel.IM.Moved(pixel.V(n.x * scale + scale/2, n.y * scale + scale/2)))
}