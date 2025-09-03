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
