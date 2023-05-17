package widget

import (
	"image"

	"github.com/ebitenui/ebitenui/input"
	"github.com/hajimehoshi/ebiten/v2"
)

type LabelOrder int

const (
	CHECKBOX_FIRST LabelOrder = iota
	LABEL_FIRST
)

type LabeledCheckbox struct {
	checkboxOpts []CheckboxOpt
	labelOpts    []LabelOpt
	spacing      int

	init       *MultiOnce
	widgetOpts []WidgetOpt
	container  *Container
	checkbox   *Checkbox
	label      *Label
	order      LabelOrder
}

type LabeledCheckboxOpt func(l *LabeledCheckbox)

type LabeledCheckboxOptions struct {
}

var LabeledCheckboxOpts LabeledCheckboxOptions

func NewLabeledCheckbox(opts ...LabeledCheckboxOpt) *LabeledCheckbox {
	l := &LabeledCheckbox{
		spacing: 8,
		order:   CHECKBOX_FIRST,
		init:    &MultiOnce{},
	}

	l.init.Append(l.createWidget)

	for _, o := range opts {
		o(l)
	}

	return l
}

func (o LabeledCheckboxOptions) WidgetOpts(opts ...WidgetOpt) LabeledCheckboxOpt {
	return func(l *LabeledCheckbox) {
		l.widgetOpts = append(l.widgetOpts, opts...)
	}
}

func (o LabeledCheckboxOptions) CheckboxOpts(opts ...CheckboxOpt) LabeledCheckboxOpt {
	return func(l *LabeledCheckbox) {
		l.checkboxOpts = append(l.checkboxOpts, opts...)
	}
}

func (o LabeledCheckboxOptions) LabelOpts(opts ...LabelOpt) LabeledCheckboxOpt {
	return func(l *LabeledCheckbox) {
		l.labelOpts = append(l.labelOpts, opts...)
	}
}

func (o LabeledCheckboxOptions) Spacing(s int) LabeledCheckboxOpt {
	return func(l *LabeledCheckbox) {
		l.spacing = s
	}
}

func (o LabeledCheckboxOptions) LabelFirst() LabeledCheckboxOpt {
	return func(l *LabeledCheckbox) {
		l.order = LABEL_FIRST
	}
}

func (l *LabeledCheckbox) SetState(state WidgetState) {
	l.init.Do()
	l.checkbox.SetState(state)
}

func (l *LabeledCheckbox) GetWidget() *Widget {
	l.init.Do()
	return l.container.GetWidget()
}

func (l *LabeledCheckbox) PreferredSize() (int, int) {
	l.init.Do()
	return l.container.PreferredSize()
}

func (l *LabeledCheckbox) SetLocation(rect image.Rectangle) {
	l.init.Do()
	l.container.SetLocation(rect)
}

func (l *LabeledCheckbox) SetupInputLayer(def input.DeferredSetupInputLayerFunc) {
	l.init.Do()
	l.checkbox.SetupInputLayer(def)
}

func (l *LabeledCheckbox) Render(screen *ebiten.Image, def DeferredRenderFunc) {
	l.init.Do()
	l.container.Render(screen, def)
}

func (l *LabeledCheckbox) Checkbox() *Checkbox {
	l.init.Do()
	return l.checkbox
}

func (l *LabeledCheckbox) Label() *Label {
	l.init.Do()
	return l.label
}
func (l *LabeledCheckbox) Focus(focused bool) {
	l.init.Do()
	l.GetWidget().FireFocusEvent(l, focused, image.Point{-1, -1})
	l.checkbox.button.focused = focused
}

func (l *LabeledCheckbox) IsFocused() bool {
	return l.checkbox.button.focused
}

func (l *LabeledCheckbox) TabOrder() int {
	l.init.Do()
	return l.checkbox.tabOrder
}
func (l *LabeledCheckbox) createWidget() {
	l.container = NewContainer(
		ContainerOpts.Layout(NewRowLayout(RowLayoutOpts.Spacing(l.spacing))),
		ContainerOpts.AutoDisableChildren(),
		ContainerOpts.WidgetOpts(l.widgetOpts...),
	)

	l.checkbox = NewCheckbox(append(l.checkboxOpts, CheckboxOpts.ButtonOpts(ButtonOpts.WidgetOpts(WidgetOpts.LayoutData(RowLayoutData{
		Position: RowLayoutPositionCenter,
	}))))...)
	l.checkboxOpts = nil

	l.label = NewLabel(append(l.labelOpts, LabelOpts.TextOpts(TextOpts.WidgetOpts(
		WidgetOpts.LayoutData(RowLayoutData{
			Position: RowLayoutPositionCenter,
		}),

		WidgetOpts.MouseButtonReleasedHandler(func(args *WidgetMouseButtonReleasedEventArgs) {
			if !args.Widget.Disabled && args.Button == ebiten.MouseButtonLeft && args.Inside {
				l.checkbox.SetState(l.checkbox.state.Advance(l.checkbox.triState))
			}
		}),
	)))...)

	if l.order == CHECKBOX_FIRST {
		l.container.AddChild(l.checkbox)
		l.container.AddChild(l.label)
	} else {
		l.container.AddChild(l.label)
		l.container.AddChild(l.checkbox)
	}
	l.labelOpts = nil
}
