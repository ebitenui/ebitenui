package widget

import (
	"testing"

	internalevent "github.com/blizzy78/ebitenui/internal/event"
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

	b := NewComboButton(append(opts, []ComboButtonOpt{
		ComboButtonOpts.ButtonOpts(ButtonOpts.Image(&ButtonImage{
			Idle: newNineSliceEmpty(t),
		})),
		ComboButtonOpts.Content(newButton(t)),
	}...)...)
	internalevent.ExecuteDeferredActions()
	render(b, t)
	return b
}
