// Package video implements constants and functions related to
// the VIC-II video interface chip, such as colors.

package gfx

import (
	"encoding/hex"
	"image/color"
	"strings"
)

type Palette struct {
	Name   string
	Colors []color.Color
}

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

var Palettes = make(map[string]*Palette)
var Pepto *Palette

func init() {
	pals := []*Palette{
		makePalette("colodore",
			"000000:ffffff:813338:75cec8:8e3c97:56ac4d:2e2c9b:edf171:8e5029:553800:c46c71:4a4a4a:7b7b7b:a9ff9f:706deb:b2b2b2"),
		makePalette("pepto",
			"000000:ffffff:68372b:70a4b2:6f3d86:588d43:352879:b8c76f:6f4f25:433900:9a6759:444444:6c6c6c:9ad284:6c5eb5:959595"),
		makePalette("levy",
			"040204:fcfefc:cc3634:84f2dc:cc5ac4:5cce34:4436cc:f4ee5c:d47e34:945e34:fc9a94:5c5a5c:8c8e8c:9cfe9c:74a2ec:c4c2c4"),
		makePalette("vice",
			"000000:fdfefc:be1a24:30e6c6:b41ae2:1fd21e:211bae:dff60a:b84104:6a3304:fe4a57:424540:70746f:59fe59:5f53fe:a4a7a2"),
		makePalette("vice_old",
			"000000:d5d5d5:72352c:659fa6:733a91:568d35:2e237d:aeb75e:774f1e:4b3c00:9c635a:474747:6b6b6b:8fc271:675db6:8f8f8f"),
		makePalette("vice_new",
			"000000:ffffff:b85438:8decff:ba56e4:79d949:553ee5:fbff79:bd7c1b:7e6400:f29580:6f716e:a2a4a1:cdff9d:a18aff:d3d5d2"),
	}
	for _, p := range pals {
		Palettes[p.Name] = p
	}
	Pepto = Palettes["pepto"]
}

func PaletteMatching(colors []color.Color) *Palette {
	var bestMatch *Palette
	var bestScore = 0
	for _, pal := range Palettes {
		score := 0
		for _, ccol := range colors {
			for _, pcol := range pal.Colors {
				if ccol == pcol {
					score++
				}
			}
		}
		if score >= bestScore {
			bestScore = score
			bestMatch = pal
		}
	}
	return bestMatch
}

func makePalette(name, values string) *Palette {
	colors := make([]color.Color, 16)
	for i, value := range strings.Split(values, ":") {
		colors[i] = hexColor(value)
	}
	return &Palette{Name: name, Colors: colors}
}

func hexColor(value string) color.Color {
	rgb, err := hex.DecodeString(value)
	if err != nil {
		panic(err)
	}
	return color.RGBA{rgb[0], rgb[1], rgb[2], 255}
}
