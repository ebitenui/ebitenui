package widget

import (
	img "image"
	"image/color"
	"testing"

	"github.com/matryer/is"
)

func newDialog(t *testing.T, opts ...DialogOpt) (*Dialog, *[]DialogChoiceEventArgs, *bool) {
	t.Helper()

	var choices []DialogChoiceEventArgs
	closed := false
	ns := newNineSliceEmpty(t)
	d := NewDialog(append([]DialogOpt{
		DialogOpts.Title("Dust storm"),
		DialogOpts.Message("A storm is coming over the landing site. Recall the rover?"),
		DialogOpts.Choices("Recall", "Keep going"),
		DialogOpts.MaxWidth(200),
		DialogOpts.TitleText(loadFont(t), color.White),
		DialogOpts.MessageText(loadFont(t), color.White),
		DialogOpts.ButtonParams(&ButtonParams{
			Image:     &ButtonImage{Idle: ns, Pressed: ns},
			TextFace:  loadFont(t),
			TextColor: &ButtonTextColor{Idle: color.White},
		}),
		DialogOpts.ChoiceHandler(func(args *DialogChoiceEventArgs) { choices = append(choices, *args) }),
	}, opts...)...)
	d.Window().SetCloseFunction(func() { closed = true })
	return d, &choices, &closed
}

// The dialog holds its title, its message wrapped at MaxWidth, and a
// button per choice under them, in order.
func TestDialog_Contents(t *testing.T) {
	is := is.New(t)

	d, _, _ := newDialog(t)
	render(d.contents, t)

	is.Equal(len(d.contents.children), 4)
	message, ok := d.contents.children[1].(*Text)
	is.True(ok)
	is.True(len(message.measurements.processedLines) > 1) // wrapped
	_, ok = d.contents.children[3].(*Button)
	is.True(ok)
}

// A choice is reported with its index and label, and closes the dialog.
func TestDialog_Choice(t *testing.T) {
	is := is.New(t)

	d, choices, closed := newDialog(t)
	render(d.contents, t)
	leftMouseButtonClick(d.contents.children[3].(*Button), t) //nolint:forcetypeassert // the second choice

	is.Equal(len(*choices), 1)
	is.Equal((*choices)[0].Choice, 1)
	is.Equal((*choices)[0].Label, "Keep going")
	is.True(*closed)
}

// With KeepOpen, a choice leaves the dialog open.
func TestDialog_KeepOpen(t *testing.T) {
	is := is.New(t)

	d, choices, closed := newDialog(t, DialogOpts.KeepOpen())
	render(d.contents, t)
	leftMouseButtonClick(d.contents.children[2].(*Button), t) //nolint:forcetypeassert // the first choice

	is.Equal(len(*choices), 1)
	is.True(!*closed)
}

// A dialog's window takes the options given through WindowOpts: with a
// transition, a choice starts it on its way out.
func TestDialog_WindowOpts(t *testing.T) {
	is := is.New(t)

	d, _, closed := newDialog(t, DialogOpts.WindowOpts(
		WindowOpts.Transition(WindowTransition{Duration: 4, Out: WindowMotion{Offset: img.Pt(0, 20)}}),
	))
	render(d.contents, t)
	leftMouseButtonClick(d.contents.children[2].(*Button), t) //nolint:forcetypeassert // the first choice

	is.True(d.Window().Closing()) // easing out
	is.True(!*closed)             // not yet removed
}

// CenterIn places the dialog at its preferred size in the middle of an area.
func TestDialog_CenterIn(t *testing.T) {
	is := is.New(t)

	d, _, _ := newDialog(t)
	d.CenterIn(img.Rect(100, 50, 500, 450))

	w, h := d.contents.PreferredSize()
	r := d.Window().GetContainer().GetWidget().Rect
	is.Equal(r.Size(), img.Pt(w, h))
	is.Equal(r.Min.Add(r.Max), img.Pt(600, 500)) // centered: twice the area's middle
}
