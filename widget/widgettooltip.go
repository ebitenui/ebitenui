package widget

import (
	"image"
	"sync/atomic"
	"time"

	"github.com/ebitenui/ebitenui/input"
	"github.com/hajimehoshi/ebiten/v2"
)

type ToolTipPosition int

const (
	TOOLTIP_POS_CURSOR_FOLLOW ToolTipPosition = iota
	TOOLTIP_POS_CURSOR_STICKY
	TOOLTIP_POS_WIDGET
)

type ToolTipAnchor int

const (
	ANCHOR_START ToolTipAnchor = iota
	ANCHOR_MIDDLE
	ANCHOR_END
)

type ToolTipDirection int

type ToolTip struct {
	Position                ToolTipPosition
	WidgetOriginVertical    ToolTipAnchor
	WidgetOriginHorizontal  ToolTipAnchor
	ContentOriginVertical   ToolTipAnchor
	ContentOriginHorizontal ToolTipAnchor
	Delay                   time.Duration
	content                 *Container
	offset                  image.Point
	window                  *Window
	visible                 bool

	state          toolTipState
	toolTipUpdater ToolTipUpdater
}
type ToolTipOpt func(t *ToolTip)
type ToolTipOptions struct {
}

var ToolTipOpts ToolTipOptions

type toolTipState func(*Widget) (toolTipState, bool)

type ToolTipUpdater func(*Container)

func NewToolTip(opts ...ToolTipOpt) *ToolTip {
	t := &ToolTip{
		offset: image.Point{0, 0},
	}
	t.state = t.idleState()
	t.WidgetOriginHorizontal = ANCHOR_END
	t.WidgetOriginVertical = ANCHOR_END
	t.ContentOriginHorizontal = ANCHOR_END
	t.ContentOriginVertical = ANCHOR_START
	for _, o := range opts {
		o(t)
	}
	t.window = NewWindow(
		WindowOpts.CloseMode(NONE),
		WindowOpts.Contents(t.content),
	)

	return t
}

func (o ToolTipOptions) Content(c *Container) ToolTipOpt {
	return func(t *ToolTip) {
		t.content = c
	}
}

func (o ToolTipOptions) Offset(off image.Point) ToolTipOpt {
	return func(t *ToolTip) {
		t.offset = off
	}
}

func (o ToolTipOptions) WidgetOriginVertical(widgetOriginVertical ToolTipAnchor) ToolTipOpt {
	return func(t *ToolTip) {
		t.WidgetOriginVertical = widgetOriginVertical
	}
}

func (o ToolTipOptions) WidgetOriginHorizontal(widgetOriginHorizontal ToolTipAnchor) ToolTipOpt {
	return func(t *ToolTip) {
		t.WidgetOriginHorizontal = widgetOriginHorizontal
	}
}

func (o ToolTipOptions) ContentOriginVertical(contentOriginVertical ToolTipAnchor) ToolTipOpt {
	return func(t *ToolTip) {
		t.ContentOriginVertical = contentOriginVertical
	}
}

func (o ToolTipOptions) ContentOriginHorizontal(contentOriginHorizontal ToolTipAnchor) ToolTipOpt {
	return func(t *ToolTip) {
		t.ContentOriginHorizontal = contentOriginHorizontal
	}
}

func (o ToolTipOptions) Position(position ToolTipPosition) ToolTipOpt {
	return func(t *ToolTip) {
		t.Position = position
	}
}

func (o ToolTipOptions) Delay(d time.Duration) ToolTipOpt {
	return func(t *ToolTip) {
		t.Delay = d
	}
}

func (o ToolTipOptions) ToolTipUpdater(toolTipUpdater ToolTipUpdater) ToolTipOpt {
	return func(t *ToolTip) {
		t.toolTipUpdater = toolTipUpdater
	}
}
func (t *ToolTip) Render(parent *Widget, screen *ebiten.Image, def DeferredRenderFunc) {
	for {
		newState, rerun := t.state(parent)
		if newState != nil {
			t.state = newState
		}
		if !rerun {
			break
		}
	}
}

func (t *ToolTip) idleState() toolTipState {
	return func(parent *Widget) (toolTipState, bool) {
		if input.MouseButtonPressed(ebiten.MouseButtonLeft) ||
			input.MouseButtonPressed(ebiten.MouseButtonMiddle) ||
			input.MouseButtonPressed(ebiten.MouseButtonRight) {
			t.visible = false
			parent.FireToolTipEvent(t.window, false)
			return nil, false
		}

		x, y := input.CursorPosition()
		p := image.Point{x, y}
		if !p.In(parent.Rect) {
			return nil, false
		}

		if t.Delay <= 0 {
			return t.showingState(p), true
		}

		return t.armedState(p, nil, nil), true
	}
}

