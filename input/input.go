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

func Update() {
	leftMouseButtonPressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	middleMouseButtonPressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle)
	rightMouseButtonPressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight)
	cursorX, cursorY = ebiten.CursorPosition()
	wheelX, wheelY = ebiten.Wheel()
}

func Draw() {
	leftMouseButtonJustPressed = leftMouseButtonPressed && leftMouseButtonPressed != lastLeftMouseButtonPressed
	middleMouseButtonJustPressed = middleMouseButtonPressed && middleMouseButtonPressed != lastMiddleMouseButtonPressed
	rightMouseButtonJustPressed = rightMouseButtonPressed && rightMouseButtonPressed != lastRightMouseButtonPressed

	lastLeftMouseButtonPressed = leftMouseButtonPressed
	lastMiddleMouseButtonPressed = middleMouseButtonPressed
	lastRightMouseButtonPressed = rightMouseButtonPressed
}

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

func MouseButtonPressedLayer(b ebiten.MouseButton, l *Layer) bool {
	if !MouseButtonPressed(b) {
		return false
	}

	x, y := CursorPosition()
	return l.ActiveFor(x, y, LayerEventTypeMouseButton)
}

func MouseButtonJustPressedLayer(b ebiten.MouseButton, l *Layer) bool {
	if !MouseButtonJustPressed(b) {
		return false
	}

	x, y := CursorPosition()
	return l.ActiveFor(x, y, LayerEventTypeMouseButton)
}

func CursorPosition() (int, int) {
	return cursorX, cursorY
}

func Wheel() (float64, float64) {
	return wheelX, wheelY
}

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
