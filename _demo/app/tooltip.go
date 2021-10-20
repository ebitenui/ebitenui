package app

import (
	"time"

	"github.com/blizzy78/ebitenui/widget"
)

type toolTipContents struct {
	tips            map[widget.HasWidget]string
	widgetsWithTime []widget.HasWidget
	showTime        bool

	res *uiResources

	text     *widget.TextToolTip
	timeText *widget.TextToolTip
}

func (t *toolTipContents) Create(w widget.HasWidget) widget.ToolTipWidget {
	if _, ok := t.tips[w]; !ok {
		return nil
	}

	c := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(t.res.toolTip.background),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(t.res.toolTip.padding),
			widget.RowLayoutOpts.Spacing(2),
		)))

	t.text = widget.NewTextToolTip(
		widget.TextToolTipOpts.TextOpts(
			widget.TextOpts.Text("", t.res.toolTip.face, t.res.toolTip.color),
		),
	)
	c.AddChild(t.text)

	if t.showTime && t.canShowTime(w) {
		t.timeText = widget.NewTextToolTip(
			widget.TextToolTipOpts.TextOpts(
				widget.TextOpts.Text("", t.res.toolTip.face, t.res.toolTip.color),
			),
		)
		c.AddChild(t.timeText)
	}

	return c
}

func (t *toolTipContents) Set(w widget.HasWidget, s string) {
	t.tips[w] = s
}

func (t *toolTipContents) Update(w widget.HasWidget) {
	t.text.Label = t.tips[w]

	if !t.showTime || !t.canShowTime(w) {
		return
	}

	t.timeText.Label = time.Now().Local().Format("2006-01-02 15:04:05")
}

func (t *toolTipContents) canShowTime(w widget.HasWidget) bool {
	for _, tw := range t.widgetsWithTime {
		if tw == w {
			return true
		}
	}
	return false
}
