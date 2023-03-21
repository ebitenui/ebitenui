package main

import (
	"github.com/ebitenui/ebitenui/widget"
)

type dragContents struct {
	res *uiResources

	text *widget.Text
}

func (d *dragContents) Create(sourceWidget widget.HasWidget) (*widget.Container, interface{}) {

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

func (d *dragContents) Update(isDroppable bool, _ widget.HasWidget, _ interface{}) {
	if isDroppable {
		d.text.Label = "* DROP ME! *"
	} else {
		d.text.Label = "Drag Me!"
	}
}
