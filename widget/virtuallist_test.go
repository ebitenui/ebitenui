package widget

import (
	img "image"
	"image/color"
	"math"
	"testing"

	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/input"
	"github.com/ebitenui/ebitenui/internal/input/inputtest"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/matryer/is"
)

// virtualRows counts the rows a VirtualList makes and the bindings it asks
// for, and remembers which entry each row shows.
type virtualRows struct {
	made, binds int
	shows       map[*simpleWidget]int
}

func (v *virtualRows) row() PreferredSizeLocateableWidget {
	v.made++
	return newSimpleWidget(10, 20, nil)
}

func (v *virtualRows) bind(row PreferredSizeLocateableWidget, index int) {
	v.binds++
	v.shows[row.(*simpleWidget)] = index //nolint:forcetypeassert // made by row
}

// at returns the entry the row laid out at y shows, or -1.
func (v *virtualRows) at(y int) int {
	for r, i := range v.shows {
		if r.GetWidget().Rect.Min.Y == y {
			return i
		}
	}
	return -1
}

// newVirtualList returns a list of count entries in rows of 20 px, 200 by
// 100 px, with opts, drawn once.
func newVirtualList(t *testing.T, count int, rows *virtualRows, opts ...VirtualListOpt) *VirtualList {
	t.Helper()
	ns := newNineSliceEmpty(t)
	l := NewVirtualList(append([]VirtualListOpt{
		VirtualListOpts.Count(count),
		VirtualListOpts.RowHeight(20),
		VirtualListOpts.RowFunc(rows.row),
		VirtualListOpts.BindFunc(rows.bind),
		VirtualListOpts.ScrollContainerImage(&ScrollContainerImage{Idle: ns, Mask: ns}),
		VirtualListOpts.SliderParams(&SliderParams{
			TrackImage:  &SliderTrackImage{},
			HandleImage: &ButtonImage{Idle: ns, Pressed: ns},
		}),
	}, opts...)...)
	l.SetLocation(img.Rect(0, 0, 200, 100))
	render(l, t)
	return l
}

// drawAgain draws l twice without validating it again, so what the slider
// reports as it draws shows too.
func drawAgain(l *VirtualList) {
	for range 2 {
		l.Render(ebiten.NewImage(1, 1))
		event.ExecuteDeferred()
	}
}

// Ten thousand entries make only the rows the view shows, plus one, and
// scrolling binds the rows to the entries now in view.
func TestVirtualList_RowsInViewOnly(t *testing.T) {
	is := is.New(t)

	rows := &virtualRows{shows: map[*simpleWidget]int{}}
	l := newVirtualList(t, 10_000, rows)
	view := l.scroll.ViewRect()
	is.True(rows.made <= view.Dy()/20+1)                 // a screenful
	is.Equal(l.content.GetWidget().Rect.Dy(), 10_000*20) // as tall as every row
	is.Equal(rows.at(view.Min.Y), 0)                     // the first entry at the top

	made := rows.made
	l.ScrollToIndex(5_000)
	drawAgain(l)
	is.Equal(rows.made, made)               // no new rows
	is.Equal(rows.at(view.Min.Y), 5_000)    // entry 5 000 at the top
	is.Equal(rows.at(view.Min.Y+20), 5_001) // and the next under it
	is.Equal(l.slider.Current, 5_000*20)    // the slider with it, in pixels

	binds := rows.binds
	l.ScrollToIndex(5_001)
	drawAgain(l)
	is.True(rows.binds-binds <= 1) // one row scrolled in, one bound

	l.ScrollToIndex(9_999)
	drawAgain(l)
	is.Equal(rows.at(view.Max.Y-20), 9_999) // the last entry at the bottom
}

