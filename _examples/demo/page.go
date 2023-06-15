package main

import (
	"fmt"
	"image"
	"time"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/input"
	"github.com/ebitenui/ebitenui/widget"
)

type page struct {
	title   string
	content widget.PreferredSizeLocateableWidget
}

func buttonPage(res *uiResources) *page {
	c := newPageContentContainer()

	bs := []*widget.Button{}
	for i := 0; i < 3; i++ {
		b := widget.NewButton(
			widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch: true,
			})),
			widget.ButtonOpts.Image(res.button.image),
			widget.ButtonOpts.Text(fmt.Sprintf("Button %d", i+1), res.button.face, res.button.text),
			widget.ButtonOpts.TextPadding(res.button.padding),
			widget.ButtonOpts.CursorEnteredHandler(func(args *widget.ButtonHoverEventArgs) { fmt.Println("Cursor Entered: " + args.Button.Text().Label) }),
			widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) { fmt.Println("Cursor Exited: " + args.Button.Text().Label) }),
		)
		c.AddChild(b)
		bs = append(bs, b)
	}

	c.AddChild(newSeparator(res, widget.RowLayoutData{
		Stretch: true,
	}))

	toggles := []*widget.Button{}
	for i := 0; i < 3; i++ {
		b := widget.NewButton(
			widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch: true,
			})),
			widget.ButtonOpts.Image(res.button.image),
			widget.ButtonOpts.Text(fmt.Sprintf("Toggle Button %d", i+1), res.button.face, res.button.text),
			widget.ButtonOpts.TextPadding(res.button.padding),
			widget.ButtonOpts.CursorEnteredHandler(func(args *widget.ButtonHoverEventArgs) { fmt.Println("Cursor Entered: " + args.Button.Text().Label) }),
			widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) { fmt.Println("Cursor Exited: " + args.Button.Text().Label) }),
		)
		c.AddChild(b)
		bs = append(bs, b)
		toggles = append(toggles, b)
	}
	elements := []widget.RadioGroupElement{}
	for _, cb := range toggles {
		elements = append(elements, cb)
	}
	widget.NewRadioGroup(widget.RadioGroupOpts.Elements(elements...))

	c.AddChild(newSeparator(res, widget.RowLayoutData{
		Stretch: true,
	}))

	c.AddChild(newCheckbox("Disabled", func(args *widget.CheckboxChangedEventArgs) {
		for _, b := range bs {
			b.GetWidget().Disabled = args.State == widget.WidgetChecked
		}
	}, res))

	return &page{
		title:   "Button",
		content: c,
	}
}

func checkboxPage(res *uiResources) *page {
	c := newPageContentContainer()

	cb1 := newCheckbox("Two-State Checkbox", nil, res)
	c.AddChild(cb1)

	cb2 := widget.NewLabeledCheckbox(
		widget.LabeledCheckboxOpts.Spacing(res.checkbox.spacing),
		widget.LabeledCheckboxOpts.CheckboxOpts(
			widget.CheckboxOpts.ButtonOpts(widget.ButtonOpts.Image(res.checkbox.image)),
			widget.CheckboxOpts.Image(res.checkbox.graphic),
			widget.CheckboxOpts.TriState()),
		widget.LabeledCheckboxOpts.LabelOpts(widget.LabelOpts.Text("Tri-State Checkbox", res.label.face, res.label.text)))
	c.AddChild(cb2)

	c.AddChild(newSeparator(res, widget.RowLayoutData{
		Stretch: true,
	}))

	c.AddChild(newCheckbox("Disabled", func(args *widget.CheckboxChangedEventArgs) {
		cb1.GetWidget().Disabled = args.State == widget.WidgetChecked
		cb2.GetWidget().Disabled = args.State == widget.WidgetChecked
	}, res))

	return &page{
		title:   "Checkbox",
		content: c,
	}
}

