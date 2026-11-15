package widget

import (
	"testing"

	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/input"
	"github.com/ebitenui/ebitenui/internal/input/inputtest"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/matryer/is"
)

// given returns actions giving acts.
func given(acts ...Action) *Actions {
	a := &Actions{}
	for _, act := range acts {
		a.Give(act)
	}
	return a
}

func label(e any) string { return e.(string) } //nolint:forcetypeassert // the tests' entries are strings

// A focused list moves through its entries with Up and Down and leaves
// those at its ends, so a gamepad moves the focus on from there; Activate
// selects the entry it is on. Unfocused, or with its default keys off, it
// takes nothing.
func TestList_TakesActions(t *testing.T) {
	is := is.New(t)

	l := newList(t, ListOpts.Entries([]any{"a", "b", "c"}), ListOpts.EntryLabelFunc(label))
	a := given(ActionDown)
	l.Update(&UpdateObject{Actions: a})
	is.True(a.Left(ActionDown)) // not focused

	l.Focus(true)
	for i, s := range []struct {
		act   Action
		taken bool
		at    int
	}{
		{ActionUp, false, 0}, // the first entry: left
		{ActionDown, true, 1},
		{ActionDown, true, 2},
		{ActionDown, false, 2}, // the last: left
		{ActionUp, true, 1},
	} {
		a := given(s.act)
		l.Update(&UpdateObject{Actions: a})
		is.Equal([]any{i, !a.Left(s.act), l.focusIndex}, []any{i, s.taken, s.at})
	}

	a = given(ActionActivate)
	l.Update(&UpdateObject{Actions: a})
	event.ExecuteDeferred()
	is.True(!a.Left(ActionActivate))
	is.Equal(l.SelectedEntry(), "b")

	off := newList(t, ListOpts.Entries([]any{"a", "b"}), ListOpts.EntryLabelFunc(label), ListOpts.DisableDefaultKeys(true))
	off.Focus(true)
	a = given(ActionDown, ActionActivate)
	off.Update(&UpdateObject{Actions: a})
	is.True(a.Left(ActionDown) && a.Left(ActionActivate)) // default keys off
}

// A focused slider moves its handle a page with the two directions of its
// orientation, also at its ends, and leaves the other two; unfocused, or
// with its default keys off, it takes nothing.
func TestSlider_TakesActions(t *testing.T) {
	is := is.New(t)

	for _, c := range []struct {
		dir               Direction
		less, more, other Action
	}{
		{DirectionHorizontal, ActionLeft, ActionRight, ActionUp},
		{DirectionVertical, ActionUp, ActionDown, ActionRight},
	} {
		s := newSlider(t, SliderOpts.Direction(c.dir), SliderOpts.MinMax(0, 20), SliderOpts.InitialCurrent(0),
			SliderOpts.PageSizeFunc(func() int { return 10 }))
		a := given(c.more)
		s.Update(&UpdateObject{Actions: a})
		is.True(a.Left(c.more)) // not focused

		s.Focus(true)
		for i, step := range []struct {
			act     Action
			taken   bool
			current int
		}{
			{c.more, true, 10},
			{c.more, true, 20},
			{c.more, true, 20}, // at the end: still taken
			{c.less, true, 10},
			{c.other, false, 10},
		} {
			a := given(step.act)
			s.Update(&UpdateObject{Actions: a})
			is.Equal([]any{c.dir, i, !a.Left(step.act), s.Current}, []any{c.dir, i, step.taken, step.current}) // within its range before it draws
		}
	}

	off := newSlider(t, SliderOpts.DisableDefaultKeys(true))
	off.Focus(true)
	a := given(ActionRight)
	off.Update(&UpdateObject{Actions: a})
	is.True(a.Left(ActionRight))

	early := NewSlider() // focused before its first layout
	early.Focus(true)
	early.Update(&UpdateObject{Actions: a})
	is.True(a.Left(ActionRight))
}

// A focused list combo button opens with Activate. Open, its list takes Up,
// Down and Activate, which selects and closes, and the combo keeps the
// directions its list leaves, so the focus stays in it. Closed, it leaves
// the directions, though the Down and Up keys still open it.
func TestListComboButton_TakesActions(t *testing.T) {
	is := is.New(t)

	l := newListComboButton(t, ListComboButtonOpts.Entries([]any{"a", "b", "c"}), ListComboButtonOpts.EntryLabelFunc(label, label))
	update := func(acts ...Action) *Actions { // drawn without validating, which closes it
		a := given(acts...)
		l.Update(&UpdateObject{Actions: a})
		l.Render(ebiten.NewImage(1, 1))
		event.ExecuteDeferred()
		return a
	}
	a := update(ActionActivate)
	is.True(a.Left(ActionActivate) && !l.ContentVisible()) // not focused

	l.Focus(true)
	a = update(ActionDown)
	is.True(a.Left(ActionDown) && !l.ContentVisible()) // closed: left

	a = update(ActionActivate)
	is.True(!a.Left(ActionActivate) && l.ContentVisible()) // opens

	for i := range 3 { // to the last entry, and once more
		a = update(ActionDown)
		is.Equal([]any{i, a.Left(ActionDown)}, []any{i, false})
	}
	a = update(ActionLeft)
	is.True(!a.Left(ActionLeft)) // open: no direction leaves it

	a = update(ActionActivate)
	is.True(!a.Left(ActionActivate))
	is.Equal(l.SelectedEntry(), "c")
	is.True(!l.ContentVisible()) // selected and closed

	d := inputtest.Use(t)
	d.Keys[ebiten.KeyArrowDown] = true
	for _, open := range []bool{true, false} {
		input.Update()
		update()
		input.AfterUpdate()
		is.Equal(l.ContentVisible(), open) // the Down key opens it, once
		l.SetContentVisible(false)
	}

	off := newListComboButton(t, ListComboButtonOpts.Entries([]any{"a"}), ListComboButtonOpts.EntryLabelFunc(label, label),
		ListComboButtonOpts.DisableDefaultKeys(true))
	off.Focus(true)
	a = given(ActionActivate)
	off.Update(&UpdateObject{Actions: a})
	is.True(a.Left(ActionActivate) && !off.ContentVisible()) // default keys off
}
