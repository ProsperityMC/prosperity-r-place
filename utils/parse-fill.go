package utils

import "errors"

var ErrInvalidFillOption = errors.New("invalid fill option")

func ParseFill(text string) (bool, error) {
	switch text {
	case "outline":
		return false, nil
	case "fill":
		return true, nil
	}
	return false, ErrInvalidFillOption
}
