package main

import (
	"fmt"
	"image"
	"time"

	"github.com/blizzy78/ebitenui"
	"github.com/blizzy78/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
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
		)
		c.AddChild(b)
		bs = append(bs, b)
	}

	c.AddChild(newSeparator(res, widget.RowLayoutData{
		Stretch: true,
	}))

	c.AddChild(newCheckbox("Disabled", func(args *widget.CheckboxChangedEventArgs) {
		for _, b := range bs {
			b.GetWidget().Disabled = args.State == widget.CheckboxChecked
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
		cb1.GetWidget().Disabled = args.State == widget.CheckboxChecked
		cb2.GetWidget().Disabled = args.State == widget.CheckboxChecked
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
		list1.GetWidget().Disabled = args.State == widget.CheckboxChecked
		list2.GetWidget().Disabled = args.State == widget.CheckboxChecked
		for _, b := range bs {
			b.GetWidget().Disabled = args.State == widget.CheckboxChecked
		}
	}, res))

	return &page{
		title:   "List",
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
		cb.GetWidget().Disabled = args.State == widget.CheckboxChecked
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
		tc := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(10))),
			widget.ContainerOpts.AutoDisableChildren())

		for j := 0; j < 3; j++ {
			b := widget.NewButton(
				widget.ButtonOpts.Image(res.button.image),
				widget.ButtonOpts.TextPadding(res.button.padding),
				widget.ButtonOpts.Text(fmt.Sprintf("Button %d on Tab %d", j+1, i+1), res.button.face, res.button.text))
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
		widget.TabBookOpts.TabButtonImage(res.tabBook.idleButton, res.tabBook.selectedButton),
		widget.TabBookOpts.TabButtonText(res.tabBook.buttonFace, res.tabBook.buttonText),
		widget.TabBookOpts.TabButtonOpts(widget.StateButtonOpts.ButtonOpts(widget.ButtonOpts.TextPadding(res.tabBook.buttonPadding))),
		widget.TabBookOpts.TabButtonSpacing(10),
		widget.TabBookOpts.Spacing(15))
	c.AddChild(t)

	c.AddChild(newSeparator(res, widget.RowLayoutData{
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

	pageSizes := []int{3, 10}
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
			})),
			widget.SliderOpts.MinMax(1, 20),
			widget.SliderOpts.Images(res.slider.trackImage, res.slider.handle),
			widget.SliderOpts.HandleSize(res.slider.handleSize),
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
			s.GetWidget().Parent().Disabled = args.State == widget.CheckboxChecked
		}
	}, res))

	return &page{
		title:   "Slider",
		content: c,
	}
}

