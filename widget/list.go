package widget

import (
	img "image"
	"image/color"
	"math"

	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/input"
	"github.com/ebitenui/ebitenui/utilities/constantutil"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/exp/slices"
)

type ListParams struct {
	EntryFace                   *text.Face
	EntryColor                  *ListEntryColor
	EntryTextPadding            *Insets
	EntryTextHorizontalPosition *TextPosition
	EntryTextVerticalPosition   *TextPosition
	ControlWidgetSpacing        *int
	MinSize                     *img.Point

	AllowReselect      *bool
	SelectFocus        *bool
	SelectPressed      *bool
	DisableDefaultKeys *bool

	Slider                 *SliderParams
	ScrollContainerImage   *ScrollContainerImage
	ScrollContainerPadding *Insets

	entryUnselectedColor     *ButtonImage
	entrySelectedColor       *ButtonImage
	entryUnselectedTextColor *ButtonTextColor
	entryTextColor           *ButtonTextColor
}

type List struct {
	definedParams  ListParams
	computedParams ListParams

	EntrySelectedEvent *event.Event

	containerOpts        []ContainerOpt
	hideHorizontalSlider bool
	hideVerticalSlider   bool

	entries        []any
	entryLabelFunc ListEntryLabelFunc

	init            *MultiOnce
	container       *Container
	listContent     *Container
	scrollContainer *ScrollContainer
	layout          *GridLayout
	vSlider         *Slider
	hSlider         *Slider
	buttons         []*Button
	selectedEntry   any

	focused        bool
	tabOrder       int
	justMoved      bool
	focusIndex     int
	prevFocusIndex int

	focusMap map[FocusDirection]Focuser
}

type ListOpt func(l *List)

type ListEntryLabelFunc func(e any) string

type ListEntryColor struct {
	Unselected                 color.Color
	Selected                   color.Color
	DisabledUnselected         color.Color
	DisabledSelected           color.Color
	SelectingBackground        color.Color
	SelectedBackground         color.Color
	FocusedBackground          color.Color
	SelectingFocusedBackground color.Color
	SelectedFocusedBackground  color.Color
	DisabledSelectedBackground color.Color
}

type ListEntrySelectedEventArgs struct {
	List          *List
	Entry         any
	PreviousEntry any
}

type ListEntrySelectedHandlerFunc func(args *ListEntrySelectedEventArgs)

type ListOptions struct {
}

var ListOpts ListOptions

func NewList(opts ...ListOpt) *List {
	l := &List{
		EntrySelectedEvent: &event.Event{},

		init:           &MultiOnce{},
		focusIndex:     0,
		prevFocusIndex: -1,
		focusMap:       make(map[FocusDirection]Focuser),
	}

	l.init.Append(l.createWidget)

	for _, o := range opts {
		o(l)
	}

	return l
}

func (l *List) Validate() {
	l.init.Do()
	l.populateComputedParams()

	if l.computedParams.EntryFace == nil {
		panic("List: EntryFontFace is required.")
	}
	if l.entryLabelFunc == nil {
		panic("List: EntryLabelFunc is required.")
	}
	l.initWidget()
	l.resetFocusIndex()
}

