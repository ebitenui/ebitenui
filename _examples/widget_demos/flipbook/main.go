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

type ListEntry struct {
	name   string
	widget *widget.Container
}

/*
Flipbook is an advanced component used to build other widgets such
as the tabbook. It is designed to only show a single widget/container at
a time set by the SetPage() function. The flipbook is an AnchorContainer
so anything that is Set should utilize the AnchorLayoutData to ensure it
is displayed properly as seen below.

This example shows Flipbook used with a combobox. Please review the
combobox widget demo prior to working with this one.
*/

func main() {
	// load images for button states: idle, hover, and pressed
	buttonImage, _ := loadButtonImage()

	// load text font
	face, _ := loadFont(16)

	// construct a new container that serves as the root of the UI hierarchy
	rootContainer := widget.NewContainer(
		// the container will use a plain color as its background
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),

		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			//A single column
			widget.GridLayoutOpts.Columns(1),
			//The objects are streched horizontally.
			//Only the second object is stretched vertically to create a header
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, true}),
		)),
	)

	// Create the first tab
	// A TabBookTab is a labelled container. The text here is what will show up in the tab button
	tabRed := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{255, 0, 0, 255})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		//Since this will be added to a FlipBook we need to tell it to fill the entire space
		// of the flipbook
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				StretchHorizontal:  true,
				StretchVertical:    true,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			})),
	)

	redBtn := widget.NewText(
		widget.TextOpts.Text("Red Tab Button", face, color.White),
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
		})),
	)
	tabRed.AddChild(redBtn)

	tabGreen := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0, 255, 0, 0xff})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		//Since this will be added to a FlipBook we need to tell it to fill the entire space
		// of the flipbook
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				StretchHorizontal:  true,
				StretchVertical:    true,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			})),
	)
	greenBtn := widget.NewText(
		widget.TextOpts.Text("Green Tab Button\nThis is configured as the initial tab.", face, color.Black),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
		})),
	)
	tabGreen.AddChild(greenBtn)

	tabBlue := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0, 0, 255, 0xff})),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(5),
		)),
		//Since this will be added to a FlipBook we need to tell it to fill the entire space
		// of the flipbook
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				StretchHorizontal:  true,
				StretchVertical:    true,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			})),
	)
	blueBtn1 := widget.NewText(
		widget.TextOpts.Text("Blue Tab Button 1", face, color.White),
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Position: widget.RowLayoutPositionCenter,
		})),
	)
	tabBlue.AddChild(blueBtn1)
	blueBtn2 := widget.NewText(
		widget.TextOpts.Text("Blue Tab Button 2", face, color.White),
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Position: widget.RowLayoutPositionCenter,
		})),
	)
	tabBlue.AddChild(blueBtn2)

	//Create the new flipbook its size is defined by its parent.
	//In this case it is defined by the GridLayout on the rootContainer
	flipBook := widget.NewFlipBook()

	//Set up the list of entries for the combobox mapping
	//the tab name to the tab container objects
	entries := []any{}
	entries = append(entries, ListEntry{"Red Tab", tabRed})
	entries = append(entries, ListEntry{"Green Tab", tabGreen})
	entries = append(entries, ListEntry{"Blue Tab", tabBlue})

	// construct a combobox
	comboBox := widget.NewListComboButton(
		widget.ListComboButtonOpts.SelectComboButtonOpts(
			widget.SelectComboButtonOpts.ComboButtonOpts(
				//Set the max height of the dropdown list
				widget.ComboButtonOpts.MaxContentHeight(150),
				//Set the parameters for the primary displayed button
				widget.ComboButtonOpts.ButtonOpts(
					widget.ButtonOpts.Image(buttonImage),
					widget.ButtonOpts.TextPadding(widget.NewInsetsSimple(5)),
					widget.ButtonOpts.Text("", face, &widget.ButtonTextColor{
						Idle:     color.White,
						Disabled: color.White,
					}),
					widget.ButtonOpts.WidgetOpts(
						//Set how wide the button should be
						widget.WidgetOpts.MinSize(150, 0),
						//Set the combobox's position
						widget.WidgetOpts.LayoutData(widget.GridLayoutData{
							HorizontalPosition: widget.GridLayoutPositionCenter,
							VerticalPosition:   widget.GridLayoutPositionCenter,
							MaxWidth:           150,
						}),
					),
				),
			),
		),
		widget.ListComboButtonOpts.ListOpts(
			//Set how wide the dropdown list should be
			widget.ListOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.MinSize(150, 0),
			)),
			//Set the entries in the list
			widget.ListOpts.Entries(entries),
			widget.ListOpts.ScrollContainerOpts(
				//Set the background images/color for the dropdown list
				widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
					Idle:     image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
					Disabled: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
					Mask:     image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				}),
			),
			widget.ListOpts.SliderOpts(
				//Set the background images/color for the background of the slider track
				widget.SliderOpts.Images(&widget.SliderTrackImage{
					Idle:  image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
					Hover: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				}, buttonImage),
				widget.SliderOpts.MinHandleSize(5),
				//Set how wide the track should be
				widget.SliderOpts.TrackPadding(widget.NewInsetsSimple(2))),
			//Set the font for the list options
			widget.ListOpts.EntryFontFace(face),
			//Set the colors for the list
			widget.ListOpts.EntryColor(&widget.ListEntryColor{
				Selected:                   color.NRGBA{254, 255, 255, 255},             //Foreground color for the unfocused selected entry
				Unselected:                 color.NRGBA{254, 255, 255, 255},             //Foreground color for the unfocused unselected entry
				SelectedBackground:         color.NRGBA{R: 130, G: 130, B: 200, A: 255}, //Background color for the unfocused selected entry
				SelectedFocusedBackground:  color.NRGBA{R: 130, G: 130, B: 170, A: 255}, //Background color for the focused selected entry
				FocusedBackground:          color.NRGBA{R: 170, G: 170, B: 180, A: 255}, //Background color for the focused unselected entry
				DisabledUnselected:         color.NRGBA{100, 100, 100, 255},             //Foreground color for the disabled unselected entry
				DisabledSelected:           color.NRGBA{100, 100, 100, 255},             //Foreground color for the disabled selected entry
				DisabledSelectedBackground: color.NRGBA{100, 100, 100, 255},             //Background color for the disabled selected entry
			}),
			//Padding for each entry
			widget.ListOpts.EntryTextPadding(widget.NewInsetsSimple(5)),
		),
		//Define how the entry is displayed
		widget.ListComboButtonOpts.EntryLabelFunc(
			func(e any) string {
				//Button Label function
				return "Button: " + e.(ListEntry).name
			},
			func(e any) string {
				//List Label function
				return "List: " + e.(ListEntry).name
			}),
		//Callback when a new entry is selected
		widget.ListComboButtonOpts.EntrySelectedHandler(func(args *widget.ListComboButtonEntrySelectedEventArgs) {
			fmt.Println("Selected Entry: ", args.Entry.(ListEntry).name)
			//On select, update the widget displayed by the flipbook
			flipBook.SetPage(args.Entry.(ListEntry).widget)
		}),
	)
	comboBox.SetSelectedEntry(entries[1])
	//The following line is needed if you dont set a selected entry in the combobox
	// since the callback for the combobox isn't called when selecting the first entry
	//flipBook.SetPage(entrys[0].(ListEntry).widget)
	rootContainer.AddChild(comboBox)
	rootContainer.AddChild(flipBook)

	// construct the UI
	ui := ebitenui.UI{
		Container: rootContainer,
	}

	// Ebiten setup
	ebiten.SetWindowSize(400, 400)
	ebiten.SetWindowTitle("Ebiten UI - Flipbook")

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

func loadButtonImage() (*widget.ButtonImage, error) {
	idle := image.NewNineSliceColor(color.NRGBA{R: 170, G: 170, B: 180, A: 255})

	hover := image.NewNineSliceColor(color.NRGBA{R: 130, G: 130, B: 150, A: 255})

	pressed := image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 120, A: 255})

	pressedHover := image.NewNineSliceColor(color.NRGBA{R: 110, G: 110, B: 110, A: 255})

	disabled := image.NewNineSliceColor(color.NRGBA{R: 80, G: 80, B: 140, A: 255})

	return &widget.ButtonImage{
		Idle:         idle,
		Hover:        hover,
		Pressed:      pressed,
		PressedHover: pressedHover,
		Disabled:     disabled,
	}, nil
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
