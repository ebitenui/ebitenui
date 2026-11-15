package widget

import (
	"image"
	"image/color"
	"testing"

	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/input"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/matryer/is"
)

func TestColorPicker_HSVRoundTrip(t *testing.T) {
	is := is.New(t)

	for _, c := range []color.NRGBA{
		{255, 0, 0, 255}, {0, 255, 0, 255}, {0, 0, 255, 255}, {128, 128, 128, 255},
		{255, 255, 255, 255}, {0, 0, 0, 255}, {12, 200, 90, 255}, {255, 128, 0, 255},
	} {
		is.Equal(hsvToRGB(rgbToHSV(c)), c) // the color back from its hue, saturation and value
	}
}

// Pressing the wheel picks the hue and saturation under the pointer: red at
// the right of the rim, white in the middle.
func TestColorPicker_PickOnWheel(t *testing.T) {
	is := is.New(t)

	p, changes := newColorPicker(t, ColorPickerOpts.InitialColor(color.Black))
	for _, c := range []struct {
		x, y int
		want color.NRGBA
	}{
		{99, 50, color.NRGBA{255, 0, 0, 255}},
		{50, 50, color.NRGBA{255, 255, 255, 255}},
	} {
		pressWheel(p, c.x, c.y)
		is.True(near(p.Color(), c.want))                       // picked
		is.Equal((*changes)[len(*changes)-1].Color, p.Color()) // reported
		is.True((*changes)[len(*changes)-1].Dragging)          // while the button is down
	}
}

// The slider sets the brightness and keeps the hue and saturation.
func TestColorPicker_Brightness(t *testing.T) {
	is := is.New(t)

	p, changes := newColorPicker(t, ColorPickerOpts.InitialColor(color.NRGBA{255, 128, 0, 255}))
	p.slider.Current = 128
	render(p, t)

	is.Equal(p.Color(), color.NRGBA{128, 64, 0, 255}) // half as bright, the same hue
	is.Equal(len(*changes), 1)
	is.Equal((*changes)[0].Color, p.Color()) // reported
}

// SetColor changes the color and the slider without reporting a change.
func TestColorPicker_SetColor(t *testing.T) {
	is := is.New(t)

	p, changes := newColorPicker(t)
	p.SetColor(color.NRGBA{0, 0, 200, 255})
	render(p, t)

	is.Equal(p.Color(), color.NRGBA{0, 0, 200, 255})
	is.Equal(p.slider.Current, 200)
	is.Equal(len(*changes), 0)
}

// Dragging across the wheel follows the pointer, and releasing the button
// reports the last color, no longer dragging.
func TestColorPicker_Drag(t *testing.T) {
	is := is.New(t)

	p, changes := newColorPicker(t)
	pressWheel(p, 50, 50)

	cursor := &pickerCursor{x: 99, y: 50, pressed: true}
	input.SetCursorUpdater(cursor)
	t.Cleanup(func() { input.SetCursorUpdater(nil) })
	p.Update(&UpdateObject{})
	event.ExecuteDeferred()
	is.True(near(p.Color(), color.NRGBA{255, 0, 0, 255})) // followed to the rim

	cursor.pressed = false
	p.Update(&UpdateObject{})
	event.ExecuteDeferred()
	last := (*changes)[len(*changes)-1]
	is.True(!last.Dragging) // released
	is.Equal(last.Color, p.Color())
}

// Pressing a swatch picks its color, done at once, and the swatch of the
// color picked is the one framed.
func TestColorPicker_Swatches(t *testing.T) {
	is := is.New(t)

	red, green := color.NRGBA{255, 0, 0, 255}, color.NRGBA{0, 255, 0, 255}
	p, changes := newColorPicker(t, ColorPickerOpts.InitialColor(red), ColorPickerOpts.Swatches(20, red, green))
	swatches := p.container.children[2].(*Container).children                        //nolint:forcetypeassert // the swatch row
	redSwatch, greenSwatch := swatches[0].(*colorSwatch), swatches[1].(*colorSwatch) //nolint:forcetypeassert // swatches
	is.True(redSwatch.chosen())

	greenSwatch.widget.MouseButtonPressedEvent.Fire(&WidgetMouseButtonPressedEventArgs{Widget: greenSwatch.widget, Button: ebiten.MouseButtonLeft})
	event.ExecuteDeferred()
	render(p, t)

	is.Equal(p.Color(), green)
	is.Equal(len(*changes), 1)       // reported once, the slider's move included
	is.True(!(*changes)[0].Dragging) // done at once
	is.True(greenSwatch.chosen() && !redSwatch.chosen())
}

