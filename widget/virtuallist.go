package widget

import (
	img "image"
	"math"

	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/input"
	"github.com/hajimehoshi/ebiten/v2"
)

// VirtualList shows a list of entries in rows of one height, keeping
// widgets only for the rows in view, plus one: as it scrolls, it binds the
// rows it has to the entries now in view. A list of thousands of entries so
// costs, in widgets, layout and drawing, what a screenful does. RowFunc
// makes a row's widget and BindFunc fills it with an entry, by index; a
// vertical slider, the mouse wheel, and with DragToScroll the pointer
// scroll the list.
//
// Entries are selected as in a List: a click on a row selects its entry,
// and focused, the list moves through its entries with Up and Down, keeping
// the one it is on in view, and selects it with Activate, Enter or Space,
// as a List's focused entry does. Each row is drawn
// on the background of its state, from EntryImage or EntryColor, or else
// from the theme's lists.
type VirtualList struct {
	EntrySelectedEvent *event.Event

	count              int
	rowHeight          int
	rowFunc            func() PreferredSizeLocateableWidget
	bindFunc           func(row PreferredSizeLocateableWidget, index int)
	sliderParams       *SliderParams
	scrollImage        *ScrollContainerImage
	entryColor         *ListEntryColor
	entryImage         *ListEntryImage
	containerOpts      []ContainerOpt
	dragToScroll       bool
	disableDefaultKeys bool
	tabOrder           int

	init      *MultiOnce
	container *Container
	scroll    *ScrollContainer
	content   *virtualContent
	slider    *Slider

	// The rows' backgrounds, unselected and selected, set as the list
	// validates.
	unselectedImage, selectedImage *ButtonImage

	selected   int // the selected entry, -1 for none
	focusIndex int // the entry the keys are on
	focused    bool
	focusMap   map[FocusDirection]Focuser
	pressed    int // the entry the pointer went down on, -1 for none
	hovering   bool
	hoverY     int // where the pointer is, while hovering
}

// VirtualListEntrySelectedEventArgs tell which entry was selected, and which
// was before, by index; -1 for none.
type VirtualListEntrySelectedEventArgs struct {
	List          *VirtualList
	Index         int
	PreviousIndex int
}

// VirtualListOpt configures a VirtualList.
type VirtualListOpt func(l *VirtualList)

type VirtualListOptions struct{}

// VirtualListOpts holds the options of NewVirtualList.
var VirtualListOpts VirtualListOptions

// NewVirtualList returns a VirtualList; RowHeight, RowFunc and BindFunc are
// required.
func NewVirtualList(opts ...VirtualListOpt) *VirtualList {
	l := &VirtualList{
		EntrySelectedEvent: &event.Event{},
		init:               &MultiOnce{},
		selected:           -1,
		pressed:            -1,
		focusMap:           map[FocusDirection]Focuser{},
	}
	l.init.Append(l.createWidget)
	for _, o := range opts {
		o(l)
	}
	return l
}

// Count sets how many entries the list shows.
func (o VirtualListOptions) Count(n int) VirtualListOpt {
	return func(l *VirtualList) { l.count = n }
}

// RowHeight sets every row's height, in pixels.
func (o VirtualListOptions) RowHeight(h int) VirtualListOpt {
	return func(l *VirtualList) { l.rowHeight = h }
}

// RowFunc sets what makes a row's widget; the list calls it as rows are
// needed, a screenful's worth at most.
func (o VirtualListOptions) RowFunc(f func() PreferredSizeLocateableWidget) VirtualListOpt {
	return func(l *VirtualList) { l.rowFunc = f }
}

// BindFunc sets what fills a row with the entry at index, called when a
// row comes to show another entry, or all rows in view after Refresh or
// SetCount.
func (o VirtualListOptions) BindFunc(f func(row PreferredSizeLocateableWidget, index int)) VirtualListOpt {
	return func(l *VirtualList) { l.bindFunc = f }
}

// SliderParams sets the vertical slider's images; the theme's list slider
// otherwise.
func (o VirtualListOptions) SliderParams(p *SliderParams) VirtualListOpt {
	return func(l *VirtualList) { l.sliderParams = p }
}

