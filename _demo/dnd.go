package main

import (
	"github.com/ebitenui/ebitenui/widget"
)

type dragContents struct {
	res *uiResources

	sources []*widget.Widget
	targets []*widget.Widget

	text *widget.Text
}

func (d *dragContents) Create(srcWidget widget.HasWidget, srcX int, srcY int) (widget.DragWidget, interface{}) {
	if !d.isSource(srcWidget.GetWidget()) {
		return nil, nil
	}

	c := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(d.res.toolTip.background),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			widget.AnchorLayoutOpts.Padding(d.res.toolTip.padding),
		)),
	)

	d.text = widget.NewText(widget.TextOpts.Text("Drag Me!", d.res.toolTip.face, d.res.toolTip.color))
	c.AddChild(d.text)

	return c, nil
}

func (d *dragContents) Update(target widget.HasWidget, _ int, _ int, _ interface{}) {
	if target != nil && d.isTarget(target.GetWidget()) {
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
