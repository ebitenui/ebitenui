package main

import (
	"image/color"
	"strconv"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	backgroundColor = "131a22"

	textIdleColor     = "dff4ff"
	textDisabledColor = "5a7a91"

	labelIdleColor     = textIdleColor
	labelDisabledColor = textDisabledColor

	buttonIdleColor     = textIdleColor
	buttonDisabledColor = labelDisabledColor

	listSelectedBackground         = "4b687a"
	listDisabledSelectedBackground = "2a3944"

	listFocusedBackground = "2a3944"

	headerColor = textIdleColor

	textInputCaretColor         = "e7c34b"
	textInputDisabledCaretColor = "766326"

	toolTipColor = backgroundColor

	separatorColor = listDisabledSelectedBackground
)

type uiResources struct {
	fonts *fonts

	background *image.NineSlice

	separatorColor color.Color

	text        *textResources
	button      *buttonResources
	label       *labelResources
	checkbox    *checkboxResources
	comboButton *comboButtonResources
	list        *listResources
	slider      *sliderResources
	progressBar *progressBarResources
	panel       *panelResources
	tabBook     *tabBookResources
	header      *headerResources
	textInput   *textInputResources
	textArea    *textAreaResources
	toolTip     *toolTipResources
}

type textResources struct {
	idleColor     color.Color
	disabledColor color.Color
	face          text.Face
	titleFace     text.Face
	bigTitleFace  text.Face
	smallFace     text.Face
}

type buttonResources struct {
	image   *widget.ButtonImage
	text    *widget.ButtonTextColor
	face    text.Face
	padding widget.Insets
}

type checkboxResources struct {
	image   *widget.CheckboxImage
	spacing int
}

type labelResources struct {
	text *widget.LabelColor
	face text.Face
}

type comboButtonResources struct {
	image   *widget.ButtonImage
	text    *widget.ButtonTextColor
	face    text.Face
	graphic *widget.GraphicImage
	padding widget.Insets
}

type listResources struct {
	image        *widget.ScrollContainerImage
	track        *widget.SliderTrackImage
	trackPadding widget.Insets
	handle       *widget.ButtonImage
	handleSize   int
	face         text.Face
	entry        *widget.ListEntryColor
	entryPadding widget.Insets
}

type sliderResources struct {
	trackImage *widget.SliderTrackImage
	handle     *widget.ButtonImage
	handleSize int
}

type progressBarResources struct {
	trackImage *widget.ProgressBarImage
	fillImage  *widget.ProgressBarImage
}

type panelResources struct {
	image    *image.NineSlice
	titleBar *image.NineSlice
	padding  widget.Insets
}

type tabBookResources struct {
	buttonFace    text.Face
	buttonText    *widget.ButtonTextColor
	buttonPadding widget.Insets
}

type headerResources struct {
	background *image.NineSlice
	padding    widget.Insets
	face       text.Face
	color      color.Color
}

type textInputResources struct {
	image   *widget.TextInputImage
	padding widget.Insets
	face    text.Face
	color   *widget.TextInputColor
}

type textAreaResources struct {
	image        *widget.ScrollContainerImage
	track        *widget.SliderTrackImage
	trackPadding widget.Insets
	handle       *widget.ButtonImage
	handleSize   int
	face         text.Face
	entryPadding widget.Insets
}

type toolTipResources struct {
	background *image.NineSlice
	padding    widget.Insets
	face       text.Face
	color      color.Color
}

