package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/themes"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

// Three hundred texts, drawn glyph by glyph or, with the checkbox ticked,
// each from an image of itself (widget.SetTextCache); the label shows how
// long the interface takes to draw.
func main() {
	root := widget.NewContainer(widget.ContainerOpts.Layout(widget.NewRowLayout(
		widget.RowLayoutOpts.Direction(widget.DirectionVertical),
		widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(16)),
		widget.RowLayoutOpts.Spacing(12))))

	g := &game{}
	g.timing = widget.NewText(widget.TextOpts.TextLabel("Drawing..."))
	root.AddChild(widget.NewCheckbox(
		widget.CheckboxOpts.TextLabel("Draw each text from an image of itself"),
		widget.CheckboxOpts.StateChangedHandler(func(args *widget.CheckboxChangedEventArgs) {
			widget.SetTextCache(args.State == widget.WidgetChecked)
		}),
	))
	root.AddChild(g.timing)

	grid := widget.NewContainer(widget.ContainerOpts.Layout(widget.NewGridLayout(
		widget.GridLayoutOpts.Columns(6),
		widget.GridLayoutOpts.Spacing(12, 4))))
	for i := range 300 {
		grid.AddChild(widget.NewText(widget.TextOpts.TextLabel(fmt.Sprintf("Text number %d", i))))
	}
	root.AddChild(grid)

	g.ui = &ebitenui.UI{Container: root, PrimaryTheme: themes.GetBasicDarkTheme()}
	ebiten.SetWindowSize(1000, 800)
	ebiten.SetWindowTitle("Ebiten UI - Text cache")
	if err := ebiten.RunGame(g); err != nil {
		log.Println(err)
	}
}

type game struct {
	ui     *ebitenui.UI
	timing *widget.Text
	spent  time.Duration
	frames int
}

func (g *game) Update() error {
	g.ui.Update()
	return nil
}

// Draw times the interface and shows the average every second.
func (g *game) Draw(screen *ebiten.Image) {
	start := time.Now()
	g.ui.Draw(screen)
	g.spent += time.Since(start)
	if g.frames++; g.frames == 60 {
		g.timing.Label = fmt.Sprintf("The interface takes %.0f µs to draw", float64(g.spent.Microseconds())/60)
		g.spent, g.frames = 0, 0
	}
}

func (g *game) Layout(w, h int) (int, int) { return w, h }
