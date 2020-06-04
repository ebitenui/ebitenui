package main

import (
	"image"
	"io/ioutil"
	"log"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"golang.org/x/image/font"

	"image/color"
	_ "image/png"

	"github.com/blizzy78/ebitenui"
	ebimage "github.com/blizzy78/ebitenui/image"
	"github.com/blizzy78/ebitenui/widget"
)

type game struct {
	ui *ebitenui.UI
}

type page struct {
	title   string
	content widget.HasWidget
}

func main() {
	ebiten.SetWindowSize(800, 500)
	ebiten.SetWindowTitle("Ebiten UI Demo")
	ebiten.SetWindowResizable(true)

	ui, fontFace, titleFontFace := createUI()

	defer func() {
		_ = fontFace.Close()
		_ = titleFontFace.Close()
	}()

	game := game{
		ui: ui,
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}

func createUI() (*ebitenui.UI, font.Face, font.Face) {
	images, err := loadImages()
	if err != nil {
		panic(err)
	}

	fontFace, titleFontFace, err := loadFonts()
	if err != nil {
		panic(err)
	}

	rootContainer := widget.NewContainer(
		widget.ContainerOpts.WithLayout(widget.NewGridLayout(
			widget.GridLayoutOpts.WithColumns(2),
			widget.GridLayoutOpts.WithStretch([]bool{false, true}, []bool{true}),
			widget.GridLayoutOpts.WithPadding(widget.NewInsetsSimple(20)),
			widget.GridLayoutOpts.WithSpacing(20, 0),
		)),
		widget.ContainerOpts.WithBackgroundImage(ebimage.NewNineSliceColor(color.White)))

	pages := []interface{}{
		buttonPage(images, fontFace),
		checkboxPage(images, fontFace),
	}

	var pageContainer *widget.Container
	var pageTitleText *widget.Text
	var pageFlipBook *widget.FlipBook

	pageList := widget.NewList(
		widget.ListOpts.WithEntries(pages),
		widget.ListOpts.WithEntryLabelFunc(func(e interface{}) string {
			return e.(*page).title
		}),
		widget.ListOpts.WithScrollContainerOpt(widget.ScrollContainerOpts.WithImage(images.scrollContainer)),
		widget.ListOpts.WithScrollContainerOpt(widget.ScrollContainerOpts.WithPadding(widget.NewInsetsSimple(2))),
		widget.ListOpts.WithEntryColor(&widget.ListEntryColor{
			Unselected:           color.Black,
			Selected:             color.Black,
			UnselectedBackground: color.White,
			SelectedBackground:   color.RGBA{128, 128, 128, 255},

			DisabledUnselected:           color.Black,
			DisabledSelected:             color.Black,
			DisabledUnselectedBackground: color.White,
			DisabledSelectedBackground:   color.RGBA{128, 128, 128, 255},
		}),
		widget.ListOpts.WithEntryFontFace(fontFace),
		widget.ListOpts.WithSliderOpt(widget.SliderOpts.WithImages(images.sliderTrack, images.button)),
		widget.ListOpts.WithHideHorizontalSlider(),
		widget.ListOpts.WithHideVerticalSlider(),
		widget.ListOpts.WithControlWidgetSpacing(2),

		widget.ListOpts.WithEntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			p := args.Entry.(*page)
			pageTitleText.Label = p.title
			pageFlipBook.SetPage(p.content)
			pageFlipBook.RequestRelayout()
		}))
	rootContainer.AddChild(pageList)

	pageContainer = widget.NewContainer(
		widget.ContainerOpts.WithLayout(widget.NewRowLayout(
			widget.RowLayoutOpts.WithDirection(widget.DirectionVertical),
			widget.RowLayoutOpts.WithSpacing(20))),
	)
	rootContainer.AddChild(pageContainer)

	pageTitleText = widget.NewText(
		widget.TextOpts.WithWidgetOpt(widget.WidgetOpts.WithLayoutData(&widget.RowLayoutData{
			Stretch: true,
		})),
		widget.TextOpts.WithText("", titleFontFace, color.Black))
	pageContainer.AddChild(pageTitleText)

	pageFlipBook = widget.NewFlipBook(
		widget.FlipBookOpts.WithContainerOpt(widget.ContainerOpts.WithWidgetOpt(widget.WidgetOpts.WithLayoutData(&widget.RowLayoutData{
			Stretch: true,
		}))))
	pageContainer.AddChild(pageFlipBook)

	pageList.SetSelectedEntry(pages[0])

	/*
		var button1 *widget.Button
		var button2 *widget.Button

		button1 = widget.NewButton(
			widget.ButtonOpts.WithLayoutData(&widget.RowLayoutData{
				Stretch: true,
			}),
			widget.ButtonOpts.WithImage(images.button),

			widget.ButtonOpts.WithClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				button2.GetWidget().Disabled = !button2.GetWidget().Disabled
			}))

		rootContainer.AddChild(button1)

		button2 = widget.NewButton(
			widget.ButtonOpts.WithLayoutData(&widget.RowLayoutData{
				Stretch: true,
			}),
			widget.ButtonOpts.WithImage(images.button),
			widget.ButtonOpts.WithText("foobar\nTy", fontFace, &widget.ButtonTextColor{
				Idle:     color.Black,
				Disabled: color.RGBA{128, 128, 128, 255},
			}),

			widget.ButtonOpts.WithClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				button1.GetWidget().Disabled = !button1.GetWidget().Disabled
			}))

		rootContainer.AddChild(button2)

		label := widget.NewText(
			widget.WithTextLayoutData(&widget.RowLayoutData{
				Stretch: true,
			}),
			widget.TextOpts.WithText("hallo", fontFace, color.Black))

		rootContainer.AddChild(label)

		rootContainer.AddChild(widget.NewButton(
			widget.ButtonOpts.WithLayoutData(&widget.RowLayoutData{
				Stretch: true,
			}),
			widget.ButtonOpts.WithImage(images.button),
			widget.ButtonOpts.WithText("bleh", fontFace, &widget.ButtonTextColor{
				Idle:     color.Black,
				Disabled: color.RGBA{128, 128, 128, 255},
			})))

		rootContainer.AddChild(widget.NewComboButton(
			widget.ComboButtonOpts.WithImage(images.button),
			widget.ComboButtonOpts.WithTextAndImage("Combo", fontFace, images.arrowDown, &widget.ButtonTextColor{
				Idle:     color.Black,
				Disabled: color.RGBA{128, 128, 128, 255},
			}),
			widget.ComboButtonOpts.WithContent(widget.NewButton(
				widget.ButtonOpts.WithImage(images.button),
				widget.ButtonOpts.WithText("foobar qux", fontFace, &widget.ButtonTextColor{
					Idle:     color.Black,
					Disabled: color.RGBA{128, 128, 128, 255},
				})))))

		entries := []interface{}{}
		for i := 1; i <= 20; i++ {
			entries = append(entries, i)
		}

		rootContainer.AddChild(widget.NewListComboButton(
			widget.ListComboButtonOpts.WithImage(images.button),
			widget.ListComboButtonOpts.WithText(fontFace, images.arrowDown, &widget.ButtonTextColor{
				Idle:     color.Black,
				Disabled: color.RGBA{128, 128, 128, 255},
			}),
			widget.ListComboButtonOpts.WithListImage(images.scrollContainer),
			widget.ListComboButtonOpts.WithListPadding(widget.NewInsetsSimple(2)),
			widget.ListComboButtonOpts.WithListSliderImages(images.sliderTrack, images.button),
			widget.ListComboButtonOpts.WithEntries(entries),
			widget.ListComboButtonOpts.WithEntryLabelFunc(
				func(e interface{}) string {
					return fmt.Sprintf("Entry %d", e.(int))
				},
				func(e interface{}) string {
					return fmt.Sprintf("Entry %d", e.(int))
				}),
			widget.ListComboButtonOpts.WithEntryFontFace(fontFace),
			widget.ListComboButtonOpts.WithEntryColor(&widget.ListEntryColor{
				Unselected:           color.Black,
				Selected:             color.Black,
				UnselectedBackground: color.White,
				SelectedBackground:   color.RGBA{128, 128, 128, 255},

				DisabledUnselected:           color.Black,
				DisabledSelected:             color.Black,
				DisabledUnselectedBackground: color.White,
				DisabledSelectedBackground:   color.RGBA{128, 128, 128, 255},
			}),

			widget.ListComboButtonOpts.WithEntrySelectedHandler(func(args *widget.ListComboButtonEntrySelectedEventArgs) {
				println("entry selected:", args.Entry.(int))
			}),
		))

		slider := widget.NewSlider(
			widget.SliderOpts.WithDirection(widget.DirectionHorizontal),
			widget.SliderOpts.WithLayoutData(&widget.RowLayoutData{
				Stretch:  true,
				MaxWidth: 250,
			}),
			widget.SliderOpts.WithImages(images.sliderTrack, images.button),

			widget.SliderOpts.WithChangedHandler(func(args *widget.SliderChangedEventArgs) {
				label.Label = fmt.Sprintf("%d", args.Current)
			}))

		rootContainer.AddChild(slider)

		rootContainer.AddChild(widget.NewButton(
			widget.ButtonOpts.WithLayoutData(&widget.RowLayoutData{
				Stretch:  true,
				MaxWidth: 250,
			}),
			widget.ButtonOpts.WithImage(images.button),
			widget.ButtonOpts.WithText("Disable", fontFace, &widget.ButtonTextColor{
				Idle:     color.Black,
				Disabled: color.RGBA{128, 128, 128, 255},
			}),

			widget.ButtonOpts.WithClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				slider.GetWidget().Disabled = !slider.GetWidget().Disabled

				if !slider.GetWidget().Disabled {
					slider.Current = 1
				}

				t := args.Button.Text()
				if slider.GetWidget().Disabled {
					t.Label = "Enable & Reset"
				} else {
					t.Label = "Disable"
				}
			})))

		rootContainer.AddChild(widget.NewList(
			widget.ListOpts.WithLayoutData(&widget.RowLayoutData{
				Stretch:   true,
				MaxHeight: 200,
			}),
			widget.ListOpts.WithEntries(entries),
			widget.ListOpts.WithEntryLabelFunc(func(e interface{}) string {
				return fmt.Sprintf("Entry %d", e.(int))
			}),
			widget.ListOpts.WithEntryColor(&widget.ListEntryColor{
				Unselected:           color.Black,
				Selected:             color.Black,
				UnselectedBackground: color.White,
				SelectedBackground:   color.RGBA{128, 128, 128, 255},

				DisabledUnselected:           color.Black,
				DisabledSelected:             color.Black,
				DisabledUnselectedBackground: color.White,
				DisabledSelectedBackground:   color.RGBA{128, 128, 128, 255},
			}),
			widget.ListOpts.WithEntryFontFace(fontFace),
			widget.ListOpts.WithImage(images.scrollContainer),
			widget.ListOpts.WithSliderImages(images.sliderTrack, images.button),
			widget.ListOpts.WithControlWidgetSpacing(2),
			widget.ListOpts.WithPadding(widget.NewInsetsSimple(2)),

			widget.ListOpts.WithEntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
				if args.Entry != args.PreviousEntry {
					println("entry selected: ", args.Entry.(int))
				}
			})))

		rootContainer.AddChild(widget.NewText(
			widget.TextOpts.WithText("test", fontFace, color.Black)))

		container := widget.NewContainer(
			widget.ContainerOpts.WithLayoutData(&widget.RowLayoutData{
				Stretch: true,
			}),
			widget.ContainerOpts.WithLayout(widget.NewGridLayout(
				widget.GridLayoutOpts.WithPadding(widget.NewInsetsSimple(10)),
				widget.GridLayoutOpts.WithColumns(3),
				widget.GridLayoutOpts.WithStretch([]bool{false, true, true}, nil),
				widget.GridLayoutOpts.WithSpacing(2, 2))),
			widget.ContainerOpts.WithBackgroundImage(ebimage.NewNineSliceColor(color.RGBA{0, 0, 255, 255})))

		for i := 0; i < 9; i++ {
			var c color.Color
			if i%2 == 0 {
				c = color.RGBA{255, 0, 0, 255}
			} else {
				c = color.RGBA{0, 255, 0, 255}
			}

			cont := widget.NewContainer(
				widget.ContainerOpts.WithBackgroundImage(ebimage.NewNineSliceColor(c)))

			if i%3 == 0 {
				cont.GetWidget().LayoutData = &widget.GridLayoutData{
					MaxWidth: 30,
				}
			}

			container.AddChild(cont)
		}

		rootContainer.AddChild(container)

		rootContainer.AddChild(widget.NewText(
			widget.TextOpts.WithText("test 2", fontFace, color.Black)))

		pageButtonsContainer := widget.NewContainer(
			widget.ContainerOpts.WithLayout(widget.NewRowLayout(
				widget.RowLayoutOpts.WithSpacing(10))))
		rootContainer.AddChild(pageButtonsContainer)

		flipBook := widget.NewFlipBook(
			widget.FlipBookOpts.WithLayoutData(&widget.RowLayoutData{
				Stretch: true,
			}))
		rootContainer.AddChild(flipBook)

		pages := []widget.HasWidget{}
		for i := 0; i < 5; i++ {
			c := widget.NewContainer(
				widget.ContainerOpts.WithLayout(widget.NewFillLayout()))
			c.AddChild(widget.NewText(
				widget.TextOpts.WithText(fmt.Sprintf("This is page %d", i+1), fontFace, color.Black)))
			pages = append(pages, c)
		}

		flipBook.SetPage(pages[0])

		for i := 0; i < 5; i++ {
			i := i
			pageButtonsContainer.AddChild(widget.NewButton(
				widget.ButtonOpts.WithImage(images.button),
				widget.ButtonOpts.WithText(fmt.Sprintf("Page %d", i+1), fontFace, &widget.ButtonTextColor{
					Idle:     color.Black,
					Disabled: color.RGBA{128, 128, 128, 255},
				}),

				widget.ButtonOpts.WithClickedHandler(func(args *widget.ButtonClickedEventArgs) {
					flipBook.SetPage(pages[i])
				})))
		}

		rootContainer.AddChild(widget.NewCheckbox(
			widget.CheckboxOpts.WithTriState(),
			widget.CheckboxOpts.WithImage(&widget.CheckboxImage{
				Button:  images.button,
				Graphic: images.checkbox,
			}),
		))
	*/

	return &ebitenui.UI{
		Container: rootContainer,
	}, fontFace, titleFontFace
}