// The slider counts pixels, so it reaches any entry of a long list, and the
// entry at the top stays there as entries are added.
func TestVirtualList_SliderInPixels(t *testing.T) {
	is := is.New(t)

	rows := &virtualRows{shows: map[*simpleWidget]int{}}
	l := newVirtualList(t, 100_000, rows)
	view := l.scroll.ViewRect()

	l.slider.Current = 12_345 * 20
	drawAgain(l)
	is.Equal(rows.at(view.Min.Y), 12_345) // the entry at that pixel

	l.SetCount(200_000)
	drawAgain(l)
	is.Equal(rows.at(view.Min.Y), 12_345) // still at the top

	l.ScrollToIndex(250_000)
	l.Render(ebiten.NewImage(1, 1))
	is.Equal(rows.at(view.Max.Y-20), 199_999) // the last entry at the bottom, in the first frame
}

// In a layout that gives it its preferred size, the list is a window on its
// rows: its minimum height, or one row without one, and only the rows that
// shows are made.
func TestVirtualList_PreferredSizeIsAWindow(t *testing.T) {
	is := is.New(t)

	ns := newNineSliceEmpty(t)
	for _, c := range []struct{ minHeight, height int }{{120, 120}, {0, 20}} {
		rows := &virtualRows{shows: map[*simpleWidget]int{}}
		l := NewVirtualList(
			VirtualListOpts.ContainerOpts(ContainerOpts.WidgetOpts(WidgetOpts.MinSize(100, c.minHeight))),
			VirtualListOpts.Count(10_000),
			VirtualListOpts.RowHeight(20),
			VirtualListOpts.RowFunc(rows.row),
			VirtualListOpts.BindFunc(rows.bind),
			VirtualListOpts.ScrollContainerImage(&ScrollContainerImage{Idle: ns, Mask: ns}),
			VirtualListOpts.SliderParams(&SliderParams{TrackImage: &SliderTrackImage{}, HandleImage: &ButtonImage{Idle: ns, Pressed: ns}}),
		)
		root := NewContainer(ContainerOpts.Layout(NewRowLayout(RowLayoutOpts.Direction(DirectionVertical))))
		root.AddChild(l)
		root.SetLocation(img.Rect(0, 0, 400, 400))
		render(root, t)

		_, h := l.PreferredSize()
		is.Equal(h, c.height)               // a window, not every row
		is.True(rows.made <= c.height/20+1) // a window's rows
	}
}

// The slider's handle is to its track what the view is to the list.
func TestVirtualList_HandleShowsTheView(t *testing.T) {
	is := is.New(t)

	rows := &virtualRows{shows: map[*simpleWidget]int{}}
	l := newVirtualList(t, 10, rows) // 200 px of rows in a 100 px view
	drawAgain(l)
	handle, track := l.slider.handle.GetWidget().Rect.Dy(), l.slider.GetWidget().Rect.Dy()
	is.True(math.Abs(float64(handle)/float64(track)-0.5) < 0.02) // half the track
}

// A count smaller than the view binds no row past the last entry, and
// Refresh binds the rows in view again.
func TestVirtualList_SetCountAndRefresh(t *testing.T) {
	is := is.New(t)

	rows := &virtualRows{shows: map[*simpleWidget]int{}}
	l := newVirtualList(t, 100, rows)
	l.ScrollToIndex(50)
	drawAgain(l)

	l.SetCount(3)
	clear(rows.shows)
	drawAgain(l)
	is.Equal(len(rows.shows), 3) // three rows bound
	for _, i := range rows.shows {
		is.True(i >= 0 && i < 3)
	}

	binds := rows.binds
	l.Refresh()
	drawAgain(l)
	is.Equal(rows.binds-binds, 3) // bound again
}

