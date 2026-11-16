package widget

import (
	"github.com/ebitenui/ebitenui/image"
)

type PanelParams struct {
	BackgroundImage *image.NineSlice
}

type Panel struct {
	Container
}

func NewPanel(opts ...ContainerOpt) *Panel {
	c := &Panel{}
	c.init = &MultiOnce{}

	c.init.Append(c.createWidget)

	for _, o := range opts {
		o(&c.Container)
	}

	return c
}

func (p *Panel) Validate() {
	p.Container.Validate()

	if p.definedParams.BackgroundImage != nil {
		p.computedParams.BackgroundImage = p.definedParams.BackgroundImage
	} else {
		theme := p.widget.GetTheme()
		if theme != nil && theme.PanelTheme != nil && theme.PanelTheme.BackgroundImage != nil {
			p.computedParams.BackgroundImage = theme.PanelTheme.BackgroundImage
		}
	}

	if p.computedParams.BackgroundImage == nil {
		panic("Panel: BackgroundImage is required.")
	}

}

// themedPanel is a container that, without a background of its own, draws
// its theme's panel background and follows the theme as it changes, as a
// Panel does, but draws none without a theme where a Panel panics. The
// on-screen keyboard and the dialog are drawn on one.
type themedPanel struct{ *Container }

func (c *themedPanel) Validate() {
	c.Container.Validate()
	if c.definedParams.BackgroundImage == nil {
		if theme := c.widget.GetTheme(); theme != nil && theme.PanelTheme != nil {
			c.computedParams.BackgroundImage = theme.PanelTheme.BackgroundImage
		}
	}
}
