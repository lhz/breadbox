package gfx

import (
	"bytes"
	img "image"
	"image/color"
	"image/png"
	"os"
)

// Image represents a complete picture converted from a PNG image
type Image struct {
	img     *img.Paletted
	palette []color.Color
	colors  []byte
	BgColor byte
}

// NewImage reads an image from a PNG file and returns a Image pointer
func NewImage(filename string, bgColor byte) *Image {
	img := pngImage(filename)
	return &Image{img, Pepto, remapIndices(img.Palette, Pepto), bgColor}
}

// Pixels return an array of pixels (color indices) from the given area in the Image
func (image *Image) Pixels(xoffset, yoffset, width, height int) [][]byte {
	pix := make([][]byte, height)
	for y := 0; y < height; y++ {
		pix[y] = make([]byte, width)
		for x := 0; x < width; x++ {
			pix[y][x] = image.colors[image.img.ColorIndexAt(xoffset + x, yoffset + y)]
		}
	}
	return pix
}

// Pixels return an array of multicolor pixels (color indices) from the given area in the Image
func (image *Image) PixelsMC(xoffset, yoffset, width, height int) [][]byte {
	pix := make([][]byte, height)
	for y := 0; y < height; y++ {
		pix[y] = make([]byte, width)
		for x := 0; x < width; x++ {
			pix[y][x] = image.colors[image.img.ColorIndexAt((xoffset + x) * 2, yoffset + y)]
		}
	}
	return pix
}

// CellMC extracts a 4x8 pixels multicolor cell as a 10-byte array,
// the first 8 bytes are bitmap data, followed by a screen byte and
// a colmap byte
func (image *Image) CellMC(xoffset, yoffset int) []byte {
	cell := make([]byte, 10)
	pixels := image.PixelsMC(xoffset, yoffset, 4, 8)
	colors := colorsUsed(pixels, image.BgColor)
	for len(colors) < 4 {
		colors = append(colors, 0)
	}
	index := make(map[byte]byte)
	for i, c := range colors {
		if _, found := index[c]; !found {
			index[c] = byte(i)
		}
	}
	var c, i byte
	for y := 0; y < 8; y++ {
		for x := 0; x < 4; x++ {
			c = pixels[y][x]
			i = index[c]
			cell[y] = (cell[y] << 2) + i
		}
	}
	cell[8] = colors[1] * 16 + colors[2]
	cell[9] = colors[3]
	return cell
}

// ToKoala extracts a full-screen 160x200 multicolor image in Koala format
func (image *Image) ToKoala(xoffset, yoffset int) *Koala {
	koala := Koala{
		Bitmap: make([]byte, 8000),
		Screen: make([]byte, 1000),
		Colmap: make([]byte, 1000),
		BgColor: image.BgColor}
	for row := 0; row < 25; row++ {
		for col := 0; col < 40; col++ {
			cell := image.CellMC(xoffset + col * 4, yoffset + row * 8)
			for i := 0; i < 8; i++ {
				koala.Bitmap[row * 320 + col * 8 + i] = cell[i]
			}
			koala.Screen[row * 40 + col] = cell[8]
			koala.Colmap[row * 40 + col] = cell[9]
		}
	}
	return &koala
}

// Koala represents a full-screen image in KoalaPainter format
type Koala struct {
	Bitmap  []byte
	Screen  []byte
	Colmap  []byte
	BgColor byte
}

// Bytes returns Koala format as raw bytes without padding
func (koala *Koala) Bytes() []byte {
	return bytes.Join([][]byte{
		koala.Bitmap, koala.Screen, koala.Colmap,
		[]byte{koala.BgColor}}, []byte{})
}

// BytesAligned returns Koala format as raw bytes with padding so that
// each segment (bitmap, screen and colmap) is aligned for direct use
func (koala *Koala) BytesAligned() []byte {
	return bytes.Join([][]byte{
		koala.Bitmap, make([]byte, 192),
		koala.Screen, make([]byte, 24),
		koala.Colmap, []byte{koala.BgColor}}, []byte{})
}

func pngImage(filename string) *img.Paletted {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decoded, err := png.Decode(file)
	if err != nil {
		panic(err)
	}
	return decoded.(*img.Paletted)
}

func remapIndices(from color.Palette, to []color.Color) []byte {
	colors := make([]byte, 16)
	for index, color := range to {
		i := from.Index(color)
		colors[i] = byte(index)
	}
	return colors
}

func histogram(pixels [][]byte) map[byte]int {
	counts := make(map[byte]int)
	for _, row := range pixels {
		for _, col := range row {
			counts[col]++
		}
	}
	return counts
}

func colorsUsed(pixels [][]byte, bgColor byte) []byte {
	counts := histogram(pixels)
	colors := []byte{bgColor}
	for key, _ := range counts {
		if key != bgColor {
			colors = append(colors, key)
		}
	}
	return colors
}