func newUIResources() (*uiResources, error) {
	background := image.NewNineSliceColor(hexToColor(backgroundColor))

	fonts, err := loadFonts()
	if err != nil {
		return nil, err
	}

	button, err := newButtonResources(fonts)
	if err != nil {
		return nil, err
	}

	checkbox, err := newCheckboxResources()
	if err != nil {
		return nil, err
	}

	comboButton, err := newComboButtonResources(fonts)
	if err != nil {
		return nil, err
	}

	list, err := newListResources(fonts)
	if err != nil {
		return nil, err
	}

	slider, err := newSliderResources()
	if err != nil {
		return nil, err
	}

	progressBar, err := newProgressBarResources()
	if err != nil {
		return nil, err
	}

	panel, err := newPanelResources()
	if err != nil {
		return nil, err
	}

	tabBook, err := newTabBookResources(fonts)
	if err != nil {
		return nil, err
	}

	header, err := newHeaderResources(fonts)
	if err != nil {
		return nil, err
	}

	textInput, err := newTextInputResources(fonts)
	if err != nil {
		return nil, err
	}
	textArea, err := newTextAreaResources(fonts)
	if err != nil {
		return nil, err
	}
	toolTip, err := newToolTipResources(fonts)
	if err != nil {
		return nil, err
	}

	return &uiResources{
		fonts: fonts,

		background: background,

		separatorColor: hexToColor(separatorColor),

		text: &textResources{
			idleColor:     hexToColor(textIdleColor),
			disabledColor: hexToColor(textDisabledColor),
			face:          fonts.face,
			titleFace:     fonts.titleFace,
			bigTitleFace:  fonts.bigTitleFace,
			smallFace:     fonts.toolTipFace,
		},

		button:      button,
		label:       newLabelResources(fonts),
		checkbox:    checkbox,
		comboButton: comboButton,
		list:        list,
		slider:      slider,
		panel:       panel,
		tabBook:     tabBook,
		header:      header,
		textInput:   textInput,
		toolTip:     toolTip,
		textArea:    textArea,
		progressBar: progressBar,
	}, nil
}

func newButtonResources(fonts *fonts) (*buttonResources, error) {
	idle, err := loadImageNineSlice("assets/graphics/button-idle.png", 12, 0)
	if err != nil {
		return nil, err
	}

	hover, err := loadImageNineSlice("assets/graphics/button-hover.png", 12, 0)
	if err != nil {
		return nil, err
	}
	pressed_hover, err := loadImageNineSlice("assets/graphics/button-selected-hover.png", 12, 0)
	if err != nil {
		return nil, err
	}
	pressed, err := loadImageNineSlice("assets/graphics/button-pressed.png", 12, 0)
	if err != nil {
		return nil, err
	}

	disabled, err := loadImageNineSlice("assets/graphics/button-disabled.png", 12, 0)
	if err != nil {
		return nil, err
	}

	i := &widget.ButtonImage{
		Idle:         idle,
		Hover:        hover,
		Pressed:      pressed,
		PressedHover: pressed_hover,
		Disabled:     disabled,
	}

	return &buttonResources{
		image: i,

		text: &widget.ButtonTextColor{
			Idle:     hexToColor(buttonIdleColor),
			Disabled: hexToColor(buttonDisabledColor),
		},

		face: fonts.face,

		padding: widget.Insets{
			Left:  30,
			Right: 30,
		},
	}, nil
}

