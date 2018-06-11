package main

import (
	"fmt"
	"os"

	"github.com/lhz/breadbox/pkg/file"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %v <target> [source]+\n", os.Args[0])
		os.Exit(1)
	}

	program := file.NewProgram()

	for _, filename := range os.Args[2:] {
		program.Inject(filename)
	}

	program.WriteFile(os.Args[1])
}
