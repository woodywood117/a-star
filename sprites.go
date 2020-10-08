package main

import (
	"github.com/faiface/pixel"
	"golang.org/x/image/colornames"
)

var pict *pixel.PictureData
var red, green, blue, white *pixel.Sprite

func init() {
	pict = pixel.MakePictureData(pixel.R(0, 0, scale, scale*4))
	for i := 0; i < len(pict.Pix); i++ {
		switch i / int(scale*scale) {
		case 0:
			pict.Pix[i] = colornames.Lightgreen
		case 1:
			pict.Pix[i] = colornames.Lightblue
		case 2:
			pict.Pix[i] = colornames.Indianred
		case 3:
			pict.Pix[i] = colornames.White
		}
	}
	green = pixel.NewSprite(pict, pixel.R(0, 0, scale, scale))
	blue = pixel.NewSprite(pict, pixel.R(0, scale, scale, scale*2))
	red = pixel.NewSprite(pict, pixel.R(0, scale*2, scale, scale*3))
	white = pixel.NewSprite(pict, pixel.R(0, scale*3, scale, scale*4))
}
