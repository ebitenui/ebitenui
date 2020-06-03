package widget

import (
	"image"

	"github.com/blizzy78/ebitenui/event"
	"github.com/blizzy78/ebitenui/input"

	"github.com/hajimehoshi/ebiten"
)

type SelectComboButton struct {
	EntrySelectedEvent *event.Event

	buttonOpts     []ComboButtonOpt
	entryLabelFunc SelectComboButtonEntryLabelFunc

	init          *MultiOnce
	button        *ComboButton
	selectedEntry interface{}
}

type SelectComboButtonOpt func(s *SelectComboButton)

type SelectComboButtonEntryLabelFunc func(e interface{}) string

type SelectComboButtonEntrySelectedEventArgs struct {
	Button        *SelectComboButton
	Entry         interface{}
	PreviousEntry interface{}
}

type SelectComboButtonEntrySelectedHandlerFunc func(args *SelectComboButtonEntrySelectedEventArgs)

const SelectComboButtonOpts = selectComboButtonOpts(true)

type selectComboButtonOpts bool

func NewSelectComboButton(opts ...SelectComboButtonOpt) *SelectComboButton {
	s := &SelectComboButton{
		EntrySelectedEvent: &event.Event{},

		init: &MultiOnce{},
	}

	s.buttonOpts = append(s.buttonOpts, ComboButtonOpts.WithMaxContentHeight(200))

	s.init.Append(s.createWidget)

	for _, o := range opts {
		o(s)
	}

	return s
}

func (o selectComboButtonOpts) WithComboButtonOpt(opt ComboButtonOpt) SelectComboButtonOpt {
	return func(s *SelectComboButton) {
		s.buttonOpts = append(s.buttonOpts, opt)
	}
}

func (o selectComboButtonOpts) WithEntrySelectedHandler(f SelectComboButtonEntrySelectedHandlerFunc) SelectComboButtonOpt {
	return func(s *SelectComboButton) {
		s.EntrySelectedEvent.AddHandler(func(args interface{}) {
			f(args.(*SelectComboButtonEntrySelectedEventArgs))
		})
	}
}

func (o selectComboButtonOpts) WithEntryLabelFunc(f SelectComboButtonEntryLabelFunc) SelectComboButtonOpt {
	return func(s *SelectComboButton) {
		s.entryLabelFunc = f
	}
}

func (s *SelectComboButton) GetWidget() *Widget {
	s.init.Do()
	return s.button.GetWidget()
}

func (s *SelectComboButton) SetLocation(rect image.Rectangle) {
	s.init.Do()
	s.button.SetLocation(rect)
}

func (s *SelectComboButton) SetupInputLayer(def input.DeferredSetupInputLayerFunc) {
	s.init.Do()
	s.button.SetupInputLayer(def)
}

func (s *SelectComboButton) PreferredSize() (int, int) {
	s.init.Do()
	return s.button.PreferredSize()
}

func (s *SelectComboButton) SetLabel(l string) {
	s.init.Do()
	s.button.SetLabel(l)
}

func (s *SelectComboButton) Render(screen *ebiten.Image, def DeferredRenderFunc) {
	s.init.Do()
	s.button.Render(screen, def)
}

func (s *SelectComboButton) createWidget() {
	s.button = NewComboButton(s.buttonOpts...)
	s.buttonOpts = nil
}

func (s *SelectComboButton) SetSelectedEntry(e interface{}) {
	if e != s.selectedEntry {
		s.init.Do()

		prev := s.selectedEntry
		s.selectedEntry = e

		s.button.SetLabel(s.entryLabelFunc(e))

		s.EntrySelectedEvent.Fire(&SelectComboButtonEntrySelectedEventArgs{
			Button:        s,
			Entry:         e,
			PreviousEntry: prev,
		})
	}
}

func (s *SelectComboButton) SelectedEntry() interface{} {
	return s.selectedEntry
}

func (s *SelectComboButton) SetContentVisible(v bool) {
	s.init.Do()
	s.button.ContentVisible = v
}

func (s *SelectComboButton) ContentVisible() bool {
	s.init.Do()
	return s.button.ContentVisible
}

func (s *SelectComboButton) Label() string {
	return s.button.Label()
}