// ScrollContainerImage sets the scrolled view's images; the theme's list
// images otherwise.
func (o VirtualListOptions) ScrollContainerImage(i *ScrollContainerImage) VirtualListOpt {
	return func(l *VirtualList) { l.scrollImage = i }
}

// ContainerOpts configures the list's outer container.
func (o VirtualListOptions) ContainerOpts(opts ...ContainerOpt) VirtualListOpt {
	return func(l *VirtualList) { l.containerOpts = append(l.containerOpts, opts...) }
}

// DragToScroll lets the pointer scroll the rows by dragging them, as
// ScrollContainerOpts.DragToScroll does; the slider follows.
func (o VirtualListOptions) DragToScroll() VirtualListOpt {
	return func(l *VirtualList) { l.dragToScroll = true }
}

// EntrySelectedHandler adds f to the handlers of EntrySelectedEvent.
func (o VirtualListOptions) EntrySelectedHandler(f func(args *VirtualListEntrySelectedEventArgs)) VirtualListOpt {
	return func(l *VirtualList) {
		l.EntrySelectedEvent.AddHandler(func(args any) {
			if a, ok := args.(*VirtualListEntrySelectedEventArgs); ok {
				f(a)
			}
		})
	}
}

// EntryColor sets the colors of the rows' backgrounds, in place of the
// theme's lists' colors and images. A List's text colors don't apply: the
// rows draw their text.
func (o VirtualListOptions) EntryColor(c *ListEntryColor) VirtualListOpt {
	return func(l *VirtualList) { l.entryColor = c }
}

// EntryImage sets images for the rows' backgrounds, each in place of its
// color; the theme's lists' otherwise. See ListEntryImage.
func (o VirtualListOptions) EntryImage(i *ListEntryImage) VirtualListOpt {
	return func(l *VirtualList) { l.entryImage = i }
}

// DisableDefaultKeys keeps the focused list from taking Up, Down and
// Activate, and from selecting with Enter and Space.
func (o VirtualListOptions) DisableDefaultKeys(val bool) VirtualListOpt {
	return func(l *VirtualList) { l.disableDefaultKeys = val }
}

// TabOrder sets the list's place in the focus order.
func (o VirtualListOptions) TabOrder(tabOrder int) VirtualListOpt {
	return func(l *VirtualList) { l.tabOrder = tabOrder }
}

// Count returns how many entries the list shows.
func (l *VirtualList) Count() int { return l.count }

// SetCount changes how many entries the list shows and binds the rows in
// view again. A selection past the new count is dropped, without the event,
// and the focus moves to the last entry if it was past it.
func (l *VirtualList) SetCount(n int) {
	l.init.Do()
	l.count = n
	if l.selected >= n {
		l.selected = -1
	}
	l.focusIndex = max(0, min(l.focusIndex, n-1))
	l.content.unbind()
}

// Selected returns the selected entry's index, or -1 for none.
func (l *VirtualList) Selected() int { return l.selected }

// SetSelected selects the entry at index, or none with -1, and fires
// EntrySelectedEvent if that changes the selection; an index past the
// entries changes nothing. The keys go on from the entry selected.
func (l *VirtualList) SetSelected(index int) {
	if index < -1 || index >= l.count || index == l.selected {
		return
	}
	previous := l.selected
	l.selected = index
	if index >= 0 {
		l.focusIndex = index
	}
	l.EntrySelectedEvent.Fire(&VirtualListEntrySelectedEventArgs{List: l, Index: index, PreviousIndex: previous})
}

// Focus focuses the list, or not, and tells the UI.
func (l *VirtualList) Focus(focused bool) {
	l.init.Do()
	l.GetWidget().FireFocusEvent(l, focused, img.Point{-1, -1})
	l.focused = focused
}

func (l *VirtualList) IsFocused() bool { return l.focused }

func (l *VirtualList) TabOrder() int { return l.tabOrder }

func (l *VirtualList) GetFocus(direction FocusDirection) Focuser { return l.focusMap[direction] }

func (l *VirtualList) AddFocus(direction FocusDirection, focus Focuser) {
	l.focusMap[direction] = focus
}

