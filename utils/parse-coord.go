package utils

import (
	"errors"
	"image"
	"strconv"
	"strings"
)

var ErrCoordinateRequiresTwoNumbers = errors.New("coordinate requires two numbers")

func ParseCoordinate(text string) (point image.Point, err error) {
	textSplit := strings.Split(text, ",")
	if len(textSplit) != 2 {
		return image.Point{}, ErrCoordinateRequiresTwoNumbers
	}
	out := make([]int, 2)
	for i := range textSplit {
		out[i], err = strconv.Atoi(textSplit[i])
		if err != nil {
			return image.Point{}, err
		}
	}
	return image.Point{X: out[0], Y: out[1]}, nil
}
