package widget

import (
	"testing"

	"github.com/ebitenui/ebitenui/event"
	"github.com/matryer/is"
)

func TestStateButton_SetState_Image(t *testing.T) {
	is := is.New(t)

	st := map[interface{}]*ButtonImage{
		1: {
			Idle: newNineSliceEmpty(t),
		},
		2: {
			Idle: newNineSliceEmpty(t),
		},
		3: {
			Idle: newNineSliceEmpty(t),
		},
	}

	s := newStateButton(t, StateButtonOpts.StateImages(st))

	s.State = 2
	render(s, t)

	is.Equal(stateButtonButton(s).Image, st[2])
}

func newStateButton(t *testing.T, opts ...StateButtonOpt) *StateButton {
	s := NewStateButton(opts...)
	event.ExecuteDeferred()
	render(s, t)
	return s
}

func stateButtonButton(s *StateButton) *Button {
	return s.button
}
