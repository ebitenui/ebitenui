package tabs

import (
	"fmt"

	"github.com/ebitenui/ebitenui/widget"
)

func NewTextAreaTab() *widget.TabBookTab {
	result := widget.NewTabBookTab("Text Area",
		widget.TabBookTabOpts.ContainerOpts(
			widget.ContainerOpts.Layout(widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(10),
				widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(20)),
			)),
		),
	)
	// construct a textarea
	textarea := widget.NewTextArea(
		widget.TextAreaOpts.ContainerOpts(
			widget.ContainerOpts.WidgetOpts(
				// Set the layout data for the textarea
				// including a max height to ensure the scroll bar is visible
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position:  widget.RowLayoutPositionCenter,
					MaxWidth:  400,
					MaxHeight: 200,
					Stretch:   true,
				}),
				// Set the minimum size for the widget
				widget.WidgetOpts.MinSize(400, 200),
			),
		),
		widget.TextAreaOpts.ProcessBBCode(true),
		// Set the initial text for the textarea
		// It will automatically line wrap and process newlines characters
		// If ProcessBBCode is true it will parse out bbcode
		widget.TextAreaOpts.Text("[link=a]Hello[/link] [color=#FFF000]World[/color] Blue bottle praxis raclette, beard try-hard paleo roof party small batch. Dreamcatcher ascot next level lomo trust fund copper mug franzen farm-to-table hashtag. Four dollar toast activated charcoal messenger bag seitan. Shaman tbh tote bag paleo franzen crucifix enamel pin cornhole taiyaki kombucha cred banh mi. Whatever JOMO four dollar toast deep v literally adaptogen, everyday carry affogato cloud bread helvetica blackbird spyplane cold-pressed."),

		// Tell the TextArea to show the vertical scrollbar
		widget.TextAreaOpts.ShowVerticalScrollbar(),

		widget.TextAreaOpts.LinkClickedEvent(func(args *widget.LinkEventArgs) {
			fmt.Println("Link Clicked Id: ", args.Id, " value: ", args.Value, " args: ", args.Args,
				" offsetX/offsetY ", args.OffsetX, "/", args.OffsetY)
		}),
		widget.TextAreaOpts.LinkCursorEnteredEvent(func(args *widget.LinkEventArgs) {
			fmt.Println("Link Entered Id: ", args.Id, " value: ", args.Value, " args: ", args.Args,
				" offsetX/offsetY ", args.OffsetX, "/", args.OffsetY)
		}),
		widget.TextAreaOpts.LinkCursorExitedEvent(func(args *widget.LinkEventArgs) {
			fmt.Println("Link Exited Id: ", args.Id, " value: ", args.Value, " args: ", args.Args,
				" offsetX/offsetY ", args.OffsetX, "/", args.OffsetY)
		}),
	)
	result.AddChild(textarea)
	return result
}
