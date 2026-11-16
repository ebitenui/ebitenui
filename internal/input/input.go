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

	LeftMouseButtonJustReleased   bool
	MiddleMouseButtonJustReleased bool
	RightMouseButtonJustReleased  bool

	LastLeftMouseButtonPressed   bool
	LastMiddleMouseButtonPressed bool
	LastRightMouseButtonPressed  bool

	InputChars []rune
	// KeyPressed holds the keys down, KeyJustPressed those that went down
	// since the last update.
	KeyPressed     [ebiten.KeyMax + 1]bool
	KeyJustPressed [ebiten.KeyMax + 1]bool
	AnyKeyPressed  bool

	// GamepadButtonPressed holds the buttons down on any gamepad with the
	// standard layout, GamepadButtonJustPressed those that went down since
	// the last update.
	GamepadButtonPressed     [ebiten.StandardGamepadButtonMax + 1]bool
	GamepadButtonJustPressed [ebiten.StandardGamepadButtonMax + 1]bool
	// GamepadAxis holds each axis as the gamepad pushing it furthest has it,
	// from -1 to 1.
	GamepadAxis [ebiten.StandardGamepadAxisMax + 1]float64

	// Devices is where the keyboard, the gamepads and the touches are read.
	Devices  Devices
	gamepads []ebiten.GamepadID // reused each update

	// TouchCount is the number of fingers on the screen. TouchGesture is
	// whether the touch going on, or the one that just lifted, has had two
	// fingers or more: a gesture, not a tap.
	TouchCount   int
	TouchGesture bool
	// KeepTouchHover, set while a touch is down, keeps the pointer where the
	// touch lifts, as a hover; otherwise on a touch screen it leaves the
	// screen the update after, so a tap leaves nothing hovered. A long press
	// sets it.
	KeepTouchHover bool
	lifted         bool // a touch lifted in the last update

	isTouched    bool
	cursorImages map[string]*ebiten.Image
	cursorOffset map[string]image.Point

	touchIDs []ebiten.TouchID
}

var InternalUIHovered = false

var InputHandler *DefaultInternalHandler = &DefaultInternalHandler{
	Devices:      Ebitengine{},
	cursorImages: make(map[string]*ebiten.Image),
	cursorOffset: make(map[string]image.Point),
}

// Update updates the input system. This is called by the UI.
func (handler *DefaultInternalHandler) Update() {
	if !handler.updateTouches() && !handler.Devices.IsTouchscreen() {
		// Only on platforms with a mouse pointer: on touch screens
		// ebiten.CursorPosition() reports (0, 0), and the pointer stays
		// where the touches put it.
		//
		// TODO: maybe get rid of this special condition when fireEvents are
		// moved to the Update() tree.
		// See issue #100.
		handler.LeftMouseButtonPressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
		handler.MiddleMouseButtonPressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle)
		handler.RightMouseButtonPressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight)
		handler.CursorX, handler.CursorY = ebiten.CursorPosition()
	}

	wx, wy := ebiten.Wheel()
	handler.WheelX += wx
	handler.WheelY += wy

	handler.InputChars = ebiten.AppendInputChars(handler.InputChars)
	handler.updateKeys()
	handler.updateGamepads()
	handler.LeftMouseButtonJustPressed = handler.LeftMouseButtonPressed && handler.LeftMouseButtonPressed != handler.LastLeftMouseButtonPressed
	handler.MiddleMouseButtonJustPressed = handler.MiddleMouseButtonPressed && handler.MiddleMouseButtonPressed != handler.LastMiddleMouseButtonPressed
	handler.RightMouseButtonJustPressed = handler.RightMouseButtonPressed && handler.RightMouseButtonPressed != handler.LastRightMouseButtonPressed

	handler.LeftMouseButtonJustReleased = !handler.LeftMouseButtonPressed && handler.LeftMouseButtonPressed != handler.LastLeftMouseButtonPressed
	handler.MiddleMouseButtonJustReleased = !handler.MiddleMouseButtonPressed && handler.MiddleMouseButtonPressed != handler.LastMiddleMouseButtonPressed
	handler.RightMouseButtonJustReleased = !handler.RightMouseButtonPressed && handler.RightMouseButtonPressed != handler.LastRightMouseButtonPressed

	handler.LastLeftMouseButtonPressed = handler.LeftMouseButtonPressed
	handler.LastMiddleMouseButtonPressed = handler.MiddleMouseButtonPressed
	handler.LastRightMouseButtonPressed = handler.RightMouseButtonPressed
}

func (handler *DefaultInternalHandler) AfterUpdate() {
	handler.InputChars = handler.InputChars[:0]
	handler.WheelX, handler.WheelY = 0, 0
}

func (handler *DefaultInternalHandler) Draw(screen *ebiten.Image) {

}

func (handler *DefaultInternalHandler) AfterDraw(screen *ebiten.Image) {

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

func (handler *DefaultInternalHandler) MouseButtonJustReleased(b ebiten.MouseButton) bool {
	switch b {
	case ebiten.MouseButtonLeft:
		return handler.LeftMouseButtonJustReleased
	case ebiten.MouseButtonMiddle:
		return handler.MiddleMouseButtonJustReleased
	case ebiten.MouseButtonRight:
		return handler.RightMouseButtonJustReleased
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
