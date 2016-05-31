package file

import (
	"io/ioutil"
	// "bytes"
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
	copy(p.memory[address:address+len(content)-2], content[2:])

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
}
