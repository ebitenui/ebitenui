package widget

import (
	"image"
	"math"

	"github.com/blizzy78/ebitenui/event"
	"github.com/blizzy78/ebitenui/input"
	"github.com/hajimehoshi/ebiten"
)

type DragAndDrop struct {
	DroppedEvent *event.Event

	container            Locater
	contentsCreater      DragContentsCreater
	minDragStartDistance int

	state      dragAndDropState
	dragWidget HasWidget
}

type DragAndDropOpt func(d *DragAndDrop)

const DragAndDropOpts = dragAndDropOpts(true)

type dragAndDropOpts bool

type DragContentsCreater interface {
	Create(HasWidget, int, int) (HasWidget, interface{})
}

type DragContentsUpdater interface {
	Update(HasWidget, int, int, interface{})
}

type DragAndDropDroppedEventArgs struct {
	Source  HasWidget
	SourceX int
	SourceY int
	Target  HasWidget
	TargetX int
	TargetY int
	Data    interface{}
}

type DragAndDropDroppedHandlerFunc func(args *DragAndDropDroppedEventArgs)

type dragAndDropState func(*ebiten.Image, DeferredRenderFunc) (dragAndDropState, bool)

func NewDragAndDrop(opts ...DragAndDropOpt) *DragAndDrop {
	d := &DragAndDrop{
		DroppedEvent: &event.Event{},

		minDragStartDistance: 15,
	}
	d.state = d.idleState()

	for _, o := range opts {
		o(d)
	}

	return d
}

func (o dragAndDropOpts) Container(c Locater) DragAndDropOpt {
	return func(d *DragAndDrop) {
		d.container = c
	}
}

func (o dragAndDropOpts) ContentsCreater(c DragContentsCreater) DragAndDropOpt {
	return func(d *DragAndDrop) {
		d.contentsCreater = c
	}
}

func (o dragAndDropOpts) MinDragStartDistance(d int) DragAndDropOpt {
	return func(dnd *DragAndDrop) {
		dnd.minDragStartDistance = d
	}
}

func (o dragAndDropOpts) DroppedHandler(f DragAndDropDroppedHandlerFunc) DragAndDropOpt {
	return func(d *DragAndDrop) {
		d.DroppedEvent.AddHandler(func(args interface{}) {
			f(args.(*DragAndDropDroppedEventArgs))
		})
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

func (d *DragAndDrop) Render(screen *ebiten.Image, def DeferredRenderFunc) {
	for {
		var rerun bool
		d.state, rerun = d.state(screen, def)
		if !rerun {
			break
		}
	}
}

func (d *DragAndDrop) idleState() dragAndDropState {
	return func(screen *ebiten.Image, def DeferredRenderFunc) (dragAndDropState, bool) {
		d.dragWidget = nil

		if !input.MouseButtonJustPressed(ebiten.MouseButtonLeft) {
			return d.idleState(), false
		}

		x, y := input.CursorPosition()
		srcWidget := d.container.WidgetAt(x, y)
		if srcWidget == nil {
			return d.idleState(), false
		}

		return d.dragArmedState(srcWidget, x, y), true
	}
}

func (d *DragAndDrop) dragArmedState(srcWidget HasWidget, srcX int, srcY int) dragAndDropState {
	return func(screen *ebiten.Image, def DeferredRenderFunc) (dragAndDropState, bool) {
		if !input.MouseButtonPressed(ebiten.MouseButtonLeft) {
			return d.idleState(), false
		}

		x, y := input.CursorPosition()
		dx, dy := math.Abs(float64(x-srcX)), math.Abs(float64(y-srcY))
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist < float64(d.minDragStartDistance) {
			return d.dragArmedState(srcWidget, srcX, srcY), false
		}

		return d.draggingState(srcWidget, srcX, srcY, nil, nil), true
	}
}

func (d *DragAndDrop) draggingState(srcWidget HasWidget, srcX int, srcY int, dragWidget HasWidget, dragData interface{}) dragAndDropState {
	return func(screen *ebiten.Image, def DeferredRenderFunc) (dragAndDropState, bool) {
		x, y := input.CursorPosition()
		w := d.container.WidgetAt(x, y)

		if !input.MouseButtonPressed(ebiten.MouseButtonLeft) {
			return d.droppingState(srcWidget, srcX, srcY, w, x, y, dragData), true
		}

		if dragWidget == nil {
			dragWidget, dragData = d.contentsCreater.Create(srcWidget, srcX, srcY)

			if dragWidget == nil {
				return d.idleState(), false
			}
		}

		defer func() {
			d.dragWidget = dragWidget
		}()

		if u, ok := d.contentsCreater.(DragContentsUpdater); ok {
			u.Update(w, x, y, dragData)
		}

		sx, sy := dragWidget.(PreferredSizer).PreferredSize()
		r := image.Rect(0, 0, sx, sy)
		r = r.Add(image.Point{x, y})
		r = r.Sub(image.Point{sx / 2, sy / 2})
		dragWidget.(Locateable).SetLocation(r)
		if rl, ok := dragWidget.(Relayoutable); ok {
			rl.RequestRelayout()
		}
		dragWidget.(Renderer).Render(screen, def)

		return d.draggingState(srcWidget, srcX, srcY, dragWidget, dragData), false
	}
}

func (d *DragAndDrop) droppingState(srcWidget HasWidget, srcX int, srcY int, targetWidget HasWidget, x int, y int, dragData interface{}) dragAndDropState {
	return func(screen *ebiten.Image, def DeferredRenderFunc) (dragAndDropState, bool) {
		d.DroppedEvent.Fire(&DragAndDropDroppedEventArgs{
			Source:  srcWidget,
			SourceX: srcX,
			SourceY: srcY,
			Target:  targetWidget,
			TargetX: x,
			TargetY: y,
			Data:    dragData,
		})

		return d.idleState(), false
	}
}
