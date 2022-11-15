package shapes

import (
	"github.com/stretchr/testify/assert"
	"image/color"
	"prosperity-r-place/utils"
	"testing"
)

func TestPixelsInSquareFilled(t *testing.T) {
	c := color.RGBA{}
	assert.Equal(t, []utils.Pixel{
		{X: 0, Y: 0, Colour: c},
		{X: 1, Y: 0, Colour: c},
		{X: 2, Y: 0, Colour: c},
		{X: 3, Y: 0, Colour: c},
		{X: 4, Y: 0, Colour: c},
		{X: 0, Y: 1, Colour: c},
		{X: 1, Y: 1, Colour: c},
		{X: 2, Y: 1, Colour: c},
		{X: 3, Y: 1, Colour: c},
		{X: 4, Y: 1, Colour: c},
		{X: 0, Y: 2, Colour: c},
		{X: 1, Y: 2, Colour: c},
		{X: 2, Y: 2, Colour: c},
		{X: 3, Y: 2, Colour: c},
		{X: 4, Y: 2, Colour: c},
	}, PixelsInSquare(0, 0, 4, 2, c, true))
}

func TestPixelsInSquare(t *testing.T) {
	c := color.RGBA{}
	assert.Equal(t, []utils.Pixel{
		{X: 0, Y: 0, Colour: c},
		{X: 1, Y: 0, Colour: c},
		{X: 2, Y: 0, Colour: c},
		{X: 3, Y: 0, Colour: c},
		{X: 4, Y: 0, Colour: c},
		{X: 4, Y: 1, Colour: c},
		{X: 4, Y: 2, Colour: c},
		{X: 3, Y: 2, Colour: c},
		{X: 2, Y: 2, Colour: c},
		{X: 1, Y: 2, Colour: c},
		{X: 0, Y: 2, Colour: c},
		{X: 0, Y: 1, Colour: c},
	}, PixelsInSquare(0, 0, 4, 2, c, false))
}
