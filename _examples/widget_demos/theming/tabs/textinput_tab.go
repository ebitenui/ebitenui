package tabs

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/ebitenui/ebitenui/widget"
)

func NewTextInputTab() *widget.TabBookTab {
	result := widget.NewTabBookTab("Text Input",
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(10),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(20)),
		)),
	)

	// construct a standard textinput widget
	standardTextInput := widget.NewTextInput(
		widget.TextInputOpts.WidgetOpts(
			// Set the layout information to center the textbox in the parent
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  true,
				MaxWidth: 400,
			}),
		),

		// This text is displayed if the input is empty
		widget.TextInputOpts.Placeholder("Standard Textbox"),

		// This is called when the user hits the "Enter" key.
		// There are other options that can configure this behavior.
		widget.TextInputOpts.SubmitHandler(func(args *widget.TextInputChangedEventArgs) {
			fmt.Println("Text Submitted: ", args.InputText)
		}),

		// This is called whenver there is a change to the text
		widget.TextInputOpts.ChangedHandler(func(args *widget.TextInputChangedEventArgs) {
			fmt.Println("Text Changed: ", args.InputText)
		}),
	)

	result.AddChild(standardTextInput)

	// construct a disabled textinput widget
	disabledTextInput := widget.NewTextInput(
		widget.TextInputOpts.WidgetOpts(
			// Set the layout information to center the textbox in the parent
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  true,
				MaxWidth: 400,
			}),
		),
		// This text is displayed if the input is empty
		widget.TextInputOpts.Placeholder("Disabled Textbox"),

		// This is called when the user hits the "Enter" key.
		// There are other options that can configure this behavior.
		widget.TextInputOpts.SubmitHandler(func(args *widget.TextInputChangedEventArgs) {
			fmt.Println("Text Submitted: ", args.InputText)
		}),

		// This is called whenver there is a change to the text
		widget.TextInputOpts.ChangedHandler(func(args *widget.TextInputChangedEventArgs) {
			fmt.Println("Text Changed: ", args.InputText)
		}),
	)
	disabledTextInput.GetWidget().Disabled = true
	result.AddChild(disabledTextInput)

	// construct a secure textinput widget
	secureTextInput := widget.NewTextInput(
		widget.TextInputOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  true,
				MaxWidth: 400,
			}),
		),

		// This parameter indicates that the inputted text should be hidden
		widget.TextInputOpts.Secure(true),

		widget.TextInputOpts.Placeholder("Password Textbox"),
		widget.TextInputOpts.SubmitHandler(func(args *widget.TextInputChangedEventArgs) {
			fmt.Println("Text Submitted: ", args.InputText)
		}),
		widget.TextInputOpts.ChangedHandler(func(args *widget.TextInputChangedEventArgs) {
			fmt.Println("Text Changed: ", args.InputText)
		}),
	)

	result.AddChild(secureTextInput)

	maxLenTextInput := widget.NewTextInput(
		widget.TextInputOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  true,
				MaxWidth: 400,
			}),
		),

		widget.TextInputOpts.Placeholder("Max length (5) Textbox"),

		// This method is called whenever there is a text change.
		// It allows the developer to allow or deny a change.
		// In this case we are limiting the string to 5 runes.
		// The first return parameter is whether or not to accept the text as is.
		// The second return parameter is what to replace the text with if it is not accepted (optional)
		widget.TextInputOpts.Validation(func(newInputText string) (bool, *string) {
			if utf8.RuneCountInString(newInputText) > 5 {
				return false, nil
			}
			return true, nil
		}),
		widget.TextInputOpts.SubmitHandler(func(args *widget.TextInputChangedEventArgs) {
			fmt.Println("Text Submitted: ", args.InputText)
		}),
		widget.TextInputOpts.ChangedHandler(func(args *widget.TextInputChangedEventArgs) {
			fmt.Println("Text Changed: ", args.InputText)
		}),
	)
	// This will do nothing because the validation above prevents this from being set.
	maxLenTextInput.SetText("123456")
	result.AddChild(maxLenTextInput)

	allCapsTextInput := widget.NewTextInput(
		widget.TextInputOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  true,
				MaxWidth: 400,
			}),
		),

		widget.TextInputOpts.Placeholder("All Caps Textbox"),

		// This method is called whenever there is a text change.
		// It allows the developer to allow or deny a change.
		// In this case we are forcing the string to be all caps.
		// The first return parameter is whether or not to accept the text as is.
		// The second return parameter is what to replace the text with if it is not accepted (optional)
		widget.TextInputOpts.Validation(func(newInputText string) (bool, *string) {
			newInputText = strings.ToUpper(newInputText)
			return false, &newInputText
		}),
		widget.TextInputOpts.SubmitHandler(func(args *widget.TextInputChangedEventArgs) {
			fmt.Println("Text Submitted: ", args.InputText)
		}),
		widget.TextInputOpts.ChangedHandler(func(args *widget.TextInputChangedEventArgs) {
			fmt.Println("Text Changed: ", args.InputText)
		}),
	)
	// This will show in all caps due to validation function above
	allCapsTextInput.SetText("Hello World")
	result.AddChild(allCapsTextInput)
	return result
}