// Rows take input: a row's button clicked reports the entry it shows.
func TestVirtualList_RowsTakeInput(t *testing.T) {
	is := is.New(t)

	c := useTestCursor(t)

	var clicked []int
	bound := map[*Button]int{}
	ns := newNineSliceEmpty(t)
	l := NewVirtualList(
		VirtualListOpts.Count(1_000),
		VirtualListOpts.RowHeight(20),
		VirtualListOpts.RowFunc(func() PreferredSizeLocateableWidget {
			var b *Button
			b = NewButton(ButtonOpts.Image(&ButtonImage{Idle: ns, Pressed: ns}),
				ButtonOpts.ClickedHandler(func(*ButtonClickedEventArgs) { clicked = append(clicked, bound[b]) }))
			return b
		}),
		VirtualListOpts.BindFunc(func(row PreferredSizeLocateableWidget, i int) { bound[row.(*Button)] = i }), //nolint:forcetypeassert // made above
		VirtualListOpts.ScrollContainerImage(&ScrollContainerImage{Idle: ns, Mask: ns}),
		VirtualListOpts.SliderParams(&SliderParams{TrackImage: &SliderTrackImage{}, HandleImage: &ButtonImage{Idle: ns, Pressed: ns}}),
	)
	l.SetLocation(img.Rect(0, 0, 200, 100))
	render(l, t)
	l.ScrollToIndex(700)
	drawAgain(l)
	input.SetupInputLayersWithDeferred([]input.Layerer{l}) // as the UI does, dropping other tests' layers
	t.Cleanup(func() { input.SetupInputLayersWithDeferred(nil) })

	top := l.scroll.ViewRect().Min.Y
	for _, step := range []struct{ pressed, just bool }{{true, true}, {false, false}} {
		c.x, c.y, c.pressed, c.just = 20, top+30, step.pressed, step.just // the second row
		l.Update(&UpdateObject{})
		event.ExecuteDeferred()
	}
	is.Equal(clicked, []int{701})
}

// Drawn again in place, and scrolled within the rows it has, the list
// allocates nothing, the rows' backgrounds included.
func TestVirtualList_DrawsWithoutAllocating(t *testing.T) {
	is := is.New(t)

	rows := &virtualRows{shows: map[*simpleWidget]int{}}
	l := newVirtualList(t, 10_000, rows, VirtualListOpts.EntryImage(newEntryImages(t)))
	l.SetSelected(1)
	l.Focus(true)
	screen := ebiten.NewImage(1, 1)
	l.Render(screen)
	is.Equal(testing.AllocsPerRun(10, func() { l.content.Render(screen) }), 0.0)
}

// The wheel scrolls the list, and rows take the list's theme.
func TestVirtualList_WheelAndTheme(t *testing.T) {
	is := is.New(t)

	rows := &virtualRows{shows: map[*simpleWidget]int{}}
	l := newVirtualList(t, 1_000, rows)
	theme := &Theme{}
	l.GetWidget().SetTheme(theme)

	l.scroll.GetWidget().ScrolledEvent.Fire(&WidgetScrolledEventArgs{Widget: l.scroll.GetWidget(), Y: -1})
	event.ExecuteDeferred()
	drawAgain(l)
	is.Equal(rows.at(l.scroll.ViewRect().Min.Y), 3) // three rows down

	for r := range rows.shows {
		is.Equal(r.GetWidget().GetTheme(), theme)
	}
}

// With DragToScroll, dragging the rows scrolls the list and the slider
// follows, so the rows stay where the drag leaves them; without it, a drag
// scrolls nothing.
func TestVirtualList_DragToScroll(t *testing.T) {
	is := is.New(t)

	c := useTestCursor(t)
	for _, drag := range []bool{true, false} {
		rows := &virtualRows{shows: map[*simpleWidget]int{}}
		var opts []VirtualListOpt
		if drag {
			opts = append(opts, VirtualListOpts.DragToScroll())
		}
		l := newVirtualList(t, 1_000, rows, opts...)
		input.SetupInputLayersWithDeferred([]input.Layerer{l}) // as the UI does, dropping other tests' layers
		t.Cleanup(func() { input.SetupInputLayersWithDeferred(nil) })
		top := l.scroll.ViewRect().Min.Y

		pointerAt(t, l, l, c, top+80, true)
		pointerAt(t, l, l, c, top+20, true) // dragged up 60 px, three rows
		pointerAt(t, l, l, c, top+20, true) // held still
		pointerAt(t, l, l, c, top+20, false)
		drawAgain(l)
		if !drag {
			is.Equal(rows.at(top), 0) // not dragged
			continue
		}
		is.Equal(l.slider.Current, 60)
		is.Equal(rows.at(top), 3)
	}
}