func newCheckboxResources() (*checkboxResources, error) {
	f1, err := embeddedAssets.Open("assets/graphics/checkbox-idle.png")
	if err != nil {
		return nil, err
	}
	defer f1.Close()
	idle, _, _ := ebitenutil.NewImageFromReader(f1)

	f2, err := embeddedAssets.Open("assets/graphics/checkbox-checked.png")
	if err != nil {
		return nil, err
	}
	defer f2.Close()
	checked, _, _ := ebitenutil.NewImageFromReader(f2)

	f3, err := embeddedAssets.Open("assets/graphics/checkbox-greyed.png")
	if err != nil {
		return nil, err
	}
	defer f3.Close()
	greyed, _, _ := ebitenutil.NewImageFromReader(f3)

	f4, err := embeddedAssets.Open("assets/graphics/checkbox-hover.png")
	if err != nil {
		return nil, err
	}
	defer f4.Close()
	idle_hovered, _, _ := ebitenutil.NewImageFromReader(f4)

	f5, err := embeddedAssets.Open("assets/graphics/checkbox-checked-hover.png")
	if err != nil {
		return nil, err
	}
	defer f5.Close()
	checked_hovered, _, _ := ebitenutil.NewImageFromReader(f5)

	f6, err := embeddedAssets.Open("assets/graphics/checkbox-greyed-hover.png")
	if err != nil {
		return nil, err
	}
	defer f6.Close()
	greyed_hovered, _, _ := ebitenutil.NewImageFromReader(f6)

	f7, err := embeddedAssets.Open("assets/graphics/checkbox-disabled.png")
	if err != nil {
		return nil, err
	}
	defer f7.Close()
	idle_disabled, _, _ := ebitenutil.NewImageFromReader(f7)

	f8, err := embeddedAssets.Open("assets/graphics/checkbox-checked-disabled.png")
	if err != nil {
		return nil, err
	}
	defer f8.Close()
	checked_disabled, _, _ := ebitenutil.NewImageFromReader(f8)

	f9, err := embeddedAssets.Open("assets/graphics/checkbox-greyed-disabled.png")
	if err != nil {
		return nil, err
	}
	defer f9.Close()
	greyed_disabled, _, _ := ebitenutil.NewImageFromReader(f9)

	return &checkboxResources{
		image: &widget.CheckboxImage{
			Unchecked:         image.NewFixedNineSlice(idle),
			Checked:           image.NewFixedNineSlice(checked),
			Greyed:            image.NewFixedNineSlice(greyed),
			UncheckedHovered:  image.NewFixedNineSlice(idle_hovered),
			CheckedHovered:    image.NewFixedNineSlice(checked_hovered),
			GreyedHovered:     image.NewFixedNineSlice(greyed_hovered),
			UncheckedDisabled: image.NewFixedNineSlice(idle_disabled),
			CheckedDisabled:   image.NewFixedNineSlice(checked_disabled),
			GreyedDisabled:    image.NewFixedNineSlice(greyed_disabled),
		},

		spacing: 10,
	}, nil
}

func newLabelResources(fonts *fonts) *labelResources {
	return &labelResources{
		text: &widget.LabelColor{
			Idle:     hexToColor(labelIdleColor),
			Disabled: hexToColor(labelDisabledColor),
		},

		face: fonts.face,
	}
}

func newComboButtonResources(fonts *fonts) (*comboButtonResources, error) {
	idle, err := loadImageNineSlice("assets/graphics/combo-button-idle.png", 12, 0)
	if err != nil {
		return nil, err
	}

	hover, err := loadImageNineSlice("assets/graphics/combo-button-hover.png", 12, 0)
	if err != nil {
		return nil, err
	}

	pressed, err := loadImageNineSlice("assets/graphics/combo-button-pressed.png", 12, 0)
	if err != nil {
		return nil, err
	}

	disabled, err := loadImageNineSlice("assets/graphics/combo-button-disabled.png", 12, 0)
	if err != nil {
		return nil, err
	}

	i := &widget.ButtonImage{
		Idle:     idle,
		Hover:    hover,
		Pressed:  pressed,
		Disabled: disabled,
	}

	arrowDown, err := loadGraphicImages("assets/graphics/arrow-down-idle.png", "assets/graphics/arrow-down-disabled.png")
	if err != nil {
		return nil, err
	}

	return &comboButtonResources{
		image: i,

		text: &widget.ButtonTextColor{
			Idle:     hexToColor(buttonIdleColor),
			Disabled: hexToColor(buttonDisabledColor),
		},

		face:    fonts.face,
		graphic: arrowDown,

		padding: widget.Insets{
			Left:  30,
			Right: 30,
		},
	}, nil
}

