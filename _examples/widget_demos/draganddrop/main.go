package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

// Game object used by ebiten
type game struct {
	ui *ebitenui.UI
}

// This object satisfies the interface DragContentsCreater
type dndWidget struct {
}

// This method is used to create the drag element. It also allows you to provide abitrary drag data.
//
// Note that you can return the same container each time if you do not need to recreate it.
//
// Inputs:
//   - parent - The widget that triggered this Drag and drop event
//   - cursorX - The X position of the cursor when the Drag and Drop event began
//   - cursorY - The Y position of the cursor when the Drag and Drop event began
func (dnd *dndWidget) Create(parent *widget.Widget, cursorX int, cursorY int) (*widget.Container, interface{}) {
	// load text font
	face, _ := loadFont(20)

	dndObj := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0, 200, 100, 255})),
	)

	dndObj.AddChild(widget.NewText(widget.TextOpts.Text("Drag FROM Here", face, color.Black), widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
		HorizontalPosition: widget.AnchorLayoutPositionCenter,
		VerticalPosition:   widget.AnchorLayoutPositionCenter,
	}))))

	// return the container to be dragged and any arbitrary data associated with this operation
	return dndObj, "Hello World"
}

func main() {
	// load text font
	face, _ := loadFont(20)
	count := 0
	// construct a new container that serves as the root of the UI hierarchy
	rootContainer := widget.NewContainer(
		// the container will use a plain color as its background
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),

		// the container will use a grid layout to layout to split the layout in half down the middle
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Padding(widget.NewInsetsSimple(10)),
			widget.GridLayoutOpts.Spacing(10, 10),
			widget.GridLayoutOpts.Stretch([]bool{true, true}, []bool{true}))),
	)

	leftSide := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255})),
		widget.ContainerOpts.WidgetOpts(
			//This command indicates this widget is the source of Drag and Drop.
			widget.WidgetOpts.EnableDragAndDrop(widget.NewDragAndDrop(widget.DragAndDropOpts.ContentsCreater(&dndWidget{}))),
		),
	)

	leftSide.AddChild(widget.NewText(widget.TextOpts.Text("Drag FROM Here", face, color.Black), widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
		HorizontalPosition: widget.AnchorLayoutPositionCenter,
		VerticalPosition:   widget.AnchorLayoutPositionCenter,
	}))))

	rootContainer.AddChild(leftSide)

	//This container splits the right side into two equal sections.
	rightSide := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(widget.GridLayoutOpts.Columns(1), widget.GridLayoutOpts.Stretch([]bool{true}, []bool{true, true}), widget.GridLayoutOpts.Spacing(10, 10))))

	var rightTop *widget.Container
	var rightTopText *widget.Text

	rightTop = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255})),
		widget.ContainerOpts.WidgetOpts(
			//This function is called to determine if the currently dragged element can be dropped on this widget
			widget.WidgetOpts.CanDrop(func(args *widget.DragAndDropDroppedEventArgs) bool {
				//This method is using the Data element (provided by the ContentsCreator above) to determine if the element can be dropped here.
				return args.Data.(string) == "Hello World"
			}),
			//This function is called if the client 'drops' the dragged element on this widget and CanDrop returns true
			widget.WidgetOpts.Dropped(func(args *widget.DragAndDropDroppedEventArgs) {
				rightTop.BackgroundImage = image.NewNineSliceColor(color.NRGBA{255, 100, 100, 255})
				count = count + 1
				rightTopText.Label = fmt.Sprintf("Drag TO Here\n(allowed)\n%d", count)
			}),
		),
	)
	rightTopText = widget.NewText(widget.TextOpts.Text(fmt.Sprintf("Drag TO Here\n(allowed)\n%d", count), face, color.Black), widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
		HorizontalPosition: widget.AnchorLayoutPositionCenter,
		VerticalPosition:   widget.AnchorLayoutPositionCenter,
	})))
	rightTop.AddChild(rightTopText)

	rightSide.AddChild(rightTop)

	var rightBottom *widget.Container
	rightBottom = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255})),
		widget.ContainerOpts.WidgetOpts(
			//This function is called to determine if the currently dragged element can be dropped on this widget
			widget.WidgetOpts.CanDrop(func(args *widget.DragAndDropDroppedEventArgs) bool {
				//This method is using the Data element (provided by the ContentsCreator above) to determine if the element can be dropped here.
				// In this example args.Data will always be "Hello World" so this will never match.
				return args.Data.(string) == "Does not match"
			}),
			//This function is called if the client 'drops' the dragged element on this widget and CanDrop returns true
			widget.WidgetOpts.Dropped(func(args *widget.DragAndDropDroppedEventArgs) {
				rightBottom.BackgroundImage = image.NewNineSliceColor(color.NRGBA{255, 100, 100, 255})
			}),
		),
	)

	rightBottom.AddChild(widget.NewText(widget.TextOpts.Text("Drag To Here\n(not allowed)", face, color.Black), widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
		HorizontalPosition: widget.AnchorLayoutPositionCenter,
		VerticalPosition:   widget.AnchorLayoutPositionCenter,
	}))))

	rightSide.AddChild(rightBottom)

	rootContainer.AddChild(rightSide)
	// construct the UI
	ui := ebitenui.UI{
		Container: rootContainer,
	}

	// Ebiten setup
	ebiten.SetWindowSize(400, 400)
	ebiten.SetWindowTitle("Ebiten UI - Drag And Drop")

	game := game{
		ui: &ui,
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
	return nil
}

// Draw implements Ebiten's Draw method.
func (g *game) Draw(screen *ebiten.Image) {
	// draw the UI onto the screen
	g.ui.Draw(screen)
}

func loadFont(size float64) (font.Face, error) {
	ttfFont, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(ttfFont, &truetype.Options{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	}), nil
}
