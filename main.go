package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"math/rand"
	"time"
)

const winx, winy float64 = 600, 600
const scale float64 = 2

//const frametime time.Duration = time.Second / 60
const frametime time.Duration = 0

func run() {
	// Set up window
	cfg := pixelgl.WindowConfig{
		Title:  "A*",
		Bounds: pixel.R(0, 0, winx, winy),
		VSync:  false,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	grid := NewGrid(int(winx/scale), int(winy/scale), 0, 0, int(winx/scale)-1, int(winy/scale)-1)

	// Game loop
	t := time.Now()
	for !win.Closed() {
		dt := time.Since(t)
		if dt < frametime {
			continue
		}
		//fmt.Println(dt)

		win.Clear(colornames.White)
		complete := grid.Step()
		grid.Draw(win)
		if complete {
			//fmt.Println("Done!")
			//return
		}

		win.Update()
		t = time.Now()
	}
}

func main() {
	rand.Seed(time.Now().Unix())
	pixelgl.Run(run)
}
