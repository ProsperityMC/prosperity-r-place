package utils

import "strconv"

func ParseIntArray(text []string) ([]int, error) {
	var err error
	out := make([]int, len(text))
	for i := range text {
		out[i], err = strconv.Atoi(text[i])
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}
