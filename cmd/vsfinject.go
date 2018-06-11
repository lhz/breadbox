package main

import (
	"fmt"
	"os"

	"github.com/lhz/breadbox/pkg/file"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Fprintf(os.Stderr, "Usage: %v <source> <target> [binaries]+\n", os.Args[0])
		os.Exit(1)
	}

	snapshot := file.NewSnapshot(os.Args[1])

	for _, filename := range os.Args[3:] {
		snapshot.Inject(filename)
	}

	snapshot.WriteFile(os.Args[2])
}
