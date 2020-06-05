package main

import (
	"image/color"
	"strconv"

	"github.com/blizzy78/ebitenui/widget"
)

type colors struct {
	textIdle     color.Color
	textDisabled color.Color

	selectedBackground         color.Color
	selectedDisabledBackground color.Color

	list       *widget.ListEntryColor
	buttonText *widget.ButtonTextColor
}

func newColors() *colors {
	c := colors{
		textIdle:     hexToColor("282828"),
		textDisabled: hexToColor("808080"),

		selectedBackground:         hexToColor("a0a0a0"),
		selectedDisabledBackground: hexToColor("c0c0c0"),
	}

	c.list = &widget.ListEntryColor{
		Unselected:         c.textIdle,
		Selected:           c.textIdle,
		SelectedBackground: c.selectedBackground,

		DisabledUnselected:         c.textDisabled,
		DisabledSelected:           c.textDisabled,
		DisabledSelectedBackground: c.selectedDisabledBackground,
	}

	c.buttonText = &widget.ButtonTextColor{
		Idle:     c.textIdle,
		Disabled: c.textDisabled,
	}

	return &c
}

func hexToColor(h string) color.Color {
	u, err := strconv.ParseUint(h, 16, 0)
	if err != nil {
		panic(err)
	}

	return color.RGBA{
		R: uint8(u & 0xff0000 >> 16),
		G: uint8(u & 0xff00 >> 8),
		B: uint8(u & 0xff),
		A: 255,
	}
}
