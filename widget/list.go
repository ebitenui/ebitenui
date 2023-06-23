package widget

import (
	img "image"
	"image/color"
	"log"
	"math"

	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/input"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

type List struct {
	EntrySelectedEvent *event.Event

	containerOpts            []ContainerOpt
	scrollContainerOpts      []ScrollContainerOpt
	sliderOpts               []SliderOpt
	entries                  []interface{}
	entryIdentifier          *func(interface{}) string
	idEntryTable             map[string]interface{}
	entriesButtonTable       map[interface{}]*Button
	entryLabelFunc           ListEntryLabelFunc
	entryFace                font.Face
	entryUnselectedColor     *ButtonImage
	entrySelectedColor       *ButtonImage
	entryUnselectedTextColor *ButtonTextColor
	entryTextColor           *ButtonTextColor
	entryTextPadding         Insets
	controlWidgetSpacing     int
	hideHorizontalSlider     bool
	hideVerticalSlider       bool
	allowReselect            bool

	init            *MultiOnce
	container       *Container
	listContent     *Container
	scrollContainer *ScrollContainer
	vSlider         *Slider
	hSlider         *Slider
	buttons         []*Button
	selectedEntry   interface{}

	disableDefaultKeys bool
	focused            bool
	tabOrder           int
	justMoved          bool
	focusIndex         int
	prevFocusIndex     int
}

type ListOpt func(l *List)

type ListEntryLabelFunc func(e interface{}) string

type ListEntryColor struct {
	Unselected                 color.Color
	Selected                   color.Color
	DisabledUnselected         color.Color
	DisabledSelected           color.Color
	SelectedBackground         color.Color
	FocusedBackground          color.Color
	SelectedFocusedBackground  color.Color
	DisabledSelectedBackground color.Color
}

type ListEntrySelectedEventArgs struct {
	List          *List
	Entry         interface{}
	PreviousEntry interface{}
}

type ListEntrySelectedHandlerFunc func(args *ListEntrySelectedEventArgs)

type ListOptions struct {
}

var ListOpts ListOptions

func NewList(opts ...ListOpt) *List {
	l := &List{
		EntrySelectedEvent: &event.Event{},

		init:               &MultiOnce{},
		entriesButtonTable: make(map[interface{}]*Button),
		idEntryTable:       make(map[string]interface{}),
		focusIndex:         0,
		prevFocusIndex:     -1,
	}

	l.init.Append(l.createWidget)

	for _, o := range opts {
		o(l)
	}

	l.resetFocusIndex()

	return l
}

func (o ListOptions) ContainerOpts(opts ...ContainerOpt) ListOpt {
	return func(l *List) {
		l.containerOpts = append(l.containerOpts, opts...)
	}
}

func (o ListOptions) ScrollContainerOpts(opts ...ScrollContainerOpt) ListOpt {
	return func(l *List) {
		l.scrollContainerOpts = append(l.scrollContainerOpts, opts...)
	}
}

func (o ListOptions) SliderOpts(opts ...SliderOpt) ListOpt {
	return func(l *List) {
		l.sliderOpts = append(l.sliderOpts, opts...)
	}
}

func (o ListOptions) ControlWidgetSpacing(s int) ListOpt {
	return func(l *List) {
		l.controlWidgetSpacing = s
	}
}

func (o ListOptions) HideHorizontalSlider() ListOpt {
	return func(l *List) {
		l.hideHorizontalSlider = true
	}
}

func (o ListOptions) HideVerticalSlider() ListOpt {
	return func(l *List) {
		l.hideVerticalSlider = true
	}
}

func (o ListOptions) Entries(e []interface{}) ListOpt {
	return func(l *List) {
		l.entries = e
	}
}

func (o ListOptions) Identifier(f func(interface{}) string) ListOpt {
	return func(l *List) {
		l.entryIdentifier = &f
	}
}

func (o ListOptions) EntryLabelFunc(f ListEntryLabelFunc) ListOpt {
	return func(l *List) {
		l.entryLabelFunc = f
	}
}

func (o ListOptions) EntryFontFace(f font.Face) ListOpt {
	return func(l *List) {
		l.entryFace = f
	}
}

func (o ListOptions) DisableDefaultKeys(val bool) ListOpt {
	return func(l *List) {
		l.disableDefaultKeys = val
	}
}

func (o ListOptions) EntryColor(c *ListEntryColor) ListOpt {
	return func(l *List) {
		l.entryUnselectedColor = &ButtonImage{
			Idle:     image.NewNineSliceColor(color.Transparent),
			Disabled: image.NewNineSliceColor(color.Transparent),
			Hover:    image.NewNineSliceColor(c.FocusedBackground),
		}

		l.entrySelectedColor = &ButtonImage{
			Idle:     image.NewNineSliceColor(c.SelectedBackground),
			Disabled: image.NewNineSliceColor(c.DisabledSelectedBackground),
			Hover:    image.NewNineSliceColor(c.SelectedFocusedBackground),
		}

		l.entryUnselectedTextColor = &ButtonTextColor{
			Idle:     c.Unselected,
			Disabled: c.DisabledUnselected,
		}

		l.entryTextColor = &ButtonTextColor{
			Idle:     c.Selected,
			Disabled: c.DisabledSelected,
		}
	}
}

func (o ListOptions) EntryTextPadding(i Insets) ListOpt {
	return func(l *List) {
		l.entryTextPadding = i
	}
}

func (o ListOptions) EntrySelectedHandler(f ListEntrySelectedHandlerFunc) ListOpt {
	return func(l *List) {
		l.EntrySelectedEvent.AddHandler(func(args interface{}) {
			f(args.(*ListEntrySelectedEventArgs))
		})
	}
}

func (o ListOptions) AllowReselect() ListOpt {
	return func(l *List) {
		l.allowReselect = true
	}
}

func (o ListOptions) TabOrder(tabOrder int) ListOpt {
	return func(l *List) {
		l.tabOrder = tabOrder
	}
}

func (l *List) GetWidget() *Widget {
	l.init.Do()
	return l.container.GetWidget()
}

func (l *List) PreferredSize() (int, int) {
	l.init.Do()
	return l.container.PreferredSize()
}

func (l *List) SetLocation(rect img.Rectangle) {
	l.init.Do()
	l.container.GetWidget().Rect = rect
}

func (l *List) RequestRelayout() {
	l.init.Do()
	l.container.RequestRelayout()
}

func (l *List) SetupInputLayer(def input.DeferredSetupInputLayerFunc) {
	l.init.Do()
	l.container.SetupInputLayer(def)
}

func (l *List) Render(screen *ebiten.Image, def DeferredRenderFunc) {
	l.init.Do()

	d := l.container.GetWidget().Disabled
	if l.vSlider != nil {
		l.vSlider.DrawTrackDisabled = d
	}
	if l.hSlider != nil {
		l.hSlider.DrawTrackDisabled = d
	}
	l.scrollContainer.GetWidget().Disabled = d

	l.handleInput()
	if l.focusIndex != l.prevFocusIndex && l.focusIndex <= len(l.buttons) {
		l.scrollVisible(l.buttons[l.focusIndex])
	}
	l.container.Render(screen, def)
}

func (l *List) Focus(focused bool) {
	l.init.Do()
	l.GetWidget().FireFocusEvent(l, focused, img.Point{-1, -1})
	l.focused = focused
}

func (l *List) IsFocused() bool {
	return l.focused
}

func (l *List) TabOrder() int {
	return l.tabOrder
}

func (l *List) handleInput() {
	if l.focused && !l.GetWidget().Disabled && len(l.buttons) > 0 {
		if input.KeyPressed(ebiten.KeyUp) || input.KeyPressed(ebiten.KeyDown) {
			if !l.justMoved {
				direction := -1
				if input.KeyPressed(ebiten.KeyDown) {
					direction = 1
				}
				l.buttons[l.focusIndex].focused = false
				l.prevFocusIndex = l.focusIndex
				l.focusIndex += direction
				if l.focusIndex < 0 {
					l.focusIndex = len(l.buttons) - 1
				}
				if l.focusIndex >= len(l.buttons) {
					l.focusIndex = 0
				}
				l.justMoved = true
			}
		} else {
			l.justMoved = false
		}

		l.buttons[l.focusIndex].focused = true
	} else if len(l.buttons) > 0 && l.focusIndex <= len(l.buttons) {
		l.buttons[l.focusIndex].focused = false
	}
}

func (l *List) resetFocusIndex() {
	if len(l.buttons) > 0 {
		if l.focusIndex != -1 {
			l.buttons[l.focusIndex].focused = false
		}
		for i := 0; i < len(l.entries); i++ {
			if l.entries[i] == l.selectedEntry {
				if i != l.focusIndex {
					l.prevFocusIndex = l.focusIndex
					l.focusIndex = i
				}
				return
			}
		}
		l.focusIndex = 0
	}
}

func (l *List) createWidget() {
	var cols int
	if l.hideVerticalSlider {
		cols = 1
	} else {
		cols = 2
	}

	l.container = NewContainer(
		append(l.containerOpts,
			ContainerOpts.Layout(NewGridLayout(
				GridLayoutOpts.Columns(cols),
				GridLayoutOpts.Stretch([]bool{true, false}, []bool{true, false}),
				GridLayoutOpts.Spacing(l.controlWidgetSpacing, l.controlWidgetSpacing))))...)

	l.listContent = NewContainer(
		ContainerOpts.Layout(NewRowLayout(
			RowLayoutOpts.Direction(DirectionVertical))),
	)

	l.buttons = make([]*Button, 0, len(l.entries))
	for _, e := range l.entries {
		e := e
		but := l.createEntry(e)

		l.buttons = append(l.buttons, but)
		l.entriesButtonTable[e] = but
		if l.entryIdentifier != nil {
			l.idEntryTable[(*l.entryIdentifier)(e)] = e
		}
		l.listContent.AddChild(but)
	}

	l.scrollContainer = NewScrollContainer(append(l.scrollContainerOpts, []ScrollContainerOpt{
		ScrollContainerOpts.Content(l.listContent),
		ScrollContainerOpts.StretchContentWidth(),
	}...)...)

	l.container.AddChild(l.scrollContainer)

	if !l.hideVerticalSlider {
		pageSizeFunc := func() int {
			return int(math.Round(float64(l.scrollContainer.ContentRect().Dy()) / float64(l.listContent.GetWidget().Rect.Dy()) * 1000))
		}

		l.vSlider = NewSlider(append(l.sliderOpts, []SliderOpt{
			SliderOpts.Direction(DirectionVertical),
			SliderOpts.MinMax(0, 1000),
			SliderOpts.PageSizeFunc(pageSizeFunc),
			SliderOpts.DisableDefaultKeys(l.disableDefaultKeys),
			SliderOpts.ChangedHandler(func(args *SliderChangedEventArgs) {
				l.scrollContainer.ScrollTop = float64(args.Slider.Current) / 1000
			}),
		}...)...)
		l.container.AddChild(l.vSlider)

		l.scrollContainer.widget.ScrolledEvent.AddHandler(func(args interface{}) {
			a := args.(*WidgetScrolledEventArgs)
			p := pageSizeFunc() / 3
			if p < 1 {
				p = 1
			}
			l.vSlider.Current -= int(math.Round(a.Y * float64(p)))
		})
	}

	if !l.hideHorizontalSlider {
		l.hSlider = NewSlider(append(l.sliderOpts, []SliderOpt{
			SliderOpts.Direction(DirectionHorizontal),
			SliderOpts.MinMax(0, 1000),
			SliderOpts.PageSizeFunc(func() int {
				return int(math.Round(float64(l.scrollContainer.ContentRect().Dx()) / float64(l.listContent.GetWidget().Rect.Dx()) * 1000))
			}),
			SliderOpts.ChangedHandler(func(args *SliderChangedEventArgs) {
				l.scrollContainer.ScrollLeft = float64(args.Slider.Current) / 1000
			}),
		}...)...)
		l.container.AddChild(l.hSlider)
	}

}

func (l *List) SetEntries(entries []interface{}) {
	l.entries = entries
	l.entriesButtonTable = make(map[interface{}]*Button, len(entries))
	l.idEntryTable = make(map[string]interface{}, len(entries))
	l.container.RemoveChildren()
	l.createWidget()
	l.resetFocusIndex()
}

func (l *List) RemoveEntry(entry interface{}) {
	if l.entryIdentifier == nil {
		log.Println("cannot add entry without identifier function")
		return
	}
	l.init.Do()

	identifier := (*l.entryIdentifier)(entry)
	targetEntry := l.idEntryTable[identifier]
	if len(l.entries) > 1 {
		for i, e := range l.entries {
			if e == targetEntry {
				l.entries = append(l.entries[:i], l.entries[i+1:]...)
				but := l.buttons[i]
				l.listContent.RemoveChild(but)
				l.focusIndex = i - 1
				delete(l.entriesButtonTable, e)
				delete(l.idEntryTable, identifier)
				l.buttons = append(l.buttons[:i], l.buttons[i+1:]...)
				break
			}
		}
	} else {
		l.entries = make([]interface{}, 0)
		l.buttons = make([]*Button, 0)
		l.focusIndex = -1
		l.listContent.RemoveChild(l.buttons[0])
		delete(l.entriesButtonTable, targetEntry)
		delete(l.idEntryTable, identifier)
	}

	l.resetFocusIndex()
}

func (l *List) AddEntry(entry interface{}) {
	if l.entryIdentifier == nil {
		log.Println("cannot add entry without identifier function")
		return
	}
	l.init.Do()

	l.entries = append(l.entries, entry)
	but := l.createEntry(entry)
	l.buttons = append(l.buttons, but)
	l.entriesButtonTable[entry] = but
	l.idEntryTable[(*l.entryIdentifier)(entry)] = entry
	l.listContent.AddChild(but)

	l.resetFocusIndex()
}

func (l *List) createEntry(entry interface{}) *Button {
	but := NewButton(
		ButtonOpts.WidgetOpts(WidgetOpts.LayoutData(RowLayoutData{
			Stretch: true,
		})),
		ButtonOpts.Image(l.entryUnselectedColor),
		ButtonOpts.TextSimpleLeft(l.entryLabelFunc(entry), l.entryFace, l.entryUnselectedTextColor, l.entryTextPadding),
		ButtonOpts.ClickedHandler(func(_ *ButtonClickedEventArgs) {
			l.setSelectedEntry(entry, true)
		}))

	return but
}

func (l *List) SetSelectedEntry(e interface{}) {
	l.setSelectedEntry(e, false)
}

func (l *List) setSelectedEntry(e interface{}, user bool) {
	if e != l.selectedEntry || (user && l.allowReselect) {
		l.init.Do()

		prev := l.selectedEntry
		l.selectedEntry = e
		l.resetFocusIndex()
		for i, b := range l.buttons {
			if l.entries[i] == e {
				b.Image = l.entrySelectedColor
				b.TextColor = l.entryTextColor
			} else {
				b.Image = l.entryUnselectedColor
				b.TextColor = l.entryUnselectedTextColor
			}
		}

		l.EntrySelectedEvent.Fire(&ListEntrySelectedEventArgs{
			Entry:         e,
			PreviousEntry: prev,
		})
	}
}

func (l *List) SelectedEntry() interface{} {
	l.init.Do()
	return l.selectedEntry
}

func (l *List) setScrollTop(t float64) {
	l.init.Do()
	if l.vSlider != nil {
		l.vSlider.Current = int(math.Round(t * 1000))
	}
	l.scrollContainer.ScrollTop = t
}

func (l *List) setScrollLeft(left float64) {
	l.init.Do()
	if l.hSlider != nil {
		l.hSlider.Current = int(math.Round(left * 1000))
	}
	l.scrollContainer.ScrollLeft = left
}

func scrollClamp(scroll float64) float64 {
	min, max := -0.1, 0.1
	if scroll < min {
		scroll = min
	} else if scroll > max {
		scroll = max
	}
	return scroll
}

func (l *List) scrollVisible(w HasWidget) {
	rect := l.scrollContainer.ContentRect()
	wrect := w.GetWidget().Rect
	if !wrect.In(rect) {
		ScrollTop := 0.0
		ScrollLeft := 0.0
		if wrect.Max.Y > rect.Max.Y {
			ScrollTop = float64(wrect.Max.Y - rect.Max.Y)
		} else if wrect.Min.Y < rect.Min.Y {
			ScrollTop = float64(wrect.Min.Y - rect.Min.Y)
		}
		if wrect.Max.X > rect.Max.X {
			ScrollLeft = float64(wrect.Max.X - rect.Max.X)
		} else if wrect.Min.X < rect.Min.X {
			ScrollLeft = -float64(wrect.Min.X - rect.Min.X)
		}
		l.setScrollTop(l.scrollContainer.ScrollTop + scrollClamp(ScrollTop/1000))
		l.setScrollLeft(l.scrollContainer.ScrollLeft + scrollClamp(ScrollLeft/1000))
	} else if wrect != rect {
		l.prevFocusIndex = l.focusIndex
	}
}