func (t *ToolTip) armedState(p image.Point, timer *time.Timer, expired *atomic.Value) toolTipState {
	return func(parent *Widget) (toolTipState, bool) {
		x, y := input.CursorPosition()
		cp := image.Point{x, y}

		if input.MouseButtonPressed(ebiten.MouseButtonLeft) ||
			input.MouseButtonPressed(ebiten.MouseButtonMiddle) ||
			input.MouseButtonPressed(ebiten.MouseButtonRight) ||
			!cp.In(parent.Rect) {
			t.visible = false
			parent.FireToolTipEvent(t.window, false)
			return t.idleState(), false
		}

		if timer != nil && expired.Load().(bool) {
			return t.showingState(cp), true
		}

		if timer == nil {
			expired = &atomic.Value{}
			expired.Store(false)

			timer = time.AfterFunc(t.Delay, func() {
				expired.Store(true)
			})

			return t.armedState(p, timer, expired), false
		}

		return nil, false
	}
}

func (t *ToolTip) showingState(p image.Point) toolTipState {
	return func(parent *Widget) (toolTipState, bool) {
		x, y := input.CursorPosition()
		cp := image.Point{x, y}
		if input.MouseButtonPressed(ebiten.MouseButtonLeft) ||
			input.MouseButtonPressed(ebiten.MouseButtonMiddle) ||
			input.MouseButtonPressed(ebiten.MouseButtonRight) ||
			!cp.In(parent.Rect) {
			t.visible = false
			parent.FireToolTipEvent(t.window, false)
			return t.idleState(), false
		}
		sx, sy := t.content.PreferredSize()

		if t.Position == TOOLTIP_POS_CURSOR_FOLLOW {
			p = cp
		} else if t.Position == TOOLTIP_POS_WIDGET {
			p = t.processWidgetPosition(parent)
			p = t.processContentPosition(p, sx, sy)
		}

		if t.toolTipUpdater != nil {
			t.toolTipUpdater(t.content)
		}

		r := image.Rect(0, 0, sx, sy)
		r = r.Add(p)
		r = r.Add(t.offset)
		t.window.SetLocation(r)
		t.content.SetLocation(r)
		t.content.RequestRelayout()
		if !t.visible {
			parent.FireToolTipEvent(t.window, true)
			t.visible = true
		}
		return t.showingState(p), false
	}
}

func (t *ToolTip) processWidgetPosition(parent *Widget) image.Point {
	p := image.Point{}
	widgetRect := parent.Rect
	if t.WidgetOriginVertical == ANCHOR_START {
		if t.WidgetOriginHorizontal == ANCHOR_START {
			p.X = widgetRect.Min.X
			p.Y = widgetRect.Min.Y
		} else if t.WidgetOriginHorizontal == ANCHOR_MIDDLE {
			p.X = widgetRect.Min.X + (widgetRect.Dx() / 2)
			p.Y = widgetRect.Min.Y
		} else {
			p.X = widgetRect.Max.X
			p.Y = widgetRect.Min.Y
		}
	} else if t.WidgetOriginVertical == ANCHOR_MIDDLE {
		if t.WidgetOriginHorizontal == ANCHOR_START {
			p.X = widgetRect.Min.X
			p.Y = widgetRect.Min.Y + (widgetRect.Dy() / 2)
		} else if t.WidgetOriginHorizontal == ANCHOR_MIDDLE {
			p.X = widgetRect.Min.X + (widgetRect.Dx() / 2)
			p.Y = widgetRect.Min.Y + (widgetRect.Dy() / 2)
		} else {
			p.X = widgetRect.Max.X
			p.Y = widgetRect.Min.Y + (widgetRect.Dy() / 2)
		}
	} else {
		if t.WidgetOriginHorizontal == ANCHOR_START {
			p.X = widgetRect.Min.X
			p.Y = widgetRect.Max.Y
		} else if t.WidgetOriginHorizontal == ANCHOR_MIDDLE {
			p.X = widgetRect.Min.X + (widgetRect.Dx() / 2)
			p.Y = widgetRect.Max.Y
		} else {
			p.X = widgetRect.Max.X
			p.Y = widgetRect.Max.Y
		}
	}
	return p
}

func (t *ToolTip) processContentPosition(p image.Point, sx int, sy int) image.Point {
	if t.ContentOriginVertical == ANCHOR_START {
		if t.ContentOriginHorizontal == ANCHOR_START {
			//Do nothing
		} else if t.ContentOriginHorizontal == ANCHOR_MIDDLE {
			p.X = p.X - (sx / 2)
		} else {
			p.X = p.X - sx
		}
	} else if t.ContentOriginVertical == ANCHOR_MIDDLE {
		if t.ContentOriginHorizontal == ANCHOR_START {
			p.Y = p.Y - (sy / 2)
		} else if t.ContentOriginHorizontal == ANCHOR_MIDDLE {
			p.X = p.X - (sx / 2)
			p.Y = p.Y - (sy / 2)
		} else {
			p.X = p.X - sx
			p.Y = p.Y - (sy / 2)
		}
	} else if t.ContentOriginVertical == ANCHOR_END {
		if t.ContentOriginHorizontal == ANCHOR_START {
			p.Y = p.Y - sy
		} else if t.ContentOriginHorizontal == ANCHOR_MIDDLE {
			p.X = p.X - (sx / 2)
			p.Y = p.Y - sy
		} else {
			p.X = p.X - sx
			p.Y = p.Y - sy
		}
	}
	return p
}
