package widget

import (
	"image"

	"github.com/blizzy78/ebitenui/event"
	"github.com/blizzy78/ebitenui/input"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

type TabBook struct {
	TabSelectedEvent *event.Event

	tabs          []*TabBookTab
	containerOpts []ContainerOpt
	buttonOpts    []StateButtonOpt
	buttonImages  map[interface{}]*ButtonImage
	buttonFace    font.Face
	buttonColor   *ButtonTextColor
	flipBookOpts  []FlipBookOpt
	buttonSpacing int
	spacing       int

	init        *MultiOnce
	container   *Container
	tabToButton map[*TabBookTab]*StateButton
	flipBook    *FlipBook
	tab         *TabBookTab
}

type TabBookTab struct {
	Disabled bool

	label  string
	widget PreferredSizeLocateableWidget
}

type TabBookOpt func(t *TabBook)

type TabBookTabSelectedEventArgs struct {
	TabBook     *TabBook
	Tab         *TabBookTab
	PreviousTab *TabBookTab
}

type TabBookTabSelectedHandlerFunc func(args *TabBookTabSelectedEventArgs)

type TabBookOptions struct {
}

var TabBookOpts TabBookOptions

func NewTabBook(opts ...TabBookOpt) *TabBook {
	t := &TabBook{
		TabSelectedEvent: &event.Event{},

		init:        &MultiOnce{},
		tabToButton: map[*TabBookTab]*StateButton{},
	}

	t.init.Append(t.createWidget)

	for _, o := range opts {
		o(t)
	}

	return t
}

func NewTabBookTab(label string, widget PreferredSizeLocateableWidget) *TabBookTab {
	return &TabBookTab{
		label:  label,
		widget: widget,
	}
}

func (o TabBookOptions) ContainerOpts(opts ...ContainerOpt) TabBookOpt {
	return func(t *TabBook) {
		t.containerOpts = append(t.containerOpts, opts...)
	}
}

func (o TabBookOptions) TabButtonOpts(opts ...StateButtonOpt) TabBookOpt {
	return func(t *TabBook) {
		t.buttonOpts = append(t.buttonOpts, opts...)
	}
}

func (o TabBookOptions) FlipBookOpts(opts ...FlipBookOpt) TabBookOpt {
	return func(t *TabBook) {
		t.flipBookOpts = append(t.flipBookOpts, opts...)
	}
}

func (o TabBookOptions) TabButtonImage(idle *ButtonImage, selected *ButtonImage) TabBookOpt {
	return func(t *TabBook) {
		t.buttonImages = map[interface{}]*ButtonImage{
			false: idle,
			true:  selected,
		}
	}
}

func (o TabBookOptions) TabButtonSpacing(s int) TabBookOpt {
	return func(t *TabBook) {
		t.buttonSpacing = s
	}
}

func (o TabBookOptions) TabButtonText(face font.Face, color *ButtonTextColor) TabBookOpt {
	return func(t *TabBook) {
		t.buttonFace = face
		t.buttonColor = color
	}
}

func (o TabBookOptions) Spacing(s int) TabBookOpt {
	return func(t *TabBook) {
		t.spacing = s
	}
}

func (o TabBookOptions) Tabs(tabs ...*TabBookTab) TabBookOpt {
	return func(t *TabBook) {
		t.tabs = append(t.tabs, tabs...)
	}
}

func (o TabBookOptions) TabSelectedHandler(f TabBookTabSelectedHandlerFunc) TabBookOpt {
	return func(t *TabBook) {
		t.TabSelectedEvent.AddHandler(func(args interface{}) {
			f(args.(*TabBookTabSelectedEventArgs))
		})
	}
}

func (t *TabBook) GetWidget() *Widget {
	t.init.Do()
	return t.container.GetWidget()
}

func (t *TabBook) PreferredSize() (int, int) {
	t.init.Do()
	return t.container.PreferredSize()
}

func (t *TabBook) SetLocation(rect image.Rectangle) {
	t.init.Do()
	t.container.SetLocation(rect)
}

func (t *TabBook) RequestRelayout() {
	t.init.Do()
	t.container.RequestRelayout()
}

func (t *TabBook) SetupInputLayer(def input.DeferredSetupInputLayerFunc) {
	t.init.Do()
	t.container.SetupInputLayer(def)
}

func (t *TabBook) Render(screen *ebiten.Image, def DeferredRenderFunc) {
	t.init.Do()

	d := t.container.GetWidget().Disabled
	for tab, b := range t.tabToButton {
		b.GetWidget().Disabled = d || tab.Disabled
	}

	t.container.Render(screen, def)
}

func (t *TabBook) createWidget() {
	t.container = NewContainer(append(t.containerOpts, []ContainerOpt{
		ContainerOpts.Layout(NewGridLayout(
			GridLayoutOpts.Columns(1),
			GridLayoutOpts.Stretch([]bool{true}, []bool{false, true}),
			GridLayoutOpts.Spacing(0, t.spacing))),
		ContainerOpts.AutoDisableChildren(),
	}...)...)
	t.containerOpts = nil

	buttonsContainer := NewContainer(
		ContainerOpts.Layout(NewRowLayout(
			RowLayoutOpts.Spacing(t.buttonSpacing))))
	t.container.AddChild(buttonsContainer)

	for _, tab := range t.tabs {
		tab := tab
		b := NewStateButton(append(t.buttonOpts, []StateButtonOpt{
			StateButtonOpts.StateImages(t.buttonImages),
			StateButtonOpts.ButtonOpts(
				ButtonOpts.Text(tab.label, t.buttonFace, t.buttonColor),
				ButtonOpts.ClickedHandler(func(args *ButtonClickedEventArgs) {
					t.SetTab(tab)
				})),
		}...)...)
		buttonsContainer.AddChild(b)

		t.tabToButton[tab] = b
	}
	t.buttonOpts = nil
	t.buttonImages = nil

	t.flipBook = NewFlipBook(append(t.flipBookOpts,
		FlipBookOpts.ContainerOpts(ContainerOpts.AutoDisableChildren()))...)
	t.container.AddChild(t.flipBook)
	t.flipBookOpts = nil

	t.setTab(t.tabs[0], false)
}

func (t *TabBook) SetTab(tab *TabBookTab) {
	t.setTab(tab, true)
}

func (t *TabBook) setTab(tab *TabBookTab, fireEvent bool) {
	if tab != t.tab {
		previousTab := t.tab

		t.tab = tab

		t.flipBook.SetPage(tab.widget)

		for bt, b := range t.tabToButton {
			b.State = bt == tab
		}

		if fireEvent {
			t.TabSelectedEvent.Fire(&TabBookTabSelectedEventArgs{
				TabBook:     t,
				Tab:         tab,
				PreviousTab: previousTab,
			})
		}
	}
}

func (t *TabBook) Tab() *TabBookTab {
	return t.tab
}
