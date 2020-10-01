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
	inputLayerers []input.Layerer
	renderers     []widget.Renderer
	windows       []*widget.Window
}

type RemoveWindowFunc func()

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

	if rect != u.lastRect {
		u.Container.RequestRelayout()
	}

	u.handleFocus()
	u.setupInputLayers()
	u.layout.LayoutRoot(rect)
	u.render(screen)
}

func (u *UI) handleFocus() {
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
}

func (u *UI) setupInputLayers() {
	num := 1 // u.Container
	if len(u.windows) > 0 {
		num += len(u.windows)
	}
	if u.DragAndDrop != nil {
		num++
	}

	if cap(u.inputLayerers) < num {
		u.inputLayerers = make([]input.Layerer, num)
	}

	u.inputLayerers = u.inputLayerers[:0]
	u.inputLayerers = append(u.inputLayerers, u.Container)
	for _, w := range u.windows {
		u.inputLayerers = append(u.inputLayerers, w)
	}
	if u.DragAndDrop != nil {
		u.inputLayerers = append(u.inputLayerers, u.DragAndDrop)
	}

	// TODO: SetupInputLayersWithDeferred should reside in "internal" subpackage
	input.SetupInputLayersWithDeferred(u.inputLayerers)
}

func (u *UI) render(screen *ebiten.Image) {
	num := 1 // u.Container
	if len(u.windows) > 0 {
		num += len(u.windows)
	}
	if u.ToolTip != nil {
		num++
	}
	if u.DragAndDrop != nil {
		num++
	}

	if cap(u.renderers) < num {
		u.renderers = make([]widget.Renderer, num)
	}

	u.renderers = u.renderers[:0]
	u.renderers = append(u.renderers, u.Container)
	for _, w := range u.windows {
		u.renderers = append(u.renderers, w)
	}
	if u.ToolTip != nil {
		u.renderers = append(u.renderers, u.ToolTip)
	}
	if u.DragAndDrop != nil {
		u.renderers = append(u.renderers, u.DragAndDrop)
	}

	// TODO: RenderWithDeferred should reside in "internal" subpackage
	widget.RenderWithDeferred(screen, u.renderers)
}

func (u *UI) AddWindow(w *widget.Window) RemoveWindowFunc {
	u.windows = append(u.windows, w)

	return func() {
		u.removeWindow(w)
	}
}

func (u *UI) removeWindow(w *widget.Window) {
	for i, uw := range u.windows {
		if uw == w {
			u.windows = append(u.windows[:i], u.windows[i+1:]...)
			break
		}
	}
}

func (u *UI) initUI() {
	u.layout = widget.NewRootLayout(u.Container)
}
