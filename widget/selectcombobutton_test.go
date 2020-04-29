package widget

import (
	"image/color"
	"testing"

	"github.com/blizzy78/ebitenui/event"

	"github.com/matryer/is"
)

func TestSelectComboButton_SetSelectedEntry(t *testing.T) {
	is := is.New(t)

	var eventArgs *SelectComboButtonEntrySelectedEventArgs
	numEvents := 0

	b := newSelectComboButton(t,
		SelectComboButtonOpts.WithEntryLabelFunc(func(e interface{}) string {
			return "label " + e.(string)
		}),

		SelectComboButtonOpts.WithEntrySelectedHandler(func(args *SelectComboButtonEntrySelectedEventArgs) {
			eventArgs = args
			numEvents++
		}))

	entry := "foo"
	b.SetSelectedEntry(entry)
	event.FireDeferredEvents()

	is.Equal(b.SelectedEntry(), entry)
	is.Equal(eventArgs.Entry, entry)
	is.Equal(b.Label(), "label foo")

	b.SetSelectedEntry(entry)
	event.FireDeferredEvents()

	is.Equal(numEvents, 1)

	entry2 := "bar"
	b.SetSelectedEntry(entry2)
	event.FireDeferredEvents()

	is.Equal(eventArgs.PreviousEntry, entry)
}

func TestSelectComboButton_ContentVisible_Click(t *testing.T) {
	is := is.New(t)

	b := newSelectComboButton(t)

	leftMouseButtonClick(b, t)
	is.True(b.ContentVisible())

	leftMouseButtonClick(b, t)
	is.True(!b.ContentVisible())
}

func TestSelectComboButton_ContentVisible_Programmatic(t *testing.T) {
	is := is.New(t)

	b := newSelectComboButton(t)

	b.SetContentVisible(true)
	event.FireDeferredEvents()

	is.True(b.ContentVisible())

	b.SetContentVisible(false)
	event.FireDeferredEvents()

	is.True(!b.ContentVisible())
}

func newSelectComboButton(t *testing.T, opts ...SelectComboButtonOpt) *SelectComboButton {
	t.Helper()

	b := NewSelectComboButton(append(opts, []SelectComboButtonOpt{
		SelectComboButtonOpts.WithTextAndImage("", loadFont(t), &ButtonImageImage{
			Idle:     newImageEmpty(t),
			Disabled: newImageEmpty(t),
		}, &ButtonTextColor{
			Idle:     color.Transparent,
			Disabled: color.Transparent,
		}),
		SelectComboButtonOpts.WithContent(newButton(t)),
	}...)...)

	event.FireDeferredEvents()
	render(b, t)
	return b
}
