package widget

import (
	"image"

	"github.com/blizzy78/ebitenui/input"
	"github.com/hajimehoshi/ebiten"
)

type LabeledCheckbox struct {
	widgetOpts   []WidgetOpt
	checkboxOpts []CheckboxOpt
	textOpts     []TextOpt

	init      *MultiOnce
	container *Container
	checkbox  *Checkbox
	text      *Text
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

func (o labeledCheckboxOpts) WithCheckboxOpt(opt CheckboxOpt) LabeledCheckboxOpt {
	return func(l *LabeledCheckbox) {
		l.checkboxOpts = append(l.checkboxOpts, opt)
	}
}

func (o labeledCheckboxOpts) WithTextOpt(opt TextOpt) LabeledCheckboxOpt {
	return func(l *LabeledCheckbox) {
		l.textOpts = append(l.textOpts, opt)
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
		CheckboxOpts.WithButtonOpt(ButtonOpts.WithWidgetOpt(WidgetOpts.WithLayoutData(&RowLayoutData{
			Position: RowLayoutPositionCenter,
		}))),
	}...)...)
	l.container.AddChild(l.checkbox)
	l.checkboxOpts = nil

	// TODO: this should really be a Label instead of a Text so that it can be disabled
	l.text = NewText(append(l.textOpts, []TextOpt{
		TextOpts.WithWidgetOpt(WidgetOpts.WithLayoutData(&RowLayoutData{
			Position: RowLayoutPositionCenter,
		})),
	}...)...)
	l.container.AddChild(l.text)
	l.textOpts = nil
}
