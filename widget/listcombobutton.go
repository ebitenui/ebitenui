package widget

import (
	"image"

	"github.com/blizzy78/ebitenui/event"
	"github.com/blizzy78/ebitenui/input"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

type ListComboButton struct {
	EntrySelectedEvent *event.Event

	buttonOpts []SelectComboButtonOpt
	listOpts   []ListOpt

	init               *MultiOnce
	button             *SelectComboButton
	list               *List
	lastContentVisible bool
}

type ListComboButtonOpt func(l *ListComboButton)

type ListComboButtonEntrySelectedEventArgs struct {
	Button        *ListComboButton
	Entry         interface{}
	PreviousEntry interface{}
}

type ListComboButtonEntrySelectedHandlerFunc func(args *ListComboButtonEntrySelectedEventArgs)

const ListComboButtonOpts = listComboButtonOpts(true)

type listComboButtonOpts bool

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

func (o listComboButtonOpts) SelectComboButtonOpts(opts ...SelectComboButtonOpt) ListComboButtonOpt {
	return func(l *ListComboButton) {
		l.buttonOpts = append(l.buttonOpts, opts...)
	}
}

func (o listComboButtonOpts) ListOpts(opts ...ListOpt) ListComboButtonOpt {
	return func(l *ListComboButton) {
		l.listOpts = append(l.listOpts, opts...)
	}
}

func (o listComboButtonOpts) Text(face font.Face, image *ButtonImageImage, color *ButtonTextColor) ListComboButtonOpt {
	return o.SelectComboButtonOpts(SelectComboButtonOpts.ComboButtonOpts(ComboButtonOpts.ButtonOpts(ButtonOpts.TextAndImage("", face, image, color))))
}

func (o listComboButtonOpts) EntryLabelFunc(button SelectComboButtonEntryLabelFunc, list ListEntryLabelFunc) ListComboButtonOpt {
	return func(l *ListComboButton) {
		l.buttonOpts = append(l.buttonOpts, SelectComboButtonOpts.EntryLabelFunc(button))
		l.listOpts = append(l.listOpts, ListOpts.EntryLabelFunc(list))
	}
}

func (o listComboButtonOpts) EntrySelectedHandler(f ListComboButtonEntrySelectedHandlerFunc) ListComboButtonOpt {
	return func(l *ListComboButton) {
		l.EntrySelectedEvent.AddHandler(func(args interface{}) {
			f(args.(*ListComboButtonEntrySelectedEventArgs))
		})
	}
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

	v := l.ContentVisible()
	if v && v != l.lastContentVisible {
		// TODO: scroll list to make current selected entry visible
		l.list.SetScrollTop(0)
	}

	l.button.Render(screen, def)

	l.lastContentVisible = v
}

func (l *ListComboButton) createWidget() {
	l.list = NewList(append(l.listOpts, []ListOpt{
		ListOpts.HideHorizontalSlider(),
		ListOpts.AllowReselect(),
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
}

func (l *ListComboButton) SelectedEntry() interface{} {
	l.init.Do()
	return l.button.SelectedEntry()
}

func (l *ListComboButton) SetContentVisible(v bool) {
	l.init.Do()
	l.button.SetContentVisible(v)
}

func (l *ListComboButton) ContentVisible() bool {
	l.init.Do()
	return l.button.ContentVisible()
}

func (l *ListComboButton) Label() string {
	l.init.Do()
	return l.button.Label()
}
