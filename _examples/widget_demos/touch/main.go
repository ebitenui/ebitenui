package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/input"
	"github.com/ebitenui/ebitenui/themes"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

// Touch on a phone, served from a computer on the same network:
//
//	go run github.com/hajimehoshi/wasmserve@latest -http=:6262 ./_examples/widget_demos/touch
//
// then open http://<the computer's address>:6262 on the phone. A tap
// clicks and leaves nothing hovered; a long press on a button opens its
// menu, or shows its tooltip after its delay, and clicks nothing; the list scrolls by
// dragging; the input brings up the phone's keyboard; two fingers pinch
// and pan the numbers at the bottom, and click nothing.
func main() {
	face := loadFace(18)
	stretch := widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true})
	// A panel: the theme's background.
	root := widget.NewPanel(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(16)),
			widget.RowLayoutOpts.Spacing(12))))
	status := widget.NewText(widget.TextOpts.WidgetOpts(stretch))
	status.Label = "Tap, long-press, drag, pinch."
	gesture := widget.NewText(widget.TextOpts.WidgetOpts(stretch))

	// A long press opens the menu, as a right click does.
	menu := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x30, 0x3a, 0x46, 0xff})),
		widget.ContainerOpts.Layout(widget.NewRowLayout(widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(6)), widget.RowLayoutOpts.Spacing(4))))
	for _, name := range []string{"Copy", "Delete"} {
		menu.AddChild(widget.NewButton(widget.ButtonOpts.TextLabel(name),
			widget.ButtonOpts.ClickedHandler(func(*widget.ButtonClickedEventArgs) { status.Label = name + " from the menu" })))
	}
	withMenu := widget.NewButton(widget.ButtonOpts.WidgetOpts(stretch, widget.WidgetOpts.ContextMenu(menu)),
		widget.ButtonOpts.TextLabel("Long-press me for a menu"),
		widget.ButtonOpts.ClickedHandler(func(*widget.ButtonClickedEventArgs) { status.Label = "Tapped the menu button" }))

	// A long press shows the tooltip, after its delay, until the next touch.
	tip := widget.NewTextToolTip("From a long press", &face, color.White,
		image.NewNineSliceColor(color.NRGBA{0x30, 0x3a, 0x46, 0xff}))
	withTip := widget.NewButton(widget.ButtonOpts.TextLabel("Long-press me for a tooltip"),
		widget.ButtonOpts.WidgetOpts(stretch, widget.WidgetOpts.ToolTip(tip)),
		widget.ButtonOpts.ClickedHandler(func(*widget.ButtonClickedEventArgs) { status.Label = "Tapped the tooltip button" }))

	entries := make([]any, 50)
	for i := range entries {
		entries[i] = fmt.Sprintf("Row %d", i+1)
	}
	list := widget.NewList(
		widget.ListOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true, MaxHeight: 220}))),
		widget.ListOpts.Entries(entries),
		widget.ListOpts.HideHorizontalSlider(),
		widget.ListOpts.EntryLabelFunc(func(e any) string { return e.(string) }), //nolint:forcetypeassert // strings
		widget.ListOpts.DragToScroll(),
		widget.ListOpts.EntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			status.Label = fmt.Sprint("Chose ", args.Entry)
		}),
	)
	name := widget.NewTextInput(widget.TextInputOpts.WidgetOpts(stretch), widget.TextInputOpts.IME(),
		widget.TextInputOpts.Placeholder("Type with the phone's keyboard"))
	root.AddChild(status, withMenu, withTip, list, name, gesture)

	ui := &ebitenui.UI{Container: root, PrimaryTheme: themes.GetBasicDarkTheme()}
	ebiten.SetWindowSize(420, 720)
	ebiten.SetWindowTitle("Ebiten UI - Touch")
	if err := ebiten.RunGame(&game{ui: ui, gesture: gesture, zoom: 1}); err != nil {
		log.Println(err)
	}
}

type game struct {
	ui         *ebitenui.UI
	gesture    *widget.Text
	zoom       float64
	panX, panY float64
}

// Update follows two fingers as a game's camera would.
func (g *game) Update() error {
	g.ui.Update()
	g.zoom *= input.TouchPinch()
	x, y := input.TouchPan()
	g.panX, g.panY = g.panX+x, g.panY+y
	g.gesture.Label = fmt.Sprintf("Fingers %d. Zoom ×%.2f, pan %.0f, %.0f", input.TouchCount(), g.zoom, g.panX, g.panY)
	return nil
}

func (g *game) Draw(screen *ebiten.Image) { g.ui.Draw(screen) }

func (g *game) Layout(w, h int) (int, int) { return w, h }

func loadFace(size float64) text.Face {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatal(err)
	}
	return &text.GoTextFace{Source: s, Size: size}
}
