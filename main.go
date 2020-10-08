package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"math/rand"
	"time"
)

// Global constants (window height/width, node size, and framerate limit)
const winx, winy float64 = 600, 600
const scale float64 = 10

const frametime time.Duration = 0

func run() {
	// Set up window
	cfg := pixelgl.WindowConfig{
		Title:  "A*",
		Bounds: pixel.R(0, 0, winx, winy),
		VSync:  false,
		Icon:   []pixel.Picture{pixel.MakePictureData(pixel.R(0, 0, 16, 16))},
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Create the grid that contains all of the nodes to search through
	grid := NewGrid(int(winx/scale), int(winy/scale), 0, 0, int(winx/scale)-1, int(winy/scale)-1)

	// Game loop
	t := time.Now()
	for !win.Closed() {
		if win.JustPressed(pixelgl.KeySpace) {
			grid = NewGrid(int(winx/scale), int(winy/scale), 0, 0, int(winx/scale)-1, int(winy/scale)-1)
		}

		if win.JustPressed(pixelgl.KeyR) {
			grid.Restart()
			grid.Step()
		}

		if win.JustPressed(pixelgl.KeyEscape) {
			return
		}

		dt := time.Since(t)
		if dt < frametime {
			continue
		}

		win.Clear(colornames.White)
		grid.Step()
		grid.Draw(win)

		win.Update()
		t = time.Now()
	}
}

func main() {
	rand.Seed(time.Now().Unix())
	pixelgl.Run(run)
}
