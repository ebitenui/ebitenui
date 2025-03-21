package widget

import (
	"image/color"
	"testing"

	"github.com/ebitenui/ebitenui/event"
	"github.com/matryer/is"
)

func TestLabeledCheckbox_SetState_User(t *testing.T) {
	is := is.New(t)

	l := newLabeledCheckbox(t)
	leftMouseButtonClick(labeledCheckboxLabel(l), t)
	is.Equal(l.Checkbox().State(), WidgetChecked)

	l2 := newLabeledCheckbox(t)
	l2.SetState(WidgetChecked)
	leftMouseButtonClick(labeledCheckboxLabel(l2), t)
	is.Equal(l2.Checkbox().State(), WidgetUnchecked)
}

func TestLabeledCheckbox_Click(t *testing.T) {
	is := is.New(t)

	c := newLabeledCheckbox(t)

	c.Click()
	event.ExecuteDeferred()

	is.Equal(c.Checkbox().State(), WidgetChecked)
}

func newLabeledCheckbox(t *testing.T, opts ...LabeledCheckboxOpt) *LabeledCheckbox {
	t.Helper()

	l := NewLabeledCheckbox(append(opts, []LabeledCheckboxOpt{
		LabeledCheckboxOpts.CheckboxOpts(
			CheckboxOpts.ButtonOpts(ButtonOpts.Image(&ButtonImage{
				Idle:    newNineSliceEmpty(t),
				Pressed: newNineSliceEmpty(t),
			})),
			CheckboxOpts.Image(&CheckboxGraphicImage{
				Unchecked: &GraphicImage{
					Idle: newImageEmpty(t),
				},
				Checked: &GraphicImage{
					Idle: newImageEmpty(t),
				},
			}),
		),
		LabeledCheckboxOpts.LabelOpts(LabelOpts.Text("", loadFont(t), &LabelColor{
			Idle: color.White,
		})),
	}...)...)
	event.ExecuteDeferred()
	render(l, t)
	return l
}

func labeledCheckboxLabel(l *LabeledCheckbox) *Label {
	return l.label
}
