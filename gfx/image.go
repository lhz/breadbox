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
	mcol    bool
	BgColor byte
}

// NewImage reads an image from a PNG file and returns a Image pointer
func NewImage(filename string, mcol bool, bgColor byte) *Image {
	img := pngImage(filename)
	return &Image{img, Pepto, remapIndices(img.Palette, Pepto), mcol, bgColor}
}

// MulticolorImage reads an image from a PNG file and returns a Image pointer
func MulticolorImage(filename string, bgColor byte) *Image {
	return NewImage(filename, true, bgColor)
}

// HiresImage reads an image from a PNG file and returns a Image pointer
func HiresImage(filename string, bgColor byte) *Image {
	return NewImage(filename, false, bgColor)
}

// KoalaImage reads an image from a PNG file and returns a Koala pointer
func KoalaImage(filename string, bgColor byte) *Koala {
	image := MulticolorImage(filename, bgColor)
	return image.Koala(0, 0)
}

func (image *Image) PixelAt(x, y int) byte {
	if (img.Point{x, y}).In(image.img.Rect) {
		return image.colors[image.img.ColorIndexAt(x, y)]
	} else {
		return image.BgColor
	}
}

// Pixels return an array of multicolor pixels (color indices) from the given area in the Image
func (image *Image) Pixels(xoffset, yoffset, width, height int) [][]byte {
	pix := make([][]byte, height)
	for y := 0; y < height; y++ {
		pix[y] = make([]byte, width)
		for x := 0; x < width; x++ {
			if image.mcol {
				pix[y][x] = image.PixelAt((xoffset + x) * 2, yoffset + y)
			} else {
				pix[y][x] = image.PixelAt(xoffset + x, yoffset + y)
			}
		}
	}
	return pix
}

// MulticolorCell extracts a 4x8 pixels multicolor cell as a 10-byte array,
// the first 8 bytes are bitmap data, followed by a screen byte and
// a colmap byte
func (image *Image) MulticolorCell(xoffset, yoffset int) []byte {
	cell := make([]byte, 10)
	pixels := image.Pixels(xoffset, yoffset, 4, 8)
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
	for y := 0; y < 8; y++ {
		for x := 0; x < 4; x++ {
			cell[y] = (cell[y] << 2) + index[pixels[y][x]]
		}
	}
	cell[8] = colors[1] * 16 + colors[2]
	cell[9] = colors[3]
	return cell
}

// Koala extracts a full-screen 160x200 multicolor image in Koala format
func (image *Image) Koala(xoffset, yoffset int) *Koala {
	koala := Koala{
		Bitmap: make([]byte, 8000),
		Screen: make([]byte, 1000),
		Colmap: make([]byte, 1000),
		BgColor: image.BgColor}
	for row := 0; row < 25; row++ {
		for col := 0; col < 40; col++ {
			cell := image.MulticolorCell(xoffset + col * 4, yoffset + row * 8)
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

// Bytes returns Koala format as raw bytes. If align is false, there
// will be no padding between data segments. If align is true, screen
// and color map data will be aligned to 1024 bytes offsets.
func (koala *Koala) Bytes(align bool) []byte {
	if align {
		return bytes.Join([][]byte{
			koala.Bitmap, make([]byte, 192),
			koala.Screen, make([]byte, 24),
			koala.Colmap, []byte{koala.BgColor}}, []byte{})
	} else {
		return bytes.Join([][]byte{
			koala.Bitmap, koala.Screen, koala.Colmap,
			[]byte{koala.BgColor}}, []byte{})
	}
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
