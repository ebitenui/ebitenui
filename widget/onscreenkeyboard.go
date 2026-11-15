package widget

import (
	"unicode"
)

// keyboardRows are an OnScreenKeyboard's character keys, row by row, lower
// case.
var keyboardRows = [...]string{
	"1234567890",
	"qwertyuiop",
	"asdfghjkl'",
	"zxcvbnm,.-",
}

// OnScreenKeyboard is a keyboard of buttons that types into a TextInput, for
// a gamepad on a platform that shows no keyboard of its own. Its keys take
// the focus as buttons do, so the D-pad or the stick moves between them and
// A presses the focused one. Shift switches the letters between lower and
// upper case, Space types a space, Erase deletes the character before the
// caret, and Done submits the input and closes the keyboard. A line above
// the keys shows the input's text, as stars for a secure input.
// ebitenui.UI.OnScreenKeyboard opens it in a window for a focused text
// input, and closes it on Back.
type OnScreenKeyboard struct {
	container *Container
	keyOpts   []ButtonOpt
	preview   *Text
	letters   []*Button
	first     *Button
	upper     bool
	target    *TextInput
	onClose   func()
}

// OnScreenKeyboardOpt configures an OnScreenKeyboard.
type OnScreenKeyboardOpt func(k *onScreenKeyboardConfig)

type onScreenKeyboardConfig struct {
	containerOpts []ContainerOpt
	keyOpts       []ButtonOpt
	previewOpts   []TextOpt
}

type OnScreenKeyboardOptions struct{}

// OnScreenKeyboardOpts holds the options of NewOnScreenKeyboard.
var OnScreenKeyboardOpts OnScreenKeyboardOptions

// ContainerOpts configures the keyboard's container: its background,
// padding, size.
func (o OnScreenKeyboardOptions) ContainerOpts(opts ...ContainerOpt) OnScreenKeyboardOpt {
	return func(k *onScreenKeyboardConfig) { k.containerOpts = append(k.containerOpts, opts...) }
}

// KeyOpts configures every key, a Button: its images, text face and colors,
// size. Without them, the keys take the theme's.
func (o OnScreenKeyboardOptions) KeyOpts(opts ...ButtonOpt) OnScreenKeyboardOpt {
	return func(k *onScreenKeyboardConfig) { k.keyOpts = append(k.keyOpts, opts...) }
}

// PreviewOpts configures the line above the keys that shows the input's
// text: its face and color, without a theme.
func (o OnScreenKeyboardOptions) PreviewOpts(opts ...TextOpt) OnScreenKeyboardOpt {
	return func(k *onScreenKeyboardConfig) { k.previewOpts = append(k.previewOpts, opts...) }
}

// NewOnScreenKeyboard returns a keyboard; GetContainer is the widget to
// show, Open makes it type into an input.
func NewOnScreenKeyboard(opts ...OnScreenKeyboardOpt) *OnScreenKeyboard {
	// Small keys by default, so ten fit across a narrow screen; the grid
	// stretches them across the keyboard's width.
	cfg := onScreenKeyboardConfig{keyOpts: []ButtonOpt{ButtonOpts.TextPadding(&Insets{Left: 8, Right: 8, Top: 6, Bottom: 6})}}
	for _, o := range opts {
		o(&cfg)
	}
	k := &OnScreenKeyboard{keyOpts: cfg.keyOpts}
	k.container = NewContainer(append([]ContainerOpt{ContainerOpts.Layout(NewRowLayout(
		RowLayoutOpts.Direction(DirectionVertical), RowLayoutOpts.Spacing(6),
		RowLayoutOpts.Padding(NewInsetsSimple(8))))}, cfg.containerOpts...)...)
	stretch := WidgetOpts.LayoutData(RowLayoutData{Stretch: true})
	k.preview = NewText(append(cfg.previewOpts, TextOpts.WidgetOpts(stretch))...)
	k.preview.Label = " " // a line's height before it shows a text

	columns := len(keyboardRows[0])
	chars := NewContainer(ContainerOpts.WidgetOpts(stretch), ContainerOpts.Layout(NewGridLayout(
		GridLayoutOpts.Columns(columns), GridLayoutOpts.Spacing(4, 4), GridLayoutOpts.Stretch(trueSlice(columns), nil))))
	for _, row := range keyboardRows {
		for _, r := range row {
			b := k.key(string(r), func() { k.typeRune(r) })
			if unicode.IsLetter(r) {
				k.letters = append(k.letters, b)
			}
			if k.first == nil {
				k.first = b
			}
			chars.AddChild(b)
		}
	}
	commands := NewContainer(ContainerOpts.WidgetOpts(stretch), ContainerOpts.Layout(NewGridLayout(
		GridLayoutOpts.Columns(4), GridLayoutOpts.Spacing(4, 4), GridLayoutOpts.Stretch(trueSlice(4), nil))))
	commands.AddChild(
		k.key("Shift", k.shift),
		k.key("Space", func() { k.typeRune(' ') }),
		k.key("Erase", k.erase),
		k.key("Done", k.done),
	)
	k.container.AddChild(k.preview, chars, commands)
	return k
}

