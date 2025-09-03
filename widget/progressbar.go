package widget

import (
	img "image"

	"github.com/ebitenui/ebitenui/image"

	"github.com/hajimehoshi/ebiten/v2"
)

type ProgressBarParams struct {
	TrackImage   *ProgressBarImage
	FillImage    *ProgressBarImage
	TrackPadding *Insets
}

type ProgressBar struct {
	definedParams  ProgressBarParams
	computedParams ProgressBarParams

	Min     int
	Max     int
	current int

	widgetOpts []WidgetOpt
	direction  Direction
	inverted   bool

	init   *MultiOnce
	widget *Widget
}

type ProgressBarImage struct {
	Idle     *image.NineSlice
	Hover    *image.NineSlice
	Disabled *image.NineSlice
}

type ProgressBarOpt func(s *ProgressBar)

type ProgressBarOptions struct {
}

var ProgressBarOpts ProgressBarOptions

func NewProgressBar(opts ...ProgressBarOpt) *ProgressBar {
	pb := &ProgressBar{
		Min:     1,
		Max:     100,
		current: 1,

		init: &MultiOnce{},
	}

	pb.init.Append(pb.createWidget)

	for _, o := range opts {
		o(pb)
	}

	return pb
}

func (pb *ProgressBar) Validate() {
	pb.init.Do()
	pb.populateComputedParams()

	if pb.computedParams.TrackImage == nil {
		panic("ProgressBar: TrackImage is required.")
	}
	if pb.computedParams.TrackImage.Idle == nil {
		panic("ProgressBar: TrackImage.Idle is required")
	}

}

func (pb *ProgressBar) populateComputedParams() {
	params := ProgressBarParams{}
	theme := pb.widget.GetTheme()
	if theme != nil {
		if theme.ProgressBarTheme != nil {
			params.FillImage = theme.ProgressBarTheme.FillImage
			params.TrackImage = theme.ProgressBarTheme.TrackImage
			params.TrackPadding = theme.ProgressBarTheme.TrackPadding
		}
	}
	if pb.definedParams.FillImage != nil {
		if params.FillImage == nil {
			params.FillImage = pb.definedParams.FillImage
		} else {
			if pb.definedParams.FillImage.Idle != nil {
				params.FillImage.Idle = pb.definedParams.FillImage.Idle
			}
			if pb.definedParams.FillImage.Hover != nil {
				params.FillImage.Hover = pb.definedParams.FillImage.Hover
			}
			if pb.definedParams.FillImage.Disabled != nil {
				params.FillImage.Disabled = pb.definedParams.FillImage.Disabled
			}
		}
	}
	if pb.definedParams.TrackImage != nil {
		if params.TrackImage == nil {
			params.TrackImage = pb.definedParams.TrackImage
		} else {

			if pb.definedParams.TrackImage.Idle != nil {
				params.TrackImage.Idle = pb.definedParams.TrackImage.Idle
			}
			if pb.definedParams.TrackImage.Hover != nil {
				params.TrackImage.Hover = pb.definedParams.TrackImage.Hover
			}
			if pb.definedParams.TrackImage.Disabled != nil {
				params.TrackImage.Disabled = pb.definedParams.TrackImage.Disabled
			}
		}
	}

	if pb.definedParams.TrackPadding != nil {
		params.TrackPadding = pb.definedParams.TrackPadding
	}

	if params.TrackPadding == nil {
		params.TrackPadding = &Insets{}
	}

	pb.computedParams = params
}

func (o ProgressBarOptions) WidgetOpts(opts ...WidgetOpt) ProgressBarOpt {
	return func(s *ProgressBar) {
		s.widgetOpts = append(s.widgetOpts, opts...)
	}
}

// Direction sets the direction of the progress bar.
// The default is horizontal.
func (o ProgressBarOptions) Direction(d Direction) ProgressBarOpt {
	return func(s *ProgressBar) {
		s.direction = d
	}
}

// Inverted sets whether the progress bar is inverted.
// The default is false, which means from left to right or top to bottom.
func (o ProgressBarOptions) Inverted(inverted bool) ProgressBarOpt {
	return func(s *ProgressBar) {
		s.inverted = inverted
	}
}

func (o ProgressBarOptions) Images(track *ProgressBarImage, fill *ProgressBarImage) ProgressBarOpt {
	return func(s *ProgressBar) {
		s.definedParams.TrackImage = track
		s.definedParams.FillImage = fill
	}
}

