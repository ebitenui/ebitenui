package widget

import (
	"image"
	"math"

	"github.com/ebitenui/ebitenui/input"
	"github.com/hajimehoshi/ebiten/v2"
)

type DragAndDrop struct {
	AvailableDropTargets []HasWidget

	contentsCreater      DragContentsCreater
	minDragStartDistance int
	state                dragAndDropState
	dragWidget           *Container
	window               *Window
}

type DragAndDropOpt func(d *DragAndDrop)

type DragAndDropOptions struct {
}

var DragAndDropOpts DragAndDropOptions

type DragContentsCreater interface {
	Create(*Widget, int, int) (*Container, interface{})
}

type DragContentsUpdater interface {
	// arg1 - X
	// arg2 - Y
	// arg3 - isDroppable
	// arg4 - HasWidget if droppable
	// arg5 - DragData
	Update(int, int, bool, HasWidget, interface{})
}

type dragAndDropState func(*Widget) (dragAndDropState, bool)

func NewDragAndDrop(opts ...DragAndDropOpt) *DragAndDrop {
	d := &DragAndDrop{
		minDragStartDistance: 15,
	}
	d.state = d.idleState()

	for _, o := range opts {
		o(d)
	}

	return d
}

func (o DragAndDropOptions) ContentsCreater(c DragContentsCreater) DragAndDropOpt {
	return func(d *DragAndDrop) {
		d.contentsCreater = c
	}
}

func (o DragAndDropOptions) MinDragStartDistance(d int) DragAndDropOpt {
	return func(dnd *DragAndDrop) {
		dnd.minDragStartDistance = d
	}
}

func (d *DragAndDrop) SetupInputLayer(def input.DeferredSetupInputLayerFunc) {
	if d.dragWidget != nil {
		d.dragWidget.GetWidget().ElevateToNewInputLayer(&input.Layer{
			DebugLabel: "drag widget",
			EventTypes: input.LayerEventTypeAll,
			BlockLower: true,
			FullScreen: true,
		})
	}
}

func (d *DragAndDrop) Render(parent *Widget, screen *ebiten.Image, def DeferredRenderFunc) {
	newState, _ := d.state(parent)
	if newState != nil {
		d.state = newState
	}
}

func (d *DragAndDrop) idleState() dragAndDropState {
	return func(parent *Widget) (dragAndDropState, bool) {
		d.dragWidget = nil

		if !input.MouseButtonJustPressed(ebiten.MouseButtonLeft) {
			if d.window != nil {
				parent.FireDragAndDropEvent(d.window, false, d)
				d.window = nil
			}
			return nil, false
		}

		x, y := input.CursorPosition()
		p := image.Point{x, y}
		if !p.In(parent.Rect) {
			return nil, false
		}

		return d.dragArmedState(x, y), true
	}
}

func (d *DragAndDrop) dragArmedState(srcX int, srcY int) dragAndDropState {
	return func(_ *Widget) (dragAndDropState, bool) {
		if !input.MouseButtonPressed(ebiten.MouseButtonLeft) {
			return d.idleState(), false
		}

		x, y := input.CursorPosition()
		dx, dy := math.Abs(float64(x-srcX)), math.Abs(float64(y-srcY))
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist < float64(d.minDragStartDistance) {
			return nil, false
		}

		return d.draggingState(srcX, srcY, nil, nil), true
	}
}

func (d *DragAndDrop) draggingState(srcX int, srcY int, dragWidget *Container, dragData interface{}) dragAndDropState {
	return func(parent *Widget) (dragAndDropState, bool) {
		x, y := input.CursorPosition()

		if !input.MouseButtonPressed(ebiten.MouseButtonLeft) {
			return d.droppingState(srcX, srcY, x, y, dragData), true
		}

		if dragWidget == nil {
			dragWidget, dragData = d.contentsCreater.Create(parent, srcX, srcY)
			if dragWidget == nil {
				return d.idleState(), false
			}
			d.window = NewWindow(
				WindowOpts.CloseMode(NONE),
				WindowOpts.Contents(dragWidget),
			)
			parent.FireDragAndDropEvent(d.window, true, d)
		}

		defer func() {
			d.dragWidget = dragWidget
		}()

		if u, ok := d.contentsCreater.(DragContentsUpdater); ok {
			droppable := false
			var element HasWidget
			args := &DragAndDropDroppedEventArgs{
				Source:  parent,
				SourceX: srcX,
				SourceY: srcY,
				TargetX: x,
				TargetY: y,
				Data:    dragData,
			}
			p := image.Point{x, y}
			for _, target := range d.AvailableDropTargets {
				if p.In(target.GetWidget().Rect) && target.GetWidget().canDrop(args) {
					droppable = true
					element = target
					break
				}
			}
			u.Update(x, y, droppable, element, dragData)
		}

		sx, sy := dragWidget.PreferredSize()
		r := image.Rect(0, 0, sx, sy)
		r = r.Add(image.Point{x, y})
		r = r.Sub(image.Point{sx / 2, sy / 2})
		d.window.SetLocation(r)
		dragWidget.SetLocation(r)
		dragWidget.RequestRelayout()

		return d.draggingState(srcX, srcY, dragWidget, dragData), false
	}
}

func (d *DragAndDrop) droppingState(srcX int, srcY int, x int, y int, dragData interface{}) dragAndDropState {
	return func(parent *Widget) (dragAndDropState, bool) {
		args := &DragAndDropDroppedEventArgs{
			Source:  parent,
			SourceX: srcX,
			SourceY: srcY,
			TargetX: x,
			TargetY: y,
			Data:    dragData,
		}
		p := image.Point{x, y}
		for _, target := range d.AvailableDropTargets {
			if p.In(target.GetWidget().Rect) && target.GetWidget().canDrop(args) {
				if target.GetWidget().drop != nil {
					target.GetWidget().drop(args)
				}
				break
			}
		}

		return d.idleState(), false
	}
}
