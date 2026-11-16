package widget

import (
	"image"
	"image/color"

	"github.com/ebitenui/ebitenui/event"
	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// Dialog is a modal window asking the player to choose: a title, a message
// and a button for each choice under them, stretched to the dialog's width.
// One choice makes an alert, two a confirmation. Choosing fires ChoiceEvent
// and closes the dialog, unless KeepOpen was set. Add Window to a UI to show
// it, after CenterIn to place it.
type Dialog struct {
	ChoiceEvent *event.Event

	title, message string
	choices        []string
	maxWidth       float64
	titleFace      *text.Face
	titleColor     color.Color
	messageFace    *text.Face
	messageColor   color.Color
	button         *ButtonParams
	background     *eimage.NineSlice
	padding        *Insets
	spacing        int
	keepOpen       bool
	windowOpts     []WindowOpt

	window   *Window
	contents *themedPanel // the theme's panel background, without one of its own
}

// DialogChoiceEventArgs reports a choice: its index among the choices and
// its label.
type DialogChoiceEventArgs struct {
	Dialog *Dialog
	Choice int
	Label  string
}

type DialogChoiceHandlerFunc func(args *DialogChoiceEventArgs)

type DialogOpt func(d *Dialog)

type DialogOptions struct {
}

var DialogOpts DialogOptions

func NewDialog(opts ...DialogOpt) *Dialog {
	d := &Dialog{
		ChoiceEvent: &event.Event{},
		padding:     NewInsetsSimple(16),
		spacing:     10,
	}
	for _, o := range opts {
		o(d)
	}
	d.build()
	return d
}

// Title sets the dialog's title, drawn above the message.
func (o DialogOptions) Title(title string) DialogOpt {
	return func(d *Dialog) {
		d.title = title
	}
}

// Message sets the dialog's message.
func (o DialogOptions) Message(message string) DialogOpt {
	return func(d *Dialog) {
		d.message = message
	}
}

// Choices sets the labels of the dialog's buttons, from top to bottom.
func (o DialogOptions) Choices(labels ...string) DialogOpt {
	return func(d *Dialog) {
		d.choices = labels
	}
}

// MaxWidth wraps the title and the message at a width, in pixels.
func (o DialogOptions) MaxWidth(width float64) DialogOpt {
	return func(d *Dialog) {
		d.maxWidth = width
	}
}

// TitleText sets the title's face and color; without it, the theme's text
// face and color.
func (o DialogOptions) TitleText(face *text.Face, c color.Color) DialogOpt {
	return func(d *Dialog) {
		d.titleFace, d.titleColor = face, c
	}
}

// MessageText sets the message's face and color; without it, the theme's
// text face and color.
func (o DialogOptions) MessageText(face *text.Face, c color.Color) DialogOpt {
	return func(d *Dialog) {
		d.messageFace, d.messageColor = face, c
	}
}

// ButtonParams sets the choice buttons' images, text and padding; without
// it, the theme's button.
func (o DialogOptions) ButtonParams(params *ButtonParams) DialogOpt {
	return func(d *Dialog) {
		d.button = params
	}
}

// BackgroundImage sets the dialog's background; without one, the dialog
// draws its theme's panel background.
func (o DialogOptions) BackgroundImage(i *eimage.NineSlice) DialogOpt {
	return func(d *Dialog) {
		d.background = i
	}
}

// WindowOpts adds options to the dialog's window, such as a Transition.
func (o DialogOptions) WindowOpts(opts ...WindowOpt) DialogOpt {
	return func(d *Dialog) {
		d.windowOpts = append(d.windowOpts, opts...)
	}
}

// Padding sets the space between the dialog's edge and its contents; 16
// pixels by default.
func (o DialogOptions) Padding(p *Insets) DialogOpt {
	return func(d *Dialog) {
		d.padding = p
	}
}

// Spacing sets the space between the title, the message and the buttons;
// 10 pixels by default.
func (o DialogOptions) Spacing(s int) DialogOpt {
	return func(d *Dialog) {
		d.spacing = s
	}
}

// KeepOpen keeps the dialog open after a choice; it closes with
// Window().Close().
func (o DialogOptions) KeepOpen() DialogOpt {
	return func(d *Dialog) {
		d.keepOpen = true
	}
}

func (o DialogOptions) ChoiceHandler(f DialogChoiceHandlerFunc) DialogOpt {
	return func(d *Dialog) {
		d.ChoiceEvent.AddHandler(func(args any) {
			if a, ok := args.(*DialogChoiceEventArgs); ok {
				f(a)
			}
		})
	}
}

// Window returns the dialog's window, to add to a UI.
func (d *Dialog) Window() *Window {
	return d.window
}

// CenterIn places the dialog, at its preferred size, in the middle of an
// area such as the screen.
func (d *Dialog) CenterIn(area image.Rectangle) {
	w, h := d.contents.PreferredSize()
	at := area.Min.Add(image.Pt((area.Dx()-w)/2, (area.Dy()-h)/2))
	d.window.SetLocation(image.Rectangle{Min: at, Max: at.Add(image.Pt(w, h))})
}

func (d *Dialog) build() {
	containerOpts := []ContainerOpt{
		ContainerOpts.Layout(NewRowLayout(
			RowLayoutOpts.Direction(DirectionVertical),
			RowLayoutOpts.Padding(d.padding),
			RowLayoutOpts.Spacing(d.spacing),
		)),
	}
	if d.background != nil {
		containerOpts = append(containerOpts, ContainerOpts.BackgroundImage(d.background))
	}
	d.contents = &themedPanel{NewContainer(containerOpts...)}

	if d.title != "" {
		d.contents.AddChild(d.text(d.title, d.titleFace, d.titleColor))
	}
	if d.message != "" {
		d.contents.AddChild(d.text(d.message, d.messageFace, d.messageColor))
	}
	for i, label := range d.choices {
		d.contents.AddChild(d.choiceButton(i, label))
	}
	d.window = NewWindow(append([]WindowOpt{WindowOpts.Contents(d.contents), WindowOpts.Modal()}, d.windowOpts...)...)
}

func (d *Dialog) text(label string, face *text.Face, c color.Color) *Text {
	opts := []TextOpt{TextOpts.TextLabel(label), TextOpts.MaxWidth(d.maxWidth)}
	if face != nil {
		opts = append(opts, TextOpts.TextFace(face))
	}
	if c != nil {
		opts = append(opts, TextOpts.TextColor(c))
	}
	return NewText(opts...)
}

func (d *Dialog) choiceButton(i int, label string) *Button {
	opts := []ButtonOpt{
		ButtonOpts.WidgetOpts(WidgetOpts.LayoutData(RowLayoutData{Stretch: true})),
		ButtonOpts.TextLabel(label),
		ButtonOpts.ClickedHandler(func(*ButtonClickedEventArgs) {
			d.ChoiceEvent.Fire(&DialogChoiceEventArgs{Dialog: d, Choice: i, Label: label})
			if !d.keepOpen {
				d.window.Close()
			}
		}),
	}
	if p := d.button; p != nil {
		if p.Image != nil {
			opts = append(opts, ButtonOpts.Image(p.Image))
		}
		if p.TextFace != nil {
			opts = append(opts, ButtonOpts.TextFace(p.TextFace))
		}
		if p.TextColor != nil {
			opts = append(opts, ButtonOpts.TextColor(p.TextColor))
		}
		if p.TextPadding != nil {
			opts = append(opts, ButtonOpts.TextPadding(p.TextPadding))
		}
	}
	return NewButton(opts...)
}
