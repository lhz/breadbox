package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/lhz/breadbox/file"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %v <memdump> <target> [address]\n", os.Args[0])
		os.Exit(1)
	}
	source := os.Args[1]
	target := os.Args[2]
	address := 0x8000
	var err error
	if len(os.Args) == 4 {
		address, err = strconv.Atoi(os.Args[3])
		if err != nil {
			log.Fatal(err)
		}
	}

	memory, err := file.ReadMemory(source)
	if err != nil {
		log.Fatal(err)
	}

	bytes := make([]byte, 0)
	bytes = append(bytes, memory.ScreenMatrix()...)
	bytes = append(bytes, memory.ColorMapCompact()...)
	bytes = append(bytes, memory.Peek(0xD021)&15)

	file.WriteBin(target, address, bytes)
}
