package widget

import (
	img "image"
	"time"

	"github.com/blizzy78/ebitenui/image"
	"github.com/blizzy78/ebitenui/input"
	"github.com/hajimehoshi/ebiten"
)

type ToolTip struct {
	textOpts  []TextOpt
	container WidgetLocator
	image     *image.NineSlice
	padding   Insets
	offset    img.Point
	sticky    bool
	delay     time.Duration

	init          *MultiOnce
	tipContainer  *Container
	tipText       *Text
	lastTipWidget HasWidget
	timer         *time.Timer
	doRender      bool
	doRenderReset bool
}

type ToolTipOpt func(t *ToolTip)

const ToolTipOpts = toolTipOpts(true)

type toolTipOpts bool

func NewToolTip(opts ...ToolTipOpt) *ToolTip {
	t := &ToolTip{
		offset: img.Point{0, 20},
		sticky: true,
		delay:  800 * time.Millisecond,

		init: &MultiOnce{},
	}

	t.init.Append(t.createWidget)

	for _, o := range opts {
		o(t)
	}

	return t
}

func (o toolTipOpts) WithTextOpts(opts ...TextOpt) ToolTipOpt {
	return func(t *ToolTip) {
		t.textOpts = append(t.textOpts, opts...)
	}
}

func (o toolTipOpts) WithContainer(c *Container) ToolTipOpt {
	return func(t *ToolTip) {
		t.container = c
	}
}

func (o toolTipOpts) WithPadding(i Insets) ToolTipOpt {
	return func(t *ToolTip) {
		t.padding = i
	}
}

func (o toolTipOpts) WithImage(i *image.NineSlice) ToolTipOpt {
	return func(t *ToolTip) {
		t.image = i
	}
}

func (o toolTipOpts) WithOffset(off img.Point) ToolTipOpt {
	return func(t *ToolTip) {
		t.offset = off
	}
}

func (o toolTipOpts) WithNoSticky() ToolTipOpt {
	return func(t *ToolTip) {
		t.sticky = false
	}
}

func (o toolTipOpts) WithDelay(d time.Duration) ToolTipOpt {
	return func(t *ToolTip) {
		t.delay = d
	}
}

func (t *ToolTip) GetWidget() *Widget {
	t.init.Do()
	return t.tipContainer.GetWidget()
}

func (t *ToolTip) Render(screen *ebiten.Image, def DeferredRenderFunc) {
	t.init.Do()

	x, y := input.CursorPosition()
	w := t.container.WidgetAt(x, y)

	defer func() {
		t.lastTipWidget = w
	}()

	if w != t.lastTipWidget ||
		input.MouseButtonPressed(ebiten.MouseButtonLeft) ||
		input.MouseButtonPressed(ebiten.MouseButtonMiddle) ||
		input.MouseButtonPressed(ebiten.MouseButtonRight) {

		t.doRender = false

		if t.timer != nil {
			if !t.timer.Stop() {
				<-t.timer.C
			}
			t.timer = nil
		}
	}

	if w == nil {
		return
	}

	if t.doRender {
		text := w.GetWidget().ToolTip
		if text == "" {
			return
		}

		if !t.sticky || t.doRenderReset || w != t.lastTipWidget || text != t.tipText.Label {
			defer func() {
				t.doRenderReset = false
			}()

			t.tipText.Label = text
			sx, sy := t.tipContainer.PreferredSize()
			r := img.Rect(x, y, x+sx, y+sy)
			r = r.Add(t.offset)
			t.tipContainer.SetLocation(r)
			t.tipContainer.RequestRelayout()
		}

		t.tipContainer.Render(screen, def)

		return
	}

	if t.timer == nil {
		t.timer = time.NewTimer(t.delay)
		go func() {
			<-t.timer.C
			t.doRenderReset = true
			t.doRender = true
			t.timer = nil
		}()
	}
}

func (t *ToolTip) createWidget() {
	t.tipContainer = NewContainer(
		ContainerOpts.WithLayout(NewFillLayout(
			FillLayoutOpts.WithPadding(t.padding),
		)),
		ContainerOpts.WithBackgroundImage(t.image))

	t.tipText = NewText(t.textOpts...)
	t.tipText.Label = ""
	t.tipContainer.AddChild(t.tipText)
	t.textOpts = nil
}
