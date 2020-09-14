package widget

import (
	img "image"
	"time"

	"github.com/blizzy78/ebitenui/input"
	"github.com/hajimehoshi/ebiten"
)

type ToolTip struct {
	container        WidgetLocator
	contentsCreater  ToolTipContentsCreater
	offset           img.Point
	sticky           bool
	delay            time.Duration
	updateEveryFrame bool

	init             *MultiOnce
	tipWidget        HasWidget
	lastTippedWidget HasWidget
	timer            *time.Timer
	doRender         bool
	doRelayout       bool
}

type ToolTipOpt func(t *ToolTip)

const ToolTipOpts = toolTipOpts(true)

type toolTipOpts bool

type ToolTipContentsCreater interface {
	Create(HasWidget) HasWidget
}

type ToolTipContentsUpdater interface {
	Update(HasWidget)
}

func NewToolTip(opts ...ToolTipOpt) *ToolTip {
	t := &ToolTip{
		offset: img.Point{0, 20},
		sticky: true,
		delay:  800 * time.Millisecond,

		init: &MultiOnce{},
	}

	for _, o := range opts {
		o(t)
	}

	return t
}

func (o toolTipOpts) WithContainer(c *Container) ToolTipOpt {
	return func(t *ToolTip) {
		t.container = c
	}
}

func (o toolTipOpts) WithContentsCreater(c ToolTipContentsCreater) ToolTipOpt {
	return func(t *ToolTip) {
		t.contentsCreater = c
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

func (o toolTipOpts) WithUpdateEveryFrame() ToolTipOpt {
	return func(t *ToolTip) {
		t.updateEveryFrame = true
	}
}

func (t *ToolTip) Render(screen *ebiten.Image, def DeferredRenderFunc) {
	t.init.Do()

	x, y := input.CursorPosition()
	w := t.container.WidgetAt(x, y)

	defer func() {
		t.lastTippedWidget = w
	}()

	if w != t.lastTippedWidget ||
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

		t.tipWidget = nil
	}

	if w == nil {
		return
	}

	if t.doRender {
		justCreated := false

		if t.tipWidget == nil {
			t.tipWidget = t.contentsCreater.Create(w)

			if t.tipWidget == nil {
				return
			}

			justCreated = true
		}

		if justCreated || t.updateEveryFrame || w != t.lastTippedWidget {
			if u, ok := t.contentsCreater.(ToolTipContentsUpdater); ok {
				u.Update(w)
			}
		}

		if !t.sticky || t.doRelayout || w != t.lastTippedWidget {
			defer func() {
				t.doRelayout = false
			}()

			sx, sy := t.tipWidget.(PreferredSizer).PreferredSize()
			r := img.Rect(x, y, x+sx, y+sy)
			r = r.Add(t.offset)
			t.tipWidget.(Locateable).SetLocation(r)
			t.tipWidget.(Relayoutable).RequestRelayout()
		}

		t.tipWidget.(Renderer).Render(screen, def)

		return
	}

	if t.timer == nil {
		if t.delay > 0 {
			t.timer = time.AfterFunc(t.delay, func() {
				t.doRelayout = true
				t.doRender = true
				t.timer = nil
			})
		} else {
			t.doRelayout = true
			t.doRender = true
		}
	}
}
