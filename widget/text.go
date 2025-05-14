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

type TextParams struct {
	Face    *text.Face
	Color   color.Color
	Padding *Insets
}

type Text struct {
	definedParams  TextParams
	computedParams TextParams
	Label          string
	MaxWidth       float64
	ProcessBBCode  bool
	LinkColor      TextLinkColor
	StripBBCode    bool

	widgetOpts         []WidgetOpt
	horizontalPosition TextPosition
	verticalPosition   TextPosition

	init         *MultiOnce
	widget       *Widget
	measurements textMeasurements
	colorList    *datastructures.Stack[color.Color]
	linkStack    *datastructures.Stack[linkData]
	currentLink  *bbCodeText
	previousLink *bbCodeText

	LinkClickedEvent       *event.Event
	LinkCursorEnteredEvent *event.Event
	LinkCursorExitedEvent  *event.Event
}

type textMeasurements struct {
	label         string
	face          *text.Face
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

type LinkEventArgs struct {
	Text    *Text
	Id      string
	Value   string
	Args    map[string]string
	OffsetX int
	OffsetY int
}

type LinkHandlerFunc func(args *LinkEventArgs)

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
		init:                   &MultiOnce{},
		LinkClickedEvent:       &event.Event{},
		LinkCursorEnteredEvent: &event.Event{},
		LinkCursorExitedEvent:  &event.Event{},
		colorList:              &datastructures.Stack[color.Color]{},
		linkStack:              &datastructures.Stack[linkData]{},
		LinkColor: TextLinkColor{
			Idle:  color.NRGBA{R: 0, G: 0, B: 250, A: 255},
			Hover: color.NRGBA{R: 255, G: 0, B: 250, A: 255},
		},
	}

	t.init.Append(t.createWidget)

	for _, o := range opts {
		o(t)
	}

	return t
}

func (t *Text) Validate() {
	t.init.Do()
	t.populateComputedParams()

	if t.computedParams.Color == nil {
		panic("Text: Color is required.")
	}
	if t.computedParams.Face == nil {
		panic("Text: Face is required.")
	}

	t.colorList.Push(&t.computedParams.Color)
}

func (t *Text) populateComputedParams() {
	txtParams := TextParams{}
	theme := t.widget.GetTheme()
	if theme != nil {
		if theme.TextTheme != nil {
			txtParams.Color = theme.TextTheme.Color
			txtParams.Face = theme.TextTheme.Face
			txtParams.Padding = theme.TextTheme.Padding
		}
	}
	if t.definedParams.Face != nil {
		txtParams.Face = t.definedParams.Face
	}

	if t.definedParams.Padding != nil {
		txtParams.Padding = t.definedParams.Padding
	}
	if t.definedParams.Color != nil {
		txtParams.Color = t.definedParams.Color
	}

	if txtParams.Padding == nil {
		txtParams.Padding = &Insets{}
	}

	t.computedParams = txtParams
}

func (o TextOptions) WidgetOpts(opts ...WidgetOpt) TextOpt {
	return func(t *Text) {
		t.widgetOpts = append(t.widgetOpts, opts...)
	}
}

// Text combines three options: TextLabel, TextFace and TextColor.
// It can be used for the inline configurations of Text object while
// separate functions are useful for a multi-step configuration.
func (o TextOptions) Text(label string, face *text.Face, color color.Color) TextOpt {
	return func(t *Text) {
		t.Label = label
		t.definedParams.Face = face
		t.definedParams.Color = color
	}
}

func (o TextOptions) TextLabel(label string) TextOpt {
	return func(t *Text) {
		t.Label = label
	}
}

func (o TextOptions) TextFace(face *text.Face) TextOpt {
	return func(t *Text) {
		t.definedParams.Face = face
	}
}

func (o TextOptions) TextColor(color color.Color) TextOpt {
	return func(t *Text) {
		t.definedParams.Color = color
	}
}

