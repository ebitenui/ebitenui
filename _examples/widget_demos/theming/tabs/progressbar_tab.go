package tabs

import (
	"image/color"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
)

func NewProgressBarTab() *widget.TabBookTab {
	result := widget.NewTabBookTab("Progress Bar",
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	// Construct a container to hold the progress bars.
	progressBarsContainer := widget.NewContainer(
		// The container will use a vertical row layout to lay out the progress
		// bars in a vertical row.
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(20),
		)),
		// Set the required anchor layout data to determine where in the root
		// container to place the progress bars.
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
	)

	// Construct a horizontal progress bar.
	hProgressbar := widget.NewProgressBar(
		widget.ProgressBarOpts.WidgetOpts(
			// Set the minimum size for the progress bar.
			// This is necessary if you wish to have the progress bar be larger than
			// the provided track image. In this exampe since we are using NineSliceColor
			// which is 1px x 1px we must set a minimum size.
			widget.WidgetOpts.MinSize(200, 20),
		),
		widget.ProgressBarOpts.Images(
			// Set the track images (Idle, Disabled).
			nil,
			// Set the progress images (Idle, Disabled).
			&widget.ProgressBarImage{
				Idle: image.NewNineSliceColor(color.NRGBA{0, 0, 255, 255}),
			},
		),
		// Set the min, max, and current values.
		widget.ProgressBarOpts.Values(0, 10, 7),
	)
	// Construct a vertical inverted progress bar.
	vProgressbar := widget.NewProgressBar(
		// Set the direction of the progress bar to vertical.
		widget.ProgressBarOpts.Direction(widget.DirectionVertical),
		// Invert the progress bar, meaning here it will fill from the bottom to the top
		// since it’s vertical.
		widget.ProgressBarOpts.Inverted(true),
		widget.ProgressBarOpts.WidgetOpts(
			// Set the minimum size for the progress bar.
			// This is necessary if you wish to have the progress bar be larger than
			// the provided track image. In this example since we are using NineSliceColor
			// which is 1px x 1px we must set a minimum size.
			widget.WidgetOpts.MinSize(20, 200),

			// Set the progress bar in the middle of the row container cell
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),
		widget.ProgressBarOpts.Images(
			// Set the track images (Idle, Hover, Disabled).
			nil,
			// Set the progress images (Idle, Hover, Disabled).
			&widget.ProgressBarImage{
				Idle:  image.NewNineSliceColor(color.NRGBA{0, 255, 0, 255}),
				Hover: image.NewNineSliceColor(color.NRGBA{0, 255, 0, 255}),
			},
		),
		// Set the min, max, and current values.
		widget.ProgressBarOpts.Values(0, 10, 4),
	)

	progressBarsContainer.AddChild(hProgressbar)
	progressBarsContainer.AddChild(vProgressbar)
	result.AddChild(progressBarsContainer)

	return result
}
