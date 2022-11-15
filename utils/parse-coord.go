package utils

import (
	"errors"
	"strconv"
	"strings"
)

var ErrCoordinateRequiresTwoNumbers = errors.New("coordinate requires two numbers")

func ParseCoordinate(text string) (x, y int, err error) {
	textSplit := strings.Split(text, ",")
	if len(textSplit) != 2 {
		return 0, 0, ErrCoordinateRequiresTwoNumbers
	}
	out := make([]int, 2)
	for i := range text {
		out[i], err = strconv.Atoi(textSplit[i])
		if err != nil {
			return 0, 0, err
		}
	}
	return out[0], out[1], nil
}
