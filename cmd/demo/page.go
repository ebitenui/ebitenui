package main

import (
	"fmt"
	"image"
	"time"

	"github.com/blizzy78/ebitenui"
	"github.com/blizzy78/ebitenui/widget"
	"github.com/hajimehoshi/ebiten"
)

type page struct {
	title   string
	content widget.PreferredSizeLocateableWidget
}

func buttonPage(res *resources) *page {
	c := newPageContentContainer()

	b1 := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(&widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ButtonOpts.Image(res.images.button),
		widget.ButtonOpts.Text("Button", res.fonts.face, res.colors.buttonText))
	c.AddChild(b1)

	b2 := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(&widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ButtonOpts.Image(res.images.button),
		widget.ButtonOpts.TextAndImage("Button with Graphic", res.fonts.face, res.images.heart, res.colors.buttonText))
	c.AddChild(b2)

	b3 := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(&widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ButtonOpts.Image(res.images.button),
		widget.ButtonOpts.Text("Multi\nLine\nButton", res.fonts.face, res.colors.buttonText))
	c.AddChild(b3)

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
		widget.LabeledCheckboxOpts.Spacing(6),
		widget.LabeledCheckboxOpts.CheckboxOpts(
			widget.CheckboxOpts.ButtonOpts(
				widget.ButtonOpts.Image(res.images.button),
				widget.ButtonOpts.GraphicPadding(widget.NewInsetsSimple(7)),
			),
			widget.CheckboxOpts.Image(res.images.checkbox)),
		widget.LabeledCheckboxOpts.LabelOpts(widget.LabelOpts.Text("Two-State Checkbox", res.fonts.face, res.colors.label)))
	c.AddChild(cb1)

	cb2 := widget.NewLabeledCheckbox(
		widget.LabeledCheckboxOpts.Spacing(6),
		widget.LabeledCheckboxOpts.CheckboxOpts(
			widget.CheckboxOpts.ButtonOpts(
				widget.ButtonOpts.Image(res.images.button),
				widget.ButtonOpts.GraphicPadding(widget.NewInsetsSimple(7)),
			),
			widget.CheckboxOpts.Image(res.images.checkbox),
			widget.CheckboxOpts.TriState()),
		widget.LabeledCheckboxOpts.LabelOpts(widget.LabelOpts.Text("Tri-State Checkbox", res.fonts.face, res.colors.label)))
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
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(&widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(3),
			widget.GridLayoutOpts.Stretch([]bool{true, false, true}, []bool{true}),
			widget.GridLayoutOpts.Spacing(10, 0))))
	c.AddChild(listsContainer)

	entries1 := []interface{}{"One", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten"}
	list1 := newList(entries1, res, widget.WidgetOpts.LayoutData(&widget.GridLayoutData{
		MaxHeight: 220,
	}))
	listsContainer.AddChild(list1)

	buttonsContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(10),
		)))
	listsContainer.AddChild(buttonsContainer)

	for i := 0; i < 3; i++ {
		buttonsContainer.AddChild(widget.NewButton(
			widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(&widget.RowLayoutData{
				Stretch: true,
			})),
			widget.ButtonOpts.Image(res.images.button),
			widget.ButtonOpts.Text(fmt.Sprintf("Action %d", i+1), res.fonts.face, res.colors.buttonText)))
	}

	entries2 := []interface{}{"Eleven", "Twelve", "Thirteen", "Fourteen", "Fifteen", "Sixteen", "Seventeen", "Eighteen", "Nineteen", "Twenty"}
	list2 := newList(entries2, res, widget.WidgetOpts.LayoutData(&widget.GridLayoutData{
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
			widget.ContainerOpts.Layout(widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(10))),
			widget.ContainerOpts.AutoDisableChildren())

		for j := 0; j < 3; j++ {
			b := widget.NewButton(
				widget.ButtonOpts.Image(res.images.button),
				widget.ButtonOpts.Text(fmt.Sprintf("Button %d on Tab %d", j+1, i+1), res.fonts.face, res.colors.buttonText))
			tc.AddChild(b)
		}

		tab := widget.NewTabBookTab(fmt.Sprintf("Tab %d", i+1), tc)
		if i == 2 {
			tab.Disabled = true
		}

		tabs = append(tabs, tab)
	}

	t := widget.NewTabBook(
		widget.TabBookOpts.Tabs(tabs...),
		widget.TabBookOpts.TabButtonImage(res.images.button, res.images.stateButtonSelected),
		widget.TabBookOpts.TabButtonText(res.fonts.face, res.colors.buttonText),
		widget.TabBookOpts.TabButtonSpacing(4),
		widget.TabBookOpts.Spacing(10))
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
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(5),
			widget.GridLayoutOpts.Stretch([]bool{true, true, true, true, true}, nil),
			widget.GridLayoutOpts.Spacing(5, 5))))

	for row := 0; row < 3; row++ {
		for col := 0; col < 5; col++ {
			i := row*5 + col
			b := widget.NewButton(
				widget.ButtonOpts.Image(res.images.button),
				widget.ButtonOpts.Text(fmt.Sprintf("Button %d", i+1), res.fonts.face, res.colors.buttonText))
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
		widget.TextOpts.Text("Horizontal", res.fonts.face, res.colors.textIdle)))

	bc := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(5))))
	c.AddChild(bc)

	for col := 0; col < 5; col++ {
		b := widget.NewButton(
			widget.ButtonOpts.Image(res.images.button),
			widget.ButtonOpts.Text(fmt.Sprintf("Button %d", col+1), res.fonts.face, res.colors.buttonText))
		bc.AddChild(b)
	}

	c.AddChild(widget.NewText(
		widget.TextOpts.Text("Vertical", res.fonts.face, res.colors.textIdle)))

	bc = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(5))))
	c.AddChild(bc)

	labels := []string{"Tiny", "Medium", "Very Large"}
	for _, l := range labels {
		b := widget.NewButton(
			widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(&widget.RowLayoutData{
				Stretch: true,
			})),
			widget.ButtonOpts.Image(res.images.button),
			widget.ButtonOpts.Text(l, res.fonts.face, res.colors.buttonText))
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
			widget.ContainerOpts.Layout(widget.NewRowLayout(
				widget.RowLayoutOpts.Spacing(10))))
		c.AddChild(sc)

		var text *widget.Text

		s := widget.NewSlider(
			widget.SliderOpts.WidgetOpts(widget.WidgetOpts.LayoutData(&widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			})),
			widget.SliderOpts.MinMax(1, 20),
			widget.SliderOpts.Images(res.images.sliderTrack, res.images.button),
			widget.SliderOpts.TrackPadding(2),
			widget.SliderOpts.HandleSize(20),
			widget.SliderOpts.PageSizeFunc(func() int {
				return ps
			}),
			widget.SliderOpts.ChangedHandler(func(args *widget.SliderChangedEventArgs) {
				text.Label = fmt.Sprintf("%d", args.Current)
			}),
		)
		sc.AddChild(s)
		sliders = append(sliders, s)

		text = widget.NewText(
			widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(&widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			})),
			widget.TextOpts.Text(fmt.Sprintf("%d", s.Current), res.fonts.face, res.colors.textIdle))
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

