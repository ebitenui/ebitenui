package main

import (
	img "image"
	"log"

	"github.com/hajimehoshi/ebiten"

	"image/color"
	_ "image/png"

	"github.com/blizzy78/ebitenui"
	"github.com/blizzy78/ebitenui/image"
	"github.com/blizzy78/ebitenui/widget"
)

type game struct {
	ui *ebitenui.UI
}

type resources struct {
	images *images
	fonts  *fonts
	colors *colors
}

type page struct {
	title   string
	content widget.HasWidget
}

func main() {
	ebiten.SetWindowSize(800, 500)
	ebiten.SetWindowTitle("Ebiten UI Demo")
	ebiten.SetWindowResizable(true)

	ui, fonts := createUI()

	defer func() {
		fonts.close()
	}()

	game := game{
		ui: ui,
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}

func createUI() (*ebitenui.UI, *fonts) {
	images, err := loadImages()
	if err != nil {
		panic(err)
	}

	fonts, err := loadFonts()
	if err != nil {
		panic(err)
	}

	colors := newColors()

	res := &resources{
		images: images,
		fonts:  fonts,
		colors: colors,
	}

	rootContainer := widget.NewContainer(
		widget.ContainerOpts.WithLayout(widget.NewGridLayout(
			widget.GridLayoutOpts.WithColumns(1),
			widget.GridLayoutOpts.WithStretch([]bool{true}, []bool{false, true}),
			widget.GridLayoutOpts.WithPadding(widget.NewInsetsSimple(20)),
			widget.GridLayoutOpts.WithSpacing(0, 20))),
		widget.ContainerOpts.WithBackgroundImage(image.NewNineSliceColor(color.White)))

	rootContainer.AddChild(widget.NewText(
		widget.TextOpts.WithText("Ebiten UI Demo", fonts.bigTitleFace, res.colors.textIdle)))

	demoContainer := widget.NewContainer(
		widget.ContainerOpts.WithLayout(widget.NewGridLayout(
			widget.GridLayoutOpts.WithColumns(2),
			widget.GridLayoutOpts.WithStretch([]bool{false, true}, []bool{true}),
			widget.GridLayoutOpts.WithSpacing(20, 0),
		)))
	rootContainer.AddChild(demoContainer)

	pages := []interface{}{
		buttonPage(res),
		checkboxPage(res),
		listPage(res),
	}

	var pageContainer *widget.Container
	var pageTitleText *widget.Text
	var pageFlipBook *widget.FlipBook

	pageList := widget.NewList(
		widget.ListOpts.WithEntries(pages),
		widget.ListOpts.WithEntryLabelFunc(func(e interface{}) string {
			return e.(*page).title
		}),
		widget.ListOpts.WithScrollContainerOpts(
			widget.ScrollContainerOpts.WithImage(images.scrollContainer),
			widget.ScrollContainerOpts.WithPadding(widget.NewInsetsSimple(2))),
		widget.ListOpts.WithEntryColor(res.colors.list),
		widget.ListOpts.WithEntryFontFace(fonts.face),
		widget.ListOpts.WithSliderOpts(widget.SliderOpts.WithImages(images.sliderTrack, images.button)),
		widget.ListOpts.WithHideHorizontalSlider(),
		widget.ListOpts.WithHideVerticalSlider(),
		widget.ListOpts.WithControlWidgetSpacing(2),

		widget.ListOpts.WithEntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			p := args.Entry.(*page)
			pageTitleText.Label = p.title
			pageFlipBook.SetPage(p.content)
			pageFlipBook.RequestRelayout()
		}))
	demoContainer.AddChild(pageList)

	pageContainer = widget.NewContainer(
		widget.ContainerOpts.WithLayout(widget.NewRowLayout(
			widget.RowLayoutOpts.WithDirection(widget.DirectionVertical),
			widget.RowLayoutOpts.WithSpacing(15))),
	)
	demoContainer.AddChild(pageContainer)

	pageTitleText = widget.NewText(
		widget.TextOpts.WithWidgetOpts(widget.WidgetOpts.WithLayoutData(&widget.RowLayoutData{
			Stretch: true,
		})),
		widget.TextOpts.WithText("", fonts.titleFace, res.colors.textIdle))
	pageContainer.AddChild(pageTitleText)

	pageFlipBook = widget.NewFlipBook(
		widget.FlipBookOpts.WithContainerOpts(widget.ContainerOpts.WithWidgetOpts(widget.WidgetOpts.WithLayoutData(&widget.RowLayoutData{
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
	}, fonts
}

func buttonPage(res *resources) *page {
	c := widget.NewContainer(
		widget.ContainerOpts.WithLayout(widget.NewRowLayout(
			widget.RowLayoutOpts.WithDirection(widget.DirectionVertical),
			widget.RowLayoutOpts.WithSpacing(10),
		)))

	b1 := widget.NewButton(
		widget.ButtonOpts.WithWidgetOpts(widget.WidgetOpts.WithLayoutData(&widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ButtonOpts.WithImage(res.images.button),
		widget.ButtonOpts.WithText("Button", res.fonts.face, res.colors.buttonText))
	c.AddChild(b1)

	b2 := widget.NewButton(
		widget.ButtonOpts.WithWidgetOpts(widget.WidgetOpts.WithLayoutData(&widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ButtonOpts.WithImage(res.images.button),
		widget.ButtonOpts.WithTextAndImage("Button with Graphic", res.fonts.face, res.images.heart, res.colors.buttonText))
	c.AddChild(b2)

	c.AddChild(newSeparator(res, widget.WidgetOpts.WithLayoutData(&widget.RowLayoutData{
		Stretch: true,
	})))

	c.AddChild(widget.NewLabeledCheckbox(
		widget.LabeledCheckboxOpts.WithCheckboxOpts(
			widget.CheckboxOpts.WithButtonOpts(widget.ButtonOpts.WithImage(res.images.button)),
			widget.CheckboxOpts.WithImage(res.images.checkbox),

			widget.CheckboxOpts.WithChangedHandler(func(args *widget.CheckboxChangedEventArgs) {
				b1.GetWidget().Disabled = args.State == widget.CheckboxChecked
				b2.GetWidget().Disabled = args.State == widget.CheckboxChecked
			})),
		widget.LabeledCheckboxOpts.WithLabelOpts(widget.LabelOpts.WithText("Disabled", res.fonts.face, res.colors.label))))

	return &page{
		title:   "Button",
		content: c,
	}
}

func checkboxPage(res *resources) *page {
	c := widget.NewContainer(
		widget.ContainerOpts.WithLayout(widget.NewRowLayout(
			widget.RowLayoutOpts.WithDirection(widget.DirectionVertical),
			widget.RowLayoutOpts.WithSpacing(10),
		)))

	cb1 := widget.NewLabeledCheckbox(
		widget.LabeledCheckboxOpts.WithCheckboxOpts(
			widget.CheckboxOpts.WithButtonOpts(widget.ButtonOpts.WithImage(res.images.button)),
			widget.CheckboxOpts.WithImage(res.images.checkbox)),
		widget.LabeledCheckboxOpts.WithLabelOpts(widget.LabelOpts.WithText("Two-State Checkbox", res.fonts.face, res.colors.label)))
	c.AddChild(cb1)

	cb2 := widget.NewLabeledCheckbox(
		widget.LabeledCheckboxOpts.WithCheckboxOpts(
			widget.CheckboxOpts.WithButtonOpts(widget.ButtonOpts.WithImage(res.images.button)),
			widget.CheckboxOpts.WithImage(res.images.checkbox),
			widget.CheckboxOpts.WithTriState()),
		widget.LabeledCheckboxOpts.WithLabelOpts(widget.LabelOpts.WithText("Tri-State Checkbox", res.fonts.face, res.colors.label)))
	c.AddChild(cb2)

	c.AddChild(newSeparator(res, widget.WidgetOpts.WithLayoutData(&widget.RowLayoutData{
		Stretch: true,
	})))

	c.AddChild(widget.NewLabeledCheckbox(
		widget.LabeledCheckboxOpts.WithCheckboxOpts(
			widget.CheckboxOpts.WithButtonOpts(widget.ButtonOpts.WithImage(res.images.button)),
			widget.CheckboxOpts.WithImage(res.images.checkbox),

			widget.CheckboxOpts.WithChangedHandler(func(args *widget.CheckboxChangedEventArgs) {
				cb1.GetWidget().Disabled = args.State == widget.CheckboxChecked
				cb2.GetWidget().Disabled = args.State == widget.CheckboxChecked
			})),
		widget.LabeledCheckboxOpts.WithLabelOpts(widget.LabelOpts.WithText("Disabled", res.fonts.face, res.colors.label))))

	return &page{
		title:   "Checkbox",
		content: c,
	}
}

func listPage(res *resources) *page {
	c := widget.NewContainer(
		widget.ContainerOpts.WithLayout(widget.NewRowLayout(
			widget.RowLayoutOpts.WithDirection(widget.DirectionVertical),
			widget.RowLayoutOpts.WithSpacing(10),
		)))

	listsContainer := widget.NewContainer(
		widget.ContainerOpts.WithWidgetOpts(widget.WidgetOpts.WithLayoutData(&widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ContainerOpts.WithLayout(widget.NewGridLayout(
			widget.GridLayoutOpts.WithColumns(3),
			widget.GridLayoutOpts.WithStretch([]bool{true, false, true}, []bool{true}),
			widget.GridLayoutOpts.WithSpacing(10, 0))))
	c.AddChild(listsContainer)

	entries1 := []interface{}{"One", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten"}
	list1 := newList(entries1, res, widget.WidgetOpts.WithLayoutData(&widget.GridLayoutData{
		MaxHeight: 220,
	}))
	listsContainer.AddChild(list1)

	buttonsContainer := widget.NewContainer(
		widget.ContainerOpts.WithLayout(widget.NewRowLayout(
			widget.RowLayoutOpts.WithDirection(widget.DirectionVertical),
			widget.RowLayoutOpts.WithSpacing(10),
		)))
	listsContainer.AddChild(buttonsContainer)

	buttonsContainer.AddChild(widget.NewButton(
		widget.ButtonOpts.WithWidgetOpts(widget.WidgetOpts.WithLayoutData(&widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ButtonOpts.WithImage(res.images.button),
		widget.ButtonOpts.WithText("Add", res.fonts.face, res.colors.buttonText)))

	buttonsContainer.AddChild(widget.NewButton(
		widget.ButtonOpts.WithWidgetOpts(widget.WidgetOpts.WithLayoutData(&widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ButtonOpts.WithImage(res.images.button),
		widget.ButtonOpts.WithText("Remove", res.fonts.face, res.colors.buttonText)))

	entries2 := []interface{}{"Eleven", "Twelve", "Thirteen", "Fourteen", "Fifteen", "Sixteen", "Seventeen", "Eighteen", "Nineteen", "Twenty"}
	list2 := newList(entries2, res, widget.WidgetOpts.WithLayoutData(&widget.GridLayoutData{
		MaxHeight: 220,
	}))
	listsContainer.AddChild(list2)

	c.AddChild(newSeparator(res, widget.WidgetOpts.WithLayoutData(&widget.RowLayoutData{
		Stretch: true,
	})))

	c.AddChild(widget.NewLabeledCheckbox(
		widget.LabeledCheckboxOpts.WithCheckboxOpts(
			widget.CheckboxOpts.WithButtonOpts(widget.ButtonOpts.WithImage(res.images.button)),
			widget.CheckboxOpts.WithImage(res.images.checkbox),

			widget.CheckboxOpts.WithChangedHandler(func(args *widget.CheckboxChangedEventArgs) {
				list1.GetWidget().Disabled = args.State == widget.CheckboxChecked
				list2.GetWidget().Disabled = args.State == widget.CheckboxChecked
			})),
		widget.LabeledCheckboxOpts.WithLabelOpts(widget.LabelOpts.WithText("Disabled", res.fonts.face, res.colors.label))))

	return &page{
		title:   "List",
		content: c,
	}
}

func newList(entries []interface{}, res *resources, widgetOpts ...widget.WidgetOpt) *widget.List {
	return widget.NewList(
		widget.ListOpts.WithContainerOpts(widget.ContainerOpts.WithWidgetOpts(widgetOpts...)),
		widget.ListOpts.WithScrollContainerOpts(
			widget.ScrollContainerOpts.WithImage(res.images.scrollContainer),
			widget.ScrollContainerOpts.WithPadding(widget.NewInsetsSimple(2))),
		widget.ListOpts.WithSliderOpts(widget.SliderOpts.WithImages(res.images.sliderTrack, res.images.button)),
		widget.ListOpts.WithHideHorizontalSlider(),
		widget.ListOpts.WithControlWidgetSpacing(2),
		widget.ListOpts.WithEntries(entries),
		widget.ListOpts.WithEntryLabelFunc(func(e interface{}) string {
			return e.(string)
		}),
		widget.ListOpts.WithEntryFontFace(res.fonts.face),
		widget.ListOpts.WithEntryColor(res.colors.list),
	)
}

func newSeparator(res *resources, widgetOpts ...widget.WidgetOpt) widget.HasWidget {
	c := widget.NewContainer(
		widget.ContainerOpts.WithLayout(widget.NewRowLayout(
			widget.RowLayoutOpts.WithDirection(widget.DirectionVertical),
			widget.RowLayoutOpts.WithPadding(widget.Insets{
				Top:    15,
				Bottom: 15,
			}))),
		widget.ContainerOpts.WithWidgetOpts(widgetOpts...))

	c.AddChild(widget.NewGraphic(
		widget.GraphicOpts.WithWidgetOpts(widget.WidgetOpts.WithLayoutData(&widget.RowLayoutData{
			Stretch:   true,
			MaxHeight: 2,
		})),
		widget.GraphicOpts.WithImageNineSlice(image.NewNineSliceColor(res.colors.selectedDisabledBackground)),
	))

	return c
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
	g.ui.Draw(screen, img.Rect(0, 0, w, h))
}