func listPage(res *uiResources) *page {
	c := newPageContentContainer()

	listsContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(3),
			widget.GridLayoutOpts.Stretch([]bool{true, false, true}, []bool{true}),
			widget.GridLayoutOpts.Spacing(10, 0))))
	c.AddChild(listsContainer)

	entries1 := []interface{}{"One", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten"}
	list1 := newList(entries1, res, widget.WidgetOpts.LayoutData(widget.GridLayoutData{
		MaxHeight: 220,
	}))
	listsContainer.AddChild(list1)

	buttonsContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(10),
		)))
	listsContainer.AddChild(buttonsContainer)

	bs := []*widget.Button{}
	for i := 0; i < 3; i++ {
		b := widget.NewButton(
			widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch: true,
			})),
			widget.ButtonOpts.Image(res.button.image),
			widget.ButtonOpts.TextPadding(res.button.padding),
			widget.ButtonOpts.Text(fmt.Sprintf("Action %d", i+1), res.button.face, res.button.text))
		buttonsContainer.AddChild(b)
		bs = append(bs, b)
	}

	entries2 := []interface{}{"Eleven", "Twelve", "Thirteen", "Fourteen", "Fifteen", "Sixteen", "Seventeen", "Eighteen", "Nineteen", "Twenty"}
	list2 := newList(entries2, res, widget.WidgetOpts.LayoutData(widget.GridLayoutData{
		MaxHeight: 220,
	}))
	listsContainer.AddChild(list2)

	c.AddChild(newSeparator(res, widget.RowLayoutData{
		Stretch: true,
	}))

	c.AddChild(newCheckbox("Disabled", func(args *widget.CheckboxChangedEventArgs) {
		list1.GetWidget().Disabled = args.State == widget.WidgetChecked
		list2.GetWidget().Disabled = args.State == widget.WidgetChecked
		for _, b := range bs {
			b.GetWidget().Disabled = args.State == widget.WidgetChecked
		}
	}, res))

	return &page{
		title:   "List",
		content: c,
	}
}

func textAreaPage(res *uiResources) *page {
	c := newPageContentContainer()

	textAreaContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch: true,
			}),
			widget.WidgetOpts.MinSize(0, 220),
		),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{true}),
			widget.GridLayoutOpts.Spacing(0, 0)),
		),
	)
	c.AddChild(textAreaContainer)

	textArea := newTextArea("Hello [color=FF0000] World! [/color]", res, widget.WidgetOpts.LayoutData(widget.GridLayoutData{
		MaxHeight: 220,
	}))
	textAreaContainer.AddChild(textArea)

	c.AddChild(newSeparator(res, widget.RowLayoutData{
		Stretch: true,
	}))
	verticalRows := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(10),
			widget.RowLayoutOpts.Direction(widget.DirectionVertical))),
		widget.ContainerOpts.AutoDisableChildren(),
	)
	c.AddChild(verticalRows)
	tOpts := []widget.TextInputOpt{
		widget.TextInputOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.TextInputOpts.Image(res.textInput.image),
		widget.TextInputOpts.Color(res.textInput.color),
		widget.TextInputOpts.Padding(widget.Insets{
			Left:   13,
			Right:  13,
			Top:    7,
			Bottom: 7,
		}),
		widget.TextInputOpts.Face(res.textInput.face),
		widget.TextInputOpts.CaretOpts(
			widget.CaretOpts.Size(res.textInput.face, 2),
		),
	}

	t := widget.NewTextInput(append(
		tOpts,
		widget.TextInputOpts.Placeholder("Enter text here"))...,
	)
	verticalRows.AddChild(t)

	row := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(10))),
		widget.ContainerOpts.AutoDisableChildren(),
	)
	verticalRows.AddChild(row)
	b := widget.NewButton(
		widget.ButtonOpts.Image(res.button.image),
		widget.ButtonOpts.TextPadding(res.button.padding),
		widget.ButtonOpts.Text("Prepend", res.button.face, res.button.text),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			textArea.PrependText(t.GetText())
		}),
	)
	row.AddChild(b)
	b = widget.NewButton(
		widget.ButtonOpts.Image(res.button.image),
		widget.ButtonOpts.TextPadding(res.button.padding),
		widget.ButtonOpts.Text("Append", res.button.face, res.button.text),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			textArea.AppendText(t.GetText())
		}),
	)
	row.AddChild(b)
	b = widget.NewButton(
		widget.ButtonOpts.Image(res.button.image),
		widget.ButtonOpts.TextPadding(res.button.padding),
		widget.ButtonOpts.Text("Set", res.button.face, res.button.text),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			textArea.SetText(t.GetText())
		}),
	)
	row.AddChild(b)

	return &page{
		title:   "Text Area",
		content: c,
	}
}

