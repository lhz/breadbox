package gfx

import (
	"bytes"
	"errors"
	"fmt"
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
	mColors []byte
	Clashes []Clash
}

type Clash struct {
	X      int
	Y      int
	Colors []byte
}

// NewImage reads an image from a PNG file and returns a Image pointer
func NewImage(filename string, mcol bool, bgColor byte) *Image {
	img := pngImage(filename)
	pal := PaletteBestMatch(img.Palette)
	return &Image{img, pal.Colors, remapIndices(img.Palette, pal.Colors), mcol, bgColor, nil, []Clash{}}
}

// MulticolorImage reads an image from a PNG file and returns a Image pointer
func MulticolorImage(filename string, bgColor byte) *Image {
	return NewImage(filename, true, bgColor)
}

// HiresImage reads an image from a PNG file and returns a Image pointer
func HiresImage(filename string, bgColor byte) *Image {
	return NewImage(filename, false, bgColor)
}

func (image *Image) Palette() []string {
	strings := make([]string, len(image.palette))
	for i, c := range image.palette {
		r, g, b, _ := c.RGBA()
		strings[i] = fmt.Sprintf("%02X%02X%02X", r, g, b)
	}
	return strings
}

func (image *Image) PixelAt(x, y int) byte {
	if (img.Point{x, y}).In(image.img.Rect) {
		//if x == 28 && y == 152 {
		//fmt.Printf("cia=%d, colors=%+v\n", image.img.ColorIndexAt(x, y), image.colors)
		//}
		return image.colors[image.img.ColorIndexAt(x, y)]
	} else {
		return image.BgColor
	}
}

func (image *Image) HiresByte(x, y, c int) byte {
	value := byte(0)
	for i := 0; i < 8; i++ {
		p := image.PixelAt(x+i, y)
		if p == byte(c) {
			value += byte(1 << uint(7-i))
		}
	}
	return value
}

// Pixels return an array of multicolor pixels (color indices) from the given area in the Image
func (image *Image) Pixels(xoffset, yoffset, width, height int) [][]byte {
	pix := make([][]byte, height)
	for y := 0; y < height; y++ {
		pix[y] = make([]byte, width)
		for x := 0; x < width; x++ {
			if image.mcol {
				pix[y][x] = image.PixelAt((xoffset+x)*2, yoffset+y)
			} else {
				pix[y][x] = image.PixelAt(xoffset+x, yoffset+y)
			}
		}
	}
	return pix
}

// MulticolorSprite extracts a multicolor sprite as a 64-byte array
func (image *Image) MulticolorSprite(xoffset, yoffset int, colors []byte) []byte {
	spr := make([]byte, 64)
	pixels := image.Pixels(xoffset, yoffset, 12, 21)
	//fmt.Printf("[%d,%d] %v\n", xoffset, yoffset, pixels)
	for y := 0; y < 21; y++ {
		for c := 0; c < 3; c++ {
			i := y*3 + c
			for x := 0; x < 4; x++ {
				n := bytes.IndexByte(colors, pixels[y][c*4+x])
				if n >= 0 {
					spr[i] = (spr[i] << 2) + byte(n)
				}
			}
		}
	}
	return spr
}

// MulticolorChar extracts a 4x8 pixels multicolor cell as a 10-byte array,
// the first 8 bytes are bitmap data, followed by a screen byte and
// a colmap byte
func (image *Image) MulticolorChar(xoffset, yoffset int) ([]byte, error) {
	char := make([]byte, 8)
	pixels := image.Pixels(xoffset, yoffset, 4, 8)
	colors := []byte{image.BgColor}
	if image.mColors == nil {
		return []byte{}, errors.New("MulticolorChar called without setting colors.")
	}
	colors = append(colors, image.mColors...)
	//fmt.Printf("colors: %+v, pixels: %+v\n", colors, pixels)
	for y := 0; y < 8; y++ {
		for x := 0; x < 4; x++ {
			char[y] = (char[y] << 2) + byte(bytes.IndexByte(colors, pixels[y][x]))
		}
	}
	return char, nil
}

// MulticolorCell extracts a 4x8 pixels multicolor cell as a 10-byte array,
// the first 8 bytes are bitmap data, followed by a screen byte and
// a colmap byte
func (image *Image) MulticolorCell(xoffset, yoffset int) ([]byte, error) {
	cell := make([]byte, 10)
	pixels := image.Pixels(xoffset, yoffset, 4, 8)
	//fmt.Printf("x=%d, y=%d, pixels: %+v\n", xoffset, yoffset, pixels)
	colors := colorsUsed(pixels, image.BgColor)
	for len(colors) < 4 {
		colors = append(colors, 0)
	}
	for y := 0; y < 8; y++ {
		for x := 0; x < 4; x++ {
			cell[y] = (cell[y] << 2) + byte(bytes.IndexByte(colors, pixels[y][x]))
		}
	}
	cell[8] = colors[1]*16 + colors[2]
	cell[9] = colors[3]
	if len(colors) > 4 {
		image.AddClash(xoffset, yoffset, colors)
		return cell, fmt.Errorf("Too many colors in cell at x=%3d, y=%3d: %v\n", xoffset, yoffset, colors)
	}
	return cell, nil
}

