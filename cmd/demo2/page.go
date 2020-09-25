package main

import (
	"fmt"
	"time"

	"github.com/blizzy78/ebitenui/widget"
)

type page struct {
	title   string
	content widget.HasWidget
}

func buttonPage(res *resources) *page {
	c := newPageContentContainer()

	b1 := widget.NewButton(
		widget.ButtonOpts.WithWidgetOpts(widget.WidgetOpts.WithLayoutData(&widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ButtonOpts.WithImage(res.images.button),
		widget.ButtonOpts.WithText("Button", res.fonts.face, res.colors.buttonText))
	c.AddChild(b1)

	b2 := widget.NewButton(
		widget.ButtonOpts.WithWidgetOpts(widget.WidgetOpts.WithLayoutData(&widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ButtonOpts.WithImage(res.images.button),
		widget.ButtonOpts.WithTextAndImage("Button with Graphic", res.fonts.face, res.images.heart, res.colors.buttonText))
	c.AddChild(b2)

	b3 := widget.NewButton(
		widget.ButtonOpts.WithWidgetOpts(widget.WidgetOpts.WithLayoutData(&widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ButtonOpts.WithImage(res.images.button),
		widget.ButtonOpts.WithText("Multi\nLine\nButton", res.fonts.face, res.colors.buttonText))
	c.AddChild(b3)

	c.AddChild(widget.NewButton(
		widget.ButtonOpts.WithWidgetOpts(widget.WidgetOpts.WithLayoutData(&widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ButtonOpts.WithImage(res.images.buttonKenney),
		widget.ButtonOpts.WithText("Button", res.fonts.face, res.colors.buttonText),
		widget.ButtonOpts.WithTextPadding(widget.Insets{
			Left:   10,
			Right:  10,
			Top:    6,
			Bottom: 10,
		})))

	c.AddChild(newSeparator(res, &widget.RowLayoutData{
		Stretch: true,
	}))

	c.AddChild(newCheckbox("Disabled", func(args *widget.CheckboxChangedEventArgs) {
		b1.GetWidget().Disabled = args.State == widget.CheckboxChecked
		b2.GetWidget().Disabled = args.State == widget.CheckboxChecked
		b3.GetWidget().Disabled = args.State == widget.CheckboxChecked
	}, res))

	return &page{
		title:   "Button",
		content: c,
	}
}

func checkboxPage(res *resources) *page {
	c := newPageContentContainer()

	cb1 := widget.NewLabeledCheckbox(
		widget.LabeledCheckboxOpts.WithCheckboxOpts(
			widget.CheckboxOpts.WithButtonOpts(widget.ButtonOpts.WithImage(res.images.button)),
			widget.CheckboxOpts.WithImage(res.images.checkbox)),
		widget.LabeledCheckboxOpts.WithLabelOpts(widget.LabelOpts.WithText("Two-State Checkbox", res.fonts.face, res.colors.label)))
	c.AddChild(cb1)

	cb2 := widget.NewLabeledCheckbox(
		widget.LabeledCheckboxOpts.WithCheckboxOpts(
			widget.CheckboxOpts.WithButtonOpts(widget.ButtonOpts.WithImage(res.images.button)),
			widget.CheckboxOpts.WithImage(res.images.checkbox),
			widget.CheckboxOpts.WithTriState()),
		widget.LabeledCheckboxOpts.WithLabelOpts(widget.LabelOpts.WithText("Tri-State Checkbox", res.fonts.face, res.colors.label)))
	c.AddChild(cb2)

	c.AddChild(newSeparator(res, &widget.RowLayoutData{
		Stretch: true,
	}))

	c.AddChild(newCheckbox("Disabled", func(args *widget.CheckboxChangedEventArgs) {
		cb1.GetWidget().Disabled = args.State == widget.CheckboxChecked
		cb2.GetWidget().Disabled = args.State == widget.CheckboxChecked
	}, res))

	return &page{
		title:   "Checkbox",
		content: c,
	}
}

func listPage(res *resources) *page {
	c := newPageContentContainer()

	listsContainer := widget.NewContainer(
		widget.ContainerOpts.WithWidgetOpts(widget.WidgetOpts.WithLayoutData(&widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ContainerOpts.WithLayout(widget.NewGridLayout(
			widget.GridLayoutOpts.WithColumns(3),
			widget.GridLayoutOpts.WithStretch([]bool{true, false, true}, []bool{true}),
			widget.GridLayoutOpts.WithSpacing(10, 0))))
	c.AddChild(listsContainer)

	entries1 := []interface{}{"One", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten"}
	list1 := newList(entries1, res, widget.WidgetOpts.WithLayoutData(&widget.GridLayoutData{
		MaxHeight: 220,
	}))
	listsContainer.AddChild(list1)

	buttonsContainer := widget.NewContainer(
		widget.ContainerOpts.WithLayout(widget.NewRowLayout(
			widget.RowLayoutOpts.WithDirection(widget.DirectionVertical),
			widget.RowLayoutOpts.WithSpacing(10),
		)))
	listsContainer.AddChild(buttonsContainer)

	for i := 0; i < 3; i++ {
		buttonsContainer.AddChild(widget.NewButton(
			widget.ButtonOpts.WithWidgetOpts(widget.WidgetOpts.WithLayoutData(&widget.RowLayoutData{
				Stretch: true,
			})),
			widget.ButtonOpts.WithImage(res.images.button),
			widget.ButtonOpts.WithText(fmt.Sprintf("Action %d", i+1), res.fonts.face, res.colors.buttonText)))
	}

	entries2 := []interface{}{"Eleven", "Twelve", "Thirteen", "Fourteen", "Fifteen", "Sixteen", "Seventeen", "Eighteen", "Nineteen", "Twenty"}
	list2 := newList(entries2, res, widget.WidgetOpts.WithLayoutData(&widget.GridLayoutData{
		MaxHeight: 220,
	}))
	listsContainer.AddChild(list2)

	c.AddChild(newSeparator(res, &widget.RowLayoutData{
		Stretch: true,
	}))

	c.AddChild(newCheckbox("Disabled", func(args *widget.CheckboxChangedEventArgs) {
		list1.GetWidget().Disabled = args.State == widget.CheckboxChecked
		list2.GetWidget().Disabled = args.State == widget.CheckboxChecked
	}, res))

	return &page{
		title:   "List",
		content: c,
	}
}

func comboButtonPage(res *resources) *page {
	c := newPageContentContainer()

	entries := []interface{}{}
	for i := 1; i <= 20; i++ {
		entries = append(entries, i)
	}

	cb := newListComboButton(
		entries,
		func(e interface{}) string {
			return fmt.Sprintf("Entry %d", e.(int))
		},
		func(e interface{}) string {
			return fmt.Sprintf("Entry %d", e.(int))
		},
		func(args *widget.ListComboButtonEntrySelectedEventArgs) {
			c.RequestRelayout()
		},
		res)
	c.AddChild(cb)

	c.AddChild(newSeparator(res, &widget.RowLayoutData{
		Stretch: true,
	}))

	c.AddChild(newCheckbox("Disabled", func(args *widget.CheckboxChangedEventArgs) {
		cb.GetWidget().Disabled = args.State == widget.CheckboxChecked
	}, res))

	return &page{
		title:   "Combo Button",
		content: c,
	}
}

func tabBookPage(res *resources) *page {
	c := newPageContentContainer()

	tabs := []*widget.TabBookTab{}

	for i := 0; i < 5; i++ {
		tc := widget.NewContainer(
			widget.ContainerOpts.WithLayout(widget.NewRowLayout(
				widget.RowLayoutOpts.WithDirection(widget.DirectionVertical),
				widget.RowLayoutOpts.WithSpacing(10))),
			widget.ContainerOpts.WithAutoDisableChildren())

		for j := 0; j < 3; j++ {
			b := widget.NewButton(
				widget.ButtonOpts.WithImage(res.images.button),
				widget.ButtonOpts.WithText(fmt.Sprintf("Button %d on Tab %d", j+1, i+1), res.fonts.face, res.colors.buttonText))
			tc.AddChild(b)
		}

		tab := widget.NewTabBookTab(fmt.Sprintf("Tab %d", i+1), tc)
		if i == 2 {
			tab.Disabled = true
		}

		tabs = append(tabs, tab)
	}

	t := widget.NewTabBook(
		widget.TabBookOpts.WithTabs(tabs...),
		widget.TabBookOpts.WithTabButtonImage(res.images.button, res.images.stateButtonSelected),
		widget.TabBookOpts.WithTabButtonText(res.fonts.face, res.colors.buttonText),
		widget.TabBookOpts.WithTabButtonSpacing(4),
		widget.TabBookOpts.WithSpacing(10))
	c.AddChild(t)

	c.AddChild(newSeparator(res, &widget.RowLayoutData{
		Stretch: true,
	}))

	c.AddChild(newCheckbox("Disabled", func(args *widget.CheckboxChangedEventArgs) {
		t.GetWidget().Disabled = args.State == widget.CheckboxChecked
	}, res))

	return &page{
		title:   "Tab Book",
		content: c,
	}
}

func gridLayoutPage(res *resources) *page {
	c := widget.NewContainer(
		widget.ContainerOpts.WithLayout(widget.NewGridLayout(
			widget.GridLayoutOpts.WithColumns(5),
			widget.GridLayoutOpts.WithStretch([]bool{true, true, true, true, true}, nil),
			widget.GridLayoutOpts.WithSpacing(5, 5))))

	for row := 0; row < 3; row++ {
		for col := 0; col < 5; col++ {
			i := row*5 + col
			b := widget.NewButton(
				widget.ButtonOpts.WithImage(res.images.button),
				widget.ButtonOpts.WithText(fmt.Sprintf("Button %d", i+1), res.fonts.face, res.colors.buttonText))
			c.AddChild(b)
		}
	}

	return &page{
		title:   "Grid Layout",
		content: c,
	}
}

func rowLayoutPage(res *resources) *page {
	c := newPageContentContainer()

	c.AddChild(widget.NewText(
		widget.TextOpts.WithText("Horizontal", res.fonts.face, res.colors.textIdle)))

	bc := widget.NewContainer(
		widget.ContainerOpts.WithLayout(widget.NewRowLayout(
			widget.RowLayoutOpts.WithSpacing(5))))
	c.AddChild(bc)

	for col := 0; col < 5; col++ {
		b := widget.NewButton(
			widget.ButtonOpts.WithImage(res.images.button),
			widget.ButtonOpts.WithText(fmt.Sprintf("Button %d", col+1), res.fonts.face, res.colors.buttonText))
		bc.AddChild(b)
	}

	c.AddChild(widget.NewText(
		widget.TextOpts.WithText("Vertical", res.fonts.face, res.colors.textIdle)))

	bc = widget.NewContainer(
		widget.ContainerOpts.WithLayout(widget.NewRowLayout(
			widget.RowLayoutOpts.WithDirection(widget.DirectionVertical),
			widget.RowLayoutOpts.WithSpacing(5))))
	c.AddChild(bc)

	labels := []string{"Tiny", "Medium", "Very Large"}
	for _, l := range labels {
		b := widget.NewButton(
			widget.ButtonOpts.WithWidgetOpts(widget.WidgetOpts.WithLayoutData(&widget.RowLayoutData{
				Stretch: true,
			})),
			widget.ButtonOpts.WithImage(res.images.button),
			widget.ButtonOpts.WithText(l, res.fonts.face, res.colors.buttonText))
		bc.AddChild(b)
	}

	return &page{
		title:   "Row Layout",
		content: c,
	}
}

func sliderPage(res *resources) *page {
	c := newPageContentContainer()

	pageSizes := []int{3, 10}
	sliders := []*widget.Slider{}

	for _, ps := range pageSizes {
		ps := ps

		sc := widget.NewContainer(
			widget.ContainerOpts.WithLayout(widget.NewRowLayout(
				widget.RowLayoutOpts.WithSpacing(10))))
		c.AddChild(sc)

		var text *widget.Text

		s := widget.NewSlider(
			widget.SliderOpts.WithWidgetOpts(widget.WidgetOpts.WithLayoutData(&widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			})),
			widget.SliderOpts.WithMinMax(1, 20),
			widget.SliderOpts.WithImages(res.images.sliderTrack, res.images.button),
			widget.SliderOpts.WithTrackPadding(3),
			widget.SliderOpts.WithPageSizeFunc(func() int {
				return ps
			}),
			widget.SliderOpts.WithChangedHandler(func(args *widget.SliderChangedEventArgs) {
				text.Label = fmt.Sprintf("%d", args.Current)
			}),
		)
		sc.AddChild(s)
		sliders = append(sliders, s)

		text = widget.NewText(
			widget.TextOpts.WithWidgetOpts(widget.WidgetOpts.WithLayoutData(&widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			})),
			widget.TextOpts.WithText(fmt.Sprintf("%d", s.Current), res.fonts.face, res.colors.textIdle))
		sc.AddChild(text)
	}

	c.AddChild(newSeparator(res, &widget.RowLayoutData{
		Stretch: true,
	}))

	c.AddChild(newCheckbox("Disabled", func(args *widget.CheckboxChangedEventArgs) {
		for _, s := range sliders {
			s.GetWidget().Disabled = args.State == widget.CheckboxChecked
		}
	}, res))

	return &page{
		title:   "Slider",
		content: c,
	}
}

func toolTipPage(res *resources, toolTips *toolTipContents) *page {
	c := newPageContentContainer()

	c.AddChild(widget.NewText(
		widget.TextOpts.WithText("Hover over these buttons to see their tool tips.", res.fonts.face, res.colors.textIdle)))

	bc := widget.NewContainer(
		widget.ContainerOpts.WithLayout(widget.NewRowLayout(
			widget.RowLayoutOpts.WithSpacing(5))))
	c.AddChild(bc)

	for col := 0; col < 5; col++ {
		b := widget.NewButton(
			widget.ButtonOpts.WithImage(res.images.button),
			widget.ButtonOpts.WithText(fmt.Sprintf("Button %d", col+1), res.fonts.face, res.colors.buttonText))

		if col == 2 {
			b.GetWidget().Disabled = true
		}

		toolTips.Set(b, fmt.Sprintf("Tool tip for button %d", col+1))
		toolTips.widgetsWithTime = append(toolTips.widgetsWithTime, b)

		bc.AddChild(b)
	}

	c.AddChild(newSeparator(res, &widget.RowLayoutData{
		Stretch: true,
	}))

	showTimeCheckbox := newCheckbox("Show additional infos in tool tips", func(args *widget.CheckboxChangedEventArgs) {
		toolTips.showTime = args.State == widget.CheckboxChecked
	}, res)
	toolTips.Set(showTimeCheckbox, "If enabled, tool tips will show system time for demonstration.")
	c.AddChild(showTimeCheckbox)

	return &page{
		title:   "Tool Tip",
		content: c,
	}
}

func dragAndDropPage(res *resources, dnd *widget.DragAndDrop, drag *dragContents) *page {
	c := newPageContentContainer()

	dndContainer := widget.NewContainer(
		widget.ContainerOpts.WithLayout(widget.NewRowLayout(
			widget.RowLayoutOpts.WithSpacing(20),
		)),
	)
	c.AddChild(dndContainer)

	sourceContainer := widget.NewContainer(
		widget.ContainerOpts.WithWidgetOpts(widget.WidgetOpts.WithLayoutData(&widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ContainerOpts.WithBackgroundImage(res.images.button.Idle),
		widget.ContainerOpts.WithLayout(widget.NewFillLayout(
			widget.FillLayoutOpts.WithPadding(widget.NewInsetsSimple(20)),
		)),
	)
	drag.addSource(sourceContainer)
	dndContainer.AddChild(sourceContainer)

	sourceContainer.AddChild(widget.NewText(
		widget.TextOpts.WithText("Drag\nFrom\nHere", res.fonts.face, res.colors.textIdle),
		widget.TextOpts.WithPosition(widget.TextPositionCenter),
	))

	targetContainer := widget.NewContainer(
		widget.ContainerOpts.WithWidgetOpts(widget.WidgetOpts.WithLayoutData(&widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ContainerOpts.WithBackgroundImage(res.images.button.Idle),
		widget.ContainerOpts.WithLayout(widget.NewFillLayout(
			widget.FillLayoutOpts.WithPadding(widget.NewInsetsSimple(20)),
		)),
	)
	drag.addTarget(targetContainer)
	dndContainer.AddChild(targetContainer)

	targetContainer.AddChild(widget.NewText(
		widget.TextOpts.WithText("Drop\nHere", res.fonts.face, res.colors.textIdle),
		widget.TextOpts.WithPosition(widget.TextPositionCenter),
	))

	dnd.DroppedEvent.AddHandler(func(args interface{}) {
		a := args.(*widget.DragAndDropDroppedEventArgs)
		if !drag.isTarget(a.Target.GetWidget()) {
			return
		}

		targetContainer.BackgroundImage = res.images.button.Pressed

		time.AfterFunc(2*time.Second, func() {
			targetContainer.BackgroundImage = res.images.button.Idle
		})
	})

	return &page{
		title:   "Drag & Drop",
		content: c,
	}
}
