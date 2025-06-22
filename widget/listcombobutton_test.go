package widget

import (
	"image/color"
	"testing"

	"github.com/ebitenui/ebitenui/event"
	"github.com/matryer/is"
)

func TestListComboButton_SelectedEntry_Initial(t *testing.T) {
	is := is.New(t)

	entries := []interface{}{"first", "second", "third"}

	l := newListComboButton(t,
		ListComboButtonOpts.ListOpts(ListOpts.Entries(entries)),

		ListComboButtonOpts.EntrySelectedHandler(func(_ *ListComboButtonEntrySelectedEventArgs) {
			is.Fail() // event fired without previous action
		}),

		ListComboButtonOpts.EntryLabelFunc(
			func(e interface{}) string {
				result, _ := e.(string)
				return "label " + result
			}, func(e interface{}) string {
				result, _ := e.(string)
				return result
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
		ListComboButtonOpts.ListOpts(ListOpts.Entries(entries)),

		ListComboButtonOpts.EntrySelectedHandler(func(args *ListComboButtonEntrySelectedEventArgs) {
			eventArgs = args
			numEvents++
		}),

		ListComboButtonOpts.EntryLabelFunc(
			func(e interface{}) string {
				result, _ := e.(string)
				return "label " + result
			}, func(e interface{}) string {
				result, _ := e.(string)
				return result
			}))

	l.SetSelectedEntry(entries[1])
	event.ExecuteDeferred()

	is.Equal(l.SelectedEntry(), entries[1])
	is.Equal(eventArgs.Entry, entries[1])
	is.Equal(eventArgs.PreviousEntry, entries[0])
	is.Equal(l.Label(), "label second")

	l.SetSelectedEntry(entries[1])
	event.ExecuteDeferred()
	is.Equal(numEvents, 1)
}

func TestListComboButton_EntrySelectedEvent_User(t *testing.T) {
	is := is.New(t)

	entries := []interface{}{"first", "second", "third"}

	var eventArgs *ListComboButtonEntrySelectedEventArgs
	numEvents := 0

	l := newListComboButton(t,
		ListComboButtonOpts.ListOpts(ListOpts.Entries(entries)),

		ListComboButtonOpts.EntrySelectedHandler(func(args *ListComboButtonEntrySelectedEventArgs) {
			eventArgs = args
			numEvents++
		}),

		ListComboButtonOpts.EntryLabelFunc(
			func(e interface{}) string {
				result, _ := e.(string)
				return "label " + result
			}, func(e interface{}) string {
				result, _ := e.(string)
				return result
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
		ListComboButtonOpts.ListOpts(ListOpts.Entries(entries)),

		ListComboButtonOpts.EntryLabelFunc(
			func(e interface{}) string {
				result, _ := e.(string)
				return result
			}, func(e interface{}) string {
				result, _ := e.(string)
				return result
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
		ListComboButtonOpts.ListOpts(ListOpts.Entries(entries)),

		ListComboButtonOpts.EntryLabelFunc(
			func(e interface{}) string {
				result, _ := e.(string)
				return result
			}, func(e interface{}) string {
				result, _ := e.(string)
				return result
			}))

	l.SetContentVisible(true)
	is.True(l.ContentVisible())

	l.SetContentVisible(false)
	is.True(!l.ContentVisible())
}

func newListComboButton(t *testing.T, opts ...ListComboButtonOpt) *ListComboButton {
	t.Helper()

	l := NewListComboButton(append(opts, []ListComboButtonOpt{
		ListComboButtonOpts.SelectComboButtonOpts(SelectComboButtonOpts.ComboButtonOpts(ComboButtonOpts.ButtonOpts(ButtonOpts.Image(&ButtonImage{
			Idle:    newNineSliceEmpty(t),
			Pressed: newNineSliceEmpty(t),
		})))),
		ListComboButtonOpts.ListOpts(
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
				},
			}),
			ListOpts.EntryColor(&ListEntryColor{
				Unselected:                 color.Transparent,
				Selected:                   color.Transparent,
				DisabledUnselected:         color.Transparent,
				DisabledSelected:           color.Transparent,
				SelectedBackground:         color.Transparent,
				DisabledSelectedBackground: color.Transparent,
			}),
			ListOpts.EntryFontFace(loadFont(t)),
		),
		ListComboButtonOpts.Text(loadFont(t), &GraphicImage{
			Idle:     newImageEmpty(t),
			Disabled: newImageEmpty(t),
		}, &ButtonTextColor{
			Idle:     color.Transparent,
			Disabled: color.Transparent,
		}),
	}...)...)

	event.ExecuteDeferred()
	render(l, t)
	return l
}

func listComboButtonContentList(l *ListComboButton) *List {
	result, _ := l.button.button.content.(*List)
	return result
}