func newColorPicker(t *testing.T, opts ...ColorPickerOpt) (*ColorPicker, *[]ColorPickerChangedEventArgs) {
	t.Helper()

	var changes []ColorPickerChangedEventArgs
	ns := newNineSliceEmpty(t)
	p := NewColorPicker(append([]ColorPickerOpt{
		ColorPickerOpts.WheelSize(100),
		ColorPickerOpts.SliderHandle(&ButtonImage{Idle: ns, Pressed: ns}, 10),
		ColorPickerOpts.ChangedHandler(func(args *ColorPickerChangedEventArgs) { changes = append(changes, *args) }),
	}, opts...)...)
	p.Validate()
	p.SetLocation(image.Rect(0, 0, 100, 130))
	render(p, t)
	return p, &changes
}

// pressWheel presses the left button on the wheel, x and y from its top left.
func pressWheel(p *ColorPicker, x, y int) {
	p.wheel.widget.MouseButtonPressedEvent.Fire(&WidgetMouseButtonPressedEventArgs{
		Widget: p.wheel.widget, Button: ebiten.MouseButtonLeft, OffsetX: x, OffsetY: y,
	})
	event.ExecuteDeferred()
}

// near reports whether two colors are within a few levels in each channel.
func near(a, b color.NRGBA) bool {
	d := func(x, y uint8) bool { return int(x)-int(y) <= 6 && int(y)-int(x) <= 6 }
	return d(a.R, b.R) && d(a.G, b.G) && d(a.B, b.B) && a.A == b.A
}

// pickerCursor is a pointer at one place, its left button down or up.
type pickerCursor struct {
	x, y    int
	pressed bool
}

func (c *pickerCursor) Update()                 {}
func (c *pickerCursor) AfterUpdate()            {}
func (c *pickerCursor) Draw(*ebiten.Image)      {}
func (c *pickerCursor) AfterDraw(*ebiten.Image) {}
func (c *pickerCursor) MouseButtonPressed(b ebiten.MouseButton) bool {
	return c.pressed && b == ebiten.MouseButtonLeft
}
func (c *pickerCursor) MouseButtonJustPressed(ebiten.MouseButton) bool  { return false }
func (c *pickerCursor) MouseButtonJustReleased(ebiten.MouseButton) bool { return false }
func (c *pickerCursor) CursorPosition() (int, int)                      { return c.x, c.y }
func (c *pickerCursor) GetCursorImage(string) *ebiten.Image             { return nil }
func (c *pickerCursor) GetCursorOffset(string) image.Point              { return image.Point{} }

// With a radius, swatches draw rounded layers: the frames the picker
// shares, rounded by the radius and 2 pixels less, and each its color,
// rounded by 4 pixels less; without one, none.
func TestColorPicker_SwatchRadius(t *testing.T) {
	is := is.New(t)

	rounded, _ := newColorPicker(t, ColorPickerOpts.Swatches(24, color.White), ColorPickerOpts.SwatchRadius(6))
	w, _ := rounded.swatchFrames[0][0].MinSize()
	is.Equal(w, 12) // the outer frame: radius 6
	w, _ = rounded.swatchFrames[1][1].MinSize()
	is.Equal(w, 8) // the inner frame: radius 4
	sw := swatchesOf(rounded)
	w, _ = sw[0].fill.MinSize()
	is.Equal(w, 4) // the color: radius 2

	square, _ := newColorPicker(t, ColorPickerOpts.Swatches(24, color.White))
	is.True(swatchesOf(square)[0].fill == nil)
}

// swatchesOf returns a picker's swatches.
func swatchesOf(p *ColorPicker) []*colorSwatch {
	var out []*colorSwatch
	for _, ch := range p.container.Children() {
		if row, ok := ch.(*Container); ok {
			for _, s := range row.Children() {
				if s, ok := s.(*colorSwatch); ok {
					out = append(out, s)
				}
			}
		}
	}
	return out
}
