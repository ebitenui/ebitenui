package widget

import (
	img "image"
	"image/color"
	"math"
	"runtime"
	"strings"
	"sync/atomic"
	"time"

	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/input"
	"github.com/ebitenui/ebitenui/internal/jsUtil"
	"github.com/ebitenui/ebitenui/utilities/mobile"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type TextInputParams struct {
	Image                *TextInputImage
	Color                *TextInputColor
	Padding              *Insets
	Face                 *text.Face
	RepeatDelay          *time.Duration
	RepeatInterval       *time.Duration
	Secure               *bool
	ClearOnSubmit        *bool
	IgnoreEmptySubmit    *bool
	AllowDuplicateSubmit *bool
	SubmitOnEnter        *bool
	ScrollSensitivity    *int
	CaretWidth           *int
	Multiline            *bool
}

type TextInput struct {
	definedParams  TextInputParams
	computedParams TextInputParams

	ChangedEvent *event.Event
	SubmitEvent  *event.Event

	inputText       string
	validationFunc  TextInputValidationFunc
	placeholderText string
	mobileInputMode mobile.InputMode

	widgetOpts            []WidgetOpt
	init                  *MultiOnce
	commandToFunc         map[textInputControlCommand]textInputCommandFunc
	widget                *Widget
	caret                 *Caret
	text                  *Text
	renderBuf             *image.MaskedRenderBuffer
	mask                  *image.NineSlice
	cursorPosition        int
	state                 textInputState
	scrollOffset          int
	scrollOffsetY         int
	lastInputText         string
	previousSubmittedText *string
	dragStartIndex        int

	tabOrder int
	focused  bool
	focusMap map[FocusDirection]Focuser
}

type TextInputOpt func(t *TextInput)

type TextInputOptions struct {
}

type TextInputChangedEventArgs struct {
	TextInput *TextInput
	InputText string
}

type TextInputChangedHandlerFunc func(args *TextInputChangedEventArgs)

type TextInputImage struct {
	Idle     *image.NineSlice
	Disabled *image.NineSlice
	// Highlight defaults to image.NewNineSliceColor(color.NRGBA{6, 67, 161, 100}).
	Highlight *image.NineSlice
}

type TextInputColor struct {
	Idle          color.Color
	Disabled      color.Color
	Caret         color.Color
	DisabledCaret color.Color
}

type TextInputValidationFunc func(newInputText string) (bool, *string)

type textInputState func() (textInputState, bool)

type textInputControlCommand int

type textInputCommandFunc func()

var TextInputOpts TextInputOptions

const (
	textInputGoLeft = textInputControlCommand(iota + 1)
	textInputGoRight
	textInputGoStart
	textInputGoEnd
	textInputBackspace
	textInputDelete
	textInputEnter
	textInputEscape
	textInputGoUp
	textInputGoDown
)

var textInputKeyToCommand = map[ebiten.Key]textInputControlCommand{
	ebiten.KeyLeft:        textInputGoLeft,
	ebiten.KeyRight:       textInputGoRight,
	ebiten.KeyHome:        textInputGoStart,
	ebiten.KeyEnd:         textInputGoEnd,
	ebiten.KeyBackspace:   textInputBackspace,
	ebiten.KeyDelete:      textInputDelete,
	ebiten.KeyEnter:       textInputEnter,
	ebiten.KeyNumpadEnter: textInputEnter,
	ebiten.KeyEscape:      textInputEscape,
	ebiten.KeyArrowUp:     textInputGoUp,
	ebiten.KeyArrowDown:   textInputGoDown,
}

func NewTextInput(opts ...TextInputOpt) *TextInput {
	t := &TextInput{
		ChangedEvent: &event.Event{},
		SubmitEvent:  &event.Event{},

		init:          &MultiOnce{},
		commandToFunc: map[textInputControlCommand]textInputCommandFunc{},
		renderBuf:     image.NewMaskedRenderBuffer(),

		mobileInputMode: mobile.TEXT,
		focusMap:        make(map[FocusDirection]Focuser),
		dragStartIndex:  -1,
	}
	t.state = t.idleState(true)

	t.commandToFunc[textInputGoLeft] = t.CursorMoveLeft
	t.commandToFunc[textInputGoRight] = t.CursorMoveRight
	t.commandToFunc[textInputGoStart] = t.CursorMoveStart
	t.commandToFunc[textInputGoEnd] = t.CursorMoveEnd
	t.commandToFunc[textInputBackspace] = t.Backspace
	t.commandToFunc[textInputDelete] = t.Delete
	t.commandToFunc[textInputEnter] = t.submitWithEnter
	t.commandToFunc[textInputEscape] = t.DeselectText
	t.commandToFunc[textInputGoUp] = t.CursorMoveUp
	t.commandToFunc[textInputGoDown] = t.CursorMoveDown

	t.init.Append(t.createWidget)

	for _, o := range opts {
		o(t)
	}

	return t
}

func (t *TextInput) Validate() {
	t.init.Do()
	t.populateComputedParams()

	if t.computedParams.Face == nil {
		panic("TextInput: Font Face is required.")
	}
	if t.computedParams.Color == nil {
		panic("TextInput: Color is required.")
	}
	if t.computedParams.Color.Idle == nil {
		panic("TextInput: Color.Idle is required.")
	}

	t.initWidget()
}

func (t *TextInput) populateComputedParams() {
	TRUE := true
	FALSE := false
	params := TextInputParams{Color: &TextInputColor{}, Image: &TextInputImage{}}

	theme := t.GetWidget().GetTheme()

	// Set theme values
	if theme != nil {
		if theme.TextInputTheme != nil {
			params.AllowDuplicateSubmit = theme.TextInputTheme.AllowDuplicateSubmit
			params.ClearOnSubmit = theme.TextInputTheme.ClearOnSubmit
			if theme.TextInputTheme.Color != nil {
				params.Color.Idle = theme.TextInputTheme.Color.Idle
				params.Color.Disabled = theme.TextInputTheme.Color.Disabled
				params.Color.Caret = theme.TextInputTheme.Color.Caret
				params.Color.DisabledCaret = theme.TextInputTheme.Color.DisabledCaret
			}
			if theme.TextInputTheme.Face != nil {
				params.Face = theme.TextInputTheme.Face
			} else {
				params.Face = theme.DefaultFace
			}
			params.IgnoreEmptySubmit = theme.TextInputTheme.IgnoreEmptySubmit
			if theme.TextInputTheme.Image != nil {
				params.Image.Idle = theme.TextInputTheme.Image.Idle
				params.Image.Disabled = theme.TextInputTheme.Image.Disabled
				params.Image.Highlight = theme.TextInputTheme.Image.Highlight
			}
			params.Padding = theme.TextInputTheme.Padding
			params.RepeatDelay = theme.TextInputTheme.RepeatDelay
			params.RepeatInterval = theme.TextInputTheme.RepeatInterval
			params.ScrollSensitivity = theme.TextInputTheme.ScrollSensitivity
			params.Secure = theme.TextInputTheme.Secure
			params.SubmitOnEnter = theme.TextInputTheme.SubmitOnEnter
			params.CaretWidth = theme.TextInputTheme.CaretWidth
		}
	}

	// Set Defined values
	if t.definedParams.AllowDuplicateSubmit != nil {
		params.AllowDuplicateSubmit = t.definedParams.AllowDuplicateSubmit
	}
	if t.definedParams.ClearOnSubmit != nil {
		params.ClearOnSubmit = t.definedParams.ClearOnSubmit
	}
	if t.definedParams.Color != nil {
		params.Color.Idle = t.definedParams.Color.Idle
		params.Color.Disabled = t.definedParams.Color.Disabled
		params.Color.Caret = t.definedParams.Color.Caret
		params.Color.DisabledCaret = t.definedParams.Color.DisabledCaret
	}
	if t.definedParams.Face != nil {
		params.Face = t.definedParams.Face
	}
	if t.definedParams.IgnoreEmptySubmit != nil {
		params.IgnoreEmptySubmit = t.definedParams.IgnoreEmptySubmit
	}
	if t.definedParams.Image != nil {
		params.Image.Idle = t.definedParams.Image.Idle
		params.Image.Disabled = t.definedParams.Image.Disabled
		params.Image.Highlight = t.definedParams.Image.Highlight
	}
	if t.definedParams.Padding != nil {
		params.Padding = t.definedParams.Padding
	}
	if t.definedParams.RepeatDelay != nil {
		params.RepeatDelay = t.definedParams.RepeatDelay
	}
	if t.definedParams.RepeatInterval != nil {
		params.RepeatDelay = t.definedParams.RepeatDelay
	}
	if t.definedParams.ScrollSensitivity != nil {
		params.ScrollSensitivity = t.definedParams.ScrollSensitivity
	}
	if t.definedParams.Secure != nil {
		params.Secure = t.definedParams.Secure
	}
	if t.definedParams.SubmitOnEnter != nil {
		params.SubmitOnEnter = t.definedParams.SubmitOnEnter
	}
	if t.definedParams.CaretWidth != nil {
		params.CaretWidth = t.definedParams.CaretWidth
	}
	if t.definedParams.Multiline != nil {
		params.Multiline = t.definedParams.Multiline
	}

	// Set Default values
	if params.Image == nil {
		params.Image = &TextInputImage{}
	}
	if params.Image.Highlight == nil {
		params.Image.Highlight = image.NewNineSliceColor(color.NRGBA{6, 67, 161, 100})
	}
	if params.RepeatDelay == nil {
		delay := 300 * time.Millisecond
		params.RepeatDelay = &delay
	}
	if params.RepeatInterval == nil {
		interval := 35 * time.Millisecond
		params.RepeatInterval = &interval
	}
	if params.ScrollSensitivity == nil {
		sensitivity := 15
		params.ScrollSensitivity = &sensitivity
	}
	if params.Padding == nil {
		params.Padding = &Insets{}
	}
	if params.Secure == nil {
		params.Secure = &FALSE
	}
	if params.ClearOnSubmit == nil {
		params.ClearOnSubmit = &FALSE
	}
	if params.IgnoreEmptySubmit == nil {
		params.IgnoreEmptySubmit = &FALSE
	}
	if params.AllowDuplicateSubmit == nil {
		params.AllowDuplicateSubmit = &FALSE
	}
	if params.SubmitOnEnter == nil {
		params.SubmitOnEnter = &TRUE
	}
	if params.CaretWidth == nil {
		width := 2
		params.CaretWidth = &width
	}
	if params.Multiline == nil {
		params.Multiline = &FALSE
	}
	if params.Color != nil && params.Color.Caret == nil {
		params.Color.Caret = params.Color.Idle
	}

	t.computedParams = params
}

func (o TextInputOptions) WidgetOpts(opts ...WidgetOpt) TextInputOpt {
	return func(t *TextInput) {
		t.widgetOpts = append(t.widgetOpts, opts...)
	}
}

func (o TextInputOptions) ChangedHandler(f TextInputChangedHandlerFunc) TextInputOpt {
	return func(t *TextInput) {
		t.ChangedEvent.AddHandler(func(args interface{}) {
			if arg, ok := args.(*TextInputChangedEventArgs); ok {
				f(arg)
			}
		})
	}
}

func (o TextInputOptions) SubmitHandler(f TextInputChangedHandlerFunc) TextInputOpt {
	return func(t *TextInput) {
		t.SubmitEvent.AddHandler(func(args interface{}) {
			if arg, ok := args.(*TextInputChangedEventArgs); ok {
				f(arg)
			}
		})
	}
}

func (o TextInputOptions) ClearOnSubmit(clearOnSubmit bool) TextInputOpt {
	return func(t *TextInput) {
		t.definedParams.ClearOnSubmit = &clearOnSubmit
	}
}

func (o TextInputOptions) IgnoreEmptySubmit(ignoreEmptySubmit bool) TextInputOpt {
	return func(t *TextInput) {
		t.definedParams.IgnoreEmptySubmit = &ignoreEmptySubmit
	}
}

func (o TextInputOptions) AllowDuplicateSubmit(allowDuplicateSubmit bool) TextInputOpt {
	return func(t *TextInput) {
		t.definedParams.AllowDuplicateSubmit = &allowDuplicateSubmit
	}
}

func (o TextInputOptions) Image(i *TextInputImage) TextInputOpt {
	return func(t *TextInput) {
		t.definedParams.Image = i
	}
}

func (o TextInputOptions) Color(c *TextInputColor) TextInputOpt {
	return func(t *TextInput) {
		t.definedParams.Color = c
	}
}

func (o TextInputOptions) Padding(i *Insets) TextInputOpt {
	return func(t *TextInput) {
		t.definedParams.Padding = i
	}
}

func (o TextInputOptions) Face(f *text.Face) TextInputOpt {
	return func(t *TextInput) {
		t.definedParams.Face = f
	}
}

func (o TextInputOptions) RepeatInterval(i time.Duration) TextInputOpt {
	return func(t *TextInput) {
		t.definedParams.RepeatInterval = &i
	}
}

func (o TextInputOptions) Validation(f TextInputValidationFunc) TextInputOpt {
	return func(t *TextInput) {
		t.validationFunc = f
	}
}

func (o TextInputOptions) Placeholder(s string) TextInputOpt {
	return func(t *TextInput) {
		t.placeholderText = s
	}
}

func (o TextInputOptions) Secure(b bool) TextInputOpt {
	return func(t *TextInput) {
		t.definedParams.Secure = &b
	}
}

func (o TextInputOptions) CaretWidth(caretWidth int) TextInputOpt {
	return func(t *TextInput) {
		t.definedParams.CaretWidth = &caretWidth
	}
}

func (o TextInputOptions) TabOrder(to int) TextInputOpt {
	return func(t *TextInput) {
		t.tabOrder = to
	}
}

// Sets if the input will submit when pressing enter or not.
func (o TextInputOptions) SubmitOnEnter(submitOnEnter bool) TextInputOpt {
	return func(t *TextInput) {
		t.definedParams.SubmitOnEnter = &submitOnEnter
	}
}

// Multiline enables multi-line editing. Enter inserts a newline instead of submitting.
// Arrow Up/Down navigate between lines.
func (o TextInputOptions) Multiline(multiline bool) TextInputOpt {
	return func(t *TextInput) {
		t.definedParams.Multiline = &multiline
	}
}

// Sets the keyboard type to use when viewed on a mobile browser.
//
// https://css-tricks.com/everything-you-ever-wanted-to-know-about-inputmode
func (o TextInputOptions) MobileInputMode(mobileInputMode mobile.InputMode) TextInputOpt {
	return func(t *TextInput) {
		t.mobileInputMode = mobileInputMode
	}
}

// Sets how many pixels from the edge the cursor must be dragged prior to it scrolling in that direction.
//
// Default: 15.
func (o TextInputOptions) ScrollSensitivity(scrollSensitivity int) TextInputOpt {
	return func(t *TextInput) {
		t.definedParams.ScrollSensitivity = &scrollSensitivity
	}
}

/*********** End of Configuration *****************/

func (t *TextInput) GetWidget() *Widget {
	t.init.Do()
	return t.widget
}

func (t *TextInput) SetLocation(rect img.Rectangle) {
	t.init.Do()
	t.widget.Rect = rect
}

func (t *TextInput) PreferredSize() (int, int) {
	t.init.Do()
	_, h := t.caret.PreferredSize()
	h = h + t.computedParams.Padding.Top + t.computedParams.Padding.Bottom
	w := 50

	if t.widget != nil && h < t.widget.MinHeight {
		h = t.widget.MinHeight
	}
	if t.widget != nil && w < t.widget.MinWidth {
		w = t.widget.MinWidth
	}

	return w, h
}

func (t *TextInput) Render(screen *ebiten.Image) {
	t.init.Do()

	t.widget.Render(screen)

	t.renderImage(screen)
	t.renderTextAndCaret(screen)
}

func (t *TextInput) Update(updObj *UpdateObject) {
	t.init.Do()
	t.text.GetWidget().Disabled = t.widget.Disabled
	if t.lastInputText != t.inputText {
		if t.validationFunc != nil {
			result, replacement := t.validationFunc(t.inputText)
			if !result {
				if replacement != nil {
					t.inputText = *replacement
				} else {
					t.inputText = t.lastInputText
				}
			}
		}
	}

	for {
		newState, rerun := t.state()
		if newState != nil {
			t.state = newState
		}
		if !rerun {
			break
		}
	}

	defer func() {
		t.lastInputText = t.inputText
	}()

	if t.inputText != t.lastInputText {
		t.ChangedEvent.Fire(&TextInputChangedEventArgs{
			TextInput: t,
			InputText: t.inputText,
		})
	}

	t.widget.Update(updObj)
	if t.text != nil {
		t.text.Update(updObj)
	}
	if t.caret != nil {
		t.caret.Update(updObj)
	}
}

func (t *TextInput) idleState(newKeyOrCommand bool) textInputState {
	return func() (textInputState, bool) {
		if !t.focused {
			t.dragStartIndex = -1
			return t.idleState(true), false
		}

		chars := input.InputChars()
		if len(chars) > 0 {
			if !ebiten.IsKeyPressed(ebiten.KeyControl) && !ebiten.IsKeyPressed(ebiten.KeyControlLeft) && !ebiten.IsKeyPressed(ebiten.KeyControlRight) {
				t.DeleteSelectedText()
				return t.charsInputState(string(chars)), true
			}
			t.DeselectText()
			return t.idleState(true), false
		}

		st := textInputCheckForCommand(t, newKeyOrCommand)
		if st != nil {
			return st, true
		}

		x, y := input.CursorPosition()
		p := img.Point{x, y}
		curIdx := 0
		tr := t.computedParams.Padding.Apply(t.widget.Rect)
		if x < tr.Min.X {
			x = tr.Min.X
		}
		if x > tr.Max.X {
			x = tr.Max.X
		}
		if *t.computedParams.Multiline {
			if p.In(t.widget.Rect) {
				curIdx = t.multilineCursorFromPoint(x, y, tr)
			} else if y < tr.Min.Y {
				curIdx = 0
			} else {
				curIdx = len([]rune(t.inputText))
			}
		} else if p.In(t.widget.Rect) {
			curIdx = fontStringIndex([]rune(t.inputText), t.computedParams.Face, x-t.scrollOffset-tr.Min.X)
		} else {
			if y < tr.Min.Y {
				curIdx = 0
			} else {
				curIdx = len([]rune(t.inputText))
			}
		}
		if input.MouseButtonJustPressedLayer(ebiten.MouseButtonLeft, t.widget.EffectiveInputLayer()) {
			t.dragStartIndex = curIdx
			t.cursorPosition = curIdx
		} else if input.MouseButtonPressed(ebiten.MouseButtonLeft) {
			if t.dragStartIndex == -1 {
				t.dragStartIndex = curIdx
			}
			t.cursorPosition = curIdx
			if !*t.computedParams.Multiline {
				textSize := tr.Dx() - fontAdvance(t.inputText, t.computedParams.Face)
				if t.scrollOffset < 0 && x < t.widget.Rect.Min.X+*t.computedParams.ScrollSensitivity {
					t.scrollOffset = min(0, t.scrollOffset+1)
				} else if t.scrollOffset > textSize && x > t.widget.Rect.Max.X-*t.computedParams.ScrollSensitivity {
					t.scrollOffset = max(t.scrollOffset-1, textSize)
				}
			}
		}

		if input.MouseButtonJustReleasedLayer(ebiten.MouseButtonLeft, t.widget.EffectiveInputLayer()) {
			t.cursorPosition = curIdx
			t.caret.ResetBlinking()

			if t.dragStartIndex == curIdx {
				t.dragStartIndex = -1
			}
		}
		if runtime.GOOS == jsUtil.JS && runtime.GOARCH == jsUtil.WASM {
			dragStartDraw := min(t.cursorPosition, t.dragStartIndex)
			dragEndDraw := max(t.cursorPosition, t.dragStartIndex)
			jsUtil.SetCursorPosition(dragStartDraw, dragEndDraw)
		}

		return t.idleState(true), false
	}
}

func textInputCheckForCommand(t *TextInput, newKeyOrCommand bool) textInputState {
	for key, cmd := range textInputKeyToCommand {
		if !input.KeyPressed(key) {
			continue
		}

		// Up/Down arrows are only used in multiline mode
		if !*t.computedParams.Multiline && (cmd == textInputGoUp || cmd == textInputGoDown) {
			continue
		}

		var delay time.Duration
		if newKeyOrCommand {
			delay = *t.computedParams.RepeatDelay
		} else {
			delay = *t.computedParams.RepeatInterval
		}

		return t.commandState(cmd, key, delay, nil, nil)
	}

	return nil
}

func (t *TextInput) charsInputState(c string) textInputState {
	return func() (textInputState, bool) {
		if !t.widget.Disabled {
			t.Insert(c)
		}

		t.caret.ResetBlinking()

		return t.idleState(true), false
	}
}

func (t *TextInput) commandState(cmd textInputControlCommand, key ebiten.Key, delay time.Duration, timer *time.Timer, expired *atomic.Value) textInputState {
	return func() (textInputState, bool) {
		if !input.KeyPressed(key) {
			return t.idleState(true), true
		}

		if timer != nil {
			if isExpired, _ := expired.Load().(bool); isExpired {
				return t.idleState(false), true
			}
		}

		if timer == nil {
			t.commandToFunc[cmd]()

			expired = &atomic.Value{}
			expired.Store(false)

			timer = time.AfterFunc(delay, func() {
				expired.Store(true)
			})

			return t.commandState(cmd, key, delay, timer, expired), false
		}

		return nil, false
	}
}

func (t *TextInput) Insert(c string) {
	t.init.Do()
	t.DeleteSelectedText()
	s := string(insertChars([]rune(t.inputText), []rune(c), t.cursorPosition))

	if t.validationFunc != nil {
		result, replacement := t.validationFunc(s)
		if !result {
			if replacement != nil {
				s = *replacement
			} else {
				return
			}
		}
	}
	t.inputText = s

	t.cursorPosition += len([]rune(c))
	if t.cursorPosition > len([]rune(t.inputText)) {
		t.cursorPosition = len([]rune(t.inputText))
	}
	t.caret.ResetBlinking()
}

func (t *TextInput) CursorMoveLeft() {
	t.init.Do()
	if t.cursorPosition > 0 {
		t.cursorPosition--
	}
	t.caret.ResetBlinking()
}

func (t *TextInput) CursorMoveRight() {
	t.init.Do()
	if t.cursorPosition < len([]rune(t.inputText)) {
		t.cursorPosition++
	}
	t.caret.ResetBlinking()
}

func (t *TextInput) CursorMoveStart() {
	t.init.Do()
	if t.computedParams.Multiline != nil && *t.computedParams.Multiline {
		maxWidth := t.contentWidth()
		wlines := wrapText(t.inputText, t.computedParams.Face, maxWidth)
		vLine, _ := cursorToWrappedLineCol(wlines, t.cursorPosition)
		t.cursorPosition = wlines[vLine].startRune
	} else {
		t.cursorPosition = 0
	}
	t.caret.ResetBlinking()
}

func (t *TextInput) CursorMoveEnd() {
	t.init.Do()
	if t.computedParams.Multiline != nil && *t.computedParams.Multiline {
		maxWidth := t.contentWidth()
		wlines := wrapText(t.inputText, t.computedParams.Face, maxWidth)
		vLine, _ := cursorToWrappedLineCol(wlines, t.cursorPosition)
		t.cursorPosition = wlines[vLine].startRune + len([]rune(wlines[vLine].text))
	} else {
		t.cursorPosition = len([]rune(t.inputText))
	}
	t.caret.ResetBlinking()
}

func (t *TextInput) CursorMoveUp() {
	t.init.Do()
	if !*t.computedParams.Multiline {
		return
	}
	maxWidth := t.contentWidth()
	wlines := wrapText(t.inputText, t.computedParams.Face, maxWidth)
	vLine, vCol := cursorToWrappedLineCol(wlines, t.cursorPosition)
	if vLine > 0 {
		prevLineRunes := len([]rune(wlines[vLine-1].text))
		col := vCol
		if col > prevLineRunes {
			col = prevLineRunes
		}
		t.cursorPosition = wlines[vLine-1].startRune + col
	}
	t.caret.ResetBlinking()
}

func (t *TextInput) CursorMoveDown() {
	t.init.Do()
	if !*t.computedParams.Multiline {
		return
	}
	maxWidth := t.contentWidth()
	wlines := wrapText(t.inputText, t.computedParams.Face, maxWidth)
	vLine, vCol := cursorToWrappedLineCol(wlines, t.cursorPosition)
	if vLine < len(wlines)-1 {
		nextLineRunes := len([]rune(wlines[vLine+1].text))
		col := vCol
		if col > nextLineRunes {
			col = nextLineRunes
		}
		t.cursorPosition = wlines[vLine+1].startRune + col
	}
	t.caret.ResetBlinking()
}

// contentWidth returns the available pixel width for text content (widget width minus padding).
func (t *TextInput) contentWidth() int {
	if t.widget == nil {
		return 0
	}
	return t.widget.Rect.Dx() - t.computedParams.Padding.Left - t.computedParams.Padding.Right
}

// multilineLineHeight returns the pixel height of one line for the given face.
func multilineLineHeight(f *text.Face) float64 {
	_, h := text.Measure(" ", *f, 0)
	return h
}

func (t *TextInput) Backspace() {
	t.init.Do()
	if !t.widget.Disabled {
		if t.dragStartIndex != -1 {
			t.DeleteSelectedText()
		} else if t.cursorPosition > 0 {
			t.inputText = string(removeChar([]rune(t.inputText), t.cursorPosition-1))
			t.cursorPosition--
		}
	}
	t.DeselectText()
	t.caret.ResetBlinking()
}

func (t *TextInput) Delete() {
	t.init.Do()
	if !t.widget.Disabled {
		if t.dragStartIndex != -1 {
			t.DeleteSelectedText()
		} else if t.cursorPosition < len([]rune(t.inputText)) {
			t.inputText = string(removeChar([]rune(t.inputText), t.cursorPosition))
		}
	}
	t.DeselectText()
	t.caret.ResetBlinking()
}

func (t *TextInput) submitWithEnter() {
	if *t.computedParams.Multiline {
		if !t.widget.Disabled {
			t.Insert("\n")
		}
		return
	}
	if *t.computedParams.SubmitOnEnter {
		t.Submit()
	}
}

func (t *TextInput) Submit() {
	t.init.Do()
	if !*t.computedParams.IgnoreEmptySubmit || len(t.inputText) > 0 {
		if *t.computedParams.AllowDuplicateSubmit || t.previousSubmittedText == nil || t.inputText != *t.previousSubmittedText {
			t.SubmitEvent.Fire(&TextInputChangedEventArgs{
				TextInput: t,
				InputText: t.inputText,
			})
			previousText := t.inputText
			t.previousSubmittedText = &previousText
		}
	}
	if *t.computedParams.ClearOnSubmit {
		t.CursorMoveStart()
		t.inputText = ""
	}

	t.DeselectText()
}

func (t *TextInput) SelectedText() string {
	t.init.Do()
	if t.dragStartIndex != -1 {
		start := min(t.dragStartIndex, t.cursorPosition)
		end := max(t.dragStartIndex, t.cursorPosition)

		return strings.Clone(t.inputText)[start:end]
	}
	return ""
}

func (t *TextInput) DeselectText() {
	t.dragStartIndex = -1
}

func (t *TextInput) SelectAll() {
	t.init.Do()
	if len(t.inputText) > 0 {
		t.dragStartIndex = 0
		t.CursorMoveEnd()
		if runtime.GOOS == jsUtil.JS && runtime.GOARCH == jsUtil.WASM {
			dragStartDraw := min(t.cursorPosition, t.dragStartIndex)
			dragEndDraw := max(t.cursorPosition, t.dragStartIndex)
			jsUtil.SetCursorPosition(dragStartDraw, dragEndDraw)
		}
	}
}

func (t *TextInput) DeleteSelectedText() {
	if t.dragStartIndex != -1 {
		start := min(t.dragStartIndex, t.cursorPosition)
		end := max(t.dragStartIndex, t.cursorPosition)
		t.inputText = strings.Replace(t.inputText, t.inputText[start:end], "", 1)
		if t.cursorPosition > t.dragStartIndex {
			t.cursorPosition -= (end - start)
		}

		t.dragStartIndex = -1
		t.caret.ResetBlinking()
	}
}

func insertChars(r []rune, c []rune, pos int) []rune {
	res := make([]rune, len(r)+len(c))
	copy(res, r[:pos])
	copy(res[pos:], c)
	copy(res[pos+len(c):], r[pos:])
	return res
}

func removeChar(r []rune, pos int) []rune {
	res := make([]rune, len(r)-1)
	copy(res, r[:pos])
	copy(res[pos:], r[pos+1:])
	return res
}

func (t *TextInput) renderImage(screen *ebiten.Image) {
	if t.computedParams.Image != nil && t.computedParams.Image.Idle != nil {
		i := t.computedParams.Image.Idle
		if t.widget.Disabled && t.computedParams.Image.Disabled != nil {
			i = t.computedParams.Image.Disabled
		}

		rect := t.widget.Rect
		i.Draw(screen, rect.Dx(), rect.Dy(), func(opts *ebiten.DrawImageOptions) {
			opts.GeoM.Translate(float64(rect.Min.X), float64(rect.Min.Y))
		})
	}
}

func (t *TextInput) renderTextAndCaret(screen *ebiten.Image) {
	t.renderBuf.Draw(screen,
		func(buf *ebiten.Image) {
			t.drawTextAndCaret(buf)
		},
		func(buf *ebiten.Image) {
			rect := t.widget.Rect
			t.mask.Draw(buf, rect.Dx()-t.computedParams.Padding.Left-t.computedParams.Padding.Right, rect.Dy()-t.computedParams.Padding.Top-t.computedParams.Padding.Bottom,
				func(opts *ebiten.DrawImageOptions) {
					opts.GeoM.Translate(float64(rect.Min.X+t.computedParams.Padding.Left), float64(rect.Min.Y+t.computedParams.Padding.Top))
					opts.CompositeMode = ebiten.CompositeModeCopy
				})
		})
}

func (t *TextInput) drawTextAndCaret(screen *ebiten.Image) {
	if *t.computedParams.Multiline {
		t.drawTextAndCaretMultiline(screen)
		return
	}

	rect := t.widget.Rect
	tr := rect
	tr = tr.Add(img.Point{t.computedParams.Padding.Left, t.computedParams.Padding.Top})

	inputStr := t.inputText
	if *t.computedParams.Secure {
		inputStr = strings.Repeat("*", len([]rune(t.inputText)))
	}

	cx := 0
	if t.focused {
		sub := string([]rune(inputStr)[:t.cursorPosition])
		cx = fontAdvance(sub, t.computedParams.Face)

		dx := tr.Min.X + t.scrollOffset + cx + t.caret.Width + t.computedParams.Padding.Right - rect.Max.X
		if dx > 0 {
			t.scrollOffset -= dx
		}

		dx = tr.Min.X + t.scrollOffset + cx - t.computedParams.Padding.Left - rect.Min.X
		if dx < 0 {
			t.scrollOffset -= dx
		}
		if t.dragStartIndex != -1 {
			dragString := string([]rune(inputStr)[:t.dragStartIndex])
			dragXStart := fontAdvance(dragString, t.computedParams.Face)

			dragStartDraw := min(dragXStart, cx)
			dragEndDraw := max(dragXStart, cx)

			// Change the Dx and the tr.Min.X based on selection
			t.computedParams.Image.Highlight.Draw(screen, dragEndDraw-dragStartDraw, tr.Dy(),
				func(opts *ebiten.DrawImageOptions) {
					opts.GeoM.Translate(float64(tr.Min.X+dragStartDraw+t.scrollOffset), float64(tr.Min.Y))
				})
		}

	}
	tr = tr.Add(img.Point{t.scrollOffset, 0})

	t.text.SetLocation(tr)
	if len([]rune(t.inputText)) > 0 {
		t.text.Label = inputStr
	} else {
		t.text.Label = t.placeholderText
	}
	if (t.widget.Disabled || len([]rune(t.inputText)) == 0) && t.computedParams.Color.Disabled != nil {
		t.text.SetColor(t.computedParams.Color.Disabled)
	} else {
		t.text.SetColor(t.computedParams.Color.Idle)
	}
	t.text.Render(screen)

	if t.focused {
		if t.widget.Disabled && t.computedParams.Color.DisabledCaret != nil {
			t.caret.Color = t.computedParams.Color.DisabledCaret
		} else {
			t.caret.Color = t.computedParams.Color.Caret
		}

		tr = tr.Add(img.Point{cx, 0})
		t.caret.SetLocation(tr)

		t.caret.Render(screen)
	}
}

// wrappedLine represents a visual line after word wrapping.
// startRune is the index into []rune(inputText) where this visual line starts.
type wrappedLine struct {
	text      string
	startRune int
}

// wrapText splits the input text into visual lines that fit within maxWidth pixels.
// Hard newlines (\n) always cause a line break.
func wrapText(s string, face *text.Face, maxWidth int) []wrappedLine {
	if maxWidth <= 0 {
		// Fallback: no wrapping, just split on newlines
		var result []wrappedLine
		pos := 0
		for _, line := range strings.Split(s, "\n") {
			result = append(result, wrappedLine{text: line, startRune: pos})
			pos += len([]rune(line)) + 1 // +1 for the \n
		}
		return result
	}

	runes := []rune(s)
	var result []wrappedLine
	lineStart := 0

	for lineStart <= len(runes) {
		// Find the end of the current logical line (up to next \n or end)
		lineEnd := len(runes)
		for i := lineStart; i < len(runes); i++ {
			if runes[i] == '\n' {
				lineEnd = i
				break
			}
		}

		// Wrap this logical line segment into visual lines
		segStart := lineStart
		for segStart <= lineEnd {
			if segStart == lineEnd {
				// Only add an empty visual line if this is a truly empty logical line
				// (e.g. from \n\n or trailing \n), not when we just finished wrapping.
				if segStart == lineStart {
					result = append(result, wrappedLine{text: "", startRune: segStart})
				}
				break
			}

			// Find how many runes fit in maxWidth
			fitEnd := segStart
			for fitEnd < lineEnd {
				w := fontAdvance(string(runes[segStart:fitEnd+1]), face)
				if w > maxWidth {
					break
				}
				fitEnd++
			}

			if fitEnd == segStart {
				// At least one character per visual line
				fitEnd = segStart + 1
			}

			// Try to break at a word boundary (space) if we didn't reach lineEnd
			if fitEnd < lineEnd {
				// Search backwards for a space to break at
				breakAt := fitEnd
				for breakAt > segStart {
					if runes[breakAt-1] == ' ' {
						break
					}
					breakAt--
				}
				if breakAt > segStart {
					fitEnd = breakAt
				}
			}

			result = append(result, wrappedLine{text: string(runes[segStart:fitEnd]), startRune: segStart})
			segStart = fitEnd
		}

		lineStart = lineEnd + 1 // skip the \n
	}

	if len(result) == 0 {
		result = append(result, wrappedLine{text: "", startRune: 0})
	}
	return result
}

// cursorToWrappedLineCol finds which visual line the cursor is on and the column within it.
func cursorToWrappedLineCol(wlines []wrappedLine, cursor int) (int, int) {
	for i := len(wlines) - 1; i >= 0; i-- {
		if cursor >= wlines[i].startRune {
			col := cursor - wlines[i].startRune
			lineRunes := len([]rune(wlines[i].text))
			if col > lineRunes {
				col = lineRunes
			}
			return i, col
		}
	}
	return 0, 0
}

func (t *TextInput) drawTextAndCaretMultiline(screen *ebiten.Image) {
	rect := t.widget.Rect
	tr := rect
	tr = tr.Add(img.Point{t.computedParams.Padding.Left, t.computedParams.Padding.Top})

	lineH := int(math.Round(multilineLineHeight(t.computedParams.Face)))
	maxWidth := t.contentWidth()
	wlines := wrapText(t.inputText, t.computedParams.Face, maxWidth)
	vLine, vCol := cursorToWrappedLineCol(wlines, t.cursorPosition)

	// Cursor pixel position within content
	var cx int
	if vLine < len(wlines) {
		lineRunes := []rune(wlines[vLine].text)
		if vCol > len(lineRunes) {
			vCol = len(lineRunes)
		}
		cx = fontAdvance(string(lineRunes[:vCol]), t.computedParams.Face)
	}
	cy := vLine * lineH

	// Vertical scroll to keep cursor visible
	visibleH := rect.Dy() - t.computedParams.Padding.Top - t.computedParams.Padding.Bottom
	if cy+lineH > -t.scrollOffsetY+visibleH {
		t.scrollOffsetY = -(cy + lineH - visibleH)
	}
	if cy < -t.scrollOffsetY {
		t.scrollOffsetY = -cy
	}

	// Determine display text and color
	textColor := t.computedParams.Color.Idle
	displayLines := wlines
	if (t.widget.Disabled || len(t.inputText) == 0) && t.computedParams.Color.Disabled != nil {
		textColor = t.computedParams.Color.Disabled
	}
	if len(t.inputText) == 0 {
		displayLines = []wrappedLine{{text: t.placeholderText, startRune: 0}}
	}

	// Render selection highlight
	if t.focused && t.dragStartIndex != -1 {
		selStart := min(t.dragStartIndex, t.cursorPosition)
		selEnd := max(t.dragStartIndex, t.cursorPosition)
		startLine, startCol := cursorToWrappedLineCol(wlines, selStart)
		endLine, endCol := cursorToWrappedLineCol(wlines, selEnd)

		for i := startLine; i <= endLine; i++ {
			lineRunes := []rune(wlines[i].text)
			var hlX0, hlX1 int
			if i == startLine {
				hlX0 = fontAdvance(string(lineRunes[:startCol]), t.computedParams.Face)
			}
			if i == endLine {
				if endCol > len(lineRunes) {
					endCol = len(lineRunes)
				}
				hlX1 = fontAdvance(string(lineRunes[:endCol]), t.computedParams.Face)
			} else {
				hlX1 = fontAdvance(string(lineRunes), t.computedParams.Face)
			}
			if hlW := hlX1 - hlX0; hlW > 0 {
				ly := i*lineH + t.scrollOffsetY
				t.computedParams.Image.Highlight.Draw(screen, hlW, lineH,
					func(opts *ebiten.DrawImageOptions) {
						opts.GeoM.Translate(float64(tr.Min.X+hlX0), float64(tr.Min.Y+ly))
					})
			}
		}
	}

	// Render each visual line
	for i, wl := range displayLines {
		ly := float64(tr.Min.Y) + float64(i*lineH) + float64(t.scrollOffsetY)
		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(tr.Min.X), ly)
		op.ColorScale.ScaleWithColor(textColor)
		text.Draw(screen, wl.text, *t.computedParams.Face, op)
	}

	// Render caret
	if t.focused {
		if t.widget.Disabled && t.computedParams.Color.DisabledCaret != nil {
			t.caret.Color = t.computedParams.Color.DisabledCaret
		} else {
			t.caret.Color = t.computedParams.Color.Caret
		}

		caretRect := tr
		caretRect = caretRect.Add(img.Point{cx, cy + t.scrollOffsetY})
		t.caret.SetLocation(caretRect)
		t.caret.Render(screen)
	}
}

