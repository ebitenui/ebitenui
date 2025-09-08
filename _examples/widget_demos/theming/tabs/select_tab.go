package tabs

import (
	"fmt"

	"github.com/ebitenui/ebitenui/widget"
)

func NewSelectTab() *widget.TabBookTab {
	result := widget.NewTabBookTab("Select",
		widget.TabBookTabOpts.ContainerOpts(
			widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		),
	)

	numEntries := 20
	entries := make([]any, 0, numEntries)
	for i := 1; i <= numEntries; i++ {
		entries = append(entries, ListEntry{i, fmt.Sprintf("Entry %d", i)})
	}
	// construct a combobox
	comboBox := widget.NewListComboButton(
		widget.ListComboButtonOpts.WidgetOpts(
			//Set the combobox's position
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionStart,
				Padding:            widget.NewInsetsSimple(20),
			}),
		),
		widget.ListComboButtonOpts.Entries(entries),

		//Define how the entry is displayed
		widget.ListComboButtonOpts.EntryLabelFunc(
			func(e any) string {
				//Button Label function
				return "Button: " + e.(ListEntry).Name
			},
			func(e any) string {
				//List Label function
				return "List: " + e.(ListEntry).Name
			}),
		//Callback when a new entry is selected
		widget.ListComboButtonOpts.EntrySelectedHandler(func(args *widget.ListComboButtonEntrySelectedEventArgs) {
			fmt.Println("Selected Entry: ", args.Entry)
		}),
	)

	// Add list to the root container
	result.AddChild(comboBox)

	return result
}