func (o ProgressBarOptions) TrackImage(track *ProgressBarImage) ProgressBarOpt {
	return func(s *ProgressBar) {
		s.definedParams.TrackImage = track
	}
}

func (o ProgressBarOptions) FillImage(fill *ProgressBarImage) ProgressBarOpt {
	return func(s *ProgressBar) {
		s.definedParams.FillImage = fill
	}
}

func (o ProgressBarOptions) TrackPadding(i *Insets) ProgressBarOpt {
	return func(s *ProgressBar) {
		s.definedParams.TrackPadding = i
	}
}

func (o ProgressBarOptions) Values(min int, max int, current int) ProgressBarOpt {
	return func(s *ProgressBar) {
		s.Min = min
		s.Max = max
		s.current = current
	}
}

func (s *ProgressBar) GetWidget() *Widget {
	s.init.Do()
	return s.widget
}

func (s *ProgressBar) PreferredSize() (int, int) {
	w, h := s.computedParams.TrackImage.Idle.MinSize()
	if s.widget != nil {
		if w < s.widget.MinWidth {
			w = s.widget.MinWidth
		}
		if h < s.widget.MinHeight {
			h = s.widget.MinHeight
		}
	}
	return w, h
}

func (s *ProgressBar) SetLocation(rect img.Rectangle) {
	s.init.Do()
	s.widget.Rect = rect
}

func (s *ProgressBar) Render(screen *ebiten.Image) {
	s.init.Do()
	s.widget.Render(screen)
	s.draw(screen)
}

func (s *ProgressBar) Update(updObj *UpdateObject) {
	s.init.Do()

	s.widget.Update(updObj)
}

func (s *ProgressBar) draw(screen *ebiten.Image) {
	i := s.computedParams.TrackImage.Idle
	var fill *image.NineSlice
	if s.computedParams.FillImage != nil {
		fill = s.computedParams.FillImage.Idle
	}
	if s.widget.Disabled {
		if s.computedParams.TrackImage.Disabled != nil {
			i = s.computedParams.TrackImage.Disabled
		}
		if s.computedParams.FillImage.Disabled != nil {
			fill = s.computedParams.FillImage.Disabled
		}
	}

	if i != nil {
		i.Draw(screen, s.widget.Rect.Dx(), s.widget.Rect.Dy(), s.widget.drawImageOptions)
	}
	if fill != nil && s.currentPercentage() > 0 {
		fillX := s.widget.Rect.Dx() - s.computedParams.TrackPadding.Left - s.computedParams.TrackPadding.Right
		fillY := s.widget.Rect.Dy() - s.computedParams.TrackPadding.Top - s.computedParams.TrackPadding.Bottom
		if s.direction == DirectionHorizontal {
			fillX = int(float64(fillX) * s.currentPercentage())
		} else {
			fillY = int(float64(fillY) * s.currentPercentage())
		}
		fill.Draw(screen, fillX, fillY, func(opts *ebiten.DrawImageOptions) {
			tx := s.widget.Rect.Min.X + s.computedParams.TrackPadding.Left
			ty := s.widget.Rect.Min.Y + s.computedParams.TrackPadding.Top
			if s.inverted {
				if s.direction == DirectionHorizontal {
					tx = s.widget.Rect.Max.X - s.computedParams.TrackPadding.Right - fillX
				} else {
					ty = s.widget.Rect.Max.Y - s.computedParams.TrackPadding.Bottom - fillY
				}
			}
			opts.GeoM.Translate(float64(tx), float64(ty))
		})
	}
}

func (s *ProgressBar) currentPercentage() float64 {
	if s.Max <= s.Min {
		return 0
	}

	return float64(s.current-s.Min) / float64(s.Max-s.Min)
}

func (s *ProgressBar) SetCurrent(value int) bool {
	oldValue := s.current
	switch {
	case value < s.Min:
		s.current = s.Min
	case value > s.Max:
		s.current = s.Max
	default:
		s.current = value
	}
	return oldValue != s.current
}

func (s *ProgressBar) GetCurrent() int {
	return s.current
}

func (s *ProgressBar) createWidget() {
	s.widget = NewWidget(append([]WidgetOpt{WidgetOpts.TrackHover(true)}, s.widgetOpts...)...)
}
