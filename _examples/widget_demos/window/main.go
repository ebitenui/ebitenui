package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/ebitenui/ebitenui"
	e_image "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

// Game object used by ebiten.
type game struct {
	ui *ebitenui.UI
}

func main() {
	// construct the UI
	ui := ebitenui.UI{}

	// construct a new container that serves as the root of the UI hierarchy
	rootContainer := widget.NewContainer(
		// the container will use a plain color as its background
		widget.ContainerOpts.BackgroundImage(e_image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),

		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	// construct a new row layout container that is centered on the page
	centerContainer := widget.NewContainer(
		// Configure the container to be centered in it's parent
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
		// Configure the container to be a Vertical Row Layout.
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(10),
			)),
	)

	window1 := createWindow(&ui, "Window 1")
	centerContainer.AddChild(createButton(&ui, window1, "Window 1", image.Pt(100, 50)))

	window2 := createWindow(&ui, "Window 2")
	centerContainer.AddChild(createButton(&ui, window2, "Window 2", image.Pt(100, 260)))

	rootContainer.AddChild(centerContainer)

	// Set Root Container
	ui.Container = rootContainer

	// Ebiten setup
	ebiten.SetWindowSize(400, 400)
	ebiten.SetWindowTitle("Ebiten UI - Window")

	game := game{
		ui: &ui,
	}

	// run Ebiten main loop
	err := ebiten.RunGame(&game)
	if err != nil {
		log.Println(err)
	}
}

func createButton(ui *ebitenui.UI, win *widget.Window, label string, winPos image.Point) *widget.Button {
	// load images for button states: idle, hover, and pressed
	buttonImage, _ := loadButtonImage() // load button text font
	face, _ := loadFont(20)
	var btn *widget.Button

	btn = widget.NewButton(
		// set general widget options
		widget.ButtonOpts.WidgetOpts(
			// instruct the container's anchor layout to center the button both horizontally and vertically
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),

		// specify the images to use
		widget.ButtonOpts.Image(buttonImage),

		// specify the button's text, the font face, and the color
		widget.ButtonOpts.Text("Open "+label, face, &widget.ButtonTextColor{
			Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
		}),

		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    5,
			Bottom: 5,
		}),

		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			if !ui.IsWindowOpen(win) {
				btn.Text().Label = "Close " + label
				// Get the preferred size of the content
				x, y := win.Contents.PreferredSize()
				// Create a rect with the preferred size of the content
				r := image.Rect(0, 0, x, y)
				// Use the Add method to move the window to the specified point
				r = r.Add(winPos)
				// Set the windows location to the rect.
				win.SetLocation(r)
				// Add the window to the UI.
				// Note: If the window is already added, this will just move the window and not add a duplicate.
				ui.AddWindow(win)
			} else {
				win.Close()
				btn.Text().Label = "Open " + label
			}
		}),
	)
	return btn
}
func createWindow(ui *ebitenui.UI, label string) *widget.Window {
	// load the font for the window title
	titleFace, _ := loadFont(12)
	face, _ := loadFont(20)

	// Create the contents of the window
	windowContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(e_image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)
	windowContainer.AddChild(widget.NewText(
		widget.TextOpts.Text("Hello from "+label, face, color.NRGBA{254, 255, 255, 255}),
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
		})),
	))

	// Create the titlebar for the window
	titleContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(e_image.NewNineSliceColor(color.NRGBA{150, 150, 150, 255})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)
	titleContainer.AddChild(widget.NewText(
		widget.TextOpts.Text(label+" Title", titleFace, color.NRGBA{254, 255, 255, 255}),
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
		})),
	))

	// Create the new window object. The window object is not tied to a container. Its location and
	// size are set manually using the SetLocation method on the window and added to the UI with ui.AddWindow()
	// Set the Button callback below to see how the window is added to the UI.
	return widget.NewWindow(
		// Set the main contents of the window
		widget.WindowOpts.Contents(windowContainer),
		// Set the titlebar for the window (Optional)
		widget.WindowOpts.TitleBar(titleContainer, 25),
		// Set the window above everything else and block input elsewhere
		// widget.WindowOpts.Modal(),
		// Set how to close the window. CLICK_OUT will close the window when clicking anywhere
		// that is not a part of the window object
		// widget.WindowOpts.CloseMode(widget.CLICK_OUT),
		// Indicates that the window is draggable. It must have a TitleBar for this to work
		widget.WindowOpts.Draggable(),
		// Set the window resizeable
		widget.WindowOpts.Resizeable(),
		// Set the minimum size the window can be
		widget.WindowOpts.MinSize(200, 100),
		// Set the maximum size a window can be
		widget.WindowOpts.MaxSize(300, 300),
		// Set the callback that triggers when a move is complete
		widget.WindowOpts.MoveHandler(func(args *widget.WindowChangedEventArgs) {
			fmt.Println("Window Moved")
		}),
		// Set the callback that triggers when a resize is complete
		widget.WindowOpts.ResizeHandler(func(args *widget.WindowChangedEventArgs) {
			fmt.Println("Window Resized")
		}),
	)
}

// Layout implements Game.
func (g *game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

// Update implements Game.
func (g *game) Update() error {
	// update the UI
	g.ui.Update()
	return nil
}

// Draw implements Ebiten's Draw method.
func (g *game) Draw(screen *ebiten.Image) {
	// draw the UI onto the screen
	g.ui.Draw(screen)
}

func loadButtonImage() (*widget.ButtonImage, error) {
	idle := e_image.NewNineSliceColor(color.NRGBA{R: 170, G: 170, B: 180, A: 255})

	hover := e_image.NewNineSliceColor(color.NRGBA{R: 130, G: 130, B: 150, A: 255})

	pressed := e_image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 120, A: 255})

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
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
