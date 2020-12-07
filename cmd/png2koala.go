package main

import (
	"github.com/lhz/breadbox/pkg/file"
	"github.com/lhz/breadbox/pkg/gfx"

	"flag"
	"fmt"
	"os"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %v [flags] <source> <target>\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {

	var align, front bool
	var address, bgCol, xOffset, yOffset int
	var clashes string
	flag.BoolVar(&align, "a", false, "Align screen and colormap to page")
	flag.IntVar(&bgCol, "b", 0, "Background color (0-15)")
	flag.StringVar(&clashes, "c", "", "Output PNG showing color clashes.")
	flag.BoolVar(&front, "f", false, "Put screen and color map data in front of bitmap data")
	flag.IntVar(&address, "s", 0x4000, "Start address of koala output")
	flag.IntVar(&xOffset, "x", 0, "Offset X-coordinate of top left corner")
	flag.IntVar(&yOffset, "y", 0, "Offset Y-coordinate of top left corner")

	flag.Parse()

	if len(flag.Args()) != 2 {
		usage()
	}

	sourceFile := flag.Arg(0)
	targetFile := flag.Arg(1)

	image := gfx.NewImage(sourceFile, true, byte(bgCol))
	koala := image.Koala(xOffset, yOffset)

	if len(clashes) > 0 && len(image.Clashes) > 0 {
		image.WriteClashesToPNG(clashes)
	}

	file.WriteBin(targetFile, address, koala.Bytes(align, front))
}
