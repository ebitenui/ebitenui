package main

import (
	"image/color"
	"log"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/themes"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Game object used by ebiten.
// Game object used by ebiten.
type game struct {
	ui         *ebitenui.UI
	btn        *widget.Button
	lightTheme *widget.Theme
	darkTheme  *widget.Theme
}

func main() {
	// construct a new container that serves as the root of the UI hierarchy
	rootContainer := widget.NewPanel(
		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	buttonTab := widget.NewTabBookTab("Button",
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{255, 255, 255, 20})),
		widget.ContainerOpts.Layout(widget.NewRowLayout(widget.RowLayoutOpts.Direction(widget.DirectionVertical))),
	)
	labelTab := widget.NewTabBookTab("Label",
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{255, 255, 255, 20})),
		widget.ContainerOpts.Layout(widget.NewRowLayout(widget.RowLayoutOpts.Direction(widget.DirectionVertical))),
	)

	tabBook := widget.NewTabBook(
		widget.TabBookOpts.ContainerOpts(
			widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				StretchHorizontal:  true,
				StretchVertical:    true,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
			),
		),

		//Set the current Tabs
		widget.TabBookOpts.Tabs(buttonTab, labelTab),
		// Set the Initial Tab
		widget.TabBookOpts.InitialTab(buttonTab),
	)
	/*

		// construct a button
		button := widget.NewButton(
			// specify the button's text
			widget.ButtonOpts.TextLabel("Hello, World!"),

			// add a handler that reacts to clicking the button.
			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				println("button clicked")
			}),
		)
		// add the button as a child of the container
		rootContainer.AddChild(button)

		rootContainer.AddChild(widget.NewLabel(
			widget.LabelOpts.LabelText("Label"),
		))

		rootContainer.AddChild(widget.NewText(
			widget.TextOpts.TextLabel("Text"),
		))
	*/

	rootContainer.AddChild(tabBook)
	lightTheme := themes.GetBasicLightTheme()
	darkTheme := themes.GetBasicDarkTheme()
	// construct the UI
	ui := ebitenui.UI{
		Container:    rootContainer,
		PrimaryTheme: darkTheme,
	}

	// Ebiten setup
	ebiten.SetWindowSize(400, 400)
	ebiten.SetWindowTitle("Ebiten UI - Theming")

	game := game{
		ui:         &ui,
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
