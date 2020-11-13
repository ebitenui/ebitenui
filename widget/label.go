package widget

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

type Label struct {
	Label string

	textOpts []TextOpt
	face     font.Face
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

	return l
}

func (o LabelOptions) TextOpts(opts ...TextOpt) LabelOpt {
	return func(l *Label) {
		l.textOpts = append(l.textOpts, opts...)
	}
}

func (o LabelOptions) Text(label string, face font.Face, color *LabelColor) LabelOpt {
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

func (l *Label) Render(screen *ebiten.Image, def DeferredRenderFunc) {
	l.init.Do()

	l.text.Label = l.Label

	if l.text.GetWidget().Disabled {
		l.text.Color = l.color.Disabled
	} else {
		l.text.Color = l.color.Idle
	}

	l.text.Render(screen, def)
}

func (l *Label) createWidget() {
	l.text = NewText(append(l.textOpts, TextOpts.Text(l.Label, l.face, l.color.Idle))...)
	l.textOpts = nil
	l.face = nil
}
