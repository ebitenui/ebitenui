package input

import (
	"image"
	"testing"

	"github.com/matryer/is"
	"github.com/stretchr/testify/mock"
)

type layererMock struct {
	mock.Mock
	setupInputLayerCall *mock.Call
}

func TestSetupInputLayersWithDeferred(t *testing.T) {
	l := newLayererMock(nil)
	SetupInputLayersWithDeferred([]Layerer{l})
	l.AssertExpectations(t)
}

func TestSetupInputLayersWithDeferred_Deferred(t *testing.T) {
	is := is.New(t)

	called := false
	l := newLayererMock(func(def DeferredSetupInputLayerFunc) {
		def(func(_ DeferredSetupInputLayerFunc) {
			called = true
		})
	})

	SetupInputLayersWithDeferred([]Layerer{l})

	is.True(called)
}

func TestLayer_ActiveFor_FullScreen(t *testing.T) {
	is := is.New(t)

	l1 := Layer{
		EventTypes: LayerEventTypeAll,
		BlockLower: false,
		FullScreen: true,
	}
	AddLayer(&l1)

	l2 := Layer{
		EventTypes: LayerEventTypeAll,
		BlockLower: false,
		FullScreen: true,
	}
	AddLayer(&l2)

	is.True(l1.ActiveFor(100, 100, LayerEventTypeMouseButton))
	is.True(l2.ActiveFor(100, 100, LayerEventTypeMouseButton))
}

func TestLayer_ActiveFor_BlockLower(t *testing.T) {
	is := is.New(t)

	l1 := Layer{
		EventTypes: LayerEventTypeAll,
		BlockLower: false,
		FullScreen: true,
	}
	AddLayer(&l1)

	l2 := Layer{
		EventTypes: LayerEventTypeAll,
		BlockLower: true,
		FullScreen: true,
	}
	AddLayer(&l2)

	is.True(!l1.ActiveFor(100, 100, LayerEventTypeMouseButton))
	is.True(l2.ActiveFor(100, 100, LayerEventTypeMouseButton))
}

func TestLayer_ActiveFor_Rect(t *testing.T) {
	is := is.New(t)

	l1 := Layer{
		EventTypes: LayerEventTypeAll,
		BlockLower: false,
		RectFunc: func() image.Rectangle {
			return image.Rect(0, 0, 50, 50)
		},
	}
	AddLayer(&l1)

	l2 := Layer{
		EventTypes: LayerEventTypeAll,
		BlockLower: false,
		RectFunc: func() image.Rectangle {
			return image.Rect(20, 20, 70, 70)
		},
	}
	AddLayer(&l2)

	is.True(l1.ActiveFor(10, 10, LayerEventTypeMouseButton))
	is.True(!l2.ActiveFor(10, 10, LayerEventTypeMouseButton))

	is.True(l1.ActiveFor(30, 30, LayerEventTypeMouseButton))
	is.True(l2.ActiveFor(30, 30, LayerEventTypeMouseButton))

	is.True(!l1.ActiveFor(100, 100, LayerEventTypeMouseButton))
	is.True(!l2.ActiveFor(100, 100, LayerEventTypeMouseButton))
}

func TestLayer_ActiveFor_EventType(t *testing.T) {
	is := is.New(t)

	l1 := Layer{
		EventTypes: LayerEventTypeAll,
		BlockLower: false,
		FullScreen: true,
	}
	AddLayer(&l1)

	l2 := Layer{
		EventTypes: LayerEventTypeMouseButton,
		BlockLower: false,
		FullScreen: true,
	}
	AddLayer(&l2)

	is.True(l1.ActiveFor(100, 100, LayerEventTypeMouseButton))
	is.True(l2.ActiveFor(100, 100, LayerEventTypeMouseButton))

	is.True(l1.ActiveFor(100, 100, LayerEventTypeWheel))
	is.True(!l2.ActiveFor(100, 100, LayerEventTypeWheel))
}

func newLayererMock(f SetupInputLayerFunc) *layererMock {
	l := layererMock{}
	l.setupInputLayerCall = l.On("SetupInputLayer", mock.Anything)
	if f != nil {
		l.setupInputLayerCall.Run(func(args mock.Arguments) {
			def := args[0].(DeferredSetupInputLayerFunc)
			f(def)
		})
	}
	return &l
}

func (l *layererMock) SetupInputLayer(def DeferredSetupInputLayerFunc) {
	l.Called(def)
}