func toolTipPage(res *resources, toolTips *toolTipContents, toolTip *widget.ToolTip) *page {
	c := newPageContentContainer()

	c.AddChild(widget.NewText(
		widget.TextOpts.Text("Hover over these buttons to see their tool tips.", res.fonts.face, res.colors.textIdle)))

	bc := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(5))))
	c.AddChild(bc)

	for col := 0; col < 5; col++ {
		b := widget.NewButton(
			widget.ButtonOpts.Image(res.images.button),
			widget.ButtonOpts.Text(fmt.Sprintf("Button %d", col+1), res.fonts.face, res.colors.buttonText))

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

	stickyDelayedCheckbox := newCheckbox("Tool tips are sticky and delayed", func(args *widget.CheckboxChangedEventArgs) {
		toolTip.Sticky = args.State == widget.CheckboxChecked
		if args.State == widget.CheckboxChecked {
			toolTip.Delay = 800 * time.Millisecond
		} else {
			toolTip.Delay = 0
		}
	}, res)
	toolTips.Set(stickyDelayedCheckbox, "If enabled, tool tips do not show immediately and will not move with the cursor.")
	c.AddChild(stickyDelayedCheckbox)

	return &page{
		title:   "Tool Tip",
		content: c,
	}
}

func dragAndDropPage(res *resources, dnd *widget.DragAndDrop, drag *dragContents) *page {
	c := newPageContentContainer()

	dndContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(20),
		)),
	)
	c.AddChild(dndContainer)

	sourceContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(&widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ContainerOpts.BackgroundImage(res.images.scrollContainer.Idle),
		widget.ContainerOpts.Layout(widget.NewFillLayout(
			widget.FillLayoutOpts.Padding(widget.NewInsetsSimple(20)),
		)),
	)
	drag.addSource(sourceContainer)
	dndContainer.AddChild(sourceContainer)

	sourceContainer.AddChild(widget.NewText(
		widget.TextOpts.Text("Drag\nFrom\nHere", res.fonts.face, res.colors.textIdle),
		widget.TextOpts.Position(widget.TextPositionCenter),
	))

	targetContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(&widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ContainerOpts.BackgroundImage(res.images.scrollContainer.Idle),
		widget.ContainerOpts.Layout(widget.NewFillLayout(
			widget.FillLayoutOpts.Padding(widget.NewInsetsSimple(20)),
		)),
	)
	drag.addTarget(targetContainer)
	dndContainer.AddChild(targetContainer)

	targetContainer.AddChild(widget.NewText(
		widget.TextOpts.Text("Drop\nHere", res.fonts.face, res.colors.textIdle),
		widget.TextOpts.Position(widget.TextPositionCenter),
	))

	dnd.DroppedEvent.AddHandler(func(args interface{}) {
		a := args.(*widget.DragAndDropDroppedEventArgs)
		if !drag.isTarget(a.Target.GetWidget()) {
			return
		}

		targetContainer.BackgroundImage = res.images.button.Pressed

		time.AfterFunc(2*time.Second, func() {
			targetContainer.BackgroundImage = res.images.scrollContainer.Idle
		})
	})

	return &page{
		title:   "Drag & Drop",
		content: c,
	}
}

