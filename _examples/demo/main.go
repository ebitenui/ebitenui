package main

import (
	"fmt"
	"log"
	"sort"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"image/color"
	_ "image/png"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
)

type game struct {
	ui *ebitenui.UI
}

type pageContainer struct {
	widget    widget.PreferredSizeLocateableWidget
	titleText *widget.Text
	flipBook  *widget.FlipBook
}

func main() {
	ebiten.SetWindowSize(900, 800)
	ebiten.SetWindowTitle("Ebiten UI Demo")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetVsyncEnabled(false)

	ui, closeUI, err := createUI()
	if err != nil {
		log.Fatal(err)
	}

	defer closeUI()

	game := game{
		ui: ui,
	}

	err = ebiten.RunGame(&game)
	if err != nil {
		log.Print(err)
	}
}

func createUI() (*ebitenui.UI, func(), error) {
	res, err := newUIResources()
	if err != nil {
		return nil, nil, err
	}

	//This creates the root container for this UI.
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			// It is using a GridLayout with a single column
			widget.GridLayoutOpts.Columns(1),
			// It uses the Stretch parameter to define how the rows will be layed out.
			// - a fixed sized header
			// - a content row that stretches to fill all remaining space
			// - a fixed sized footer
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, true, false}),
			// Padding defines how much space to put around the outside of the grid.
			widget.GridLayoutOpts.Padding(widget.Insets{
				Top:    20,
				Bottom: 20,
			}),
			// Spacing defines how much space to put between each column and row
			widget.GridLayoutOpts.Spacing(0, 20))),
		widget.ContainerOpts.BackgroundImage(res.background))

	rootContainer.AddChild(headerContainer(res))

	var ui *ebitenui.UI
	rootContainer.AddChild(demoContainer(res, func() *ebitenui.UI {
		return ui
	}))

	footerContainer := widget.NewContainer(widget.ContainerOpts.Layout(widget.NewRowLayout(
		widget.RowLayoutOpts.Padding(widget.Insets{
			Left:  25,
			Right: 25,
		}),
	)))
	rootContainer.AddChild(footerContainer)

	footerContainer.AddChild(widget.NewText(
		widget.TextOpts.Text("github.com/ebitenui/ebitenui", res.text.smallFace, res.text.disabledColor)))

	ui = &ebitenui.UI{
		Container: rootContainer,
	}

	return ui, func() {
		res.close()
	}, nil
}

func headerContainer(res *uiResources) widget.PreferredSizeLocateableWidget {
	c := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(15))),
	)

	c.AddChild(header("Ebiten UI Demo", res,
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
	))

	c2 := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.Insets{
				Left:  25,
				Right: 25,
			}),
		)),
	)
	c.AddChild(c2)

	c2.AddChild(widget.NewText(
		widget.TextOpts.Text("This program is a showcase of Ebiten UI widgets and layouts.", res.text.face, res.text.idleColor)))

	return c
}

func header(label string, res *uiResources, opts ...widget.ContainerOpt) widget.PreferredSizeLocateableWidget {
	c := widget.NewContainer(append(opts, []widget.ContainerOpt{
		widget.ContainerOpts.BackgroundImage(res.header.background),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(widget.AnchorLayoutOpts.Padding(res.header.padding))),
	}...)...)

	c.AddChild(widget.NewText(
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionStart,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
		})),
		widget.TextOpts.Text(label, res.header.face, res.header.color),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionCenter),
	))

	return c
}

func demoContainer(res *uiResources, ui func() *ebitenui.UI) widget.PreferredSizeLocateableWidget {

	demoContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Padding(widget.Insets{
				Left:  25,
				Right: 25,
			}),
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
		toolTipPage(res),
		dragAndDropPage(res),
		textInputPage(res),
		radioGroupPage(res),
		windowPage(res, ui),
		anchorLayoutPage(res),
		textAreaPage(res),
		progressBarPage(res),
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
		widget.ListOpts.ScrollContainerOpts(widget.ScrollContainerOpts.Image(res.list.image)),
		widget.ListOpts.SliderOpts(
			widget.SliderOpts.Images(res.list.track, res.list.handle),
			widget.SliderOpts.MinHandleSize(res.list.handleSize),
			widget.SliderOpts.TrackPadding(res.list.trackPadding),
		),
		widget.ListOpts.EntryColor(res.list.entry),
		widget.ListOpts.EntryFontFace(res.list.face),
		widget.ListOpts.EntryTextPadding(res.list.entryPadding),
		widget.ListOpts.HideHorizontalSlider(),

		widget.ListOpts.EntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			pageContainer.setPage(args.Entry.(*page))
		}))
	demoContainer.AddChild(pageList)

	demoContainer.AddChild(pageContainer.widget)

	pageList.SetSelectedEntry(pages[0])

	return demoContainer
}

func newPageContainer(res *uiResources) *pageContainer {
	c := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(res.panel.image),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(res.panel.padding),
			widget.RowLayoutOpts.Spacing(15))),
	)

	titleText := widget.NewText(
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.TextOpts.Text("", res.text.titleFace, res.text.idleColor))
	c.AddChild(titleText)

	flipBook := widget.NewFlipBook(
		widget.FlipBookOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		}))),
	)
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

