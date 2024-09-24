package main

import (
	"bytes"
	"image/color"
	"log"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

func GetLightTheme() *widget.Theme {
	// load button text font
	face, _ := loadFont(20)

	// load images for button states: idle, hover, and pressed
	buttonImage, _ := loadButtonImage()

	return &widget.Theme{
		ButtonTheme: &widget.ButtonParams{
			TextColor: &widget.ButtonTextColor{
				Idle:    color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
				Hover:   color.NRGBA{0, 255, 128, 255},
				Pressed: color.NRGBA{255, 0, 0, 255},
			},
			TextFace: &face,
			Image:    buttonImage,
			TextPadding: &widget.Insets{
				Left:   30,
				Right:  30,
				Top:    5,
				Bottom: 5,
			},
			HTextPosition: widget.TextPositionStart,
		},
		PanelTheme: &widget.PanelParams{
			BackgroundImage: image.NewNineSliceColor(color.NRGBA{212, 208, 200, 255}),
		},
	}
}

func GetDarkTheme() *widget.Theme {
	// load button text font
	face, _ := loadFont(20)

	return &widget.Theme{
		ButtonTheme: &widget.ButtonParams{
			TextColor: &widget.ButtonTextColor{
				Idle:    color.NRGBA{0x00, 0xf4, 0xff, 0xff},
				Hover:   color.NRGBA{0, 255, 128, 255},
				Pressed: color.NRGBA{255, 0, 0, 255},
			},
			TextFace: &face,
			Image: &widget.ButtonImage{
				Idle:    image.NewNineSliceColor(color.NRGBA{R: 100, G: 50, B: 50, A: 255}),
				Hover:   image.NewNineSliceColor(color.NRGBA{R: 100, G: 30, B: 30, A: 255}),
				Pressed: image.NewNineSliceColor(color.NRGBA{R: 80, G: 30, B: 30, A: 255}),
			},
			TextPadding: &widget.Insets{
				Left:   60,
				Right:  60,
				Top:    5,
				Bottom: 5,
			},
			HTextPosition: widget.TextPositionStart,
		},
		PanelTheme: &widget.PanelParams{
			BackgroundImage: image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff}),
		},
	}
}

func loadButtonImage() (*widget.ButtonImage, error) {
	idle := image.NewNineSliceColor(color.NRGBA{R: 170, G: 170, B: 180, A: 255})

	hover := image.NewNineSliceColor(color.NRGBA{R: 130, G: 130, B: 150, A: 255})

	pressed := image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 120, A: 255})

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}, nil
}

func loadFont(size float64) (text.Face, error) {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &text.GoTextFace{
		Source: s,
		Size:   size,
	}, nil
}
