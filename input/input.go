package input

import (
	"image"

	internalinput "github.com/ebitenui/ebitenui/internal/input"
	"github.com/hajimehoshi/ebiten/v2"
)

type CursorUpdater interface {
	// Called every Update call from Ebiten
	// Note that before this is called the current cursor shape is reset to DEFAULT every cycle
	Update()
	// Called at the beginning of every Draw call.
	Draw(screen *ebiten.Image)
	// Called at the end of every Draw call
	AfterDraw(screen *ebiten.Image)
	// MouseButtonPressed returns whether mouse button b is currently pressed.
	MouseButtonPressed(b ebiten.MouseButton) bool
	// MouseButtonJustPressed returns whether mouse button b has just been pressed.
	// It only returns true during the first frame that the button is pressed.
	MouseButtonJustPressed(b ebiten.MouseButton) bool
	// MouseButtonJustReleased returns whether mouse button b has just been released.
	// It only returns true during the first frame that the button is released.
	MouseButtonJustReleased(b ebiten.MouseButton) bool
	// CursorPosition returns the current cursor position.
	// If you define a CursorPosition that doesn't align with a system cursor you will need to
	// set the CursorDrawMode to Custom. This is because ebiten doesn't have a way to set the
	// cursor location manually
	CursorPosition() (int, int)
	// Returns the image to use as the cursor
	// EbitenUI by default will look for the following cursors:
	//  "EWResize"
	//  "NSResize"
	//  "Default"
	GetCursorImage(name string) *ebiten.Image
	// Returns how far from the CursorPosition to offset the cursor image.
	// This is best used with cursors such as resizing.
	GetCursorOffset(name string) image.Point
}

// This flag allows you to disable ebitenui's cursor management
var CursorManagementEnabled = true

// This variable indicates if the UI has currently being hovered over
var UIHovered = false

var currentCursorUpdater CursorUpdater = internalinput.InputHandler
var windowSize image.Point

// If the system cannot find a cursor image, it will revert to the system defaults.
// If cursorUpdater is nil the system will revert to the standard InputHandler system
//
// EbitenUI by default will look for the following cursors:
//
//	CURSOR_EWRESIZE  : "Cursor_EWResize"
//	CURSOR_NSRESIZE  : "Cursor_NSResize"
//	CURSOR_DEFAULT   : "Cursor_Default"
//	CURSOR_POINTER   : "Cursor_Pointer"
//	CURSOR_TEXT      : "Cursor_Text"
//	CURSOR_CROSSHAIR : "Cursor_Crosshair"
func SetCursorUpdater(cursorUpdater CursorUpdater) {
	if cursorUpdater == nil {
		cursorUpdater = internalinput.InputHandler
	}
	currentCursorUpdater = cursorUpdater
}

const (
	CURSOR_DEFAULT   = "Cursor_Default"
	CURSOR_EWRESIZE  = "Cursor_EWResize"
	CURSOR_NSRESIZE  = "Cursor_NSResize"
	CURSOR_POINTER   = "Cursor_Pointer"
	CURSOR_TEXT      = "Cursor_Text"
	CURSOR_CROSSHAIR = "Cursor_Crosshair"
	CURSOR_NONE      = "Cursor_None"
)

var currentCursor string = CURSOR_DEFAULT

func SetCursorShape(name string) {
	currentCursor = name
}

func SetCursorImage(name string, cursorImage *ebiten.Image) {
	internalinput.InputHandler.SetCursorImage(name, cursorImage, image.Point{})
}

func SetCursorImageWithOffset(name string, cursorImage *ebiten.Image, offset image.Point) {
	internalinput.InputHandler.SetCursorImage(name, cursorImage, offset)
}

// MouseButtonPressed returns whether mouse button b is currently pressed.
func MouseButtonPressed(b ebiten.MouseButton) bool {
	return currentCursorUpdater.MouseButtonPressed(b)
}

// MouseButtonJustPressed returns whether mouse button b has just been pressed.
// It only returns true during the first frame that the button is pressed.
func MouseButtonJustPressed(b ebiten.MouseButton) bool {
	return currentCursorUpdater.MouseButtonJustPressed(b)
}