func (o TextOptions) Padding(padding *Insets) TextOpt {
	return func(t *Text) {
		t.definedParams.Padding = padding
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

// Set whether or not the text area should automatically strip out BBCodes from being displayed.
func (o TextOptions) StripBBCode(stripBBCode bool) TextOpt {
	return func(t *Text) {
		t.StripBBCode = stripBBCode
	}
}

// This option sets the idle and hover color for text that is wrapped in a
// [link][/link] bbcode.
//
// Note: this is only used if ProcessBBCode is true.
func (o TextOptions) LinkColor(linkColor *TextLinkColor) TextOpt {
	return func(t *Text) {
		if linkColor != nil {
			t.LinkColor = *linkColor
		}
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
func (o TextOptions) LinkClickedHandler(f LinkHandlerFunc) TextOpt {
	return func(b *Text) {
		if f != nil {
			b.LinkClickedEvent.AddHandler(func(args any) {
				if arg, ok := args.(*LinkEventArgs); ok {
					f(arg)
				}
			})
		}
	}
}

// Defines the handler to be called when the cursor enters a BBCode defined link.
//
// Note: this is only used if ProcessBBCode is true.
func (o TextOptions) LinkCursorEnteredHandler(f LinkHandlerFunc) TextOpt {
	return func(b *Text) {
		if f != nil {
			b.LinkCursorEnteredEvent.AddHandler(func(args any) {
				if arg, ok := args.(*LinkEventArgs); ok {
					f(arg)
				}
			})
		}
	}
}

// Defines the handler to be called when the cursor enters a BBCode defined link.
//
// Note: this is only used if ProcessBBCode is true.
func (o TextOptions) LinkCursorExitedHandler(f LinkHandlerFunc) TextOpt {
	return func(b *Text) {
		if f != nil {
			b.LinkCursorExitedEvent.AddHandler(func(args any) {
				if arg, ok := args.(*LinkEventArgs); ok {
					f(arg)
				}
			})
		}
	}
}

func (t *Text) SetColor(color color.Color) {
	t.definedParams.Color = color
	if t.definedParams.Color == nil {
		t.Validate()
	} else {
		t.computedParams.Color = color
	}
}

func (t *Text) SetFace(face *text.Face) {
	t.definedParams.Face = face
	if t.definedParams.Face == nil {
		t.Validate()
	} else {
		t.computedParams.Face = face
	}
}

func (t *Text) SetPadding(padding *Insets) {
	t.definedParams.Padding = padding
	t.Validate()
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
	w := int(math.Ceil(t.measurements.boundingBoxWidth)) + t.computedParams.Padding.Left + t.computedParams.Padding.Right
	h := int(math.Ceil(t.measurements.boundingBoxHeight)) + t.computedParams.Padding.Top + t.computedParams.Padding.Bottom

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
	if t.ProcessBBCode {
		t.handleLinkEvents()
	}
}

func (t *Text) handleLinkEvents() {
	if t.previousLink != nil && (t.currentLink == nil || t.currentLink.linkValue != t.previousLink.linkValue) {
		if t.LinkCursorExitedEvent != nil {
			off := t.getCursorOffset()
			t.LinkCursorExitedEvent.Fire(&LinkEventArgs{
				Text:    t,
				Value:   t.previousLink.linkValue.text,
				Id:      t.previousLink.linkValue.id,
				Args:    t.previousLink.linkValue.args,
				OffsetX: off.X,
				OffsetY: off.Y,
			})
		}
	}

	if t.currentLink != nil && (t.previousLink == nil || t.currentLink.linkValue != t.previousLink.linkValue) {
		if t.LinkCursorEnteredEvent != nil {
			off := t.getCursorOffset()
			t.LinkCursorEnteredEvent.Fire(&LinkEventArgs{
				Text:    t,
				Value:   t.currentLink.linkValue.text,
				Id:      t.currentLink.linkValue.id,
				Args:    t.currentLink.linkValue.args,
				OffsetX: off.X,
				OffsetY: off.Y,
			})
		}
	}

	if t.currentLink != nil && input.MouseButtonJustPressed(ebiten.MouseButton0) {
		if t.LinkClickedEvent != nil {
			off := t.getCursorOffset()
			t.LinkClickedEvent.Fire(&LinkEventArgs{
				Text:    t,
				Value:   t.currentLink.linkValue.text,
				Id:      t.currentLink.linkValue.id,
				Args:    t.currentLink.linkValue.args,
				OffsetX: off.X,
				OffsetY: off.Y,
			})
		}
	}

	t.previousLink = t.currentLink
	t.currentLink = nil
}

func (t *Text) draw(screen *ebiten.Image) {
	t.measure()

	r := t.widget.Rect
	w := r.Dx()
	p := r.Min

	switch t.verticalPosition {
	case TextPositionStart:
		p = p.Add(image.Point{0, t.computedParams.Padding.Top})
	case TextPositionCenter:
		p = p.Add(image.Point{0, int((float64(r.Dy())-t.measurements.boundingBoxHeight)/2 + float64(t.computedParams.Padding.Top))})
	case TextPositionEnd:
		p = p.Add(image.Point{0, int(float64(r.Dy())-t.measurements.boundingBoxHeight) - t.computedParams.Padding.Bottom})
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
					lx += ((float64(w) - t.measurements.processedLineWidths[linesIdx]) / 2) + float64(t.computedParams.Padding.Left)
				case TextPositionEnd:
					lx += float64(w) - t.measurements.processedLineWidths[linesIdx] - float64(t.computedParams.Padding.Right)
				case TextPositionStart:
					lx += float64(t.computedParams.Padding.Left)
				}
				hoverLX := lx
				for idx := range t.measurements.processedLines[linesIdx] {
					wordWidth, _ := text.Measure(t.measurements.processedLines[linesIdx][idx].text, *t.computedParams.Face, 0)

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
			lx += ((float64(w) - t.measurements.processedLineWidths[linesIdx]) / 2) + float64(t.computedParams.Padding.Left)
		case TextPositionEnd:
			lx += float64(w) - t.measurements.processedLineWidths[linesIdx] - float64(t.computedParams.Padding.Right)
		case TextPositionStart:
			lx += float64(t.computedParams.Padding.Left)
		}

		if t.ProcessBBCode {
			for _, piece := range t.measurements.processedLines[linesIdx] {
				wordWidth, _ := text.Measure(piece.text, *t.computedParams.Face, 0)

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
				text.Draw(screen, piece.text, *t.computedParams.Face, op)
				lx += float64(wordWidth)
			}

		} else {
			op := &text.DrawOptions{}
			op.GeoM.Translate(lx, ly)
			op.ColorScale.ScaleWithColor(t.computedParams.Color)
			lineStr := ""

			for _, piece := range t.measurements.processedLines[linesIdx] {
				lineStr += piece.text
			}

			text.Draw(screen, lineStr, *t.computedParams.Face, op)
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
			switch {
			case nodeVal.Name == COLOR_TAG && (t.StripBBCode || t.ProcessBBCode):
				if t.colorList.Size() > 1 {
					t.colorList.Pop()
				}
				newColor = *t.colorList.Top()
			case nodeVal.Name == LINK_TAG && (t.StripBBCode || t.ProcessBBCode):
				t.linkStack.Pop()
				linkVal = t.linkStack.Top()
			case !t.StripBBCode:
				if nodeVal, ok := node.Value.(bbcode.BBClosingTag); ok {
					tb := bbCodeText{text: nodeVal.Raw, color: newColor, linkValue: linkVal}
					if linkVal != nil {
						linkVal.textBlocks = append(linkVal.textBlocks, &tb)
						linkVal.text += nodeVal.Raw
					}
					result = append(result, &tb)
				}
			}
		}
		for _, child := range node.Children {
			var iresult []*bbCodeText
			iresult, newColor, linkVal = t.processTree(child, newColor, linkVal)
			result = append(result, iresult...)
		}
	default:
		if node.GetOpeningTag() != nil {
			switch {
			case node.GetOpeningTag().Name == COLOR_TAG && (t.StripBBCode || t.ProcessBBCode):
				c, err := colorutil.HexToColor(node.GetOpeningTag().Value)
				if err == nil {
					t.colorList.Push(&c)
					newColor = c
				}
			case node.GetOpeningTag().Name == LINK_TAG && (t.StripBBCode || t.ProcessBBCode):
				linkVal = &linkData{id: node.GetOpeningTag().Value, args: node.GetOpeningTag().Args, textBlocks: []*bbCodeText{}}
				t.linkStack.Push(linkVal)
			case !t.StripBBCode:
				if nodeVal, ok := node.Value.(bbcode.BBOpeningTag); ok {
					tb := bbCodeText{text: nodeVal.Raw, color: newColor, linkValue: linkVal}
					if linkVal != nil {
						linkVal.textBlocks = append(linkVal.textBlocks, &tb)
						linkVal.text += nodeVal.Raw
					}
					result = append(result, &tb)
				}
			}
		}

		for _, child := range node.Children {
			var iresult []*bbCodeText
			iresult, newColor, linkVal = t.processTree(child, newColor, linkVal)
			result = append(result, iresult...)
		}
		if node.ClosingTag != nil {
			switch {
			case node.ClosingTag.Name == COLOR_TAG && (t.StripBBCode || t.ProcessBBCode):
				if t.colorList.Size() > 1 {
					t.colorList.Pop()
				}
				newColor = *t.colorList.Top()
			case node.ClosingTag.Name == LINK_TAG && (t.StripBBCode || t.ProcessBBCode):
				t.linkStack.Pop()
				linkVal = t.linkStack.Top()
			case !t.StripBBCode:
				tb := bbCodeText{text: node.ClosingTag.Raw, color: newColor, linkValue: linkVal}
				if linkVal != nil {
					linkVal.textBlocks = append(linkVal.textBlocks, &tb)
					linkVal.text += node.ClosingTag.Raw
				}
				result = append(result, &tb)
			}
		}

	}

	return result, newColor, linkVal
}

func (t *Text) measure() {
	if t.Label == t.measurements.label && t.computedParams.Face == t.measurements.face && t.MaxWidth == t.measurements.maxWidth && t.ProcessBBCode == t.measurements.ProcessBBCode {
		return
	}
	m := (*t.computedParams.Face).Metrics()
	t.measurements = textMeasurements{
		label:         t.Label,
		face:          t.computedParams.Face,
		ProcessBBCode: t.ProcessBBCode,
		ascent:        m.HAscent,
		maxWidth:      t.MaxWidth,
	}

	sWidth, sHeight := text.Measure(" ", *t.measurements.face, 0)

	fh := m.HAscent + m.HDescent
	t.measurements.lineHeight = sHeight
	ld := t.measurements.lineHeight - fh

	s := bufio.NewScanner(strings.NewReader(t.Label))
	for s.Scan() {
		if t.MaxWidth > 0 || t.ProcessBBCode {
			var newLine []*bbCodeText
			newLineWidth := float64(t.computedParams.Padding.Left + t.computedParams.Padding.Right)

			blocks, _, _ := t.handleBBCodeColor(s.Text())

			for idx := range blocks {
				words := strings.Split(blocks[idx].text, " ")
				for i, word := range words {
					var wordWidth float64
					wordWidth, _ = text.Measure(word, *t.computedParams.Face, 0)

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
						newLineWidth = wordWidth + float64(t.computedParams.Padding.Left+t.computedParams.Padding.Right)
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
			lw, _ := text.Measure(line, *t.computedParams.Face, 0)
			lw += float64(t.computedParams.Padding.Left + t.computedParams.Padding.Right)
			t.measurements.processedLineWidths = append(t.measurements.processedLineWidths, lw)

			if lw > t.measurements.boundingBoxWidth {
				t.measurements.boundingBoxWidth = lw
			}
		}
	}

	t.measurements.boundingBoxHeight = float64(len(t.measurements.processedLines))*t.measurements.lineHeight - ld
}

func (t *Text) getCursorOffset() image.Point {
	x, y := input.CursorPosition()
	p := image.Point{x, y}
	return p.Sub(t.widget.Rect.Min)
}

func (t *Text) createWidget() {
	t.widget = NewWidget(t.widgetOpts...)
	t.widgetOpts = nil
}
