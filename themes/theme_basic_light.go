package themes

import (
	img "image"
	"image/color"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/utilities/constantutil"
	"github.com/ebitenui/ebitenui/widget"
)

func GetBasicLightTheme() *widget.Theme {
	// load button text font
	face, _ := loadFont(20)

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
			Image: &widget.ButtonImage{
				Idle:    image.NewNineSliceColor(color.NRGBA{R: 170, G: 170, B: 180, A: 255}),
				Hover:   image.NewNineSliceColor(color.NRGBA{R: 150, G: 150, B: 170, A: 255}),
				Pressed: image.NewNineSliceColor(color.NRGBA{R: 120, G: 120, B: 140, A: 255}),
			},
			TextPadding: &widget.Insets{
				Left:   30,
				Right:  30,
				Top:    5,
				Bottom: 5,
			},
			TextPosition: &widget.TextPositioning{
				VTextPosition: widget.TextPositionCenter,
				HTextPosition: widget.TextPositionCenter,
			},
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
			Padding: &widget.Insets{Top: 10},
		},
		TextTheme: &widget.TextParams{
			Face:    &face,
			Color:   color.NRGBA{0, 255, 0, 255},
			Padding: &widget.Insets{Left: 10, Top: 20},
		},
		/*
			widget.TabBookOpts.TabButtonImage(buttonImage),
			widget.TabBookOpts.TabButtonText(&face, &widget.ButtonTextColor{Idle: color.White, Disabled: color.White}),
			widget.TabBookOpts.TabButtonSpacing(5),
			widget.TabBookOpts.ContentPadding(widget.NewInsetsSimple(5)),
			widget.TabBookOpts.ContentSpacing(10),
			widget.TabBookOpts.TabButtonMinSize(&image.Point{98, 40}),

		*/
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
			TabSpacing:     constantutil.ConstantToPointer(1),
			ContentSpacing: constantutil.ConstantToPointer(5),
		},
	}
}
