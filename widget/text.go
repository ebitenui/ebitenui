package widget

import (
	"bufio"
	"image"
	"image/color"
	"math"
	"regexp"
	"strings"

	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/input"
	"github.com/ebitenui/ebitenui/utilities/colorutil"
	"github.com/ebitenui/ebitenui/utilities/datastructures"
	"github.com/frustra/bbcode"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var bbcodeRegex = regexp.MustCompile(`\[color=#[0-9a-fA-F]{6}]|\[\/color]|\[link]|\[\/link]|\[link=[^\]]*]`)

const COLOR_TAG = "color"
const LINK_TAG = "link"

type Text struct {
	Label         string
	Face          text.Face
	Color         color.Color
	MaxWidth      float64
	Inset         Insets
	Padding       Insets
	ProcessBBCode bool
	LinkColor     TextLinkColor

	widgetOpts         []WidgetOpt
	horizontalPosition TextPosition
	verticalPosition   TextPosition

	init         *MultiOnce
	widget       *Widget
	measurements textMeasurements
	colorList    *datastructures.Stack[color.Color]
	linkStack    *datastructures.Stack[linkData]
	currentLink  *bbCodeText

	LinkClickedEvent *event.Event
}

type textMeasurements struct {
	label         string
	face          text.Face
	maxWidth      float64
	ProcessBBCode bool

	processedLines      [][]*bbCodeText
	processedLineWidths []float64
	lineHeight          float64
	ascent              float64
	boundingBoxWidth    float64
	boundingBoxHeight   float64
}

type bbCodeText struct {
	text      string
	color     color.Color
	linkValue *linkData
	hovered   bool
}

type linkData struct {
	id         string
	text       string
	args       map[string]string
	textBlocks []*bbCodeText
}

type TextLinkColor struct {
	Idle  color.Color
	Hover color.Color
}

type LinkClickedEventArgs struct {
	Text    *Text
	Id      string
	Value   string
	Args    map[string]string
	OffsetX int
	OffsetY int
}
type LinkClickedHandlerFunc func(args *LinkClickedEventArgs)

type TextOpt func(t *Text)

type TextPosition int

const (
	TextPositionStart = TextPosition(iota)
	TextPositionCenter
	TextPositionEnd
)

type TextOptions struct {
}

var TextOpts TextOptions

func NewText(opts ...TextOpt) *Text {
	t := &Text{
		init:             &MultiOnce{},
		LinkClickedEvent: &event.Event{},
		colorList:        &datastructures.Stack[color.Color]{},
		linkStack:        &datastructures.Stack[linkData]{},
		LinkColor: TextLinkColor{
			Idle:  color.NRGBA{R: 0, G: 0, B: 250, A: 255},
			Hover: color.NRGBA{R: 255, G: 0, B: 250, A: 255},
		},
	}

	t.init.Append(t.createWidget)

	for _, o := range opts {
		o(t)
	}
	t.validate()

	t.colorList.Push(&t.Color)
	return t
}

func (t *Text) validate() {
	if t.Color == nil {
		panic("Text: Color is required.")
	}
	if t.Face == nil {
		panic("Text: Face is required.")
	}
}

func (o TextOptions) WidgetOpts(opts ...WidgetOpt) TextOpt {
	return func(t *Text) {
		t.widgetOpts = append(t.widgetOpts, opts...)
	}
}

// Text combines three options: TextLabel, TextFace and TextColor.
// It can be used for the inline configurations of Text object while
// separate functions are useful for a multi-step configuration.
func (o TextOptions) Text(label string, face text.Face, color color.Color) TextOpt {
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

func (o TextOptions) TextFace(face text.Face) TextOpt {
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
func (o TextOptions) Padding(padding Insets) TextOpt {
	return func(t *Text) {
		t.Padding = padding
	}
}
func (o TextOptions) Position(h TextPosition, v TextPosition) TextOpt {
	return func(t *Text) {
		t.horizontalPosition = h
		t.verticalPosition = v
	}
}

// This option tells the text object to process BBCodes.
//
// Currently the system supports the following BBCodes:
//   - color - [color=#FFFFFF] text [/color] - defines a color code for the enclosed text
//   - link - [link=id arg1:value1 ... argX:valueX] text [/link] - defines a clickable section of text,
//     that will trigger a callback.
func (o TextOptions) ProcessBBCode(processBBCode bool) TextOpt {
	return func(t *Text) {
		t.ProcessBBCode = processBBCode
	}
}

// This option sets the idle and hover color for text that is wrapped in a
// [link][/link] bbcode.
//
// Note: this is only used if ProcessBBCode is true.
func (o TextOptions) LinkColor(linkColor TextLinkColor) TextOpt {
	return func(t *Text) {
		t.LinkColor = linkColor
	}
}

// MaxWidth sets the max width the text will allow before wrapping to the next line.
func (o TextOptions) MaxWidth(maxWidth float64) TextOpt {
	return func(t *Text) {
		t.MaxWidth = maxWidth
	}
}

// Defines the handler to be called when a BBCode defined link is clicked.
//
// Note: this is only used if ProcessBBCode is true.
func (o TextOptions) LinkClickedHandler(f LinkClickedHandlerFunc) TextOpt {
	return func(b *Text) {
		b.LinkClickedEvent.AddHandler(func(args any) {
			if arg, ok := args.(*LinkClickedEventArgs); ok {
				f(arg)
			}
		})
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
	w := int(math.Ceil(t.measurements.boundingBoxWidth)) + t.Padding.Left + t.Padding.Right
	h := int(math.Ceil(t.measurements.boundingBoxHeight)) + t.Padding.Top + t.Padding.Bottom

	if t.widget != nil && h < t.widget.MinHeight {
		h = t.widget.MinHeight
	}
	if t.widget != nil && w < t.widget.MinWidth {
		w = t.widget.MinWidth
	}
	return w, h
}

func (t *Text) Render(screen *ebiten.Image) {
	t.init.Do()
	t.widget.Render(screen)
	t.draw(screen)
}

func (t *Text) Update() {
	t.init.Do()

	t.widget.Update()
	if t.currentLink != nil && input.MouseButtonJustPressed(ebiten.MouseButton0) {
		if t.LinkClickedEvent != nil {
			t.LinkClickedEvent.Fire(&LinkClickedEventArgs{
				Text:  t,
				Value: t.currentLink.linkValue.text,
				Id:    t.currentLink.linkValue.id,
				Args:  t.currentLink.linkValue.args,
			})
		}
	}
	t.currentLink = nil
}

func (t *Text) draw(screen *ebiten.Image) {
	t.measure()

	r := t.widget.Rect
	w := r.Dx()
	p := r.Min

	switch t.verticalPosition {
	case TextPositionStart:
		p = p.Add(image.Point{0, t.Inset.Top})
	case TextPositionCenter:
		p = p.Add(image.Point{0, int((float64(r.Dy())-t.measurements.boundingBoxHeight)/2 + float64(t.Inset.Top))})
	case TextPositionEnd:
		p = p.Add(image.Point{0, int(float64(r.Dy())-t.measurements.boundingBoxHeight) - t.Inset.Bottom})
	}

	if t.ProcessBBCode {
		// Reset Hovered
		for linesIdx := range t.measurements.processedLines {
			for idx := range t.measurements.processedLines[linesIdx] {
				t.measurements.processedLines[linesIdx][idx].hovered = false
			}
		}

		// Process Hovered
		cursorX, cursorY := input.CursorPosition()
		cursorPoint := image.Point{X: cursorX, Y: cursorY}
		drawnRectangle := t.widget.Rect
		if t.widget.parent != nil {
			drawnRectangle = t.widget.parent.Rect
		}

		if cursorPoint.In(drawnRectangle) {
			for linesIdx := range t.measurements.processedLines {
				ly := float64(p.Y) + t.measurements.lineHeight*float64(linesIdx)
				lx := float64(p.X)
				switch t.horizontalPosition {
				case TextPositionCenter:
					lx += ((float64(w) - t.measurements.processedLineWidths[linesIdx]) / 2) + float64(t.Inset.Left)
				case TextPositionEnd:
					lx += float64(w) - t.measurements.processedLineWidths[linesIdx] - float64(t.Inset.Right)
				case TextPositionStart:
					lx += float64(t.Inset.Left)
				}
				hoverLX := lx
				for idx := range t.measurements.processedLines[linesIdx] {
					wordWidth, _ := text.Measure(t.measurements.processedLines[linesIdx][idx].text, t.Face, 0)

					if t.measurements.processedLines[linesIdx][idx].linkValue != nil {
						if cursorPoint.In(image.Rect(int(hoverLX), int(ly), int(hoverLX+wordWidth), int(ly+t.measurements.lineHeight))) {
							input.SetCursorShape(input.CURSOR_POINTER)
							t.currentLink = t.measurements.processedLines[linesIdx][idx]
							t.measurements.processedLines[linesIdx][idx].hovered = true
							for additionalIdx := range t.measurements.processedLines[linesIdx][idx].linkValue.textBlocks {
								t.measurements.processedLines[linesIdx][idx].linkValue.textBlocks[additionalIdx].hovered = true
							}
						}
					}
					hoverLX += float64(wordWidth)
				}
			}
		}
	}

	// Draw text
	for linesIdx := range t.measurements.processedLines {
		ly := float64(p.Y) + t.measurements.lineHeight*float64(linesIdx)
		if ly > float64(screen.Bounds().Max.Y) {
			return
		}
		if ly < -t.measurements.lineHeight {
			continue
		}
		if t.widget.parent != nil {
			if ly < float64(t.widget.parent.Rect.Min.Y)-t.measurements.lineHeight {
				continue
			}
			if ly-t.measurements.lineHeight > float64(t.widget.parent.Rect.Max.Y) {
				return
			}
		}

		lx := float64(p.X)
		switch t.horizontalPosition {
		case TextPositionCenter:
			lx += ((float64(w) - t.measurements.processedLineWidths[linesIdx]) / 2) + float64(t.Inset.Left)
		case TextPositionEnd:
			lx += float64(w) - t.measurements.processedLineWidths[linesIdx] - float64(t.Inset.Right)
		case TextPositionStart:
			lx += float64(t.Inset.Left)
		}

		if t.ProcessBBCode {
			for _, piece := range t.measurements.processedLines[linesIdx] {
				wordWidth, _ := text.Measure(piece.text, t.Face, 0)

				op := &text.DrawOptions{}
				op.GeoM.Translate(lx, ly)
				if piece.linkValue != nil {
					if piece.hovered {
						op.ColorScale.ScaleWithColor(t.LinkColor.Hover)
					} else {
						op.ColorScale.ScaleWithColor(t.LinkColor.Idle)
					}
				} else {
					op.ColorScale.ScaleWithColor(piece.color)
				}
				text.Draw(screen, piece.text, t.Face, op)
				lx += float64(wordWidth)
			}

		} else {
			op := &text.DrawOptions{}
			op.GeoM.Translate(lx, ly)
			op.ColorScale.ScaleWithColor(t.Color)
			lineStr := ""

			for _, piece := range t.measurements.processedLines[linesIdx] {
				lineStr += piece.text
			}

			text.Draw(screen, lineStr, t.Face, op)
		}
	}
}

func (t *Text) handleBBCodeColor(line string) ([]*bbCodeText, color.Color, *linkData) {
	var newColor = *t.colorList.Top()
	var link = t.linkStack.Top()
	tokens := bbcode.Lex(line)
	tree := bbcode.Parse(tokens)
	return t.processTree(tree, newColor, link)
}

func (t *Text) processTree(node *bbcode.BBCodeNode, newColor color.Color, linkVal *linkData) ([]*bbCodeText, color.Color, *linkData) {
	var result []*bbCodeText

	switch node.ID {
	case bbcode.TEXT:
		if nodeVal, ok := node.Value.(string); ok {
			tb := bbCodeText{text: nodeVal, color: newColor, linkValue: linkVal}
			if linkVal != nil {
				linkVal.textBlocks = append(linkVal.textBlocks, &tb)
				linkVal.text += nodeVal
			}
			result = append(result, &tb)
		}

		for _, child := range node.Children {
			var iresult []*bbCodeText
			iresult, newColor, linkVal = t.processTree(child, newColor, linkVal)
			result = append(result, iresult...)
		}
	case bbcode.CLOSING_TAG:
		// Handle changing color back.
		if nodeVal, ok := node.Value.(bbcode.BBClosingTag); ok {
			switch nodeVal.Name {
			case COLOR_TAG:
				if t.colorList.Size() > 1 {
					t.colorList.Pop()
				}
				newColor = *t.colorList.Top()
			case LINK_TAG:
				t.linkStack.Pop()
				linkVal = t.linkStack.Top()
			}
		}
		for _, child := range node.Children {
			var iresult []*bbCodeText
			iresult, newColor, linkVal = t.processTree(child, newColor, linkVal)
			result = append(result, iresult...)
		}
	default:
		if node.GetOpeningTag() != nil {
			switch node.GetOpeningTag().Name {
			case COLOR_TAG:
				c, err := colorutil.HexToColor(node.GetOpeningTag().Value)
				if err == nil {
					t.colorList.Push(&c)
					newColor = c
				}
			case LINK_TAG:
				linkVal = &linkData{id: node.GetOpeningTag().Value, args: node.GetOpeningTag().Args, textBlocks: []*bbCodeText{}}
				t.linkStack.Push(linkVal)
			}
		}

		for _, child := range node.Children {
			var iresult []*bbCodeText
			iresult, newColor, linkVal = t.processTree(child, newColor, linkVal)
			result = append(result, iresult...)
		}
		if node.ClosingTag != nil {
			switch node.ClosingTag.Name {
			case COLOR_TAG:
				if t.colorList.Size() > 1 {
					t.colorList.Pop()
				}
				newColor = *t.colorList.Top()
			case LINK_TAG:
				t.linkStack.Pop()
				linkVal = t.linkStack.Top()
			}
		}

	}

	return result, newColor, linkVal
}

func (t *Text) measure() {
	if t.Label == t.measurements.label && t.Face == t.measurements.face && t.MaxWidth == t.measurements.maxWidth && t.ProcessBBCode == t.measurements.ProcessBBCode {
		return
	}
	m := t.Face.Metrics()
	t.measurements = textMeasurements{
		label:         t.Label,
		face:          t.Face,
		ProcessBBCode: t.ProcessBBCode,
		ascent:        m.HAscent,
		maxWidth:      t.MaxWidth,
	}

	sWidth, sHeight := text.Measure(" ", t.measurements.face, 0)

	fh := m.HAscent + m.HDescent
	t.measurements.lineHeight = sHeight
	ld := t.measurements.lineHeight - fh

	s := bufio.NewScanner(strings.NewReader(t.Label))
	for s.Scan() {
		if t.MaxWidth > 0 || t.ProcessBBCode {
			var newLine []*bbCodeText
			newLineWidth := float64(t.Inset.Left + t.Inset.Right)

			blocks, _, _ := t.handleBBCodeColor(s.Text())

			for idx := range blocks {
				words := strings.Split(blocks[idx].text, " ")
				for i, word := range words {
					var wordWidth float64
					wordWidth, _ = text.Measure(word, t.Face, 0)

					// Don't add the space to the last chunk.
					if i != len(words)-1 {
						wordWidth += sWidth
					}

					// If the new word doesn't push this past the max width continue adding to the current line
					if t.MaxWidth == 0 || newLineWidth+wordWidth < t.MaxWidth {
						wordBlock := bbCodeText{text: word, color: blocks[idx].color, linkValue: blocks[idx].linkValue}
						if i != len(words)-1 {
							wordBlock.text += " "
						}
						if wordBlock.linkValue != nil {
							wordBlock.linkValue.textBlocks = append(wordBlock.linkValue.textBlocks, &wordBlock)
						}
						newLine = append(newLine, &wordBlock)

						newLineWidth += wordWidth
					} else {
						// If the new word would push this past the max width save off the current line and start a new one
						if len(newLine) != 0 {
							t.measurements.processedLines = append(t.measurements.processedLines, newLine)
							t.measurements.processedLineWidths = append(t.measurements.processedLineWidths, newLineWidth)

							if newLineWidth > t.measurements.boundingBoxWidth {
								t.measurements.boundingBoxWidth = newLineWidth
							}
						}
						wordBlock := bbCodeText{text: word, color: blocks[idx].color, linkValue: blocks[idx].linkValue}
						if wordBlock.linkValue != nil {
							wordBlock.linkValue.textBlocks = append(wordBlock.linkValue.textBlocks, &wordBlock)
						}
						if i != len(words)-1 {
							wordBlock.text += " "
						}
						newLine = []*bbCodeText{&wordBlock}
						newLineWidth = wordWidth + float64(t.Inset.Left+t.Inset.Right)
					}
				}
			}

			// Save the final line
			if len(newLine) != 0 {
				t.measurements.processedLines = append(t.measurements.processedLines, newLine)
				t.measurements.processedLineWidths = append(t.measurements.processedLineWidths, newLineWidth)

				if newLineWidth > t.measurements.boundingBoxWidth {
					t.measurements.boundingBoxWidth = newLineWidth
				}
			}
		} else {
			line := s.Text()
			t.measurements.processedLines = append(t.measurements.processedLines, []*bbCodeText{{text: line}})
			lw, _ := text.Measure(line, t.Face, 0)
			lw += float64(t.Inset.Left + t.Inset.Right)
			t.measurements.processedLineWidths = append(t.measurements.processedLineWidths, lw)

			if lw > t.measurements.boundingBoxWidth {
				t.measurements.boundingBoxWidth = lw
			}
		}
	}

	t.measurements.boundingBoxHeight = float64(len(t.measurements.processedLines))*t.measurements.lineHeight - ld
}

func (t *Text) createWidget() {
	t.widget = NewWidget(t.widgetOpts...)
	t.widgetOpts = nil
}
