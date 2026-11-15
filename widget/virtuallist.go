package widget

import (
	img "image"
	"math"

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
type VirtualList struct {
	count         int
	rowHeight     int
	rowFunc       func() PreferredSizeLocateableWidget
	bindFunc      func(row PreferredSizeLocateableWidget, index int)
	sliderParams  *SliderParams
	scrollImage   *ScrollContainerImage
	containerOpts []ContainerOpt
	dragToScroll  bool

	init      *MultiOnce
	container *Container
	scroll    *ScrollContainer
	content   *virtualContent
	slider    *Slider
}

// VirtualListOpt configures a VirtualList.
type VirtualListOpt func(l *VirtualList)

type VirtualListOptions struct{}

// VirtualListOpts holds the options of NewVirtualList.
var VirtualListOpts VirtualListOptions

// NewVirtualList returns a VirtualList; RowHeight, RowFunc and BindFunc are
// required.
func NewVirtualList(opts ...VirtualListOpt) *VirtualList {
	l := &VirtualList{init: &MultiOnce{}}
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

// Count returns how many entries the list shows.
func (l *VirtualList) Count() int { return l.count }

// SetCount changes how many entries the list shows and binds the rows in
// view again.
func (l *VirtualList) SetCount(n int) {
	l.init.Do()
	l.count = n
	l.content.unbind()
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
	theme := l.container.GetWidget().GetTheme()
	if l.scrollImage == nil && theme != nil && theme.ListTheme != nil {
		l.scrollImage = theme.ListTheme.ScrollContainerImage
	}
	if l.sliderParams == nil && theme != nil && theme.ListTheme != nil {
		l.sliderParams = theme.ListTheme.Slider
	}
	l.scroll.image = l.scrollImage
	if p := l.sliderParams; p != nil {
		if l.slider.definedParams.TrackImage == nil {
			l.slider.definedParams.TrackImage = p.TrackImage
		}
		if l.slider.definedParams.HandleImage == nil {
			l.slider.definedParams.HandleImage = p.HandleImage
		}
	}
	l.container.Validate()
}

func (l *VirtualList) Render(screen *ebiten.Image) {
	l.init.Do()
	l.sync()
	l.container.Render(screen)
}

func (l *VirtualList) Update(updObj *UpdateObject) {
	l.init.Do()
	l.container.Update(updObj)
}

func (l *VirtualList) createWidget() {
	l.container = NewContainer(append([]ContainerOpt{
		ContainerOpts.Layout(NewGridLayout(
			GridLayoutOpts.Columns(2),
			GridLayoutOpts.Stretch([]bool{true, false}, []bool{true})))}, l.containerOpts...)...)
	l.content = &virtualContent{list: l, widget: NewWidget()}
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
	for i := c.first; i < c.last; i++ {
		if r, ok := c.rows[i%len(c.rows)].widget.(Renderer); ok {
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