// handleActions moves the focused list through its entries with Up and
// Down, leaving them at its ends so a gamepad moves the focus on from there,
// and selects the entry it is on with Activate, Enter or Space.
func (l *VirtualList) handleActions(actions *Actions) {
	if !l.focused || l.disableDefaultKeys || l.GetWidget().Disabled || l.count == 0 {
		return
	}
	if l.focusIndex < l.count-1 && actions.Take(ActionDown) {
		l.focusIndex++
		l.scrollToFocus()
	}
	if l.focusIndex > 0 && actions.Take(ActionUp) {
		l.focusIndex--
		l.scrollToFocus()
	}
	if actions.Take(ActionActivate) || input.KeyJustPressed(ebiten.KeyEnter) || input.KeyJustPressed(ebiten.KeySpace) {
		l.SetSelected(l.focusIndex)
	}
}

// scrollToFocus scrolls as little as shows the focused entry's row whole.
func (l *VirtualList) scrollToFocus() {
	top, view := l.focusIndex*l.rowHeight, l.scroll.ViewRect().Dy()
	switch {
	case top < l.slider.Current:
		l.slider.Current = top
	case top+l.rowHeight > l.slider.Current+view:
		l.slider.Current = top + l.rowHeight - view
	}
}

// entryAt returns the entry whose row holds the content's y, -1 for none.
func (l *VirtualList) entryAt(y int) int {
	i := (y - l.content.widget.Rect.Min.Y) / l.rowHeight
	if y < l.content.widget.Rect.Min.Y || i >= l.count {
		return -1
	}
	return i
}

// hovered returns the entry under the pointer, -1 for none.
func (l *VirtualList) hovered() int {
	if !l.hovering {
		return -1
	}
	return l.entryAt(l.hoverY)
}

// background returns what the row of entry i is drawn on, hover being the
// entry under the pointer: the selected or unselected background in the
// row's state.
func (l *VirtualList) background(i int, hover int) *image.NineSlice {
	set := l.unselectedImage
	if i == l.selected {
		set = l.selectedImage
	}
	switch {
	case l.container.GetWidget().Disabled:
		return set.Disabled
	case i == l.pressed:
		return set.Pressed
	case i == hover || l.focused && i == l.focusIndex:
		return set.Hover
	}
	return set.Idle
}

// Refresh binds the rows in view again, after the entries they show changed.
func (l *VirtualList) Refresh() {
	l.init.Do()
	l.content.unbind()
}

// ScrollToIndex scrolls so the entry at index shows at the top of the view,
// or as near as the list's end allows.
func (l *VirtualList) ScrollToIndex(index int) {
	l.init.Do()
	l.slider.Current = index * l.rowHeight
}

// sync scrolls the view to the slider. The slider counts pixels from the
// top, so it reaches every row of a long list, and keeps the entry at the
// top in place as the count or the view changes. Both clamp as they draw.
func (l *VirtualList) sync() {
	span := max(0, l.count*l.rowHeight-l.scroll.ViewRect().Dy())
	l.slider.Max = span
	l.scroll.ScrollTop = 0
	if span > 0 {
		l.scroll.ScrollTop = float64(l.slider.Current) / float64(span)
	}
}

func (l *VirtualList) GetWidget() *Widget {
	l.init.Do()
	return l.container.GetWidget()
}

// PreferredSize is the list's minimum height, or one row's without one:
// the list is a window on its rows, so it doesn't ask for all their height.
func (l *VirtualList) PreferredSize() (int, int) {
	l.init.Do()
	w, _ := l.container.PreferredSize()
	return w, max(l.rowHeight, l.container.GetWidget().MinHeight)
}

func (l *VirtualList) SetLocation(rect img.Rectangle) {
	l.init.Do()
	l.container.SetLocation(rect)
}

func (l *VirtualList) RequestRelayout() {
	l.init.Do()
	l.container.RequestRelayout()
}

func (l *VirtualList) SetupInputLayer(def input.DeferredSetupInputLayerFunc) {
	l.init.Do()
	l.container.SetupInputLayer(def)
}

