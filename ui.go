package ebitenui

import (
	"image"
	"sort"

	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/input"
	"github.com/ebitenui/ebitenui/widget"

	"github.com/hajimehoshi/ebiten/v2"
)

// UI encapsulates a complete user interface that can be rendered onto the screen.
// There should only be exactly one UI per application.
type UI struct {
	// Container is the root container of the UI hierarchy.
	Container *widget.Container

	//If true the default tab/shift-tab to focus will be disabled
	DisableDefaultFocus bool

	// This exposes a Render call before the Container is drawn,
	// but after the Windows with DrawLayer < 0 are drawn.
	PreRenderHook widget.RenderFunc

	// This exposes a Render call after the Container is drawn,
	// but before the Windows with DrawLayer >= 0 (all by default) are drawn.
	PostRenderHook widget.RenderFunc

	//Theme
	PrimaryTheme  *widget.Theme
	previousTheme *widget.Theme

	focusedWidget widget.HasWidget
	inputLayerers []input.Layerer
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
		// Close all Ephemeral Windows (tooltip/dnd/etc)
		u.closeEphemeralWindows(0)
	}

	if (u.Container.GetWidget().GetTheme() == nil && u.PrimaryTheme != nil) || u.PrimaryTheme != u.previousTheme {
		u.Container.GetWidget().SetTheme(u.PrimaryTheme)
		u.Container.Validate()
		u.previousTheme = u.PrimaryTheme
	}

	u.handleFocusChangeRequest()

	index := 0
	for ; index < len(u.windows); index++ {
		if u.windows[index].DrawLayer < 0 {
			u.windows[index].Update()
		} else {
			break
		}
	}
	u.Container.Update()

	for ; index < len(u.windows); index++ {
		u.windows[index].Update()
	}

	event.ExecuteDeferred()
}

// Draw renders u onto screen. This function should be called in the Ebiten Draw function.
func (u *UI) Draw(screen *ebiten.Image) {
	input.Draw(screen)
	defer input.AfterDraw(screen)
	x, y := screen.Bounds().Dx(), screen.Bounds().Dy()
	rect := image.Rect(0, 0, x, y)
	u.setupInputLayers()
	u.Container.SetLocation(rect)
	u.render(screen)
	//Render elements that pop up (like combobox) on top of everything else
	widget.RenderDeferred(screen)
}

func (u *UI) setupInputLayers() {
	num := 1
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

	input.SetupInputLayersWithDeferred(u.inputLayerers)
}

func (u *UI) render(screen *ebiten.Image) {
	index := 0
	for ; index < len(u.windows); index++ {
		if u.windows[index].DrawLayer < 0 {
			u.windows[index].Render(screen)
		} else {
			break
		}
	}

	if u.PreRenderHook != nil {
		u.PreRenderHook(screen)
	}

	u.Container.Render(screen)

	if u.PostRenderHook != nil {
		u.PostRenderHook(screen)
	}

	for ; index < len(u.windows); index++ {
		u.windows[index].Render(screen)
	}
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
					u.ChangeFocus(widget.FOCUS_PREVIOUS)
				} else {
					u.ChangeFocus(widget.FOCUS_NEXT)
				}
			}
		} else {
			u.tabWasPressed = false
		}
	}
}

func (u *UI) ChangeFocus(direction widget.FocusDirection) {
	if u.focusedWidget != nil {
		if next := u.focusedWidget.(widget.Focuser).GetFocus(direction); next != nil {
			if !next.GetWidget().Disabled {
				next.Focus(true)
			}
		}
	}

	if direction == widget.FOCUS_NEXT || direction == widget.FOCUS_PREVIOUS {
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
				if direction == widget.FOCUS_PREVIOUS {
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
}

// AddWindow adds window w to ui for rendering. It returns a function to remove w from ui.
func (u *UI) AddWindow(w *widget.Window) widget.RemoveWindowFunc {
	if u.addWindow(w) {
		w.GetContainer().GetWidget().ContextMenuEvent.AddHandler(u.handleContextMenu)
		w.GetContainer().GetWidget().FocusEvent.AddHandler(u.handleFocusEvent)
		w.GetContainer().GetWidget().ToolTipEvent.AddHandler(u.handleToolTipEvent)
		w.GetContainer().GetWidget().DragAndDropEvent.AddHandler(u.handleDragAndDropEvent)

		if w.Modal && u.focusedWidget != nil {
			u.focusedWidget.(widget.Focuser).Focus(false)
		}
		// Close all Ephemeral Windows (tooltip/dnd/etc)
		u.closeEphemeralWindows(0)
	}

	return w.GetCloseFunction()
}

func (u *UI) addWindow(w *widget.Window) bool {
	if u.IsWindowOpen(w) {
		return false
	}

	if w.Contents.GetWidget().GetTheme() == nil {
		w.Contents.GetWidget().SetTheme(u.PrimaryTheme)
	}
	w.Contents.Validate()

	closeFunc := func() {
		u.removeWindow(w)
	}
	w.SetCloseFunction(closeFunc)

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
	if windowIdx != -1 && !w.Ephemeral {
		u.closeEphemeralWindows(windowIdx)
	}
}

// Used to close tooltips/dnd etc
func (u *UI) closeEphemeralWindows(windowIdx int) {
	for i := len(u.windows) - 1; i >= windowIdx; i-- {
		if u.windows[i].Ephemeral {
			u.windows = append(u.windows[:i], u.windows[i+1:]...)
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
