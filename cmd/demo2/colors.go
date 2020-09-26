package main

import (
	"image/color"
	"strconv"

	"github.com/blizzy78/ebitenui/widget"
)

type colors struct {
	background color.Color

	textIdle     color.Color
	textDisabled color.Color
	textToolTip  color.Color

	selectedBackground         color.Color
	selectedDisabledBackground color.Color

	list       *widget.ListEntryColor
	buttonText *widget.ButtonTextColor
	label      *widget.LabelColor
}

func newColors() *colors {
	c := colors{
		background: hexToColor("666666"),

		textIdle:     color.White,
		textDisabled: hexToColor("808080"),
		textToolTip:  color.Black,

		selectedBackground:         hexToColor("a0a0a0"),
		selectedDisabledBackground: hexToColor("707070"),
	}

	c.list = &widget.ListEntryColor{
		Unselected:         c.textIdle,
		Selected:           c.textIdle,
		SelectedBackground: c.selectedBackground,

		DisabledUnselected:         c.textDisabled,
		DisabledSelected:           hexToColor("d0d0d0"),
		DisabledSelectedBackground: c.selectedDisabledBackground,
	}

	c.buttonText = &widget.ButtonTextColor{
		Idle:     color.Black,
		Disabled: hexToColor("555555"),
	}

	c.label = &widget.LabelColor{
		Idle:     c.textIdle,
		Disabled: hexToColor("a0a0a0"),
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
