package main

import (
	"fmt"
	"log"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/themes"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

// A list of a hundred thousand entries: the VirtualList keeps widgets only
// for the rows in view, and binds them to the entries as it scrolls.
func main() {
	root := widget.NewContainer(widget.ContainerOpts.Layout(widget.NewAnchorLayout(
		widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(30)))))

	list := widget.NewVirtualList(
		widget.VirtualListOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{StretchHorizontal: true, StretchVertical: true}))),
		widget.VirtualListOpts.Count(100_000),
		widget.VirtualListOpts.RowHeight(24),
		// A row's widget, made only for as many rows as the view shows.
		widget.VirtualListOpts.RowFunc(func() widget.PreferredSizeLocateableWidget {
			return widget.NewText(widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionCenter))
		}),
		// Fills a row with the entry it now shows.
		widget.VirtualListOpts.BindFunc(func(row widget.PreferredSizeLocateableWidget, index int) {
			row.(*widget.Text).Label = fmt.Sprintf("Entry %d", index) //nolint:forcetypeassert // made by RowFunc
		}),
		// Dragging the rows, with a mouse or a finger, scrolls them too.
		widget.VirtualListOpts.DragToScroll(),
	)
	root.AddChild(list)

	ui := &ebitenui.UI{Container: root, PrimaryTheme: themes.GetBasicDarkTheme()}
	ebiten.SetWindowSize(400, 500)
	ebiten.SetWindowTitle("Ebiten UI - Virtual List")
	if err := ebiten.RunGame(&game{ui: ui}); err != nil {
		log.Println(err)
	}
}

type game struct{ ui *ebitenui.UI }

func (g *game) Update() error {
	g.ui.Update()
	return nil
}

func (g *game) Draw(screen *ebiten.Image) { g.ui.Draw(screen) }

func (g *game) Layout(w, h int) (int, int) { return w, h }