func comboButtonPage(res *uiResources) *page {
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

	c.AddChild(newSeparator(res, widget.RowLayoutData{
		Stretch: true,
	}))

	c.AddChild(newCheckbox("Disabled", func(args *widget.CheckboxChangedEventArgs) {
		cb.GetWidget().Disabled = args.State == widget.WidgetChecked
	}, res))

	return &page{
		title:   "Combo Button",
		content: c,
	}
}

func tabBookPage(res *uiResources) *page {
	c := newPageContentContainer()

	tabs := []*widget.TabBookTab{}

	for i := 0; i < 4; i++ {
		tab := widget.NewTabBookTab(fmt.Sprintf("Tab %d", i+1),
			widget.ContainerOpts.Layout(widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(10))),
			widget.ContainerOpts.AutoDisableChildren())

		for j := 0; j < 3; j++ {
			b := widget.NewButton(
				widget.ButtonOpts.Image(res.button.image),
				widget.ButtonOpts.TextPadding(res.button.padding),
				widget.ButtonOpts.Text(fmt.Sprintf("Button %d on Tab %d", j+1, i+1), res.button.face, res.button.text))
			tab.AddChild(b)
		}

		if i == 2 {
			tab.Disabled = true
		}

		tabs = append(tabs, tab)
	}

	t := widget.NewTabBook(
		widget.TabBookOpts.Tabs(tabs...),
		widget.TabBookOpts.TabButtonImage(res.button.image),
		widget.TabBookOpts.TabButtonText(res.tabBook.buttonFace, res.tabBook.buttonText),
		widget.TabBookOpts.TabButtonOpts(widget.ButtonOpts.TextPadding(res.tabBook.buttonPadding)),
		widget.TabBookOpts.TabButtonSpacing(10),
		widget.TabBookOpts.Spacing(15))
	c.AddChild(t)

	c.AddChild(newSeparator(res, widget.RowLayoutData{
		Stretch: true,
	}))

	c.AddChild(newCheckbox("Disabled", func(args *widget.CheckboxChangedEventArgs) {
		t.GetWidget().Disabled = args.State == widget.WidgetChecked
	}, res))

	return &page{
		title:   "Tab Book",
		content: c,
	}
}

func gridLayoutPage(res *uiResources) *page {
	c := newPageContentContainer()

	bc := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(4),
			widget.GridLayoutOpts.Stretch([]bool{true, true, true, true}, nil),
			widget.GridLayoutOpts.Spacing(10, 10))))
	c.AddChild(bc)

	i := 0
	for row := 0; row < 3; row++ {
		for col := 0; col < 4; col++ {
			b := widget.NewButton(
				widget.ButtonOpts.Image(res.button.image),
				widget.ButtonOpts.Text(fmt.Sprintf("%s %d", string(rune('A'+i)), i+1), res.button.face, res.button.text))
			bc.AddChild(b)

			i++
		}
	}

	return &page{
		title:   "Grid Layout",
		content: c,
	}
}

func rowLayoutPage(res *uiResources) *page {
	c := newPageContentContainer()

	c.AddChild(widget.NewText(
		widget.TextOpts.Text("Horizontal", res.text.face, res.text.idleColor)))

	bc := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(5))))
	c.AddChild(bc)

	for col := 0; col < 5; col++ {
		b := widget.NewButton(
			widget.ButtonOpts.Image(res.button.image),
			widget.ButtonOpts.TextPadding(res.button.padding),
			widget.ButtonOpts.Text(fmt.Sprintf("%s %d", string(rune('A'+col)), col+1), res.button.face, res.button.text))
		bc.AddChild(b)
	}

	c.AddChild(widget.NewText(
		widget.TextOpts.Text("Vertical", res.text.face, res.text.idleColor)))

	bc = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(5))))
	c.AddChild(bc)

	labels := []string{"Tiny", "Medium", "Very Large"}
	for _, l := range labels {
		b := widget.NewButton(
			widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch: true,
			})),
			widget.ButtonOpts.Image(res.button.image),
			widget.ButtonOpts.TextPadding(res.button.padding),
			widget.ButtonOpts.Text(l, res.button.face, res.button.text))
		bc.AddChild(b)
	}

	return &page{
		title:   "Row Layout",
		content: c,
	}
}

