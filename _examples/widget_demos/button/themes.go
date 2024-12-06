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
		DefaultFace:      &face,
		DefaultTextColor: color.Black,
		ButtonTheme: &widget.ButtonParams{
			TextColor: &widget.ButtonTextColor{
				Idle:    color.NRGBA{40, 40, 40, 255},
				Hover:   color.NRGBA{40, 40, 40, 255},
				Pressed: color.NRGBA{40, 40, 40, 255},
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
		LabelTheme: &widget.LabelParams{
			Face: &face,
			Color: &widget.LabelColor{
				Idle:     color.Black,
				Disabled: color.NRGBA{222, 222, 222, 255},
			},
			Insets: &widget.Insets{Top: 10},
		},
		TextTheme: &widget.TextParams{
			Face:   &face,
			Color:  color.NRGBA{0, 255, 0, 255},
			Insets: &widget.Insets{Top: 10},
		},
	}
}

func GetDarkTheme() *widget.Theme {
	// load button text font
	face, _ := loadFont(20)

	return &widget.Theme{
		DefaultFace:      &face,
		DefaultTextColor: color.White,
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
		LabelTheme: &widget.LabelParams{
			Face: &face,
			Color: &widget.LabelColor{
				Idle:     color.White,
				Disabled: color.NRGBA{222, 222, 222, 255},
			},
			Insets: &widget.Insets{Top: 5},
		},
		TextTheme: &widget.TextParams{
			Face:   &face,
			Color:  color.NRGBA{255, 0, 0, 255},
			Insets: &widget.Insets{Top: 5},
		},
	}
}

func loadButtonImage() (*widget.ButtonImage, error) {
	idle := image.NewNineSliceColor(color.NRGBA{R: 170, G: 170, B: 180, A: 255})

	hover := image.NewNineSliceColor(color.NRGBA{R: 150, G: 150, B: 170, A: 255})

	pressed := image.NewNineSliceColor(color.NRGBA{R: 120, G: 120, B: 140, A: 255})

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