func (t *TextInput) multilineCursorFromPoint(x, y int, tr img.Rectangle) int {
	lineH := int(math.Round(multilineLineHeight(t.computedParams.Face)))
	if lineH <= 0 {
		lineH = 1
	}
	maxWidth := t.contentWidth()
	wlines := wrapText(t.inputText, t.computedParams.Face, maxWidth)
	clickLine := (y - tr.Min.Y - t.scrollOffsetY) / lineH
	if clickLine < 0 {
		clickLine = 0
	}
	if clickLine >= len(wlines) {
		clickLine = len(wlines) - 1
	}
	col := fontStringIndex([]rune(wlines[clickLine].text), t.computedParams.Face, x-tr.Min.X)
	return wlines[clickLine].startRune + col
}

func (t *TextInput) GetText() string {
	return t.inputText
}

func (t *TextInput) SetText(text string) {
	t.setText(text, false)
}

func (t *TextInput) setJSText(text string) string {
	t.setText(text, true)
	return t.inputText
}

func (t *TextInput) setText(text string, isJS bool) {
	t.init.Do()
	t.DeselectText()
	t.inputText = text
	if t.validationFunc != nil {
		result, replacement := t.validationFunc(t.inputText)
		if !result {
			if replacement != nil {
				t.inputText = *replacement
			} else {
				t.inputText = t.lastInputText
			}
		}
	}
	if t.inputText != t.lastInputText {

		t.ChangedEvent.Fire(&TextInputChangedEventArgs{
			TextInput: t,
			InputText: t.inputText,
		})
		t.lastInputText = t.inputText

		if isJS {
			t.cursorPosition = jsUtil.GetCursorPosition()
		} else {
			t.CursorMoveEnd()
		}
	}
}