func sliderPage(res *uiResources) *page {
	c := newPageContentContainer()

	pageSizes := []int{3, 1}
	sliders := []*widget.Slider{}

	for _, ps := range pageSizes {
		ps := ps

		sc := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewRowLayout(
				widget.RowLayoutOpts.Spacing(10))),
			widget.ContainerOpts.AutoDisableChildren(),
		)
		c.AddChild(sc)

		var text *widget.Label

		s := widget.NewSlider(
			widget.SliderOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}), widget.WidgetOpts.MinSize(200, 6)),
			widget.SliderOpts.MinMax(1, 10),
			widget.SliderOpts.Images(res.slider.trackImage, res.slider.handle),
			widget.SliderOpts.FixedHandleSize(res.slider.handleSize),
			widget.SliderOpts.TrackOffset(5),
			widget.SliderOpts.PageSizeFunc(func() int {
				return ps
			}),
			widget.SliderOpts.ChangedHandler(func(args *widget.SliderChangedEventArgs) {
				text.Label = fmt.Sprintf("%d", args.Current)
			}),
		)
		sc.AddChild(s)
		sliders = append(sliders, s)

		text = widget.NewLabel(
			widget.LabelOpts.TextOpts(widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}))),
			widget.LabelOpts.Text(fmt.Sprintf("%d", s.Current), res.label.face, res.label.text),
		)
		sc.AddChild(text)
	}

	c.AddChild(newSeparator(res, widget.RowLayoutData{
		Stretch: true,
	}))

	c.AddChild(newCheckbox("Disabled", func(args *widget.CheckboxChangedEventArgs) {
		for _, s := range sliders {
			s.GetWidget().Parent().Disabled = args.State == widget.WidgetChecked
		}
	}, res))

	return &page{
		title:   "Slider",
		content: c,
	}
}

func progressBarPage(res *uiResources) *page {
	c := newPageContentContainer()

	sc := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(10))),
		widget.ContainerOpts.AutoDisableChildren(),
	)
	c.AddChild(sc)

	var text *widget.Label

	progressBar := widget.NewProgressBar(
		widget.ProgressBarOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter}),
			widget.WidgetOpts.MinSize(200, 30),
		),
		widget.ProgressBarOpts.Values(0, 20, 20),
		widget.ProgressBarOpts.TrackPadding(widget.Insets{
			Top:    3,
			Bottom: 3,
			Left:   2,
			Right:  2,
		}),
		widget.ProgressBarOpts.Images(res.progressBar.trackImage, res.progressBar.fillImage),
	)

	text = widget.NewLabel(
		widget.LabelOpts.TextOpts(widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter)),
		widget.LabelOpts.Text(fmt.Sprintf("%d", progressBar.GetCurrent()), res.label.face, res.label.text),
	)
	stackedLayout := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewStackedLayout()),
	)
	stackedLayout.AddChild(progressBar)
	stackedLayout.AddChild(text)
	sc.AddChild(stackedLayout)
	c.AddChild(newSeparator(res, widget.RowLayoutData{
		Stretch: true,
	}))
	sc = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(10))),
		widget.ContainerOpts.AutoDisableChildren(),
	)
	c.AddChild(sc)
	b := widget.NewButton(
		widget.ButtonOpts.Image(res.button.image),
		widget.ButtonOpts.TextPadding(res.button.padding),
		widget.ButtonOpts.Text("+ 1", res.button.face, res.button.text),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			progressBar.SetCurrent(progressBar.GetCurrent() + 1)
			text.Label = fmt.Sprintf("%d", progressBar.GetCurrent())
		}),
	)
	sc.AddChild(b)
	b = widget.NewButton(
		widget.ButtonOpts.Image(res.button.image),
		widget.ButtonOpts.TextPadding(res.button.padding),
		widget.ButtonOpts.Text("- 1", res.button.face, res.button.text),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			progressBar.SetCurrent(progressBar.GetCurrent() - 1)
			text.Label = fmt.Sprintf("%d", progressBar.GetCurrent())
		}),
	)
	sc.AddChild(b)
	return &page{
		title:   "Progress Bar",
		content: c,
	}
}

