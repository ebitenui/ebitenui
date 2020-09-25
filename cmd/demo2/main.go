package main

import (
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

type pageContainer struct {
	widget    widget.HasWidget
	titleText *widget.Text
	flipBook  *widget.FlipBook
}

func main() {
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Ebiten UI Demo")
	ebiten.SetWindowResizable(true)

	ui, closeUI, err := createUI()
	if err != nil {
		log.Fatal(err)
	}

	defer closeUI()

	game := game{
		ui: ui,
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}

func createUI() (*ebitenui.UI, func(), error) {
	res, err := newResources()
	if err != nil {
		return nil, nil, err
	}

	toolTips := toolTipContents{
		tips: map[widget.HasWidget]string{},
		res:  res,
	}

	drag := newTextDragContents(res)

	rootContainer := widget.NewContainer(
		widget.ContainerOpts.WithLayout(widget.NewGridLayout(
			widget.GridLayoutOpts.WithColumns(1),
			widget.GridLayoutOpts.WithStretch([]bool{true}, []bool{false, true}),
			widget.GridLayoutOpts.WithPadding(widget.NewInsetsSimple(20)),
			widget.GridLayoutOpts.WithSpacing(0, 20))),
		widget.ContainerOpts.WithBackgroundImage(image.NewNineSliceColor(color.White)))

	dnd := widget.NewDragAndDrop(
		widget.DragAndDropOpts.WithContainer(rootContainer),
		widget.DragAndDropOpts.WithContentsCreater(drag),
	)

	rootContainer.AddChild(newInfoContainer(res))
	rootContainer.AddChild(newDemoContainer(res, &toolTips, dnd, drag))

	return &ebitenui.UI{
			Container: rootContainer,

			ToolTip: widget.NewToolTip(
				widget.ToolTipOpts.WithContainer(rootContainer),
				widget.ToolTipOpts.WithContentsCreater(&toolTips),
				widget.ToolTipOpts.WithUpdateEveryFrame(),
				widget.ToolTipOpts.WithNoSticky(),
				widget.ToolTipOpts.WithDelay(0)),

			DragAndDrop: dnd,
		},
		func() {
			res.close()
		},
		nil
}

func newResources() (*resources, error) {
	images, err := loadImages()
	if err != nil {
		return nil, err
	}

	fonts, err := loadFonts()
	if err != nil {
		return nil, err
	}

	colors := newColors()

	return &resources{
		images: images,
		fonts:  fonts,
		colors: colors,
	}, nil
}

func newInfoContainer(res *resources) widget.HasWidget {
	infoContainer := widget.NewContainer(
		widget.ContainerOpts.WithLayout(widget.NewRowLayout(
			widget.RowLayoutOpts.WithDirection(widget.DirectionVertical),
			widget.RowLayoutOpts.WithSpacing(0))))

	infoContainer.AddChild(widget.NewText(
		widget.TextOpts.WithText("Ebiten UI Demo", res.fonts.bigTitleFace, res.colors.textIdle)))

	infoContainer.AddChild(widget.NewText(
		widget.TextOpts.WithText("This program is a showcase of Ebiten UI widgets and layouts.", res.fonts.face, res.colors.textIdle)))

	return infoContainer
}

func newDemoContainer(res *resources, toolTips *toolTipContents, dnd *widget.DragAndDrop, drag *dragContents) widget.HasWidget {
	demoContainer := widget.NewContainer(
		widget.ContainerOpts.WithLayout(widget.NewGridLayout(
			widget.GridLayoutOpts.WithColumns(2),
			widget.GridLayoutOpts.WithStretch([]bool{false, true}, []bool{true}),
			widget.GridLayoutOpts.WithSpacing(20, 0),
		)))

	pages := []interface{}{
		buttonPage(res),
		checkboxPage(res),
		listPage(res),
		comboButtonPage(res),
		tabBookPage(res),
		gridLayoutPage(res),
		rowLayoutPage(res),
		sliderPage(res),
		toolTipPage(res, toolTips),
		dragAndDropPage(res, dnd, drag),
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
			widget.ScrollContainerOpts.WithImage(res.images.scrollContainer),
			widget.ScrollContainerOpts.WithPadding(widget.NewInsetsSimple(2))),
		widget.ListOpts.WithEntryColor(res.colors.list),
		widget.ListOpts.WithEntryFontFace(res.fonts.face),
		widget.ListOpts.WithSliderOpts(widget.SliderOpts.WithImages(res.images.sliderTrack, res.images.button)),
		widget.ListOpts.WithHideHorizontalSlider(),
		widget.ListOpts.WithHideVerticalSlider(),
		widget.ListOpts.WithControlWidgetSpacing(2),

		widget.ListOpts.WithEntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			pageContainer.setPage(args.Entry.(*page))
		}))
	demoContainer.AddChild(pageList)

	demoContainer.AddChild(pageContainer.widget)

	pageList.SetSelectedEntry(pages[0])

	return demoContainer
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
			widget.ListOpts.WithSliderOpts(
				widget.SliderOpts.WithImages(res.images.sliderTrack, res.images.button),
				widget.SliderOpts.WithTrackPadding(3)),
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
		widget.ListOpts.WithSliderOpts(
			widget.SliderOpts.WithImages(res.images.sliderTrack, res.images.button),
			widget.SliderOpts.WithTrackPadding(3)),
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

func (r *resources) close() {
	r.fonts.close()
}
