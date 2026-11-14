package widget

import (
	"image"
	"image/color"
	"math"

	"github.com/ebitenui/ebitenui/event"
	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/input"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// ColorPicker picks an opaque color: a wheel of every hue around, red to the
// right, and every saturation outward, white at the middle, and under it a
// slider of the color's brightness, from black to full. The color can be
// picked by pressing the wheel or dragging across it.
type ColorPicker struct {
	ChangedEvent *event.Event

	widgetOpts  []WidgetOpt
	wheelSize   int
	spacing     int
	handleImage *ButtonImage
	handleSize  int

	init      *MultiOnce
	container *Container
	wheel     *colorWheel
	slider    *Slider
	track     *ebiten.Image // the slider's track, black to the hue at full brightness

	h, s, v        float64 // the color: hue from 0 to 1 around from red, saturation and value
	sliderDragging bool
}

// ColorPickerChangedEventArgs reports a color picked, and whether the wheel
// or the slider is still being dragged.
type ColorPickerChangedEventArgs struct {
	ColorPicker *ColorPicker
	Color       color.NRGBA
	Dragging    bool
}

type ColorPickerChangedHandlerFunc func(args *ColorPickerChangedEventArgs)

type ColorPickerOpt func(p *ColorPicker)

type ColorPickerOptions struct {
}

var ColorPickerOpts ColorPickerOptions

// trackSteps is how many shades the brightness slider's track holds.
const trackSteps = 64

func NewColorPicker(opts ...ColorPickerOpt) *ColorPicker {
	p := &ColorPicker{
		ChangedEvent: &event.Event{},
		init:         &MultiOnce{},
		wheelSize:    160,
		spacing:      8,
		v:            1,
	}
	p.init.Append(p.createWidget)
	for _, o := range opts {
		o(p)
	}
	return p
}

func (o ColorPickerOptions) WidgetOpts(opts ...WidgetOpt) ColorPickerOpt {
	return func(p *ColorPicker) {
		p.widgetOpts = append(p.widgetOpts, opts...)
	}
}

// WheelSize sets the wheel's diameter, in pixels; 160 by default.
func (o ColorPickerOptions) WheelSize(size int) ColorPickerOpt {
	return func(p *ColorPicker) {
		p.wheelSize = size
	}
}

// Spacing sets the space between the wheel and the slider, in pixels; 8 by
// default.
func (o ColorPickerOptions) Spacing(spacing int) ColorPickerOpt {
	return func(p *ColorPicker) {
		p.spacing = spacing
	}
}

// InitialColor sets the color the picker starts with; white by default.
// Its alpha is ignored.
func (o ColorPickerOptions) InitialColor(c color.Color) ColorPickerOpt {
	return func(p *ColorPicker) {
		p.h, p.s, p.v = rgbToHSV(color.NRGBAModel.Convert(c).(color.NRGBA)) //nolint:forcetypeassert // NRGBAModel returns color.NRGBA
	}
}

// SliderHandle sets the image and the fixed size of the brightness slider's
// handle; without it, the slider takes them from the theme.
func (o ColorPickerOptions) SliderHandle(handle *ButtonImage, size int) ColorPickerOpt {
	return func(p *ColorPicker) {
		p.handleImage = handle
		p.handleSize = size
	}
}

func (o ColorPickerOptions) ChangedHandler(f ColorPickerChangedHandlerFunc) ColorPickerOpt {
	return func(p *ColorPicker) {
		p.ChangedEvent.AddHandler(func(args any) {
			if arg, ok := args.(*ColorPickerChangedEventArgs); ok {
				f(arg)
			}
		})
	}
}

// Color returns the color picked.
func (p *ColorPicker) Color() color.NRGBA {
	return hsvToRGB(p.h, p.s, p.v)
}

// SetColor sets the color picked, without firing ChangedEvent. Its alpha is
// ignored.
func (p *ColorPicker) SetColor(c color.Color) {
	p.init.Do()
	p.h, p.s, p.v = rgbToHSV(color.NRGBAModel.Convert(c).(color.NRGBA)) //nolint:forcetypeassert // NRGBAModel returns color.NRGBA
	p.updateTrack()
	p.slider.Current = p.sliderValue()
}

func (p *ColorPicker) GetWidget() *Widget {
	p.init.Do()
	return p.container.GetWidget()
}

func (p *ColorPicker) PreferredSize() (int, int) {
	p.init.Do()
	return p.container.PreferredSize()
}

func (p *ColorPicker) SetLocation(rect image.Rectangle) {
	p.init.Do()
	p.container.SetLocation(rect)
}

func (p *ColorPicker) RequestRelayout() {
	p.init.Do()
	p.container.RequestRelayout()
}

func (p *ColorPicker) SetupInputLayer(def input.DeferredSetupInputLayerFunc) {
	p.init.Do()
	p.container.SetupInputLayer(def)
}

func (p *ColorPicker) Validate() {
	p.init.Do()
	p.container.Validate()
}

func (p *ColorPicker) Render(screen *ebiten.Image) {
	p.init.Do()
	p.container.Render(screen)
}

func (p *ColorPicker) Update(updObj *UpdateObject) {
	p.init.Do()
	p.container.Update(updObj)
}

func (p *ColorPicker) createWidget() {
	p.container = NewContainer(
		ContainerOpts.WidgetOpts(p.widgetOpts...),
		ContainerOpts.Layout(NewRowLayout(RowLayoutOpts.Direction(DirectionVertical), RowLayoutOpts.Spacing(p.spacing))),
	)
	p.widgetOpts = nil

	p.wheel = newColorWheel(p)
	p.container.AddChild(p.wheel)

	p.track = ebiten.NewImage(trackSteps, 1)
	p.updateTrack()
	track := eimage.NewNineSlice(p.track, [3]int{0, trackSteps, 0}, [3]int{0, 1, 0})
	sliderOpts := []SliderOpt{
		SliderOpts.Orientation(DirectionHorizontal),
		SliderOpts.MinMax(0, 255),
		SliderOpts.InitialCurrent(p.sliderValue()),
		SliderOpts.TrackImage(&SliderTrackImage{Idle: track, Hover: track}),
		SliderOpts.WidgetOpts(WidgetOpts.LayoutData(RowLayoutData{Stretch: true})),
		SliderOpts.ChangedHandler(p.sliderChanged),
	}
	if p.handleImage != nil {
		sliderOpts = append(sliderOpts, SliderOpts.HandleImage(p.handleImage), SliderOpts.FixedHandleSize(p.handleSize))
	}
	p.slider = NewSlider(sliderOpts...)
	p.container.AddChild(p.slider)
}

// sliderValue is the slider's place for the color's brightness.
func (p *ColorPicker) sliderValue() int {
	return int(math.Round(p.v * 255))
}

// sliderChanged takes a brightness from the slider. The slider also reports
// values set by SetColor, which change nothing.
func (p *ColorPicker) sliderChanged(args *SliderChangedEventArgs) {
	v := float64(args.Current) / 255
	ended := p.sliderDragging && !args.Dragging
	p.sliderDragging = args.Dragging
	if v == p.v && !ended {
		return
	}
	p.v = v
	p.fireChanged(args.Dragging)
}

// pick takes a hue and saturation from the wheel, keeping the brightness,
// or full brightness from black, where no hue shows.
func (p *ColorPicker) pick(h, s float64, dragging bool) {
	p.h, p.s = h, s
	if p.v == 0 {
		p.v = 1
		p.slider.Current = p.sliderValue()
	}
	p.updateTrack()
	p.fireChanged(dragging)
}

func (p *ColorPicker) fireChanged(dragging bool) {
	p.ChangedEvent.Fire(&ColorPickerChangedEventArgs{ColorPicker: p, Color: p.Color(), Dragging: dragging})
}

// updateTrack shades the slider's track from black to the color's hue and
// saturation at full brightness.
func (p *ColorPicker) updateTrack() {
	if p.track == nil {
		return
	}
	pix := make([]byte, 4*trackSteps)
	for x := range trackSteps {
		c := hsvToRGB(p.h, p.s, float64(x)/(trackSteps-1))
		pix[4*x], pix[4*x+1], pix[4*x+2], pix[4*x+3] = c.R, c.G, c.B, 0xff
	}
	p.track.WritePixels(pix)
}

// colorWheel is the picker's disc of hues and saturations, drawn once, with
// a ring where the color lies.
type colorWheel struct {
	widget   *Widget
	picker   *ColorPicker
	disc     *ebiten.Image
	dragging bool
	op       ebiten.DrawImageOptions
}

func newColorWheel(p *ColorPicker) *colorWheel {
	w := &colorWheel{widget: NewWidget(), picker: p, disc: wheelDisc(p.wheelSize)}
	w.widget.MouseButtonPressedEvent.AddHandler(func(args any) {
		a, ok := args.(*WidgetMouseButtonPressedEventArgs)
		if !ok || a.Button != ebiten.MouseButtonLeft || w.widget.Disabled {
			return
		}
		w.dragging = true
		w.pickAt(a.OffsetX, a.OffsetY, true)
	})
	return w
}

func (w *colorWheel) GetWidget() *Widget { return w.widget }

func (w *colorWheel) PreferredSize() (int, int) { return w.picker.wheelSize, w.picker.wheelSize }

func (w *colorWheel) SetLocation(rect image.Rectangle) { w.widget.Rect = rect }

func (w *colorWheel) Validate() {}

func (w *colorWheel) Update(updObj *UpdateObject) {
	w.widget.Update(updObj)
	if !w.dragging {
		return
	}
	x, y := input.CursorPosition()
	at := image.Pt(x, y).Sub(w.widget.Rect.Min)
	if !input.MouseButtonPressed(ebiten.MouseButtonLeft) {
		w.dragging = false
		w.pickAt(at.X, at.Y, false)
		return
	}
	w.pickAt(at.X, at.Y, true)
}

func (w *colorWheel) Render(screen *ebiten.Image) {
	w.widget.Render(screen)
	r := w.widget.Rect
	w.op.GeoM.Reset()
	w.op.GeoM.Translate(float64(r.Min.X), float64(r.Min.Y))
	w.op.ColorScale.Reset()
	if w.widget.Disabled {
		w.op.ColorScale.ScaleAlpha(0.35)
	}
	screen.DrawImage(w.disc, &w.op)

	radius := float64(w.picker.wheelSize) / 2
	h, s := w.picker.h, w.picker.s
	x := float64(r.Min.X) + radius + s*radius*math.Cos(2*math.Pi*h)
	y := float64(r.Min.Y) + radius - s*radius*math.Sin(2*math.Pi*h)
	ring := float32(max(4, radius/12))
	vector.StrokeCircle(screen, float32(x), float32(y), ring, 3, color.NRGBA{A: 0xff}, true)
	vector.StrokeCircle(screen, float32(x), float32(y), ring, 1.5, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}, true)
}