func toolTipPage(res *uiResources) *page {
	c := newPageContentContainer()

	c.AddChild(widget.NewText(
		widget.TextOpts.Text("Hover over these buttons to see their tool tips.", res.text.face, res.text.idleColor)))

	bc := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(15))))
	c.AddChild(bc)

	showTimeCheckbox := newCheckbox("Show additional infos in tool tips", func(args *widget.CheckboxChangedEventArgs) {

	}, res)

	for col := 0; col < 4; col++ {
		tt := widget.NewContainer(
			widget.ContainerOpts.BackgroundImage(res.toolTip.background),
			widget.ContainerOpts.Layout(widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Padding(res.toolTip.padding),
				widget.RowLayoutOpts.Spacing(2),
			)))

		text := widget.NewText(
			widget.TextOpts.Text(fmt.Sprintf("Tool tip for button %d", col+1), res.toolTip.face, res.toolTip.color),
		)
		tt.AddChild(text)
		timeTxt := widget.NewText(
			widget.TextOpts.Text("", res.toolTip.face, res.toolTip.color),
		)
		tt.AddChild(timeTxt)
		b := widget.NewButton(
			widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.ToolTip(widget.NewToolTip(
				widget.ToolTipOpts.Content(tt),
				widget.ToolTipOpts.ToolTipUpdater(func(c *widget.Container) {
					if showTimeCheckbox.Checkbox().State() == widget.WidgetChecked {
						c.Children()[1].(*widget.Text).Label = time.Now().Local().Format("2006-01-02 15:04:05")
					} else {
						c.Children()[1].(*widget.Text).Label = ""
					}
				}),
			))),
			widget.ButtonOpts.Image(res.button.image),
			widget.ButtonOpts.TextPadding(res.button.padding),
			widget.ButtonOpts.Text(fmt.Sprintf("%s %d", string(rune('A'+col)), col+1), res.button.face, res.button.text))

		if col == 2 {
			b.GetWidget().Disabled = true
		}

		bc.AddChild(b)
	}

	c.AddChild(newSeparator(res, widget.RowLayoutData{
		Stretch: true,
	}))

	c.AddChild(showTimeCheckbox)
	stickyDelayedCheckbox := newCheckbox("Tool tips are sticky and delayed", func(args *widget.CheckboxChangedEventArgs) {
		children := bc.Children()
		for i := range children {
			if args.State == widget.WidgetChecked {
				children[i].GetWidget().ToolTip.Delay = 800 * time.Millisecond
				children[i].GetWidget().ToolTip.Position = widget.TOOLTIP_POS_CURSOR_STICKY
			} else {
				children[i].GetWidget().ToolTip.Delay = 0
				children[i].GetWidget().ToolTip.Position = widget.TOOLTIP_POS_CURSOR_FOLLOW
			}
		}
	}, res)
	c.AddChild(stickyDelayedCheckbox)

	return &page{
		title:   "Tool Tip",
		content: c,
	}
}

func dragAndDropPage(res *uiResources) *page {
	c := newPageContentContainer()

	dnd := widget.NewDragAndDrop(
		widget.DragAndDropOpts.ContentsCreater(&dragContents{
			res: res,
		}),
	)

	dndContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(30),
		)),
	)
	c.AddChild(dndContainer)

	sourcePanel := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(200, 200),
			widget.WidgetOpts.EnableDragAndDrop(dnd),
		),
		widget.ContainerOpts.BackgroundImage(res.panel.image),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(widget.AnchorLayoutOpts.Padding(res.panel.padding))),
	)
	dndContainer.AddChild(sourcePanel)

	sourcePanel.AddChild(widget.NewText(
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
		})),
		widget.TextOpts.Text("Drag\nFrom\nHere", res.text.face, res.text.disabledColor),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
	))

	targetText := widget.NewText(
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
		})),
		widget.TextOpts.Text("Drop\nHere", res.text.face, res.text.disabledColor),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
	)

	targetPanel := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(200, 200),
			widget.WidgetOpts.CanDrop(func(args *widget.DragAndDropDroppedEventArgs) bool {
				return true
			}),
			widget.WidgetOpts.Dropped(func(args *widget.DragAndDropDroppedEventArgs) {
				targetText.Label = "Thanks!"
				targetText.Color = res.text.idleColor

				time.AfterFunc(2500*time.Millisecond, func() {
					targetText.Label = "Drop\nHere"
					targetText.Color = res.text.disabledColor
				})
			}),
		),
		widget.ContainerOpts.BackgroundImage(res.panel.image),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(widget.AnchorLayoutOpts.Padding(res.panel.padding))),
	)
	dndContainer.AddChild(targetPanel)

	targetPanel.AddChild(targetText)

	return &page{
		title:   "Drag & Drop",
		content: c,
	}
}