func toolTipPage(res *uiResources, toolTips *toolTipContents, toolTip *widget.ToolTip) *page {
	c := newPageContentContainer()

	c.AddChild(widget.NewText(
		widget.TextOpts.Text("Hover over these buttons to see their tool tips.", res.text.face, res.text.idleColor)))

	bc := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(15))))
	c.AddChild(bc)

	for col := 0; col < 4; col++ {
		b := widget.NewButton(
			widget.ButtonOpts.Image(res.button.image),
			widget.ButtonOpts.TextPadding(res.button.padding),
			widget.ButtonOpts.Text(fmt.Sprintf("%s %d", string(rune('A'+col)), col+1), res.button.face, res.button.text))

		if col == 2 {
			b.GetWidget().Disabled = true
		}

		toolTips.Set(b, fmt.Sprintf("Tool tip for button %d", col+1))
		toolTips.widgetsWithTime = append(toolTips.widgetsWithTime, b)

		bc.AddChild(b)
	}

	c.AddChild(newSeparator(res, widget.RowLayoutData{
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

func dragAndDropPage(res *uiResources, dnd *widget.DragAndDrop, drag *dragContents) *page {
	c := newPageContentContainer()

	dndContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(30),
		)),
	)
	c.AddChild(dndContainer)

	sourcePanel := newSizedPanel(200, 200,
		widget.ContainerOpts.BackgroundImage(res.panel.image),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(widget.AnchorLayoutOpts.Padding(res.panel.padding))),
	)
	drag.addSource(sourcePanel)
	dndContainer.AddChild(sourcePanel)

	sourcePanel.Container().AddChild(widget.NewText(
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
		})),
		widget.TextOpts.Text("Drag\nFrom\nHere", res.text.face, res.text.disabledColor),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
	))

	targetPanel := newSizedPanel(200, 200,
		widget.ContainerOpts.BackgroundImage(res.panel.image),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(widget.AnchorLayoutOpts.Padding(res.panel.padding))),
	)
	drag.addTarget(targetPanel)
	dndContainer.AddChild(targetPanel)

	targetText := widget.NewText(
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
		})),
		widget.TextOpts.Text("Drop\nHere", res.text.face, res.text.disabledColor),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
	)

	targetPanel.Container().AddChild(targetText)

	dnd.DroppedEvent.AddHandler(func(args interface{}) {
		a := args.(*widget.DragAndDropDroppedEventArgs)
		if !drag.isTarget(a.Target.GetWidget()) {
			return
		}

		targetText.Label = "Thanks!"
		targetText.Color = res.text.idleColor

		time.AfterFunc(2500*time.Millisecond, func() {
			targetText.Label = "Drop\nHere"
			targetText.Color = res.text.disabledColor
		})
	})

	return &page{
		title:   "Drag & Drop",
		content: c,
	}
}

func textInputPage(res *uiResources, ui func() *ebitenui.UI) *page {
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
		widget.TextInputOpts.EnterFunc(func(text string, enable widget.TextInputEnable) {
			println("Enter:", text)
			enable(false) // llint: disable the TextInput widget, until the window (3) is closed
			openWindow3(res, ui, text, enable)
		}),
	}

	t := widget.NewTextInput(append(
		tOpts,
		widget.TextInputOpts.Placeholder("Enter text here"))...,
	)
	c.AddChild(t)

	tSecure := widget.NewTextInput(append(
		tOpts,
		widget.TextInputOpts.Placeholder("Enter secure text here"),
		widget.TextInputOpts.Secure(true))...,
	)
	c.AddChild(tSecure)

	c.AddChild(newSeparator(res, widget.RowLayoutData{
		Stretch: true,
	}))

	c.AddChild(newCheckbox("Disabled", func(args *widget.CheckboxChangedEventArgs) {
		t.GetWidget().Disabled = args.State == widget.CheckboxChecked
		tSecure.GetWidget().Disabled = args.State == widget.CheckboxChecked
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

	widget.NewRadioGroup(widget.RadioGroupOpts.Checkboxes(cbs...))

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
	var rw ebitenui.RemoveWindowFunc

	c := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(res.panel.image),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(res.panel.padding),
			widget.RowLayoutOpts.Spacing(15),
		)),
	)

	c.AddChild(widget.NewText(
		widget.TextOpts.Text("Modal Window", res.text.bigTitleFace, res.text.idleColor),
	))

	c.AddChild(widget.NewText(
		widget.TextOpts.Text("This window blocks all input to widgets below it.", res.text.face, res.text.idleColor),
	))

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

	w := widget.NewWindow(
		widget.WindowOpts.Modal(),
		widget.WindowOpts.Contents(c),
	)

	ww, wh := ebiten.WindowSize()
	r := image.Rect(0, 0, ww*3/4, wh/3)
	r = r.Add(image.Point{ww / 4 / 2, wh * 2 / 3 / 2})
	w.SetLocation(r)

	rw = ui().AddWindow(w)
}

func openWindow2(res *uiResources, ui func() *ebitenui.UI) {
	var rw ebitenui.RemoveWindowFunc

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
	)

	ww, wh := ebiten.WindowSize()
	r := image.Rect(0, 0, ww/2, wh/2)
	r = r.Add(image.Point{ww * 4 / 10, wh / 2 / 2})
	w.SetLocation(r)

	rw = ui().AddWindow(w)
}