func (l *VirtualList) Validate() {
	l.init.Do()
	if l.rowHeight <= 0 || l.rowFunc == nil || l.bindFunc == nil {
		panic("VirtualList: RowHeight, RowFunc and BindFunc are required.")
	}
	// The images given, or else the theme's lists', taken again on each
	// validation so that a theme switch shows: the scroll container falls
	// back to the theme's list images itself.
	l.scroll.definedImage = l.scrollImage
	slider := l.sliderParams
	if theme := l.container.GetWidget().GetTheme(); slider == nil && theme != nil && theme.ListTheme != nil {
		slider = theme.ListTheme.Slider
	}
	if slider != nil {
		l.slider.definedParams.TrackImage, l.slider.definedParams.HandleImage = slider.TrackImage, slider.HandleImage
	}
	colors, images := l.entryColor, l.entryImage
	if theme := l.container.GetWidget().GetTheme(); theme != nil && theme.ListTheme != nil {
		if colors == nil {
			colors = theme.ListTheme.EntryColor
		}
		if images == nil && l.entryColor == nil {
			images = theme.ListTheme.EntryImage // colors given show, not the theme's images
		}
	}
	l.unselectedImage, l.selectedImage = entryBackgrounds(colors, images)
	l.container.Validate()
}

func (l *VirtualList) Render(screen *ebiten.Image) {
	l.init.Do()
	l.sync()
	l.container.Render(screen)
}

func (l *VirtualList) Update(updObj *UpdateObject) {
	l.init.Do()
	l.handleActions(updObj.Actions)
	l.container.Update(updObj)
}

func (l *VirtualList) createWidget() {
	l.container = NewContainer(append([]ContainerOpt{
		ContainerOpts.Layout(NewGridLayout(
			GridLayoutOpts.Columns(2),
			GridLayoutOpts.Stretch([]bool{true, false}, []bool{true})))}, l.containerOpts...)...)
	l.content = &virtualContent{list: l, widget: NewWidget()}
	l.watchPointer()
	scrollOpts := []ScrollContainerOpt{ScrollContainerOpts.Content(l.content), ScrollContainerOpts.StretchContentWidth()}
	if l.dragToScroll {
		scrollOpts = append(scrollOpts, ScrollContainerOpts.DragToScroll())
	}
	l.scroll = NewScrollContainer(scrollOpts...)
	// A page is the view's height in slider steps, scaled so the handle is
	// to the track what the view is to the list.
	pageSize := func() int {
		total := l.count * l.rowHeight
		if total <= 0 {
			return 1
		}
		view := float64(l.scroll.ViewRect().Dy())
		return max(1, int(math.Round(view*float64(l.slider.Max+1)/float64(total))))
	}
	l.slider = NewSlider(
		SliderOpts.Direction(DirectionVertical),
		SliderOpts.MinMax(0, 0), // the max set as the list draws
		SliderOpts.InitialCurrent(0),
		SliderOpts.FixedHandleSize(0),
		SliderOpts.PageSizeFunc(pageSize),
	)
	// The wheel moves three rows a notch.
	l.scroll.GetWidget().ScrolledEvent.AddHandler(func(args any) {
		if a, ok := args.(*WidgetScrolledEventArgs); ok {
			l.slider.Current -= int(math.Round(a.Y * float64(3*l.rowHeight)))
		}
	})
	// A drag moves the view; the slider, which the view follows as the list
	// draws, moves with it.
	l.scroll.ScrollMovedEvent.AddHandler(func(args any) {
		if a, ok := args.(*ScrollContainerScrollMovedEventArgs); ok {
			l.slider.Current = int(math.Round(a.Top * float64(l.slider.Max)))
		}
	})
	l.container.AddChild(l.scroll, l.slider)
}

