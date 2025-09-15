package widget

import (
	"image"

	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/input"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type TabBookParams struct {
	TabButton      *ButtonParams
	TabSpacing     *int
	ContentSpacing *int
	ContentPadding *Insets
}

type TabBook struct {
	definedParams  TabBookParams
	computedParams TabBookParams

	TabSelectedEvent *event.Event

	tabs          []*TabBookTab
	containerOpts []ContainerOpt

	init         *MultiOnce
	container    *Container
	btnContainer *Container
	gridLayout   *GridLayout
	tabToButton  map[*TabBookTab]*Button
	flipBook     *FlipBook
	tab          *TabBookTab
	initialTab   *TabBookTab
}

type TabBookOpt func(t *TabBook)

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

func (t *TabBook) Validate() {
	t.init.Do()
	t.populateComputedParams()
	if len(t.tabs) == 0 {
		panic("TabBook: At least one tab is required.")
	}
	if t.computedParams.TabButton == nil {
		panic("TabBook: TabButton parameters are required.")
	}
	if t.computedParams.TabButton.TextColor == nil {
		panic("TabBook: TabButtonText Color is required.")
	}
	if t.computedParams.TabButton.TextColor.Idle == nil {
		panic("TabBook: TabButtonText Color.Idle is required.")
	}
	if t.computedParams.TabButton.TextFace == nil {
		panic("TabBook: TabButtonText Font Face is required.")
	}
	if t.computedParams.TabButton.Image == nil {
		panic("TabBook: TabButtonImage is required.")
	}
	if t.computedParams.TabButton.Image.Idle == nil {
		panic("TabBook: TabButtonImage.Idle is required.")
	}
	if t.computedParams.TabButton.Image.Pressed == nil {
		panic("TabBook: TabButtonImage.Pressed is required.")
	}

	t.initTabBook()
}

func (t *TabBook) populateComputedParams() {
	ZERO := 0
	params := TabBookParams{ContentSpacing: &ZERO, TabSpacing: &ZERO, ContentPadding: &Insets{}}

	theme := t.GetWidget().GetTheme()

	if theme != nil {
		if theme.TabbookTheme != nil {
			params.ContentPadding = theme.TabbookTheme.ContentPadding
			params.ContentSpacing = theme.TabbookTheme.ContentSpacing
			params.TabSpacing = theme.TabbookTheme.TabSpacing
			params.TabButton = theme.TabbookTheme.TabButton
		}
		if params.TabButton == nil {
			params.TabButton = &ButtonParams{TextColor: &ButtonTextColor{}}
		}
		if params.TabButton.TextFace == nil {
			params.TabButton.TextFace = theme.DefaultFace
		}
		if params.TabButton.TextColor.Idle == nil {
			params.TabButton.TextColor.Idle = theme.DefaultTextColor
		}
	}

	if t.definedParams.ContentPadding != nil {
		params.ContentPadding = t.definedParams.ContentPadding
	}
	if t.definedParams.ContentSpacing != nil {
		params.ContentSpacing = t.definedParams.ContentSpacing
	}
	if t.definedParams.TabSpacing != nil {
		params.TabSpacing = t.definedParams.TabSpacing
	}
	if params.TabButton == nil {
		params.TabButton = &ButtonParams{TextColor: &ButtonTextColor{}}
	}
	if t.definedParams.TabButton != nil {
		if t.definedParams.TabButton.Image != nil {
			if params.TabButton.Image == nil {
				params.TabButton.Image = t.definedParams.TabButton.Image
			} else {
				if t.definedParams.TabButton.Image.Idle != nil {
					params.TabButton.Image.Idle = t.definedParams.TabButton.Image.Idle
				}
				if t.definedParams.TabButton.Image.Hover != nil {
					params.TabButton.Image.Hover = t.definedParams.TabButton.Image.Hover
				}
				if t.definedParams.TabButton.Image.Pressed != nil {
					params.TabButton.Image.Pressed = t.definedParams.TabButton.Image.Pressed
				}
				if t.definedParams.TabButton.Image.PressedHover != nil {
					params.TabButton.Image.PressedHover = t.definedParams.TabButton.Image.PressedHover
				}
				if t.definedParams.TabButton.Image.Disabled != nil {
					params.TabButton.Image.Disabled = t.definedParams.TabButton.Image.Disabled
				}
			}
		}

		if t.definedParams.TabButton.GraphicImage != nil {
			if params.TabButton.GraphicImage == nil {
				params.TabButton.GraphicImage = t.definedParams.TabButton.GraphicImage
			} else {
				if t.definedParams.TabButton.GraphicImage.Idle != nil {
					params.TabButton.GraphicImage.Idle = t.definedParams.TabButton.GraphicImage.Idle
				}
				if t.definedParams.TabButton.GraphicImage.Disabled != nil {
					params.TabButton.GraphicImage.Disabled = t.definedParams.TabButton.GraphicImage.Disabled
				}
			}
		}
		if t.definedParams.TabButton.GraphicPadding != nil {
			params.TabButton.GraphicPadding = t.definedParams.TabButton.GraphicPadding
		}
		if t.definedParams.TabButton.TextPosition != nil {
			params.TabButton.TextPosition.HTextPosition = t.definedParams.TabButton.TextPosition.HTextPosition
			params.TabButton.TextPosition.VTextPosition = t.definedParams.TabButton.TextPosition.VTextPosition
		}
		if t.definedParams.TabButton.TextFace != nil {
			params.TabButton.TextFace = t.definedParams.TabButton.TextFace
		}
		if t.definedParams.TabButton.TextPadding != nil {
			params.TabButton.TextPadding = t.definedParams.TabButton.TextPadding
		}
		if t.definedParams.TabButton.MinSize != nil {
			params.TabButton.MinSize = t.definedParams.TabButton.MinSize
		}

		if t.definedParams.TabButton.TextColor != nil {
			if params.TabButton.TextColor == nil {
				params.TabButton.TextColor = t.definedParams.TabButton.TextColor
			} else {
				if t.definedParams.TabButton.TextColor.Disabled != nil {
					params.TabButton.TextColor.Disabled = t.definedParams.TabButton.TextColor.Disabled
				}
				if t.definedParams.TabButton.TextColor.Hover != nil {
					params.TabButton.TextColor.Hover = t.definedParams.TabButton.TextColor.Hover
				}
				if t.definedParams.TabButton.TextColor.Idle != nil {
					params.TabButton.TextColor.Idle = t.definedParams.TabButton.TextColor.Idle
				}
				if t.definedParams.TabButton.TextColor.Pressed != nil {
					params.TabButton.TextColor.Pressed = t.definedParams.TabButton.TextColor.Pressed
				}
			}
		}
	}
	if params.ContentSpacing == nil {
		params.ContentSpacing = &ZERO
	}

	t.computedParams = params

}

