package ebitenui

import (
	"image"
	"sort"

	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/input"
	"github.com/ebitenui/ebitenui/widget"

	"github.com/hajimehoshi/ebiten/v2"
)

type FocusDirection int

const (
	FOCUS_NEXT FocusDirection = iota
	FOCUS_PREVIOUS
)

// UI encapsulates a complete user interface that can be rendered onto the screen.
// There should only be exactly one UI per application.
type UI struct {
	// Container is the root container of the UI hierarchy.
	Container *widget.Container

	//If true the default tab/shift-tab to focus will be disabled
	DisableDefaultFocus bool

	focusedWidget widget.HasWidget
	inputLayerers []input.Layerer
	renderers     []widget.Renderer
	windows       []*widget.Window

	previousContainer *widget.Container
	tabWasPressed     bool
}

// Update updates u. This method should be called in the Ebiten Update function.
func (u *UI) Update() {
	input.Update()
	if u.previousContainer == nil || u.previousContainer != u.Container {
		u.Container.GetWidget().ContextMenuEvent.AddHandler(u.handleContextMenu)
		u.Container.GetWidget().FocusEvent.AddHandler(u.handleFocusEvent)
		u.Container.GetWidget().ToolTipEvent.AddHandler(u.handleToolTipEvent)
		u.Container.GetWidget().DragAndDropEvent.AddHandler(u.handleDragAndDropEvent)

		u.previousContainer = u.Container
	}
}

// Draw renders u onto screen. This function should be called in the Ebiten Draw function.
func (u *UI) Draw(screen *ebiten.Image) {
	event.ExecuteDeferred()

	x, y := screen.Bounds().Dx(), screen.Bounds().Dy()
	rect := image.Rect(0, 0, x, y)

	u.handleFocusChangeRequest()
	u.setupInputLayers()
	u.Container.SetLocation(rect)
	u.render(screen)
	input.Draw(screen)
}

func (u *UI) handleContextMenu(args interface{}) {
	a := args.(*widget.WidgetContextMenuEventArgs)

	x, y := a.Widget.ContextMenu.PreferredSize()
	r := image.Rect(0, 0, x, y)
	r = r.Add(a.Location)
	a.Widget.ContextMenuWindow = widget.NewWindow(
		widget.WindowOpts.Contents(a.Widget.ContextMenu),
		widget.WindowOpts.CloseMode(a.Widget.ContextMenuCloseMode),
		widget.WindowOpts.Modal(),
		widget.WindowOpts.Location(r),
	)
	u.AddWindow(a.Widget.ContextMenuWindow)
}

func (u *UI) handleFocusEvent(args interface{}) {
	a := args.(*widget.WidgetFocusEventArgs)
	if a.Focused { //New widget focused
		if u.focusedWidget != nil && u.focusedWidget != a.Widget {
			u.focusedWidget.(widget.Focuser).Focus(false)
		}
		u.focusedWidget = a.Widget
	} else if a.Widget == u.focusedWidget { //Current widget focus removed
		u.focusedWidget = nil
	} else if a.Widget == nil { //Clicked out of focusable widgets
		if u.focusedWidget != nil {
			//If we didnt just click on the same widget
			if !a.Location.In(u.focusedWidget.GetWidget().Rect) {
				u.focusedWidget.(widget.Focuser).Focus(false)
				u.focusedWidget = nil
			}
		}
	}
}

func (u *UI) handleToolTipEvent(args interface{}) {
	a := args.(*widget.WidgetToolTipEventArgs)
	a.Window.Ephemeral = true
	if a.Show {
		u.addWindow(a.Window)
	} else {
		u.removeWindow(a.Window)
	}
}

func (u *UI) handleDragAndDropEvent(args interface{}) {
	a := args.(*widget.WidgetDragAndDropEventArgs)
	if a.Show {
		a.Window.Ephemeral = true
		a.DnD.AvailableDropTargets = u.getDropTargets()
		u.addWindow(a.Window)
	} else {
		a.DnD.AvailableDropTargets = nil
		u.removeWindow(a.Window)
	}
}

func (u *UI) getDropTargets() []widget.HasWidget {
	dropTargets := u.Container.GetDropTargets()
	//Loop through the windows array in reverse. If we find a modal window, only loop through its droppable widgets
	for i := len(u.windows) - 1; i >= 0; i-- {
		if !u.windows[i].Modal {
			dropTargets = append(dropTargets, u.windows[i].GetContainer().GetDropTargets()...)
		} else {
			dropTargets = u.windows[i].GetContainer().GetDropTargets()
			break
		}
	}
	return dropTargets
}

func (u *UI) handleFocusChangeRequest() {
	if !u.DisableDefaultFocus {
		if input.KeyPressed(ebiten.KeyTab) {
			if !u.tabWasPressed {
				u.tabWasPressed = true
				if input.KeyPressed(ebiten.KeyShift) {
					u.ChangeFocus(FOCUS_PREVIOUS)
				} else {
					u.ChangeFocus(FOCUS_NEXT)
				}
			}
		} else {
			u.tabWasPressed = false
		}
	}
}

