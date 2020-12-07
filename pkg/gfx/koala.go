package gfx

import (
	"bytes"
)

// Koala represents a full-screen image in KoalaPainter format
type Koala struct {
	Bitmap  []byte
	Screen  []byte
	Colmap  []byte
	BgColor byte
}

// Bytes returns Koala format as raw bytes. If align is false, there
// will be no padding between data segments. If align is true, screen
// and color map data will be aligned to 1024 bytes offsets.
func (koala *Koala) Bytes(align bool, front bool) []byte {
	if align {
		if (front) {
			//fmt.Printf("Putting screen and colmap data in front, aligned.\n")
			return bytes.Join([][]byte{
				koala.Screen, make([]byte, 24),
				koala.Colmap, make([]byte, 23),
				[]byte{koala.BgColor},
				koala.Bitmap,
			}, []byte{})
		} else {
			//fmt.Printf("Putting screen and colmap data in back, aligned.\n")
			return bytes.Join([][]byte{
				koala.Bitmap, make([]byte, 192),
				koala.Screen, make([]byte, 24),
				koala.Colmap, []byte{koala.BgColor},
			}, []byte{})
		}
	} else {
		if (front) {
			//fmt.Printf("Putting screen and colmap data in front, unaligned.\n")
			return bytes.Join([][]byte{
				koala.Screen, koala.Colmap,
				[]byte{koala.BgColor}, koala.Bitmap,
			}, []byte{})
		} else {
			//fmt.Printf("Putting screen and colmap data in back, unaligned.\n")
			return bytes.Join([][]byte{
				koala.Bitmap, koala.Screen, koala.Colmap,
				[]byte{koala.BgColor},
			}, []byte{})
		}
	}
}

// KoalaImage reads an image from a PNG file and returns a Koala pointer
func KoalaImage(filename string, bgColor byte) *Koala {
	image := MulticolorImage(filename, bgColor)
	return image.Koala(0, 0)
}
