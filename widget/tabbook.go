package widget

import (
	"image"

	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/input"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

type TabBook struct {
	TabSelectedEvent *event.Event

	tabs          []*TabBookTab
	containerOpts []ContainerOpt
	buttonOpts    []ButtonOpt
	buttonImages  *ButtonImage
	buttonFace    font.Face
	buttonColor   *ButtonTextColor
	flipBookOpts  []FlipBookOpt
	buttonSpacing int
	spacing       int

	init        *MultiOnce
	container   *Container
	tabToButton map[*TabBookTab]*Button
	flipBook    *FlipBook
	tab         *TabBookTab
	initialTab  *TabBookTab
}

type TabBookTab struct {
	Container
	Disabled bool
	label    string
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
		tabToButton: map[*TabBookTab]*Button{},
	}

	t.init.Append(t.createWidget)

	for _, o := range opts {
		o(t)
	}

	return t
}

func NewTabBookTab(label string, opts ...ContainerOpt) *TabBookTab {
	c := &TabBookTab{
		label: label,
	}
	c.init = &MultiOnce{}
	c.init.Append(c.createWidget)

	//Set a default layout so that tabs use the full container
	c.widgetOpts = append(c.widgetOpts, WidgetOpts.LayoutData(AnchorLayoutData{
		StretchHorizontal:  true,
		StretchVertical:    true,
		HorizontalPosition: AnchorLayoutPositionCenter,
		VerticalPosition:   AnchorLayoutPositionCenter,
	}))

	for _, o := range opts {
		o(&c.Container)
	}
	return c
}

func (o TabBookOptions) ContainerOpts(opts ...ContainerOpt) TabBookOpt {
	return func(t *TabBook) {
		t.containerOpts = append(t.containerOpts, opts...)
	}
}

func (o TabBookOptions) TabButtonOpts(opts ...ButtonOpt) TabBookOpt {
	return func(t *TabBook) {
		t.buttonOpts = append(t.buttonOpts, opts...)
	}
}

func (o TabBookOptions) FlipBookOpts(opts ...FlipBookOpt) TabBookOpt {
	return func(t *TabBook) {
		t.flipBookOpts = append(t.flipBookOpts, opts...)
	}
}

func (o TabBookOptions) TabButtonImage(buttonImages *ButtonImage) TabBookOpt {
	return func(t *TabBook) {
		t.buttonImages = buttonImages
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

func (o TabBookOptions) InitialTab(tab *TabBookTab) TabBookOpt {
	return func(t *TabBook) {
		t.initialTab = tab
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

func (t *TabBook) GetDropTargets() []HasWidget {
	return t.container.GetDropTargets()
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
	}...)...)
	t.containerOpts = nil

	buttonsContainer := NewContainer(
		ContainerOpts.Layout(NewRowLayout(
			RowLayoutOpts.Spacing(t.buttonSpacing))))
	t.container.AddChild(buttonsContainer)

	btnElements := []RadioGroupElement{}
	var firstTab *TabBookTab

	for _, tab := range t.tabs {
		tab := tab
		btn := NewButton(
			append(t.buttonOpts,
				ButtonOpts.Image(t.buttonImages),
				ButtonOpts.Text(tab.label, t.buttonFace, t.buttonColor),
				ButtonOpts.WidgetOpts(WidgetOpts.CustomData(tab)),
			)...,
		)
		btnElements = append(btnElements, btn)
		buttonsContainer.AddChild(btn)
		t.tabToButton[tab] = btn
		if firstTab == nil {
			if t.initialTab == nil && !tab.Disabled {
				firstTab = tab
			} else if t.initialTab == tab && !tab.Disabled {
				firstTab = t.initialTab
			}
		}
	}
	//If we cannot find an initial tab default to to the first one
	if firstTab == nil {
		firstTab = t.tabs[0]
	}

	NewRadioGroup(
		RadioGroupOpts.Elements(btnElements...),
		RadioGroupOpts.InitialElement(t.tabToButton[firstTab]),
		RadioGroupOpts.ChangedHandler(func(args *RadioGroupChangedEventArgs) {
			tab := args.Active.(*Button).GetWidget().CustomData.(*TabBookTab)
			t.SetTab(tab)
		}))

	t.buttonOpts = nil
	t.buttonImages = nil

	t.flipBook = NewFlipBook(append(t.flipBookOpts,
		FlipBookOpts.ContainerOpts(ContainerOpts.AutoDisableChildren()))...)
	t.container.AddChild(t.flipBook)
	t.flipBookOpts = nil

}

// Set the current tab for the tab book.
//
//		Note: This method should only be called after the
//	 ui is running. To set the initial tab please use the
//	 TabBookOptions.InitialTab method during tabbook creation.
func (t *TabBook) SetTab(tab *TabBookTab) {
	if tab.Disabled {
		return
	}
	t.init.Do()

	if tab != t.tab {
		previousTab := t.tab

		t.tab = tab

		t.flipBook.SetPage(tab)

		for bt, b := range t.tabToButton {
			state := WidgetUnchecked
			if bt == tab {
				state = WidgetChecked
			}
			b.SetState(state)
		}

		t.TabSelectedEvent.Fire(&TabBookTabSelectedEventArgs{
			TabBook:     t,
			Tab:         tab,
			PreviousTab: previousTab,
		})
	}

}

// Return the currently selected tab
func (t *TabBook) Tab() *TabBookTab {
	return t.tab
}
