package widget

import (
	"image/color"
	"testing"

	"github.com/ebitenui/ebitenui/event"
	"github.com/matryer/is"
)

func TestList_SelectedEntry_Initial(t *testing.T) {
	is := is.New(t)

	entries := []interface{}{"first", "second", "third"}

	list := newList(t,
		ListOpts.Entries(entries),

		ListOpts.EntryLabelFunc(func(e interface{}) string {
			result, _ := e.(string)
			return result
		}),

		ListOpts.EntrySelectedHandler(func(_ *ListEntrySelectedEventArgs) {
			is.Fail() // event fired without previous action
		}))

	is.Equal(list.SelectedEntry(), nil)
}

func TestList_NoSliderOpts(t *testing.T) {
	entries := []interface{}{"first", "second", "third"}
	_ = NewList(
		ListOpts.Entries(entries),

		ListOpts.EntryLabelFunc(func(e interface{}) string {
			result, _ := e.(string)
			return result
		}),

		ListOpts.EntrySelectedHandler(func(_ *ListEntrySelectedEventArgs) {
		}),
		ListOpts.ScrollContainerImage(&ScrollContainerImage{
			Idle:     newNineSliceEmpty(t),
			Disabled: newNineSliceEmpty(t),
			Mask:     newNineSliceEmpty(t),
		}),

		ListOpts.HideHorizontalSlider(),
		ListOpts.HideVerticalSlider(),

		ListOpts.EntryFontFace(loadFont(t)),

		ListOpts.EntryColor(&ListEntryColor{
			Unselected:                 color.Transparent,
			Selected:                   color.Transparent,
			DisabledUnselected:         color.Transparent,
			DisabledSelected:           color.Transparent,
			SelectedBackground:         color.Transparent,
			DisabledSelectedBackground: color.Transparent,
		}),
	)
}

func TestList_SetSelectedEntry(t *testing.T) {
	is := is.New(t)

	entries := []interface{}{"first", "second", "third"}

	var eventArgs *ListEntrySelectedEventArgs
	numEvents := 0

	list := newList(t,
		ListOpts.Entries(entries),

		ListOpts.EntryLabelFunc(func(e interface{}) string {
			result, _ := e.(string)
			return result
		}),

		ListOpts.EntrySelectedHandler(func(args *ListEntrySelectedEventArgs) {
			eventArgs = args
			numEvents++
		}))

	list.SetSelectedEntry(entries[1])
	event.ExecuteDeferred()

	is.Equal(eventArgs.Entry, entries[1])
	is.Equal(list.SelectedEntry(), entries[1])

	list.SetSelectedEntry(entries[1])
	event.ExecuteDeferred()

	is.Equal(numEvents, 1)
}

func TestList_UpdateEntry(t *testing.T) {
	is := is.New(t)

	type listEntry struct {
		label string
	}

	entries := []interface{}{
		&listEntry{label: "first"},
		&listEntry{label: "second"},
		&listEntry{label: "third"},
	}

	list := newList(t,
		ListOpts.Entries(entries),

		ListOpts.EntryLabelFunc(func(e interface{}) string {
			return e.(*listEntry).label
		}),
	)

	oldButton := list.buttons[1]
	entry := entries[1].(*listEntry)

	is.Equal(list.buttons[1].Text().Label, "second")

	entry.label = "updated"

	is.Equal(list.buttons[1].Text().Label, "second")

	list.UpdateEntry(entry)

	is.True(list.buttons[1] != oldButton)
	is.Equal(list.buttons[1].Text().Label, "updated")
}

func TestList_EntrySelectedEvent_User(t *testing.T) {
	is := is.New(t)

	entries := []interface{}{"first", "second", "third"}

	var eventArgs *ListEntrySelectedEventArgs
	numEvents := 0

	list := newList(t,
		ListOpts.Entries(entries),

		ListOpts.EntryLabelFunc(func(e interface{}) string {
			result, _ := e.(string)
			return result
		}),

		ListOpts.EntrySelectedHandler(func(args *ListEntrySelectedEventArgs) {
			eventArgs = args
			numEvents++
		}))

	leftMouseButtonClick(list.buttons[1], t)

	is.Equal(eventArgs.Entry, entries[1])
	is.Equal(list.SelectedEntry(), entries[1])

	leftMouseButtonClick(list.buttons[1], t)

	is.Equal(numEvents, 1)
}

func TestList_EntrySelectedEvent_User_AllowReselect(t *testing.T) {
	is := is.New(t)

	entries := []interface{}{"first", "second", "third"}

	var eventArgs *ListEntrySelectedEventArgs
	numEvents := 0

	list := newList(t,
		ListOpts.Entries(entries),

		ListOpts.EntryLabelFunc(func(e interface{}) string {
			result, _ := e.(string)
			return result
		}),

		ListOpts.AllowReselect(),

		ListOpts.EntrySelectedHandler(func(args *ListEntrySelectedEventArgs) {
			eventArgs = args
			numEvents++
		}))

	leftMouseButtonClick(list.buttons[1], t)

	is.Equal(eventArgs.Entry, entries[1])
	is.Equal(list.SelectedEntry(), entries[1])

	leftMouseButtonClick(list.buttons[1], t)

	is.Equal(eventArgs.Entry, entries[1])
	is.Equal(eventArgs.PreviousEntry, entries[1])
	is.Equal(list.SelectedEntry(), entries[1])

	is.Equal(numEvents, 2)
}

func newList(t *testing.T, opts ...ListOpt) *List {
	t.Helper()

	l := NewList(append(opts, []ListOpt{
		ListOpts.ScrollContainerImage(&ScrollContainerImage{
			Idle:     newNineSliceEmpty(t),
			Disabled: newNineSliceEmpty(t),
			Mask:     newNineSliceEmpty(t),
		}),

		ListOpts.SliderParams(&SliderParams{
			TrackImage: &SliderTrackImage{},
			HandleImage: &ButtonImage{
				Idle:    newNineSliceEmpty(t),
				Pressed: newNineSliceEmpty(t),
			}}),

		ListOpts.EntryFontFace(loadFont(t)),

		ListOpts.EntryColor(&ListEntryColor{
			Unselected:                 color.Transparent,
			Selected:                   color.Transparent,
			DisabledUnselected:         color.Transparent,
			DisabledSelected:           color.Transparent,
			SelectedBackground:         color.Transparent,
			DisabledSelectedBackground: color.Transparent,
		}),
	}...)...)

	event.ExecuteDeferred()
	render(l, t)
	return l
}

func listEntryButtons(l *List) []*Button {
	return l.buttons
}
