package main

import (
	img "image"
	"log"
	"sort"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"

	"github.com/hajimehoshi/ebiten"

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
		log.Print(err)
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
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, true, false}),
			widget.GridLayoutOpts.Padding(widget.NewInsetsSimple(20)),
			widget.GridLayoutOpts.Spacing(0, 20))),
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(res.colors.background)))

	dnd := widget.NewDragAndDrop(
		widget.DragAndDropOpts.Container(rootContainer),
		widget.DragAndDropOpts.ContentsCreater(drag),
	)

	rootContainer.AddChild(newInfoContainer(res))
	rootContainer.AddChild(newDemoContainer(res, &toolTips, dnd, drag))

	rootContainer.AddChild(widget.NewText(
		widget.TextOpts.Text("github.com/blizzy78/ebitenui", res.fonts.toolTipFace, res.colors.textIdle)))

	return &ebitenui.UI{
			Container: rootContainer,

			ToolTip: widget.NewToolTip(
				widget.ToolTipOpts.Container(rootContainer),
				widget.ToolTipOpts.ContentsCreater(&toolTips),
			),

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
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(0))))

	infoContainer.AddChild(widget.NewText(
		widget.TextOpts.Text("Ebiten UI Demo", res.fonts.bigTitleFace, res.colors.textIdle)))

	infoContainer.AddChild(widget.NewText(
		widget.TextOpts.Text("This program is a showcase of Ebiten UI widgets and layouts.", res.fonts.face, res.colors.textIdle)))

	return infoContainer
}

func newDemoContainer(res *resources, toolTips *toolTipContents, dnd *widget.DragAndDrop, drag *dragContents) widget.HasWidget {
	demoContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Stretch([]bool{false, true}, []bool{true}),
			widget.GridLayoutOpts.Spacing(20, 0),
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
		widget.ListOpts.Entries(pages),
		widget.ListOpts.EntryLabelFunc(func(e interface{}) string {
			return e.(*page).title
		}),
		widget.ListOpts.ScrollContainerOpts(
			widget.ScrollContainerOpts.Image(res.images.scrollContainer),
			widget.ScrollContainerOpts.Padding(widget.NewInsetsSimple(2))),
		widget.ListOpts.EntryColor(res.colors.list),
		widget.ListOpts.EntryFontFace(res.fonts.face),
		widget.ListOpts.SliderOpts(widget.SliderOpts.Images(res.images.sliderTrack, res.images.button)),
		widget.ListOpts.HideHorizontalSlider(),
		widget.ListOpts.HideVerticalSlider(),
		widget.ListOpts.ControlWidgetSpacing(2),

		widget.ListOpts.EntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			pageContainer.setPage(args.Entry.(*page))
		}))
	demoContainer.AddChild(pageList)

	demoContainer.AddChild(pageContainer.widget)

	pageList.SetSelectedEntry(pages[0])

	return demoContainer
}

func newPageContainer(res *resources) *pageContainer {
	c := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(15))),
	)

	titleText := widget.NewText(
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(&widget.RowLayoutData{
			Stretch: true,
		})),
		widget.TextOpts.Text("", res.fonts.titleFace, res.colors.textIdle))
	c.AddChild(titleText)

	flipBook := widget.NewFlipBook(
		widget.FlipBookOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(&widget.RowLayoutData{
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
		widget.LabeledCheckboxOpts.Spacing(6),
		widget.LabeledCheckboxOpts.CheckboxOpts(
			widget.CheckboxOpts.ButtonOpts(
				widget.ButtonOpts.Image(res.images.button),
				widget.ButtonOpts.GraphicPadding(widget.NewInsetsSimple(7)),
			),
			widget.CheckboxOpts.Image(res.images.checkbox),

			widget.CheckboxOpts.ChangedHandler(changedHandler)),
		widget.LabeledCheckboxOpts.LabelOpts(widget.LabelOpts.Text(label, res.fonts.face, res.colors.label)))
}

func newPageContentContainer() *widget.Container {
	return widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(10),
		)))
}

func newListComboButton(entries []interface{}, buttonLabel widget.SelectComboButtonEntryLabelFunc, entryLabel widget.ListEntryLabelFunc,
	entrySelectedHandler widget.ListComboButtonEntrySelectedHandlerFunc, res *resources) *widget.ListComboButton {

	return widget.NewListComboButton(
		widget.ListComboButtonOpts.SelectComboButtonOpts(widget.SelectComboButtonOpts.ComboButtonOpts(widget.ComboButtonOpts.ButtonOpts(widget.ButtonOpts.Image(res.images.button)))),
		widget.ListComboButtonOpts.Text(res.fonts.face, res.images.arrowDown, res.colors.buttonText),
		widget.ListComboButtonOpts.ListOpts(
			widget.ListOpts.Entries(entries),
			widget.ListOpts.ScrollContainerOpts(
				widget.ScrollContainerOpts.Image(res.images.scrollContainer),
				widget.ScrollContainerOpts.Padding(widget.NewInsetsSimple(2))),
			widget.ListOpts.SliderOpts(
				widget.SliderOpts.Images(res.images.sliderTrack, res.images.button),
				widget.SliderOpts.HandleSize(20),
				widget.SliderOpts.TrackPadding(2)),
			widget.ListOpts.ControlWidgetSpacing(0),
			widget.ListOpts.EntryFontFace(res.fonts.face),
			widget.ListOpts.EntryColor(res.colors.list)),
		widget.ListComboButtonOpts.EntryLabelFunc(buttonLabel, entryLabel),
		widget.ListComboButtonOpts.EntrySelectedHandler(entrySelectedHandler))
}

func newList(entries []interface{}, res *resources, widgetOpts ...widget.WidgetOpt) *widget.List {
	return widget.NewList(
		widget.ListOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(widgetOpts...)),
		widget.ListOpts.ScrollContainerOpts(
			widget.ScrollContainerOpts.Image(res.images.scrollContainer),
			widget.ScrollContainerOpts.Padding(widget.NewInsetsSimple(2))),
		widget.ListOpts.SliderOpts(
			widget.SliderOpts.Images(res.images.sliderTrack, res.images.button),
			widget.SliderOpts.HandleSize(20),
			widget.SliderOpts.TrackPadding(2)),
		widget.ListOpts.HideHorizontalSlider(),
		widget.ListOpts.Entries(entries),
		widget.ListOpts.EntryLabelFunc(func(e interface{}) string {
			return e.(string)
		}),
		widget.ListOpts.EntryFontFace(res.fonts.face),
		widget.ListOpts.EntryColor(res.colors.list),
	)
}

func newSeparator(res *resources, ld interface{}) widget.HasWidget {
	c := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.Insets{
				Top:    15,
				Bottom: 15,
			}))),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(ld)))

	c.AddChild(widget.NewGraphic(
		widget.GraphicOpts.WidgetOpts(widget.WidgetOpts.LayoutData(&widget.RowLayoutData{
			Stretch:   true,
			MaxHeight: 2,
		})),
		widget.GraphicOpts.ImageNineSlice(image.NewNineSliceColor(res.colors.separator)),
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
