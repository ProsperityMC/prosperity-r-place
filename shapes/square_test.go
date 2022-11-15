package shapes

import (
	"github.com/stretchr/testify/assert"
	"image"
	"image/color"
	"prosperity-r-place/utils"
	"testing"
)

func TestPixelsInSquareFilled(t *testing.T) {
	c := color.RGBA{}
	assert.Equal(t, []utils.Pixel{
		{image.Point{}, c},
		{image.Point{X: 1}, c},
		{image.Point{X: 2}, c},
		{image.Point{X: 3}, c},
		{image.Point{X: 4}, c},
		{image.Point{Y: 1}, c},
		{image.Pt(1, 1), c},
		{image.Pt(2, 1), c},
		{image.Pt(3, 1), c},
		{image.Pt(4, 1), c},
		{image.Point{Y: 2}, c},
		{image.Pt(1, 2), c},
		{image.Pt(2, 2), c},
		{image.Pt(3, 2), c},
		{image.Pt(4, 2), c},
	}, PixelsInSquare(image.Rect(0, 0, 4, 2), c, true))
}

func TestPixelsInSquare(t *testing.T) {
	c := color.RGBA{}
	assert.Equal(t, []utils.Pixel{
		{image.Point{}, c},
		{image.Point{X: 1}, c},
		{image.Point{X: 2}, c},
		{image.Point{X: 3}, c},
		{image.Point{X: 4}, c},
		{image.Pt(4, 1), c},
		{image.Pt(4, 2), c},
		{image.Pt(3, 2), c},
		{image.Pt(2, 2), c},
		{image.Pt(1, 2), c},
		{image.Point{Y: 2}, c},
		{image.Point{Y: 1}, c},
	}, PixelsInSquare(image.Rect(0, 0, 4, 2), c, false))
}
