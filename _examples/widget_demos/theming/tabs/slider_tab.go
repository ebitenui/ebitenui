package tabs

import (
	"fmt"

	"github.com/ebitenui/ebitenui/widget"
)

func NewSliderTab() *widget.TabBookTab {
	result := widget.NewTabBookTab("Slider",
		widget.TabBookTabOpts.ContainerOpts(
			widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		),
	)

	// construct a slider
	slider := widget.NewSlider(
		widget.SliderOpts.WidgetOpts(
			// Set the Widget to layout in the center on the screen
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
			// Set the widget's dimensions
			widget.WidgetOpts.MinSize(6, 200),
		),
		// Set the slider orientation - n/s vs e/w
		widget.SliderOpts.Orientation(widget.DirectionVertical),
		// Set the minimum and maximum value for the slider
		widget.SliderOpts.MinMax(0, 10),
		// Set the current value of the slider, without triggering a change event
		widget.SliderOpts.InitialCurrent(5),

		// Set the callback to call when the slider value is changed
		widget.SliderOpts.ChangedHandler(func(args *widget.SliderChangedEventArgs) {
			fmt.Println(args.Current, "dragging", args.Dragging)
		}),
	)
	// add the slider as a child of the container
	result.AddChild(slider)

	return result
}
