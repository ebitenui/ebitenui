package widget

import (
	"image"
	"image/color"
	"sync/atomic"
	"time"

	e_image "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/input"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type ToolTipPosition int

const (
	// The tooltip will follow the cursor around while visible.
	TOOLTIP_POS_CURSOR_FOLLOW ToolTipPosition = iota
	// The tooltip will stick to where the cursor was when the tooltip was made visible.
	TOOLTIP_POS_CURSOR_STICKY
	// The tooltip will display based on the Widget and Content anchor settings.
	// It defaults to opening right aligned and directly under the widget.
	TOOLTIP_POS_WIDGET
	// The tooltip will display based on x/y (offset is required)
	TOOLTIP_POS_ABSOLUTE
	// The tooltip will display based on the Widget and Content anchor settings.
	// It defaults to opening right aligned and directly under the x: 0, y: 0.
	TOOLTIP_POS_SCREEN
)

type ToolTipAnchor int

const (
	// Anchor at the start of the element.
	TOOLTIP_ANCHOR_START ToolTipAnchor = iota
	// Anchor in the middle of the element.
	TOOLTIP_ANCHOR_MIDDLE
	// Anchor at the end of the element.
	TOOLTIP_ANCHOR_END
)

type ToolTipDirection int

type ToolTip struct {
	Position ToolTipPosition
	// WidgetOriginVertical was renamed to AnchorOriginVertical to make the it more generic and reuse it for TOOLTIP_POS_SCREEN
	AnchorOriginVertical ToolTipAnchor
	// WidgetOriginHorizontal was renamed to AnchorOriginHorizontal to make the it more generic and reuse it for TOOLTIP_POS_SCREEN
	AnchorOriginHorizontal  ToolTipAnchor
	ContentOriginVertical   ToolTipAnchor
	ContentOriginHorizontal ToolTipAnchor
	Delay                   time.Duration
	Offset                  image.Point
	content                 Containerer
	window                  *Window
	visible                 bool

	state          toolTipState
	ToolTipUpdater ToolTipUpdater
}
type ToolTipOpt func(t *ToolTip)
type ToolTipOptions struct {
}

var ToolTipOpts ToolTipOptions

type toolTipState func(*Widget) toolTipState

type ToolTipUpdater func(Containerer)

// Create a new Tooltip. This method allows you to specify
// every aspect of the displayed tooltip's content.
func NewToolTip(opts ...ToolTipOpt) *ToolTip {
	t := &ToolTip{
		Offset: image.Point{10, 10},
	}
	t.state = t.idleState()
	t.AnchorOriginHorizontal = TOOLTIP_ANCHOR_END
	t.AnchorOriginVertical = TOOLTIP_ANCHOR_END
	t.ContentOriginHorizontal = TOOLTIP_ANCHOR_END
	t.ContentOriginVertical = TOOLTIP_ANCHOR_START
	for _, o := range opts {
		o(t)
	}

	t.Validate()

	t.window = NewWindow(
		WindowOpts.CloseMode(NONE),
		WindowOpts.Contents(t.content),
	)

	return t
}

func (t *ToolTip) Validate() {
	if t.content == nil {
		panic("ToolTip: Contents is required.")
	}
}

// Create a new Text Tooltip with the following defaults:
//   - ProcessBBCode = true
//   - Padding = Top/Bottom: 5px Left/Right: 10px
//   - Delay = 800ms
//   - Offset = 0, 20
//   - ContentOriginHorizontal = TOOLTIP_ANCHOR_END
//   - ContentOriginVertical = TOOLTIP_ANCHOR_START
func NewTextToolTip(label string, face *text.Face, color color.Color, background *e_image.NineSlice) *ToolTip {
	if color == nil {
		panic("TextToolTip: color is required.")
	}
	if face == nil {
		panic("TextToolTip: face is required.")
	}

	c := NewContainer(
		ContainerOpts.BackgroundImage(background),
		ContainerOpts.AutoDisableChildren(),
		ContainerOpts.Layout(NewAnchorLayout(AnchorLayoutOpts.Padding(&Insets{
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
		ToolTipOpts.ContentOriginHorizontal(TOOLTIP_ANCHOR_START),
		ToolTipOpts.ContentOriginVertical(TOOLTIP_ANCHOR_START),
	)
}

// The container to be displayed.
func (o ToolTipOptions) Content(c Containerer) ToolTipOpt {
	return func(t *ToolTip) {
		t.content = c
	}
}

// The X/Y offsets from the Tooltip anchor point.
func (o ToolTipOptions) Offset(off image.Point) ToolTipOpt {
	return func(t *ToolTip) {
		t.Offset = off
	}
}

// The vertical position of the anchor on the widget. Only used when Postion = WIDGET.
func (o ToolTipOptions) AnchorOriginVertical(anchorOriginVertical ToolTipAnchor) ToolTipOpt {
	return func(t *ToolTip) {
		t.AnchorOriginVertical = anchorOriginVertical
	}
}

// The horizontal position of the anchor on the widget. Only used when Postion = WIDGET.
func (o ToolTipOptions) AnchorOriginHorizontal(anchorOriginHorizontal ToolTipAnchor) ToolTipOpt {
	return func(t *ToolTip) {
		t.AnchorOriginHorizontal = anchorOriginHorizontal
	}
}

// The vertical position of the anchor on the tooltip.
func (o ToolTipOptions) ContentOriginVertical(contentOriginVertical ToolTipAnchor) ToolTipOpt {
	return func(t *ToolTip) {
		t.ContentOriginVertical = contentOriginVertical
	}
}

// The horizontal position of the anchor on the tooltip.
func (o ToolTipOptions) ContentOriginHorizontal(contentOriginHorizontal ToolTipAnchor) ToolTipOpt {
	return func(t *ToolTip) {
		t.ContentOriginHorizontal = contentOriginHorizontal
	}
}

// Where to display the tooltip.
func (o ToolTipOptions) Position(position ToolTipPosition) ToolTipOpt {
	return func(t *ToolTip) {
		t.Position = position
	}
}

// How long to wait before displaying the tooltip.
func (o ToolTipOptions) Delay(d time.Duration) ToolTipOpt {
	return func(t *ToolTip) {
		t.Delay = d
	}
}

// A method that is called every draw call that the tooltip is visible.
// This allows you to hook into the draw loop to update the tooltip if necessary.
func (o ToolTipOptions) ToolTipUpdater(toolTipUpdater ToolTipUpdater) ToolTipOpt {
	return func(t *ToolTip) {
		t.ToolTipUpdater = toolTipUpdater
	}
}

func (t *ToolTip) Update(parent *Widget) {
	newState := t.state(parent)
	if newState != nil {
		t.state = newState
	}
}

func (t *ToolTip) idleState() toolTipState {
	return func(parent *Widget) toolTipState {
		if input.MouseButtonPressed(ebiten.MouseButtonLeft) ||
			input.MouseButtonPressed(ebiten.MouseButtonMiddle) ||
			input.MouseButtonPressed(ebiten.MouseButtonRight) ||
			!parent.IsVisible() {
			t.visible = false
			parent.FireToolTipEvent(t.window, false)
			return nil
		}

		x, y := input.CursorPosition()
		p := image.Point{x, y}
		if !p.In(parent.Rect) {
			return nil
		}
		if !parent.EffectiveInputLayer().ActiveFor(x, y, input.LayerEventTypeAny) {
			return nil
		}

		if t.Delay <= 0 {
			return t.showingState(p)
		}

		return t.armedState(p, nil, nil)
	}
}

func (t *ToolTip) armedState(p image.Point, timer *time.Timer, expired *atomic.Value) toolTipState {
	return func(parent *Widget) toolTipState {
		x, y := input.CursorPosition()
		cp := image.Point{x, y}

		if input.MouseButtonPressed(ebiten.MouseButtonLeft) ||
			input.MouseButtonPressed(ebiten.MouseButtonMiddle) ||
			input.MouseButtonPressed(ebiten.MouseButtonRight) ||
			!cp.In(parent.Rect) {
			t.visible = false
			parent.FireToolTipEvent(t.window, false)
			return t.idleState()
		}
		if timer != nil {
			if isExpired, _ := expired.Load().(bool); isExpired {
				return t.showingState(cp)
			}
		}

		if timer == nil {
			expired = &atomic.Value{}
			expired.Store(false)

			timer = time.AfterFunc(t.Delay, func() {
				expired.Store(true)
			})

			return t.armedState(p, timer, expired)
		}

		return nil
	}
}

func (t *ToolTip) showingState(p image.Point) toolTipState {
	return func(parent *Widget) toolTipState {
		x, y := input.CursorPosition()
		cp := image.Point{x, y}
		if input.MouseButtonPressed(ebiten.MouseButtonLeft) ||
			input.MouseButtonPressed(ebiten.MouseButtonMiddle) ||
			input.MouseButtonPressed(ebiten.MouseButtonRight) ||
			!cp.In(parent.Rect) {
			t.visible = false
			parent.FireToolTipEvent(t.window, false)
			return t.idleState()
		}
		sx, sy := t.content.PreferredSize()

		position := p
		switch t.Position {
		case TOOLTIP_POS_CURSOR_FOLLOW:
			position = cp
		case TOOLTIP_POS_WIDGET:
			position = t.processWidgetPosition(parent.Rect)
		case TOOLTIP_POS_ABSOLUTE:
			position = image.Point{}
		case TOOLTIP_POS_SCREEN:
			position = t.processScreenPosition()
		}
		position = position.Add(t.Offset)
		position = t.processContentPosition(position, sx, sy, parent.Rect)

		if t.ToolTipUpdater != nil {
			t.ToolTipUpdater(t.content)
		}

		r := image.Rect(0, 0, sx, sy)
		r = r.Add(position)
		t.window.SetLocation(r)
		t.content.SetLocation(r)
		if !t.visible {
			parent.FireToolTipEvent(t.window, true)
			t.visible = true
		}
		return t.showingState(p)
	}
}

func (t *ToolTip) processWidgetPosition(widgetRect image.Rectangle) image.Point {
	p := image.Point{}
	switch t.AnchorOriginVertical {
	case TOOLTIP_ANCHOR_START:
		switch t.AnchorOriginHorizontal {
		case TOOLTIP_ANCHOR_START:
			p.X = widgetRect.Min.X
			p.Y = widgetRect.Min.Y
		case TOOLTIP_ANCHOR_MIDDLE:
			p.X = widgetRect.Min.X + (widgetRect.Dx() / 2)
			p.Y = widgetRect.Min.Y
		case TOOLTIP_ANCHOR_END:
			p.X = widgetRect.Max.X
			p.Y = widgetRect.Min.Y
		}
	case TOOLTIP_ANCHOR_MIDDLE:
		switch t.AnchorOriginHorizontal {
		case TOOLTIP_ANCHOR_START:
			p.X = widgetRect.Min.X
			p.Y = widgetRect.Min.Y + (widgetRect.Dy() / 2)
		case TOOLTIP_ANCHOR_MIDDLE:
			p.X = widgetRect.Min.X + (widgetRect.Dx() / 2)
			p.Y = widgetRect.Min.Y + (widgetRect.Dy() / 2)
		case TOOLTIP_ANCHOR_END:
			p.X = widgetRect.Max.X
			p.Y = widgetRect.Min.Y + (widgetRect.Dy() / 2)
		}
	case TOOLTIP_ANCHOR_END:
		switch t.AnchorOriginHorizontal {
		case TOOLTIP_ANCHOR_START:
			p.X = widgetRect.Min.X
			p.Y = widgetRect.Max.Y
		case TOOLTIP_ANCHOR_MIDDLE:
			p.X = widgetRect.Min.X + (widgetRect.Dx() / 2)
			p.Y = widgetRect.Max.Y
		case TOOLTIP_ANCHOR_END:
			p.X = widgetRect.Max.X
			p.Y = widgetRect.Max.Y
		}
	}
	return p
}

func (t *ToolTip) processScreenPosition() image.Point {
	windowSize := input.GetWindowSize()
	p := image.Point{}
	switch t.AnchorOriginHorizontal {
	case TOOLTIP_ANCHOR_START:
		p.X = 0
	case TOOLTIP_ANCHOR_MIDDLE:
		p.X = windowSize.X / 2
	case TOOLTIP_ANCHOR_END:
		p.X = windowSize.X
	}
	switch t.AnchorOriginVertical {
	case TOOLTIP_ANCHOR_START:
		p.Y = 0
	case TOOLTIP_ANCHOR_MIDDLE:
		p.Y = windowSize.Y / 2
	case TOOLTIP_ANCHOR_END:
		p.Y = windowSize.Y
	}
	return p
}

func (t *ToolTip) processContentPosition(p image.Point, sx, sy int, widgetRect image.Rectangle) image.Point {
	result := processContentPositionWorker(p, sx, sy, t.ContentOriginHorizontal, t.ContentOriginVertical)
	windowSize := input.GetWindowSize()
	horizontalAnchor := t.ContentOriginHorizontal
	if result.X+sx > windowSize.X {
		horizontalAnchor = TOOLTIP_ANCHOR_END
		if t.Position == TOOLTIP_POS_WIDGET {
			p.X = widgetRect.Min.X
		}
		p.X -= 2 * t.Offset.X
		result = processContentPositionWorker(p, sx, sy, horizontalAnchor, t.ContentOriginVertical)
	} else if result.X < 0 {
		p.X -= 2 * t.Offset.X
		horizontalAnchor = TOOLTIP_ANCHOR_START
		result = processContentPositionWorker(p, sx, sy, horizontalAnchor, t.ContentOriginVertical)
	}

	if result.Y+sy > windowSize.Y {
		if t.Position == TOOLTIP_POS_WIDGET {
			p.Y = widgetRect.Min.Y
		}
		p.Y -= 2 * t.Offset.Y
		result = processContentPositionWorker(p, sx, sy, horizontalAnchor, TOOLTIP_ANCHOR_END)
	} else if result.Y < 0 {
		p.Y -= 2 * t.Offset.Y
		result = processContentPositionWorker(p, sx, sy, horizontalAnchor, TOOLTIP_ANCHOR_START)
	}

	return result
}

func processContentPositionWorker(p image.Point, sx int, sy int, originHorizontal ToolTipAnchor, originVertical ToolTipAnchor) image.Point {
	switch originVertical {
	case TOOLTIP_ANCHOR_START:
		switch originHorizontal {
		case TOOLTIP_ANCHOR_START:
			// Do nothing
		case TOOLTIP_ANCHOR_MIDDLE:
			p.X -= (sx / 2)
		case TOOLTIP_ANCHOR_END:
			p.X -= sx
		}
	case TOOLTIP_ANCHOR_MIDDLE:
		switch originHorizontal {
		case TOOLTIP_ANCHOR_START:
			p.Y -= (sy / 2)
		case TOOLTIP_ANCHOR_MIDDLE:
			p.X -= (sx / 2)
			p.Y -= (sy / 2)
		case TOOLTIP_ANCHOR_END:
			p.X -= sx
			p.Y -= (sy / 2)
		}
	case TOOLTIP_ANCHOR_END:
		switch originHorizontal {
		case TOOLTIP_ANCHOR_START:
			p.Y -= sy
		case TOOLTIP_ANCHOR_MIDDLE:
			p.X -= (sx / 2)
			p.Y -= sy
		case TOOLTIP_ANCHOR_END:
			p.X -= sx
			p.Y -= sy
		}
	}
	return p
}
