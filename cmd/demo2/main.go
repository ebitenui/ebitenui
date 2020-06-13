package main

import (
	"fmt"
	img "image"
	"log"
	"sort"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"

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

type pageContainer struct {
	widget    widget.HasWidget
	titleText *widget.Text
	flipBook  *widget.FlipBook
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
		comboButtonPage(res),
		tabBookPage(res),
	}

	collator := collate.New(language.English)
	sort.Slice(pages, func(a int, b int) bool {
		p1 := pages[a].(*page)
		p2 := pages[b].(*page)
		return collator.CompareString(p1.title, p2.title) < 0
	})

	pageContainer := newPageContainer(res)

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
			pageContainer.setPage(args.Entry.(*page))
		}))
	demoContainer.AddChild(pageList)

	demoContainer.AddChild(pageContainer.widget)

	pageList.SetSelectedEntry(pages[0])

	return &ebitenui.UI{
		Container: rootContainer,
	}, fonts
}

func newPageContainer(res *resources) *pageContainer {
	c := widget.NewContainer(
		widget.ContainerOpts.WithLayout(widget.NewRowLayout(
			widget.RowLayoutOpts.WithDirection(widget.DirectionVertical),
			widget.RowLayoutOpts.WithSpacing(15))),
	)

	titleText := widget.NewText(
		widget.TextOpts.WithWidgetOpts(widget.WidgetOpts.WithLayoutData(&widget.RowLayoutData{
			Stretch: true,
		})),
		widget.TextOpts.WithText("", res.fonts.titleFace, res.colors.textIdle))
	c.AddChild(titleText)

	flipBook := widget.NewFlipBook(
		widget.FlipBookOpts.WithContainerOpts(widget.ContainerOpts.WithWidgetOpts(widget.WidgetOpts.WithLayoutData(&widget.RowLayoutData{
			Stretch: true,
		}))))
	c.AddChild(flipBook)

	return &pageContainer{
		widget:    c,
		titleText: titleText,
		flipBook:  flipBook,
	}
}

func (p *pageContainer) setPage(page *page) {
	p.titleText.Label = page.title
	p.flipBook.SetPage(page.content)
	p.flipBook.RequestRelayout()
}

