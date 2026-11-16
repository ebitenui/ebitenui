package widget

import (
	img "image"
	"strings"
	"testing"

	"github.com/ebitenui/ebitenui/input"
	internalinput "github.com/ebitenui/ebitenui/internal/input"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/exp/textinput"
	"github.com/matryer/is"
)

// imeStub stands in for the platform's IME: Update reports handled and
// counts its calls; ending a session records whether it committed.
type imeStub struct {
	handled bool
	updates int
	ends    []bool // true for a commit, false for a cancel
}

func stubIME(t *testing.T) *imeStub {
	t.Helper()
	s := &imeStub{}
	oldUpdate, oldEnd := imeUpdate, imeEnd
	imeUpdate = func(*textinput.Composer) (bool, error) { s.updates++; return s.handled, nil }
	imeEnd = func(_ *textinput.Composer, commit bool) { s.ends = append(s.ends, commit) }
	t.Cleanup(func() { imeUpdate, imeEnd = oldUpdate, oldEnd })
	return s
}

// placedInput returns an input taking its text through the IME, 200 by 30
// px at the origin.
func placedInput(t *testing.T, opts ...TextInputOpt) *TextInput {
	t.Helper()
	ti := newTextInput(t, append(opts, TextInputOpts.IME())...)
	ti.SetLocation(img.Rect(0, 0, 200, 30))
	return ti
}

// drawInput draws ti on a screen its size, without validating it again.
func drawInput(ti *TextInput) {
	ti.Render(ebiten.NewImage(200, 30))
}

func pressKey(t *testing.T, k ebiten.Key) {
	t.Helper()
	internalinput.InputHandler.KeyPressed[k] = true
	t.Cleanup(func() { internalinput.InputHandler.KeyPressed[k] = false })
}

func typeChars(t *testing.T, s string) {
	t.Helper()
	internalinput.InputHandler.InputChars = append(internalinput.InputHandler.InputChars[:0], []rune(s)...)
	t.Cleanup(func() { internalinput.InputHandler.InputChars = internalinput.InputHandler.InputChars[:0] })
}

// Committed text goes in at the caret, on one line.
func TestTextInput_IMECommit(t *testing.T) {
	is := is.New(t)

	ti := newTextInput(t, TextInputOpts.IME())
	ti.SetText("foo")
	ti.cursorPosition = 1

	ti.commitIME("ab", false, "", "")
	is.Equal(ti.GetText(), "faboo")
	is.Equal(ti.cursorPosition, 3) // after what was committed

	ti.commitIME("x\r\ny", false, "", "")
	is.Equal(ti.GetText(), "fabxyoo") // line breaks dropped
}

// When the IME rewrites the text around the caret, its version replaces the
// input's, the caret after what it committed, unless validation refuses it.
func TestTextInput_IMEReplacesSurroundings(t *testing.T) {
	is := is.New(t)

	ti := newTextInput(t, TextInputOpts.IME(), TextInputOpts.Validation(func(s string) (bool, *string) {
		return !strings.Contains(s, "!"), nil
	}))
	ti.SetText("abcd")
	ti.cursorPosition = 2

	ti.commitIME("X", true, "a", "d")
	is.Equal(ti.GetText(), "aXd")
	is.Equal(ti.cursorPosition, 2)

	ti.commitIME("!", true, "", "")
	is.Equal(ti.GetText(), "aXd") // refused
}

// The composition shows at the caret, with the caret inside it where the IME
// puts it, and the IME learns where the caret is and the text around it.
func TestTextInput_IMEComposition(t *testing.T) {
	is := is.New(t)

	ti := placedInput(t)
	ti.SetText("foo")
	ti.cursorPosition = 1
	ti.Focus(true)
	ti.compose("かな", len("か"))
	drawInput(ti)

	is.Equal(ti.text.Label, "fかなoo")
	is.Equal(ti.caretRect.Min.X-ti.text.GetWidget().Rect.Min.X, fontAdvanceAt("fかなoo", 2, ti.computedParams.Face)) // after "fか"
	is.True(ti.caretRect.Dy() > 0)

	o := ti.imeSession()
	is.Equal(o.CaretBounds, ti.caretRect)
	is.Equal(o.TextBeforeCaret, "f")
	is.Equal(o.TextAfterCaret, "oo")

	ti.compose("", 0)
	drawInput(ti)
	is.Equal(ti.text.Label, "foo") // gone when the IME clears it
}

// A secure input shows the composition as stars too, and keeps its text
// from the IME.
func TestTextInput_IMESecure(t *testing.T) {
	is := is.New(t)

	ti := placedInput(t, TextInputOpts.Secure(true))
	ti.SetText("pw")
	ti.cursorPosition = 1
	ti.Focus(true)
	ti.compose("ab", 0)
	drawInput(ti)

	is.Equal(ti.text.Label, "****")
	o := ti.imeSession()
	is.Equal(o.TextBeforeCaret+o.TextAfterCaret, "")
}

// Keys the IME takes leave the input alone; other keys end the IME's session
// first, then act. Typed characters come only through the IME.
func TestTextInput_IMEKeys(t *testing.T) {
	is := is.New(t)

	ti := newTextInput(t, TextInputOpts.IME())
	ti.SetText("foo")
	ti.Focus(true)
	ime := stubIME(t)
	pressKey(t, ebiten.KeyBackspace)
	typeChars(t, "x")

	ime.handled = true
	ti.Update(&UpdateObject{})
	is.Equal(ti.GetText(), "foo") // the IME took the key
	is.Equal(ime.updates, 1)

	ime.handled = false
	ti.Update(&UpdateObject{})
	is.Equal(ti.GetText(), "fo")     // Backspace, and no "x"
	is.Equal(ime.ends, []bool{true}) // committed first
}

// Pressing in the input commits the composition before the caret moves;
// losing focus commits it; a disabled input cancels it and asks the IME
// nothing; SetText cancels it.
func TestTextInput_IMESessionEnds(t *testing.T) {
	is := is.New(t)

	ti := placedInput(t)
	ti.SetText("foo")
	ti.Focus(true)
	ime := stubIME(t)
	r := ti.GetWidget().Rect
	input.SetupInputLayersWithDeferred(nil) // drops layers other tests left
	c := useTestCursor(t)
	c.x, c.y, c.just = r.Min.X+1, r.Min.Y+r.Dy()/2, true // a press in the input

	ti.Update(&UpdateObject{})
	is.Equal(ime.ends, []bool{true}) // before the caret moved
	is.Equal(ti.cursorPosition, 0)

	c.x, c.y, c.just = 0, 0, false
	ti.GetWidget().Disabled = true
	ti.Update(&UpdateObject{})
	is.Equal(ime.updates, 1)                // not asked while disabled
	is.Equal(ime.ends, []bool{true, false}) // cancelled

	ti.Focus(false)
	is.Equal(ime.ends, []bool{true, false, true}) // committed on losing focus

	ti.SetText("bar")
	is.Equal(ime.ends, []bool{true, false, true, false}) // replaced, so cancelled
}

// Dismissing the virtual keyboard ends the input's focus.
func TestTextInput_IMEEndByUser(t *testing.T) {
	is := is.New(t)

	ti := newTextInput(t, TextInputOpts.IME())
	ti.Focus(true)
	ti.composer.OnEndByUser()
	is.True(!ti.IsFocused())
}