// Without images of its own, the list takes its theme's list images, for
// its view and its slider, and follows the theme as it changes; given
// images win.
func TestVirtualList_FollowsTheTheme(t *testing.T) {
	is := is.New(t)

	theme := func() *Theme {
		ns := newNineSliceEmpty(t)
		return &Theme{ListTheme: &ListParams{
			ScrollContainerImage: &ScrollContainerImage{Idle: ns, Mask: ns},
			Slider:               &SliderParams{TrackImage: &SliderTrackImage{Idle: ns}, HandleImage: &ButtonImage{Idle: ns, Pressed: ns}},
		}}
	}
	dark, light := theme(), theme()
	rows := &virtualRows{shows: map[*simpleWidget]int{}}
	l := NewVirtualList(VirtualListOpts.Count(10), VirtualListOpts.RowHeight(20),
		VirtualListOpts.RowFunc(rows.row), VirtualListOpts.BindFunc(rows.bind))
	for _, th := range []*Theme{dark, light} {
		l.GetWidget().SetTheme(th)
		l.Validate()
		is.True(l.scroll.image == th.ListTheme.ScrollContainerImage)
		is.True(l.slider.definedParams.TrackImage == th.ListTheme.Slider.TrackImage)
	}

	given := newVirtualList(t, 10, rows)
	given.GetWidget().SetTheme(dark)
	given.Validate()
	is.True(given.scroll.image != dark.ListTheme.ScrollContainerImage)
	is.True(given.slider.definedParams.TrackImage != dark.ListTheme.Slider.TrackImage)
}

// clickRow presses and releases the pointer over the row y pixels into the
// view of l, updating l each time.
func clickRow(l *VirtualList, c *testCursor, y int) {
	for _, step := range []struct{ pressed, just bool }{{true, true}, {false, false}} {
		c.x, c.y, c.pressed, c.just = 20, l.scroll.ViewRect().Min.Y+y, step.pressed, step.just
		l.Update(&UpdateObject{})
		event.ExecuteDeferred()
	}
}

// A click on a row selects its entry, once: the list fires its event with
// the entry selected before and moves its focus there. The row under the
// pointer, and the one pressed while the button is down, are known as the
// rows draw.
func TestVirtualList_ClickSelects(t *testing.T) {
	is := is.New(t)

	c := useTestCursor(t)
	var got []VirtualListEntrySelectedEventArgs
	rows := &virtualRows{shows: map[*simpleWidget]int{}}
	l := newVirtualList(t, 1_000, rows,
		VirtualListOpts.EntrySelectedHandler(func(a *VirtualListEntrySelectedEventArgs) { got = append(got, *a) }))
	l.ScrollToIndex(700)
	drawAgain(l)
	input.SetupInputLayersWithDeferred([]input.Layerer{l}) // as the UI does, dropping other tests' layers
	t.Cleanup(func() { input.SetupInputLayersWithDeferred(nil) })

	c.x, c.y = 20, l.scroll.ViewRect().Min.Y+30
	l.Update(&UpdateObject{})
	event.ExecuteDeferred()
	is.Equal(l.hovered(), 701) // the second row in view
	top := l.content.GetWidget().Rect.Min.Y
	is.Equal([]int{l.entryAt(top - 1), l.entryAt(top + 1_000*20)}, []int{-1, -1}) // above and below the entries
	c.pressed, c.just = true, true
	l.Update(&UpdateObject{})
	event.ExecuteDeferred()
	is.Equal(l.pressed, 701)

	clickRow(l, c, 30)
	is.Equal(l.pressed, -1) // released
	is.Equal(l.Selected(), 701)
	is.Equal(len(got), 1)
	is.True(got[0].List == l)
	is.Equal([]int{got[0].Index, got[0].PreviousIndex}, []int{701, -1})
	is.Equal(l.focusIndex, 701) // the keys go on from there

	clickRow(l, c, 30)
	is.Equal(len(got), 1) // the selected entry: not again
	clickRow(l, c, 50)
	is.Equal([]int{got[1].Index, got[1].PreviousIndex}, []int{702, 701})

	l.GetWidget().Disabled = true
	clickRow(l, c, 10)
	is.Equal(l.Selected(), 702) // disabled
}

