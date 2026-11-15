package widget

import (
	"fmt"
	img "image"
	"math"
	"testing"

	"github.com/ebitenui/ebitenui/event"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/matryer/is"
)

// pointerAt moves the pointer to y over the middle of a 50 px wide view,
// with the button down or not, and runs an update and a draw. The draw
// doesn't validate again: a List rebuilds itself when validated.
func pointerAt(t *testing.T, r ValidatedRenderer, u interface{ Update(*UpdateObject) }, c *testCursor, y int, pressed bool) {
	t.Helper()
	c.just = pressed && !c.pressed
	c.x, c.y, c.pressed = 25, y, pressed
	u.Update(&UpdateObject{})
	event.ExecuteDeferred()
	r.Render(ebiten.NewImage(1, 1))
	event.ExecuteDeferred()
}

// newDragScroll returns a scroll container over ten rows of 20 px in a
// view of 50 px (150 px to scroll), with opts.
func newDragScroll(t *testing.T, rows []PreferredSizeLocateableWidget, opts ...ScrollContainerOpt) *ScrollContainer {
	t.Helper()
	content := NewContainer(ContainerOpts.Layout(NewRowLayout(RowLayoutOpts.Direction(DirectionVertical))))
	for range 10 - len(rows) {
		rows = append(rows, newSimpleWidget(50, 20, nil))
	}
	content.AddChild(rows...)
	ns := newNineSliceEmpty(t)
	s := NewScrollContainer(append(opts, ScrollContainerOpts.Content(content), ScrollContainerOpts.Image(&ScrollContainerImage{Idle: ns, Mask: ns}))...)
	s.SetLocation(img.Rect(0, 0, 50, 50))
	render(s, t)
	return s
}

func nearly(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

// A press in the view that moves past a few pixels drags the content with
// the pointer, reporting each move; a smaller one scrolls nothing.
func TestScrollContainer_DragToScroll(t *testing.T) {
	is := is.New(t)

	c := useTestCursor(t)
	s := newDragScroll(t, nil, ScrollContainerOpts.DragToScroll())
	var moves []float64
	s.ScrollMovedEvent.AddHandler(func(args any) {
		moves = append(moves, args.(*ScrollContainerScrollMovedEventArgs).Top) //nolint:forcetypeassert // the event's own args
	})

	pointerAt(t, s, s, c, 60, true) // below the view
	pointerAt(t, s, s, c, 20, true)
	pointerAt(t, s, s, c, 20, false)
	is.Equal(s.ScrollTop, 0.0)

	pointerAt(t, s, s, c, 40, true)
	pointerAt(t, s, s, c, 35, true) // 5 px: not yet a drag
	is.Equal(s.ScrollTop, 0.0)
	pointerAt(t, s, s, c, 30, true) // 10 px up: dragged 10
	is.True(nearly(s.ScrollTop, 10.0/150))
	pointerAt(t, s, s, c, 0, true)
	is.True(nearly(s.ScrollTop, 40.0/150))
	is.Equal(len(moves), 2) // each move reported
	is.True(nearly(moves[0], 10.0/150) && nearly(moves[1], 40.0/150))

	pointerAt(t, s, s, c, 0, true) // held still
	pointerAt(t, s, s, c, 0, false)
	is.True(nearly(s.ScrollTop, 40.0/150)) // let go standing still: no fling
}

// Let go while moving, the content keeps going, slower each update, and
// stops; pressed again, it stops at once.
func TestScrollContainer_DragKeepsMoving(t *testing.T) {
	is := is.New(t)

	c := useTestCursor(t)
	s := newDragScroll(t, nil, ScrollContainerOpts.DragToScroll())
	pointerAt(t, s, s, c, 49, true)
	for y := 39; y >= 19; y -= 10 { // 10 px an update
		pointerAt(t, s, s, c, y, true)
	}
	pointerAt(t, s, s, c, 19, false)
	before := s.ScrollTop

	last, step := s.ScrollTop, math.Inf(1)
	for range 10 {
		pointerAt(t, s, s, c, 19, false)
		d := s.ScrollTop - last
		is.True(d > 0 && d < step) // still moving, slower
		last, step = s.ScrollTop, d
	}
	is.True(s.ScrollTop > before)

	pointerAt(t, s, s, c, 19, true) // caught
	caught := s.ScrollTop
	pointerAt(t, s, s, c, 19, true)
	is.Equal(s.ScrollTop, caught)

	for range 300 { // and left alone, it comes to rest
		pointerAt(t, s, s, c, 19, false)
	}
	rest := s.ScrollTop
	pointerAt(t, s, s, c, 19, false)
	is.Equal(s.ScrollTop, rest)
}

// A press that becomes a drag clicks nothing; one that stays put clicks.
func TestScrollContainer_DraggedPressClicksNothing(t *testing.T) {
	is := is.New(t)

	c := useTestCursor(t)
	clicks := 0
	b := newButton(t, ButtonOpts.WidgetOpts(WidgetOpts.MinSize(50, 20)),
		ButtonOpts.ClickedHandler(func(*ButtonClickedEventArgs) { clicks++ }))
	s := newDragScroll(t, []PreferredSizeLocateableWidget{newSimpleWidget(50, 20, nil), b}, ScrollContainerOpts.DragToScroll())

	pointerAt(t, s, s, c, 30, true) // on the button
	pointerAt(t, s, s, c, 30, false)
	is.Equal(clicks, 1)

	pointerAt(t, s, s, c, 30, true)
	pointerAt(t, s, s, c, 20, true) // dragged 10: the button follows the pointer
	pointerAt(t, s, s, c, 20, false)
	is.Equal(clicks, 1)
	is.True(s.ScrollTop > 0)
}

// Without DragToScroll, a drag moves nothing.
func TestScrollContainer_NoDragByDefault(t *testing.T) {
	is := is.New(t)

	c := useTestCursor(t)
	s := newDragScroll(t, nil)
	pointerAt(t, s, s, c, 40, true)
	pointerAt(t, s, s, c, 0, true)
	is.Equal(s.ScrollTop, 0.0)
}

// A list that drags to scroll keeps its slider with the content.
func TestList_DragToScroll(t *testing.T) {
	is := is.New(t)

	c := useTestCursor(t)
	entries := make([]any, 30)
	for i := range entries {
		entries[i] = fmt.Sprint("entry ", i)
	}
	l := newList(t, ListOpts.Entries(entries), ListOpts.EntryLabelFunc(func(e any) string { return e.(string) }), //nolint:forcetypeassert // strings
		ListOpts.DragToScroll())
	l.SetLocation(img.Rect(0, 0, 200, 100))
	render(l, t)

	pointerAt(t, l, l, c, 80, true) // in the view, above the bottom slider
	pointerAt(t, l, l, c, 30, true)
	is.True(l.scrollContainer.ScrollTop > 0)
	is.Equal(l.vSlider.Current, int(math.Round(l.scrollContainer.ScrollTop*1000)))
}
