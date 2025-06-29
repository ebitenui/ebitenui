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
				Idle:    color.Black,
				Hover:   color.Black,
				Pressed: color.Black,
			},
			TextFace: &face,
			Image: &widget.ButtonImage{
				Idle:    image.NewNineSliceColor(color.NRGBA{233, 231, 231, 255}),
				Hover:   image.NewNineSliceColor(color.NRGBA{223, 220, 220, 255}),
				Pressed: image.NewNineSliceColor(color.NRGBA{197, 192, 196, 255}),
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
			BackgroundImage: image.NewNineSliceColor(color.NRGBA{245, 246, 247, 255}),
		},
		LabelTheme: &widget.LabelParams{
			Face: &face,
			Color: &widget.LabelColor{
				Idle:     color.Black,
				Disabled: color.NRGBA{122, 122, 122, 255},
			},
		},
		TextTheme: &widget.TextParams{
			Face:  &face,
			Color: color.Black,
		},
		TabbookTheme: &widget.TabBookParams{
			TabButton: &widget.ButtonParams{
				TextColor: &widget.ButtonTextColor{
					Idle:    color.Black,
					Hover:   color.Black,
					Pressed: color.Black,
				},
				TextFace: &face,
				Image: &widget.ButtonImage{
					Idle:    image.NewNineSliceColor(color.NRGBA{233, 231, 231, 255}),
					Hover:   image.NewNineSliceColor(color.NRGBA{223, 220, 220, 255}),
					Pressed: image.NewNineSliceColor(color.NRGBA{197, 192, 196, 255}),
				},
				TextPadding: widget.NewInsetsSimple(5),
				MinSize:     &img.Point{98, 40},
			},
			TabSpacing: constantutil.ConstantToPointer(1),
		},
		TabTheme: &widget.TabParams{
			BackgroundImage: image.NewNineSliceColor(color.NRGBA{197, 197, 197, 255}),
		},
	}
}
