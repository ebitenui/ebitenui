package widget

import (
	"image/color"
	"testing"

	"github.com/matryer/is"
	"github.com/mcarpenter622/ebitenui/event"
)

func TestLabeledCheckbox_SetState_User(t *testing.T) {
	is := is.New(t)

	l := newLabeledCheckbox(t)
	leftMouseButtonClick(labeledCheckboxLabel(l), t)

	is.Equal(l.Checkbox().State(), WidgetChecked)
}

func newLabeledCheckbox(t *testing.T, opts ...LabeledCheckboxOpt) *LabeledCheckbox {
	t.Helper()

	l := NewLabeledCheckbox(append(opts, []LabeledCheckboxOpt{
		LabeledCheckboxOpts.CheckboxOpts(
			CheckboxOpts.ButtonOpts(ButtonOpts.Image(&ButtonImage{
				Idle: newNineSliceEmpty(t),
			})),
			CheckboxOpts.Image(&CheckboxGraphicImage{
				Unchecked: &ButtonImageImage{
					Idle: newImageEmpty(t),
				},
				Checked: &ButtonImageImage{
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