func (t *List) populateComputedParams() {
	params := ListParams{}

	theme := t.GetWidget().GetTheme()

	// Set theme values
	if theme != nil {
		if theme.ListTheme != nil {
			params.AllowReselect = theme.ListTheme.AllowReselect
			params.ControlWidgetSpacing = theme.ListTheme.ControlWidgetSpacing
			params.DisableDefaultKeys = theme.ListTheme.DisableDefaultKeys
			if theme.ListTheme.EntryFace != nil {
				params.EntryFace = theme.ListTheme.EntryFace
			} else {
				params.EntryFace = theme.DefaultFace
			}
			params.EntryColor = theme.ListTheme.EntryColor
			params.EntryTextHorizontalPosition = theme.ListTheme.EntryTextHorizontalPosition
			params.EntryTextPadding = theme.ListTheme.EntryTextPadding
			params.EntryTextVerticalPosition = theme.ListTheme.EntryTextVerticalPosition
			params.ScrollContainerImage = theme.ListTheme.ScrollContainerImage
			params.ScrollContainerPadding = theme.ListTheme.ScrollContainerPadding
			params.SelectFocus = theme.ListTheme.SelectFocus
			params.SelectPressed = theme.ListTheme.SelectPressed
			params.Slider = theme.ListTheme.Slider
			params.MinSize = theme.ListTheme.MinSize
		}
	}

	// Set definedParam values
	if t.definedParams.AllowReselect != nil {
		params.AllowReselect = t.definedParams.AllowReselect
	}
	if t.definedParams.ControlWidgetSpacing != nil {
		params.ControlWidgetSpacing = t.definedParams.ControlWidgetSpacing
	}
	if t.definedParams.DisableDefaultKeys != nil {
		params.DisableDefaultKeys = t.definedParams.DisableDefaultKeys
	}
	if t.definedParams.EntryFace != nil {
		params.EntryFace = t.definedParams.EntryFace
	}
	if t.definedParams.EntryColor != nil {
		params.EntryColor = t.definedParams.EntryColor
	}
	if t.definedParams.EntryTextHorizontalPosition != nil {
		params.EntryTextHorizontalPosition = t.definedParams.EntryTextHorizontalPosition
	}
	if t.definedParams.EntryTextPadding != nil {
		params.EntryTextPadding = t.definedParams.EntryTextPadding
	}
	if t.definedParams.EntryTextVerticalPosition != nil {
		params.EntryTextVerticalPosition = t.definedParams.EntryTextVerticalPosition
	}
	if t.definedParams.ScrollContainerImage != nil {
		params.ScrollContainerImage = t.definedParams.ScrollContainerImage
	}
	if t.definedParams.ScrollContainerPadding != nil {
		params.ScrollContainerPadding = t.definedParams.ScrollContainerPadding
	}
	if t.definedParams.SelectFocus != nil {
		params.SelectFocus = t.definedParams.SelectFocus
	}
	if t.definedParams.SelectPressed != nil {
		params.SelectPressed = t.definedParams.SelectPressed
	}
	if t.definedParams.MinSize != nil {
		params.MinSize = t.definedParams.MinSize
	}

	if t.definedParams.Slider != nil {
		if params.Slider == nil {
			params.Slider = &SliderParams{}
		}
		if t.definedParams.Slider.FixedHandleSize != nil {
			params.Slider.FixedHandleSize = t.definedParams.Slider.FixedHandleSize
		}
		if t.definedParams.Slider.HandleImage != nil {
			params.Slider.HandleImage = t.definedParams.Slider.HandleImage
		}
		if t.definedParams.Slider.MinHandleSize != nil {
			params.Slider.MinHandleSize = t.definedParams.Slider.MinHandleSize
		}
		if t.definedParams.Slider.TrackImage != nil {
			params.Slider.TrackImage = t.definedParams.Slider.TrackImage
		}
		if t.definedParams.Slider.TrackOffset != nil {
			params.Slider.TrackOffset = t.definedParams.Slider.TrackOffset
		}
		if t.definedParams.Slider.TrackPadding != nil {
			params.Slider.TrackPadding = t.definedParams.Slider.TrackPadding
		}
	}

	// Set defaults
	if params.EntryTextPadding == nil {
		params.EntryTextPadding = &Insets{}
	}
	if params.EntryTextHorizontalPosition == nil {
		pos := TextPositionCenter
		params.EntryTextHorizontalPosition = &pos
	}
	if params.EntryTextVerticalPosition == nil {
		pos := TextPositionCenter
		params.EntryTextVerticalPosition = &pos
	}
	if params.ControlWidgetSpacing == nil {
		spacing := 0
		params.ControlWidgetSpacing = &spacing
	}
	FALSE := false
	if params.AllowReselect == nil {
		params.AllowReselect = &FALSE
	}
	if params.SelectFocus == nil {
		params.SelectFocus = &FALSE
	}
	if params.SelectPressed == nil {
		params.SelectPressed = &FALSE
	}
	if params.DisableDefaultKeys == nil {
		params.DisableDefaultKeys = &FALSE
	}
	if params.ScrollContainerPadding == nil {
		params.ScrollContainerPadding = &Insets{}
	}

	params.entryUnselectedColor = &ButtonImage{
		Idle:     image.NewNineSliceColor(color.Transparent),
		Disabled: image.NewNineSliceColor(color.Transparent),
		Hover:    image.NewNineSliceColor(params.EntryColor.FocusedBackground),
		Pressed:  image.NewNineSliceColor(params.EntryColor.SelectingBackground),
	}

	params.entrySelectedColor = &ButtonImage{
		Idle:     image.NewNineSliceColor(params.EntryColor.SelectedBackground),
		Disabled: image.NewNineSliceColor(params.EntryColor.DisabledSelectedBackground),
		Hover:    image.NewNineSliceColor(params.EntryColor.SelectedFocusedBackground),
		Pressed:  image.NewNineSliceColor(params.EntryColor.SelectingFocusedBackground),
	}

	params.entryUnselectedTextColor = &ButtonTextColor{
		Idle:     params.EntryColor.Unselected,
		Disabled: params.EntryColor.DisabledUnselected,
	}

	params.entryTextColor = &ButtonTextColor{
		Idle:     params.EntryColor.Selected,
		Disabled: params.EntryColor.DisabledSelected,
	}

	t.computedParams = params
}

