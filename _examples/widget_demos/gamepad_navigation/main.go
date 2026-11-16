package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/themes"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

// The arrow keys and a gamepad navigate the same way: the D-pad or the left
// stick moves through the list and the slider while they have the focus,
// and a gamepad moves the focus on past them; A activates, B or Escape goes
// back. A direction held repeats.
func main() {
	stretch := widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true})
	root := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(20)),
			widget.RowLayoutOpts.Spacing(12))))

	help := widget.NewText(widget.TextOpts.WidgetOpts(stretch), widget.TextOpts.MaxWidth(460))
	help.Label = "Arrow keys, or a gamepad's D-pad or left stick. A activates, B or Escape goes back."
	status := widget.NewText(widget.TextOpts.WidgetOpts(stretch))
	label := func(e any) string { return e.(string) } //nolint:forcetypeassert // the entries are strings

	entries := make([]any, 30)
	for i := range entries {
		entries[i] = fmt.Sprintf("Save %d", i+1)
	}
	list := widget.NewList(
		widget.ListOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true, MaxHeight: 160}))),
		widget.ListOpts.Entries(entries),
		widget.ListOpts.HideHorizontalSlider(),
		widget.ListOpts.EntryLabelFunc(label),
		widget.ListOpts.EntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			status.Label = "Loaded " + label(args.Entry)
		}),
	)
	slider := widget.NewSlider(
		widget.SliderOpts.WidgetOpts(stretch, widget.WidgetOpts.MinSize(0, 24)),
		widget.SliderOpts.MinMax(0, 100),
		widget.SliderOpts.InitialCurrent(50),
		widget.SliderOpts.PageSizeFunc(func() int { return 5 }),
		widget.SliderOpts.ChangedHandler(func(args *widget.SliderChangedEventArgs) {
			status.Label = fmt.Sprintf("Volume %d", args.Current)
		}),
	)
	combo := widget.NewListComboButton(
		widget.ListComboButtonOpts.WidgetOpts(stretch),
		widget.ListComboButtonOpts.Entries([]any{"Easy", "Normal", "Hard"}),
		widget.ListComboButtonOpts.EntryLabelFunc(label, label),
		widget.ListComboButtonOpts.EntrySelectedHandler(func(args *widget.ListComboButtonEntrySelectedEventArgs) {
			status.Label = "Difficulty " + label(args.Entry)
		}),
	)
	subtitles := widget.NewCheckbox(
		widget.CheckboxOpts.TextLabel("Subtitles"),
		widget.CheckboxOpts.StateChangedHandler(func(args *widget.CheckboxChangedEventArgs) {
			status.Label = fmt.Sprintf("Subtitles %v", args.State == widget.WidgetChecked)
		}),
	)
	buttons := widget.NewContainer(widget.ContainerOpts.Layout(widget.NewRowLayout(widget.RowLayoutOpts.Spacing(12))))
	for _, name := range []string{"Play", "Quit"} {
		buttons.AddChild(widget.NewButton(
			widget.ButtonOpts.TextLabel(name),
			widget.ButtonOpts.ClickedHandler(func(*widget.ButtonClickedEventArgs) { status.Label = name }),
		))
	}
	root.AddChild(help, list, slider, combo, subtitles, buttons, status)

	backs := 0
	ui := &ebitenui.UI{
		Container:    root,
		PrimaryTheme: themes.GetBasicDarkTheme(),
		GamepadFocus: true,
		// Back no widget takes: a game closes its menu here.
		BackHandler: func() {
			backs++
			status.Label = fmt.Sprintf("Back pressed: %d", backs)
		},
	}
	ui.SetFocusedWidget(list)

	ebiten.SetWindowSize(500, 560)
	ebiten.SetWindowTitle("Ebiten UI - Gamepad Navigation")
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
