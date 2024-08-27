package widget

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Label struct {
	Label string

	textOpts []TextOpt
	face     text.Face
	color    *LabelColor

	init *MultiOnce
	text *Text
}

type LabelOpt func(l *Label)

type LabelColor struct {
	Idle     color.Color
	Disabled color.Color
}

type LabelOptions struct {
}

var LabelOpts LabelOptions

func NewLabel(opts ...LabelOpt) *Label {
	l := &Label{
		init: &MultiOnce{},
	}

	l.init.Append(l.createWidget)

	for _, o := range opts {
		o(l)
	}

	l.validate()

	return l
}

func (l *Label) validate() {
	if l.color == nil {
		panic("Label: LabelColor is required.")
	}
	if l.color.Idle == nil {
		panic("Label: LabelColor.Idle is required.")
	}
	if l.face == nil {
		panic("Label: Font Face is required.")
	}
}

func (o LabelOptions) TextOpts(opts ...TextOpt) LabelOpt {
	return func(l *Label) {
		l.textOpts = append(l.textOpts, opts...)
	}
}

func (o LabelOptions) Text(label string, face text.Face, color *LabelColor) LabelOpt {
	return func(l *Label) {
		l.Label = label
		l.face = face
		l.color = color
	}
}

func (l *Label) GetWidget() *Widget {
	l.init.Do()
	return l.text.GetWidget()
}

func (l *Label) SetLocation(rect image.Rectangle) {
	l.init.Do()
	l.text.SetLocation(rect)
}

func (l *Label) PreferredSize() (int, int) {
	l.init.Do()
	return l.text.PreferredSize()
}

func (l *Label) Render(screen *ebiten.Image) {
	l.init.Do()

	l.text.Label = l.Label

	if l.text.GetWidget().Disabled && l.color.Disabled != nil {
		l.text.Color = l.color.Disabled
	} else {
		l.text.Color = l.color.Idle
	}

	l.text.Render(screen)
}

func (l *Label) Update() {
	l.init.Do()

	l.text.Update()
}

func (l *Label) createWidget() {
	l.text = NewText(append(l.textOpts, TextOpts.Text(l.Label, l.face, l.color.Idle))...)
	l.textOpts = nil
	l.face = nil
}