func (u *UI) ChangeFocus(direction FocusDirection) {
	focusableWidgets := u.Container.GetFocusers()
	//Loop through the windows array in reverse. If we find a modal window, only loop through its focusable widgets
	for i := len(u.windows) - 1; i >= 0; i-- {
		if !u.windows[i].Modal {
			focusableWidgets = append(focusableWidgets, u.windows[i].GetContainer().GetFocusers()...)
		} else {
			focusableWidgets = u.windows[i].GetContainer().GetFocusers()
			break
		}
	}
	len := len(focusableWidgets)
	if len == 1 {
		if u.focusedWidget != nil && u.focusedWidget.(widget.Focuser) != focusableWidgets[0] {
			u.focusedWidget.(widget.Focuser).Focus(false)
		}
		focusableWidgets[0].Focus(true)
	} else if len > 0 {
		sort.SliceStable(focusableWidgets, func(i, j int) bool {
			return focusableWidgets[i].TabOrder() < focusableWidgets[j].TabOrder()
		})
		if u.focusedWidget != nil {
			if direction == FOCUS_PREVIOUS {
				for i := 0; i < len; i++ {
					if focusableWidgets[i] == u.focusedWidget.(widget.Focuser) {
						u.focusedWidget.(widget.Focuser).Focus(false)
						if i == 0 {
							focusableWidgets[len-1].Focus(true)
						} else {
							focusableWidgets[i-1].Focus(true)
						}
						return
					}
				}
			} else {
				for i := 0; i < len-1; i++ {
					if focusableWidgets[i] == u.focusedWidget.(widget.Focuser) {
						u.focusedWidget.(widget.Focuser).Focus(false)
						focusableWidgets[i+1].Focus(true)
						return
					}
				}
			}
			u.focusedWidget.(widget.Focuser).Focus(false)
		}
		focusableWidgets[0].Focus(true)
	}
}

func (u *UI) setupInputLayers() {
	num := 1 // u.Container
	if len(u.windows) > 0 {
		num += len(u.windows)
	}

	if cap(u.inputLayerers) < num {
		u.inputLayerers = make([]input.Layerer, num)
	}

	u.inputLayerers = u.inputLayerers[:0]
	u.inputLayerers = append(u.inputLayerers, u.Container)
	for _, w := range u.windows {
		u.inputLayerers = append(u.inputLayerers, w)
	}

	// TODO: SetupInputLayersWithDeferred should reside in "internal" subpackage
	input.SetupInputLayersWithDeferred(u.inputLayerers)
}

func (u *UI) render(screen *ebiten.Image) {
	num := 1 // u.Container
	if len(u.windows) > 0 {
		num += len(u.windows)
	}

	if cap(u.renderers) < num {
		u.renderers = make([]widget.Renderer, num)
	}
	u.renderers = u.renderers[:0]

	index := 0
	for ; index < len(u.windows); index++ {
		if u.windows[index].DrawLayer < 0 {
			u.renderers = append(u.renderers, u.windows[index])
		} else {
			break
		}
	}
	u.renderers = append(u.renderers, u.Container)

	for ; index < len(u.windows); index++ {
		u.renderers = append(u.renderers, u.windows[index])
	}

	// TODO: RenderWithDeferred should reside in "internal" subpackage
	widget.RenderWithDeferred(screen, u.renderers)
}

// AddWindow adds window w to u for rendering. It returns a function to remove w from u.
func (u *UI) AddWindow(w *widget.Window) widget.RemoveWindowFunc {
	closeFunc := func() {
		u.removeWindow(w)
	}

	if u.addWindow(w) {
		w.GetContainer().GetWidget().ContextMenuEvent.AddHandler(u.handleContextMenu)
		w.GetContainer().GetWidget().FocusEvent.AddHandler(u.handleFocusEvent)
		w.GetContainer().GetWidget().ToolTipEvent.AddHandler(u.handleToolTipEvent)
		w.GetContainer().GetWidget().DragAndDropEvent.AddHandler(u.handleDragAndDropEvent)

		if w.Modal && u.focusedWidget != nil {
			u.focusedWidget.(widget.Focuser).Focus(false)
		}

		w.SetCloseFunction(closeFunc)
	}

	return closeFunc
}

func (u *UI) addWindow(w *widget.Window) bool {
	if u.IsWindowOpen(w) {
		return false
	}
	u.windows = append(u.windows, w)

	sort.SliceStable(u.windows, func(i, j int) bool {
		return u.windows[i].DrawLayer < u.windows[j].DrawLayer
	})
	return true
}

func (u *UI) removeWindow(w *widget.Window) {
	windowIdx := -1
	for i := range u.windows {
		if u.windows[i] == w {
			u.windows = append(u.windows[:i], u.windows[i+1:]...)
			windowIdx = i
			break
		}
	}
	if windowIdx != -1 {
		for i := len(u.windows) - 1; i >= windowIdx; i-- {
			if u.windows[i].Ephemeral {
				u.windows = append(u.windows[:i], u.windows[i+1:]...)
			}
		}
	}
}

func (u *UI) IsWindowOpen(w *widget.Window) bool {
	for i := range u.windows {
		if u.windows[i] == w {
			return true
		}
	}
	return false
}

func (u *UI) HasFocus() bool {
	for i := len(u.windows) - 1; i >= 0; i-- {
		if u.windows[i].Modal {
			return true
		}
	}
	return u.focusedWidget != nil
}

func (u *UI) ClearFocus() {
	if u.focusedWidget != nil {
		u.focusedWidget.(widget.Focuser).Focus(false)
	}
}

func (u *UI) GetFocusedWidget() widget.Focuser {
	if u.focusedWidget != nil {
		return u.focusedWidget.(widget.Focuser)
	}
	return nil
}