func buttonPage(res *resources) *page {
	c := newPageContentContainer()

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

	c.AddChild(widget.NewButton(
		widget.ButtonOpts.WithWidgetOpts(widget.WidgetOpts.WithLayoutData(&widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ButtonOpts.WithImage(res.images.buttonKenney),
		widget.ButtonOpts.WithText("Button", res.fonts.face, res.colors.buttonText),
		widget.ButtonOpts.WithTextPadding(widget.Insets{
			Left:   10,
			Right:  10,
			Top:    6,
			Bottom: 10,
		})))

	c.AddChild(newSeparator(res, &widget.RowLayoutData{
		Stretch: true,
	}))

	c.AddChild(newCheckbox("Disabled", func(args *widget.CheckboxChangedEventArgs) {
		b1.GetWidget().Disabled = args.State == widget.CheckboxChecked
		b2.GetWidget().Disabled = args.State == widget.CheckboxChecked
	}, res))

	return &page{
		title:   "Button",
		content: c,
	}
}

func checkboxPage(res *resources) *page {
	c := newPageContentContainer()

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

	c.AddChild(newSeparator(res, &widget.RowLayoutData{
		Stretch: true,
	}))

	c.AddChild(newCheckbox("Disabled", func(args *widget.CheckboxChangedEventArgs) {
		cb1.GetWidget().Disabled = args.State == widget.CheckboxChecked
		cb2.GetWidget().Disabled = args.State == widget.CheckboxChecked
	}, res))

	return &page{
		title:   "Checkbox",
		content: c,
	}
}

func listPage(res *resources) *page {
	c := newPageContentContainer()

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

	c.AddChild(newSeparator(res, &widget.RowLayoutData{
		Stretch: true,
	}))

	c.AddChild(newCheckbox("Disabled", func(args *widget.CheckboxChangedEventArgs) {
		list1.GetWidget().Disabled = args.State == widget.CheckboxChecked
		list2.GetWidget().Disabled = args.State == widget.CheckboxChecked
	}, res))

	return &page{
		title:   "List",
		content: c,
	}
}

func comboButtonPage(res *resources) *page {
	c := newPageContentContainer()

	entries := []interface{}{}
	for i := 1; i <= 20; i++ {
		entries = append(entries, i)
	}

	cb := newListComboButton(
		entries,
		func(e interface{}) string {
			return fmt.Sprintf("Entry %d", e.(int))
		},
		func(e interface{}) string {
			return fmt.Sprintf("Entry %d", e.(int))
		},
		func(args *widget.ListComboButtonEntrySelectedEventArgs) {
			c.RequestRelayout()
		},
		res)
	c.AddChild(cb)

	c.AddChild(newSeparator(res, &widget.RowLayoutData{
		Stretch: true,
	}))

	c.AddChild(newCheckbox("Disabled", func(args *widget.CheckboxChangedEventArgs) {
		cb.GetWidget().Disabled = args.State == widget.CheckboxChecked
	}, res))

	return &page{
		title:   "Combo Button",
		content: c,
	}
}

func tabBookPage(res *resources) *page {
	c := newPageContentContainer()

	tabs := []*widget.TabBookTab{}

	for i := 0; i < 5; i++ {
		tc := widget.NewContainer(
			widget.ContainerOpts.WithLayout(widget.NewRowLayout(
				widget.RowLayoutOpts.WithDirection(widget.DirectionVertical),
				widget.RowLayoutOpts.WithSpacing(10))),
			widget.ContainerOpts.WithAutoDisableChildren())

		for j := 0; j < 3; j++ {
			b := widget.NewButton(
				widget.ButtonOpts.WithImage(res.images.button),
				widget.ButtonOpts.WithText(fmt.Sprintf("Button %d on Tab %d", j+1, i+1), res.fonts.face, res.colors.buttonText))
			tc.AddChild(b)
		}

		tab := widget.NewTabBookTab(fmt.Sprintf("Tab %d", i+1), tc)
		if i == 2 {
			tab.Disabled = true
		}

		tabs = append(tabs, tab)
	}

	t := widget.NewTabBook(
		widget.TabBookOpts.WithTabs(tabs...),
		widget.TabBookOpts.WithTabButtonImage(res.images.button, res.images.stateButtonSelected),
		widget.TabBookOpts.WithTabButtonText(res.fonts.face, res.colors.buttonText),
		widget.TabBookOpts.WithTabButtonSpacing(4),
		widget.TabBookOpts.WithSpacing(10))
	c.AddChild(t)

	c.AddChild(newSeparator(res, &widget.RowLayoutData{
		Stretch: true,
	}))

	c.AddChild(newCheckbox("Disabled", func(args *widget.CheckboxChangedEventArgs) {
		t.GetWidget().Disabled = args.State == widget.CheckboxChecked
	}, res))

	return &page{
		title:   "Tab Book",
		content: c,
	}
}

func newCheckbox(label string, changedHandler widget.CheckboxChangedHandlerFunc, res *resources) *widget.LabeledCheckbox {
	return widget.NewLabeledCheckbox(
		widget.LabeledCheckboxOpts.WithCheckboxOpts(
			widget.CheckboxOpts.WithButtonOpts(widget.ButtonOpts.WithImage(res.images.button)),
			widget.CheckboxOpts.WithImage(res.images.checkbox),

			widget.CheckboxOpts.WithChangedHandler(changedHandler)),
		widget.LabeledCheckboxOpts.WithLabelOpts(widget.LabelOpts.WithText(label, res.fonts.face, res.colors.label)))
}

func newPageContentContainer() *widget.Container {
	return widget.NewContainer(
		widget.ContainerOpts.WithLayout(widget.NewRowLayout(
			widget.RowLayoutOpts.WithDirection(widget.DirectionVertical),
			widget.RowLayoutOpts.WithSpacing(10),
		)))
}

func newListComboButton(entries []interface{}, buttonLabel widget.SelectComboButtonEntryLabelFunc, entryLabel widget.ListEntryLabelFunc,
	entrySelectedHandler widget.ListComboButtonEntrySelectedHandlerFunc, res *resources) *widget.ListComboButton {

	return widget.NewListComboButton(
		widget.ListComboButtonOpts.WithSelectComboButtonOpts(widget.SelectComboButtonOpts.WithComboButtonOpts(widget.ComboButtonOpts.WithButtonOpts(widget.ButtonOpts.WithImage(res.images.button)))),
		widget.ListComboButtonOpts.WithText(res.fonts.face, res.images.arrowDown, res.colors.buttonText),
		widget.ListComboButtonOpts.WithListOpts(
			widget.ListOpts.WithEntries(entries),
			widget.ListOpts.WithScrollContainerOpts(
				widget.ScrollContainerOpts.WithImage(res.images.scrollContainer),
				widget.ScrollContainerOpts.WithPadding(widget.NewInsetsSimple(2))),
			widget.ListOpts.WithSliderOpts(widget.SliderOpts.WithImages(res.images.sliderTrack, res.images.button)),
			widget.ListOpts.WithEntryFontFace(res.fonts.face),
			widget.ListOpts.WithEntryColor(res.colors.list)),
		widget.ListComboButtonOpts.WithEntryLabelFunc(buttonLabel, entryLabel),
		widget.ListComboButtonOpts.WithEntrySelectedHandler(entrySelectedHandler))
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

func newSeparator(res *resources, ld interface{}) widget.HasWidget {
	c := widget.NewContainer(
		widget.ContainerOpts.WithLayout(widget.NewRowLayout(
			widget.RowLayoutOpts.WithDirection(widget.DirectionVertical),
			widget.RowLayoutOpts.WithPadding(widget.Insets{
				Top:    15,
				Bottom: 15,
			}))),
		widget.ContainerOpts.WithWidgetOpts(widget.WidgetOpts.WithLayoutData(ld)))

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