func (o ListOptions) ContainerOpts(opts ...ContainerOpt) ListOpt {
	return func(l *List) {
		l.containerOpts = append(l.containerOpts, opts...)
	}
}

// Specify the images for the scroll container.
func (o ListOptions) ScrollContainerImage(image *ScrollContainerImage) ListOpt {
	return func(l *List) {
		l.definedParams.ScrollContainerImage = image
	}
}

// Specify the padding for the scroll container.
func (o ListOptions) ScrollContainerPadding(padding *Insets) ListOpt {
	return func(l *List) {
		l.definedParams.ScrollContainerPadding = padding
	}
}

// Specify the options for the scroll bars.
func (o ListOptions) SliderParams(sliderParams *SliderParams) ListOpt {
	return func(l *List) {
		l.definedParams.Slider = sliderParams
	}
}

func (o ListOptions) ControlWidgetSpacing(s int) ListOpt {
	return func(l *List) {
		l.definedParams.ControlWidgetSpacing = &s
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

func (o ListOptions) Entries(e []any) ListOpt {
	return func(l *List) {
		l.entries = slices.CompactFunc(e, func(a any, b any) bool { return a == b })
	}
}

func (o ListOptions) EntryLabelFunc(f ListEntryLabelFunc) ListOpt {
	return func(l *List) {
		l.entryLabelFunc = f
	}
}

func (o ListOptions) EntryFontFace(f *text.Face) ListOpt {
	return func(l *List) {
		l.definedParams.EntryFace = f
	}
}

func (o ListOptions) DisableDefaultKeys(val bool) ListOpt {
	return func(l *List) {
		l.definedParams.DisableDefaultKeys = &val
	}
}

func (o ListOptions) EntryColor(c *ListEntryColor) ListOpt {
	return func(l *List) {
		l.definedParams.EntryColor = c
	}
}

func (o ListOptions) EntryTextPadding(i *Insets) ListOpt {
	return func(l *List) {
		l.definedParams.EntryTextPadding = i
	}
}

// EntryTextPosition sets the position of the text for entries.
// Defaults to both TextPositionCenter.
func (o ListOptions) EntryTextPosition(h TextPosition, v TextPosition) ListOpt {
	return func(l *List) {
		l.definedParams.EntryTextHorizontalPosition = &h
		l.definedParams.EntryTextVerticalPosition = &v
	}
}

func (o ListOptions) EntrySelectedHandler(f ListEntrySelectedHandlerFunc) ListOpt {
	return func(l *List) {
		l.EntrySelectedEvent.AddHandler(func(args any) {
			if arg, ok := args.(*ListEntrySelectedEventArgs); ok {
				f(arg)
			}
		})
	}
}

func (o ListOptions) AllowReselect() ListOpt {
	return func(l *List) {
		l.definedParams.AllowReselect = constantutil.ConstantToPointer(true)
	}
}

// SelectFocus automatically selects each focused entry.
func (o ListOptions) SelectFocus() ListOpt {
	return func(l *List) {
		l.definedParams.SelectFocus = constantutil.ConstantToPointer(true)
	}
}

// SelectPressed selects entries when pressing instead of releasing (the default).
func (o ListOptions) SelectPressed() ListOpt {
	return func(l *List) {
		l.definedParams.SelectPressed = constantutil.ConstantToPointer(true)
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

func (l *List) Render(screen *ebiten.Image) {
	l.init.Do()

	d := l.container.GetWidget().Disabled
	if l.vSlider != nil {
		l.vSlider.DrawTrackDisabled = d
	}
	if l.hSlider != nil {
		l.hSlider.DrawTrackDisabled = d
	}
	l.scrollContainer.GetWidget().Disabled = d

	if l.focusIndex != l.prevFocusIndex && l.focusIndex >= 0 && l.focusIndex < len(l.buttons) {
		l.scrollVisible(l.buttons[l.focusIndex])
	}

	if *l.computedParams.SelectFocus {
		l.SelectFocused()
	}

	l.container.Render(screen)
}

func (l *List) Update() {
	l.init.Do()

	l.handleInput()
	l.container.Update()
}

/** Focuser Interface - Start **/

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

func (l *List) GetFocus(direction FocusDirection) Focuser {
	return l.focusMap[direction]
}

func (l *List) AddFocus(direction FocusDirection, focus Focuser) {
	l.focusMap[direction] = focus
}

/** Focuser Interface - End **/

func (l *List) handleInput() {
	if l.focused && !l.GetWidget().Disabled && len(l.buttons) > 0 {
		if !*l.computedParams.DisableDefaultKeys && (input.KeyPressed(ebiten.KeyUp) || input.KeyPressed(ebiten.KeyDown)) {
			if !l.justMoved {
				if input.KeyPressed(ebiten.KeyDown) {
					l.FocusNext()
				} else {
					l.FocusPrevious()
				}
			}
		} else {
			l.justMoved = false
		}
		l.buttons[l.focusIndex].focused = true
	} else if len(l.buttons) > 0 && l.focusIndex <= len(l.buttons) {
		l.buttons[l.focusIndex].focused = false
	}
}

func (l *List) FocusNext() {
	if len(l.buttons) > 0 {
		direction := 1
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
		l.buttons[l.focusIndex].focused = true
	}
}

func (l *List) FocusPrevious() {
	if len(l.buttons) > 0 {
		direction := -1
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
		l.buttons[l.focusIndex].focused = true
	}
}

func (l *List) SelectFocused() {
	if l.focusIndex >= 0 && l.focusIndex < len(l.buttons) {
		l.scrollVisible(l.buttons[l.focusIndex])
		l.setSelectedEntry(l.entries[l.focusIndex], false)
	}
}

func (l *List) resetFocusIndex() {
	if len(l.buttons) > 0 {
		if l.focusIndex != -1 && l.focusIndex < len(l.buttons) {
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
	l.layout = NewGridLayout(
		GridLayoutOpts.Columns(cols),
		GridLayoutOpts.Stretch([]bool{true, false}, []bool{true, false}))

	l.container = NewContainer(
		append([]ContainerOpt{
			ContainerOpts.WidgetOpts(WidgetOpts.TrackHover(true)),
			ContainerOpts.Layout(l.layout)}, l.containerOpts...)...,
	)

}

func (l *List) initWidget() {
	l.container.RemoveChildren()

	l.layout.columnSpacing = *l.computedParams.ControlWidgetSpacing
	l.layout.rowSpacing = *l.computedParams.ControlWidgetSpacing

	if l.computedParams.MinSize != nil {
		l.container.GetWidget().MinWidth = l.computedParams.MinSize.X
		l.container.GetWidget().MinHeight = l.computedParams.MinSize.Y
	}

	l.listContent = NewContainer(
		ContainerOpts.Layout(NewRowLayout(
			RowLayoutOpts.Direction(DirectionVertical))),
		ContainerOpts.AutoDisableChildren(),
	)

	l.buttons = make([]*Button, 0, len(l.entries))
	for _, e := range l.entries {
		e := e
		but := l.createEntry(e)

		l.buttons = append(l.buttons, but)
		l.listContent.AddChild(but)
	}

	l.scrollContainer = NewScrollContainer(
		ScrollContainerOpts.Content(l.listContent),
		ScrollContainerOpts.StretchContentWidth(),
		ScrollContainerOpts.Image(l.computedParams.ScrollContainerImage),
		ScrollContainerOpts.Padding(l.computedParams.ScrollContainerPadding),
	)

	l.container.AddChild(l.scrollContainer)

	var sliderOpts []SliderOpt
	if l.computedParams.Slider != nil {
		if l.computedParams.Slider.FixedHandleSize != nil {
			sliderOpts = append(sliderOpts, SliderOpts.FixedHandleSize(*l.computedParams.Slider.FixedHandleSize))
		} else {
			sliderOpts = append(sliderOpts, SliderOpts.FixedHandleSize(0))
		}
		if l.computedParams.Slider.HandleImage != nil {
			sliderOpts = append(sliderOpts, SliderOpts.HandleImage(l.computedParams.Slider.HandleImage))
		}
		if l.computedParams.Slider.TrackImage != nil {
			sliderOpts = append(sliderOpts, SliderOpts.TrackImage(l.computedParams.Slider.TrackImage))
		}
		if l.computedParams.Slider.MinHandleSize != nil {
			sliderOpts = append(sliderOpts, SliderOpts.MinHandleSize(*l.computedParams.Slider.MinHandleSize))
		}
		if l.computedParams.Slider.TrackOffset != nil {
			sliderOpts = append(sliderOpts, SliderOpts.TrackOffset(*l.computedParams.Slider.TrackOffset))
		}
		if l.computedParams.Slider.TrackPadding != nil {
			sliderOpts = append(sliderOpts, SliderOpts.TrackPadding(l.computedParams.Slider.TrackPadding))
		}
	}
	if l.computedParams.DisableDefaultKeys != nil {
		sliderOpts = append(sliderOpts, SliderOpts.DisableDefaultKeys(*l.computedParams.DisableDefaultKeys))
	}

	if !l.hideVerticalSlider {
		pageSizeFunc := func() int {
			return int(math.Round(float64(l.scrollContainer.ViewRect().Dy()) / float64(l.listContent.GetWidget().Rect.Dy()) * 1000))
		}

		l.vSlider = NewSlider(append(sliderOpts, []SliderOpt{
			SliderOpts.Direction(DirectionVertical),
			SliderOpts.MinMax(0, 1000),
			SliderOpts.PageSizeFunc(pageSizeFunc),
			SliderOpts.ChangedHandler(func(args *SliderChangedEventArgs) {
				current := args.Slider.Current
				if pageSizeFunc() >= 1000 {
					current = 0
				}
				l.scrollContainer.ScrollTop = float64(current) / 1000
			}),
		}...)...)
		l.container.AddChild(l.vSlider)

		l.scrollContainer.widget.ScrolledEvent.AddHandler(func(args any) {
			if a, ok := args.(*WidgetScrolledEventArgs); ok {
				p := pageSizeFunc() / 3
				if p < 1 {
					p = 1
				}
				l.vSlider.Current -= int(math.Round(a.Y * float64(p)))
			}
		})
	}

	if !l.hideHorizontalSlider {
		l.hSlider = NewSlider(append(sliderOpts, []SliderOpt{
			SliderOpts.Direction(DirectionHorizontal),
			SliderOpts.MinMax(0, 1000),
			SliderOpts.PageSizeFunc(func() int {
				return int(math.Round(float64(l.scrollContainer.ViewRect().Dx()) / float64(l.listContent.GetWidget().Rect.Dx()) * 1000))
			}),
			SliderOpts.ChangedHandler(func(args *SliderChangedEventArgs) {
				l.scrollContainer.ScrollLeft = float64(args.Slider.Current) / 1000
			}),
		}...)...)
		l.container.AddChild(l.hSlider)
	}

	if l.selectedEntry != nil {
		se := l.selectedEntry
		l.selectedEntry = nil
		l.SetSelectedEntry(se)
	}
}

// Updates the entries in the list.
// Note: Duplicates will be removed.
func (l *List) SetEntries(newEntries []any) {
	// Remove old entries
	for i := range l.entries {
		but := l.buttons[i]
		l.listContent.RemoveChild(but)
	}
	l.entries = nil
	l.buttons = nil

	// Add new Entries
	for idx := range newEntries {
		if !slices.ContainsFunc(l.entries, func(cmp any) bool {
			return cmp == newEntries[idx]
		}) {
			l.entries = append(l.entries, newEntries[idx])
			but := l.createEntry(newEntries[idx])
			l.buttons = append(l.buttons, but)
			l.listContent.AddChild(but)
		}
	}
	l.selectedEntry = nil
	l.resetFocusIndex()
}

// Remove the passed in entry from the list if it exists.
func (l *List) RemoveEntry(entry any) {
	l.init.Do()

	if len(l.entries) > 0 && entry != nil {
		for i, e := range l.entries {
			if e == entry {
				but := l.buttons[i]
				l.entries = append(l.entries[:i], l.entries[i+1:]...)
				l.buttons = append(l.buttons[:i], l.buttons[i+1:]...)
				l.listContent.RemoveChild(but)

				entryLen := len(l.entries)
				if l.focusIndex >= entryLen {
					l.focusIndex = i - 1
				}

				if l.focusIndex >= 0 && l.focusIndex < entryLen {
					l.setSelectedEntry(l.entries[l.focusIndex], false)
				}
				break
			}
		}
		l.resetFocusIndex()
	}
}

// Add a new entry to the end of the list
// Note: Duplicates will not be added.
func (l *List) AddEntry(entry any) {
	l.init.Do()
	if !l.checkForDuplicates(l.entries, entry) {
		l.entries = append(l.entries, entry)
		but := l.createEntry(entry)
		l.buttons = append(l.buttons, but)
		l.listContent.AddChild(but)
	}
	l.resetFocusIndex()
}

// Return the current entries in the list.
func (l *List) Entries() []any {
	l.init.Do()
	return l.entries
}

// Return the currently selected entry in the list.
func (l *List) SelectedEntry() any {
	l.init.Do()
	return l.selectedEntry
}

// Set the Selected Entry to e if it is found.
func (l *List) SetSelectedEntry(entry any) {
	l.setSelectedEntry(entry, false)
}

func (l *List) setSelectedEntry(e any, user bool) {
	if e != l.selectedEntry || (user && *l.computedParams.AllowReselect) {
		l.init.Do()

		prev := l.selectedEntry
		l.selectedEntry = e
		l.resetFocusIndex()
		for i := range l.buttons {
			if l.entries[i] == e {
				l.buttons[i].definedParams.Image = l.computedParams.entrySelectedColor
				l.buttons[i].definedParams.TextColor = l.computedParams.entryTextColor
				l.buttons[i].computedParams.Image = l.computedParams.entrySelectedColor
				l.buttons[i].computedParams.TextColor = l.computedParams.entryTextColor
			} else {
				l.buttons[i].definedParams.Image = l.computedParams.entryUnselectedColor
				l.buttons[i].definedParams.TextColor = l.computedParams.entryUnselectedTextColor
				l.buttons[i].computedParams.Image = l.computedParams.entryUnselectedColor
				l.buttons[i].computedParams.TextColor = l.computedParams.entryUnselectedTextColor
			}
		}

		l.EntrySelectedEvent.Fire(&ListEntrySelectedEventArgs{
			Entry:         e,
			PreviousEntry: prev,
		})
	}
}

func (l *List) checkForDuplicates(entries []any, entry any) bool {
	return slices.ContainsFunc(entries, func(cmp any) bool {
		return cmp == entry
	})
}

func (l *List) createEntry(entry any) *Button {
	but := NewButton(
		ButtonOpts.WidgetOpts(WidgetOpts.LayoutData(RowLayoutData{
			Stretch: true,
		})),
		ButtonOpts.Image(l.computedParams.entryUnselectedColor),
		ButtonOpts.Text(l.entryLabelFunc(entry), l.computedParams.EntryFace, l.computedParams.entryUnselectedTextColor),
		ButtonOpts.TextPadding(l.computedParams.EntryTextPadding),
		ButtonOpts.TextPosition(*l.computedParams.EntryTextHorizontalPosition, *l.computedParams.EntryTextVerticalPosition),
	)
	events := but.ClickedEvent
	if *l.computedParams.SelectPressed {
		events = but.PressedEvent
	}
	events.AddHandler(func(_ interface{}) {
		l.setSelectedEntry(entry, true)
	})
	but.Validate()
	return but
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

func (l *List) scrollVisible(w HasWidget) {
	vrect := l.scrollContainer.ViewRect()
	wrect := w.GetWidget().Rect
	if !wrect.In(vrect) {
		crect := l.scrollContainer.ContentRect()
		scrollTop := l.scrollContainer.ScrollTop
		scrollHeight := crect.Dy() - vrect.Dy()
		if wrect.Max.Y > vrect.Max.Y {
			scrollTop = float64(wrect.Max.Y-vrect.Dy()-crect.Min.Y) / float64(scrollHeight)
		} else if wrect.Min.Y < vrect.Min.Y {
			scrollTop = float64(wrect.Min.Y-crect.Min.Y) / float64(scrollHeight)
		}
		scrollLeft := l.scrollContainer.ScrollLeft
		scrollWidth := crect.Dx() - vrect.Dx()
		if wrect.Max.X > vrect.Max.X {
			scrollLeft = float64(wrect.Max.X-vrect.Dx()-crect.Min.X) / float64(scrollWidth)
		} else if wrect.Min.X < vrect.Min.X {
			scrollLeft = float64(wrect.Min.X-crect.Min.X) / float64(scrollWidth)
		}
		l.setScrollTop(scrollClamp(scrollTop, l.scrollContainer.ScrollTop))
		l.setScrollLeft(scrollClamp(scrollLeft, l.scrollContainer.ScrollLeft))
	} else if wrect != vrect {
		l.prevFocusIndex = l.focusIndex
	}
}

func scrollClamp(targetScroll, currentScroll float64) float64 {
	const maxScrollStep = 0.1
	minScroll := currentScroll - maxScrollStep
	maxScroll := currentScroll + maxScrollStep
	return math.Max(minScroll, math.Min(targetScroll, maxScroll))
}