func trueSlice(n int) []bool {
	s := make([]bool, n)
	for i := range s {
		s[i] = true
	}
	return s
}

// key returns a key labelled label that calls press when clicked.
func (k *OnScreenKeyboard) key(label string, press func()) *Button {
	return NewButton(append(k.keyOpts[:len(k.keyOpts):len(k.keyOpts)],
		ButtonOpts.TextLabel(label),
		ButtonOpts.ClickedHandler(func(*ButtonClickedEventArgs) { press() }))...)
}

// GetContainer returns the keyboard's widget: the preview line and the
// keys.
func (k *OnScreenKeyboard) GetContainer() *Container { return k.container }

// Open makes the keyboard type into target until Done or Close, which call
// onClose; ebitenui.UI opens it in a window and closes that window there.
func (k *OnScreenKeyboard) Open(target *TextInput, onClose func()) {
	k.target, k.onClose = target, onClose
	k.showText()
}

// IsOpen reports whether the keyboard types into an input.
func (k *OnScreenKeyboard) IsOpen() bool { return k.target != nil }

// Close ends the typing: the keyboard forgets its input and calls the close
// function Open was given.
func (k *OnScreenKeyboard) Close() {
	onClose := k.onClose
	k.target, k.onClose = nil, nil
	if onClose != nil {
		onClose()
	}
}

// FirstKey returns the key that takes the focus as the keyboard opens.
func (k *OnScreenKeyboard) FirstKey() Focuser { return k.first }

// showText shows the input's text above the keys, as the input shows it; a
// space while it is empty keeps the line's height.
func (k *OnScreenKeyboard) showText() {
	if k.target == nil {
		return
	}
	k.preview.Label = k.target.shown()
	if k.preview.Label == "" {
		k.preview.Label = " "
	}
}

func (k *OnScreenKeyboard) typeRune(r rune) {
	if k.target == nil {
		return
	}
	if k.upper {
		r = unicode.ToUpper(r)
	}
	k.target.Insert(string(r))
	k.showText()
}

func (k *OnScreenKeyboard) erase() {
	if k.target == nil {
		return
	}
	k.target.Backspace()
	k.showText()
}

// shift switches the letters between lower and upper case.
func (k *OnScreenKeyboard) shift() {
	k.upper = !k.upper
	for _, b := range k.letters {
		r := []rune(b.textLabel)[0]
		if k.upper {
			b.SetText(string(unicode.ToUpper(r)))
		} else {
			b.SetText(string(unicode.ToLower(r)))
		}
	}
}

// done submits the input and closes the keyboard.
func (k *OnScreenKeyboard) done() {
	t := k.target
	k.Close()
	if t != nil {
		t.Submit()
	}
}
