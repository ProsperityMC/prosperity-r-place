package shapes

import (
	"fmt"
	"image/color"
	"prosperity-r-place/utils"
)

func PixelsInSquare(x1, y1, x2, y2 int, colour color.RGBA, fill bool) []utils.Pixel {
	width := x2 - x1 + 1
	height := y2 - y1 + 1
	if fill {
		// #####
		// #####
		// #####
		// size of array width * height
		pixels := make([]utils.Pixel, width*height)
		for j := 0; j < height; j++ {
			for i := 0; i < width; i++ {
				pixels[j*width+i] = utils.Pixel{X: i + x1, Y: j + y1, Colour: colour}
			}
		}
		return pixels
	}
	// size of array
	// 11112
	// 4   2
	// 43333
	// total pixels = twice width and twice height
	// pixel order is clockwise from the top left
	pixels := make([]utils.Pixel, (width-1)*2+(height-1)*2)
	for i := 0; i < width-1; i++ {
		// order of items
		// 1234-
		// -   -
		// -8765
		pixels[i] = utils.Pixel{X: i + x1, Y: y1, Colour: colour}
		pixels[(width-1)+(height-1)+i] = utils.Pixel{X: (width - i - 1) + x1, Y: y2, Colour: colour}
	}
	for i := 0; i < height-1; i++ {
		// order of items
		// ----1
		// 4   2
		// 3----
		pixels[(width-1)+i] = utils.Pixel{X: x2, Y: i + y1, Colour: colour}
		pixels[(width-1)*2+(height-1)+i] = utils.Pixel{X: x1, Y: (height - i - 1) + y1, Colour: colour}
	}
	for _, i := range pixels {
		fmt.Printf("[A] %d, %d\n", i.X, i.Y)
	}
	return pixels
}
