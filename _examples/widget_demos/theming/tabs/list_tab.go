package tabs

import (
	"fmt"

	"github.com/ebitenui/ebitenui/widget"
)

type ListEntry struct {
	id   int
	Name string
}

func NewListTab() *widget.TabBookTab {
	result := widget.NewTabBookTab(
		widget.TabBookTabOpts.Label("List"),
		widget.TabBookTabOpts.ContainerOpts(
			widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		),
	)
	// Create array of list entries
	numEntries := 20
	var id int
	entries := make([]any, 0, numEntries)
	for id = 1; id <= numEntries; id++ {
		entries = append(entries, ListEntry{id, fmt.Sprintf("Entry %d", id)})
	}

	// Construct a list. This is one of the more complicated widgets to use since
	// it is composed of multiple widget types
	list := widget.NewList(
		// Set how wide the list should be
		widget.ListOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(150, 0),
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionStart,
				StretchVertical:    true,
				Padding:            widget.NewInsetsSimple(50),
			}),
		)),
		// Set the entries in the list
		widget.ListOpts.Entries(entries),

		// Hide the horizontal slider
		widget.ListOpts.HideHorizontalSlider(),

		// This required function returns the string displayed in the list
		widget.ListOpts.EntryLabelFunc(func(e interface{}) string {
			return e.(ListEntry).Name
		}),
		// Padding for each entry
		widget.ListOpts.EntryTextPadding(widget.NewInsetsSimple(5)),
		// Text position for each entry
		widget.ListOpts.EntryTextPosition(widget.TextPositionStart, widget.TextPositionCenter),
		// This handler defines what function to run when a list item is selected.
		widget.ListOpts.EntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			entry := args.Entry.(ListEntry)
			fmt.Println("Entry Selected: ", entry)
		}),
	)

	// Add list to the root container
	result.AddChild(list)

	return result
}