// pickAt picks the hue and saturation under a point, from the wheel's top
// left corner; points past the rim take the rim's.
func (w *colorWheel) pickAt(x, y int, dragging bool) {
	r := float64(w.picker.wheelSize) / 2
	dx, dy := float64(x)+0.5-r, r-(float64(y)+0.5)
	h := math.Atan2(dy, dx) / (2 * math.Pi)
	w.picker.pick(h-math.Floor(h), min(1, math.Hypot(dx, dy)/r), dragging)
}

// wheelDisc draws a disc size pixels wide of every hue around and every
// saturation outward at full brightness, its rim smoothed over a pixel.
func wheelDisc(size int) *ebiten.Image {
	pix := make([]byte, 4*size*size)
	r := float64(size) / 2
	for y := range size {
		for x := range size {
			dx, dy := float64(x)+0.5-r, r-(float64(y)+0.5)
			d := math.Hypot(dx, dy) / r
			if d > 1 {
				continue
			}
			c := hsvToRGB(math.Atan2(dy, dx)/(2*math.Pi), d, 1)
			edge := min(1, (1-d)*r)
			i := 4 * (y*size + x)
			pix[i], pix[i+1], pix[i+2], pix[i+3] = uint8(float64(c.R)*edge), uint8(float64(c.G)*edge), uint8(float64(c.B)*edge), uint8(255*edge)
		}
	}
	disc := ebiten.NewImage(size, size)
	disc.WritePixels(pix)
	return disc
}

