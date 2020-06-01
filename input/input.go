package input

import "github.com/hajimehoshi/ebiten"

var (
	leftMouseButtonPressed   bool
	middleMouseButtonPressed bool
	rightMouseButtonPressed  bool
	cursorX                  int
	cursorY                  int
	wheelX                   float64
	wheelY                   float64

	leftMouseButtonJustPressed   bool
	middleMouseButtonJustPressed bool
	rightMouseButtonJustPressed  bool

	lastLeftMouseButtonPressed   bool
	lastMiddleMouseButtonPressed bool
	lastRightMouseButtonPressed  bool
)

// Update updates the input system. This function should not be called directly.
func Update() {
	leftMouseButtonPressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	middleMouseButtonPressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle)
	rightMouseButtonPressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight)
	cursorX, cursorY = ebiten.CursorPosition()
	wheelX, wheelY = ebiten.Wheel()
}

// Draw updates the input system. This function should not be called directly.
func Draw() {
	leftMouseButtonJustPressed = leftMouseButtonPressed && leftMouseButtonPressed != lastLeftMouseButtonPressed
	middleMouseButtonJustPressed = middleMouseButtonPressed && middleMouseButtonPressed != lastMiddleMouseButtonPressed
	rightMouseButtonJustPressed = rightMouseButtonPressed && rightMouseButtonPressed != lastRightMouseButtonPressed

	lastLeftMouseButtonPressed = leftMouseButtonPressed
	lastMiddleMouseButtonPressed = middleMouseButtonPressed
	lastRightMouseButtonPressed = rightMouseButtonPressed
}

// MouseButtonPressed returns whether mouse button b is currently pressed.
func MouseButtonPressed(b ebiten.MouseButton) bool {
	switch b {
	case ebiten.MouseButtonLeft:
		return leftMouseButtonPressed
	case ebiten.MouseButtonMiddle:
		return middleMouseButtonPressed
	case ebiten.MouseButtonRight:
		return rightMouseButtonPressed
	default:
		return false
	}
}

// MouseButtonJustPressed returns whether mouse button b has just been pressed.
// It only returns true during the first frame that the button is pressed.
func MouseButtonJustPressed(b ebiten.MouseButton) bool {
	switch b {
	case ebiten.MouseButtonLeft:
		return leftMouseButtonJustPressed
	case ebiten.MouseButtonMiddle:
		return middleMouseButtonJustPressed
	case ebiten.MouseButtonRight:
		return rightMouseButtonJustPressed
	default:
		return false
	}
}

// MouseButtonPressedLayer returns whether mouse button b is currently pressed and input layer l is
// eligible to handle it.
func MouseButtonPressedLayer(b ebiten.MouseButton, l *Layer) bool {
	if !MouseButtonPressed(b) {
		return false
	}

	x, y := CursorPosition()
	return l.ActiveFor(x, y, LayerEventTypeMouseButton)
}

// MouseButtonJustPressedLayer returns whether mouse button b has just been pressed and input layer l
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
	return cursorX, cursorY
}

// Wheel returns current mouse wheel movement.
func Wheel() (float64, float64) {
	return wheelX, wheelY
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
