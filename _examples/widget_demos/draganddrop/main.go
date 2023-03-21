package main

import (
	"fmt"
	img "image"
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
	dndObj         *widget.Container
	text           *widget.Text
	targetedWidget widget.HasWidget
}

// This method is used to create the drag element. It also allows you to provide abitrary drag data.
//
// Note that you can return the same container each time if you do not need to recreate it.
//
// Inputs:
//   - parent - The widget that triggered this Drag and drop event
func (dnd *dndWidget) Create(parent widget.HasWidget) (*widget.Container, interface{}) {
	// For this example we do not need to recreate the Dragged element each time. We can re-use it.
	if dnd.dndObj == nil {
		// load text font
		face, _ := loadFont(20)
		dnd.dndObj = widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
			widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0, 200, 100, 255})),
		)

		dnd.text = widget.NewText(widget.TextOpts.Text("Cannot Drop", face, color.Black), widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
		})))

		dnd.dndObj.AddChild(dnd.text)
	}
	// return the container to be dragged and any arbitrary data associated with this operation
	return dnd.dndObj, "Hello World"
}

// This method is optional for Drag and Drop
// It will be called every draw cycle that the Drag and Drop is active.
// Inputs:
//   - canDrop - if the cursor is over a widget that allows this object to be dropped
//   - targetWidget - The widget that will allow this object to be dropped.
//   - dragData - The drag data provided by the Create method above.
func (dnd *dndWidget) Update(canDrop bool, targetWidget widget.HasWidget, dragData interface{}) {
	if canDrop {
		dnd.text.Label = "* Can Drop *"
		if targetWidget != nil {
			targetWidget.(*widget.Container).BackgroundImage = image.NewNineSliceColor(color.NRGBA{100, 100, 255, 255})
			dnd.targetedWidget = targetWidget
		}
	} else {
		dnd.text.Label = "Cannot Drop"
		if dnd.targetedWidget != nil {
			dnd.targetedWidget.(*widget.Container).BackgroundImage = image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255})
			dnd.targetedWidget = nil
		}
	}
}

// This method is optional for Drag and Drop
// It will be called when the Drag and Drop is completed.
// Inputs:
//   - dropped - if drop was completed successfully
//   - targetWidget - The widget that will allow this object to be dropped.
//   - dragData - The drag data provided by the Create method above.
func (dnd *dndWidget) EndDrag(dropped bool, sourceWidget widget.HasWidget, dragData interface{}) {
	if dropped {
		fmt.Println("Dropped Successful")
	} else {
		fmt.Println("Drop Cancelled")
	}
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
			widget.WidgetOpts.EnableDragAndDrop(
				widget.NewDragAndDrop(
					//The object which will create/update the dragged element. Required.
					widget.DragAndDropOpts.ContentsCreater(&dndWidget{}),
					//How many pixels the user must drag their cursor before the drag begins.
					//This is an optional parameter that defaults to 15 pixels
					widget.DragAndDropOpts.MinDragStartDistance(15),
					//This sets where to anchor the widget to the cursor - vertical orientation
					widget.DragAndDropOpts.ContentsOriginVertical(widget.DND_ANCHOR_END),
					//This sets where to anchor the widget to the cursor - horizontal orientation
					widget.DragAndDropOpts.ContentsOriginHorizontal(widget.DND_ANCHOR_END),
					//This sets of far off the cursor to offset the dragged element
					widget.DragAndDropOpts.Offset(img.Point{-5, -5}),
					//This will turn of Drag to initiate drag and drop
					//Primary use case will be click to drag
					//widget.DragAndDropOpts.DisableDrag(),
				),
			),
			widget.WidgetOpts.MouseButtonReleasedHandler(func(args *widget.WidgetMouseButtonReleasedEventArgs) {
				if args.Inside && args.Button == ebiten.MouseButtonLeft && ebiten.IsKeyPressed(ebiten.KeyControl) {
					args.Widget.DragAndDrop.StartDrag()
				}
				if args.Button == ebiten.MouseButtonRight {
					args.Widget.DragAndDrop.StopDrag()
				}
			}),
		),
	)

	leftSide.AddChild(widget.NewText(widget.TextOpts.Text("Drag from Here\nOr Ctrl-Click", face, color.Black), widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
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
				rightTopText.Label = fmt.Sprintf("Drag to here\n(allowed)\n%d", count)
			}),
		),
	)
	rightTopText = widget.NewText(widget.TextOpts.Text(fmt.Sprintf("Drag to here\n(allowed)\n%d", count), face, color.Black), widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
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

	rightBottom.AddChild(widget.NewText(widget.TextOpts.Text("Drag to here\n(not allowed)", face, color.Black), widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
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