func newListResources(fonts *fonts) (*listResources, error) {
	idle, err := newImageFromFile("assets/graphics/list-idle.png")
	if err != nil {
		return nil, err
	}

	disabled, err := newImageFromFile("assets/graphics/list-disabled.png")
	if err != nil {
		return nil, err
	}

	mask, err := newImageFromFile("assets/graphics/list-mask.png")
	if err != nil {
		return nil, err
	}

	trackIdle, err := newImageFromFile("assets/graphics/list-track-idle.png")
	if err != nil {
		return nil, err
	}

	trackDisabled, err := newImageFromFile("assets/graphics/list-track-disabled.png")
	if err != nil {
		return nil, err
	}

	handleIdle, err := newImageFromFile("assets/graphics/slider-handle-idle.png")
	if err != nil {
		return nil, err
	}

	handleHover, err := newImageFromFile("assets/graphics/slider-handle-hover.png")
	if err != nil {
		return nil, err
	}

	return &listResources{
		image: &widget.ScrollContainerImage{
			Idle:     image.NewNineSlice(idle, [3]int{25, 12, 22}, [3]int{25, 12, 25}),
			Disabled: image.NewNineSlice(disabled, [3]int{25, 12, 22}, [3]int{25, 12, 25}),
			Mask:     image.NewNineSlice(mask, [3]int{26, 10, 23}, [3]int{26, 10, 26}),
		},

		track: &widget.SliderTrackImage{
			Idle:     image.NewNineSlice(trackIdle, [3]int{5, 0, 0}, [3]int{25, 12, 25}),
			Hover:    image.NewNineSlice(trackIdle, [3]int{5, 0, 0}, [3]int{25, 12, 25}),
			Disabled: image.NewNineSlice(trackDisabled, [3]int{0, 5, 0}, [3]int{25, 12, 25}),
		},

		trackPadding: widget.Insets{
			Top:    5,
			Bottom: 24,
		},

		handle: &widget.ButtonImage{
			Idle:     image.NewNineSliceSimple(handleIdle, 0, 5),
			Hover:    image.NewNineSliceSimple(handleHover, 0, 5),
			Pressed:  image.NewNineSliceSimple(handleHover, 0, 5),
			Disabled: image.NewNineSliceSimple(handleIdle, 0, 5),
		},

		handleSize: 5,
		face:       fonts.face,

		entry: &widget.ListEntryColor{
			Unselected:         hexToColor(textIdleColor),
			DisabledUnselected: hexToColor(textDisabledColor),

			Selected:         hexToColor(textIdleColor),
			DisabledSelected: hexToColor(textDisabledColor),

			SelectedBackground:         hexToColor(listSelectedBackground),
			DisabledSelectedBackground: hexToColor(listDisabledSelectedBackground),

			FocusedBackground:         hexToColor(listFocusedBackground),
			SelectedFocusedBackground: hexToColor(listSelectedBackground),
		},

		entryPadding: widget.Insets{
			Left:   30,
			Right:  30,
			Top:    2,
			Bottom: 2,
		},
	}, nil
}

func newSliderResources() (*sliderResources, error) {
	idle, err := newImageFromFile("assets/graphics/slider-track-idle.png")
	if err != nil {
		return nil, err
	}

	disabled, err := newImageFromFile("assets/graphics/slider-track-disabled.png")
	if err != nil {
		return nil, err
	}

	handleIdle, err := newImageFromFile("assets/graphics/slider-handle-idle.png")
	if err != nil {
		return nil, err
	}

	handleHover, err := newImageFromFile("assets/graphics/slider-handle-hover.png")
	if err != nil {
		return nil, err
	}

	handleDisabled, err := newImageFromFile("assets/graphics/slider-handle-disabled.png")
	if err != nil {
		return nil, err
	}

	return &sliderResources{
		trackImage: &widget.SliderTrackImage{
			Idle:     image.NewNineSlice(idle, [3]int{0, 19, 0}, [3]int{6, 0, 0}),
			Hover:    image.NewNineSlice(idle, [3]int{0, 19, 0}, [3]int{6, 0, 0}),
			Disabled: image.NewNineSlice(disabled, [3]int{0, 19, 0}, [3]int{6, 0, 0}),
		},

		handle: &widget.ButtonImage{
			Idle:     image.NewNineSliceSimple(handleIdle, 0, 5),
			Hover:    image.NewNineSliceSimple(handleHover, 0, 5),
			Pressed:  image.NewNineSliceSimple(handleHover, 0, 5),
			Disabled: image.NewNineSliceSimple(handleDisabled, 0, 5),
		},

		handleSize: 6,
	}, nil
}

