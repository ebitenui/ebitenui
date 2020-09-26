package main

import (
	"github.com/blizzy78/ebitenui/image"
	"github.com/blizzy78/ebitenui/widget"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type images struct {
	button              *widget.ButtonImage
	stateButtonSelected *widget.ButtonImage
	sliderTrack         *widget.SliderTrackImage
	arrowDown           *widget.ButtonImageImage
	scrollContainer     *widget.ScrollContainerImage
	checkbox            *widget.CheckboxGraphicImage
	heart               *widget.ButtonImageImage
	toolTip             *image.NineSlice
}

func loadImages() (*images, error) {
	button, err := loadButtonImages(
		"graphics/button-idle.png",
		"graphics/button-hover.png",
		"graphics/button-pressed.png",
		"graphics/button-disabled.png",
		6, 12)
	if err != nil {
		return nil, err
	}

	stateButtonSelected, err := loadButtonImages(
		"graphics/button-pressed.png",
		"graphics/button-pressed.png",
		"graphics/button-pressed.png",
		"graphics/button-pressed-disabled.png",
		6, 12)
	if err != nil {
		return nil, err
	}

	arrowDown, err := loadGraphicImages(
		"graphics/arrow-down-idle.png",
		"graphics/arrow-down-disabled.png")
	if err != nil {
		return nil, err
	}

	checkboxUnchecked, err := loadGraphicImages(
		"graphics/checkbox-unchecked.png",
		"graphics/checkbox-unchecked.png")
	if err != nil {
		return nil, err
	}

	checkboxChecked, err := loadGraphicImages(
		"graphics/checkbox-checked-idle.png",
		"graphics/checkbox-checked-disabled.png")
	if err != nil {
		return nil, err
	}

	checkboxGreyed, err := loadGraphicImages(
		"graphics/checkbox-greyed-idle.png",
		"graphics/checkbox-greyed-disabled.png")
	if err != nil {
		return nil, err
	}

	heart, err := loadGraphicImages(
		"graphics/heart-idle.png",
		"graphics/heart-disabled.png")
	if err != nil {
		return nil, err
	}

	list, err := loadImageNineSlice("graphics/list.png", 6, 12)
	if err != nil {
		return nil, err
	}

	listMask, err := loadImageNineSlice("graphics/list-mask.png", 6, 12)
	if err != nil {
		return nil, err
	}

	toolTip, err := loadImageNineSlice("graphics/tooltip.png", 6, 12)
	if err != nil {
		return nil, err
	}

	return &images{
		button:              button,
		stateButtonSelected: stateButtonSelected,
		sliderTrack: &widget.SliderTrackImage{
			Idle:     list,
			Hover:    list,
			Disabled: list,
		},
		arrowDown: arrowDown,
		scrollContainer: &widget.ScrollContainerImage{
			Idle:     list,
			Disabled: list,
			Mask:     listMask,
		},
		checkbox: &widget.CheckboxGraphicImage{
			Unchecked: checkboxUnchecked,
			Checked:   checkboxChecked,
			Greyed:    checkboxGreyed,
		},
		heart:   heart,
		toolTip: toolTip,
	}, nil
}

func loadButtonImages(idle string, hover string, pressed string, disabled string, border int, center int) (*widget.ButtonImage, error) {
	idleImage, err := loadImageNineSlice(idle, border, center)
	if err != nil {
		return nil, err
	}

	hoverImage, err := loadImageNineSlice(hover, border, center)
	if err != nil {
		return nil, err
	}

	pressedImage, err := loadImageNineSlice(pressed, border, center)
	if err != nil {
		return nil, err
	}

	disabledImage, err := loadImageNineSlice(disabled, border, center)
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

func loadGraphicImages(idle string, disabled string) (*widget.ButtonImageImage, error) {
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

	return &widget.ButtonImageImage{
		Idle:     idleImage,
		Disabled: disabledImage,
	}, nil
}

func loadImageNineSlice(path string, border int, center int) (*image.NineSlice, error) {
	i, _, err := ebitenutil.NewImageFromFile(path, ebiten.FilterDefault)
	if err != nil {
		return nil, err
	}

	return image.NewNineSliceSimple(i, border, center), nil
}
