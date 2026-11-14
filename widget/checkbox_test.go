package widget

import (
	img "image"
	"testing"

	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/input"
	"github.com/matryer/is"
)

func TestCheckbox_State_Initial(t *testing.T) {
	is := is.New(t)

	c := newCheckbox(t,
		CheckboxOpts.StateChangedHandler(func(_ *CheckboxChangedEventArgs) {
			is.Fail() // event fired without previous action
		}))

	is.Equal(c.State(), WidgetUnchecked)
}

func TestCheckbox_State_SetInitial(t *testing.T) {
	is := is.New(t)

	c := newCheckbox(t,
		CheckboxOpts.InitialState(WidgetChecked),
		CheckboxOpts.StateChangedHandler(func(_ *CheckboxChangedEventArgs) {
			is.Fail() // event fired without previous action
		}))

	is.Equal(c.State(), WidgetChecked)
}

func TestCheckbox_State_SetInitialTri(t *testing.T) {
	is := is.New(t)

	c := newCheckbox(t,
		CheckboxOpts.TriState(),
		CheckboxOpts.InitialState(WidgetGreyed),
		CheckboxOpts.StateChangedHandler(func(_ *CheckboxChangedEventArgs) {
			is.Fail() // event fired without previous action
		}))

	is.Equal(c.State(), WidgetGreyed)
}

func TestCheckbox_State_SetInitialPanic(t *testing.T) {
	is := is.New(t)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	newCheckbox(t,
		CheckboxOpts.InitialState(WidgetGreyed),
		CheckboxOpts.StateChangedHandler(func(_ *CheckboxChangedEventArgs) {
			is.Fail() // event fired without previous action
		}))
}

func TestCheckbox_ChangedEvent_User(t *testing.T) {
	is := is.New(t)

	var eventArgs *CheckboxChangedEventArgs

	c := newCheckbox(t,
		CheckboxOpts.StateChangedHandler(func(args *CheckboxChangedEventArgs) {
			eventArgs = args
		}))

	leftMouseButtonClick(c, t)

	is.Equal(eventArgs.State, WidgetChecked)
	is.Equal(c.State(), WidgetChecked)
}

func TestCheckbox_Click(t *testing.T) {
	is := is.New(t)

	var eventArgs *CheckboxChangedEventArgs

	c := newCheckbox(t,
		CheckboxOpts.StateChangedHandler(func(args *CheckboxChangedEventArgs) {
			eventArgs = args
		}))

	c.Click()
	event.ExecuteDeferred()

	is.Equal(eventArgs.State, WidgetChecked)
	is.Equal(c.State(), WidgetChecked)
}

func TestCheckbox_SetState(t *testing.T) {
	is := is.New(t)

	var eventArgs *CheckboxChangedEventArgs
	numEvents := 0

	c := newCheckbox(t,
		CheckboxOpts.StateChangedHandler(func(args *CheckboxChangedEventArgs) {
			eventArgs = args
			numEvents++
		}))

	c.SetState(WidgetChecked)
	event.ExecuteDeferred()

	is.Equal(eventArgs.State, WidgetChecked)
	is.Equal(c.State(), WidgetChecked)

	c.SetState(WidgetChecked)
	event.ExecuteDeferred()

	is.Equal(numEvents, 1)
}

func TestCheckbox_State_Cycle(t *testing.T) {
	is := is.New(t)

	c := newCheckbox(t)
	leftMouseButtonClick(c, t)
	is.Equal(c.State(), WidgetChecked)
	leftMouseButtonClick(c, t)
	is.Equal(c.State(), WidgetUnchecked)
}

func TestCheckbox_State_Cycle_TriState(t *testing.T) {
	is := is.New(t)

	c := newCheckbox(t, CheckboxOpts.TriState())
	leftMouseButtonClick(c, t)
	is.Equal(c.State(), WidgetChecked)
	leftMouseButtonClick(c, t)
	is.Equal(c.State(), WidgetGreyed)
	leftMouseButtonClick(c, t)
	is.Equal(c.State(), WidgetUnchecked)
}

func newCheckbox(t *testing.T, opts ...CheckboxOpt) *Checkbox {
	t.Helper()

	c := NewCheckbox(append(opts, []CheckboxOpt{
		CheckboxOpts.Image(&CheckboxImage{
			Unchecked: newNineSliceEmpty(t),
			Checked:   newNineSliceEmpty(t),
			Greyed:    newNineSliceEmpty(t),
		})}...)...)

	event.ExecuteDeferred()
	render(c, t)
	return c
}