func textInputPage(res *uiResources) *page {
	c := newPageContentContainer()

	tOpts := []widget.TextInputOpt{
		widget.TextInputOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.TextInputOpts.Image(res.textInput.image),
		widget.TextInputOpts.Color(res.textInput.color),
		widget.TextInputOpts.Padding(widget.Insets{
			Left:   13,
			Right:  13,
			Top:    7,
			Bottom: 7,
		}),
		widget.TextInputOpts.Face(res.textInput.face),
		widget.TextInputOpts.CaretOpts(
			widget.CaretOpts.Size(res.textInput.face, 2),
		),
	}

	label := widget.NewLabel(widget.LabelOpts.Text("", res.label.face, res.label.text))
	t := widget.NewTextInput(append(
		tOpts,
		widget.TextInputOpts.Placeholder("Enter text here"),
		widget.TextInputOpts.AllowDuplicateSubmit(true),
		widget.TextInputOpts.SubmitHandler(func(args *widget.TextInputChangedEventArgs) {
			label.Label = fmt.Sprint("Unsecured Text Box Submitted: ", args.InputText)
			fmt.Println(label.Label)
		}))...,
	)
	c.AddChild(t)

	tSecure := widget.NewTextInput(append(
		tOpts,
		widget.TextInputOpts.Placeholder("Enter secure text here"),
		widget.TextInputOpts.Secure(true),
		widget.TextInputOpts.ClearOnSubmit(true),
		widget.TextInputOpts.IgnoreEmptySubmit(true),
		widget.TextInputOpts.SubmitHandler(func(args *widget.TextInputChangedEventArgs) {
			label.Label = fmt.Sprint("Secured Text Box Submitted: ", args.InputText)
			fmt.Println(label.Label)
		}))...,
	)
	c.AddChild(tSecure)
	c.AddChild(label)

	c.AddChild(newSeparator(res, widget.RowLayoutData{
		Stretch: true,
	}))

	c.AddChild(newCheckbox("Disabled", func(args *widget.CheckboxChangedEventArgs) {
		t.GetWidget().Disabled = args.State == widget.WidgetChecked
		tSecure.GetWidget().Disabled = args.State == widget.WidgetChecked
	}, res))

	return &page{
		title:   "Text Input",
		content: c,
	}
}

func radioGroupPage(res *uiResources) *page {
	c := newPageContentContainer()

	var cbs []*widget.Checkbox
	for i := 0; i < 5; i++ {
		cb := newCheckbox(fmt.Sprintf("Checkbox %d", i+1), nil, res)
		c.AddChild(cb)
		cbs = append(cbs, cb.Checkbox())
	}

	elements := []widget.RadioGroupElement{}
	for _, cb := range cbs {
		elements = append(elements, cb)
	}
	widget.NewRadioGroup(
		widget.RadioGroupOpts.Elements(elements...),
		widget.RadioGroupOpts.InitialElement(elements[2]),
	)

	return &page{
		title:   "Radio Group",
		content: c,
	}
}

func windowPage(res *uiResources, ui func() *ebitenui.UI) *page {
	c := newPageContentContainer()

	b := widget.NewButton(
		widget.ButtonOpts.Image(res.button.image),
		widget.ButtonOpts.TextPadding(res.button.padding),
		widget.ButtonOpts.Text("Open Window", res.button.face, res.button.text),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			openWindow(res, ui)
		}),
	)
	c.AddChild(b)

	return &page{
		title:   "Window",
		content: c,
	}
}

