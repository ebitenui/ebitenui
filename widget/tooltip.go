package widget

import (
	"image"
	"image/color"
	"sync/atomic"
	"time"

	e_image "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/input"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

type ToolTipPosition int

const (
	// The tooltip will follow the cursor around while visible
	TOOLTIP_POS_CURSOR_FOLLOW ToolTipPosition = iota
	// The tooltip will stick to where the cursor was when the tooltip was made visible
	TOOLTIP_POS_CURSOR_STICKY
	// The tooltip will display based on the Widget and Content anchor settings.
	// It defaults to opening right aligned and directly under the widget.
	TOOLTIP_POS_WIDGET
)

type ToolTipAnchor int

const (
	// Anchor at the start of the element
	ANCHOR_START ToolTipAnchor = iota
	// Anchor in the middle of the element
	ANCHOR_MIDDLE
	// Anchor at the end of the element
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
	Offset                  image.Point
	content                 *Container
	window                  *Window
	visible                 bool

	state          toolTipState
	ToolTipUpdater ToolTipUpdater
}
type ToolTipOpt func(t *ToolTip)
type ToolTipOptions struct {
}

var ToolTipOpts ToolTipOptions

type toolTipState func(*Widget) (toolTipState, bool)

type ToolTipUpdater func(*Container)

// Create a new Tooltip. This method allows you to specify
// every aspect of the displayed tooltip's content.
func NewToolTip(opts ...ToolTipOpt) *ToolTip {
	t := &ToolTip{
		Offset: image.Point{0, 0},
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

// Create a new Text Tooltip with the following defaults:
//   - ProcessBBCode = true
//   - Padding = Top/Bottom: 5px Left/Right: 10px
//   - Delay = 800ms
//   - Offset = 0, 20
func NewTextToolTip(label string, face font.Face, color color.Color, background *e_image.NineSlice) *ToolTip {
	c := NewContainer(
		ContainerOpts.BackgroundImage(background),
		ContainerOpts.AutoDisableChildren(),
		ContainerOpts.Layout(NewAnchorLayout(AnchorLayoutOpts.Padding(Insets{
			Top:    5,
			Bottom: 5,
			Left:   10,
			Right:  10,
		}))),
	)

	c.AddChild(NewText(TextOpts.ProcessBBCode(true), TextOpts.Text(label, face, color)))

	return NewToolTip(
		ToolTipOpts.Content(c),
		ToolTipOpts.Delay(800*time.Millisecond),
		ToolTipOpts.Offset(image.Point{0, 20}),
	)
}

// The container to be displayed
func (o ToolTipOptions) Content(c *Container) ToolTipOpt {
	return func(t *ToolTip) {
		t.content = c
	}
}

// The X/Y offsets from the Tooltip anchor point
func (o ToolTipOptions) Offset(off image.Point) ToolTipOpt {
	return func(t *ToolTip) {
		t.Offset = off
	}
}

// The vertical position of the anchor on the widget. Only used when Postion = WIDGET
func (o ToolTipOptions) WidgetOriginVertical(widgetOriginVertical ToolTipAnchor) ToolTipOpt {
	return func(t *ToolTip) {
		t.WidgetOriginVertical = widgetOriginVertical
	}
}

// The horizontal position of the anchor on the widget. Only used when Postion = WIDGET
func (o ToolTipOptions) WidgetOriginHorizontal(widgetOriginHorizontal ToolTipAnchor) ToolTipOpt {
	return func(t *ToolTip) {
		t.WidgetOriginHorizontal = widgetOriginHorizontal
	}
}

// The vertical position of the anchor on the tooltip. Only used when Postion = WIDGET
func (o ToolTipOptions) ContentOriginVertical(contentOriginVertical ToolTipAnchor) ToolTipOpt {
	return func(t *ToolTip) {
		t.ContentOriginVertical = contentOriginVertical
	}
}

// The horizontal position of the anchor on the tooltip. Only used when Postion = WIDGET
func (o ToolTipOptions) ContentOriginHorizontal(contentOriginHorizontal ToolTipAnchor) ToolTipOpt {
	return func(t *ToolTip) {
		t.ContentOriginHorizontal = contentOriginHorizontal
	}
}

// Where to display the tooltip
func (o ToolTipOptions) Position(position ToolTipPosition) ToolTipOpt {
	return func(t *ToolTip) {
		t.Position = position
	}
}

// How long to wait before displaying the tooltip
func (o ToolTipOptions) Delay(d time.Duration) ToolTipOpt {
	return func(t *ToolTip) {
		t.Delay = d
	}
}

// A method that is called every draw call that the tooltip is visible.
// This allows you to hook into the draw loop to update the tooltip if necessary
func (o ToolTipOptions) ToolTipUpdater(toolTipUpdater ToolTipUpdater) ToolTipOpt {
	return func(t *ToolTip) {
		t.ToolTipUpdater = toolTipUpdater
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

		if t.ToolTipUpdater != nil {
			t.ToolTipUpdater(t.content)
		}

		r := image.Rect(0, 0, sx, sy)
		r = r.Add(p)
		r = r.Add(t.Offset)
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