func TestCheckbox_HoveredImage(t *testing.T) {
	for _, state := range []WidgetState{WidgetUnchecked, WidgetChecked, WidgetGreyed} {
		is := is.New(t)
		images := newCheckboxImages(t)
		c := newCheckboxWithImages(t, images, state)
		normal, _, hovered := checkboxImagesFor(images, state)

		moveCursor(c, 10, 10, t)
		is.Equal(c.currentImage(), hovered) // cursor over the checkbox

		moveCursor(c, 30, 10, t)
		is.Equal(c.currentImage(), normal) // cursor gone

		c.Focus(true)
		is.Equal(c.currentImage(), hovered) // focused, as when tabbed to
	}
}

func TestCheckbox_DisabledImage(t *testing.T) {
	for _, state := range []WidgetState{WidgetUnchecked, WidgetChecked, WidgetGreyed} {
		is := is.New(t)
		images := newCheckboxImages(t)
		normal, disabled, _ := checkboxImagesFor(images, state)

		c := newCheckboxWithImages(t, images, state)
		moveCursor(c, 10, 10, t)
		c.GetWidget().Disabled = true
		is.Equal(c.currentImage(), disabled) // disabled, though hovered

		// Without disabled images: the normal image, not the hovered one.
		images.UncheckedDisabled, images.CheckedDisabled, images.GreyedDisabled = nil, nil, nil
		c = newCheckboxWithImages(t, images, state)
		moveCursor(c, 10, 10, t)
		c.GetWidget().Disabled = true
		is.Equal(c.currentImage(), normal)
	}
}

// Only the disabled image for the checkbox's state counts.
func TestCheckbox_DisabledImageOfItsState(t *testing.T) {
	is := is.New(t)

	images := newCheckboxImages(t)
	images.UncheckedDisabled = nil
	c := newCheckboxWithImages(t, images, WidgetGreyed)
	c.GetWidget().Disabled = true
	is.Equal(c.currentImage(), images.GreyedDisabled)

	images = newCheckboxImages(t)
	images.GreyedDisabled = nil
	c = newCheckboxWithImages(t, images, WidgetGreyed)
	c.GetWidget().Disabled = true
	is.Equal(c.currentImage(), images.Greyed)
}

// newCheckboxWithImages returns a 20x20 tri-state checkbox at the origin,
// drawn with images, in state.
func newCheckboxWithImages(t *testing.T, images *CheckboxImage, state WidgetState) *Checkbox {
	t.Helper()
	c := NewCheckbox(CheckboxOpts.TriState(), CheckboxOpts.InitialState(state), CheckboxOpts.Image(images))
	event.ExecuteDeferred()
	render(c, t)
	c.SetLocation(img.Rect(0, 0, 20, 20))
	return c
}

// newCheckboxImages returns checkbox images that are all different.
func newCheckboxImages(t *testing.T) *CheckboxImage {
	t.Helper()
	return &CheckboxImage{
		Unchecked: newNineSliceEmpty(t), UncheckedHovered: newNineSliceEmpty(t), UncheckedDisabled: newNineSliceEmpty(t),
		Checked: newNineSliceEmpty(t), CheckedHovered: newNineSliceEmpty(t), CheckedDisabled: newNineSliceEmpty(t),
		Greyed: newNineSliceEmpty(t), GreyedHovered: newNineSliceEmpty(t), GreyedDisabled: newNineSliceEmpty(t),
	}
}

// checkboxImagesFor returns the normal, disabled and hovered images for state.
func checkboxImagesFor(i *CheckboxImage, state WidgetState) (normal, disabled, hovered *image.NineSlice) {
	switch state {
	case WidgetUnchecked:
		return i.Unchecked, i.UncheckedDisabled, i.UncheckedHovered
	case WidgetChecked:
		return i.Checked, i.CheckedDisabled, i.CheckedHovered
	case WidgetGreyed:
		return i.Greyed, i.GreyedDisabled, i.GreyedHovered
	}
	return nil, nil, nil
}

// moveCursor puts the cursor at x, y and updates w, so that it sees the
// cursor enter or leave.
func moveCursor(w Updater, x, y int, t *testing.T) {
	t.Helper()
	input.SetCursorUpdater(&fixedCursor{x: x, y: y})
	t.Cleanup(func() { input.SetCursorUpdater(nil) })
	w.Update(&UpdateObject{})
	event.ExecuteDeferred()
}
