package input

import "github.com/hajimehoshi/ebiten/v2"

// Devices is where the handler reads the keyboard and the gamepads:
// Ebitengine, or a stand-in in tests (package inputtest).
type Devices interface {
	IsKeyPressed(k ebiten.Key) bool
	AppendGamepadIDs(ids []ebiten.GamepadID) []ebiten.GamepadID
	// IsStandardGamepadButtonPressed reports false for a gamepad without
	// the standard layout, as Ebitengine does.
	IsStandardGamepadButtonPressed(id ebiten.GamepadID, b ebiten.StandardGamepadButton) bool
}

// Ebitengine reads the keyboard and the gamepads from Ebitengine.
type Ebitengine struct{}

func (Ebitengine) IsKeyPressed(k ebiten.Key) bool { return ebiten.IsKeyPressed(k) }

func (Ebitengine) AppendGamepadIDs(ids []ebiten.GamepadID) []ebiten.GamepadID {
	return ebiten.AppendGamepadIDs(ids)
}

func (Ebitengine) IsStandardGamepadButtonPressed(id ebiten.GamepadID, b ebiten.StandardGamepadButton) bool {
	return ebiten.IsStandardGamepadButtonPressed(id, b)
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

// updateGamepads takes the buttons of every gamepad with the standard
// layout. A button is down while any gamepad holds it, so one press seen
// from two gamepads (a platform's virtual pad beside the real one) counts
// once.
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
}