func newCheckbox(label string, changedHandler widget.CheckboxChangedHandlerFunc, res *uiResources) *widget.LabeledCheckbox {
	return widget.NewLabeledCheckbox(
		widget.LabeledCheckboxOpts.Spacing(res.checkbox.spacing),
		widget.LabeledCheckboxOpts.CheckboxOpts(
			widget.CheckboxOpts.ButtonOpts(widget.ButtonOpts.Image(res.checkbox.image)),
			widget.CheckboxOpts.Image(res.checkbox.graphic),
			widget.CheckboxOpts.StateChangedHandler(func(args *widget.CheckboxChangedEventArgs) {
				if changedHandler != nil {
					changedHandler(args)
				}
			})),
		widget.LabeledCheckboxOpts.LabelOpts(widget.LabelOpts.Text(label, res.label.face, res.label.text)))
}

func newPageContentContainer() *widget.Container {
	return widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
		})),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(10),
		)))
}

func newListComboButton(entries []interface{}, buttonLabel widget.SelectComboButtonEntryLabelFunc, entryLabel widget.ListEntryLabelFunc,
	entrySelectedHandler widget.ListComboButtonEntrySelectedHandlerFunc, res *uiResources) *widget.ListComboButton {

	return widget.NewListComboButton(
		widget.ListComboButtonOpts.SelectComboButtonOpts(
			widget.SelectComboButtonOpts.ComboButtonOpts(
				widget.ComboButtonOpts.ButtonOpts(
					widget.ButtonOpts.Image(res.comboButton.image),
					widget.ButtonOpts.TextPadding(res.comboButton.padding),
				),
			),
		),
		widget.ListComboButtonOpts.Text(res.comboButton.face, res.comboButton.graphic, res.comboButton.text),
		widget.ListComboButtonOpts.ListOpts(
			widget.ListOpts.Entries(entries),
			widget.ListOpts.ScrollContainerOpts(
				widget.ScrollContainerOpts.Image(res.list.image),
			),
			widget.ListOpts.SliderOpts(
				widget.SliderOpts.Images(res.list.track, res.list.handle),
				widget.SliderOpts.MinHandleSize(res.list.handleSize),
				widget.SliderOpts.TrackPadding(res.list.trackPadding)),
			widget.ListOpts.EntryFontFace(res.list.face),
			widget.ListOpts.EntryColor(res.list.entry),
			widget.ListOpts.EntryTextPadding(res.list.entryPadding),
		),
		widget.ListComboButtonOpts.EntryLabelFunc(buttonLabel, entryLabel),
		widget.ListComboButtonOpts.EntrySelectedHandler(entrySelectedHandler))
}

func newList(entries []interface{}, res *uiResources, widgetOpts ...widget.WidgetOpt) *widget.List {
	return widget.NewList(
		widget.ListOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(widgetOpts...)),
		widget.ListOpts.ScrollContainerOpts(widget.ScrollContainerOpts.Image(res.list.image)),
		widget.ListOpts.SliderOpts(
			widget.SliderOpts.Images(res.list.track, res.list.handle),
			widget.SliderOpts.MinHandleSize(res.list.handleSize),
			widget.SliderOpts.TrackPadding(res.list.trackPadding),
		),
		widget.ListOpts.HideHorizontalSlider(),
		widget.ListOpts.Entries(entries),
		widget.ListOpts.EntryLabelFunc(func(e interface{}) string {
			return e.(string)
		}),
		widget.ListOpts.EntryFontFace(res.list.face),
		widget.ListOpts.EntryColor(res.list.entry),
		widget.ListOpts.EntryTextPadding(res.list.entryPadding),
	)
}
func newTextArea(text string, res *uiResources, widgetOpts ...widget.WidgetOpt) *widget.TextArea {
	return widget.NewTextArea(
		widget.TextAreaOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(widgetOpts...)),
		widget.TextAreaOpts.ScrollContainerOpts(widget.ScrollContainerOpts.Image(res.list.image)),
		widget.TextAreaOpts.SliderOpts(
			widget.SliderOpts.Images(res.list.track, res.list.handle),
			widget.SliderOpts.MinHandleSize(res.list.handleSize),
			widget.SliderOpts.TrackPadding(res.list.trackPadding),
		),
		widget.TextAreaOpts.ShowVerticalScrollbar(),
		widget.TextAreaOpts.VerticalScrollMode(widget.PositionAtEnd),
		widget.TextAreaOpts.ProcessBBCode(true),
		widget.TextAreaOpts.FontFace(res.textArea.face),
		widget.TextAreaOpts.FontColor(color.NRGBA{R: 200, G: 100, B: 0, A: 255}),
		widget.TextAreaOpts.TextPadding(res.textArea.entryPadding),
		widget.TextAreaOpts.Text(text),
	)
}

func newSeparator(res *uiResources, ld interface{}) widget.PreferredSizeLocateableWidget {
	c := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.Insets{
				Top:    20,
				Bottom: 20,
			}))),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(ld)))

	c.AddChild(widget.NewGraphic(
		widget.GraphicOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch:   true,
			MaxHeight: 2,
		})),
		widget.GraphicOpts.ImageNineSlice(image.NewNineSliceColor(res.separatorColor)),
	))

	return c
}

func (g *game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (g *game) Update() error {
	g.ui.Update()
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	g.ui.Draw(screen)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %f", ebiten.ActualFPS()))
}
