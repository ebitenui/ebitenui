package main

import (
	"time"

	"github.com/blizzy78/ebitenui/widget"
)

type toolTipContents struct {
	tips            map[widget.HasWidget]string
	widgetsWithTime []widget.HasWidget
	showTime        bool

	res *resources

	text     *widget.TextToolTip
	timeText *widget.TextToolTip
}

func (t *toolTipContents) Create(w widget.HasWidget) widget.HasWidget {
	if _, ok := t.tips[w]; !ok {
		return nil
	}

	c := widget.NewContainer(
		widget.ContainerOpts.WithLayout(widget.NewRowLayout(
			widget.RowLayoutOpts.WithDirection(widget.DirectionVertical),
			widget.RowLayoutOpts.WithSpacing(2),
		)))

	t.text = widget.NewTextToolTip(
		widget.TextToolTipOpts.WithUpdater(t),
		widget.TextToolTipOpts.WithContainerOpts(
			widget.ContainerOpts.WithBackgroundImage(t.res.images.button.Disabled),
			widget.ContainerOpts.WithWidgetOpts(
				widget.WidgetOpts.WithLayoutData(&widget.RowLayoutData{
					Stretch: true,
				}),
			),
		),
		widget.TextToolTipOpts.WithPadding(widget.Insets{
			Left:   8,
			Right:  8,
			Top:    4,
			Bottom: 4,
		}),
		widget.TextToolTipOpts.WithTextOpts(
			widget.TextOpts.WithText("", t.res.fonts.toolTipFace, t.res.colors.textToolTip)))
	c.AddChild(t.text)

	canShowTime := false
	for _, tw := range t.widgetsWithTime {
		if tw == w {
			canShowTime = true
			break
		}
	}

	if t.showTime && canShowTime {
		t.timeText = widget.NewTextToolTip(
			widget.TextToolTipOpts.WithContainerOpts(
				widget.ContainerOpts.WithBackgroundImage(t.res.images.button.Disabled),
				widget.ContainerOpts.WithWidgetOpts(
					widget.WidgetOpts.WithLayoutData(&widget.RowLayoutData{
						Stretch: true,
					}),
				),
			),
			widget.TextToolTipOpts.WithPadding(widget.Insets{
				Left:   8,
				Right:  8,
				Top:    4,
				Bottom: 4,
			}),
			widget.TextToolTipOpts.WithTextOpts(
				widget.TextOpts.WithText("", t.res.fonts.toolTipFace, t.res.colors.textToolTip)))
		c.AddChild(t.timeText)
	}

	return c
}

func (t *toolTipContents) Set(w widget.HasWidget, s string) {
	t.tips[w] = s
}

func (t *toolTipContents) Update(w widget.HasWidget) {
	t.text.Label = t.tips[w]

	canShowTime := false
	for _, tw := range t.widgetsWithTime {
		if tw == w {
			canShowTime = true
			break
		}
	}

	if t.showTime && canShowTime {
		t.timeText.Label = time.Now().Local().Format("2006-01-02 15:04:05")
	}
}
