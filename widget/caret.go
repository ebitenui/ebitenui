package widget

import (
	img "image"
	"image/color"
	"math"
	"sync/atomic"
	"time"

	"github.com/ebitenui/ebitenui/image"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Caret struct {
	Width int
	Color color.Color

	face          text.Face
	blinkInterval time.Duration

	init    *MultiOnce
	widget  *Widget
	image   *image.NineSlice
	height  int
	state   caretBlinkState
	visible bool
}

type CaretOpt func(c *Caret)

type CaretOptions struct {
}

var CaretOpts CaretOptions

type caretBlinkState func() caretBlinkState

func NewCaret(opts ...CaretOpt) *Caret {
	c := &Caret{
		blinkInterval: 450 * time.Millisecond,

		init: &MultiOnce{},
	}
	c.resetBlinking()

	c.init.Append(c.createWidget)

	for _, o := range opts {
		o(c)
	}

	c.validate()

	return c
}

func (c *Caret) validate() {
	if c.face == nil {
		panic("Caret: Font Face is required.")
	}
}

func (o CaretOptions) Color(c color.Color) CaretOpt {
	return func(ca *Caret) {
		ca.Color = c
	}
}

func (o CaretOptions) Size(face text.Face, width int) CaretOpt {
	return func(c *Caret) {
		c.face = face
		c.Width = width
	}
}

func (c *Caret) GetWidget() *Widget {
	c.init.Do()
	return c.widget
}

func (c *Caret) SetLocation(rect img.Rectangle) {
	c.init.Do()
	c.widget.Rect = rect
}

func (c *Caret) PreferredSize() (int, int) {
	c.init.Do()
	return c.Width, c.height
}

func (c *Caret) Render(screen *ebiten.Image) {
	c.init.Do()

	c.state = c.state()

	c.widget.Render(screen)

	if !c.visible {
		return
	}

	c.image = image.NewNineSliceColor(c.Color)

	c.image.Draw(screen, c.Width, c.height, func(opts *ebiten.DrawImageOptions) {
		p := c.widget.Rect.Min
		opts.GeoM.Translate(float64(p.X), float64(p.Y))
	})
}

func (c *Caret) Update() {
	c.init.Do()

	c.widget.Update()
}

func (c *Caret) ResetBlinking() {
	c.init.Do()
	c.resetBlinking()
}

func (c *Caret) resetBlinking() {
	c.state = c.blinkState(true, nil, nil)
}

func (c *Caret) blinkState(visible bool, timer *time.Timer, expired *atomic.Value) caretBlinkState {
	return func() caretBlinkState {
		c.visible = visible

		if timer != nil {
			if isExpired, _ := expired.Load().(bool); isExpired {
				return c.blinkState(!visible, nil, nil)
			}
		}

		if timer == nil {
			expired = &atomic.Value{}
			expired.Store(false)

			timer = time.AfterFunc(c.blinkInterval, func() {
				expired.Store(true)
			})
		}

		return c.blinkState(visible, timer, expired)
	}
}

func (c *Caret) createWidget() {
	c.widget = NewWidget()

	_, height := text.Measure(" ", c.face, 0)
	c.height = int(math.Round(height))
	c.face = nil
}
