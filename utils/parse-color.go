package utils

import (
	"errors"
	"image/color"
	"strconv"
)

var ErrMissingColorStartingHash = errors.New("missing color starting #")

func ParseColor(text string) (color.RGBA, error) {
	if text[0] != '#' {
		return color.RGBA{}, ErrMissingColorStartingHash
	}
	colourParse, err := strconv.ParseInt(text[1:], 16, 64)
	if err != nil {
		return color.RGBA{}, err
	}
	r1 := (colourParse >> 16) & 0xff
	g1 := (colourParse >> 8) & 0xff
	b1 := colourParse & 0xff
	return color.RGBA{R: uint8(r1), G: uint8(g1), B: uint8(b1), A: uint8(0xff)}, nil
}
