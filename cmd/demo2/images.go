package main

import (
	"github.com/blizzy78/ebitenui/image"
	"github.com/blizzy78/ebitenui/widget"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type images struct {
	button              *widget.ButtonImage
	buttonKenney        *widget.ButtonImage
	buttonFlatLeft      *widget.ButtonImage
	buttonNoLeft        *widget.ButtonImage
	stateButtonSelected *widget.ButtonImage
	sliderTrack         *widget.SliderTrackImage
	arrowDown           *widget.ButtonImageImage
	scrollContainer     *widget.ScrollContainerImage
	checkbox            *widget.CheckboxGraphicImage
	heart               *widget.ButtonImageImage
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

	buttonKenney, err := loadButtonImages(
		"graphics/button-kenney-idle.png",
		"graphics/button-kenney-idle.png",
		"graphics/button-kenney-pressed.png",
		"graphics/button-kenney-idle.png",
		9, 2)
	if err != nil {
		return nil, err
	}

	stateButtonSelected, err := loadButtonImages(
		"graphics/button-2px-pressed.png",
		"graphics/button-2px-pressed.png",
		"graphics/button-2px-pressed.png",
		"graphics/button-2px-disabled.png",
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

	containerIdle, err := loadImageNineSlice("graphics/container-2px-idle.png", 5, 6)
	if err != nil {
		return nil, err
	}

	containerDisabled, err := loadImageNineSlice("graphics/container-2px-disabled.png", 5, 6)
	if err != nil {
		return nil, err
	}

	return &images{
		button:              button,
		buttonFlatLeft:      buttonFlatLeft,
		buttonNoLeft:        buttonNoLeft,
		buttonKenney:        buttonKenney,
		stateButtonSelected: stateButtonSelected,
		sliderTrack: &widget.SliderTrackImage{
			Idle:     button.Idle,
			Hover:    button.Hover,
			Disabled: button.Disabled,
		},
		arrowDown: arrowDown,
		scrollContainer: &widget.ScrollContainerImage{
			Idle:     containerIdle,
			Disabled: containerDisabled,
			Mask:     mask,
		},
		checkbox: &widget.CheckboxGraphicImage{
			Unchecked: checkboxUnchecked,
			Checked:   checkboxChecked,
			Greyed:    checkboxGreyed,
		},
		heart: heart,
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
