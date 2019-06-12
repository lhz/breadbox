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

	var align bool
	var address, xOffset, yOffset int
	flag.BoolVar(&align, "a", false, "Align screen to page")
	flag.IntVar(&address, "s", 0x4000, "Start address of koala output")
	flag.IntVar(&xOffset, "x", 0, "Offset X-coordinate of top left corner")
	flag.IntVar(&yOffset, "y", 0, "Offset Y-coordinate of top left corner")

	flag.Parse()

	if len(flag.Args()) != 2 {
		usage()
	}

	sourceFile := flag.Arg(0)
	targetFile := flag.Arg(1)

	image := gfx.HiresImage(sourceFile, byte(0))
	hires := image.Hires(xOffset, yOffset)

	file.WriteBin(targetFile, address, hires.Bytes(align))
}
