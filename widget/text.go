package widget

import (
	"bufio"
	"image"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type Text struct {
	Label string
	Face  font.Face
	Color color.Color

	widgetOpts []WidgetOpt
	position   TextPosition

	init                      *MultiOnce
	widget                    *Widget
	lastLabelForPreferredSize string
	lastFaceForPreferredSize  font.Face
	preferredWidth            int
	preferredHeight           int
}

type TextOpt func(t *Text)

type TextPosition int

const (
	TextPositionStart = TextPosition(iota)
	TextPositionCenter
	TextPositionEnd
)

const TextOpts = textOpts(true)

type textOpts bool

func NewText(opts ...TextOpt) *Text {
	t := &Text{
		init: &MultiOnce{},
	}

	t.init.Append(t.createWidget)

	for _, o := range opts {
		o(t)
	}

	return t
}

func (o textOpts) WithWidgetOpt(opt WidgetOpt) TextOpt {
	return func(t *Text) {
		t.widgetOpts = append(t.widgetOpts, opt)
	}
}

func (o textOpts) WithText(label string, face font.Face, color color.Color) TextOpt {
	return func(t *Text) {
		t.Label = label
		t.Face = face
		t.Color = color
	}
}

func (o textOpts) WithPosition(p TextPosition) TextOpt {
	return func(t *Text) {
		t.position = p
	}
}

func (t *Text) GetWidget() *Widget {
	t.init.Do()
	return t.widget
}

func (t *Text) SetLocation(rect image.Rectangle) {
	t.init.Do()
	t.widget.Rect = rect
}

func (t *Text) PreferredSize() (int, int) {
	t.init.Do()

	if t.Label == t.lastLabelForPreferredSize && t.Face == t.lastFaceForPreferredSize {
		return t.preferredWidth, t.preferredHeight
	}

	m := t.Face.Metrics()
	fh := m.Ascent + m.Descent
	lh := m.Height
	ld := lh - fh

	lines := 0
	w := 0
	s := bufio.NewScanner(strings.NewReader(t.Label))
	for s.Scan() {
		lines++

		lw := font.MeasureString(t.Face, s.Text()).Ceil()
		if lw > w {
			w = lw
		}
	}

	t.preferredWidth, t.preferredHeight = w, (fixed.I(lines).Mul(lh) - ld).Ceil()

	t.lastLabelForPreferredSize = t.Label
	t.lastFaceForPreferredSize = t.Face

	return t.preferredWidth, t.preferredHeight
}

func (t *Text) Render(screen *ebiten.Image, def DeferredRenderFunc) {
	t.init.Do()
	t.widget.Render(screen, def)
	t.draw(screen)
}

func (t *Text) draw(screen *ebiten.Image) {
	w, h := t.PreferredSize()

	r := t.widget.Rect
	p := r.Min

	x := p.X
	switch t.position {
	case TextPositionCenter:
		x += (r.Dx() - w) / 2
	case TextPositionEnd:
		x += r.Dx() - w
	}
	y := p.Y + (r.Dy()-h)/2 + t.Face.Metrics().Ascent.Round()

	text.Draw(screen, t.Label, t.Face, x, y, t.Color)
}

func (t *Text) createWidget() {
	t.widget = NewWidget(t.widgetOpts...)
	t.widgetOpts = nil
}
