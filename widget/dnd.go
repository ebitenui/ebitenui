package widget

import (
	"image"
	"math"

	"github.com/ebitenui/ebitenui/input"
	"github.com/hajimehoshi/ebiten/v2"
)

type DragAndDropAnchor int

const (
	// Anchor at the start of the element
	DND_ANCHOR_START DragAndDropAnchor = iota
	// Anchor in the middle of the element
	DND_ANCHOR_MIDDLE
	// Anchor at the end of the element
	DND_ANCHOR_END
)

type DragAndDrop struct {
	ContentsOriginVertical   DragAndDropAnchor
	ContentsOriginHorizontal DragAndDropAnchor
	Offset                   image.Point

	AvailableDropTargets []HasWidget
	contentsCreater      DragContentsCreater
	minDragStartDistance int
	state                dragAndDropState
	dragWidget           *Container
	window               *Window
	dndTriggered         bool
	dndStopped           bool
	dragDisabled         bool
}

type DragAndDropOpt func(d *DragAndDrop)

type DragAndDropOptions struct {
}

var DragAndDropOpts DragAndDropOptions

type DragContentsCreater interface {
	Create(HasWidget) (*Container, interface{})
}

type DragContentsUpdater interface {
	// arg1 - isDroppable
	// arg2 - HasWidget if droppable
	// arg3 - DragData
	Update(bool, HasWidget, interface{})
}

type DragContentsEnder interface {
	// arg1 - Drop was successful
	// arg2 - Source Widget
	// arg3 - DragData
	EndDrag(bool, HasWidget, interface{})
}

type dragAndDropState func(HasWidget) (dragAndDropState, bool)

