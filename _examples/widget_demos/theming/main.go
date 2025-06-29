package main

import (
	"log"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/_examples/widget_demos/theming/tabs"
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

	tabList := []*widget.TabBookTab{
		tabs.NewButtonTab(),
		tabs.NewLabelTab(),
		tabs.NewTextInputTab(),
		tabs.NewTextAreaTab(),
		tabs.NewProgressBarTab(),
		tabs.NewSliderTab(),
		tabs.NewListTab(),
		tabs.NewSelectTab(),
	}

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

		// Set the current Tabs.
		widget.TabBookOpts.Tabs(tabList...),
	)

	rootContainer.AddChild(tabBook)
	lightTheme := themes.GetBasicLightTheme()
	darkTheme := themes.GetBasicDarkTheme()
	// construct the UI
	ui := ebitenui.UI{
		Container:    rootContainer,
		PrimaryTheme: darkTheme,
	}

	// Ebiten setup
	ebiten.SetWindowSize(1200, 400)
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
