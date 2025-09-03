package themes

import (
	img "image"
	"image/color"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/utilities/constantutil"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func getLightCheckbox() *widget.CheckboxImage {
	border := image.NewBorderedNineSliceColor(color.White, color.Black, 2)
	hover_border := image.NewBorderedNineSliceColor(color.NRGBA{214, 242, 255, 255}, color.Black, 2)
	idle := ebiten.NewImage(32, 32)
	border.Draw(idle, 32, 32, nil)
	idle9s := image.NewFixedNineSlice(idle)

	idle_hover := ebiten.NewImage(32, 32)
	hover_border.Draw(idle_hover, 32, 32, nil)
	idle_hover9s := image.NewFixedNineSlice(idle_hover)
	// Create checked image
	checked := ebiten.NewImage(32, 32)
	border.Draw(checked, 32, 32, nil)

	var checkmark vector.Path
	checkmark.MoveTo(28, 7)
	checkmark.LineTo(14, 24)
	checkmark.LineTo(6, 15)

	vertices := []ebiten.Vertex{}
	indices := []uint16{}
	vertices, indices = checkmark.AppendVerticesAndIndicesForStroke(vertices, indices, &vector.StrokeOptions{
		Width:    3,
		LineJoin: vector.LineJoinBevel,
	})
	blackImg := image.NewImageColor(color.Black)
	checked.DrawTriangles(vertices, indices, blackImg, &ebiten.DrawTrianglesOptions{AntiAlias: true})
	checked9s := image.NewFixedNineSlice(checked)

	checked_hover := ebiten.NewImage(32, 32)
	hover_border.Draw(checked_hover, 32, 32, nil)
	checked_hover.DrawTriangles(vertices, indices, blackImg, &ebiten.DrawTrianglesOptions{AntiAlias: true})
	checked_hover9s := image.NewFixedNineSlice(checked_hover)

	// Create greyed image
	greyed := ebiten.NewImage(32, 32)
	border.Draw(greyed, 32, 32, nil)
	vector.StrokeLine(greyed, 5, 16, 27, 16, 3, color.Black, true)
	greyed9s := image.NewFixedNineSlice(greyed)

	greyed_hover := ebiten.NewImage(32, 32)
	hover_border.Draw(greyed_hover, 32, 32, nil)
	vector.StrokeLine(greyed_hover, 5, 16, 27, 16, 3, color.Black, true)
	greyed_hover9s := image.NewFixedNineSlice(greyed_hover)

	return &widget.CheckboxImage{
		Unchecked:         idle9s,
		Checked:           checked9s,
		Greyed:            greyed9s,
		UncheckedHovered:  idle_hover9s,
		CheckedHovered:    checked_hover9s,
		GreyedHovered:     greyed_hover9s,
		UncheckedDisabled: idle9s,
		CheckedDisabled:   checked9s,
		GreyedDisabled:    greyed9s,
	}
}

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
		ListTheme: &widget.ListParams{
			EntryFace:                   &face,
			EntryTextPadding:            widget.NewInsetsSimple(5),
			EntryTextHorizontalPosition: constantutil.ConstantToPointer(widget.TextPositionStart),
			EntryTextVerticalPosition:   constantutil.ConstantToPointer(widget.TextPositionCenter),
			MinSize:                     &img.Point{150, 0},
			EntryColor: &widget.ListEntryColor{
				Unselected:         color.Black,
				Selected:           color.Black,
				DisabledUnselected: color.NRGBA{127, 122, 126, 255},
				DisabledSelected:   color.NRGBA{127, 122, 126, 255},

				SelectedBackground:        color.NRGBA{204, 232, 255, 255},
				SelectedFocusedBackground: color.NRGBA{214, 242, 255, 255},

				SelectingBackground:        color.NRGBA{229, 243, 255, 255},
				FocusedBackground:          color.NRGBA{229, 243, 255, 255},
				SelectingFocusedBackground: color.NRGBA{229, 243, 255, 255},
				DisabledSelectedBackground: color.NRGBA{229, 243, 255, 255},
			},
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
		ListComboButtonTheme: &widget.ListComboButtonParams{
			MaxContentHeight: constantutil.ConstantToPointer(200),
			List: &widget.ListParams{
				EntryFace:                   &face,
				EntryTextPadding:            widget.NewInsetsSimple(5),
				EntryTextHorizontalPosition: constantutil.ConstantToPointer(widget.TextPositionStart),
				EntryTextVerticalPosition:   constantutil.ConstantToPointer(widget.TextPositionCenter),
				MinSize:                     &img.Point{200, 0},
				EntryColor: &widget.ListEntryColor{
					Unselected:         color.Black,
					Selected:           color.Black,
					DisabledUnselected: color.NRGBA{127, 122, 126, 255},
					DisabledSelected:   color.NRGBA{127, 122, 126, 255},

					SelectedBackground:        color.NRGBA{204, 232, 255, 255},
					SelectedFocusedBackground: color.NRGBA{214, 242, 255, 255},

					SelectingBackground:        color.NRGBA{229, 243, 255, 255},
					FocusedBackground:          color.NRGBA{229, 243, 255, 255},
					SelectingFocusedBackground: color.NRGBA{229, 243, 255, 255},
					DisabledSelectedBackground: color.NRGBA{229, 243, 255, 255},
				},
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
			Button: &widget.ButtonParams{
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
				MinSize: &img.Point{200, 0},
			},
		},
		CheckboxTheme: &widget.CheckboxParams{
			Image: getLightCheckbox(),
			Label: &widget.LabelParams{
				Face: &face,
				Color: &widget.LabelColor{
					Idle:     color.Black,
					Disabled: color.NRGBA{122, 122, 122, 255},
				},
			},
		},
	}
}