// MouseButtonJustPressed returns whether mouse button b has just been pressed.
// It only returns true during the first frame that the button is pressed.
func MouseButtonJustReleased(b ebiten.MouseButton) bool {
	return currentCursorUpdater.MouseButtonJustReleased(b)
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

// MouseButtonJustPressedLayer returns whether mouse button b has just been pressed if input layer l
// is eligible to handle it. It only returns true during the first frame that the button is pressed.
func MouseButtonJustReleasedLayer(b ebiten.MouseButton, l *Layer) bool {
	if !MouseButtonJustReleased(b) {
		return false
	}

	x, y := CursorPosition()
	return l.ActiveFor(x, y, LayerEventTypeMouseButton)
}

// CursorPosition returns the current cursor position.
func CursorPosition() (int, int) {
	return currentCursorUpdater.CursorPosition()
}

// Wheel returns current mouse wheel movement.
func Wheel() (float64, float64) {
	return internalinput.InputHandler.WheelX, internalinput.InputHandler.WheelY
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
	return internalinput.InputHandler.InputChars
}

// KeyPressed returns whether key k is currently pressed.
func KeyPressed(k ebiten.Key) bool {
	p, ok := internalinput.InputHandler.KeyPressed[k]
	return ok && p
}

// AnyKeyPressed returns whether any key is currently pressed.
func AnyKeyPressed() bool {
	return internalinput.InputHandler.AnyKeyPressed
}

// This method returns the drawable screen size whether it is fullscreen or not.
func GetWindowSize() image.Point {
	return windowSize
}

func Update() {
	SetCursorShape(CURSOR_DEFAULT)
	currentCursorUpdater.Update()
	internalinput.InternalUIHovered = false
}

func Draw(screen *ebiten.Image) {
	windowSize = screen.Bounds().Max
	currentCursorUpdater.Draw(screen)
}

func AfterDraw(screen *ebiten.Image) {
	currentCursorUpdater.AfterDraw(screen)
	UIHovered = internalinput.InternalUIHovered
	if CursorManagementEnabled {
		// Process Cursor
		posX, posY := currentCursorUpdater.CursorPosition()
		//If cursor outside the window do nothing
		if posX < 0 || posY < 0 || posX > windowSize.X || posY > windowSize.Y {
			return
		}
		cursorImage := currentCursorUpdater.GetCursorImage(currentCursor)
		// If we have a cursor image hide current cursor and use it
		if cursorImage != nil {
			if ebiten.CursorMode() != ebiten.CursorModeHidden {
				ebiten.SetCursorMode(ebiten.CursorModeHidden)
			}
			cursorOffset := currentCursorUpdater.GetCursorOffset(currentCursor)
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(posX+cursorOffset.X), float64(posY+cursorOffset.Y))
			screen.DrawImage(cursorImage, op)
			// If we don't have an image use the system shapes.
		} else {
			cursorMode := ebiten.CursorModeVisible
			switch currentCursor {
			case CURSOR_DEFAULT:
				ebiten.SetCursorShape(ebiten.CursorShapeDefault)
			case CURSOR_EWRESIZE:
				ebiten.SetCursorShape(ebiten.CursorShapeEWResize)
			case CURSOR_NSRESIZE:
				ebiten.SetCursorShape(ebiten.CursorShapeNSResize)
			case CURSOR_TEXT:
				ebiten.SetCursorShape(ebiten.CursorShapeText)
			case CURSOR_CROSSHAIR:
				ebiten.SetCursorShape(ebiten.CursorShapeCrosshair)
			case CURSOR_POINTER:
				ebiten.SetCursorShape(ebiten.CursorShapePointer)
			case CURSOR_NONE:
				cursorMode = ebiten.CursorModeHidden
			default:
				ebiten.SetCursorShape(ebiten.CursorShapeDefault)
			}
			if ebiten.CursorMode() != cursorMode {
				ebiten.SetCursorMode(cursorMode)
			}

		}
	}
}