// Each row is drawn on its entry's background in the state it is in:
// selected or not, and disabled, pressed, or under the pointer or the keys'
// focus.
func TestVirtualList_RowBackgrounds(t *testing.T) {
	is := is.New(t)

	rows := &virtualRows{shows: map[*simpleWidget]int{}}
	images := newEntryImages(t)
	l := newVirtualList(t, 100, rows, VirtualListOpts.EntryImage(images))
	l.SetSelected(2)

	none := -1
	is.True(l.background(1, none) == images.Unselected)
	is.True(l.background(2, none) == images.Selected)
	is.True(l.background(1, 1) == images.Focused) // the pointer over it
	is.True(l.background(2, 2) == images.SelectedFocused)

	l.Focus(true)
	l.focusIndex = 1
	is.True(l.background(1, none) == images.Focused) // the keys' focus
	l.Focus(false)
	is.True(l.background(1, none) == images.Unselected)

	l.pressed = 1
	is.True(l.background(1, 1) == images.Selecting)
	l.pressed = 2
	is.True(l.background(2, 2) == images.SelectingFocused)

	l.GetWidget().Disabled = true
	is.True(l.background(1, 1) == images.DisabledUnselected)
	is.True(l.background(2, 2) == images.DisabledSelected)
}

// Without colors or images given, the rows take the theme's lists' ones,
// again as the theme changes; with colors given, not the theme's images.
func TestVirtualList_BackgroundsFromTheTheme(t *testing.T) {
	is := is.New(t)

	rows := &virtualRows{shows: map[*simpleWidget]int{}}
	themed := newEntryImages(t)
	l := newVirtualList(t, 10, rows)
	l.SetSelected(2)
	for _, th := range []*Theme{{ListTheme: &ListParams{EntryImage: themed}}, {ListTheme: &ListParams{EntryImage: newEntryImages(t)}}} {
		l.GetWidget().SetTheme(th)
		l.Validate()
		is.True(l.background(2, -1) == th.ListTheme.EntryImage.Selected)
	}

	selected := color.NRGBA{1, 2, 3, 255}
	l.GetWidget().SetTheme(&Theme{ListTheme: &ListParams{EntryColor: &ListEntryColor{SelectedBackground: selected}}})
	l.Validate()
	is.True(l.background(2, -1) == image.NewNineSliceColor(selected)) // the theme's colors

	given := newVirtualList(t, 10, rows, VirtualListOpts.EntryColor(&ListEntryColor{SelectedBackground: selected}))
	given.GetWidget().SetTheme(&Theme{ListTheme: &ListParams{EntryImage: themed}})
	given.Validate()
	given.SetSelected(2)
	is.True(given.background(2, -1) == image.NewNineSliceColor(selected))          // the color given
	is.True(given.background(1, -1) == image.NewNineSliceColor(color.Transparent)) // not the theme's image
}

// SetSelected selects an entry, or none with -1, and fires the event; an
// index past the entries changes nothing. SetCount drops a selection past
// the new count, and keeps the focus on an entry.
func TestVirtualList_SetSelected(t *testing.T) {
	is := is.New(t)

	var got []int
	rows := &virtualRows{shows: map[*simpleWidget]int{}}
	l := newVirtualList(t, 10, rows,
		VirtualListOpts.EntrySelectedHandler(func(a *VirtualListEntrySelectedEventArgs) { got = append(got, a.Index) }))
	is.Equal(l.Selected(), -1) // none at first

	l.SetSelected(3)
	l.SetSelected(3)
	l.SetSelected(10)
	l.SetSelected(-2)
	event.ExecuteDeferred()
	is.Equal(got, []int{3}) // once, and nothing past the entries
	l.SetSelected(-1)
	is.Equal(l.Selected(), -1)

	l.SetSelected(8)
	l.focusIndex = 9
	l.SetCount(5)
	is.Equal(l.Selected(), -1)
	is.Equal(l.focusIndex, 4)
}

