package widget

import (
	"bufio"
	"image"
	"image/color"
	"math"
	"regexp"
	"strings"

	"github.com/ebitenui/ebitenui/utilities/colorutil"
	"github.com/ebitenui/ebitenui/utilities/datastructures"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

const bbcodeRegEx = `\[color=[0-9a-fA-F]{6}\]|\[\/color\]`
const COLOR_OPEN = "color="
const COLOR_CLOSE = "/color]"

type Text struct {
	Label              string
	Face               font.Face
	Color              color.Color
	MaxWidth           float64
	Inset              Insets
	widgetOpts         []WidgetOpt
	horizontalPosition TextPosition
	verticalPosition   TextPosition

	init          *MultiOnce
	widget        *Widget
	measurements  textMeasurements
	bbcodeRegex   *regexp.Regexp
	processBBCode bool
	colorList     *datastructures.Stack[color.Color]
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

	lines             [][]string
	lineWidths        []float64
	lineHeight        float64
	ascent            float64
	boundingBoxWidth  float64
	boundingBoxHeight float64
}

type bbCodeText struct {
	text  string
	color color.Color
}

var TextOpts TextOptions

func NewText(opts ...TextOpt) *Text {
	t := &Text{
		init: &MultiOnce{},
	}
	t.bbcodeRegex, _ = regexp.Compile(bbcodeRegEx)

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

// Text combines three options: TextLabel, TextFace and TextColor.
// It can be used for the inline configurations of Text object while
// separate functions are useful for a multi-step configuration.
func (o TextOptions) Text(label string, face font.Face, color color.Color) TextOpt {
	return func(t *Text) {
		t.Label = label
		t.Face = face
		t.Color = color
	}
}

func (o TextOptions) TextLabel(label string) TextOpt {
	return func(t *Text) {
		t.Label = label
	}
}

func (o TextOptions) TextFace(face font.Face) TextOpt {
	return func(t *Text) {
		t.Face = face
	}
}

func (o TextOptions) TextColor(color color.Color) TextOpt {
	return func(t *Text) {
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

func (o TextOptions) ProcessBBCode(processBBCode bool) TextOpt {
	return func(t *Text) {
		t.processBBCode = processBBCode
	}
}

// MaxWidth sets the max width the text will allow before wrapping to the next line
func (o TextOptions) MaxWidth(maxWidth float64) TextOpt {
	return func(t *Text) {
		t.MaxWidth = maxWidth
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
		p = p.Add(image.Point{0, int(float64(r.Dy()) - t.measurements.boundingBoxHeight)})
	}

	t.colorList = &datastructures.Stack[color.Color]{}
	t.colorList.Push(&t.Color)

	for i, line := range t.measurements.lines {
		ly := int(math.Round(float64(p.Y) + t.measurements.lineHeight*float64(i) + t.measurements.ascent))
		if ly > screen.Bounds().Max.Y+int(math.Round(t.measurements.ascent)) {
			return
		}
		if ly < -int(math.Round(t.measurements.lineHeight-t.measurements.ascent)) {
			continue
		}
		if t.widget.parent != nil {
			if ly < t.widget.parent.Rect.Min.Y {
				continue
			}
			if ly-int(math.Round(t.measurements.lineHeight)) > t.widget.parent.Rect.Max.Y {
				return
			}
		}

		lx := p.X
		switch t.horizontalPosition {
		case TextPositionCenter:
			lx += int(math.Round((float64(w) - t.measurements.lineWidths[i]) / 2))
		case TextPositionEnd:
			lx += int(math.Ceil(float64(w)-t.measurements.lineWidths[i])) - t.Inset.Right
		default:
			lx += t.Inset.Left
		}

		if t.processBBCode {
			spaceWidth := font.MeasureString(t.Face, " ").Round()
			for _, word := range line {
				pieces, updatedColor := t.handleBBCodeColor(word)
				for _, piece := range pieces {
					text.Draw(screen, piece.text, t.Face, lx, ly, piece.color)
					wordWidth := font.MeasureString(t.Face, piece.text)
					lx += wordWidth.Round()
				}
				text.Draw(screen, " ", t.Face, lx, ly, updatedColor)
				lx += spaceWidth
			}
		} else {
			text.Draw(screen, strings.Join(line, " "), t.Face, lx, ly, t.Color)
		}
	}
}

func (t *Text) handleBBCodeColor(word string) ([]bbCodeText, color.Color) {
	var result []bbCodeText
	tags := t.bbcodeRegex.FindAllStringIndex(word, -1)
	var newColor = *t.colorList.Top()
	if len(tags) > 0 {
		resultStr := ""
		isTag := false
		for idx := range word {
			if len(tags) > 0 {
				if tags[0][0] > idx || (isTag && idx < tags[0][1]) {
					resultStr = resultStr + string(word[idx])
				} else if tags[0][1] == idx {
					if strings.HasPrefix(resultStr, COLOR_OPEN) {
						c, err := colorutil.HexToColor(resultStr[6:12])
						if err == nil {
							t.colorList.Push(&c)
							newColor = c
						}
					} else if resultStr == COLOR_CLOSE {
						if t.colorList.Size() > 1 {
							t.colorList.Pop()
						}
						newColor = *t.colorList.Top()
					}
					tags = tags[1:]
					if len(tags) > 0 && tags[0][0] == idx {
						resultStr = ""
						isTag = true
					} else {
						resultStr = string(word[idx])
						isTag = false
					}
				} else {
					result = append(result, bbCodeText{text: resultStr, color: newColor})
					resultStr = ""
					isTag = true
				}
			} else {
				resultStr = resultStr + string(word[idx])
			}
		}
		if len(resultStr) > 0 {
			if isTag {
				if strings.HasPrefix(resultStr, COLOR_OPEN) {
					c, err := colorutil.HexToColor(resultStr[6:12])
					if err == nil {
						t.colorList.Push(&c)
						newColor = c
					}
				} else if resultStr == COLOR_CLOSE {
					if t.colorList.Size() > 1 {
						t.colorList.Pop()
					}
					newColor = *t.colorList.Top()
				}
			} else {
				result = append(result, bbCodeText{text: resultStr, color: newColor})
			}
		}
	} else {
		result = append(result, bbCodeText{text: word, color: newColor})
	}

	return result, newColor
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
		if t.MaxWidth > 0 {
			var newLine []string
			newLineWidth := float64(t.Inset.Left + t.Inset.Right)
			words := strings.Split(s.Text(), " ")
			for _, word := range words {
				wordWidth := fixedInt26_6ToFloat64(font.MeasureString(t.Face, word+" "))

				// Strip out any bbcodes from size calculation
				if t.processBBCode {
					if t.bbcodeRegex.MatchString(word) {
						cleaned := t.bbcodeRegex.ReplaceAllString(word, "")
						wordWidth = fixedInt26_6ToFloat64(font.MeasureString(t.Face, cleaned+" "))
					}
				}

				// If the new word doesn't push this past the max width continue adding to the current line
				if newLineWidth+wordWidth < t.MaxWidth {
					newLine = append(newLine, word)
					newLineWidth += wordWidth
				} else {
					// If the new word would push this past the max width save off the current line and start a new one
					if len(newLine) != 0 {
						t.measurements.lines = append(t.measurements.lines, newLine)
						t.measurements.lineWidths = append(t.measurements.lineWidths, newLineWidth)

						if newLineWidth > t.measurements.boundingBoxWidth {
							t.measurements.boundingBoxWidth = newLineWidth
						}
					}
					newLine = []string{word}
					newLineWidth = wordWidth + float64(t.Inset.Left+t.Inset.Right)
				}
			}
			// Save the final line
			if len(newLine) != 0 {
				t.measurements.lines = append(t.measurements.lines, newLine)
				t.measurements.lineWidths = append(t.measurements.lineWidths, newLineWidth)

				if newLineWidth > t.measurements.boundingBoxWidth {
					t.measurements.boundingBoxWidth = newLineWidth
				}
			}
		} else {
			line := s.Text()
			t.measurements.lines = append(t.measurements.lines, []string{line})
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
