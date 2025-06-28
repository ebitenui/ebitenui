package themes

import (
	img "image"
	"image/color"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/utilities/constantutil"
	"github.com/ebitenui/ebitenui/widget"
)

func GetBasicDarkTheme() *widget.Theme {
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
			TextPosition: &widget.TextPositioning{
				VTextPosition: widget.TextPositionCenter,
				HTextPosition: widget.TextPositionCenter,
			},
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
			Padding: &widget.Insets{Top: 5},
		},
		TextTheme: &widget.TextParams{
			Face:    &face,
			Color:   color.NRGBA{255, 0, 0, 255},
			Padding: &widget.Insets{Top: 5},
		},
		TabbookTheme: &widget.TabBookParams{
			TabButton: &widget.ButtonParams{
				TextColor: &widget.ButtonTextColor{
					Idle:    color.NRGBA{40, 40, 40, 255},
					Hover:   color.NRGBA{40, 40, 40, 255},
					Pressed: color.NRGBA{40, 40, 40, 255},
				},
				TextFace: &face,
				Image: &widget.ButtonImage{
					Idle:    image.NewNineSliceColor(color.NRGBA{R: 170, G: 170, B: 180, A: 255}),
					Hover:   image.NewNineSliceColor(color.NRGBA{R: 150, G: 150, B: 170, A: 255}),
					Pressed: image.NewNineSliceColor(color.NRGBA{R: 120, G: 120, B: 140, A: 255}),
				},
				TextPadding: widget.NewInsetsSimple(5),
				MinSize:     &img.Point{98, 40},
			},
			TabSpacing: constantutil.ConstantToPointer(1),
		},
		TabTheme: &widget.TabParams{
			BackgroundImage: image.NewNineSliceColor(color.NRGBA{162, 158, 150, 255}),
		},
	}
}
