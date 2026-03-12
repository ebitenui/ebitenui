package widget

import (
	"image/color"
	"strconv"
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

func TestList_Entries_SortedOnCreate(t *testing.T) {
	is := is.New(t)

	list := newList(t,
		ListOpts.Entries([]interface{}{"third", "first", "second"}),
		ListOpts.EntrySortFunc(func(a, b any) int {
			left := a.(string)
			right := b.(string)
			switch {
			case left < right:
				return -1
			case left > right:
				return 1
			default:
				return 0
			}
		}),
		ListOpts.EntryLabelFunc(func(e interface{}) string {
			return e.(string)
		}),
	)

	is.Equal(list.Entries()[0], "first")
	is.Equal(list.Entries()[1], "second")
	is.Equal(list.Entries()[2], "third")
	is.Equal(list.buttons[0].Text().Label, "first")
	is.Equal(list.buttons[1].Text().Label, "second")
	is.Equal(list.buttons[2].Text().Label, "third")
}

func TestList_AddEntry_ResortsAndReusesButtons(t *testing.T) {
	is := is.New(t)

	list := newList(t,
		ListOpts.Entries([]interface{}{2, 4}),
		ListOpts.EntrySortFunc(func(a, b any) int {
			return a.(int) - b.(int)
		}),
		ListOpts.EntryLabelFunc(func(e interface{}) string {
			return strconv.Itoa(e.(int))
		}),
	)

	button2 := list.buttons[0]
	button4 := list.buttons[1]

	list.AddEntry(3)

	is.Equal(list.Entries()[0], 2)
	is.Equal(list.Entries()[1], 3)
	is.Equal(list.Entries()[2], 4)
	is.Equal(list.buttons[0], button2)
	is.Equal(list.buttons[2], button4)
	is.Equal(list.buttons[1].Text().Label, "3")

	children := list.listContent.Children()
	is.Equal(children[0], button2)
	is.Equal(children[1], list.buttons[1])
	is.Equal(children[2], button4)
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
