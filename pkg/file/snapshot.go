package file

import (
	"io/ioutil"
	"bytes"
)

type Snapshot struct {
	content []byte
	c64mem  int
}

func NewSnapshot(source string) *Snapshot {
	content, err := ioutil.ReadFile(source)
	if err != nil {
		panic(err)
	}

	header := bytes.Index(content, []byte("VICE Snapshot File"))
	if header != 0 {
		panic("No snapshot header found in source file.")
	}

	c64mem := bytes.Index(content, []byte("C64MEM"))
	if c64mem < 0 {
		panic("No C64MEM marker found in source file.")
	}

	return &Snapshot{content, c64mem + 26}
}

func (s *Snapshot) Inject(filename string) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	address := int(content[0]) + 256 * int(content[1])
	offset  := s.c64mem + address
	copy(s.content[offset:offset+len(content)-2], content[2:])
	return
}

func (s *Snapshot) WriteFile(target string) {
	ioutil.WriteFile(target, s.content, 0644)
}
