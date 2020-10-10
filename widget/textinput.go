package widget

import (
	img "image"
	"image/color"
	"math"
	"sync/atomic"
	"time"

	"github.com/blizzy78/ebitenui/event"
	"github.com/blizzy78/ebitenui/image"
	"github.com/blizzy78/ebitenui/input"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

type TextInput struct {
	ChangedEvent *event.Event

	InputText string

	widgetOpts      []WidgetOpt
	caretOpts       []CaretOpt
	image           *TextInputImage
	color           *TextInputColor
	padding         Insets
	face            font.Face
	repeatDelay     time.Duration
	repeatInterval  time.Duration
	validationFunc  TextInputValidationFunc
	placeholderText string

	init           *MultiOnce
	commandToFunc  map[textInputControlCommand]textInputCommandFunc
	widget         *Widget
	caret          *Caret
	text           *Text
	renderBuf      *image.MaskedRenderBuffer
	mask           *image.NineSlice
	cursorPosition int
	state          textInputState
	scrollOffset   int
	focused        bool
	lastInputText  string
}

type TextInputOpt func(t *TextInput)

const TextInputOpts = textInputOpts(true)

type textInputOpts bool

type TextInputChangedEventArgs struct {
	TextInput *TextInput
	InputText string
}

type TextInputChangedHandlerFunc func(args *TextInputChangedEventArgs)

type TextInputImage struct {
	Idle     *image.NineSlice
	Disabled *image.NineSlice
}

type TextInputColor struct {
	Idle     color.Color
	Disabled color.Color
	Caret    color.Color
}

type TextInputValidationFunc func(newInputText string) bool

type textInputState func() (textInputState, bool)

type textInputControlCommand int

type textInputCommandFunc func()

const (
	textInputGoLeft = textInputControlCommand(iota + 1)
	textInputGoRight
	textInputGoStart
	textInputGoEnd
	textInputBackspace
	textInputDelete
)

var textInputKeyToCommand = map[ebiten.Key]textInputControlCommand{
	ebiten.KeyLeft:      textInputGoLeft,
	ebiten.KeyRight:     textInputGoRight,
	ebiten.KeyHome:      textInputGoStart,
	ebiten.KeyEnd:       textInputGoEnd,
	ebiten.KeyBackspace: textInputBackspace,
	ebiten.KeyDelete:    textInputDelete,
}

func NewTextInput(opts ...TextInputOpt) *TextInput {
	t := &TextInput{
		ChangedEvent: &event.Event{},

		repeatDelay:    300 * time.Millisecond,
		repeatInterval: 35 * time.Millisecond,

		init:          &MultiOnce{},
		commandToFunc: map[textInputControlCommand]textInputCommandFunc{},
		renderBuf:     image.NewMaskedRenderBuffer(),
	}
	t.state = t.idleState(true)

	t.commandToFunc[textInputGoLeft] = t.doGoLeft
	t.commandToFunc[textInputGoRight] = t.doGoRight
	t.commandToFunc[textInputGoStart] = t.doGoStart
	t.commandToFunc[textInputGoEnd] = t.doGoEnd
	t.commandToFunc[textInputBackspace] = t.doBackspace
	t.commandToFunc[textInputDelete] = t.doDelete

	t.init.Append(t.createWidget)

	for _, o := range opts {
		o(t)
	}

	return t
}

func (o textInputOpts) WidgetOpts(opts ...WidgetOpt) TextInputOpt {
	return func(t *TextInput) {
		t.widgetOpts = append(t.widgetOpts, opts...)
	}
}

func (o textInputOpts) CaretOpts(opts ...CaretOpt) TextInputOpt {
	return func(t *TextInput) {
		t.caretOpts = append(t.caretOpts, opts...)
	}
}

func (o textInputOpts) ChangedHandler(f TextInputChangedHandlerFunc) TextInputOpt {
	return func(t *TextInput) {
		t.ChangedEvent.AddHandler(func(args interface{}) {
			f(args.(*TextInputChangedEventArgs))
		})
	}
}

func (o textInputOpts) Image(i *TextInputImage) TextInputOpt {
	return func(t *TextInput) {
		t.image = i
	}
}

func (o textInputOpts) Color(c *TextInputColor) TextInputOpt {
	return func(t *TextInput) {
		t.color = c
	}
}

func (o textInputOpts) Padding(i Insets) TextInputOpt {
	return func(t *TextInput) {
		t.padding = i
	}
}

func (o textInputOpts) Face(f font.Face) TextInputOpt {
	return func(t *TextInput) {
		t.face = f
	}
}

func (o textInputOpts) RepeatInterval(i time.Duration) TextInputOpt {
	return func(t *TextInput) {
		t.repeatInterval = i
	}
}

func (o textInputOpts) Validation(f TextInputValidationFunc) TextInputOpt {
	return func(t *TextInput) {
		t.validationFunc = f
	}
}

func (o textInputOpts) Placeholder(s string) TextInputOpt {
	return func(t *TextInput) {
		t.placeholderText = s
	}
}

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
	return 50, h + t.padding.Top + t.padding.Bottom
}