/** Focuser Interface - Start **/

func (t *TextInput) Focus(focused bool) {
	t.init.Do()
	t.GetWidget().FireFocusEvent(t, focused, img.Point{-1, -1})
	t.caret.resetBlinking()
	t.focused = focused

	if focused && runtime.GOOS == jsUtil.JS && runtime.GOARCH == jsUtil.WASM {
		jsUtil.Prompt(t.mobileInputMode, "Please enter a value.", t.inputText, t.cursorPosition, t.widget.Rect.Min.Y, t.setJSText, t.SelectAll)
	}
	if !focused {
		t.dragStartIndex = -1
	}
}

func (t *TextInput) IsFocused() bool {
	return t.focused
}

func (t *TextInput) TabOrder() int {
	return t.tabOrder
}

func (t *TextInput) GetFocus(direction FocusDirection) Focuser {
	return t.focusMap[direction]
}

func (t *TextInput) AddFocus(direction FocusDirection, focus Focuser) {
	t.focusMap[direction] = focus
}

/** Focuser Interface - End **/

func (t *TextInput) createWidget() {
	t.widget = NewWidget(append([]WidgetOpt{WidgetOpts.TrackHover(true)}, t.widgetOpts...)...)
	t.widget.focusable = t
	t.caret = NewCaret()
	t.mask = image.NewNineSliceColor(color.NRGBA{255, 0, 255, 255})
}

