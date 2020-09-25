package main

import (
	"fmt"
	img "image"
	"io/ioutil"
	"log"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"golang.org/x/image/font"

	"image/color"
	_ "image/png"

	"github.com/blizzy78/ebitenui"
	"github.com/blizzy78/ebitenui/image"
	"github.com/blizzy78/ebitenui/widget"
)

type game struct {
	ui *ebitenui.UI
}

// var stopProfiler func()
// var firstRender = true

func main() {
	// defer func() {
	// 	stopProfiler()
	// }()

	ebiten.SetWindowSize(640, 1000)
	ebiten.SetWindowTitle("Hello, World!")
	ebiten.SetWindowResizable(true)

	game := game{
		ui: createUI(),
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}

func createUI() *ebitenui.UI {
	images, err := loadImages()
	if err != nil {
		panic(err)
	}

	fontFace, err := loadFont()
	if err != nil {
		panic(err)
	}

	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(20)),
			widget.RowLayoutOpts.Spacing(10))),
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.White)))

	var button1 *widget.Button
	var button2 *widget.Button

	button1 = widget.NewButton(
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(&widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ButtonOpts.Image(images.button),

		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			button2.GetWidget().Disabled = !button2.GetWidget().Disabled
		}))

	rootContainer.AddChild(button1)

	button2 = widget.NewButton(
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(&widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ButtonOpts.Image(images.button),
		widget.ButtonOpts.Text("foobar\nTy", fontFace, &widget.ButtonTextColor{
			Idle:     color.Black,
			Disabled: color.RGBA{128, 128, 128, 255},
		}),

		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			button1.GetWidget().Disabled = !button1.GetWidget().Disabled
		}))

	rootContainer.AddChild(button2)

	label := widget.NewText(
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(&widget.RowLayoutData{
			Stretch: true,
		})),
		widget.TextOpts.Text("hallo", fontFace, color.Black))

	rootContainer.AddChild(label)

	rootContainer.AddChild(widget.NewButton(
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(&widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ButtonOpts.Image(images.button),
		widget.ButtonOpts.Text("bleh", fontFace, &widget.ButtonTextColor{
			Idle:     color.Black,
			Disabled: color.RGBA{128, 128, 128, 255},
		})))

	rootContainer.AddChild(widget.NewComboButton(
		widget.ComboButtonOpts.ButtonOpts(widget.ButtonOpts.Image(images.button)),
		widget.ComboButtonOpts.ButtonOpts(widget.ButtonOpts.TextAndImage("Combo", fontFace, images.arrowDown, &widget.ButtonTextColor{
			Idle:     color.Black,
			Disabled: color.RGBA{128, 128, 128, 255},
		})),
		widget.ComboButtonOpts.Content(widget.NewButton(
			widget.ButtonOpts.Image(images.button),
			widget.ButtonOpts.Text("foobar qux", fontFace, &widget.ButtonTextColor{
				Idle:     color.Black,
				Disabled: color.RGBA{128, 128, 128, 255},
			})))))

	entries := []interface{}{}
	for i := 1; i <= 20; i++ {
		entries = append(entries, i)
	}

	rootContainer.AddChild(widget.NewListComboButton(
		widget.ListComboButtonOpts.SelectComboButtonOpts(widget.SelectComboButtonOpts.ComboButtonOpts(widget.ComboButtonOpts.ButtonOpts(widget.ButtonOpts.Image(images.button)))),
		widget.ListComboButtonOpts.Text(fontFace, images.arrowDown, &widget.ButtonTextColor{
			Idle:     color.Black,
			Disabled: color.RGBA{128, 128, 128, 255},
		}),
		widget.ListComboButtonOpts.ListOpts(widget.ListOpts.ScrollContainerOpts(widget.ScrollContainerOpts.Image(images.scrollContainer))),
		widget.ListComboButtonOpts.ListOpts(widget.ListOpts.ScrollContainerOpts(widget.ScrollContainerOpts.Padding(widget.NewInsetsSimple(2)))),
		widget.ListComboButtonOpts.ListOpts(widget.ListOpts.SliderOpts(widget.SliderOpts.Images(images.sliderTrack, images.button))),
		widget.ListComboButtonOpts.ListOpts(widget.ListOpts.Entries(entries)),
		widget.ListComboButtonOpts.EntryLabelFunc(
			func(e interface{}) string {
				return fmt.Sprintf("Entry %d", e.(int))
			},
			func(e interface{}) string {
				return fmt.Sprintf("Entry %d", e.(int))
			}),
		widget.ListComboButtonOpts.ListOpts(widget.ListOpts.EntryFontFace(fontFace)),
		widget.ListComboButtonOpts.ListOpts(widget.ListOpts.EntryColor(&widget.ListEntryColor{
			Unselected:         color.Black,
			Selected:           color.Black,
			SelectedBackground: color.RGBA{128, 128, 128, 255},

			DisabledUnselected:         color.Black,
			DisabledSelected:           color.Black,
			DisabledSelectedBackground: color.RGBA{128, 128, 128, 255},
		})),

		widget.ListComboButtonOpts.EntrySelectedHandler(func(args *widget.ListComboButtonEntrySelectedEventArgs) {
			println("entry selected:", args.Entry.(int))
		}),
	))

	slider := widget.NewSlider(
		widget.SliderOpts.WidgetOpts(widget.WidgetOpts.LayoutData(&widget.RowLayoutData{
			Stretch:  true,
			MaxWidth: 250,
		})),
		widget.SliderOpts.Direction(widget.DirectionHorizontal),
		widget.SliderOpts.Images(images.sliderTrack, images.button),

		widget.SliderOpts.ChangedHandler(func(args *widget.SliderChangedEventArgs) {
			label.Label = fmt.Sprintf("%d", args.Current)
		}))

	rootContainer.AddChild(slider)

	rootContainer.AddChild(widget.NewButton(
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(&widget.RowLayoutData{
			Stretch:  true,
			MaxWidth: 250,
		})),
		widget.ButtonOpts.Image(images.button),
		widget.ButtonOpts.Text("Disable", fontFace, &widget.ButtonTextColor{
			Idle:     color.Black,
			Disabled: color.RGBA{128, 128, 128, 255},
		}),

		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			slider.GetWidget().Disabled = !slider.GetWidget().Disabled

			if !slider.GetWidget().Disabled {
				slider.Current = 1
			}

			t := args.Button.Text()
			if slider.GetWidget().Disabled {
				t.Label = "Enable & Reset"
			} else {
				t.Label = "Disable"
			}
		})))

	rootContainer.AddChild(widget.NewList(
		widget.ListOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(&widget.RowLayoutData{
			Stretch:   true,
			MaxHeight: 200,
		}))),
		widget.ListOpts.Entries(entries),
		widget.ListOpts.EntryLabelFunc(func(e interface{}) string {
			return fmt.Sprintf("Entry %d", e.(int))
		}),
		widget.ListOpts.EntryColor(&widget.ListEntryColor{
			Unselected:         color.Black,
			Selected:           color.Black,
			SelectedBackground: color.RGBA{128, 128, 128, 255},

			DisabledUnselected:         color.Black,
			DisabledSelected:           color.Black,
			DisabledSelectedBackground: color.RGBA{128, 128, 128, 255},
		}),
		widget.ListOpts.EntryFontFace(fontFace),
		widget.ListOpts.ScrollContainerOpts(widget.ScrollContainerOpts.Image(images.scrollContainer)),
		widget.ListOpts.SliderOpts(widget.SliderOpts.Images(images.sliderTrack, images.button)),
		widget.ListOpts.ControlWidgetSpacing(2),
		widget.ListOpts.ScrollContainerOpts(widget.ScrollContainerOpts.Padding(widget.NewInsetsSimple(2))),

		widget.ListOpts.EntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			if args.Entry != args.PreviousEntry {
				println("entry selected: ", args.Entry.(int))
			}
		})))

	rootContainer.AddChild(widget.NewText(
		widget.TextOpts.Text("test", fontFace, color.Black)))

	container := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(&widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Padding(widget.NewInsetsSimple(10)),
			widget.GridLayoutOpts.Columns(3),
			widget.GridLayoutOpts.Stretch([]bool{false, true, true}, nil),
			widget.GridLayoutOpts.Spacing(2, 2))),
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.RGBA{0, 0, 255, 255})))

	for i := 0; i < 9; i++ {
		var c color.Color
		if i%2 == 0 {
			c = color.RGBA{255, 0, 0, 255}
		} else {
			c = color.RGBA{0, 255, 0, 255}
		}

		cont := widget.NewContainer(
			widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(c)))

		if i%3 == 0 {
			cont.GetWidget().LayoutData = &widget.GridLayoutData{
				MaxWidth: 30,
			}
		}

		container.AddChild(cont)
	}

	rootContainer.AddChild(container)

	rootContainer.AddChild(widget.NewText(
		widget.TextOpts.Text("test 2", fontFace, color.Black)))

	pageButtonsContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(10))))
	rootContainer.AddChild(pageButtonsContainer)

	flipBook := widget.NewFlipBook(
		widget.FlipBookOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(&widget.RowLayoutData{
			Stretch: true,
		}))))
	rootContainer.AddChild(flipBook)

	pages := []widget.HasWidget{}
	for i := 0; i < 5; i++ {
		c := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewFillLayout()))
		c.AddChild(widget.NewText(
			widget.TextOpts.Text(fmt.Sprintf("This is page %d", i+1), fontFace, color.Black)))
		pages = append(pages, c)
	}

	flipBook.SetPage(pages[0])

	for i := 0; i < 5; i++ {
		i := i
		pageButtonsContainer.AddChild(widget.NewButton(
			widget.ButtonOpts.Image(images.button),
			widget.ButtonOpts.Text(fmt.Sprintf("Page %d", i+1), fontFace, &widget.ButtonTextColor{
				Idle:     color.Black,
				Disabled: color.RGBA{128, 128, 128, 255},
			}),

			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				flipBook.SetPage(pages[i])
			})))
	}

	rootContainer.AddChild(widget.NewCheckbox(
		widget.CheckboxOpts.TriState(),
		widget.CheckboxOpts.ButtonOpts(widget.ButtonOpts.Image(images.button)),
		widget.CheckboxOpts.Image(images.checkbox)))

	return &ebitenui.UI{
		Container: rootContainer,
	}
}

func loadFont() (font.Face, error) {
	fontData, err := ioutil.ReadFile("fonts/JetBrainsMonoNL-Regular.ttf")
	if err != nil {
		return nil, err
	}

	ttfFont, err := truetype.Parse(fontData)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(ttfFont, &truetype.Options{
		Size:    20,
		DPI:     72,
		Hinting: font.HintingFull,
	}), nil
}

func (g *game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (g *game) Update(screen *ebiten.Image) error {
	g.ui.Update()
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	w, h := screen.Size()
	g.ui.Draw(screen, img.Rect(0, 0, w, h))

	// if firstRender {
	// 	stopper := profile.Start(profile.CPUProfile, profile.ProfilePath("."), profile.NoShutdownHook)
	// 	stopProfiler = stopper.Stop

	// 	firstRender = false
	// }
}
