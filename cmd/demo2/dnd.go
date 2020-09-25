package main

import (
	"github.com/blizzy78/ebitenui/widget"
)

type dragContents struct {
	res *resources

	sources []*widget.Widget
	targets []*widget.Widget

	text *widget.Text
}

func newTextDragContents(res *resources) *dragContents {
	return &dragContents{
		res: res,
	}
}

func (d *dragContents) Create(srcWidget widget.HasWidget, srcX int, srcY int) (widget.HasWidget, interface{}) {
	if !d.isSource(srcWidget.GetWidget()) {
		return nil, nil
	}

	c := widget.NewContainer(
		widget.ContainerOpts.WithBackgroundImage(d.res.images.button.Disabled),
		widget.ContainerOpts.WithWidgetOpts(widget.WidgetOpts.WithLayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ContainerOpts.WithLayout(widget.NewFillLayout(
			widget.FillLayoutOpts.WithPadding(widget.Insets{
				Left:   8,
				Right:  8,
				Top:    4,
				Bottom: 4,
			}),
		)),
	)

	d.text = widget.NewText(widget.TextOpts.WithText("Drag Me!", d.res.fonts.face, d.res.colors.textIdle))
	c.AddChild(d.text)

	return c, nil
}

func (d *dragContents) Update(target widget.HasWidget, x int, y int, dragData interface{}) {
	if d.isTarget(target.GetWidget()) {
		d.text.Label = "* DROP ME! *"
	} else {
		d.text.Label = "Drag Me!"
	}
}

func (d *dragContents) addSource(s widget.HasWidget) {
	d.sources = append(d.sources, s.GetWidget())
}

func (d *dragContents) addTarget(t widget.HasWidget) {
	d.targets = append(d.targets, t.GetWidget())
}

func (d *dragContents) isSource(w *widget.Widget) bool {
	for _, s := range d.sources {
		if s == w {
			return true
		}
	}

	p := w.Parent()
	if p == nil {
		return false
	}

	return d.isSource(p)
}

func (d *dragContents) isTarget(w *widget.Widget) bool {
	for _, t := range d.targets {
		if t == w {
			return true
		}
	}

	p := w.Parent()
	if p == nil {
		return false
	}

	return d.isTarget(p)
}
