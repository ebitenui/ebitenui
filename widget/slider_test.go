package widget

import (
	"testing"

	"github.com/ebitenui/ebitenui/event"
	"github.com/matryer/is"
)

func TestSlider_Current_Initial(t *testing.T) {
	is := is.New(t)

	var eventArgs *SliderChangedEventArgs
	s := newSlider(t,
		SliderOpts.MinMax(10, 20),
		SliderOpts.ChangedHandler(func(args *SliderChangedEventArgs) {
			eventArgs = args
		}))

	is.Equal(s.Current, 10)
	is.Equal(eventArgs.Current, 10)
}

func newSlider(t *testing.T, opts ...SliderOpt) *Slider {
	s := NewSlider(append(opts, SliderOpts.Images(&SliderTrackImage{
		Idle: newNineSliceEmpty(t),
	}, &ButtonImage{
		Idle: newNineSliceEmpty(t),
	}))...)
	event.ExecuteDeferred()
	render(s, t)
	return s
}
