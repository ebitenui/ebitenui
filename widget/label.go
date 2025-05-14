package widget

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type LabelParams struct {
	Face    *text.Face
	Color   *LabelColor
	Padding *Insets
}

type Label struct {
	Label string

	definedParams  LabelParams
	computedParams LabelParams

	textOpts []TextOpt
	init     *MultiOnce
	text     *Text
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

func (l *Label) Validate() {
	l.init.Do()
	l.populateComputedParams()

	if l.computedParams.Color == nil {
		panic("Label: LabelColor is required.")
	}
	if l.computedParams.Color.Idle == nil {
		panic("Label: LabelColor.Idle is required.")
	}
	if l.computedParams.Face == nil {
		panic("Label: Font Face is required.")
	}
	l.text.Validate()
}

func (l *Label) populateComputedParams() {
	lblParams := LabelParams{}

	theme := l.text.GetWidget().GetTheme()

	if theme != nil {
		if theme.LabelTheme != nil {
			lblParams.Face = theme.LabelTheme.Face
			lblParams.Color = theme.LabelTheme.Color
			lblParams.Padding = theme.LabelTheme.Padding
		}
	}
	if l.definedParams.Color != nil {
		if lblParams.Color == nil {
			lblParams.Color = l.definedParams.Color
		} else {
			if l.definedParams.Color.Idle != nil {
				lblParams.Color.Idle = l.definedParams.Color.Idle
			}
			if l.definedParams.Color.Disabled != nil {
				lblParams.Color.Disabled = l.definedParams.Color.Disabled
			}
		}
	}
	if l.definedParams.Face != nil {
		lblParams.Face = l.definedParams.Face
	}
	if l.definedParams.Padding != nil {
		lblParams.Padding = l.definedParams.Padding
	}

	l.computedParams = lblParams
	l.setComputedParams()
}

func (o LabelOptions) TextOpts(opts ...TextOpt) LabelOpt {
	return func(l *Label) {
		l.textOpts = append(l.textOpts, opts...)
	}
}

// Set the label text, font, and font colors.
func (o LabelOptions) Text(label string, face *text.Face, color *LabelColor) LabelOpt {
	return func(l *Label) {
		l.Label = label
		l.definedParams.Face = face
		l.definedParams.Color = color
	}
}

// Set the label text.
func (o LabelOptions) LabelText(label string) LabelOpt {
	return func(l *Label) {
		l.Label = label
	}
}

// Set the label font.
func (o LabelOptions) LabelFace(face *text.Face) LabelOpt {
	return func(l *Label) {
		l.definedParams.Face = face
	}
}

// Set the label font colors.
func (o LabelOptions) LabelColor(color *LabelColor) LabelOpt {
	return func(l *Label) {
		l.definedParams.Color = color
	}
}

// Set the label padding.
func (o LabelOptions) LabelPadding(padding *Insets) LabelOpt {
	return func(l *Label) {
		l.definedParams.Padding = padding
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

	if l.text.GetWidget().Disabled && l.computedParams.Color.Disabled != nil {
		l.text.SetColor(l.computedParams.Color.Disabled)
	} else {
		l.text.SetColor(l.computedParams.Color.Idle)
	}

	l.text.Render(screen)
}

func (l *Label) Update() {
	l.init.Do()
	l.text.Update()
}

func (l *Label) createWidget() {
	l.text = NewText(append(l.textOpts, TextOpts.TextLabel(l.Label))...)
}

func (l *Label) setComputedParams() {
	l.text.SetFace(l.computedParams.Face)
	if l.computedParams.Padding != nil {
		l.text.SetPadding(l.computedParams.Padding)
	}
	l.text.SetColor(l.computedParams.Color.Idle)
}