// watchPointer follows the pointer over the rows: the entry it is over, the
// one it goes down on, and a click, which selects the entry clicked unless
// the list is disabled.
func (l *VirtualList) watchPointer() {
	w := l.content.widget
	w.CursorEnterEvent.AddHandler(func(args any) {
		if a, ok := args.(*WidgetCursorEnterEventArgs); ok {
			l.hovering, l.hoverY = true, w.Rect.Min.Y+a.OffsetY
		}
	})
	w.CursorMoveEvent.AddHandler(func(args any) {
		if a, ok := args.(*WidgetCursorMoveEventArgs); ok {
			l.hoverY = w.Rect.Min.Y + a.OffsetY
		}
	})
	w.CursorExitEvent.AddHandler(func(any) { l.hovering = false })
	w.MouseButtonPressedEvent.AddHandler(func(args any) {
		if a, ok := args.(*WidgetMouseButtonPressedEventArgs); ok && a.Button == ebiten.MouseButtonLeft {
			l.pressed = l.entryAt(w.Rect.Min.Y + a.OffsetY)
		}
	})
	w.MouseButtonReleasedEvent.AddHandler(func(args any) {
		if a, ok := args.(*WidgetMouseButtonReleasedEventArgs); ok && a.Button == ebiten.MouseButtonLeft {
			l.pressed = -1
		}
	})
	w.MouseButtonClickedEvent.AddHandler(func(args any) {
		if a, ok := args.(*WidgetMouseButtonClickedEventArgs); ok && a.Button == ebiten.MouseButtonLeft && !l.GetWidget().Disabled {
			if i := l.entryAt(w.Rect.Min.Y + a.OffsetY); i >= 0 {
				l.SetSelected(i)
			}
		}
	})
}

// virtualContent is a VirtualList's scrolled content: as tall as every row,
// holding widgets only for the rows in view. Entry i shows in row i modulo
// the number of rows, so scrolling by a row binds one row.
type virtualContent struct {
	list        *VirtualList
	widget      *Widget
	rows        []virtualRow
	first, last int // the entries in view, last excluded
}

type virtualRow struct {
	widget PreferredSizeLocateableWidget
	index  int // the entry it shows, -1 for none
}

func (c *virtualContent) GetWidget() *Widget { return c.widget }

func (c *virtualContent) PreferredSize() (int, int) {
	return 0, c.list.count * c.list.rowHeight
}

func (c *virtualContent) SetLocation(rect img.Rectangle) { c.widget.Rect = rect }

func (c *virtualContent) Validate() {}

// unbind marks every row as showing nothing, to be bound again.
func (c *virtualContent) unbind() {
	for i := range c.rows {
		c.rows[i].index = -1
	}
}

// layOut finds the entries in view, makes rows when the view needs more,
// binds the rows that come to show another entry, and places them.
func (c *virtualContent) layOut() {
	l, rect, view := c.list, c.widget.Rect, c.list.scroll.ViewRect()
	h := l.rowHeight
	c.first = max(0, (view.Min.Y-rect.Min.Y)/h)
	c.last = min(l.count, (view.Max.Y-rect.Min.Y+h-1)/h)
	if c.last <= c.first {
		c.last = c.first
		return
	}
	for len(c.rows) < c.last-c.first {
		w := l.rowFunc()
		w.GetWidget().parent = c.widget
		w.Validate()
		c.rows = append(c.rows, virtualRow{widget: w, index: -1})
	}
	for i := c.first; i < c.last; i++ {
		r := &c.rows[i%len(c.rows)]
		if r.index != i {
			l.bindFunc(r.widget, i)
			r.index = i
		}
		r.widget.SetLocation(img.Rect(rect.Min.X, rect.Min.Y+i*h, rect.Max.X, rect.Min.Y+(i+1)*h))
	}
}

func (c *virtualContent) Render(screen *ebiten.Image) {
	c.widget.Render(screen)
	c.layOut()
	hover := c.list.hovered()
	for i := c.first; i < c.last; i++ {
		row := c.rows[i%len(c.rows)].widget
		rect := row.GetWidget().Rect
		c.list.background(i, hover).Draw(screen, rect.Dx(), rect.Dy(), row.GetWidget().drawImageOptions)
		if r, ok := row.(Renderer); ok {
			r.Render(screen)
		}
	}
}

func (c *virtualContent) Update(updObj *UpdateObject) {
	c.widget.Update(updObj)
	for i := c.first; i < c.last; i++ {
		if u, ok := c.rows[i%len(c.rows)].widget.(Updater); ok {
			u.Update(updObj)
		}
	}
}

func (c *virtualContent) SetupInputLayer(def input.DeferredSetupInputLayerFunc) {
	for _, r := range c.rows {
		if l, ok := r.widget.(input.Layerer); ok {
			l.SetupInputLayer(def)
		}
	}
}