func textInputPage(res *resources) *page {
	c := newPageContentContainer()

	t := widget.NewTextInput(
		widget.TextInputOpts.Placeholder("Enter text here"),
		widget.TextInputOpts.WidgetOpts(widget.WidgetOpts.LayoutData(&widget.RowLayoutData{
			Stretch: true,
		})),
		widget.TextInputOpts.Image(&widget.TextInputImage{
			Idle:     res.images.scrollContainer.Idle,
			Disabled: res.images.scrollContainer.Disabled,
		}),
		widget.TextInputOpts.Color(&widget.TextInputColor{
			Idle:     res.colors.textIdle,
			Disabled: res.colors.textDisabled,
			Caret:    res.colors.textIdle,
		}),
		widget.TextInputOpts.Padding(widget.Insets{
			Left:   13,
			Right:  13,
			Top:    7,
			Bottom: 7,
		}),
		widget.TextInputOpts.Face(res.fonts.face),
		widget.TextInputOpts.CaretOpts(
			widget.CaretOpts.Size(res.fonts.face, 2),
		),
	)
	c.AddChild(t)

	c.AddChild(newSeparator(res, &widget.RowLayoutData{
		Stretch: true,
	}))

	c.AddChild(newCheckbox("Disabled", func(args *widget.CheckboxChangedEventArgs) {
		t.GetWidget().Disabled = args.State == widget.CheckboxChecked
	}, res))

	return &page{
		title:   "Text Input",
		content: c,
	}
}

func radioGroupPage(res *resources) *page {
	c := newPageContentContainer()

	var cbs []*widget.Checkbox
	for i := 0; i < 5; i++ {
		cb := widget.NewLabeledCheckbox(
			widget.LabeledCheckboxOpts.Spacing(6),
			widget.LabeledCheckboxOpts.CheckboxOpts(
				widget.CheckboxOpts.ButtonOpts(
					widget.ButtonOpts.Image(res.images.button),
					widget.ButtonOpts.GraphicPadding(widget.NewInsetsSimple(7)),
				),
				widget.CheckboxOpts.Image(res.images.checkbox)),
			widget.LabeledCheckboxOpts.LabelOpts(widget.LabelOpts.Text(fmt.Sprintf("Checkbox %d", i+1), res.fonts.face, res.colors.label)))
		c.AddChild(cb)
		cbs = append(cbs, cb.Checkbox())
	}

	widget.NewRadioGroup(widget.RadioGroupOpts.Checkboxes(cbs...))

	return &page{
		title:   "Radio Group",
		content: c,
	}
}

func windowPage(res *resources, ui func() *ebitenui.UI) *page {
	c := newPageContentContainer()

	b := widget.NewButton(
		widget.ButtonOpts.Image(res.images.button),
		widget.ButtonOpts.Text("Open Window", res.fonts.face, res.colors.buttonText),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			var rw ebitenui.RemoveWindowFunc

			wc := widget.NewContainer(
				widget.ContainerOpts.BackgroundImage(res.images.button.Disabled),
				widget.ContainerOpts.Layout(widget.NewRowLayout(
					widget.RowLayoutOpts.Direction(widget.DirectionVertical),
					widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(20)),
					widget.RowLayoutOpts.Spacing(15),
				)),
			)

			wc.AddChild(widget.NewText(
				widget.TextOpts.Text("Modal Window", res.fonts.titleFace, res.colors.textIdle),
			))

			wc.AddChild(widget.NewText(
				widget.TextOpts.Text("This window blocks all input to widgets below it.", res.fonts.face, res.colors.textIdle),
			))

			cb := widget.NewButton(
				widget.ButtonOpts.Image(res.images.button),
				widget.ButtonOpts.Text("Close", res.fonts.face, res.colors.buttonText),
				widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
					rw()
				}),
			)
			wc.AddChild(cb)

			w := widget.NewWindow(
				widget.WindowOpts.Modal(),
				widget.WindowOpts.Contents(wc),
			)

			ww, wh := ebiten.WindowSize()
			r := image.Rect(0, 0, ww*3/4, wh/3)
			r = r.Add(image.Point{ww / 4 / 2, wh * 2 / 3 / 2})
			w.SetLocation(r)

			rw = ui().AddWindow(w)
		}),
	)
	c.AddChild(b)

	return &page{
		title:   "Window",
		content: c,
	}
}