// A focused virtual list moves through its entries with Up and Down and
// leaves those at its ends, so a gamepad moves the focus on from there;
// Activate, Enter or Space selects the entry it is on. Unfocused, or with
// its default keys off, it takes nothing.
func TestVirtualList_TakesActions(t *testing.T) {
	is := is.New(t)

	rows := &virtualRows{shows: map[*simpleWidget]int{}}
	l := newVirtualList(t, 3, rows)
	a := given(ActionDown)
	l.Update(&UpdateObject{Actions: a})
	is.True(a.Left(ActionDown)) // not focused

	l.Focus(true)
	is.True(l.IsFocused())
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
	is.True(!a.Left(ActionActivate))
	is.Equal(l.Selected(), 1)

	off := newVirtualList(t, 3, rows, VirtualListOpts.DisableDefaultKeys(true))
	off.Focus(true)
	a = given(ActionDown, ActionActivate)
	off.Update(&UpdateObject{Actions: a})
	is.True(a.Left(ActionDown) && a.Left(ActionActivate)) // default keys off

	d := inputtest.Use(t)
	for _, k := range []ebiten.Key{ebiten.KeyEnter, ebiten.KeySpace} {
		l.SetSelected(-1)
		d.Keys[k] = true
		input.Update()
		l.Update(&UpdateObject{Actions: &Actions{}})
		off.Update(&UpdateObject{Actions: &Actions{}})
		input.AfterUpdate()
		d.Keys[k] = false
		is.Equal([]int{l.Selected(), off.Selected()}, []int{1, -1}) // the key selects, but not with default keys off
	}
	d.Keys[ebiten.KeyEnter] = true
	for _, act := range []Action{ActionUp, ActionDown} {
		input.Update()
		l.Update(&UpdateObject{Actions: given(act)})
		input.AfterUpdate()
	}
	is.Equal([]int{l.focusIndex, l.Selected()}, []int{1, 0}) // held, it selected once, not as the focus moved on

	l.GetWidget().Disabled = true
	a = given(ActionUp)
	l.Update(&UpdateObject{Actions: a})
	is.True(a.Left(ActionUp)) // disabled
}

// The keys keep the focused entry in view: past the view's bottom the list
// scrolls down until its row's bottom shows, past the view's top up to it.
func TestVirtualList_FocusStaysInView(t *testing.T) {
	is := is.New(t)

	rows := &virtualRows{shows: map[*simpleWidget]int{}}
	l := newVirtualList(t, 1_000, rows)
	l.Focus(true)
	view := l.scroll.ViewRect().Dy()
	for range 6 {
		l.Update(&UpdateObject{Actions: given(ActionDown)})
	}
	is.Equal(l.slider.Current, 7*20-view) // the seventh row's bottom at the view's

	l.ScrollToIndex(500)
	l.Update(&UpdateObject{Actions: given(ActionUp)})
	is.Equal(l.slider.Current, 5*20) // up to the sixth row
}

// A virtual list is a focus stop with the tab order given, and keeps the
// focus links set on it.
func TestVirtualList_Focuser(t *testing.T) {
	is := is.New(t)

	rows := &virtualRows{shows: map[*simpleWidget]int{}}
	l := newVirtualList(t, 3, rows, VirtualListOpts.TabOrder(4))
	var f Focuser = l
	is.Equal(f.TabOrder(), 4)
	other := newVirtualList(t, 3, rows)
	l.AddFocus(FOCUS_SOUTH, other)
	is.True(l.GetFocus(FOCUS_SOUTH) == other)

	var events []bool
	l.GetWidget().FocusEvent.AddHandler(func(args any) {
		if a, ok := args.(*WidgetFocusEventArgs); ok {
			events = append(events, a.Focused)
		}
	})
	l.Focus(true)
	l.Focus(false)
	event.ExecuteDeferred()
	is.Equal(events, []bool{true, false}) // the UI follows it
}