func openWindow(res *uiResources, ui func() *ebitenui.UI) {
	var rw widget.RemoveWindowFunc
	var window *widget.Window

	titleBar := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(res.panel.titleBar),
		widget.ContainerOpts.Layout(widget.NewGridLayout(widget.GridLayoutOpts.Columns(3), widget.GridLayoutOpts.Stretch([]bool{true, false, false}, []bool{true}), widget.GridLayoutOpts.Padding(widget.Insets{
			Left:   30,
			Right:  5,
			Top:    6,
			Bottom: 5,
		}))))

	titleBar.AddChild(widget.NewText(
		widget.TextOpts.Text("Modal Window", res.text.titleFace, res.textInput.color.Idle),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionCenter),
	))

	titleBar.AddChild(widget.NewButton(
		widget.ButtonOpts.Image(res.button.image),
		widget.ButtonOpts.TextPadding(res.button.padding),
		widget.ButtonOpts.Text("X", res.button.face, res.button.text),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			rw()
		}),
		widget.ButtonOpts.TabOrder(99),
	))

	c := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(res.panel.image),
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(1),
				widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, true, false}),
				widget.GridLayoutOpts.Padding(res.panel.padding),
				widget.GridLayoutOpts.Spacing(0, 15),
			),
		),
	)

	c.AddChild(widget.NewText(
		widget.TextOpts.Text("This window blocks all input to widgets below it.", res.text.face, res.text.idleColor),
	))

	tOpts := []widget.TextInputOpt{
		widget.TextInputOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.TextInputOpts.Image(res.textInput.image),
		widget.TextInputOpts.Color(res.textInput.color),
		widget.TextInputOpts.Padding(widget.Insets{
			Left:   13,
			Right:  13,
			Top:    7,
			Bottom: 7,
		}),
		widget.TextInputOpts.Face(res.textInput.face),
		widget.TextInputOpts.CaretOpts(
			widget.CaretOpts.Size(res.textInput.face, 2),
		),
	}

	t := widget.NewTextInput(append(
		tOpts,
		widget.TextInputOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
		})),
		widget.TextInputOpts.Placeholder("Enter text here"))...,
	)
	textContainer := widget.NewContainer(widget.ContainerOpts.Layout(widget.NewAnchorLayout()))
	textContainer.AddChild(t)
	c.AddChild(textContainer)

	bc := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(15),
		)),
	)
	c.AddChild(bc)

	o2b := widget.NewButton(
		widget.ButtonOpts.Image(res.button.image),
		widget.ButtonOpts.TextPadding(res.button.padding),
		widget.ButtonOpts.Text("Open Another", res.button.face, res.button.text),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			openWindow2(res, ui)
		}),
	)
	bc.AddChild(o2b)

	cb := widget.NewButton(
		widget.ButtonOpts.Image(res.button.image),
		widget.ButtonOpts.TextPadding(res.button.padding),
		widget.ButtonOpts.Text("Close", res.button.face, res.button.text),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			rw()
		}),
	)
	bc.AddChild(cb)

	window = widget.NewWindow(
		widget.WindowOpts.Modal(),
		widget.WindowOpts.Contents(c),
		widget.WindowOpts.TitleBar(titleBar, 30),
		widget.WindowOpts.Draggable(),
		widget.WindowOpts.Resizeable(),
		widget.WindowOpts.MinSize(500, 200),
		widget.WindowOpts.MaxSize(700, 400),
		widget.WindowOpts.ResizeHandler(func(args *widget.WindowChangedEventArgs) {
			fmt.Println("Resize: ", args.Rect)
		}),
		widget.WindowOpts.MoveHandler(func(args *widget.WindowChangedEventArgs) {
			fmt.Println("Move: ", args.Rect)
		}),
	)
	windowSize := input.GetWindowSize()
	r := image.Rect(0, 0, 550, 250)
	r = r.Add(image.Point{windowSize.X / 4 / 2, windowSize.Y * 2 / 3 / 2})
	window.SetLocation(r)

	rw = ui().AddWindow(window)
}

func openWindow2(res *uiResources, ui func() *ebitenui.UI) {
	var rw widget.RemoveWindowFunc

	c := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(res.panel.image),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(res.panel.padding),
			widget.RowLayoutOpts.Spacing(15),
		)),
	)

	c.AddChild(widget.NewText(
		widget.TextOpts.Text("Second Window", res.text.bigTitleFace, res.text.idleColor),
	))

	cb := widget.NewButton(
		widget.ButtonOpts.Image(res.button.image),
		widget.ButtonOpts.TextPadding(res.button.padding),
		widget.ButtonOpts.Text("Close", res.button.face, res.button.text),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			rw()
		}),
	)
	c.AddChild(cb)

	w := widget.NewWindow(
		widget.WindowOpts.Modal(),
		widget.WindowOpts.Contents(c),
		widget.WindowOpts.CloseMode(widget.CLICK_OUT),
	)

	windowSize := input.GetWindowSize()
	r := image.Rect(0, 0, windowSize.X/2, windowSize.Y/2)
	r = r.Add(image.Point{windowSize.X * 4 / 10, windowSize.Y / 2 / 2})
	w.SetLocation(r)

	rw = ui().AddWindow(w)
}

