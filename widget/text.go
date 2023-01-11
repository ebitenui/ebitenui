package widget

import (
	"bufio"
	"image"
	"image/color"
	"math"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type Text struct {
	Label              string
	Face               font.Face
	Color              color.Color
	MaxWidth           float64
	Inset              Insets
	widgetOpts         []WidgetOpt
	horizontalPosition TextPosition
	verticalPosition   TextPosition

	init         *MultiOnce
	widget       *Widget
	measurements textMeasurements
}

type TextOpt func(t *Text)

type TextPosition int

const (
	TextPositionStart = TextPosition(iota)
	TextPositionCenter
	TextPositionEnd
)

type TextOptions struct {
}

type textMeasurements struct {
	label    string
	face     font.Face
	maxWidth float64

	lines             []string
	lineWidths        []float64
	lineHeight        float64
	ascent            float64
	boundingBoxWidth  float64
	boundingBoxHeight float64
}

var TextOpts TextOptions

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

func (o TextOptions) WidgetOpts(opts ...WidgetOpt) TextOpt {
	return func(t *Text) {
		t.widgetOpts = append(t.widgetOpts, opts...)
	}
}

func (o TextOptions) Text(label string, face font.Face, color color.Color) TextOpt {
	return func(t *Text) {
		t.Label = label
		t.Face = face
		t.Color = color
	}
}

func (o TextOptions) Insets(inset Insets) TextOpt {
	return func(t *Text) {
		t.Inset = inset
	}
}

func (o TextOptions) Position(h TextPosition, v TextPosition) TextOpt {
	return func(t *Text) {
		t.horizontalPosition = h
		t.verticalPosition = v
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
	t.measure()
	w := int(math.Ceil(t.measurements.boundingBoxWidth))
	h := int(math.Ceil(t.measurements.boundingBoxHeight))

	if t.widget != nil && h < t.widget.MinHeight {
		h = t.widget.MinHeight
	}
	if t.widget != nil && w < t.widget.MinWidth {
		w = t.widget.MinWidth
	}
	return w, h
}

func (t *Text) Render(screen *ebiten.Image, def DeferredRenderFunc) {
	t.init.Do()
	t.widget.Render(screen, def)
	t.draw(screen)
}

func (t *Text) draw(screen *ebiten.Image) {
	t.measure()

	r := t.widget.Rect
	w := r.Dx()
	p := r.Min

	switch t.verticalPosition {
	case TextPositionCenter:
		p = p.Add(image.Point{0, int((float64(r.Dy()) - t.measurements.boundingBoxHeight) / 2)})
	case TextPositionEnd:
		p = p.Add(image.Point{0, int((float64(r.Dy()) - t.measurements.boundingBoxHeight))})
	}

	for i, line := range t.measurements.lines {
		lx := p.X
		switch t.horizontalPosition {
		case TextPositionCenter:
			lx += int(math.Round((float64(w) - t.measurements.lineWidths[i]) / 2))
		case TextPositionEnd:
			lx += int(math.Ceil(float64(w)-t.measurements.lineWidths[i])) - t.Inset.Right
		default:
			lx += t.Inset.Left
		}

		ly := int(math.Round(float64(p.Y) + t.measurements.lineHeight*float64(i) + t.measurements.ascent))

		text.Draw(screen, line, t.Face, lx, ly, t.Color)
	}
}

func (t *Text) measure() {
	if t.Label == t.measurements.label && t.Face == t.measurements.face && t.MaxWidth == t.measurements.maxWidth {
		return
	}

	m := t.Face.Metrics()

	t.measurements = textMeasurements{
		label:    t.Label,
		face:     t.Face,
		ascent:   fixedInt26_6ToFloat64(m.Ascent),
		maxWidth: t.MaxWidth,
	}

	fh := fixedInt26_6ToFloat64(m.Ascent + m.Descent)
	t.measurements.lineHeight = fixedInt26_6ToFloat64(m.Height)
	ld := t.measurements.lineHeight - fh

	s := bufio.NewScanner(strings.NewReader(t.Label))
	for s.Scan() {
		line := s.Text()
		var newLine string = ""
		var newLineWidth float64 = float64(t.Inset.Left + t.Inset.Right)
		words := strings.Split(line, " ")
		if t.MaxWidth > 0 {
			for _, word := range words {
				word = word + " "
				wordWidth := fixedInt26_6ToFloat64(font.MeasureString(t.Face, word))

				if newLineWidth+wordWidth < t.MaxWidth {
					newLine = newLine + word
					newLineWidth += wordWidth
				} else {
					if newLine != "" {
						t.measurements.lines = append(t.measurements.lines, newLine)
						t.measurements.lineWidths = append(t.measurements.lineWidths, newLineWidth)

						if newLineWidth > t.measurements.boundingBoxWidth {
							t.measurements.boundingBoxWidth = newLineWidth
						}
					}
					newLine = word
					newLineWidth = wordWidth + float64(t.Inset.Left+t.Inset.Right)
				}
			}
			if newLine != "" {
				t.measurements.lines = append(t.measurements.lines, newLine)
				t.measurements.lineWidths = append(t.measurements.lineWidths, newLineWidth)

				if newLineWidth > t.measurements.boundingBoxWidth {
					t.measurements.boundingBoxWidth = newLineWidth
				}
			}
		} else {
			line := s.Text()
			t.measurements.lines = append(t.measurements.lines, line)

			lw := fixedInt26_6ToFloat64(font.MeasureString(t.Face, line)) + float64(t.Inset.Left+t.Inset.Right)
			t.measurements.lineWidths = append(t.measurements.lineWidths, lw)

			if lw > t.measurements.boundingBoxWidth {
				t.measurements.boundingBoxWidth = lw
			}
		}

	}

	t.measurements.boundingBoxHeight = float64(len(t.measurements.lines))*t.measurements.lineHeight - ld
}

func (t *Text) createWidget() {
	t.widget = NewWidget(t.widgetOpts...)
	t.widgetOpts = nil
}

func fixedInt26_6ToFloat64(i fixed.Int26_6) float64 {
	return float64(i) / (1 << 6)
}
