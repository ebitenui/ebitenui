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
				Idle:    image.NewBorderedNineSliceColor(color.NRGBA{233, 231, 231, 255}, color.NRGBA{223, 220, 220, 255}, 2),
				Hover:   image.NewBorderedNineSliceColor(color.NRGBA{223, 220, 220, 255}, color.NRGBA{197, 192, 196, 255}, 2),
				Pressed: image.NewBorderedNineSliceColor(color.NRGBA{197, 192, 196, 255}, color.NRGBA{177, 172, 176, 255}, 2),
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
		TextInputTheme: &widget.TextInputParams{
			Face: &face,
			Image: &widget.TextInputImage{
				Idle:     image.NewBorderedNineSliceColor(color.White, color.NRGBA{177, 177, 177, 255}, 1),
				Disabled: image.NewBorderedNineSliceColor(color.NRGBA{223, 220, 220, 255}, color.NRGBA{177, 177, 177, 255}, 1),
			},
			Color: &widget.TextInputColor{
				Idle:          color.Black,
				Caret:         color.Black,
				Disabled:      color.NRGBA{122, 122, 122, 255},
				DisabledCaret: color.NRGBA{122, 122, 122, 255},
			},
			Padding: widget.NewInsetsSimple(5),
		},
		TextAreaTheme: &widget.TextAreaParams{
			Face:                   &face,
			StripBBCode:            constantutil.ConstantToPointer(true),
			ControlWidgetSpacing:   constantutil.ConstantToPointer(2),
			TextPadding:            &widget.Insets{Right: 18},
			ForegroundColor:        color.Black,
			ScrollContainerPadding: widget.NewInsetsSimple(4),
			ScrollContainerImage: &widget.ScrollContainerImage{
				Idle:     image.NewBorderedNineSliceColor(color.White, color.NRGBA{177, 177, 177, 255}, 1),
				Disabled: image.NewBorderedNineSliceColor(color.NRGBA{223, 220, 220, 255}, color.NRGBA{177, 177, 177, 255}, 1),
				Mask:     image.NewBorderedNineSliceColor(color.White, color.NRGBA{177, 177, 177, 255}, 1),
			},
			Slider: &widget.SliderParams{
				TrackImage: &widget.SliderTrackImage{
					Idle:     image.NewBorderedNineSliceColor(color.NRGBA{233, 231, 231, 255}, color.NRGBA{223, 220, 220, 255}, 2),
					Disabled: image.NewBorderedNineSliceColor(color.NRGBA{223, 220, 220, 255}, color.NRGBA{177, 177, 177, 255}, 1),
				},
				HandleImage: &widget.ButtonImage{
					Idle:    image.NewBorderedNineSliceColor(color.White, color.NRGBA{177, 177, 177, 255}, 1),
					Hover:   image.NewBorderedNineSliceColor(color.NRGBA{235, 235, 235, 255}, color.NRGBA{177, 177, 177, 255}, 2),
					Pressed: image.NewBorderedNineSliceColor(color.NRGBA{210, 210, 210, 255}, color.NRGBA{177, 177, 177, 255}, 2),
				},
			},
		},
		ProgressBarTheme: &widget.ProgressBarParams{
			TrackPadding: widget.NewInsetsSimple(2),
			TrackImage: &widget.ProgressBarImage{
				Idle:     image.NewBorderedNineSliceColor(color.White, color.NRGBA{177, 177, 177, 255}, 1),
				Hover:    image.NewBorderedNineSliceColor(color.White, color.NRGBA{177, 177, 177, 255}, 1),
				Disabled: image.NewBorderedNineSliceColor(color.NRGBA{223, 220, 220, 255}, color.NRGBA{177, 177, 177, 255}, 1),
			},
		},
		SliderTheme: &widget.SliderParams{
			TrackPadding:    widget.NewInsetsSimple(2),
			FixedHandleSize: constantutil.ConstantToPointer(6),
			TrackOffset:     constantutil.ConstantToPointer(0),
			PageSizeFunc: func() int {
				return 1
			},
			TrackImage: &widget.SliderTrackImage{
				Idle:     image.NewBorderedNineSliceColor(color.NRGBA{233, 231, 231, 255}, color.NRGBA{223, 220, 220, 255}, 2),
				Disabled: image.NewBorderedNineSliceColor(color.NRGBA{223, 220, 220, 255}, color.NRGBA{177, 177, 177, 255}, 1),
			},
			HandleImage: &widget.ButtonImage{
				Idle:    image.NewBorderedNineSliceColor(color.NRGBA{77, 77, 77, 255}, color.NRGBA{51, 51, 51, 255}, 2),
				Hover:   image.NewBorderedNineSliceColor(color.NRGBA{99, 99, 99, 255}, color.NRGBA{77, 77, 77, 255}, 2),
				Pressed: image.NewBorderedNineSliceColor(color.NRGBA{99, 99, 99, 255}, color.NRGBA{77, 77, 77, 255}, 2),
			},
		},
	}
}
