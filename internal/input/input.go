package input

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type DefaultInternalHandler struct {
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

	InputChars     []rune
	KeyPressed     map[ebiten.Key]bool
	AnyKeyPressed  bool
	isTouched      bool
	cursorImages   map[string]*ebiten.Image
	cursorOffset   map[string]image.Point
}

var InputHandler *DefaultInternalHandler = &DefaultInternalHandler{
	KeyPressed:     make(map[ebiten.Key]bool),
	cursorImages:   make(map[string]*ebiten.Image),
	cursorOffset:   make(map[string]image.Point)}

// Update updates the input system. This is called by the UI.
func (handler *DefaultInternalHandler) Update() {
	touches := ebiten.TouchIDs()
	if len(touches) > 0 {
		handler.isTouched = true
	}
	if handler.isTouched {
		if len(touches) > 0 {
			handler.LeftMouseButtonPressed = true
			handler.CursorX, handler.CursorY = ebiten.TouchPosition(touches[0])
		} else {
			handler.LeftMouseButtonPressed = false
			handler.isTouched = false
		}
	} else {
		handler.LeftMouseButtonPressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
		handler.MiddleMouseButtonPressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle)
		handler.RightMouseButtonPressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight)
		handler.CursorX, handler.CursorY = ebiten.CursorPosition()
	}

	wx, wy := ebiten.Wheel()
	handler.WheelX += wx
	handler.WheelY += wy

	handler.InputChars = ebiten.AppendInputChars(handler.InputChars)
	handler.AnyKeyPressed = false
	for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
		p := ebiten.IsKeyPressed(k)
		handler.KeyPressed[k] = p
		if p {
			handler.AnyKeyPressed = true
		}
	}

}

func (handler *DefaultInternalHandler) Draw(screen *ebiten.Image) {
	handler.LeftMouseButtonJustPressed = handler.LeftMouseButtonPressed && handler.LeftMouseButtonPressed != handler.LastLeftMouseButtonPressed
	handler.MiddleMouseButtonJustPressed = handler.MiddleMouseButtonPressed && handler.MiddleMouseButtonPressed != handler.LastMiddleMouseButtonPressed
	handler.RightMouseButtonJustPressed = handler.RightMouseButtonPressed && handler.RightMouseButtonPressed != handler.LastRightMouseButtonPressed

	handler.LastLeftMouseButtonPressed = handler.LeftMouseButtonPressed
	handler.LastMiddleMouseButtonPressed = handler.MiddleMouseButtonPressed
	handler.LastRightMouseButtonPressed = handler.RightMouseButtonPressed

}

func (handler *DefaultInternalHandler) AfterDraw(screen *ebiten.Image) {
	handler.InputChars = handler.InputChars[:0]
	handler.WheelX, handler.WheelY = 0, 0
}

func (handler *DefaultInternalHandler) MouseButtonPressed(b ebiten.MouseButton) bool {
	switch b {
	case ebiten.MouseButtonLeft:
		return handler.LeftMouseButtonPressed
	case ebiten.MouseButtonMiddle:
		return handler.MiddleMouseButtonPressed
	case ebiten.MouseButtonRight:
		return handler.RightMouseButtonPressed
	default:
		return false
	}
}
func (handler *DefaultInternalHandler) MouseButtonJustPressed(b ebiten.MouseButton) bool {
	switch b {
	case ebiten.MouseButtonLeft:
		return handler.LeftMouseButtonJustPressed
	case ebiten.MouseButtonMiddle:
		return handler.MiddleMouseButtonJustPressed
	case ebiten.MouseButtonRight:
		return handler.RightMouseButtonJustPressed
	default:
		return false
	}
}

func (handler *DefaultInternalHandler) CursorPosition() (int, int) {
	return handler.CursorX, handler.CursorY
}

func (handler *DefaultInternalHandler) GetCursorImage(name string) *ebiten.Image {
	return handler.cursorImages[name]
}

func (handler *DefaultInternalHandler) GetCursorOffset(name string) image.Point {
	return handler.cursorOffset[name]
}
func (handler *DefaultInternalHandler) SetCursorImage(name string, cursorImage *ebiten.Image, offset image.Point) {
	handler.cursorImages[name] = cursorImage
	handler.cursorOffset[name] = offset
}
