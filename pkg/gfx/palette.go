// Package video implements constants and functions related to
// the VIC-II video interface chip, such as colors.

package gfx

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"os/user"
	"path/filepath"
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

var PaletteMap = map[string]*Palette{}

func init() {
	var paletteConfig map[string]string
	usr, _ := user.Current()
	configFile := filepath.Join(usr.HomeDir, ".config/vic-palettes.json")
	paletteJSON, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Printf("Could not read palette config file", err)
		return
	}
	json.Unmarshal(paletteJSON, &paletteConfig)
	for name, valueList := range paletteConfig {
		MakePalette(name, strings.Replace(valueList, ",", ":", -1))
	}
}

func (p *Palette) Color(index int) color.Color {
	return p.Colors[index]
}

func PaletteByName(name string) *Palette {
	palette, ok := PaletteMap[name]
	if !ok {
		log.Printf("Invalid palette name %q, defaulting to %q.\n", name, "colodore")
		return PaletteMap["colodore"]
	}
	return palette
}

func PaletteBestMatch(colors []color.Color) *Palette {
	var bestMatch *Palette
	var bestScore = 0
	for _, pal := range PaletteMap {
		score := 0
		for _, ccol := range colors {
			for _, pcol := range pal.Colors {
				if ccol == pcol {
					score++
				}
			}
		}
		if pal.Name == "FORCE" {
			score = 9999
		}
		if score > bestScore {
			bestScore = score
			bestMatch = pal
		}
	}
	fmt.Printf("Palette %q won with a score of %d.\n", bestMatch.Name, bestScore)
	return bestMatch
}

func MakePalette(name string, values ...string) *Palette {
	valueStr := strings.Join(values, ":")
	colors := make([]color.Color, 16)
	for i, value := range strings.Split(valueStr, ":") {
		colors[i] = hexColor(value)
	}
	palette := &Palette{Name: name, Colors: colors}
	_, exists := PaletteMap[name]
	if !exists {
		PaletteMap[name] = palette
	}
	return palette
}

func hexColor(value string) color.Color {
	rgb, err := hex.DecodeString(value)
	if err != nil {
		panic(err)
	}
	return color.RGBA{rgb[0], rgb[1], rgb[2], 255}
}
