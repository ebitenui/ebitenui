package widget

import (
	"image"

	"github.com/blizzy78/ebitenui/event"
	"github.com/blizzy78/ebitenui/input"
	"github.com/hajimehoshi/ebiten"
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
	widget HasWidget
}

type TabBookOpt func(t *TabBook)

type TabBookTabSelectedEventArgs struct {
	TabBook     *TabBook
	Tab         *TabBookTab
	PreviousTab *TabBookTab
}

type TabBookTabSelectedHandlerFunc func(args *TabBookTabSelectedEventArgs)

const TabBookOpts = tabBookOpts(true)

type tabBookOpts bool

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

func NewTabBookTab(label string, widget HasWidget) *TabBookTab {
	return &TabBookTab{
		label:  label,
		widget: widget,
	}
}

func (o tabBookOpts) WithContainerOpts(opts ...ContainerOpt) TabBookOpt {
	return func(t *TabBook) {
		t.containerOpts = append(t.containerOpts, opts...)
	}
}

func (o tabBookOpts) WithTabButtonOpts(opts ...StateButtonOpt) TabBookOpt {
	return func(t *TabBook) {
		t.buttonOpts = append(t.buttonOpts, opts...)
	}
}

func (o tabBookOpts) WithFlipBookOpts(opts ...FlipBookOpt) TabBookOpt {
	return func(t *TabBook) {
		t.flipBookOpts = append(t.flipBookOpts, opts...)
	}
}

func (o tabBookOpts) WithTabButtonImage(idle *ButtonImage, selected *ButtonImage) TabBookOpt {
	return func(t *TabBook) {
		t.buttonImages = map[interface{}]*ButtonImage{
			false: idle,
			true:  selected,
		}
	}
}

func (o tabBookOpts) WithTabButtonSpacing(s int) TabBookOpt {
	return func(t *TabBook) {
		t.buttonSpacing = s
	}
}

func (o tabBookOpts) WithTabButtonText(face font.Face, color *ButtonTextColor) TabBookOpt {
	return func(t *TabBook) {
		t.buttonFace = face
		t.buttonColor = color
	}
}

func (o tabBookOpts) WithSpacing(s int) TabBookOpt {
	return func(t *TabBook) {
		t.spacing = s
	}
}

func (o tabBookOpts) WithTabs(tabs ...*TabBookTab) TabBookOpt {
	return func(t *TabBook) {
		t.tabs = append(t.tabs, tabs...)
	}
}

func (o tabBookOpts) WithTabSelectedHandler(f TabBookTabSelectedHandlerFunc) TabBookOpt {
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
		ContainerOpts.WithLayout(NewGridLayout(
			GridLayoutOpts.WithColumns(1),
			GridLayoutOpts.WithStretch([]bool{true}, []bool{false, true}),
			GridLayoutOpts.WithSpacing(0, t.spacing))),
		ContainerOpts.WithAutoDisableChildren(),
	}...)...)
	t.containerOpts = nil

	buttonsContainer := NewContainer(
		ContainerOpts.WithLayout(NewRowLayout(
			RowLayoutOpts.WithSpacing(t.buttonSpacing))))
	t.container.AddChild(buttonsContainer)

	for _, tab := range t.tabs {
		tab := tab
		b := NewStateButton(append(t.buttonOpts, []StateButtonOpt{
			StateButtonOpts.WithStateImages(t.buttonImages),
			StateButtonOpts.WithButtonOpts(
				ButtonOpts.WithText(tab.label, t.buttonFace, t.buttonColor),
				ButtonOpts.WithClickedHandler(func(args *ButtonClickedEventArgs) {
					t.SetTab(tab)
				})),
		}...)...)
		buttonsContainer.AddChild(b)

		t.tabToButton[tab] = b
	}
	t.buttonOpts = nil
	t.buttonImages = nil

	t.flipBook = NewFlipBook(append(t.flipBookOpts,
		FlipBookOpts.WithContainerOpts(ContainerOpts.WithAutoDisableChildren()))...)
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
