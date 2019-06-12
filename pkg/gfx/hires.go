package gfx

import "bytes"

// Hires represents a full-screen image in Hires format
type Hires struct {
	Bitmap []byte
	Screen []byte
}

// Bytes returns Hires format as raw bytes. If align is false, there
// will be no padding between data segments. If align is true, screen
// data will be aligned to 1024 bytes offset.
func (hires *Hires) Bytes(align bool) []byte {
	if align {
		return bytes.Join([][]byte{
			hires.Bitmap, make([]byte, 192),
			hires.Screen, make([]byte, 24)}, []byte{})
	} else {
		return bytes.Join([][]byte{
			hires.Bitmap, hires.Screen}, []byte{})
	}
}