// hsvToRGB returns the opaque color of a hue, from 0 to 1 around from red,
// a saturation and a value, both from 0 to 1.
func hsvToRGB(h, s, v float64) color.NRGBA {
	h = (h - math.Floor(h)) * 6
	sector := math.Floor(h)
	f := h - sector
	p, q, t := v*(1-s), v*(1-s*f), v*(1-s*(1-f))
	var r, g, b float64
	switch int(sector) % 6 {
	case 0:
		r, g, b = v, t, p
	case 1:
		r, g, b = q, v, p
	case 2:
		r, g, b = p, v, t
	case 3:
		r, g, b = p, q, v
	case 4:
		r, g, b = t, p, v
	default:
		r, g, b = v, p, q
	}
	channel := func(x float64) uint8 { return uint8(math.Round(x * 255)) }
	return color.NRGBA{R: channel(r), G: channel(g), B: channel(b), A: 0xff}
}

// rgbToHSV returns a color's hue, from 0 to 1 around from red, its
// saturation and its value; gray has hue 0.
func rgbToHSV(c color.NRGBA) (h, s, v float64) {
	r, g, b := float64(c.R)/255, float64(c.G)/255, float64(c.B)/255
	hi, lo := max(r, g, b), min(r, g, b)
	v = hi
	if hi == 0 {
		return 0, 0, 0
	}
	s = (hi - lo) / hi
	if hi == lo {
		return 0, s, v
	}
	switch hi {
	case r:
		h = (g - b) / (hi - lo)
	case g:
		h = 2 + (b-r)/(hi-lo)
	default:
		h = 4 + (r-g)/(hi-lo)
	}
	h /= 6
	return h - math.Floor(h), s, v
}
