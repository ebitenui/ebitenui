package tabs

import (
	"fmt"

	"github.com/ebitenui/ebitenui/widget"
)

func NewCheckboxTab() *widget.TabBookTab {
	result := widget.NewTabBookTab("Checkbox",
		widget.TabBookTabOpts.ContainerOpts(
			widget.ContainerOpts.Layout(
				widget.NewRowLayout(
					widget.RowLayoutOpts.Direction(widget.DirectionVertical),
					widget.RowLayoutOpts.Spacing(35),
					widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(30)),
				),
			),
		),
	)
	labeledCheckBox1 := widget.NewCheckbox(
		// Set the labeled checkbox's position
		widget.CheckboxOpts.WidgetOpts(
			// Set the location of the checkbox
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  false,
			}),
			widget.WidgetOpts.MinSize(30, 30),
		),

		// Set the label
		widget.CheckboxOpts.TextLabel("Labeled Checkbox1"),
		// Set the state change handler
		widget.CheckboxOpts.StateChangedHandler(func(args *widget.CheckboxChangedEventArgs) {
			if args.State == widget.WidgetChecked {
				fmt.Println("Checkbox1 is Checked")
			} else {
				fmt.Println("Checkbox1 is Unchecked")
			}
		}),
	)
	result.AddChild(labeledCheckBox1)

	labeledCheckBox2 := widget.NewCheckbox(
		// Set the labeled checkbox's position
		widget.CheckboxOpts.WidgetOpts(
			// Set the location of the checkbox
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  false,
			}),
			// Set the minimum size of the checkbox
			widget.WidgetOpts.MinSize(30, 30),
		),

		// Set the label
		widget.CheckboxOpts.TextLabel("Labeled Checkbox2"),

		widget.CheckboxOpts.LabelFirst(),

		// Set the state change handler
		widget.CheckboxOpts.StateChangedHandler(func(args *widget.CheckboxChangedEventArgs) {
			if args.State == widget.WidgetChecked {
				fmt.Println("Checkbox2 is Checked")
			} else {
				fmt.Println("Checkbox2 is Unchecked")
			}
		}),
	)
	// Set this checkbox as Checked by default
	labeledCheckBox2.SetState(widget.WidgetChecked)

	result.AddChild(labeledCheckBox2)

	labeledCheckBox3 := widget.NewCheckbox(
		// Set the labeled checkbox's position
		widget.CheckboxOpts.WidgetOpts(
			// Set the location of the checkbox
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  false,
			}),
			// Set the minimum size of the checkbox
			widget.WidgetOpts.MinSize(30, 30),
		),

		// Set the label
		widget.CheckboxOpts.TextLabel("Labeled Tristate Checkbox"),
		// Set this checkbox to be tri-state
		widget.CheckboxOpts.TriState(),

		widget.CheckboxOpts.LabelFirst(),

		// Set the state change handler
		widget.CheckboxOpts.StateChangedHandler(func(args *widget.CheckboxChangedEventArgs) {
			if args.State == widget.WidgetChecked {
				fmt.Println("Checkbox3 is Checked")
			} else if args.State == widget.WidgetGreyed {
				fmt.Println("Checkbox3 is Greyed")
			} else {
				fmt.Println("Checkbox3 is Unchecked")
			}
		}),
	)
	// Set this checkbox as Checked by default
	labeledCheckBox3.SetState(widget.WidgetGreyed)

	result.AddChild(labeledCheckBox3)
	return result
}