func (o TabBookOptions) ContainerOpts(opts ...ContainerOpt) TabBookOpt {
	return func(t *TabBook) {
		t.containerOpts = append(t.containerOpts, opts...)
	}
}

func (o TabBookOptions) ContentPadding(contentPadding *Insets) TabBookOpt {
	return func(t *TabBook) {
		t.definedParams.ContentPadding = contentPadding
	}
}

func (o TabBookOptions) TabButtonImage(buttonImages *ButtonImage) TabBookOpt {
	return func(t *TabBook) {
		if t.definedParams.TabButton == nil {
			t.definedParams.TabButton = &ButtonParams{}
		}
		t.definedParams.TabButton.Image = buttonImages
	}
}

func (o TabBookOptions) TabButtonSpacing(s int) TabBookOpt {
	return func(t *TabBook) {
		t.definedParams.TabSpacing = &s
	}
}

func (o TabBookOptions) TabButtonText(face *text.Face, color *ButtonTextColor) TabBookOpt {
	return func(t *TabBook) {
		if t.definedParams.TabButton == nil {
			t.definedParams.TabButton = &ButtonParams{}
		}
		t.definedParams.TabButton.TextFace = face
		t.definedParams.TabButton.TextColor = color
	}
}

func (o TabBookOptions) TabButtonTextPadding(textPadding *Insets) TabBookOpt {
	return func(t *TabBook) {
		if t.definedParams.TabButton == nil {
			t.definedParams.TabButton = &ButtonParams{}
		}
		t.definedParams.TabButton.TextPadding = textPadding
	}
}
func (o TabBookOptions) TabButtonMinSize(minSize *image.Point) TabBookOpt {
	return func(t *TabBook) {
		if t.definedParams.TabButton == nil {
			t.definedParams.TabButton = &ButtonParams{}
		}
		t.definedParams.TabButton.MinSize = minSize
	}
}

func (o TabBookOptions) ContentSpacing(s int) TabBookOpt {
	return func(t *TabBook) {
		t.definedParams.ContentSpacing = &s
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
			if arg, ok := args.(*TabBookTabSelectedEventArgs); ok {
				f(arg)
			}
		})
	}
}

func (t *TabBook) GetWidget() *Widget {
	t.init.Do()
	return t.container.GetWidget()
}

