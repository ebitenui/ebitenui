package widget

import (
	"testing"

	"github.com/blizzy78/ebitenui/event"

	"github.com/matryer/is"
)

func TestComboButton_ContentVisible_Click(t *testing.T) {
	is := is.New(t)

	b := newComboButton(t)

	leftMouseButtonClick(b, t)
	is.True(b.ContentVisible)

	leftMouseButtonClick(b, t)
	is.True(!b.ContentVisible)
}

func newComboButton(t *testing.T, opts ...ComboButtonOpt) *ComboButton {
	t.Helper()

	b := NewComboButton(append(opts, ComboButtonOpts.WithContent(newButton(t)))...)
	event.FireDeferredEvents()
	render(b, t)
	return b
}
