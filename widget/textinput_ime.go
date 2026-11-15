package widget

import (
	"strings"
	"unicode/utf8"

	"github.com/hajimehoshi/ebiten/v2/exp/textinput"
)

// imeUpdate and imeEnd reach the platform's text input; tests stand in for
// them.
var (
	imeUpdate = (*textinput.Composer).Update
	imeEnd    = func(c *textinput.Composer, commit bool) {
		if commit {
			c.Confirm()
		} else {
			c.Cancel()
		}
	}
)

// IME has the input take its text through the platform's text input
// (Ebitengine's exp/textinput) instead of typed characters. An IME composes
// in place: its text shows highlighted at the caret, its candidate window
// by the caret. On iOS and Android a virtual keyboard shows while the input
// has focus, the screen shifting so it doesn't cover the caret, and
// dismissing it takes the focus away; in browsers it replaces the prompt.
func (o TextInputOptions) IME() TextInputOpt {
	return func(t *TextInput) {
		t.ime = true
	}
}

// shownKey is what the text as drawn depends on.
type shownKey struct {
	text, composition string
	cursor            int
	secure            bool
}

// initIME hooks the input to its composer.
func (t *TextInput) initIME() {
	t.composer.OnNewSession = t.imeSession
	t.composer.OnComposition = func(c *textinput.Composition) {
		start, _ := c.SelectionRangeInBytes()
		t.compose(c.Text(), start)
	}
	t.composer.OnCommit = func(c *textinput.Commit) {
		replacedBefore, replacedAfter := c.IsSurroundingTextReplaced()
		before, after := c.SurroundingText()
		t.commitIME(c.Text(), replacedBefore || replacedAfter, before, after)
	}
	t.composer.OnEndByUser = func() { t.Focus(false) }
}

// updateIME runs the IME for a tick while the input has focus and reports
// whether it took the tick's keys. A disabled input cancels its session
// instead. An error from the platform has dropped the session; the next
// tick starts another.
func (t *TextInput) updateIME() bool {
	if t.widget.Disabled {
		imeEnd(&t.composer, false)
		return false
	}
	handled, _ := imeUpdate(&t.composer)
	return handled
}

// endIME ends the IME's session, committing its composition or not, before
// the text or the caret changes from outside the IME.
func (t *TextInput) endIME(commit bool) {
	if t.ime {
		imeEnd(&t.composer, commit)
	}
}

// imeSession tells the IME where the caret is and the text around it; a
// secure input keeps its text to itself.
func (t *TextInput) imeSession() *textinput.SessionOptions {
	o := &textinput.SessionOptions{CaretBounds: t.caretRect}
	if !*t.computedParams.Secure {
		i := runeOffset(t.inputText, t.cursorPosition)
		o.TextBeforeCaret, o.TextAfterCaret = t.inputText[:i], t.inputText[i:]
	}
	return o
}

// compose keeps the IME's composition, shown at the caret, and the IME's
// caret within it, in bytes.
func (t *TextInput) compose(text string, caret int) {
	t.composition, t.compositionCaret = text, caret
}

// commitIME applies text the IME committed, on one line: inserted at the
// caret, or, when the IME rewrote the text around the caret, with what it
// left before and after in place of the input's text, validated the same
// way.
func (t *TextInput) commitIME(text string, replaced bool, before, after string) {
	if strings.ContainsAny(text, "\r\n") {
		text = strings.Map(func(r rune) rune {
			if r == '\r' || r == '\n' {
				return -1
			}
			return r
		}, text)
	}
	if !replaced {
		t.Insert(text)
		return
	}
	s := before + text + after
	if t.validationFunc != nil {
		if ok, replacement := t.validationFunc(s); !ok {
			if replacement == nil {
				return
			}
			s = *replacement
		}
	}
	t.DeselectText()
	t.inputText = s
	t.cursorPosition = min(utf8.RuneCountInString(before)+utf8.RuneCountInString(text), utf8.RuneCountInString(s))
}

// shown returns the text as drawn: the composition at the caret, and stars
// for a secure input. It makes a new string only when one of those changed.
func (t *TextInput) shown() string {
	k := shownKey{t.inputText, t.composition, t.cursorPosition, *t.computedParams.Secure}
	if k != t.shownFor {
		t.shownFor = k
		s := t.inputText
		if t.composition != "" {
			i := runeOffset(s, t.cursorPosition)
			s = s[:i] + t.composition + s[i:]
		}
		if k.secure {
			s = strings.Repeat("*", utf8.RuneCountInString(s))
		}
		t.shownText = s
	}
	return t.shownText
}

// shownCaret returns the caret's rune in the text as drawn: inside the
// composition where the IME puts it.
func (t *TextInput) shownCaret() int {
	if t.composition == "" {
		return t.cursorPosition
	}
	return t.cursorPosition + utf8.RuneCountInString(t.composition[:min(t.compositionCaret, len(t.composition))])
}
