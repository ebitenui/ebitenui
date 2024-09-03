package main

import (
	"bytes"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/gofont/goregular"
	"log"
)

const (
	Width  = 800
	Height = 600
)

func main() {
	// Set up Ebiten
	ebiten.SetWindowSize(Width, Height)
	ebiten.SetWindowTitle("Ebitenui Toolbar Demo")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	res, err := loadResources()
	if err != nil {
		log.Fatal(err)
	}

	// Construct a new container that serves as the root of the UI hierarchy.
	root := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	ui := ebitenui.UI{
		Container: root,
	}

	// Create a toolbar and add it to the UI.
	toolbar := newToolbar(&ui, res)
	root.AddChild(toolbar.container)

	// Set up the ebiten game struct.
	game := game{
		ui: &ui,
	}

	// Event handling
	//
	// Example 1: Configure the "Help" button to display a message in console when it's pressed.
	//
	toolbar.helpButton.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			println("The help button was pressed!")
		}),
	)

	// Example 2: Configure the "Quit" menu entry to end the program when it's pressed.
	toolbar.quitButton.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			game.exit = true
		}),
	)

	// Run the game.
	if err := ebiten.RunGame(&game); err != nil {
		log.Println(err)
	}
}

type game struct {
	ui   *ebitenui.UI
	exit bool
}

// Update implements ebiten.Game.
func (g *game) Update() error {
	// Exit the game if the exit flag is set.
	if g.exit {
		return ebiten.Termination
	}

	// Update the UI
	g.ui.Update()

	return nil
}

// Draw implements ebiten.Game.
func (g *game) Draw(screen *ebiten.Image) {
	// Clear the screen with the color teal
	screen.Fill(colornames.Teal)

	// Draw the UI onto the screen
	g.ui.Draw(screen)
}

// Layout implements ebiten.Game.
func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

type resources struct {
	font text.Face
}

func loadResources() (*resources, error) {
	fnt, err := loadFont(16)
	if err != nil {
		return nil, err
	}
	return &resources{
		font: fnt,
	}, nil
}

func loadFont(size float64) (text.Face, error) {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &text.GoTextFace{
		Source: s,
		Size:   size,
	}, nil
}
