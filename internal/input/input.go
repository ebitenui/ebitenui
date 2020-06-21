package input

import "github.com/hajimehoshi/ebiten"

var (
	LeftMouseButtonPressed   bool
	MiddleMouseButtonPressed bool
	RightMouseButtonPressed  bool
	CursorX                  int
	CursorY                  int
	WheelX                   float64
	WheelY                   float64

	LeftMouseButtonJustPressed   bool
	MiddleMouseButtonJustPressed bool
	RightMouseButtonJustPressed  bool

	LastLeftMouseButtonPressed   bool
	LastMiddleMouseButtonPressed bool
	LastRightMouseButtonPressed  bool
)

// Update updates the input system.
func Update() {
	LeftMouseButtonPressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	MiddleMouseButtonPressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle)
	RightMouseButtonPressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight)
	CursorX, CursorY = ebiten.CursorPosition()
	WheelX, WheelY = ebiten.Wheel()
}

// Draw updates the input system.
func Draw() {
	LeftMouseButtonJustPressed = LeftMouseButtonPressed && LeftMouseButtonPressed != LastLeftMouseButtonPressed
	MiddleMouseButtonJustPressed = MiddleMouseButtonPressed && MiddleMouseButtonPressed != LastMiddleMouseButtonPressed
	RightMouseButtonJustPressed = RightMouseButtonPressed && RightMouseButtonPressed != LastRightMouseButtonPressed

	LastLeftMouseButtonPressed = LeftMouseButtonPressed
	LastMiddleMouseButtonPressed = MiddleMouseButtonPressed
	LastRightMouseButtonPressed = RightMouseButtonPressed
}
