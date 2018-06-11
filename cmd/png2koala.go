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

	var bgCol, xOffset, yOffset int
	flag.IntVar(&bgCol, "b", 0, "Background color (0-15)")
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

	file.WriteBin(targetFile, 0x4000, koala.Bytes(false))
}
