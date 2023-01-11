package widget

import (
	"image/color"
	"testing"

	"github.com/ebitenui/ebitenui/event"
	"github.com/matryer/is"
)

func TestLabel_SetLabel(t *testing.T) {
	is := is.New(t)

	l := newLabel(t)

	l.Label = "foo"
	render(l, t)

	is.Equal(labelText(l).Label, "foo")
}

func TestLabel_SetDisabled_Color(t *testing.T) {
	is := is.New(t)

	l := newLabel(t)

	l.GetWidget().Disabled = true
	render(l, t)

	is.Equal(labelText(l).Color, color.Black)
}

func newLabel(t *testing.T, opts ...LabelOpt) *Label {
	t.Helper()

	l := NewLabel(append(opts, LabelOpts.Text("", loadFont(t), &LabelColor{
		Idle:     color.White,
		Disabled: color.Black,
	}))...)
	event.ExecuteDeferred()
	render(l, t)
	return l
}

func labelText(l *Label) *Text {
	return l.text
}
