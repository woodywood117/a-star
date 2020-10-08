package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
	"math"
)

type Node struct {
	x, y, g, h            float64
	last                  *Node
	left, right, down, up bool
}

func NewNode(x, y float64) *Node {
	n := &Node{
		x:     x,
		y:     y,
		g:     -1,
		h:     0,
		last:  nil,
		left:  true,
		right: true,
		down:  true,
		up:    true,
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

func (n *Node) Draw(win pixel.Target, tile *pixel.Sprite) {
	tile.Draw(win, pixel.IM.Moved(pixel.V(n.x*scale+scale/2, n.y*scale+scale/2)))

	imd := imdraw.New(nil)
	thickness := scale / 10

	if n.left {
		imd.Color = colornames.Black
		imd.Push(pixel.V(n.x*scale, n.y*scale+scale), pixel.V(n.x*scale, n.y*scale))
		imd.Line(thickness)
		imd.Draw(win)
		imd.Reset()
	}
	if n.right {
		imd.Color = colornames.Black
		imd.Push(pixel.V(n.x*scale+scale, n.y*scale+scale), pixel.V(n.x*scale+scale, n.y*scale))
		imd.Line(thickness)
		imd.Draw(win)
		imd.Reset()
	}
	if n.down {
		imd.Color = colornames.Black
		imd.Push(pixel.V(n.x*scale, n.y*scale), pixel.V(n.x*scale+scale, n.y*scale))
		imd.Line(thickness)
		imd.Draw(win)
		imd.Reset()
	}
	if n.up {
		imd.Color = colornames.Black
		imd.Push(pixel.V(n.x*scale, n.y*scale+scale), pixel.V(n.x*scale+scale, n.y*scale+scale))
		imd.Line(thickness)
		imd.Draw(win)
		imd.Reset()
	}
}
