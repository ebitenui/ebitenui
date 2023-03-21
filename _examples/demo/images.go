package main

import (
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func newImageFromFile(path string) (*ebiten.Image, error) {
	f, err := embeddedAssets.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	i, _, err := ebitenutil.NewImageFromReader(f)
	return i, err
}

func loadGraphicImages(idle string, disabled string) (*widget.ButtonImageImage, error) {
	idleImage, err := newImageFromFile(idle)
	if err != nil {
		return nil, err
	}

	var disabledImage *ebiten.Image
	if disabled != "" {
		disabledImage, err = newImageFromFile(disabled)
		if err != nil {
			return nil, err
		}
	}

	return &widget.ButtonImageImage{
		Idle:     idleImage,
		Disabled: disabledImage,
	}, nil
}

func loadImageNineSlice(path string, centerWidth int, centerHeight int) (*image.NineSlice, error) {
	i, err := newImageFromFile(path)
	if err != nil {
		return nil, err
	}
	w := i.Bounds().Dx()
	h := i.Bounds().Dy()
	return image.NewNineSlice(i,
			[3]int{(w - centerWidth) / 2, centerWidth, w - (w-centerWidth)/2 - centerWidth},
			[3]int{(h - centerHeight) / 2, centerHeight, h - (h-centerHeight)/2 - centerHeight}),
		nil
}
