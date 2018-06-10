package main

import (
	"github.com/lhz/breadbox/file"
	"github.com/lhz/breadbox/gfx"

	"flag"
	"fmt"
	//"log"
	"os"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %v [flags] <source> <target>\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {

	var bgCol int
	flag.IntVar(&bgCol, "b", 0, "Background color (0-15)")

	flag.Parse()

	if len(flag.Args()) != 2 {
		usage()
	}

	sourceFile := flag.Arg(0)
	targetFile := flag.Arg(1)

	image := gfx.NewImage(sourceFile, true, byte(bgCol))
	koala := image.Koala(0, 0)

	file.WriteBin(targetFile, 0x4000, koala.Bytes(false))
}
