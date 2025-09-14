package tabs

import (
	"github.com/ebitenui/ebitenui/widget"
)

func NewLabelTab() *widget.TabBookTab {
	result := widget.NewTabBookTab(
		widget.TabBookTabOpts.Label("Label"),
		widget.TabBookTabOpts.ContainerOpts(
			widget.ContainerOpts.Layout(widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(10),
				widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(20)),
			)),
		),
	)

	/**
	There are two ways to create a label within ebitenui.
	1) widget.NewText - This is a simple method to create a text label.
	   You will want to use this one unless you have a need to provide a
		 separate Disabled color for the text label you're displaying.

	2) widget.NewLabel - This is a more complex method to create a text label that can be disabled.
		You will want to use this one when you want to be able to provide a separate
		Disabled color for the text label you're displaying

	*/

	label1 := widget.NewText(
		widget.TextOpts.TextLabel("Label 1 (NewText)"),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),
	)
	// Set the widget as Disabled. This does not affect NewText
	label1.GetWidget().Disabled = true
	// Add the first Text as a child of the container
	result.AddChild(label1)

	// Create a new label
	label2 := widget.NewLabel(
		widget.LabelOpts.LabelText("Label 2 (NewLabel - Enabled)"),
		widget.LabelOpts.TextOpts(
			widget.TextOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position: widget.RowLayoutPositionCenter,
				}),
			),
		),
	)
	// Add the first label as a child of the container
	result.AddChild(label2)

	// Create a new label
	label3 := widget.NewLabel(
		widget.LabelOpts.LabelText("Label 3 (NewLabel - Disabled)"),
		widget.LabelOpts.TextOpts(
			widget.TextOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position: widget.RowLayoutPositionCenter,
				}),
			),
		),
	)
	// Add the second label as a child of the container
	result.AddChild(label3)
	// Set this label as disabled and tells it to use the
	// Disabled color.
	label3.GetWidget().Disabled = true
	return result
}
