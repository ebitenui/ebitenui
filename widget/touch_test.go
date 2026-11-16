package widget

import (
	img "image"
	"testing"

	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/input"
	"github.com/ebitenui/ebitenui/internal/input/inputtest"
	"github.com/matryer/is"
)

// touchButton returns a button at (20, 20)-(120, 60), away from where an
// absent mouse points, counting its clicks and context menus opened, on a
// touch screen or a desktop with one, and step, which runs an update with
// the fingers given down.
func touchButton(t *testing.T, touchscreen bool) (b *Button, clicks, menus *int, step func(fingers ...img.Point)) {
	t.Helper()
	d := inputtest.Use(t)
	d.Touchscreen = touchscreen
	input.SetupInputLayersWithDeferred(nil) // drops layers other tests left
	clicks, menus = new(int), new(int)
	b = newButton(t, ButtonOpts.ClickedHandler(func(*ButtonClickedEventArgs) { *clicks++ }),
		ButtonOpts.WidgetOpts(WidgetOpts.ContextMenu(NewContainer())))
	b.SetLocation(img.Rect(20, 20, 120, 60))
	b.GetWidget().ContextMenuEvent.AddHandler(func(any) { *menus++ })
	step = func(fingers ...img.Point) {
		d.Touches = fingers
		input.Update()
		b.Update(&UpdateObject{})
		input.AfterUpdate()
		event.ExecuteDeferred()
	}
	return b, clicks, menus, step
}

// A tap clicks, released where the finger lifts, on a desktop with a touch
// screen too. On a touch screen the pointer then leaves the screen, so the
// tapped widget doesn't stay hovered.
func TestTouch_TapLeavesNoHover(t *testing.T) {
	is := is.New(t)

	_, clicks, _, step := touchButton(t, false)
	step(img.Pt(70, 40))
	step()
	is.Equal(*clicks, 1) // on a desktop

	b, clicks, menus, step := touchButton(t, true)
	at := img.Pt(70, 40)
	step(at)
	step(at)
	step() // lifted: released where it was
	is.Equal(*clicks, 1)
	is.True(b.hovering) // the pointer over it, up
	step()              // the next update
	x, y := input.CursorPosition()
	is.Equal([]int{x, y}, []int{-1, -1})
	is.True(!b.hovering)
	is.Equal(*menus, 0)
}

// Held still for the long press's time, a touch opens the widget's context
// menu, as a right click does, clicks nothing when it lifts, and stays as a
// hover, so a tooltip shows, until the next touch ends.
func TestTouch_LongPress(t *testing.T) {
	is := is.New(t)

	b, clicks, menus, step := touchButton(t, true)
	at := img.Pt(70, 40)
	for range b.GetWidget().longPressDuration + 1 {
		step(at)
	}
	is.Equal(*menus, 1)
	step()
	step()
	is.Equal(*clicks, 0)
	is.True(b.hovering) // still hovered

	step(at) // the next touch, a tap on it
	step()
	step()
	is.True(!b.hovering)
	is.Equal(*clicks, 1)
}

// A second finger makes the touch a gesture: lifting the fingers clicks
// nothing, and holding them opens no context menu.
func TestTouch_SecondFingerCancels(t *testing.T) {
	is := is.New(t)

	b, clicks, menus, step := touchButton(t, true)
	at, other := img.Pt(70, 40), img.Pt(100, 50)
	step(at)
	for range b.GetWidget().longPressDuration + 1 {
		step(at, other)
	}
	step(at) // one lifts first
	step()
	is.Equal(*clicks, 0)
	is.Equal(*menus, 0)

	step(at) // a tap again is a tap
	step()
	is.Equal(*clicks, 1)
}

// The mouse's long press keeps its meaning: it opens no context menu, and
// the release still clicks.
func TestTouch_MouseLongPressUnchanged(t *testing.T) {
	is := is.New(t)

	c := useTestCursor(t)
	input.SetupInputLayersWithDeferred(nil)
	clicks, menus := 0, 0
	b := newButton(t, ButtonOpts.ClickedHandler(func(*ButtonClickedEventArgs) { clicks++ }),
		ButtonOpts.WidgetOpts(WidgetOpts.ContextMenu(NewContainer())))
	b.SetLocation(img.Rect(0, 0, 100, 40))
	b.GetWidget().ContextMenuEvent.AddHandler(func(any) { menus++ })
	c.x, c.y = 50, 20
	for i := range b.GetWidget().longPressDuration + 2 {
		c.pressed, c.just = true, i == 0
		b.Update(&UpdateObject{})
		event.ExecuteDeferred()
	}
	c.pressed, c.just = false, false
	b.Update(&UpdateObject{})
	event.ExecuteDeferred()
	is.Equal(menus, 0)
	is.Equal(clicks, 1)
}