func (t *TextInput) Render(screen *ebiten.Image, def DeferredRenderFunc) {
	t.init.Do()

	if t.InputText != t.lastInputText {
		t.ChangedEvent.Fire(&TextInputChangedEventArgs{
			TextInput: t,
			InputText: t.InputText,
		})
	}

	t.text.GetWidget().Disabled = t.widget.Disabled

	if t.cursorPosition > len([]rune(t.InputText)) {
		t.cursorPosition = len([]rune(t.InputText))
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

	t.widget.Render(screen, def)

	t.renderImage(screen)
	t.renderTextAndCaret(screen, def)
}

func (t *TextInput) idleState(newKeyOrCommand bool) textInputState {
	return func() (textInputState, bool) {
		if !t.focused {
			return t.idleState(true), false
		}

		chars := input.InputChars()
		if len(chars) > 0 {
			return t.charInputState(chars[0]), true
		}

		if st := textInputCheckForCommand(t, newKeyOrCommand); st != nil {
			return st, true
		}

		if input.MouseButtonJustPressedLayer(ebiten.MouseButtonLeft, t.widget.EffectiveInputLayer()) {
			t.doGoXY(input.CursorPosition())
		}

		return t.idleState(true), false
	}
}

func textInputCheckForCommand(t *TextInput, newKeyOrCommand bool) textInputState {
	for key, cmd := range textInputKeyToCommand {
		if !input.KeyPressed(key) {
			continue
		}

		var delay time.Duration
		if newKeyOrCommand {
			delay = t.repeatDelay
		} else {
			delay = t.repeatInterval
		}

		return t.commandState(cmd, key, delay, nil, nil)
	}

	return nil
}

func (t *TextInput) charInputState(c rune) textInputState {
	return func() (textInputState, bool) {
		if !t.widget.Disabled {
			t.doInsert(c)
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

		if timer != nil && expired.Load().(bool) {
			return t.idleState(false), true
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

func (t *TextInput) doInsert(c rune) {
	s := insertChar(t.InputText, c, t.cursorPosition)

	if t.validationFunc != nil && !t.validationFunc(s) {
		return
	}

	t.InputText = s
	t.cursorPosition++
}

func (t *TextInput) doGoLeft() {
	if t.cursorPosition > 0 {
		t.cursorPosition--
	}
	t.caret.ResetBlinking()
}

func (t *TextInput) doGoRight() {
	if t.cursorPosition < len([]rune(t.InputText)) {
		t.cursorPosition++
	}
	t.caret.ResetBlinking()
}

func (t *TextInput) doGoStart() {
	t.cursorPosition = 0
	t.caret.ResetBlinking()
}

func (t *TextInput) doGoEnd() {
	t.cursorPosition = len([]rune(t.InputText))
	t.caret.ResetBlinking()
}

func (t *TextInput) doGoXY(x int, y int) {
	p := img.Point{x, y}
	if p.In(t.widget.Rect) {
		tr := t.padding.Apply(t.widget.Rect)
		if x < tr.Min.X {
			x = tr.Min.X
		}
		if x > tr.Max.X {
			x = tr.Max.X
		}

		i := fontStringIndex(t.InputText, t.face, x-t.scrollOffset-tr.Min.X)
		t.cursorPosition = i
		t.caret.ResetBlinking()
	}
}

func (t *TextInput) doBackspace() {
	if t.cursorPosition > 0 {
		t.InputText = removeChar(t.InputText, t.cursorPosition-1)
		t.cursorPosition--
	}
	t.caret.ResetBlinking()
}

func (t *TextInput) doDelete() {
	if t.cursorPosition < len([]rune(t.InputText)) {
		t.InputText = removeChar(t.InputText, t.cursorPosition)
	}
	t.caret.ResetBlinking()
}

func insertChar(s string, r rune, pos int) string {
	b := string([]rune(s)[:pos])
	a := string([]rune(s)[pos:])
	return b + string(r) + a
}

func removeChar(s string, pos int) string {
	return string([]rune(s)[:pos]) + string([]rune(s)[pos+1:])
}

func (t *TextInput) renderImage(screen *ebiten.Image) {
	if t.image != nil {
		i := t.image.Idle
		if t.widget.Disabled && t.image.Disabled != nil {
			i = t.image.Disabled
		}

		rect := t.widget.Rect
		i.Draw(screen, rect.Dx(), rect.Dy(), func(opts *ebiten.DrawImageOptions) {
			opts.GeoM.Translate(float64(rect.Min.X), float64(rect.Min.Y))
		})
	}
}

func (t *TextInput) renderTextAndCaret(screen *ebiten.Image, def DeferredRenderFunc) {
	t.renderBuf.Draw(screen,
		func(buf *ebiten.Image) {
			t.drawTextAndCaret(buf, def)
		},
		func(buf *ebiten.Image) {
			rect := t.widget.Rect
			t.mask.Draw(buf, rect.Dx()-t.padding.Left-t.padding.Right, rect.Dy()-t.padding.Top-t.padding.Bottom,
				func(opts *ebiten.DrawImageOptions) {
					opts.GeoM.Translate(float64(rect.Min.X+t.padding.Left), float64(rect.Min.Y+t.padding.Top))
					opts.CompositeMode = ebiten.CompositeModeCopy
				})
		})
}

func (t *TextInput) drawTextAndCaret(screen *ebiten.Image, def DeferredRenderFunc) {
	rect := t.widget.Rect
	tr := rect
	tr = tr.Add(img.Point{t.padding.Left, t.padding.Top})

	cx := 0
	if t.focused {
		sub := string([]rune(t.InputText)[:t.cursorPosition])
		cx = fontAdvance(sub, t.face)

		dx := tr.Min.X + t.scrollOffset + cx + t.caret.Width + t.padding.Right - rect.Max.X
		if dx > 0 {
			t.scrollOffset -= dx
		}

		dx = tr.Min.X + t.scrollOffset + cx - t.padding.Left - rect.Min.X
		if dx < 0 {
			t.scrollOffset -= dx
		}
	}

	tr = tr.Add(img.Point{t.scrollOffset, 0})

	t.text.SetLocation(tr)
	if len([]rune(t.InputText)) > 0 {
		t.text.Label = t.InputText
	} else {
		t.text.Label = t.placeholderText
	}
	if t.widget.Disabled || len([]rune(t.InputText)) == 0 {
		t.text.Color = t.color.Disabled
	} else {
		t.text.Color = t.color.Idle
	}
	t.text.Render(screen, def)

	if t.focused {
		tr = tr.Add(img.Point{cx, 0})
		t.caret.SetLocation(tr)
		t.caret.Render(screen, def)
	}
}

func (t *TextInput) Focus(focused bool) {
	t.init.Do()
	WidgetFireFocusEvent(t.widget, focused)
	t.caret.resetBlinking()
	t.focused = focused
}

func (t *TextInput) createWidget() {
	t.widget = NewWidget(t.widgetOpts...)
	t.widgetOpts = nil

	t.caret = NewCaret(append(t.caretOpts, CaretOpts.Color(t.color.Caret))...)
	t.caretOpts = nil

	t.text = NewText(TextOpts.Text("", t.face, color.White))

	t.mask = image.NewNineSliceColor(color.RGBA{255, 0, 255, 255})
}

func fontAdvance(s string, f font.Face) int {
	_, a := font.BoundString(f, s)
	return int(math.Round(fixedInt26_6ToFloat64(a)))
}

func fontStringIndex(s string, f font.Face, x int) int {
	start := 0
	end := len(s)
	var p int
	for {
		p = start + (end-start)/2
		sub := string([]rune(s)[:p])
		a := fontAdvance(sub, f)

		if x-a == 0 {
			return p
		}

		if x-a > 0 {
			if p == start {
				break
			}

			start = p
		} else {
			if end == p {
				break
			}

			end = p
		}
	}

	if len(s) > 0 {
		sub := string([]rune(s)[:p])
		a1 := fontAdvance(sub, f)
		sub = string([]rune(s)[:p+1])
		a2 := fontAdvance(sub, f)
		if math.Abs(float64(x-a2)) < math.Abs(float64(x-a1)) {
			p++
		}
	}

	return p
}