func newProgressBarResources() (*progressBarResources, error) {
	idle, err := newImageFromFile("assets/graphics/progressbar-track-idle.png")
	if err != nil {
		return nil, err
	}
	fill_idle, err := newImageFromFile("assets/graphics/progressbar-fill-idle.png")
	if err != nil {
		return nil, err
	}
	disabled, err := newImageFromFile("assets/graphics/slider-track-disabled.png")
	if err != nil {
		return nil, err
	}

	return &progressBarResources{
		trackImage: &widget.ProgressBarImage{
			Idle:     image.NewNineSlice(idle, [3]int{4, 11, 4}, [3]int{2, 2, 2}),
			Hover:    image.NewNineSlice(idle, [3]int{4, 11, 4}, [3]int{2, 2, 2}),
			Disabled: image.NewNineSlice(disabled, [3]int{4, 11, 4}, [3]int{2, 2, 2}),
		},

		fillImage: &widget.ProgressBarImage{
			Idle:     image.NewNineSlice(fill_idle, [3]int{4, 11, 4}, [3]int{2, 2, 2}),
			Hover:    image.NewNineSlice(fill_idle, [3]int{4, 11, 4}, [3]int{2, 2, 2}),
			Disabled: image.NewNineSlice(fill_idle, [3]int{4, 11, 4}, [3]int{2, 2, 2}),
		},
	}, nil
}
func newPanelResources() (*panelResources, error) {
	i, err := loadImageNineSlice("assets/graphics/panel-idle.png", 10, 10)
	if err != nil {
		return nil, err
	}
	t, err := loadImageNineSlice("assets/graphics/titlebar-idle.png", 10, 10)
	if err != nil {
		return nil, err
	}
	return &panelResources{
		image:    i,
		titleBar: t,
		padding: widget.Insets{
			Left:   30,
			Right:  30,
			Top:    20,
			Bottom: 20,
		},
	}, nil
}

func newTabBookResources(fonts *fonts) (*tabBookResources, error) {

	return &tabBookResources{
		buttonFace: fonts.face,

		buttonText: &widget.ButtonTextColor{
			Idle:     hexToColor(buttonIdleColor),
			Disabled: hexToColor(buttonDisabledColor),
		},

		buttonPadding: widget.Insets{
			Left:  30,
			Right: 30,
		},
	}, nil
}

func newHeaderResources(fonts *fonts) (*headerResources, error) {
	bg, err := loadImageNineSlice("assets/graphics/header.png", 446, 9)
	if err != nil {
		return nil, err
	}

	return &headerResources{
		background: bg,

		padding: widget.Insets{
			Left:   25,
			Right:  25,
			Top:    4,
			Bottom: 4,
		},

		face:  fonts.bigTitleFace,
		color: hexToColor(headerColor),
	}, nil
}

func newTextInputResources(fonts *fonts) (*textInputResources, error) {
	idle, err := newImageFromFile("assets/graphics/text-input-idle.png")
	if err != nil {
		return nil, err
	}

	disabled, err := newImageFromFile("assets/graphics/text-input-disabled.png")
	if err != nil {
		return nil, err
	}

	return &textInputResources{
		image: &widget.TextInputImage{
			Idle:     image.NewNineSlice(idle, [3]int{9, 14, 6}, [3]int{9, 14, 6}),
			Disabled: image.NewNineSlice(disabled, [3]int{9, 14, 6}, [3]int{9, 14, 6}),
		},

		padding: widget.Insets{
			Left:   8,
			Right:  8,
			Top:    4,
			Bottom: 4,
		},

		face: fonts.face,

		color: &widget.TextInputColor{
			Idle:          hexToColor(textIdleColor),
			Disabled:      hexToColor(textDisabledColor),
			Caret:         hexToColor(textInputCaretColor),
			DisabledCaret: hexToColor(textInputDisabledCaretColor),
		},
	}, nil
}