func anchorLayoutPage(res *uiResources) *page {
	c := newPageContentContainer()

	p := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(300, 220)),
		widget.ContainerOpts.BackgroundImage(res.panel.image),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			widget.AnchorLayoutOpts.Padding(res.panel.padding),
		)),
	)
	c.AddChild(p)

	sp := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(50, 50)),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{})),
		widget.ContainerOpts.BackgroundImage(res.panel.image),
	)
	p.AddChild(sp)

	c.AddChild(newSeparator(res, widget.RowLayoutData{
		Stretch: true,
	}))

	posC := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(50),
		)),
	)
	c.AddChild(posC)

	hPosC := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(10),
		)),
		widget.ContainerOpts.AutoDisableChildren(),
	)
	posC.AddChild(hPosC)

	hPosC.AddChild(widget.NewLabel(widget.LabelOpts.Text("Horizontal", res.label.face, res.label.text)))

	labels := []string{"Start", "Center", "End"}
	hCBs := []*widget.Checkbox{}
	for _, l := range labels {
		cb := newCheckbox(l, nil, res)
		hPosC.AddChild(cb)
		hCBs = append(hCBs, cb.Checkbox())
	}
	elements := []widget.RadioGroupElement{}
	for _, cb := range hCBs {
		elements = append(elements, cb)
	}
	widget.NewRadioGroup(
		widget.RadioGroupOpts.Elements(elements...),
		widget.RadioGroupOpts.ChangedHandler(func(args *widget.RadioGroupChangedEventArgs) {
			ald := sp.GetWidget().LayoutData.(widget.AnchorLayoutData)
			ald.HorizontalPosition = widget.AnchorLayoutPosition(indexCheckbox(hCBs, args.Active))
			sp.GetWidget().LayoutData = ald
			p.RequestRelayout()
		}),
	)

	vPosC := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(10),
		)),
		widget.ContainerOpts.AutoDisableChildren(),
	)
	posC.AddChild(vPosC)

	vPosC.AddChild(widget.NewLabel(widget.LabelOpts.Text("Vertical", res.label.face, res.label.text)))

	vCBs := []*widget.Checkbox{}
	for _, l := range labels {
		cb := newCheckbox(l, nil, res)
		vPosC.AddChild(cb)
		vCBs = append(vCBs, cb.Checkbox())
	}
	vElements := []widget.RadioGroupElement{}
	for _, cb := range vCBs {
		vElements = append(vElements, cb)
	}
	widget.NewRadioGroup(
		widget.RadioGroupOpts.Elements(vElements...),
		widget.RadioGroupOpts.ChangedHandler(func(args *widget.RadioGroupChangedEventArgs) {
			ald := sp.GetWidget().LayoutData.(widget.AnchorLayoutData)
			ald.VerticalPosition = widget.AnchorLayoutPosition(indexCheckbox(vCBs, args.Active))
			sp.GetWidget().LayoutData = ald
			p.RequestRelayout()
		}),
	)

	stretchC := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(10),
		)),
	)
	posC.AddChild(stretchC)

	stretchC.AddChild(widget.NewText(widget.TextOpts.Text("Stretch", res.text.face, res.text.idleColor)))

	stretchHorizontalCheckbox := newCheckbox("Horizontal", func(args *widget.CheckboxChangedEventArgs) {
		ald := sp.GetWidget().LayoutData.(widget.AnchorLayoutData)
		ald.StretchHorizontal = args.State == widget.WidgetChecked
		sp.GetWidget().LayoutData = ald
		p.RequestRelayout()

		hPosC.GetWidget().Disabled = args.State == widget.WidgetChecked
	}, res)
	stretchC.AddChild(stretchHorizontalCheckbox)

	stretchVerticalCheckbox := newCheckbox("Vertical", func(args *widget.CheckboxChangedEventArgs) {
		ald := sp.GetWidget().LayoutData.(widget.AnchorLayoutData)
		ald.StretchVertical = args.State == widget.WidgetChecked
		sp.GetWidget().LayoutData = ald
		p.RequestRelayout()

		vPosC.GetWidget().Disabled = args.State == widget.WidgetChecked
	}, res)
	stretchC.AddChild(stretchVerticalCheckbox)

	return &page{
		title:   "Anchor Layout",
		content: c,
	}
}

func indexCheckbox(cs []*widget.Checkbox, c widget.RadioGroupElement) int {
	for i, cb := range cs {
		if cb == c {
			return i
		}
	}
	return -1
}
