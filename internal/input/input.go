package input

import (
	"github.com/hajimehoshi/ebiten/v2"
)

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

	InputChars    []rune
	KeyPressed    = map[ebiten.Key]bool{}
	AnyKeyPressed bool
	isTouched     bool
)

// Update updates the input system. This is called by the UI.
func Update() {
	touches := ebiten.TouchIDs()
	if len(touches) > 0 {
		isTouched = true
	}
	if isTouched {
		if len(touches) > 0 {
			LeftMouseButtonPressed = true
			CursorX, CursorY = ebiten.TouchPosition(touches[0])
		} else {
			LeftMouseButtonPressed = false
			isTouched = false
		}
	} else {
		LeftMouseButtonPressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
		MiddleMouseButtonPressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle)
		RightMouseButtonPressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight)
		CursorX, CursorY = ebiten.CursorPosition()
	}

	wx, wy := ebiten.Wheel()
	WheelX += wx
	WheelY += wy

	InputChars = ebiten.AppendInputChars(InputChars)
	AnyKeyPressed = false
	for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
		p := ebiten.IsKeyPressed(k)
		KeyPressed[k] = p

		if p {
			AnyKeyPressed = true
		}
	}
}

// Draw updates the input system. This is called by the UI.
func Draw() {
	LeftMouseButtonJustPressed = LeftMouseButtonPressed && LeftMouseButtonPressed != LastLeftMouseButtonPressed
	MiddleMouseButtonJustPressed = MiddleMouseButtonPressed && MiddleMouseButtonPressed != LastMiddleMouseButtonPressed
	RightMouseButtonJustPressed = RightMouseButtonPressed && RightMouseButtonPressed != LastRightMouseButtonPressed

	LastLeftMouseButtonPressed = LeftMouseButtonPressed
	LastMiddleMouseButtonPressed = MiddleMouseButtonPressed
	LastRightMouseButtonPressed = RightMouseButtonPressed
}

// AfterDraw updates the input system after the Ebiten Draw function has been called. This is called by the UI.
func AfterDraw() {
	InputChars = InputChars[:0]
	WheelX, WheelY = 0, 0
}
