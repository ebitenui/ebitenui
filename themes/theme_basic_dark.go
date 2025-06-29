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
				Idle:    image.NewBorderedNineSliceColor(color.NRGBA{51, 51, 51, 255}, color.NRGBA{81, 81, 81, 255}, 2),
				Hover:   image.NewBorderedNineSliceColor(color.NRGBA{77, 77, 77, 255}, color.NRGBA{51, 51, 51, 255}, 2),
				Pressed: image.NewBorderedNineSliceColor(color.NRGBA{119, 119, 119, 255}, color.NRGBA{77, 77, 77, 255}, 2),
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
		TextInputTheme: &widget.TextInputParams{
			Face: &face,
			Image: &widget.TextInputImage{
				Idle:     image.NewBorderedNineSliceColor(color.NRGBA{77, 77, 77, 255}, color.NRGBA{177, 177, 177, 255}, 1),
				Disabled: image.NewBorderedNineSliceColor(color.NRGBA{47, 47, 47, 255}, color.NRGBA{177, 177, 177, 255}, 1),
			},
			Color: &widget.TextInputColor{
				Idle:          color.White,
				Caret:         color.White,
				Disabled:      color.NRGBA{127, 122, 126, 255},
				DisabledCaret: color.NRGBA{127, 122, 126, 255},
			},
			Padding: widget.NewInsetsSimple(5),
		},
		TextAreaTheme: &widget.TextAreaParams{
			Face:                   &face,
			StripBBCode:            constantutil.ConstantToPointer(true),
			ControlWidgetSpacing:   constantutil.ConstantToPointer(2),
			TextPadding:            &widget.Insets{Right: 18},
			ForegroundColor:        color.White,
			ScrollContainerPadding: widget.NewInsetsSimple(4),
			ScrollContainerImage: &widget.ScrollContainerImage{
				Idle:     image.NewBorderedNineSliceColor(color.NRGBA{77, 77, 77, 255}, color.NRGBA{177, 177, 177, 255}, 1),
				Disabled: image.NewBorderedNineSliceColor(color.NRGBA{47, 47, 47, 255}, color.NRGBA{177, 177, 177, 255}, 1),
				Mask:     image.NewBorderedNineSliceColor(color.NRGBA{77, 77, 77, 255}, color.NRGBA{177, 177, 177, 255}, 1),
			},
			Slider: &widget.SliderParams{
				TrackImage: &widget.SliderTrackImage{
					Idle:     image.NewBorderedNineSliceColor(color.NRGBA{77, 77, 77, 255}, color.NRGBA{177, 177, 177, 255}, 1),
					Disabled: image.NewBorderedNineSliceColor(color.NRGBA{47, 47, 47, 255}, color.NRGBA{177, 177, 177, 255}, 1),
				},
				HandleImage: &widget.ButtonImage{
					Idle:    image.NewBorderedNineSliceColor(color.NRGBA{77, 77, 77, 255}, color.NRGBA{51, 51, 51, 255}, 2),
					Hover:   image.NewBorderedNineSliceColor(color.NRGBA{99, 99, 99, 255}, color.NRGBA{77, 77, 77, 255}, 2),
					Pressed: image.NewBorderedNineSliceColor(color.NRGBA{99, 99, 99, 255}, color.NRGBA{77, 77, 77, 255}, 2),
				},
			},
		},
	}
}
