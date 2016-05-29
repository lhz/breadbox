package gfx

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

// Pic represents a complete picture converted from a PNG image
type Pic struct {
	image   *image.Paletted
	palette []color.Color
	colors  []byte
}

// Pixels return an array of pixels (color indices) from the given area in the Pic
func (pic *Pic) Pixels(x, y, w, h int) [][]byte {
	pix := make([][]byte, h)
	for yy := 0; yy < h; yy++ {
		pix[yy] = make([]byte, w)
		for xx := 0; xx < w; xx++ {
			pix[yy][xx] = pic.colors[pic.image.ColorIndexAt(x + xx, y + yy)]
		}
	}
	return pix
}

// Pixels return an array of multicolor pixels (color indices) from the given area in the Pic
func (pic *Pic) PixelsMC(x, y, w, h int) [][]byte {
	pix := make([][]byte, h)
	for yy := 0; yy < h; yy++ {
		pix[yy] = make([]byte, w)
		for xx := 0; xx < w; xx++ {
			pix[yy][xx] = pic.colors[pic.image.ColorIndexAt(2 * (x + xx), y + yy)]
		}
	}
	return pix
}

// ReadPNG reads an image from a PNG file and returns a two-dimensional
// array of bytes that are color indices in the VIC-II palette (0-15)
func ReadPNG(filename string) *Pic {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decoded, err := png.Decode(file)
	if err != nil {
		panic(err)
	}

	var pic Pic
	pic.image   = decoded.(*image.Paletted)
	pic.palette = Pepto 
	pic.colors  = remapIndices(pic.image.Palette, pic.palette)
	return &pic
}

func remapIndices(from color.Palette, to []color.Color) []byte {
	colors := make([]byte, 16)
	for index, color := range to {
		i := from.Index(color)
		colors[i] = byte(index)
	}
	return colors
}