func openWindow3(res *uiResources, ui func() *ebitenui.UI, text string, enable widget.TextInputEnable) {
	var rw ebitenui.RemoveWindowFunc

	c := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(res.panel.image),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(res.panel.padding),
			widget.RowLayoutOpts.Spacing(15),
		)),
	)

	c.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Text Entered: %v", text), res.text.bigTitleFace, res.text.idleColor),
	))

	cb := widget.NewButton(
		widget.ButtonOpts.Image(res.button.image),
		widget.ButtonOpts.TextPadding(res.button.padding),
		widget.ButtonOpts.Text("Close", res.button.face, res.button.text),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			rw()
			enable(true) // llint: re-enable the TextInput after the window is closed
		}),
	)
	c.AddChild(cb)

	w := widget.NewWindow(
		widget.WindowOpts.Modal(),
		widget.WindowOpts.Contents(c),
	)

	ww, wh := ebiten.WindowSize()
	r := image.Rect(0, 0, ww/2, wh/2)
	r = r.Add(image.Point{ww * 4 / 10, wh / 2 / 2})
	w.SetLocation(r)

	rw = ui().AddWindow(w)
}

func anchorLayoutPage(res *uiResources) *page {
	c := newPageContentContainer()

	p := newSizedPanel(300, 220,
		widget.ContainerOpts.BackgroundImage(res.panel.image),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			widget.AnchorLayoutOpts.Padding(res.panel.padding),
		)),
	)
	c.AddChild(p)

	sp := newSizedPanel(50, 50,
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{})),
		widget.ContainerOpts.BackgroundImage(res.panel.image),
	)
	p.Container().AddChild(sp.Container())

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

	widget.NewRadioGroup(
		widget.RadioGroupOpts.Checkboxes(hCBs...),
		widget.RadioGroupOpts.ChangedHandler(func(args *widget.RadioGroupChangedEventArgs) {
			ald := sp.Container().GetWidget().LayoutData.(widget.AnchorLayoutData)
			ald.HorizontalPosition = widget.AnchorLayoutPosition(indexCheckbox(hCBs, args.Active))
			sp.Container().GetWidget().LayoutData = ald
			p.Container().RequestRelayout()
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

	widget.NewRadioGroup(
		widget.RadioGroupOpts.Checkboxes(vCBs...),
		widget.RadioGroupOpts.ChangedHandler(func(args *widget.RadioGroupChangedEventArgs) {
			ald := sp.Container().GetWidget().LayoutData.(widget.AnchorLayoutData)
			ald.VerticalPosition = widget.AnchorLayoutPosition(indexCheckbox(vCBs, args.Active))
			sp.Container().GetWidget().LayoutData = ald
			p.Container().RequestRelayout()
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
		ald := sp.Container().GetWidget().LayoutData.(widget.AnchorLayoutData)
		ald.StretchHorizontal = args.State == widget.CheckboxChecked
		sp.Container().GetWidget().LayoutData = ald
		p.Container().RequestRelayout()

		hPosC.GetWidget().Disabled = args.State == widget.CheckboxChecked
	}, res)
	stretchC.AddChild(stretchHorizontalCheckbox)

	stretchVerticalCheckbox := newCheckbox("Vertical", func(args *widget.CheckboxChangedEventArgs) {
		ald := sp.Container().GetWidget().LayoutData.(widget.AnchorLayoutData)
		ald.StretchVertical = args.State == widget.CheckboxChecked
		sp.Container().GetWidget().LayoutData = ald
		p.Container().RequestRelayout()

		vPosC.GetWidget().Disabled = args.State == widget.CheckboxChecked
	}, res)
	stretchC.AddChild(stretchVerticalCheckbox)

	return &page{
		title:   "Anchor Layout",
		content: c,
	}
}

func indexCheckbox(cs []*widget.Checkbox, c *widget.Checkbox) int {
	for i, cb := range cs {
		if cb == c {
			return i
		}
	}
	return -1
}