func (t *TabBook) PreferredSize() (int, int) {
	t.init.Do()
	x, y := t.container.PreferredSize()
	_, bcY := t.btnContainer.PreferredSize()
	for tab := range t.tabs {
		tx, ty := t.tabs[tab].PreferredSize()
		ty += bcY
		if tx > x {
			x = tx
		}
		if ty > y {
			y = ty
		}
	}
	return x, y
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

func (t *TabBook) Render(screen *ebiten.Image) {
	t.init.Do()

	d := t.container.GetWidget().Disabled
	for tab, b := range t.tabToButton {
		b.GetWidget().Disabled = d || tab.Disabled
	}

	t.container.Render(screen)
}

func (t *TabBook) Update(updObj *UpdateObject) {
	t.init.Do()

	t.container.Update(updObj)
}

func (t *TabBook) createWidget() {
	t.gridLayout = NewGridLayout(
		GridLayoutOpts.Columns(1),
		GridLayoutOpts.Stretch([]bool{true}, []bool{false, true}))

	t.container = NewContainer(append(t.containerOpts, []ContainerOpt{ContainerOpts.Layout(t.gridLayout)}...)...)
	t.containerOpts = nil
}

func (t *TabBook) initTabBook() {
	t.init.Do()
	t.container.RemoveChildren()

	t.gridLayout.rowSpacing = *t.computedParams.ContentSpacing
	t.btnContainer = NewContainer(
		ContainerOpts.Layout(NewRowLayout(
			RowLayoutOpts.Spacing(*t.computedParams.TabSpacing))))
	t.container.AddChild(t.btnContainer)

	btnElements := []RadioGroupElement{}
	var currentTab *TabBookTab = t.tab

	for i := range t.tabs {

		t.tabs[i].GetWidget().parent = t.GetWidget()
		t.tabs[i].Validate()
		btnOpts := []ButtonOpt{
			ButtonOpts.Image(t.computedParams.TabButton.Image),
		}
		if t.tabs[i].image != nil {
			btnOpts = append(btnOpts, ButtonOpts.TextAndImage(t.tabs[i].label, t.computedParams.TabButton.TextFace, t.tabs[i].image, t.computedParams.TabButton.TextColor))
		} else {
			btnOpts = append(btnOpts, ButtonOpts.Text(t.tabs[i].label, t.computedParams.TabButton.TextFace, t.computedParams.TabButton.TextColor))
		}
		if t.computedParams.TabButton.MinSize != nil {
			btnOpts = append(btnOpts, ButtonOpts.WidgetOpts(WidgetOpts.MinSize(t.computedParams.TabButton.MinSize.X, t.computedParams.TabButton.MinSize.Y)))
		}
		if t.computedParams.TabButton.TextPadding != nil {
			btnOpts = append(btnOpts, ButtonOpts.TextPadding(t.computedParams.TabButton.TextPadding))
		}
		btn := NewButton(append(btnOpts, ButtonOpts.WidgetOpts(WidgetOpts.CustomData(t.tabs[i])))...)
		btnElements = append(btnElements, btn)
		t.btnContainer.AddChild(btn)
		t.tabToButton[t.tabs[i]] = btn
		if currentTab == nil {
			if t.initialTab == nil && !t.tabs[i].Disabled {
				currentTab = t.tabs[i]
			} else if t.initialTab == t.tabs[i] && !t.tabs[i].Disabled {
				currentTab = t.initialTab
			}
		}
	}
	// If we cannot find an initial tab default to to the first one
	if currentTab == nil {
		currentTab = t.tabs[0]
	}

	NewRadioGroup(
		RadioGroupOpts.Elements(btnElements...),
		RadioGroupOpts.InitialElement(t.tabToButton[currentTab]),
		RadioGroupOpts.ChangedHandler(func(args *RadioGroupChangedEventArgs) {
			if hasWidget, ok := args.Active.(HasWidget); ok {
				if tab, ok := hasWidget.GetWidget().CustomData.(*TabBookTab); ok {
					t.SetTab(tab)
				}
			}
		}))

	t.flipBook = NewFlipBook(
		FlipBookOpts.ContainerOpts(ContainerOpts.AutoDisableChildren()),
		FlipBookOpts.Padding(t.computedParams.ContentPadding),
	)
	t.container.AddChild(t.flipBook)

	t.tab = nil
	t.SetTab(currentTab)
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
		btn := t.GetTabButton(tab)
		if btn != nil {
			previousTab := t.tab

			tab.widget.parent = t.GetWidget()
			t.tab = tab
			t.flipBook.SetPage(tab)

			btn.SetState(WidgetChecked)

			t.TabSelectedEvent.Fire(&TabBookTabSelectedEventArgs{
				TabBook:     t,
				Tab:         tab,
				PreviousTab: previousTab,
			})
		}
	}
}

// Return the currently selected tab.
func (t *TabBook) Tab() *TabBookTab {
	return t.tab
}

// Return the button associated with the provided TabBookTab if not exists else nil.
func (t *TabBook) GetTabButton(tab *TabBookTab) *Button {
	t.init.Do()

	return t.tabToButton[tab]
}
