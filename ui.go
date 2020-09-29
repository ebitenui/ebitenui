package ebitenui

import (
	"image"
	"sync"

	"github.com/blizzy78/ebitenui/input"
	internalevent "github.com/blizzy78/ebitenui/internal/event"
	internalinput "github.com/blizzy78/ebitenui/internal/input"
	"github.com/blizzy78/ebitenui/widget"

	"github.com/hajimehoshi/ebiten"
)

// UI encapsulates a complete user interface that can be rendered onto the screen.
// There should only be exactly one UI per application.
type UI struct {
	Container   *widget.Container
	ToolTip     *widget.ToolTip
	DragAndDrop *widget.DragAndDrop

	init          sync.Once
	layout        *widget.RootLayout
	lastRect      image.Rectangle
	focusedWidget widget.HasWidget
}

// Update updates u. This function should be called in the Ebiten Update function.
func (u *UI) Update() {
	u.init.Do(u.initUI)
	internalinput.Update()
}

// Draw renders u onto screen, with rect as the area reserved for rendering.
// This function should be called in the Ebiten Draw function.
//
// If rect changes from one frame to the next, u.Container.RequestRelayout is called.
func (u *UI) Draw(screen *ebiten.Image, rect image.Rectangle) {
	u.init.Do(u.initUI)

	internalevent.ExecuteDeferredActions()

	internalinput.Draw()
	defer internalinput.AfterDraw()

	defer func() {
		u.lastRect = rect
	}()

	if input.MouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if u.focusedWidget != nil {
			u.focusedWidget.(widget.Focuser).Focus(false)
			u.focusedWidget = nil
		}

		x, y := input.CursorPosition()
		w := u.Container.WidgetAt(x, y)
		if w != nil {
			if f, ok := w.(widget.Focuser); ok {
				f.Focus(true)
				u.focusedWidget = w
			}
		}
	}

	if rect != u.lastRect {
		u.Container.RequestRelayout()
	}

	// TODO: SetupInputLayersWithDeferred should reside in "internal" subpackage
	input.SetupInputLayersWithDeferred(u.Container, u.DragAndDrop)

	u.layout.LayoutRoot(rect)

	// TODO: RenderWithDeferred should reside in "internal" subpackage
	widget.RenderWithDeferred(screen, u.Container, u.ToolTip, u.DragAndDrop)
}

func (u *UI) initUI() {
	u.layout = widget.NewRootLayout(u.Container)
}
