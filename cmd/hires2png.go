package main

import (
	"github.com/lhz/breadbox/pkg/gfx"

	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
)

var palette *gfx.Palette

var mask = []byte{128, 64, 32, 16, 8, 4, 2, 1}
var c byte

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %v [flags] <source> <target>\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	paletteName := flag.String("p", "pepto", "Name of palette to use [pepto|colodore|vice|levy]")

	flag.Parse()

	if len(flag.Args()) != 2 {
		usage()
	}

	palette = gfx.Palettes[*paletteName]

	sourceFile := flag.Arg(0)
	targetFile := flag.Arg(1)

	hires := make([]byte, 9194)
	f, err := os.Open(sourceFile)
	if err != nil {
		log.Fatalf("Can't open file %s for reading: %v", sourceFile, err)
		return
	}

	size, err := f.Read(hires)
	if err != nil {
		log.Fatalf("Can't read from file %s: %v", sourceFile, err)
		return
	}

	screen := 8002
	if size >= 9194 {
		screen = 8194
	}

	img := image.NewPaletted(image.Rect(0, 0, 320, 200), gfx.Pepto.Colors)

	for row := 0; row < 25; row++ {
		for col := 0; col < 40; col++ {
			scr := hires[screen+row*40+col]
			for y := 0; y < 8; y++ {
				byte := hires[2+row*320+col*8+y]
				for x := 0; x < 8; x++ {
					if byte&mask[x] > 0 {
						c = (scr & 0xF0) >> 4
					} else {
						c = scr & 0x0F
					}
					img.Set(col*8+x, row*8+y, palette.Colors[c])
				}
			}
		}
	}

	f, _ = os.OpenFile(targetFile, os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	png.Encode(f, img)
}
