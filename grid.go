package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"math/rand"
)

var pict *pixel.PictureData
var red, green, blue, black *pixel.Sprite

func init() {
	pict = pixel.MakePictureData(pixel.R(0, 0, scale, scale*4))
	for i := 0; i < len(pict.Pix); i++ {
		switch i / int(scale*scale) {
		case 0:
			pict.Pix[i] = colornames.Green
		case 1:
			pict.Pix[i] = colornames.Blue
		case 2:
			pict.Pix[i] = colornames.Red
		case 3:
			pict.Pix[i] = colornames.Black
		}
	}
	green = pixel.NewSprite(pict, pixel.R(0, 0, scale, scale))
	blue = pixel.NewSprite(pict, pixel.R(0, scale, scale, scale*2))
	red = pixel.NewSprite(pict, pixel.R(0, scale*2, scale, scale*3))
	black = pixel.NewSprite(pict, pixel.R(0, scale*3, scale, scale*4))
}

type Grid struct {
	width, height, src_x, src_y, dst_x, dst_y int
	opened, closed                            map[int]map[int]*Node
	current                                   *Node
	grid                                      [][]*Node
	batch                                     *pixel.Batch
}

func NewGrid(width, height, src_x, src_y, dst_x, dst_y int) *Grid {
	g := &Grid{
		width:  width,
		height: height,
		src_x:  src_x,
		src_y:  src_y,
		dst_x:  dst_x,
		dst_y:  dst_y,
		batch:  pixel.NewBatch(&pixel.TrianglesData{}, pict),
		opened: make(map[int]map[int]*Node),
		closed: make(map[int]map[int]*Node),
	}

	for x := 0; x < width; x++ {
		g.grid = append(g.grid, make([]*Node, height))
		for y := 0; y < height; y++ {
			chance := rand.Float64()
			traversable := false
			if chance > 0.4 || (x == src_x && y == src_y) || (x == dst_x && y == dst_y) {
				traversable = true
			}
			g.grid[x][y] = NewNode(float64(x), float64(y), traversable)
			g.grid[x][y].Hscore(float64(dst_x), float64(dst_y))
		}
		g.opened[x] = make(map[int]*Node)
		g.closed[x] = make(map[int]*Node)
	}

	g.opened[src_x][src_y] = g.grid[src_x][src_y]
	return g
}

func (g *Grid) Restart() {
	for x := 0; x < g.width; x++ {
		for y := 0; y < g.height; y++ {
			g.grid[x][y].g = -1
			g.grid[x][y].last = nil
		}
		g.opened[x] = make(map[int]*Node)
		g.closed[x] = make(map[int]*Node)
	}
	g.current = nil
	g.opened[g.src_x][g.src_y] = g.grid[g.src_x][g.src_y]
}

func (g *Grid) Draw(win *pixelgl.Window) {
	g.batch.Clear()
	for _, v := range g.closed {
		for _, iv := range v {
			iv.Draw(g.batch, red)
		}
	}
	for _, v := range g.opened {
		for _, iv := range v {
			iv.Draw(g.batch, green)
		}
	}

	current := g.current
	for current != nil {
		current.Draw(g.batch, blue)
		current = current.last
	}

	for x := range g.grid {
		for y := range g.grid[x] {
			if !g.grid[x][y].traversable {
				g.grid[x][y].Draw(g.batch, black)
			}
		}
	}

	g.batch.Draw(win)
}

func (g *Grid) GetNextOpen() *Node {
	sum := 0
	for _, v := range g.opened {
		sum += len(v)
	}
	if sum == 0 {
		return nil
	}

	var minnode *Node
	for _, v := range g.opened {
		for _, iv := range v {
			if iv.Cost() < minnode.Cost() {
				minnode = iv
			}
		}
	}

	// remove minnode from open set
	delete(g.opened[int(minnode.x)], int(minnode.y))

	return minnode
}

func (g *Grid) UpdateNeighbor(current, neighbor *Node) {
	if g.closed[int(neighbor.x)][int(neighbor.y)] != nil {
		return
	}

	if neighbor.g > current.g+current.dist(neighbor.x, neighbor.y) || neighbor.g == -1 {
		neighbor.g = current.g + current.dist(neighbor.x, neighbor.y)
		neighbor.last = current
	}

	g.opened[int(neighbor.x)][int(neighbor.y)] = neighbor
}

func (g *Grid) Step() (complete bool) {
	if g.current != nil {
		if int(g.current.x) == g.dst_x && int(g.current.y) == g.dst_y {
			return true
		}
		g.closed[int(g.current.x)][int(g.current.y)] = g.current
	}

	current := g.GetNextOpen()
	if current == nil {
		return true
	}
	g.current = current

	if g.current.x > 0 {
		left := g.grid[int(g.current.x)-1][int(g.current.y)]
		if left.traversable {
			g.UpdateNeighbor(g.current, left)
		}
	}
	if int(g.current.x) < g.width-1 {
		right := g.grid[int(g.current.x)+1][int(g.current.y)]
		if right.traversable {
			g.UpdateNeighbor(g.current, right)
		}
	}
	if g.current.y > 0 {
		down := g.grid[int(g.current.x)][int(g.current.y)-1]
		if down.traversable {
			g.UpdateNeighbor(g.current, down)
		}
	}
	if int(g.current.y) < g.height-1 {
		up := g.grid[int(g.current.x)][int(g.current.y)+1]
		if up.traversable {
			g.UpdateNeighbor(g.current, up)
		}
	}
	if g.current.x > 0 && g.current.y > 0 {
		left_down := g.grid[int(g.current.x)-1][int(g.current.y)-1]
		if left_down.traversable {
			g.UpdateNeighbor(g.current, left_down)
		}
	}
	if g.current.x > 0 && int(g.current.y) < g.height-1 {
		left_up := g.grid[int(g.current.x)-1][int(g.current.y)+1]
		if left_up.traversable {
			g.UpdateNeighbor(g.current, left_up)
		}
	}
	if int(g.current.x) < g.width-1 && g.current.y > 0 {
		left_down := g.grid[int(g.current.x)+1][int(g.current.y)-1]
		if left_down.traversable {
			g.UpdateNeighbor(g.current, left_down)
		}
	}
	if int(g.current.x) < g.width-1 && int(g.current.y) < g.height-1 {
		left_up := g.grid[int(g.current.x)+1][int(g.current.y)+1]
		if left_up.traversable {
			g.UpdateNeighbor(g.current, left_up)
		}
	}

	return false
}