// HiresCell extracts a 4x8 pixels multicolor cell as a 9-byte array,
// the first 8 bytes are bitmap data, followed by a screen byte
func (image *Image) HiresCell(xoffset, yoffset int) ([]byte, error) {
	cell := make([]byte, 9)
	pixels := image.Pixels(xoffset, yoffset, 8, 8)
	colors := colorsUsedNoBg(pixels)
	for len(colors) < 2 {
		colors = append(colors, 0)
	}
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			cell[y] = (cell[y] << 1) + byte(bytes.IndexByte(colors, pixels[y][x]))
		}
	}
	cell[8] = colors[0]*16 + colors[1]
	if len(colors) > 2 {
		image.AddClash(xoffset, yoffset, colors)
		return cell, fmt.Errorf("Too many colors in cell at x=%3d, y=%3d: %v\n", xoffset, yoffset, colors)
	}
	return cell, nil
}

// Koala extracts a full-screen 160x200 multicolor image in Koala format
func (image *Image) Hires(xoffset, yoffset int) *Hires {
	hires := Hires{
		Bitmap: make([]byte, 8000),
		Screen: make([]byte, 1000)}
	for row := 0; row < 25; row++ {
		for col := 0; col < 40; col++ {
			cell, err := image.HiresCell(xoffset+col*8, yoffset+row*8)
			if err != nil {
				os.Stderr.WriteString(err.Error())
			}
			copy(hires.Bitmap[row*320+col*8:], cell[0:8])
			hires.Screen[row*40+col] = cell[8]
		}
	}
	return &hires
}

// Koala extracts a full-screen 160x200 multicolor image in Koala format
func (image *Image) Koala(xoffset, yoffset int) *Koala {
	koala := Koala{
		Bitmap:  make([]byte, 8000),
		Screen:  make([]byte, 1000),
		Colmap:  make([]byte, 1000),
		BgColor: image.BgColor}
	for row := 0; row < 25; row++ {
		for col := 0; col < 40; col++ {
			cell, err := image.MulticolorCell(xoffset+col*4, yoffset+row*8)
			if err != nil {
				os.Stderr.WriteString(err.Error())
			}
			copy(koala.Bitmap[row*320+col*8:], cell[0:8])
			koala.Screen[row*40+col] = cell[8]
			koala.Colmap[row*40+col] = cell[9]
		}
	}
	return &koala
}

func (image *Image) SetMultiColors(main, mcol1, mcol2 byte) {
	if image.mcol {
		image.mColors = []byte{main, mcol1, mcol2}
	} else {
		panic("Can't call SetMultiColors on hires image.")
	}
}

func (image *Image) AddClash(xoffset, yoffset int, colors []byte) {
	image.Clashes = append(image.Clashes,
		Clash{X: xoffset, Y: yoffset, Colors: colors})
}

func (image *Image) WriteClashesToPNG(filename string) {
	t := img.NewRGBA(img.Rectangle{img.Point{0, 0}, img.Point{1042, 652}})
	for y := 0; y < 652; y++ {
		for x := 0; x < 1042; x++ {
			// Space between cells
			if x%26 < 2 || y%26 < 2 {
				t.Set(x, y, color.RGBA{0x14, 0x14, 0x14, 0xff})
			} else {
				xx := ((x-2)/26)*8 + ((x-2)%26)/3
				yy := ((y-2)/26)*8 + ((y-2)%26)/3
				c := image.PixelAt(xx, yy)
				t.Set(x, y, image.palette[c])
			}
		}
	}
	ccol := color.RGBA{255, 0, 0, 0xff}
	for _, clash := range image.Clashes {
		x := (clash.X / 4) * 26
		y := (clash.Y / 8) * 26
		for i := 0; i < 28; i += 2 {
			t.Set(x+i, y, ccol)
			t.Set(x+i, y+1, ccol)
			t.Set(x+i, y+26, ccol)
			t.Set(x+i, y+27, ccol)
			t.Set(x, y+i, ccol)
			t.Set(x+1, y+i, ccol)
			t.Set(x+26, y+i, ccol)
			t.Set(x+27, y+i, ccol)
		}
	}

	f, err := os.Create(filename)
	if err != nil {
		fmt.Errorf("Failed to create file %v: %v", filename, err)
		return
	}
	png.Encode(f, t)
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
	//fmt.Printf("remapIndices:\n  from: %+v\n  to: %+v\n", from, to)
	colors := make([]byte, 16)
	for index, color := range to {
		//i := from.Index(color) // BUGGY! Finds colors not used
		i := -1
		for fi, fc := range from {
			if fc == color {
				i = fi
			}
		}
		if i >= 0 {
			//fmt.Printf("    Found color %+v at from[%d], vic-index is %d\n", color, i, index)
			colors[i] = byte(index)
		} else {
			//fmt.Printf("    Color with vic-index %d was not found in image.", index)
		}
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

func colorsUsedNoBg(pixels [][]byte) []byte {
	counts := histogram(pixels)
	colors := []byte{}
	for key, _ := range counts {
		colors = append(colors, key)
	}
	return colors
}