func newTextAreaResources(fonts *fonts) (*textAreaResources, error) {
	idle, err := newImageFromFile("assets/graphics/list-idle.png")
	if err != nil {
		return nil, err
	}

	disabled, err := newImageFromFile("assets/graphics/list-disabled.png")
	if err != nil {
		return nil, err
	}

	mask, err := newImageFromFile("assets/graphics/list-mask.png")
	if err != nil {
		return nil, err
	}

	trackIdle, err := newImageFromFile("assets/graphics/list-track-idle.png")
	if err != nil {
		return nil, err
	}

	trackDisabled, err := newImageFromFile("assets/graphics/list-track-disabled.png")
	if err != nil {
		return nil, err
	}

	handleIdle, err := newImageFromFile("assets/graphics/slider-handle-idle.png")
	if err != nil {
		return nil, err
	}

	handleHover, err := newImageFromFile("assets/graphics/slider-handle-hover.png")
	if err != nil {
		return nil, err
	}

	return &textAreaResources{
		image: &widget.ScrollContainerImage{
			Idle:     image.NewNineSlice(idle, [3]int{25, 12, 22}, [3]int{25, 12, 25}),
			Disabled: image.NewNineSlice(disabled, [3]int{25, 12, 22}, [3]int{25, 12, 25}),
			Mask:     image.NewNineSlice(mask, [3]int{26, 10, 23}, [3]int{26, 10, 26}),
		},

		track: &widget.SliderTrackImage{
			Idle:     image.NewNineSlice(trackIdle, [3]int{5, 0, 0}, [3]int{25, 12, 25}),
			Hover:    image.NewNineSlice(trackIdle, [3]int{5, 0, 0}, [3]int{25, 12, 25}),
			Disabled: image.NewNineSlice(trackDisabled, [3]int{0, 5, 0}, [3]int{25, 12, 25}),
		},

		trackPadding: widget.Insets{
			Top:    5,
			Bottom: 24,
		},

		handle: &widget.ButtonImage{
			Idle:     image.NewNineSliceSimple(handleIdle, 0, 5),
			Hover:    image.NewNineSliceSimple(handleHover, 0, 5),
			Pressed:  image.NewNineSliceSimple(handleHover, 0, 5),
			Disabled: image.NewNineSliceSimple(handleIdle, 0, 5),
		},

		handleSize: 5,
		face:       fonts.face,

		entryPadding: widget.Insets{
			Left:   30,
			Right:  30,
			Top:    2,
			Bottom: 2,
		},
	}, nil
}

func newToolTipResources(fonts *fonts) (*toolTipResources, error) {
	bg, err := newImageFromFile("assets/graphics/tool-tip.png")
	if err != nil {
		return nil, err
	}

	return &toolTipResources{
		background: image.NewNineSlice(bg, [3]int{19, 6, 13}, [3]int{19, 5, 13}),

		padding: widget.Insets{
			Left:   15,
			Right:  15,
			Top:    10,
			Bottom: 10,
		},

		face:  fonts.toolTipFace,
		color: hexToColor(toolTipColor),
	}, nil
}

func hexToColor(h string) color.Color {
	u, err := strconv.ParseUint(h, 16, 0)
	if err != nil {
		panic(err)
	}

	return color.NRGBA{
		R: uint8(u & 0xff0000 >> 16),
		G: uint8(u & 0xff00 >> 8),
		B: uint8(u & 0xff),
		A: 255,
	}
}
