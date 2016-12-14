package main

import (
	"github.com/lhz/breadbox/gfx"

	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
)

var palette = gfx.Pepto

var mask1 = []byte{128, 32, 8, 2}
var mask2 = []byte{64, 16, 4, 1}
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

	switch *paletteName {
	case "pepto":
		palette = gfx.Pepto
	case "colodore":
		palette = gfx.Colodore
	case "vice":
		palette = gfx.Vice
	case "levy":
		palette = gfx.Levy
	default:
		usage()
	}

	sourceFile := flag.Arg(0)
	targetFile := flag.Arg(1)

	koala := make([]byte, 10003)
	f, err := os.Open(sourceFile)
	if err != nil {
		log.Fatalf("Can't open file %s for reading: %v", sourceFile, err)
		return
	}

	_, err = f.Read(koala)
	if err != nil {
		log.Fatalf("Can't read from file %s: %v", sourceFile, err)
		return
	}

	img := image.NewPaletted(image.Rect(0, 0, 320, 200), palette)

	bkg := koala[10002] & 0x0F

	for row := 0; row < 25; row++ {
		for col := 0; col < 40; col++ {
			scr := koala[8002+row*40+col]
			cmp := koala[9002+row*40+col]
			for y := 0; y < 8; y++ {
				byte := koala[2+row*320+col*8+y]
				for x := 0; x < 4; x++ {
					b1, b2 := byte&mask1[x], byte&mask2[x]
					if b1 > 0 && b2 > 0 {
						c = cmp & 0x0F
					} else if b1 > 0 {
						c = scr & 0x0F
					} else if b2 > 0 {
						c = (scr & 0xF0) >> 4
					} else {
						c = bkg
					}
					img.Set(col*8+x*2, row*8+y, palette[c])
					img.Set(col*8+x*2+1, row*8+y, palette[c])
				}
			}
		}
	}

	f, _ = os.OpenFile(targetFile, os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	png.Encode(f, img)
}
