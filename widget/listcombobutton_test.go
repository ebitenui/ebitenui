package widget

import (
	"image/color"
	"testing"

	"github.com/blizzy78/ebitenui/event"

	"github.com/matryer/is"
)

func TestListComboButton_SelectedEntry_Initial(t *testing.T) {
	is := is.New(t)

	entries := []interface{}{"first", "second", "third"}

	l := newListComboButton(t,
		ListComboButtonOpts.WithListOpt(ListOpts.WithEntries(entries)),

		ListComboButtonOpts.WithEntrySelectedHandler(func(args *ListComboButtonEntrySelectedEventArgs) {
			is.Fail() // event fired without previous action
		}),

		ListComboButtonOpts.WithEntryLabelFunc(
			func(e interface{}) string {
				return "label " + e.(string)
			}, func(e interface{}) string {
				return e.(string)
			}))

	is.Equal(l.SelectedEntry(), entries[0])
	is.Equal(l.Label(), "label first")
}

func TestListComboButton_SetSelectedEntry(t *testing.T) {
	is := is.New(t)

	entries := []interface{}{"first", "second", "third"}

	var eventArgs *ListComboButtonEntrySelectedEventArgs
	numEvents := 0

	l := newListComboButton(t,
		ListComboButtonOpts.WithListOpt(ListOpts.WithEntries(entries)),

		ListComboButtonOpts.WithEntrySelectedHandler(func(args *ListComboButtonEntrySelectedEventArgs) {
			eventArgs = args
			numEvents++
		}),

		ListComboButtonOpts.WithEntryLabelFunc(
			func(e interface{}) string {
				return "label " + e.(string)
			}, func(e interface{}) string {
				return e.(string)
			}))

	l.SetSelectedEntry(entries[1])
	event.FireDeferredEvents()

	is.Equal(l.SelectedEntry(), entries[1])
	is.Equal(eventArgs.Entry, entries[1])
	is.Equal(eventArgs.PreviousEntry, entries[0])
	is.Equal(l.Label(), "label second")

	l.SetSelectedEntry(entries[1])
	event.FireDeferredEvents()
	is.Equal(numEvents, 1)
}

func TestListComboButton_EntrySelectedEvent_User(t *testing.T) {
	is := is.New(t)

	entries := []interface{}{"first", "second", "third"}

	var eventArgs *ListComboButtonEntrySelectedEventArgs
	numEvents := 0

	l := newListComboButton(t,
		ListComboButtonOpts.WithListOpt(ListOpts.WithEntries(entries)),

		ListComboButtonOpts.WithEntrySelectedHandler(func(args *ListComboButtonEntrySelectedEventArgs) {
			eventArgs = args
			numEvents++
		}),

		ListComboButtonOpts.WithEntryLabelFunc(
			func(e interface{}) string {
				return "label " + e.(string)
			}, func(e interface{}) string {
				return e.(string)
			}))

	l.SetContentVisible(true)
	render(l, t)

	leftMouseButtonClick(listEntryButtons(listComboButtonContentList(l))[1], t)

	is.Equal(l.SelectedEntry(), entries[1])
	is.Equal(eventArgs.Entry, entries[1])
	is.Equal(eventArgs.PreviousEntry, entries[0])
	is.Equal(l.Label(), "label second")

	l.SetContentVisible(true)
	render(l, t)

	leftMouseButtonClick(listEntryButtons(listComboButtonContentList(l))[1], t)

	is.Equal(numEvents, 1)
}

func TestListComboButton_ContentVisible_Click(t *testing.T) {
	is := is.New(t)

	entries := []interface{}{"first", "second", "third"}

	l := newListComboButton(t,
		ListComboButtonOpts.WithListOpt(ListOpts.WithEntries(entries)),

		ListComboButtonOpts.WithEntryLabelFunc(
			func(e interface{}) string {
				return e.(string)
			}, func(e interface{}) string {
				return e.(string)
			}))

	leftMouseButtonClick(l, t)
	is.True(l.ContentVisible())

	leftMouseButtonClick(l, t)
	is.True(!l.ContentVisible())
}

func TestListComboButton_ContentVisible_Programmatic(t *testing.T) {
	is := is.New(t)

	entries := []interface{}{"first", "second", "third"}

	l := newListComboButton(t,
		ListComboButtonOpts.WithListOpt(ListOpts.WithEntries(entries)),

		ListComboButtonOpts.WithEntryLabelFunc(
			func(e interface{}) string {
				return e.(string)
			}, func(e interface{}) string {
				return e.(string)
			}))

	l.SetContentVisible(true)
	is.True(l.ContentVisible())

	l.SetContentVisible(false)
	is.True(!l.ContentVisible())
}

func newListComboButton(t *testing.T, opts ...ListComboButtonOpt) *ListComboButton {
	t.Helper()

	l := NewListComboButton(append(opts, []ListComboButtonOpt{
		ListComboButtonOpts.WithListOpt(ListOpts.WithScrollContainerOpt(ScrollContainerOpts.WithImage(&ScrollContainerImage{
			Idle:     newNineSliceEmpty(t),
			Disabled: newNineSliceEmpty(t),
			Mask:     newNineSliceEmpty(t),
		}))),
		ListComboButtonOpts.WithListOpt(ListOpts.WithEntryColor(&ListEntryColor{
			Unselected:                 color.Transparent,
			Selected:                   color.Transparent,
			DisabledUnselected:         color.Transparent,
			DisabledSelected:           color.Transparent,
			SelectedBackground:         color.Transparent,
			DisabledSelectedBackground: color.Transparent,
		})),
		ListComboButtonOpts.WithListOpt(ListOpts.WithEntryFontFace(loadFont(t))),
		ListComboButtonOpts.WithText(loadFont(t), &ButtonImageImage{
			Idle:     newImageEmpty(t),
			Disabled: newImageEmpty(t),
		}, &ButtonTextColor{
			Idle:     color.Transparent,
			Disabled: color.Transparent,
		}),
	}...)...)

	event.FireDeferredEvents()
	render(l, t)
	return l
}

func listComboButtonContentList(l *ListComboButton) *List {
	return l.button.button.content.(*List)
}
