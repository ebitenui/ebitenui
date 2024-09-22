package main

import (
	"image/color"
	"log"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Game object used by ebiten
type game struct {
	ui         *ebitenui.UI
	btn        *widget.Button
	lightTheme *widget.Theme
	darkTheme  *widget.Theme
}

func main() {
	// construct a new container that serves as the root of the UI hierarchy
	rootContainer := widget.NewContainer(
		// the container will use a plain color as its background
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),

		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	// construct a button
	button := widget.NewButton(
		// set general widget options
		widget.ButtonOpts.WidgetOpts(
			// instruct the container's anchor layout to center the button both horizontally and vertically
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
		// specify the button's text
		widget.ButtonOpts.TextLabel("Hello, World!"),

		widget.ButtonOpts.TextProcessBBCode(true),

		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			println("button clicked")
		}),
	)

	// add the button as a child of the container
	rootContainer.AddChild(button)

	lightTheme := GetLightTheme()
	darkTheme := GetDarkTheme()
	// construct the UI
	ui := ebitenui.UI{
		Container:    rootContainer,
		PrimaryTheme: darkTheme,
	}

	// Ebiten setup
	ebiten.SetWindowSize(400, 400)
	ebiten.SetWindowTitle("Ebiten UI - Buttons")

	game := game{
		ui:         &ui,
		btn:        button,
		lightTheme: lightTheme,
		darkTheme:  darkTheme,
	}

	// run Ebiten main loop
	err := ebiten.RunGame(&game)
	if err != nil {
		log.Println(err)
	}
}

// Layout implements Game.
func (g *game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

// Update implements Game.
func (g *game) Update() error {
	// update the UI
	g.ui.Update()
	if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		g.btn.Click()
	}

	//Test that you can call Click on the focused widget.
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		if btn, ok := g.ui.GetFocusedWidget().(*widget.Button); ok {
			btn.Click()
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if g.ui.PrimaryTheme == g.lightTheme {
			g.ui.PrimaryTheme = g.darkTheme
		} else {
			g.ui.PrimaryTheme = g.lightTheme
		}
	}

	return nil
}

// Draw implements Ebiten's Draw method.
func (g *game) Draw(screen *ebiten.Image) {
	// draw the UI onto the screen
	g.ui.Draw(screen)
}
