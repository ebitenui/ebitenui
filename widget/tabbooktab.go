package widget

import "github.com/ebitenui/ebitenui/image"

type TabBookTab struct {
	Container
	Disabled bool
	label    string
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

func NewTabBookTab(label string, opts ...ContainerOpt) *TabBookTab {
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
		o(&c.Container)
	}
	return c
}

func (t *TabBookTab) Validate() {
	t.Container.Validate()

	if t.definedParams.BackgroundImage != nil {
		t.computedParams.BackgroundImage = t.definedParams.BackgroundImage
	} else {
		theme := t.widget.GetTheme()
		if theme != nil && theme.PanelTheme != nil && theme.PanelTheme.BackgroundImage != nil {
			t.computedParams.BackgroundImage = theme.PanelTheme.BackgroundImage
		}
	}
}
