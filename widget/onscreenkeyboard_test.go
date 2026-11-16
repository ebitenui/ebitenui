package widget

import (
	"image/color"
	"testing"

	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/image"
	"github.com/matryer/is"
)

// newKeyboard returns a keyboard whose keys and preview have images and a
// face of their own, with opts, laid out.
func newKeyboard(t *testing.T, opts ...OnScreenKeyboardOpt) *OnScreenKeyboard {
	t.Helper()
	ns := newNineSliceEmpty(t)
	k := NewOnScreenKeyboard(append([]OnScreenKeyboardOpt{
		OnScreenKeyboardOpts.KeyOpts(
			ButtonOpts.Image(&ButtonImage{Idle: ns, Pressed: ns}),
			ButtonOpts.TextFace(loadFont(t)),
			ButtonOpts.TextColor(&ButtonTextColor{Idle: color.White}),
		),
		OnScreenKeyboardOpts.PreviewOpts(TextOpts.Text("", loadFont(t), color.White)),
	}, opts...)...)
	render(k.GetContainer(), t)
	return k
}

// The caller's key options win over the keyboard's own: wider padding
// makes wider keys.
func TestOnScreenKeyboard_KeyOptsWin(t *testing.T) {
	is := is.New(t)

	narrow, _ := key(t, newKeyboard(t), "1").PreferredSize()
	wide, _ := key(t, newKeyboard(t, OnScreenKeyboardOpts.KeyOpts(ButtonOpts.TextPadding(&Insets{Left: 30, Right: 30}))), "1").PreferredSize()
	is.True(wide >= narrow+40)
}

// withChildren is a container of any type.
type withChildren interface {
	Children() []PreferredSizeLocateableWidget
}

// key returns the keyboard's key labelled label.
func key(t *testing.T, k *OnScreenKeyboard, label string) *Button {
	t.Helper()
	var found *Button
	var walk func(c withChildren)
	walk = func(c withChildren) {
		for _, ch := range c.Children() {
			switch w := ch.(type) {
			case *Button:
				if w.textLabel == label {
					found = w
				}
			case *Container:
				walk(w)
			}
		}
	}
	walk(k.GetContainer())
	if found == nil {
		t.Fatalf("no key %q", label)
	}
	return found
}

// press clicks the keys labelled labels, one after the other.
func press(t *testing.T, k *OnScreenKeyboard, labels ...string) {
	t.Helper()
	for _, l := range labels {
		key(t, k, l).Click()
		event.ExecuteDeferred()
	}
}

// The keys type into the input the keyboard was opened for: characters,
// Shift switching the letters' case, Space, and Erase deleting the
// character before the caret. The line above the keys shows the input's
// text as the input does, a secure one as stars, from the moment it opens.
func TestOnScreenKeyboard_Types(t *testing.T) {
	is := is.New(t)

	k := newKeyboard(t)
	ti := newTextInput(t)
	ti.SetText("z")
	k.Open(ti, func() {})
	is.Equal(k.preview.Label, "z")
	press(t, k, "a", "Shift", "B", "1", "Space", "Shift", "c", "Erase")
	is.Equal(ti.GetText(), "zaB1 ")
	is.Equal(k.preview.Label, "zaB1 ")
	is.Equal(key(t, k, "q").textLabel, "q") // lower case again

	secure := newTextInput(t, TextInputOpts.Secure(true))
	k.Open(secure, func() {})
	_, empty := k.GetContainer().PreferredSize()
	press(t, k, "x", "y")
	is.Equal(k.preview.Label, "**")
	_, typed := k.GetContainer().PreferredSize()
	is.Equal(typed, empty) // as tall empty as typed in, so its window stays right
}

// Done submits the input and closes the keyboard, which then types
// nothing; Close closes it without submitting.
func TestOnScreenKeyboard_DoneAndClose(t *testing.T) {
	is := is.New(t)

	k := newKeyboard(t)
	submits := 0
	ti := newTextInput(t, TextInputOpts.SubmitHandler(func(*TextInputChangedEventArgs) { submits++ }))
	closes := 0
	k.Open(ti, func() { closes++ })
	press(t, k, "a", "Done")
	is.Equal(closes, 1)
	is.Equal(submits, 1)
	is.True(!k.IsOpen())
	press(t, k, "b", "Erase")
	is.Equal(ti.GetText(), "a") // closed: types nothing

	k.Open(ti, func() { closes++ })
	is.True(k.IsOpen())
	k.Close()
	event.ExecuteDeferred()
	is.Equal(closes, 2)
	is.Equal(submits, 1) // not submitted
}

// Without a background of its own, the keyboard takes its theme's panel
// background and follows the theme when it changes; given one, it keeps it;
// without either, it has none.
func TestOnScreenKeyboard_ThemeBackground(t *testing.T) {
	is := is.New(t)

	dark := &Theme{PanelTheme: &PanelParams{BackgroundImage: newNineSliceEmpty(t)}}
	light := &Theme{PanelTheme: &PanelParams{BackgroundImage: newNineSliceEmpty(t)}}
	background := func(k *OnScreenKeyboard) any {
		k.GetContainer().Validate()
		return k.container.computedParams.BackgroundImage
	}

	k := newKeyboard(t)
	k.GetContainer().GetWidget().SetTheme(dark)
	is.True(background(k) == dark.PanelTheme.BackgroundImage)
	k.GetContainer().GetWidget().SetTheme(light)
	is.True(background(k) == light.PanelTheme.BackgroundImage)

	own := newNineSliceEmpty(t)
	k = newKeyboard(t, OnScreenKeyboardOpts.ContainerOpts(ContainerOpts.BackgroundImage(own)))
	k.GetContainer().GetWidget().SetTheme(dark)
	is.True(background(k) == own)

	k = newKeyboard(t)
	is.True(background(k) == (*image.NineSlice)(nil))
}
