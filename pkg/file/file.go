package file

import (
	"io/ioutil"
)

func WriteBin(filename string, address int, content []byte) {
	err := ioutil.WriteFile(filename, append(addressHeader(address), content...), 0644)
	if err != nil {
		panic(err)
	}
}

func addressHeader(address int) []byte {
	return []byte{byte(address % 256), byte(address >> 8)}
}
