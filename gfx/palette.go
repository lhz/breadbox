// Package video implements constants and functions related to
// the VIC-II video interface chip, such as colors.

package gfx

import (
	"encoding/hex"
	"image/color"
	"strings"
)

// Color index names
const (
	Black = iota
	White
	Red
	Cyan
	Purple
	Green
	Blue
	Yellow
	Orange
	Brown
	LightRed
	DarkGrey
	MediumGrey
	LightGreen
	LightBlue
	LightGrey
)

// Palettes each consisting of an array of 16 Color values
var (
	Pepto   = palette("000000:ffffff:68372b:70a4b2:6f3d86:588d43:352879:b8c76f:6f4f25:433900:9a6759:444444:6c6c6c:9ad284:6c5eb5:959595")
	Levy    = palette("040204:fcfefc:cc3634:84f2dc:cc5ac4:5cce34:4436cc:f4ee5c:d47e34:945e34:fc9a94:5c5a5c:8c8e8c:9cfe9c:74a2ec:c4c2c4")
	Vice    = palette("000000:fdfefc:be1a24:30e6c6:b41ae2:1fd21e:211bae:dff60a:b84104:6a3304:fe4a57:424540:70746f:59fe59:5f53fe:a4a7a2")
	ViceOld = palette("000000:d5d5d5:72352c:659fa6:733a91:568d35:2e237d:aeb75e:774f1e:4b3c00:9c635a:474747:6b6b6b:8fc271:675db6:8f8f8f")
)

func palette(values string) []color.Color {
	colors := make([]color.Color, 16)
	for i, value := range strings.Split(values, ":") {
		colors[i] = hexColor(value)
	}
	return colors
}

func hexColor(value string) color.Color {
	rgb, err := hex.DecodeString(value)
	if err != nil {
		panic(err)
	}
	return color.RGBA{rgb[0], rgb[1], rgb[2], 255}
}