package widget

import (
	"image"

	"github.com/blizzy78/ebitenui/input"
	"github.com/hajimehoshi/ebiten"
	"golang.org/x/image/font"
)

type TabBook struct {
	tabs          []*TabBookTab
	containerOpts []ContainerOpt
	buttonOpts    []ButtonOpt
	buttonFace    font.Face
	buttonColor   *ButtonTextColor
	flipBookOpts  []FlipBookOpt
	buttonSpacing int
	spacing       int

	init          *MultiOnce
	container     *Container
	buttonsToTabs map[*Button]*TabBookTab
}

type TabBookTab struct {
	Disabled bool

	label  string
	widget HasWidget
}

type TabBookOpt func(t *TabBook)

const TabBookOpts = tabBookOpts(true)

type tabBookOpts bool

func NewTabBook(opts ...TabBookOpt) *TabBook {
	t := &TabBook{
		init:          &MultiOnce{},
		buttonsToTabs: map[*Button]*TabBookTab{},
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

func (o tabBookOpts) WithTabButtonOpts(opts ...ButtonOpt) TabBookOpt {
	return func(t *TabBook) {
		t.buttonOpts = append(t.buttonOpts, opts...)
	}
}

func (o tabBookOpts) WithFlipBookOpts(opts ...FlipBookOpt) TabBookOpt {
	return func(t *TabBook) {
		t.flipBookOpts = append(t.flipBookOpts, opts...)
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
	for b, tab := range t.buttonsToTabs {
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

	var f *FlipBook

	for _, tab := range t.tabs {
		b := NewButton(append(t.buttonOpts, []ButtonOpt{
			ButtonOpts.WithText(tab.label, t.buttonFace, t.buttonColor),
			ButtonOpts.WithClickedHandler(func(args *ButtonClickedEventArgs) {
				tab := t.buttonsToTabs[args.Button]
				f.SetPage(tab.widget)
			}),
		}...)...)
		buttonsContainer.AddChild(b)

		t.buttonsToTabs[b] = tab
	}
	t.buttonOpts = nil

	f = NewFlipBook(append(t.flipBookOpts,
		FlipBookOpts.WithContainerOpts(ContainerOpts.WithAutoDisableChildren()))...)
	t.container.AddChild(f)
	t.flipBookOpts = nil

	f.SetPage(t.tabs[0].widget)
}
