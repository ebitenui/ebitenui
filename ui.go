package ebitenui

import (
	"image"

	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/input"
	internalinput "github.com/ebitenui/ebitenui/internal/input"
	"github.com/ebitenui/ebitenui/widget"

	"github.com/hajimehoshi/ebiten/v2"
)

// UI encapsulates a complete user interface that can be rendered onto the screen.
// There should only be exactly one UI per application.
type UI struct {
	// Container is the root container of the UI hierarchy.
	Container *widget.Container

	// ToolTip is used to render mouse hover tool tips. It may be nil to disable rendering.
	ToolTip *widget.ToolTip

	// DragAndDrop is used to render drag widgets while dragging and dropping. It may be nil to disable rendering.
	DragAndDrop *widget.DragAndDrop

	lastRect      image.Rectangle
	focusedWidget widget.HasWidget
	inputLayerers []input.Layerer
	renderers     []widget.Renderer
	windows       []*widget.Window
}

// Update updates u. This method should be called in the Ebiten Update function.
func (u *UI) Update() {
	internalinput.Update()
}

// Draw renders u onto screen. This function should be called in the Ebiten Draw function.
//
// If screen's size changes from one frame to the next, u.Container.RequestRelayout is called.
func (u *UI) Draw(screen *ebiten.Image) {
	event.ExecuteDeferred()

	internalinput.Draw()
	defer internalinput.AfterDraw()

	w, h := screen.Size()
	rect := image.Rect(0, 0, w, h)

	defer func() {
		u.lastRect = rect
	}()

	if rect != u.lastRect {
		u.Container.RequestRelayout()
	}

	u.handleFocus()
	u.setupInputLayers()
	u.Container.SetLocation(rect)
	u.render(screen)
}

func (u *UI) handleFocus() {
	if input.MouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if u.focusedWidget != nil {
			u.focusedWidget.(widget.Focuser).Focus(false)
			u.focusedWidget = nil
		}
		x, y := input.CursorPosition()

		for _, window := range u.windows {
			w := window.Contents.WidgetAt(x, y)
			if w != nil && w.GetWidget().EffectiveInputLayer().ActiveFor(x, y, input.LayerEventTypeAll) {
				if f, ok := w.(widget.Focuser); ok {
					f.Focus(true)
					u.focusedWidget = w
					return
				}
			}
		}

		w := u.Container.WidgetAt(x, y)
		if w != nil && w.GetWidget().EffectiveInputLayer().ActiveFor(x, y, input.LayerEventTypeAll) {
			if f, ok := w.(widget.Focuser); ok {
				f.Focus(true)
				u.focusedWidget = w
				return
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

// AddWindow adds window w to u for rendering. It returns a function to remove w from u.
func (u *UI) AddWindow(w *widget.Window) widget.RemoveWindowFunc {
	u.windows = append(u.windows, w)
	closeFunc := func() {
		u.removeWindow(w)
	}

	w.SetCloseFunction(closeFunc)
	return closeFunc
}

func (u *UI) removeWindow(w *widget.Window) {
	for i, uw := range u.windows {
		if uw == w {
			u.windows = append(u.windows[:i], u.windows[i+1:]...)
			break
		}
	}
}

func (u *UI) HasFocus() bool {
	return u.focusedWidget != nil
}