func buttonPage(images *images, fontFace font.Face) *page {
	c := widget.NewContainer(
		widget.ContainerOpts.WithLayout(widget.NewRowLayout(
			widget.RowLayoutOpts.WithDirection(widget.DirectionVertical),
			widget.RowLayoutOpts.WithSpacing(10),
		)))

	b1 := widget.NewButton(
		widget.ButtonOpts.WithWidgetOpt(widget.WidgetOpts.WithLayoutData(&widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ButtonOpts.WithImage(images.button),
		widget.ButtonOpts.WithText("Button", fontFace, &widget.ButtonTextColor{
			Idle:     color.Black,
			Disabled: color.RGBA{128, 128, 128, 255},
		}))
	c.AddChild(b1)

	b2 := widget.NewButton(
		widget.ButtonOpts.WithWidgetOpt(widget.WidgetOpts.WithLayoutData(&widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ButtonOpts.WithImage(images.button),
		widget.ButtonOpts.WithTextAndImage("Button with Graphic", fontFace, images.heart, &widget.ButtonTextColor{
			Idle:     color.Black,
			Disabled: color.RGBA{128, 128, 128, 255},
		}))
	c.AddChild(b2)

	c.AddChild(newSeparator(&widget.RowLayoutData{
		Stretch: true,
	}))

	c.AddChild(widget.NewLabeledCheckbox(
		widget.LabeledCheckboxOpts.WithCheckboxOpt(widget.CheckboxOpts.WithButtonOpt(widget.ButtonOpts.WithImage(images.button))),
		widget.LabeledCheckboxOpts.WithCheckboxOpt(widget.CheckboxOpts.WithImage(images.checkbox)),
		widget.LabeledCheckboxOpts.WithTextOpt(widget.TextOpts.WithText("Disabled", fontFace, color.Black)),

		widget.LabeledCheckboxOpts.WithCheckboxOpt(widget.CheckboxOpts.WithChangedHandler(func(args *widget.CheckboxChangedEventArgs) {
			b1.GetWidget().Disabled = args.State == widget.CheckboxChecked
			b2.GetWidget().Disabled = args.State == widget.CheckboxChecked
		})),
	))

	return &page{
		title:   "Button",
		content: c,
	}
}

func checkboxPage(images *images, fontFace font.Face) *page {
	c := widget.NewContainer(
		widget.ContainerOpts.WithLayout(widget.NewRowLayout(
			widget.RowLayoutOpts.WithDirection(widget.DirectionVertical),
			widget.RowLayoutOpts.WithSpacing(10),
		)))

	cb := widget.NewLabeledCheckbox(
		widget.LabeledCheckboxOpts.WithCheckboxOpt(widget.CheckboxOpts.WithButtonOpt(widget.ButtonOpts.WithImage(images.button))),
		widget.LabeledCheckboxOpts.WithCheckboxOpt(widget.CheckboxOpts.WithImage(images.checkbox)),
		widget.LabeledCheckboxOpts.WithTextOpt(widget.TextOpts.WithText("Disabled", fontFace, color.Black)))
	c.AddChild(cb)

	c.AddChild(newSeparator(&widget.RowLayoutData{
		Stretch: true,
	}))

	c.AddChild(widget.NewLabeledCheckbox(
		widget.LabeledCheckboxOpts.WithCheckboxOpt(widget.CheckboxOpts.WithButtonOpt(widget.ButtonOpts.WithImage(images.button))),
		widget.LabeledCheckboxOpts.WithCheckboxOpt(widget.CheckboxOpts.WithImage(images.checkbox)),
		widget.LabeledCheckboxOpts.WithTextOpt(widget.TextOpts.WithText("Disabled", fontFace, color.Black)),

		widget.LabeledCheckboxOpts.WithCheckboxOpt(widget.CheckboxOpts.WithChangedHandler(func(args *widget.CheckboxChangedEventArgs) {
			cb.GetWidget().Disabled = args.State == widget.CheckboxChecked
		})),
	))

	return &page{
		title:   "Checkbox",
		content: c,
	}
}

func newSeparator(ld interface{}) widget.HasWidget {
	c := widget.NewContainer(
		widget.ContainerOpts.WithLayout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.WithDirection(widget.DirectionVertical),
				widget.RowLayoutOpts.WithPadding(widget.Insets{
					Top:    15,
					Bottom: 15,
				}))),
		widget.ContainerOpts.WithWidgetOpt(widget.WidgetOpts.WithLayoutData(ld)))

	c.AddChild(widget.NewGraphic(
		widget.GraphicOpts.WithWidgetOpt(widget.WidgetOpts.WithLayoutData(&widget.RowLayoutData{
			Stretch:   true,
			MaxHeight: 2,
		})),
		widget.GraphicOpts.WithImageNineSlice(ebimage.NewNineSliceColor(color.RGBA{192, 192, 192, 255})),
	))

	return c
}

func loadFonts() (font.Face, font.Face, error) {
	fontFace, err := loadFont("fonts/IBMPlexSans-Regular.ttf", 20)
	if err != nil {
		return nil, nil, err
	}

	titleFontFace, err := loadFont("fonts/IBMPlexSans-SemiBold.ttf", 26)
	if err != nil {
		return nil, nil, err
	}

	return fontFace, titleFontFace, nil
}

func loadFont(path string, size float64) (font.Face, error) {
	fontData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	ttfFont, err := truetype.Parse(fontData)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(ttfFont, &truetype.Options{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	}), nil
}

func (g *game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (g *game) Update(screen *ebiten.Image) error {
	g.ui.Update()
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	w, h := screen.Size()
	g.ui.Draw(screen, image.Rect(0, 0, w, h))
}
