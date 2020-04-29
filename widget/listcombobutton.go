package widget

import (
	"image"

	"github.com/blizzy78/ebitenui/event"
	"github.com/blizzy78/ebitenui/input"

	"github.com/hajimehoshi/ebiten"
	"golang.org/x/image/font"
)

type ListComboButton struct {
	EntrySelectedEvent *event.Event

	buttonOpts []SelectComboButtonOpt
	listOpts   []ListOpt
	firstEntry interface{}

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

func (o listComboButtonOpts) WithLayoutData(ld interface{}) ListComboButtonOpt {
	return func(l *ListComboButton) {
		l.buttonOpts = append(l.buttonOpts, SelectComboButtonOpts.WithLayoutData(ld))
	}
}

func (o listComboButtonOpts) WithImage(i *ButtonImage) ListComboButtonOpt {
	return func(l *ListComboButton) {
		l.buttonOpts = append(l.buttonOpts, SelectComboButtonOpts.WithImage(i))
	}
}

func (o listComboButtonOpts) WithText(face font.Face, image *ButtonImageImage, color *ButtonTextColor) ListComboButtonOpt {
	return func(l *ListComboButton) {
		l.buttonOpts = append(l.buttonOpts, SelectComboButtonOpts.WithTextAndImage("", face, image, color))
	}
}

func (o listComboButtonOpts) WithListImage(i *ScrollContainerImage) ListComboButtonOpt {
	return func(l *ListComboButton) {
		l.listOpts = append(l.listOpts, ListOpts.WithImage(i))
	}
}

func (o listComboButtonOpts) WithListPadding(p Insets) ListComboButtonOpt {
	return func(l *ListComboButton) {
		l.listOpts = append(l.listOpts, ListOpts.WithPadding(p))
	}
}

func (o listComboButtonOpts) WithListSliderImages(track *SliderTrackImage, handle *ButtonImage) ListComboButtonOpt {
	return func(l *ListComboButton) {
		l.listOpts = append(l.listOpts, ListOpts.WithSliderImages(track, handle))
	}
}

func (o listComboButtonOpts) WithEntries(e []interface{}) ListComboButtonOpt {
	return func(l *ListComboButton) {
		l.listOpts = append(l.listOpts, ListOpts.WithEntries(e))
		if len(e) > 0 {
			l.firstEntry = e[0]
		}
	}
}

func (o listComboButtonOpts) WithEntryLabelFunc(button SelectComboButtonEntryLabelFunc, list ListEntryLabelFunc) ListComboButtonOpt {
	return func(l *ListComboButton) {
		l.buttonOpts = append(l.buttonOpts, SelectComboButtonOpts.WithEntryLabelFunc(button))
		l.listOpts = append(l.listOpts, ListOpts.WithEntryLabelFunc(list))
	}
}

func (o listComboButtonOpts) WithEntryFontFace(f font.Face) ListComboButtonOpt {
	return func(l *ListComboButton) {
		l.listOpts = append(l.listOpts, ListOpts.WithEntryFontFace(f))
	}
}

func (o listComboButtonOpts) WithEntryColor(c *ListEntryColor) ListComboButtonOpt {
	return func(l *ListComboButton) {
		l.listOpts = append(l.listOpts, ListOpts.WithEntryColor(c))
	}
}

func (o listComboButtonOpts) WithEntrySelectedHandler(f ListComboButtonEntrySelectedHandlerFunc) ListComboButtonOpt {
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
		ListOpts.WithControlWidgetSpacing(2),
		ListOpts.WithHideHorizontalSlider(),
		ListOpts.WithAllowReselect(),
	}...)...)
	l.listOpts = nil

	l.button = NewSelectComboButton(append(l.buttonOpts,
		SelectComboButtonOpts.WithContent(l.list),
	)...)
	l.buttonOpts = nil

	l.button.SetSelectedEntry(l.firstEntry)
	l.list.SetSelectedEntry(l.firstEntry)

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
