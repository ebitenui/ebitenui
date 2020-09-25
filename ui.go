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

	init     sync.Once
	layout   *widget.RootLayout
	lastRect image.Rectangle
}

// Update updates u. This function should be called in the Ebiten Update function.
func (u *UI) Update() {
	u.init.Do(u.initUI)

	internalevent.ExecuteDeferredActions()
	internalinput.Update()
}

// Draw renders u onto screen, with rect as the area reserved for rendering.
// This function should be called in the Ebiten Draw function.
//
// If rect changes from one frame to the next, u.Container.RequestRelayout is called.
func (u *UI) Draw(screen *ebiten.Image, rect image.Rectangle) {
	u.init.Do(u.initUI)

	internalinput.Draw()

	defer func() {
		u.lastRect = rect
	}()

	if rect != u.lastRect {
		u.Container.RequestRelayout()
	}

	// TODO: SetupInputLayersWithDeferred should reside in "internal" subpackage
	input.SetupInputLayersWithDeferred(u.Container)
	if u.DragAndDrop != nil {
		input.SetupInputLayersWithDeferred(u.DragAndDrop)
	}

	u.layout.LayoutRoot(rect)

	// TODO: RenderWithDeferred should reside in "internal" subpackage
	widget.RenderWithDeferred(u.Container, screen)
	if u.ToolTip != nil {
		widget.RenderWithDeferred(u.ToolTip, screen)
	}
	if u.DragAndDrop != nil {
		widget.RenderWithDeferred(u.DragAndDrop, screen)
	}
}

func (u *UI) initUI() {
	u.layout = widget.NewRootLayout(u.Container)
}
