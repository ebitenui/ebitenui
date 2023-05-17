package widget

import (
	"image"

	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/input"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

type ListComboButton struct {
	EntrySelectedEvent *event.Event

	buttonOpts []SelectComboButtonOpt
	listOpts   []ListOpt

	init   *MultiOnce
	button *SelectComboButton
	list   *List

	tabOrder           int
	disableDefaultKeys bool
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

		init: &MultiOnce{},
	}

	l.init.Append(l.createWidget)

	for _, o := range opts {
		o(l)
	}

	return l
}

func (o ListComboButtonOptions) SelectComboButtonOpts(opts ...SelectComboButtonOpt) ListComboButtonOpt {
	return func(l *ListComboButton) {
		l.buttonOpts = append(l.buttonOpts, opts...)
	}
}

func (o ListComboButtonOptions) ListOpts(opts ...ListOpt) ListComboButtonOpt {
	return func(l *ListComboButton) {
		l.listOpts = append(l.listOpts, opts...)
	}
}

func (o ListComboButtonOptions) Text(face font.Face, image *ButtonImageImage, color *ButtonTextColor) ListComboButtonOpt {
	return o.SelectComboButtonOpts(SelectComboButtonOpts.ComboButtonOpts(ComboButtonOpts.ButtonOpts(ButtonOpts.TextAndImage("", face, image, color))))
}

func (o ListComboButtonOptions) EntryLabelFunc(button SelectComboButtonEntryLabelFunc, list ListEntryLabelFunc) ListComboButtonOpt {
	return func(l *ListComboButton) {
		l.buttonOpts = append(l.buttonOpts, SelectComboButtonOpts.EntryLabelFunc(button))
		l.listOpts = append(l.listOpts, ListOpts.EntryLabelFunc(list))
	}
}

func (o ListComboButtonOptions) EntrySelectedHandler(f ListComboButtonEntrySelectedHandlerFunc) ListComboButtonOpt {
	return func(l *ListComboButton) {
		l.EntrySelectedEvent.AddHandler(func(args interface{}) {
			f(args.(*ListComboButtonEntrySelectedEventArgs))
		})
	}
}

func (o ListComboButtonOptions) TabOrder(tabOrder int) ListComboButtonOpt {
	return func(sl *ListComboButton) {
		sl.tabOrder = tabOrder
	}
}

func (o ListComboButtonOptions) DisableDefaultKeys(val bool) ListComboButtonOpt {
	return func(sl *ListComboButton) {
		sl.disableDefaultKeys = val
	}
}
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

func (l *ListComboButton) GetWidget() *Widget {
	l.init.Do()
	return l.button.GetWidget()
}

func (l *ListComboButton) PreferredSize() (int, int) {
	l.init.Do()
	return l.button.PreferredSize()
}

func (l *ListComboButton) SetLocation(rect image.Rectangle) {
	l.init.Do()
	l.button.SetLocation(rect)
}

func (l *ListComboButton) SetupInputLayer(def input.DeferredSetupInputLayerFunc) {
	l.init.Do()
	l.button.SetupInputLayer(def)
}

func (l *ListComboButton) Render(screen *ebiten.Image, def DeferredRenderFunc) {
	l.init.Do()
	if l.button.button.button.focused {
		if !l.disableDefaultKeys {
			if input.KeyPressed(ebiten.KeyDown) || input.KeyPressed(ebiten.KeyUp) {
				l.SetContentVisible(true)
			}
		}
	}

	l.button.Render(screen, def)

}

func (l *ListComboButton) createWidget() {
	l.list = NewList(append(l.listOpts, []ListOpt{
		ListOpts.HideHorizontalSlider(),
		ListOpts.AllowReselect(),
		ListOpts.DisableDefaultKeys(l.disableDefaultKeys),
	}...)...)
	l.listOpts = nil

	l.button = NewSelectComboButton(append(l.buttonOpts,
		SelectComboButtonOpts.ComboButtonOpts(ComboButtonOpts.Content(l.list)),
	)...)
	l.buttonOpts = nil

	if len(l.list.entries) > 0 {
		firstEntry := l.list.entries[0]
		l.button.SetSelectedEntry(firstEntry)
		l.list.SetSelectedEntry(firstEntry)
	}

	l.button.EntrySelectedEvent.AddHandler(func(args interface{}) {
		a := args.(*SelectComboButtonEntrySelectedEventArgs)
		l.EntrySelectedEvent.Fire(&ListComboButtonEntrySelectedEventArgs{
			Button:        l,
			Entry:         a.Entry,
			PreviousEntry: a.PreviousEntry,
		})
	})

	l.list.EntrySelectedEvent.AddHandler(func(args interface{}) {
		a := args.(*ListEntrySelectedEventArgs)
		l.SetContentVisible(false)
		l.SetSelectedEntry(a.Entry)
	})
}

func (l *ListComboButton) SetSelectedEntry(e interface{}) {
	l.init.Do()
	l.button.SetSelectedEntry(e)
	l.list.setSelectedEntry(e, false)
}

func (l *ListComboButton) SelectedEntry() interface{} {
	l.init.Do()
	return l.button.SelectedEntry()
}

func (l *ListComboButton) SetContentVisible(v bool) {
	l.init.Do()
	l.list.Focus(v)
	l.button.SetContentVisible(v)
	if !v {
		l.list.resetFocusIndex()
	}
}

func (l *ListComboButton) ContentVisible() bool {
	l.init.Do()
	return l.button.ContentVisible()
}

func (l *ListComboButton) Label() string {
	l.init.Do()
	return l.button.Label()
}