func (t *TextInput) initWidget() {
	_, height := text.Measure(" ", *t.computedParams.Face, 0)
	h := int(math.Round(height))

	t.caret.Color = t.computedParams.Color.Caret
	t.caret.Height = h
	t.caret.Width = *t.computedParams.CaretWidth
	t.caret.Validate()

	t.text = NewText(TextOpts.Text("", t.computedParams.Face, color.White))
	t.text.Validate()
}

func fontAdvance(s string, f *text.Face) int {
	a := text.Advance(s, *f)
	return int(math.Round(a))
}

// fontStringIndex returns an index into r that corresponds closest to pixel position x
// when string(r) is drawn using f. Pixel position x==0 corresponds to r[0].
func fontStringIndex(r []rune, f *text.Face, x int) int {
	start := 0
	end := len(r)
	p := 0
loop:
	for {
		p = start + (end-start)/2
		sub := string(r[:p])
		a := fontAdvance(sub, f)

		switch {
		// x is right of advance
		case x > a:
			if p == start {
				break loop
			}

			start = p

		// x is left of advance
		case x < a:
			if end == p {
				break loop
			}

			end = p

		// x matches advance exactly
		default:
			return p
		}
	}

	if len(r) > 0 {
		a1 := fontAdvance(string(r[:p]), f)
		a2 := fontAdvance(string(r[:p+1]), f)
		if math.Abs(float64(x-a2)) < math.Abs(float64(x-a1)) {
			p++
		}
	}

	return p
}
