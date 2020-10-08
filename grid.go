package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"math/rand"
)

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
			g.grid[x][y] = NewNode(float64(x), float64(y))
			g.grid[x][y].Hscore(float64(dst_x), float64(dst_y))
		}
		g.opened[x] = make(map[int]*Node)
		g.closed[x] = make(map[int]*Node)
	}

	g.InitMaze()
	g.Restart()

	g.opened[src_x][src_y] = g.grid[src_x][src_y]
	return g
}

func (g *Grid) InitMaze() {
	stack := []*Node{}

	current := g.grid[g.src_x][g.src_y]
	current.last = current
	stack = append(stack, current)

	for {
		if len(stack) == 0 {
			break
		}

		// Pop from stack
		current = stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// Select random unvisited neighbor
		unvisited := []*Node{}
		if current.x > 0 {
			left := g.grid[int(current.x)-1][int(current.y)]
			if left.last == nil {
				unvisited = append(unvisited, left)
			}
		}
		if int(current.x) < g.width-1 {
			right := g.grid[int(current.x)+1][int(current.y)]
			if right.last == nil {
				unvisited = append(unvisited, right)
			}
		}
		if current.y > 0 {
			down := g.grid[int(current.x)][int(current.y)-1]
			if down.last == nil {
				unvisited = append(unvisited, down)
			}
		}
		if int(current.y) < g.height-1 {
			up := g.grid[int(current.x)][int(current.y)+1]
			if up.last == nil {
				unvisited = append(unvisited, up)
			}
		}
		if len(unvisited) == 0 {
			continue
		}
		stack = append(stack, current)
		next := unvisited[rand.Intn(len(unvisited))]

		// Remove wall between cells
		if current.x > next.x {
			current.left = false
			next.right = false
		}
		if current.x < next.x {
			current.right = false
			next.left = false
		}
		if current.y > next.y {
			current.down = false
			next.up = false
		}
		if current.y < next.y {
			current.up = false
			next.down = false
		}

		next.last = next
		stack = append(stack, next)
	}
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
	for x := range g.grid {
		for y := range g.grid[x] {
			g.grid[x][y].Draw(g.batch, white)
		}
	}

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

	// Left neighbor
	if !g.current.left {
		left := g.grid[int(g.current.x)-1][int(g.current.y)]
		g.UpdateNeighbor(g.current, left)
	}

	// Right neighbor
	if !g.current.right {
		right := g.grid[int(g.current.x)+1][int(g.current.y)]
		g.UpdateNeighbor(g.current, right)
	}

	// Bottom neighbor
	if !g.current.down {
		down := g.grid[int(g.current.x)][int(g.current.y)-1]
		g.UpdateNeighbor(g.current, down)
	}

	// Top neighbor
	if !g.current.up {
		up := g.grid[int(g.current.x)][int(g.current.y)+1]
		g.UpdateNeighbor(g.current, up)
	}

	return false
}
