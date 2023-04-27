package helpers

import (
	"errors"
	"strconv"
)

var (
	EMPTY_REGION = Rect{0, 0, 0, 0}
)

// Converts a CSV of floats to a Rect type.
func ParseRect(txt string) (Rect, error) {
	values, err := GetCSV(txt)
	if err != nil {
		return EMPTY_REGION, err
	}
	var region Rect
	if len(values) == 2 {
		width, err := strconv.ParseFloat(values[0], 64)
		if err != nil {
			return EMPTY_REGION, err
		}
		height, err := strconv.ParseFloat(values[1], 64)
		if err != nil {
			return EMPTY_REGION, err
		}
		region.X = 0
		region.Y = 0
		region.Width = width
		region.Height = height
		return region, nil
	} else if len(values) == 4 {
		x, err := strconv.ParseFloat(values[0], 64)
		if err != nil {
			return EMPTY_REGION, err
		}
		y, err := strconv.ParseFloat(values[1], 64)
		if err != nil {
			return EMPTY_REGION, err
		}
		width, err := strconv.ParseFloat(values[2], 64)
		if err != nil {
			return EMPTY_REGION, err
		}
		height, err := strconv.ParseFloat(values[3], 64)
		if err != nil {
			return EMPTY_REGION, err
		}
		region.X = x
		region.Y = y
		region.Width = width
		region.Height = height
		return region, nil
	}
	return EMPTY_REGION, errors.New("Invalid rect")
}
