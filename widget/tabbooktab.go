package widget

import (
	"github.com/ebitenui/ebitenui/image"
)

type TabBookTab struct {
	Container
	Disabled bool
	label    string
	image   *GraphicImage
}

type TabBookTabSelectedEventArgs struct {
	TabBook     *TabBook
	Tab         *TabBookTab
	PreviousTab *TabBookTab
}

type TabBookTabSelectedHandlerFunc func(args *TabBookTabSelectedEventArgs)

type TabParams struct {
	BackgroundImage *image.NineSlice
}

type TabBookTabOptions struct {
}

type TabBookTabOpt func(o *TabBookTab)

var TabBookTabOpts TabBookTabOptions

func (o *TabBookTabOptions) ContainerOpts(opts ...ContainerOpt) TabBookTabOpt {
	return func(t *TabBookTab) {
		for _, o := range opts {
			o(&t.Container)
		}
	}
}

func (o *TabBookTabOptions) Image(img *GraphicImage) TabBookTabOpt {
	return func(t *TabBookTab) {
		t.image = img
	}
}

func NewTabBookTab(label string, opts... TabBookTabOpt) *TabBookTab {
	c := &TabBookTab{
		label: label,
	}
	c.init = &MultiOnce{}
	c.init.Append(c.createWidget)

	// Set a default layout so that tabs use the full container
	c.widgetOpts = append(c.widgetOpts, WidgetOpts.LayoutData(AnchorLayoutData{
		StretchHorizontal:  true,
		StretchVertical:    true,
		HorizontalPosition: AnchorLayoutPositionCenter,
		VerticalPosition:   AnchorLayoutPositionCenter,
	}))

	for _, o := range opts {
		o(c)
	}
	return c
}

func (t *TabBookTab) Validate() {
	t.Container.Validate()

	if t.definedParams.BackgroundImage != nil {
		t.computedParams.BackgroundImage = t.definedParams.BackgroundImage
	} else {
		theme := t.widget.GetTheme()
		if theme != nil && theme.TabTheme != nil && theme.TabTheme.BackgroundImage != nil {
			t.computedParams.BackgroundImage = theme.TabTheme.BackgroundImage
		}
	}
}
