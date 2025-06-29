package themes

import (
	img "image"
	"image/color"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/utilities/constantutil"
	"github.com/ebitenui/ebitenui/widget"
	"golang.org/x/image/colornames"
)

func GetBasicDarkTheme() *widget.Theme {
	// load button text font
	face, _ := loadFont(20)

	return &widget.Theme{
		DefaultFace:      &face,
		DefaultTextColor: color.White,
		ButtonTheme: &widget.ButtonParams{
			TextColor: &widget.ButtonTextColor{
				Idle:    color.White,
				Hover:   color.White,
				Pressed: color.White,
			},
			TextFace: &face,
			Image: &widget.ButtonImage{
				Idle:    image.NewNineSliceColor(color.NRGBA{51, 51, 51, 255}),
				Hover:   image.NewNineSliceColor(color.NRGBA{77, 77, 77, 255}),
				Pressed: image.NewNineSliceColor(color.NRGBA{119, 119, 119, 255}),
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
			BackgroundImage: image.NewNineSliceColor(colornames.Black),
		},
		LabelTheme: &widget.LabelParams{
			Face: &face,
			Color: &widget.LabelColor{
				Idle:     color.White,
				Disabled: color.NRGBA{122, 122, 122, 255},
			},
		},
		TextTheme: &widget.TextParams{
			Face:  &face,
			Color: color.White,
		},
		TabbookTheme: &widget.TabBookParams{
			TabButton: &widget.ButtonParams{
				TextColor: &widget.ButtonTextColor{
					Idle:    color.White,
					Hover:   color.White,
					Pressed: color.White,
				},
				TextFace: &face,
				Image: &widget.ButtonImage{
					Idle:    image.NewNineSliceColor(color.NRGBA{51, 51, 51, 255}),
					Hover:   image.NewNineSliceColor(color.NRGBA{77, 77, 77, 255}),
					Pressed: image.NewNineSliceColor(color.NRGBA{119, 119, 119, 255}),
				},
				TextPadding: widget.NewInsetsSimple(5),
				MinSize:     &img.Point{98, 40},
			},
			TabSpacing: constantutil.ConstantToPointer(1),
		},
		TabTheme: &widget.TabParams{
			BackgroundImage: image.NewNineSliceColor(color.NRGBA{32, 32, 32, 255}),
		},
	}
}
