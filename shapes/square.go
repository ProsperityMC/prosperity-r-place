package shapes

import (
	"image"
	"image/color"
	"prosperity-r-place/utils"
)

func PixelsInSquare(rect image.Rectangle, colour color.RGBA, fill bool) []utils.Pixel {
	// this correctly forms the rectangle points
	rect = rect.Canon()
	width := rect.Dx() + 1
	height := rect.Dy() + 1
	if fill {
		// #####
		// #####
		// #####
		// size of array width * height
		pixels := make([]utils.Pixel, width*height)
		for j := 0; j < height; j++ {
			for i := 0; i < width; i++ {
				pixels[j*width+i] = utils.Pixel{Point: rect.Min.Add(image.Point{X: i, Y: j}), Colour: colour}
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
		pixels[i] = utils.Pixel{Point: rect.Min.Add(image.Point{X: i}), Colour: colour}
		pixels[(width-1)+(height-1)+i] = utils.Pixel{Point: rect.Max.Add(image.Point{X: -i}), Colour: colour}
	}
	for i := 0; i < height-1; i++ {
		// order of items
		// ----1
		// 4   2
		// 3----
		pixels[(width-1)+i] = utils.Pixel{Point: image.Point{X: rect.Max.X, Y: i + rect.Min.Y}, Colour: colour}
		pixels[(width-1)*2+(height-1)+i] = utils.Pixel{Point: image.Point{X: rect.Min.X, Y: rect.Max.Y - i}, Colour: colour}
	}
	return pixels
}
