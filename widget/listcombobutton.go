package widget

import (
	"image"

	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/input"
	"golang.org/x/exp/slices"
	"github.com/ebitenui/ebitenui/utilities/constantutil"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type ListComboButtonParams struct {
	List               *ListParams
	Button             *ButtonParams
	DisableDefaultKeys *bool
	MaxContentHeight   *int
}

type ListComboButton struct {
	definedParams  ListComboButtonParams
	computedParams ListComboButtonParams

	EntrySelectedEvent *event.Event

	init       *MultiOnce
	widget     *Widget
	widgetOpts []WidgetOpt
	entries    []any

	button          *SelectComboButton
	list            *List
	buttonLabelFunc SelectComboButtonEntryLabelFunc
	listLabelFunc   ListEntryLabelFunc

	selectedEntry any

	tabOrder int
	focusMap map[FocusDirection]Focuser
}

type ListComboButtonOpt func(l *ListComboButton)

type ListComboButtonEntrySelectedEventArgs struct {
	Button        *ListComboButton
	Entry         interface{}
	PreviousEntry interface{}
}

type ListComboButtonEntrySelectedHandlerFunc func(args *ListComboButtonEntrySelectedEventArgs)

type ListComboButtonOptions struct {
}

var ListComboButtonOpts ListComboButtonOptions

func NewListComboButton(opts ...ListComboButtonOpt) *ListComboButton {
	l := &ListComboButton{
		EntrySelectedEvent: &event.Event{},

		init:     &MultiOnce{},
		focusMap: make(map[FocusDirection]Focuser),
	}

	l.init.Append(l.createWidget)

	for _, o := range opts {
		o(l)
	}

	return l
}

func (l *ListComboButton) Validate() {
	l.init.Do()
	l.populateComputedParams()

	l.initWidget()
}

func (t *ListComboButton) populateComputedParams() {
	params := ListComboButtonParams{Button: &ButtonParams{}, List: &ListParams{}}

	theme := t.GetWidget().GetTheme()

	// Set theme values
	if theme != nil {
		if theme.ListComboButtonTheme != nil {
			params.DisableDefaultKeys = theme.ListComboButtonTheme.DisableDefaultKeys
			params.MaxContentHeight = theme.ListComboButtonTheme.MaxContentHeight
			if theme.ListComboButtonTheme.Button != nil {
				params.Button.GraphicImage = theme.ListComboButtonTheme.Button.GraphicImage
				params.Button.GraphicPadding = theme.ListComboButtonTheme.Button.GraphicPadding
				params.Button.Image = theme.ListComboButtonTheme.Button.Image
				params.Button.MinSize = theme.ListComboButtonTheme.Button.MinSize
				params.Button.TextColor = theme.ListComboButtonTheme.Button.TextColor
				if theme.ListComboButtonTheme.Button.TextFace != nil {
					params.Button.TextFace = theme.ListComboButtonTheme.Button.TextFace
				} else {
					params.Button.TextFace = theme.DefaultFace
				}
				params.Button.TextPadding = theme.ListComboButtonTheme.Button.TextPadding
				params.Button.TextPosition = theme.ListComboButtonTheme.Button.TextPosition
			}
			if theme.ListComboButtonTheme.List != nil {
				params.List.AllowReselect = theme.ListComboButtonTheme.List.AllowReselect
				params.List.ControlWidgetSpacing = theme.ListComboButtonTheme.List.ControlWidgetSpacing
				params.List.DisableDefaultKeys = theme.ListComboButtonTheme.List.DisableDefaultKeys
				if theme.ListComboButtonTheme.List.EntryFace != nil {
					params.List.EntryFace = theme.ListComboButtonTheme.List.EntryFace
				} else {
					params.List.EntryFace = theme.DefaultFace
				}
				params.List.EntryColor = theme.ListComboButtonTheme.List.EntryColor
				params.List.EntryTextHorizontalPosition = theme.ListComboButtonTheme.List.EntryTextHorizontalPosition
				params.List.EntryTextPadding = theme.ListComboButtonTheme.List.EntryTextPadding
				params.List.EntryTextVerticalPosition = theme.ListComboButtonTheme.List.EntryTextVerticalPosition
				params.List.ScrollContainerImage = theme.ListComboButtonTheme.List.ScrollContainerImage
				params.List.ScrollContainerPadding = theme.ListComboButtonTheme.List.ScrollContainerPadding
				params.List.SelectFocus = theme.ListComboButtonTheme.List.SelectFocus
				params.List.SelectPressed = theme.ListComboButtonTheme.List.SelectPressed
				params.List.Slider = theme.ListComboButtonTheme.List.Slider
				params.List.MinSize = theme.ListComboButtonTheme.List.MinSize
			}
		}
	}

	// Set definedParam values
	if t.definedParams.DisableDefaultKeys != nil {
		params.DisableDefaultKeys = t.definedParams.DisableDefaultKeys
	}
	if t.definedParams.MaxContentHeight != nil {
		params.MaxContentHeight = t.definedParams.MaxContentHeight
	}
	if t.definedParams.Button != nil {
		if t.definedParams.Button.GraphicImage != nil {
			params.Button.GraphicImage = t.definedParams.Button.GraphicImage
		}
		if t.definedParams.Button.GraphicPadding != nil {
			params.Button.GraphicPadding = t.definedParams.Button.GraphicPadding
		}
		if t.definedParams.Button.Image != nil {
			params.Button.Image = t.definedParams.Button.Image
		}
		if t.definedParams.Button.MinSize != nil {
			params.Button.MinSize = t.definedParams.Button.MinSize
		}
		if t.definedParams.Button.TextColor != nil {
			params.Button.TextColor = t.definedParams.Button.TextColor
		}
		if t.definedParams.Button.TextFace != nil {
			params.Button.TextFace = t.definedParams.Button.TextFace
		}
		if t.definedParams.Button.TextPadding != nil {
			params.Button.TextPadding = t.definedParams.Button.TextPadding
		}
		if t.definedParams.Button.TextPosition != nil {
			params.Button.TextPosition = t.definedParams.Button.TextPosition
		}
	}
	if t.definedParams.List != nil {
		if t.definedParams.List.AllowReselect != nil {
			params.List.AllowReselect = t.definedParams.List.AllowReselect
		}
		if t.definedParams.List.ControlWidgetSpacing != nil {
			params.List.ControlWidgetSpacing = t.definedParams.List.ControlWidgetSpacing
		}
		if t.definedParams.List.DisableDefaultKeys != nil {
			params.List.DisableDefaultKeys = t.definedParams.List.DisableDefaultKeys
		}
		if t.definedParams.List.EntryFace != nil {
			params.List.EntryFace = t.definedParams.List.EntryFace
		}
		if t.definedParams.List.EntryColor != nil {
			params.List.EntryColor = t.definedParams.List.EntryColor
		}

		if t.definedParams.List.EntryTextHorizontalPosition != nil {
			params.List.EntryTextHorizontalPosition = t.definedParams.List.EntryTextHorizontalPosition
		}
		if t.definedParams.List.EntryTextPadding != nil {
			params.List.EntryTextPadding = t.definedParams.List.EntryTextPadding
		}
		if t.definedParams.List.EntryTextVerticalPosition != nil {
			params.List.EntryTextVerticalPosition = t.definedParams.List.EntryTextVerticalPosition
		}
		if t.definedParams.List.EntryColor != nil {
			params.List.EntryColor = t.definedParams.List.EntryColor
		}
		if t.definedParams.List.ScrollContainerImage != nil {
			params.List.ScrollContainerImage = t.definedParams.List.ScrollContainerImage
		}
		if t.definedParams.List.ScrollContainerPadding != nil {
			params.List.ScrollContainerPadding = t.definedParams.List.ScrollContainerPadding
		}
		if t.definedParams.List.SelectFocus != nil {
			params.List.SelectFocus = t.definedParams.List.SelectFocus
		}
		if t.definedParams.List.SelectPressed != nil {
			params.List.SelectPressed = t.definedParams.List.SelectPressed
		}
		if t.definedParams.List.Slider != nil {
			params.List.Slider = t.definedParams.List.Slider
		}
		if t.definedParams.List.MinSize != nil {
			params.List.MinSize = t.definedParams.List.MinSize
		}
	}
	// Set defaults
	if params.MaxContentHeight == nil {
		params.MaxContentHeight = constantutil.ConstantToPointer(200)
	}
	if params.DisableDefaultKeys == nil {
		params.DisableDefaultKeys = constantutil.ConstantToPointer(false)
	}

	t.computedParams = params
}

func (o ListComboButtonOptions) WidgetOpts(opts ...WidgetOpt) ListComboButtonOpt {
	return func(l *ListComboButton) {
		l.widgetOpts = append(l.widgetOpts, opts...)
	}
}

func (o ListComboButtonOptions) ButtonParams(buttonParams *ButtonParams) ListComboButtonOpt {
	return func(l *ListComboButton) {
		l.definedParams.Button = buttonParams
	}
}

func (o ListComboButtonOptions) ListParams(listParams *ListParams) ListComboButtonOpt {
	return func(l *ListComboButton) {
		l.definedParams.List = listParams
	}
}

func (o ListComboButtonOptions) Text(face *text.Face, image *GraphicImage, color *ButtonTextColor) ListComboButtonOpt {
	return func(l *ListComboButton) {
		if l.definedParams.Button == nil {
			l.definedParams.Button = &ButtonParams{}
		}
		l.definedParams.Button.GraphicImage = image
		l.definedParams.Button.TextColor = color
		l.definedParams.Button.TextFace = face
	}
}

func (o ListComboButtonOptions) Entries(e []any) ListComboButtonOpt {
	return func(l *ListComboButton) {
		l.entries = slices.CompactFunc(e, func(a any, b any) bool { return a == b })
	}
}

func (o ListComboButtonOptions) EntryLabelFunc(button SelectComboButtonEntryLabelFunc, list ListEntryLabelFunc) ListComboButtonOpt {
	return func(l *ListComboButton) {
		l.buttonLabelFunc = button
		l.listLabelFunc = list
	}
}

func (o ListComboButtonOptions) EntrySelectedHandler(f ListComboButtonEntrySelectedHandlerFunc) ListComboButtonOpt {
	return func(l *ListComboButton) {
		l.EntrySelectedEvent.AddHandler(func(args interface{}) {
			if arg, ok := args.(*ListComboButtonEntrySelectedEventArgs); ok {
				f(arg)
			}
		})
	}
}

func (o ListComboButtonOptions) MaxContentHeight(h int) ListComboButtonOpt {
	return func(c *ListComboButton) {
		c.definedParams.MaxContentHeight = &h
	}
}

func (o ListComboButtonOptions) TabOrder(tabOrder int) ListComboButtonOpt {
	return func(sl *ListComboButton) {
		sl.tabOrder = tabOrder
	}
}

func (o ListComboButtonOptions) DisableDefaultKeys(val bool) ListComboButtonOpt {
	return func(sl *ListComboButton) {
		sl.definedParams.DisableDefaultKeys = &val
	}
}

/** Focuser Interface - Start **/

func (l *ListComboButton) Focus(focused bool) {
	l.init.Do()
	l.GetWidget().FireFocusEvent(l, focused, image.Point{-1, -1})
	l.button.button.button.focused = focused
	if !focused {
		l.SetContentVisible(false)
	}
}

func (l *ListComboButton) IsFocused() bool {
	return l.button.button.button.focused
}

func (l *ListComboButton) TabOrder() int {
	return l.tabOrder
}

func (l *ListComboButton) GetFocus(direction FocusDirection) Focuser {
	return l.focusMap[direction]
}

func (l *ListComboButton) AddFocus(direction FocusDirection, focus Focuser) {
	l.focusMap[direction] = focus
}

/** Focuser Interface - End **/

func (l *ListComboButton) FocusNext() {
	if l.list != nil {
		l.SetContentVisible(true)
		l.list.FocusNext()
	}
}

func (l *ListComboButton) FocusPrevious() {
	if l.list != nil {
		l.SetContentVisible(true)
		l.list.FocusPrevious()
	}
}

func (l *ListComboButton) SelectFocused() {
	if l.list != nil {
		l.SetContentVisible(true)
		l.list.SelectFocused()
	}
}

func (l *ListComboButton) GetWidget() *Widget {
	l.init.Do()
	return l.widget
}

func (l *ListComboButton) PreferredSize() (int, int) {
	l.init.Do()
	return l.button.PreferredSize()
}

func (l *ListComboButton) SetLocation(rect image.Rectangle) {
	l.init.Do()
	l.widget.Rect = rect
	l.button.SetLocation(rect)
}

func (l *ListComboButton) SetupInputLayer(def input.DeferredSetupInputLayerFunc) {
	l.init.Do()
	l.button.SetupInputLayer(def)
}

func (l *ListComboButton) Render(screen *ebiten.Image) {
	l.init.Do()

	l.button.Render(screen)

}

func (l *ListComboButton) Update() {
	l.init.Do()

	l.button.Update()

	if l.button.button.button.focused {
		if !*l.computedParams.DisableDefaultKeys {
			if input.KeyPressed(ebiten.KeyDown) || input.KeyPressed(ebiten.KeyUp) {
				l.SetContentVisible(true)
			}
		}
	}

}

func (l *ListComboButton) createWidget() {
	l.widget = NewWidget(l.widgetOpts...)
}

func (l *ListComboButton) initWidget() {

	l.list = NewList(ListOpts.HideHorizontalSlider(), ListOpts.Entries(l.entries), ListOpts.EntryLabelFunc(l.listLabelFunc))
	l.list.definedParams = *l.computedParams.List
	l.list.definedParams.DisableDefaultKeys = l.computedParams.DisableDefaultKeys
	l.list.definedParams.AllowReselect = constantutil.ConstantToPointer(true)
	l.list.Validate()
	btnOpts := []ButtonOpt{
		ButtonOpts.Image(l.computedParams.Button.Image),
		ButtonOpts.Text("", l.computedParams.Button.TextFace, l.computedParams.Button.TextColor),
	}
	if l.computedParams.Button.MinSize != nil {
		btnOpts = append(btnOpts, ButtonOpts.WidgetOpts(WidgetOpts.MinSize(l.computedParams.Button.MinSize.X, l.computedParams.Button.MinSize.Y)))
	}
	if l.computedParams.Button.TextPadding != nil {
		btnOpts = append(btnOpts, ButtonOpts.TextPadding(l.computedParams.Button.TextPadding))
	}

	l.button = NewSelectComboButton(
		SelectComboButtonOpts.ComboButtonOpts(
			ComboButtonOpts.Content(l.list),
			ComboButtonOpts.MaxContentHeight(*l.computedParams.MaxContentHeight),
			ComboButtonOpts.ButtonOpts(btnOpts...),
		),
		SelectComboButtonOpts.EntryLabelFunc(l.buttonLabelFunc),
	)
	l.button.Validate()

	if len(l.list.entries) > 0 {
		firstEntry := l.list.entries[0]
		l.button.SetSelectedEntry(firstEntry)
		l.list.SetSelectedEntry(firstEntry)
	}

	l.button.EntrySelectedEvent.AddHandler(func(args interface{}) {
		if a, ok := args.(*SelectComboButtonEntrySelectedEventArgs); ok {
			l.EntrySelectedEvent.Fire(&ListComboButtonEntrySelectedEventArgs{
				Button:        l,
				Entry:         a.Entry,
				PreviousEntry: a.PreviousEntry,
			})
		}
	})

	l.list.EntrySelectedEvent.AddHandler(func(args interface{}) {
		if a, ok := args.(*ListEntrySelectedEventArgs); ok {
			l.SetContentVisible(false)
			l.SetSelectedEntry(a.Entry)
		}
	})

	if l.selectedEntry != nil {
		l.button.SetSelectedEntry(l.selectedEntry)
		l.list.setSelectedEntry(l.selectedEntry, false)
	}
}

func (l *ListComboButton) SetSelectedEntry(e interface{}) {
	l.init.Do()
	if l.button == nil || l.list == nil {
		l.selectedEntry = e
	} else {
		l.button.SetSelectedEntry(e)
		l.list.setSelectedEntry(e, false)
	}
}

func (l *ListComboButton) SelectedEntry() interface{} {
	l.init.Do()
	if l.button == nil {
		return l.selectedEntry
	} else {
		return l.button.SelectedEntry()
	}
}

func (l *ListComboButton) SetContentVisible(v bool) {
	if l.list != nil {
		l.init.Do()
		l.list.Focus(v)
		l.button.SetContentVisible(v)
		if !v {
			l.list.resetFocusIndex()
		}
	}
}

func (l *ListComboButton) ContentVisible() bool {
	l.init.Do()
	if l.button == nil {
		return false
	} else {
		return l.button.ContentVisible()
	}
}

func (l *ListComboButton) Label() string {
	l.init.Do()
	if l.button == nil {
		return ""
	} else {
		return l.button.Label()
	}
}
