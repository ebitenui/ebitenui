package main

import (
	"github.com/blizzy78/ebitenui/image"
	"github.com/blizzy78/ebitenui/widget"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type images struct {
	button          *widget.ButtonImage
	buttonFlatLeft  *widget.ButtonImage
	buttonNoLeft    *widget.ButtonImage
	sliderTrack     *widget.SliderTrackImage
	arrowDown       *graphicImage
	scrollContainer *widget.ScrollContainerImage
}

type graphicImage struct {
	idle     *ebiten.Image
	disabled *ebiten.Image
}

func loadImages() (*images, error) {
	button, err := loadButtonImages(
		"graphics/button-2px-idle.png",
		"graphics/button-2px-hover.png",
		"graphics/button-2px-pressed.png",
		"graphics/button-2px-disabled.png",
		5, 6)
	if err != nil {
		return nil, err
	}

	buttonFlatLeft, err := loadButtonImages(
		"graphics/button-2px-flat-left-idle.png",
		"graphics/button-2px-flat-left-hover.png",
		"graphics/button-2px-flat-left-pressed.png",
		"graphics/button-2px-flat-left-disabled.png",
		5, 6)
	if err != nil {
		return nil, err
	}

	buttonNoLeft, err := loadButtonImages(
		"graphics/button-2px-no-left-idle.png",
		"graphics/button-2px-no-left-hover.png",
		"graphics/button-2px-no-left-pressed.png",
		"graphics/button-2px-no-left-disabled.png",
		5, 6)
	if err != nil {
		return nil, err
	}

	arrowDown, err := loadGraphicImages(
		"graphics/arrow-down-idle.png",
		"graphics/arrow-down-disabled.png")
	if err != nil {
		return nil, err
	}

	mask, err := loadImageNineSlice("graphics/mask.png", 5, 6)
	if err != nil {
		return nil, err
	}

	return &images{
		button:         button,
		buttonFlatLeft: buttonFlatLeft,
		buttonNoLeft:   buttonNoLeft,
		sliderTrack: &widget.SliderTrackImage{
			Idle:     button.Idle,
			Hover:    button.Hover,
			Disabled: button.Disabled,
		},
		arrowDown: arrowDown,
		scrollContainer: &widget.ScrollContainerImage{
			Idle:     button.Idle,
			Disabled: button.Disabled,
			Mask:     mask,
		},
	}, nil
}

func loadButtonImages(idle string, hover string, pressed string, disabled string, w int, h int) (*widget.ButtonImage, error) {
	idleImage, err := loadImageNineSlice(idle, w, h)
	if err != nil {
		return nil, err
	}

	hoverImage, err := loadImageNineSlice(hover, w, h)
	if err != nil {
		return nil, err
	}

	pressedImage, err := loadImageNineSlice(pressed, w, h)
	if err != nil {
		return nil, err
	}

	disabledImage, err := loadImageNineSlice(disabled, w, h)
	if err != nil {
		return nil, err
	}

	return &widget.ButtonImage{
		Idle:     idleImage,
		Hover:    hoverImage,
		Pressed:  pressedImage,
		Disabled: disabledImage,
	}, nil
}

func loadGraphicImages(idle string, disabled string) (*graphicImage, error) {
	idleImage, _, err := ebitenutil.NewImageFromFile(idle, ebiten.FilterDefault)
	if err != nil {
		return nil, err
	}

	var disabledImage *ebiten.Image
	if disabled != "" {
		var err error
		if disabledImage, _, err = ebitenutil.NewImageFromFile(disabled, ebiten.FilterDefault); err != nil {
			return nil, err
		}
	}

	return &graphicImage{
		idle:     idleImage,
		disabled: disabledImage,
	}, nil
}

func loadImageNineSlice(path string, w int, h int) (*image.NineSlice, error) {
	i, _, err := ebitenutil.NewImageFromFile(path, ebiten.FilterDefault)
	if err != nil {
		return nil, err
	}

	return image.NewNineSliceSimple(i, w, h), nil
}
