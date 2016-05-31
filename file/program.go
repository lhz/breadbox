package file

import (
	"fmt"
	"io/ioutil"
)

const (
	MemSize = 1<<16
)

type Program struct {
	memory []byte
	min    int
	max    int
}

func NewProgram() *Program {
	return &Program{make([]byte, MemSize), MemSize-1, 0}
}

func (p *Program) Inject(filename string) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	address := int(content[0]) + 256 * int(content[1])
	copy(p.memory[address:MemSize], content[2:])
	fmt.Printf("Copying to $%04x-$%04x from file %q\n", address, address+len(content)-3, filename)

	if address < p.min {
		p.min = address
	}
	if address + len(content) - 2 > p.max {
		p.max = address + len(content) - 2
	}
	return
}

func (p *Program) WriteFile(target string) {
	WriteBin(target, p.min, p.memory[p.min:p.max+1])
	fmt.Printf("Result file %q spans range $%04x-$%04x.\n", target, p.min, p.max)
}
