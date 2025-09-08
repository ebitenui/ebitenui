package tabs

import (
	"github.com/ebitenui/ebitenui/widget"
)

func NewButtonTab() *widget.TabBookTab {
	result := widget.NewTabBookTab("Button",
		widget.TabBookTabOpts.ContainerOpts(
			widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		),
	)
	var button *widget.Button
	button = widget.NewButton(
		// set general widget options
		widget.ButtonOpts.WidgetOpts(
			// instruct the container's anchor layout to center the button both horizontally and vertically.
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),

		// specify the button's text, the font face, and the color.
		widget.ButtonOpts.TextLabel("Hello, World!"),

		// Move the text down and right on press
		widget.ButtonOpts.PressedHandler(func(args *widget.ButtonPressedEventArgs) {
			button.Text().SetPadding(&widget.Insets{Top: 1, Bottom: -1})
			button.GetWidget().CustomData = true
		}),
		// Move the text back to start on press released
		widget.ButtonOpts.ReleasedHandler(func(args *widget.ButtonReleasedEventArgs) {
			button.Text().SetPadding(&widget.Insets{})
			button.GetWidget().CustomData = false
		}),

		// add a handler that reacts to clicking the button.
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			println("button clicked")
		}),

		// add a handler that reacts to entering the button with the cursor
		widget.ButtonOpts.CursorEnteredHandler(func(args *widget.ButtonHoverEventArgs) {
			println("cursor entered button: entered =", args.Entered, "offsetX =", args.OffsetX, "offsetY =", args.OffsetY)
			// If we moved the Text because we clicked on this button previously, move the text down and right
			if button.GetWidget().CustomData == true {
				button.Text().SetPadding(&widget.Insets{Top: 1, Bottom: -1})
			}
		}),

		// add a handler that reacts to entering the button with the cursor.
		widget.ButtonOpts.CursorEnteredHandler(func(args *widget.ButtonHoverEventArgs) {
			println("cursor entered button: entered =", args.Entered, "offsetX =", args.OffsetX, "offsetY =", args.OffsetY)
		}),

		// add a handler that reacts to moving the cursor on the button.
		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) {
			println("cursor moved on button: entered =", args.Entered, "offsetX =", args.OffsetX, "offsetY =", args.OffsetY, "diffX =", args.DiffX, "diffY =", args.DiffY)
		}),

		// add a handler that reacts to exiting the button with the cursor.
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) {
			println("cursor exited button: entered =", args.Entered, "offsetX =", args.OffsetX, "offsetY =", args.OffsetY)
			// Reset the Text inset if the cursor is no longer over the button
			button.Text().SetPadding(&widget.Insets{})
		}),
	)

	// add the button as a child of the container.
	result.AddChild(button)

	return result
}
