package file

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
)

type Memory struct {
	content []byte
}

func ReadMemory(source string) (*Memory, error) {
	content, err := ioutil.ReadFile(source)
	if err != nil {
		panic(err)
	}
	switch len(content) {
	case 0x10000:
		return &Memory{content}, nil
	case 0x10002:
		return &Memory{content[2:]}, nil
	default:
		return nil, errors.New(fmt.Sprintf("File %s does not look like a memory dump.", source))
	}
}

func (m *Memory) Peek(address int) byte {
	if address < 0 || address > 0x10000 {
		return byte(0)
	}
	value := m.content[address]
	log.Printf("Peeked at $%03X and found $%02X", address, value)
	return value
}

func (m *Memory) Read(start int, length int) ([]byte, error) {
	if start+length > 0x10000 {
		return nil, errors.New("Can't read past end of memory.")
	}
	return m.content[start : start+length], nil
}

func (m *Memory) ScreenMatrix() []byte {
	bank := 0x4000 * (3 - (int(m.Peek(0xDD00)) % 4))
	screen := bank + 0x0400*(int(m.Peek(0xD018))>>4)
	log.Printf("Screen is at $%04X", screen)
	bytes, err := m.Read(screen, 1000)
	if err != nil {
		panic(err)
	}
	return bytes
}

func (m *Memory) ColorMap() []byte {
	bytes, err := m.Read(0xD800, 1000)
	if err != nil {
		panic(err)
	}
	return bytes
}

func (m *Memory) ColorMapCompact() []byte {
	bytes := m.ColorMap()
	log.Printf("bytes is %d bytes long", len(bytes))
	compact := make([]byte, 500)
	for i := 0; i < 500; i++ {
		log.Println(i)
		compact[i] = (bytes[i*2]&15)*16 + (bytes[i*2+1] & 15)
	}
	return compact
}