func NewDragAndDrop(opts ...DragAndDropOpt) *DragAndDrop {
	d := &DragAndDrop{
		minDragStartDistance:     15,
		ContentsOriginVertical:   DND_ANCHOR_MIDDLE,
		ContentsOriginHorizontal: DND_ANCHOR_MIDDLE,
		Offset:                   image.Point{0, 0},
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

// The minimum distance in pixels a user must drag their cursor to display the dragged element.
//
//	Optional - Defaults to 15 pixels
func (o DragAndDropOptions) MinDragStartDistance(d int) DragAndDropOpt {
	return func(dnd *DragAndDrop) {
		dnd.minDragStartDistance = d
	}
}

// The vertical position of the anchor on the tooltip.
//
//	Optional - Defaults to DND_ANCHOR_MIDDLE
func (o DragAndDropOptions) ContentsOriginVertical(contentsOriginVertical DragAndDropAnchor) DragAndDropOpt {
	return func(t *DragAndDrop) {
		t.ContentsOriginVertical = contentsOriginVertical
	}
}

// The horizontal position of the anchor on the tooltip.
//
//	Optional - Defaults to DND_ANCHOR_MIDDLE
func (o DragAndDropOptions) ContentsOriginHorizontal(contentsOriginHorizontal DragAndDropAnchor) DragAndDropOpt {
	return func(t *DragAndDrop) {
		t.ContentsOriginHorizontal = contentsOriginHorizontal
	}
}

// The X/Y offsets from the Tooltip anchor point
func (o DragAndDropOptions) Offset(off image.Point) DragAndDropOpt {
	return func(t *DragAndDrop) {
		t.Offset = off
	}
}

// Disable Drag to start Drag and Drop.
// You may use the "StartDrag()" method on this object or the
// to begin the drag operation.
//
//	Expected use-case: click to pick up, click to drop.
func (o DragAndDropOptions) DisableDrag() DragAndDropOpt {
	return func(t *DragAndDrop) {
		t.dragDisabled = true
	}
}

// To avoid conflicting with dragging, if you trigger it on left click you should put the trigger in the button released event
func (d *DragAndDrop) StartDrag() {
	d.dndTriggered = true
}

func (d *DragAndDrop) StopDrag() {
	d.dndStopped = true
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

func (d *DragAndDrop) Render(parent HasWidget, screen *ebiten.Image, def DeferredRenderFunc) {
	newState, _ := d.state(parent)
	if newState != nil {
		d.state = newState
	}
}

func (d *DragAndDrop) idleState() dragAndDropState {
	return func(parent HasWidget) (dragAndDropState, bool) {
		d.dragWidget = nil
		if (!input.MouseButtonJustPressed(ebiten.MouseButtonLeft) && !d.dndTriggered) || d.dndStopped {
			d.dndStopped = false
			if d.window != nil {
				parent.GetWidget().FireDragAndDropEvent(d.window, false, d)
				d.window = nil
			}
			return nil, false
		}

		x, y := input.CursorPosition()
		p := image.Point{x, y}
		if !p.In(parent.GetWidget().Rect) && !d.dndTriggered {
			return nil, false
		}
		if !parent.GetWidget().EffectiveInputLayer().ActiveFor(x, y, input.LayerEventTypeAny) {
			return nil, false
		}

		return d.dragArmedState(x, y), true
	}
}

func (d *DragAndDrop) dragArmedState(srcX int, srcY int) dragAndDropState {
	return func(_ HasWidget) (dragAndDropState, bool) {
		if !input.MouseButtonPressed(ebiten.MouseButtonLeft) && !d.dndTriggered {
			return d.idleState(), false
		}
		if !d.dndTriggered {
			x, y := input.CursorPosition()
			dx, dy := math.Abs(float64(x-srcX)), math.Abs(float64(y-srcY))
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist < float64(d.minDragStartDistance) || d.dragDisabled {
				return nil, false
			}
		}
		return d.draggingState(srcX, srcY, nil, nil, !d.dndTriggered), true
	}
}

func (d *DragAndDrop) draggingState(srcX int, srcY int, dragWidget *Container, dragData interface{}, mousePressed bool) dragAndDropState {
	return func(parent HasWidget) (dragAndDropState, bool) {
		x, y := input.CursorPosition()

		d.dndTriggered = false

		if input.MouseButtonPressed(ebiten.MouseButtonLeft) != mousePressed {
			return d.droppingState(srcX, srcY, x, y, dragData), true
		}

		if dragWidget == nil {
			dragWidget, dragData = d.contentsCreater.Create(parent)
			if dragWidget == nil {
				return d.idleState(), false
			}
			d.window = NewWindow(
				WindowOpts.CloseMode(NONE),
				WindowOpts.Contents(dragWidget),
			)
			parent.GetWidget().FireDragAndDropEvent(d.window, true, d)
		}

		defer func() {
			d.dragWidget = dragWidget
		}()

		if u, ok := d.contentsCreater.(DragContentsUpdater); ok {
			droppable := false
			var element HasWidget

			if !input.KeyPressed(ebiten.KeyEscape) && !d.dndStopped {
				p := image.Point{x, y}
				args := &DragAndDropDroppedEventArgs{
					Source:  parent,
					SourceX: srcX,
					SourceY: srcY,
					TargetX: x,
					TargetY: y,
					Data:    dragData,
				}
				for _, target := range d.AvailableDropTargets {
					if target.GetWidget().Visibility == Visibility_Hide {
						continue
					}
					if !p.In(target.GetWidget().Rect) {
						continue
					}
					if !target.GetWidget().EffectiveInputLayer().ActiveFor(x, y, input.LayerEventTypeAny) {
						continue
					}
					if target.GetWidget().canDrop(args) {
						droppable = true
						element = target
						break
					}
				}
			}
			u.Update(droppable, element, dragData)
		}

		if input.KeyPressed(ebiten.KeyEscape) || d.dndStopped {
			if e, ok := d.contentsCreater.(DragContentsEnder); ok {
				e.EndDrag(false, parent, dragData)
			}

			return d.idleState(), false
		}

		sx, sy := dragWidget.PreferredSize()
		r := image.Rect(0, 0, sx, sy)
		r = r.Add(d.processContentsPosition(image.Point{x, y}, sx, sy))
		r = r.Add(d.Offset)
		d.window.SetLocation(r)
		dragWidget.SetLocation(r)

		return d.draggingState(srcX, srcY, dragWidget, dragData, mousePressed), false
	}
}

func (d *DragAndDrop) droppingState(srcX int, srcY int, x int, y int, dragData interface{}) dragAndDropState {
	return func(parent HasWidget) (dragAndDropState, bool) {
		args := &DragAndDropDroppedEventArgs{
			Source:  parent,
			SourceX: srcX,
			SourceY: srcY,
			TargetX: x,
			TargetY: y,
			Data:    dragData,
		}
		p := image.Point{x, y}
		dropSuccessful := false
		for _, target := range d.AvailableDropTargets {
			if target.GetWidget().Visibility == Visibility_Hide {
				continue
			}
			if !p.In(target.GetWidget().Rect) {
				continue
			}
			if !target.GetWidget().EffectiveInputLayer().ActiveFor(x, y, input.LayerEventTypeAny) {
				continue
			}
			if target.GetWidget().canDrop(args) {
				if target.GetWidget().drop != nil {
					args.Target = target
					target.GetWidget().drop(args)
					dropSuccessful = true
				}
				break
			}
		}

		if e, ok := d.contentsCreater.(DragContentsEnder); ok {
			e.EndDrag(dropSuccessful, parent, dragData)
		}

		d.dndStopped = false
		if d.window != nil {
			parent.GetWidget().FireDragAndDropEvent(d.window, false, d)
			d.window = nil
		}

		return d.idleState(), false
	}
}
func (d *DragAndDrop) processContentsPosition(p image.Point, sx int, sy int) image.Point {
	if d.ContentsOriginVertical == DND_ANCHOR_START {
		if d.ContentsOriginHorizontal == DND_ANCHOR_START {
			//Do nothing
		} else if d.ContentsOriginHorizontal == DND_ANCHOR_MIDDLE {
			p.X = p.X - (sx / 2)
		} else {
			p.X = p.X - sx
		}
	} else if d.ContentsOriginVertical == DND_ANCHOR_MIDDLE {
		if d.ContentsOriginHorizontal == DND_ANCHOR_START {
			p.Y = p.Y - (sy / 2)
		} else if d.ContentsOriginHorizontal == DND_ANCHOR_MIDDLE {
			p.X = p.X - (sx / 2)
			p.Y = p.Y - (sy / 2)
		} else {
			p.X = p.X - sx
			p.Y = p.Y - (sy / 2)
		}
	} else if d.ContentsOriginVertical == DND_ANCHOR_END {
		if d.ContentsOriginHorizontal == DND_ANCHOR_START {
			p.Y = p.Y - sy
		} else if d.ContentsOriginHorizontal == DND_ANCHOR_MIDDLE {
			p.X = p.X - (sx / 2)
			p.Y = p.Y - sy
		} else {
			p.X = p.X - sx
			p.Y = p.Y - sy
		}
	}
	return p
}
