package widget

import (
	"image"

	"github.com/blizzy78/ebitenui/input"
	"github.com/hajimehoshi/ebiten"
)

type LabeledCheckbox struct {
	checkboxOpts []CheckboxOpt
	labelOpts    []LabelOpt

	init      *MultiOnce
	container *Container
	checkbox  *Checkbox
	label     *Label
}

type LabeledCheckboxOpt func(l *LabeledCheckbox)

const LabeledCheckboxOpts = labeledCheckboxOpts(true)

type labeledCheckboxOpts bool

func NewLabeledCheckbox(opts ...LabeledCheckboxOpt) *LabeledCheckbox {
	l := &LabeledCheckbox{
		init: &MultiOnce{},
	}

	l.init.Append(l.createWidget)

	for _, o := range opts {
		o(l)
	}

	return l
}

func (o labeledCheckboxOpts) WithCheckboxOpts(opts ...CheckboxOpt) LabeledCheckboxOpt {
	return func(l *LabeledCheckbox) {
		l.checkboxOpts = append(l.checkboxOpts, opts...)
	}
}

func (o labeledCheckboxOpts) WithLabelOpts(opts ...LabelOpt) LabeledCheckboxOpt {
	return func(l *LabeledCheckbox) {
		l.labelOpts = append(l.labelOpts, opts...)
	}
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

func (l *LabeledCheckbox) createWidget() {
	l.container = NewContainer(
		ContainerOpts.WithLayout(NewRowLayout(
			RowLayoutOpts.WithSpacing(10))),
		ContainerOpts.WithAutoDisableChildren())

	l.checkbox = NewCheckbox(append(l.checkboxOpts, []CheckboxOpt{
		CheckboxOpts.WithButtonOpts(ButtonOpts.WithWidgetOpts(WidgetOpts.WithLayoutData(&RowLayoutData{
			Position: RowLayoutPositionCenter,
		}))),
	}...)...)
	l.container.AddChild(l.checkbox)
	l.checkboxOpts = nil

	l.label = NewLabel(append(l.labelOpts, []LabelOpt{
		LabelOpts.WithTextOpts(
			TextOpts.WithWidgetOpts(
				WidgetOpts.WithLayoutData(&RowLayoutData{
					Position: RowLayoutPositionCenter,
				}),

				WidgetOpts.WithMouseButtonReleasedHandler(func(args *WidgetMouseButtonReleasedEventArgs) {
					if !args.Widget.Disabled && args.Button == ebiten.MouseButtonLeft && args.Inside {
						l.checkbox.SetState(l.checkbox.state.Advance(l.checkbox.triState))
					}
				}),
			),
		),
	}...)...)
	l.container.AddChild(l.label)
	l.labelOpts = nil
}
