package widget

import (
	img "image"
	"math"
	"testing"

	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/input"
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
// allocates nothing.
func TestVirtualList_DrawsWithoutAllocating(t *testing.T) {
	is := is.New(t)

	rows := &virtualRows{shows: map[*simpleWidget]int{}}
	l := newVirtualList(t, 10_000, rows)
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
