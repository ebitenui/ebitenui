package input

import (
	internalinput "github.com/blizzy78/ebitenui/internal/input"
	"github.com/hajimehoshi/ebiten"
)

// MouseButtonPressed returns whether mouse button b is currently pressed.
func MouseButtonPressed(b ebiten.MouseButton) bool {
	switch b {
	case ebiten.MouseButtonLeft:
		return internalinput.LeftMouseButtonPressed
	case ebiten.MouseButtonMiddle:
		return internalinput.MiddleMouseButtonPressed
	case ebiten.MouseButtonRight:
		return internalinput.RightMouseButtonPressed
	default:
		return false
	}
}

// MouseButtonJustPressed returns whether mouse button b has just been pressed.
// It only returns true during the first frame that the button is pressed.
func MouseButtonJustPressed(b ebiten.MouseButton) bool {
	switch b {
	case ebiten.MouseButtonLeft:
		return internalinput.LeftMouseButtonJustPressed
	case ebiten.MouseButtonMiddle:
		return internalinput.MiddleMouseButtonJustPressed
	case ebiten.MouseButtonRight:
		return internalinput.RightMouseButtonJustPressed
	default:
		return false
	}
}

// MouseButtonPressedLayer returns whether mouse button b is currently pressed if input layer l is
// eligible to handle it.
func MouseButtonPressedLayer(b ebiten.MouseButton, l *Layer) bool {
	if !MouseButtonPressed(b) {
		return false
	}

	x, y := CursorPosition()
	return l.ActiveFor(x, y, LayerEventTypeMouseButton)
}

// MouseButtonJustPressedLayer returns whether mouse button b has just been pressed if input layer l
// is eligible to handle it. It only returns true during the first frame that the button is pressed.
func MouseButtonJustPressedLayer(b ebiten.MouseButton, l *Layer) bool {
	if !MouseButtonJustPressed(b) {
		return false
	}

	x, y := CursorPosition()
	return l.ActiveFor(x, y, LayerEventTypeMouseButton)
}

// CursorPosition returns the current cursor position.
func CursorPosition() (int, int) {
	return internalinput.CursorX, internalinput.CursorY
}

// Wheel returns current mouse wheel movement.
func Wheel() (float64, float64) {
	return internalinput.WheelX, internalinput.WheelY
}

// WheelLayer returns current mouse wheel movement if input layer l is eligible to handle it.
// If l is not eligible, it returns 0, 0.
func WheelLayer(l *Layer) (float64, float64) {
	x, y := Wheel()
	if x == 0 && y == 0 {
		return 0, 0
	}

	cx, cy := CursorPosition()
	if !l.ActiveFor(cx, cy, LayerEventTypeWheel) {
		return 0, 0
	}

	return x, y
}

// InputChars returns user keyboard input.
func InputChars() []rune { //nolint:golint
	return internalinput.InputChars
}

// KeyPressed returns whether key k is currently pressed.
func KeyPressed(k ebiten.Key) bool {
	p, ok := internalinput.KeyPressed[k]
	return ok && p
}

// AnyKeyPressed returns whether any key is currently pressed.
func AnyKeyPressed() bool {
	return internalinput.AnyKeyPressed
}
