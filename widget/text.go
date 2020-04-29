package widget

import (
	"bufio"
	"image"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
)

type Text struct {
	Label string
	Face  font.Face
	Color color.Color

	widget                    *Widget
	lastLabelForPreferredSize string
	lastFaceForPreferredSize  font.Face
	preferredWidth            int
	preferredHeight           int
}

type TextOpt func(t *Text)

const TextOpts = textOpts(true)

type textOpts bool

func NewText(opts ...TextOpt) *Text {
	t := &Text{
		widget: NewWidget(),
	}

	for _, o := range opts {
		o(t)
	}

	return t
}

func WithTextLayoutData(ld interface{}) TextOpt {
	return func(t *Text) {
		WidgetOpts.WithLayoutData(ld)(t.widget)
	}
}

func (o textOpts) WithText(label string, face font.Face, color color.Color) TextOpt {
	return func(t *Text) {
		t.Label = label
		t.Face = face
		t.Color = color
	}
}

func (t *Text) GetWidget() *Widget {
	return t.widget
}

func (t *Text) SetLocation(rect image.Rectangle) {
	t.widget.Rect = rect
}

func (t *Text) PreferredSize() (int, int) {
	if t.Label == t.lastLabelForPreferredSize && t.Face == t.lastFaceForPreferredSize {
		return t.preferredWidth, t.preferredHeight
	}

	lh := t.Face.Metrics().Height.Round()

	lines := 0
	w := 0
	s := bufio.NewScanner(strings.NewReader(t.Label))
	for s.Scan() {
		lines++

		lw := font.MeasureString(t.Face, s.Text()).Round()
		if lw > w {
			w = lw
		}
	}

	t.preferredWidth, t.preferredHeight = w, lh*lines
	t.lastLabelForPreferredSize = t.Label
	t.lastFaceForPreferredSize = t.Face

	return t.preferredWidth, t.preferredHeight
}

func (t *Text) Render(screen *ebiten.Image, def DeferredRenderFunc) {
	t.widget.Render(screen, def)
	t.draw(screen)
}

func (t *Text) draw(screen *ebiten.Image) {
	w, h := t.PreferredSize()

	r := t.widget.Rect

	// TODO: add alignment options
	x := (r.Dx()-w)/2 + r.Min.X - 1
	y := (r.Dy()-h)/2 + r.Min.Y - 1

	text.Draw(screen, t.Label, t.Face, x, y+t.Face.Metrics().Ascent.Round(), t.Color)
}
