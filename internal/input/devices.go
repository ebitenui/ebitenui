package input

import (
	"image"
	"math"
	"runtime"

	"github.com/ebitenui/ebitenui/internal/jsUtil"
	"github.com/hajimehoshi/ebiten/v2"
)

// Devices is where the handler reads the keyboard, the gamepads and the
// touches: Ebitengine, or a stand-in in tests (package inputtest).
type Devices interface {
	IsKeyPressed(k ebiten.Key) bool
	AppendGamepadIDs(ids []ebiten.GamepadID) []ebiten.GamepadID
	// IsStandardGamepadButtonPressed and StandardGamepadAxisValue report
	// nothing for a gamepad without the standard layout, as Ebitengine does.
	IsStandardGamepadButtonPressed(id ebiten.GamepadID, b ebiten.StandardGamepadButton) bool
	StandardGamepadAxisValue(id ebiten.GamepadID, a ebiten.StandardGamepadAxis) float64
	AppendTouchIDs(ids []ebiten.TouchID) []ebiten.TouchID
	TouchPosition(id ebiten.TouchID) (int, int)
	// IsTouchscreen reports whether the platform has a touch screen and no
	// mouse pointer: Android, iOS, or a browser on one of them.
	IsTouchscreen() bool
}

// touchscreen is whether the platform has a touch screen and no mouse
// pointer. Ebitengine reports no cursor position on Android and iOS.
var touchscreen = jsUtil.IsMobileBrowser() || runtime.GOOS == "android" || runtime.GOOS == "ios"

// Ebitengine reads the devices from Ebitengine.
type Ebitengine struct{}

func (Ebitengine) AppendTouchIDs(ids []ebiten.TouchID) []ebiten.TouchID {
	return ebiten.AppendTouchIDs(ids)
}

func (Ebitengine) TouchPosition(id ebiten.TouchID) (int, int) { return ebiten.TouchPosition(id) }

func (Ebitengine) IsTouchscreen() bool { return touchscreen }

func (Ebitengine) IsKeyPressed(k ebiten.Key) bool { return ebiten.IsKeyPressed(k) }

func (Ebitengine) AppendGamepadIDs(ids []ebiten.GamepadID) []ebiten.GamepadID {
	return ebiten.AppendGamepadIDs(ids)
}

func (Ebitengine) IsStandardGamepadButtonPressed(id ebiten.GamepadID, b ebiten.StandardGamepadButton) bool {
	return ebiten.IsStandardGamepadButtonPressed(id, b)
}

func (Ebitengine) StandardGamepadAxisValue(id ebiten.GamepadID, a ebiten.StandardGamepadAxis) float64 {
	return ebiten.StandardGamepadAxisValue(id, a)
}

// updateTouches makes the first finger on the screen the pointer, its left
// button down while the finger is. It reports whether touches placed the
// pointer this update, the update a touch lifts in included. The update
// after a touch lifts, the pointer leaves the screen, unless KeepTouchHover;
// on a platform with a mouse, the mouse then places it.
func (handler *DefaultInternalHandler) updateTouches() bool {
	handler.touchIDs = handler.Devices.AppendTouchIDs(handler.touchIDs[:0])
	handler.TouchCount = len(handler.touchIDs)
	handler.updatePair()
	switch {
	case handler.TouchCount > 0:
		if !handler.isTouched { // a new touch
			handler.TouchGesture, handler.KeepTouchHover = false, false
		}
		handler.isTouched = true
		handler.TouchGesture = handler.TouchGesture || handler.TouchCount > 1
		handler.LeftMouseButtonPressed = true
		handler.CursorX, handler.CursorY = handler.Devices.TouchPosition(handler.touchIDs[0])
		return true
	case handler.isTouched: // lifted: released where it was
		handler.isTouched = false
		handler.LeftMouseButtonPressed = false
		handler.lifted = true
		return true
	}
	if handler.lifted && !handler.KeepTouchHover {
		handler.CursorX, handler.CursorY = -1, -1 // with a mouse, the mouse puts it back
	}
	handler.lifted = false
	return false
}

// updatePair measures how the first two fingers moved since the last
// update, when the same two were down then: how much they spread apart,
// as a ratio, and how far their midpoint moved.
func (handler *DefaultInternalHandler) updatePair() {
	handler.TouchPinch, handler.TouchPanX, handler.TouchPanY = 1, 0, 0
	if len(handler.touchIDs) < 2 {
		handler.paired = false
		return
	}
	ids := [2]ebiten.TouchID{handler.touchIDs[0], handler.touchIDs[1]}
	var at [2]image.Point
	for i, id := range ids {
		at[i].X, at[i].Y = handler.Devices.TouchPosition(id)
	}
	if handler.paired && ids == handler.pairIDs {
		if before := spread(handler.pair); before > 0 {
			handler.TouchPinch = spread(at) / before
		}
		handler.TouchPanX = float64(at[0].X+at[1].X-handler.pair[0].X-handler.pair[1].X) / 2
		handler.TouchPanY = float64(at[0].Y+at[1].Y-handler.pair[0].Y-handler.pair[1].Y) / 2
	}
	handler.pair, handler.pairIDs, handler.paired = at, ids, true
}

func spread(p [2]image.Point) float64 {
	return math.Hypot(float64(p[1].X-p[0].X), float64(p[1].Y-p[0].Y))
}

// updateKeys takes the state of every key, and which went down since the
// last update.
func (handler *DefaultInternalHandler) updateKeys() {
	handler.AnyKeyPressed = false
	for k := range handler.KeyPressed {
		p := handler.Devices.IsKeyPressed(ebiten.Key(k))
		handler.KeyJustPressed[k] = p && !handler.KeyPressed[k]
		handler.KeyPressed[k] = p
		handler.AnyKeyPressed = handler.AnyKeyPressed || p
	}
}

// updateGamepads takes the buttons and axes of every gamepad with the
// standard layout. A button is down while any gamepad holds it, so one
// press seen from two gamepads (a platform's virtual pad beside the real
// one) counts once; an axis is as the gamepad pushing it furthest has it.
func (handler *DefaultInternalHandler) updateGamepads() {
	handler.gamepads = handler.Devices.AppendGamepadIDs(handler.gamepads[:0])
	for b := range handler.GamepadButtonPressed {
		p := false
		for _, id := range handler.gamepads {
			if handler.Devices.IsStandardGamepadButtonPressed(id, ebiten.StandardGamepadButton(b)) {
				p = true
				break
			}
		}
		handler.GamepadButtonJustPressed[b] = p && !handler.GamepadButtonPressed[b]
		handler.GamepadButtonPressed[b] = p
	}
	for a := range handler.GamepadAxis {
		v := 0.0
		for _, id := range handler.gamepads {
			if w := handler.Devices.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxis(a)); math.Abs(w) > math.Abs(v) {
				v = w
			}
		}
		handler.GamepadAxis[a] = v
	}
}
