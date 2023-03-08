package colorutil

import (
	"image/color"
	"strconv"
)

func HexToColor(h string) (color.Color, error) {
	u, err := strconv.ParseUint(h, 16, 0)
	if err != nil {
		return nil, err
	}

	return color.NRGBA{
		R: uint8(u & 0xff0000 >> 16),
		G: uint8(u & 0xff00 >> 8),
		B: uint8(u & 0xff),
		A: uint8(255),
	}, nil
}
