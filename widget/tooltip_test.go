package widget

import (
	img "image"
	"testing"

	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/input"
	"github.com/matryer/is"
)

// A list combo button's own widget gets its updates: a tooltip on it shows
// while the pointer rests on the button.
func TestListComboButton_ToolTip(t *testing.T) {
	is := is.New(t)

	c := useTestCursor(t)
	input.SetupInputLayersWithDeferred(nil) // drops layers other tests left
	tip := NewToolTip(ToolTipOpts.Content(NewContainer()))
	l := newListComboButton(t, ListComboButtonOpts.Entries([]any{"a"}), ListComboButtonOpts.EntryLabelFunc(label, label),
		ListComboButtonOpts.WidgetOpts(WidgetOpts.ToolTip(tip)))
	l.SetLocation(img.Rect(0, 0, 100, 30))
	shown := 0
	l.GetWidget().ToolTipEvent.AddHandler(func(a any) {
		if a.(*WidgetToolTipEventArgs).Show { //nolint:forcetypeassert // the event's own args
			shown++
		}
	})

	c.x, c.y = 50, 15
	for range 2 {
		l.Update(&UpdateObject{})
		event.ExecuteDeferred()
	}
	is.Equal(shown, 1)
}
