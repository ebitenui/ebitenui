package input

import (
	"image"
)

type InputLayerer interface {
	SetupInputLayer(def DeferredSetupInputLayerFunc)
}

type DeferredSetupInputLayerFunc func(s SetupInputLayerFunc)

type SetupInputLayerFunc func(def DeferredSetupInputLayerFunc)

type Layer struct {
	DebugLabel string
	EventTypes LayerEventType
	BlockLower bool
	FullScreen bool
	RectFunc   LayerRectFunc

	invalid bool
}

type LayerRectFunc func() image.Rectangle

type LayerEventType uint16

const (
	LayerEventTypeAny         = LayerEventType(0b0000000000000000)
	LayerEventTypeAll         = LayerEventType(0b1111111111111111)
	LayerEventTypeMouseButton = LayerEventType(0b0000000000000001)
	LayerEventTypeWheel       = LayerEventType(0b0000000000000010)
)

var DefaultLayer = Layer{
	DebugLabel: "default",
	EventTypes: LayerEventTypeAll,
	BlockLower: true,
	FullScreen: true,
}

var layers []*Layer

var deferredSetupInputLayers []SetupInputLayerFunc

func AddLayer(l *Layer) {
	if !l.Valid() {
		panic("invalid layer")
	}

	if l.EventTypes == LayerEventTypeAny {
		panic("LayerEventTypeAny is invalid for an input layer, perhaps you meant to use LayerEventTypeAll instead")
	}

	layers = append(layers, l)
}

func (l *Layer) Valid() bool {
	return !l.invalid
}

func (l *Layer) ActiveFor(x int, y int, eventType LayerEventType) bool {
	if !l.Valid() {
		return false
	}

	for i := len(layers) - 1; i >= 0; i-- {
		layer := layers[i]

		if !layer.contains(x, y) {
			continue
		}

		if eventType != LayerEventTypeAny && layer.EventTypes&eventType != eventType {
			continue
		}

		if layer != l {
			if layer.BlockLower {
				return false
			}
			continue
		}

		return true
	}

	return l == &DefaultLayer
}

func (l *Layer) contains(x int, y int) bool {
	if l.FullScreen {
		return true
	}
	return image.Point{x, y}.In(l.RectFunc())
}

func SetupInputLayersWithDeferred(i InputLayerer) {
	for _, l := range layers {
		l.invalid = true
	}
	layers = layers[:0]

	appendToDeferredSetupInputLayerQueue(i.SetupInputLayer)
	setupDeferredInputLayers()
}

func setupDeferredInputLayers() {
	for len(deferredSetupInputLayers) > 0 {
		s := deferredSetupInputLayers[0]
		deferredSetupInputLayers = deferredSetupInputLayers[1:]

		s(appendToDeferredSetupInputLayerQueue)
	}
}

func appendToDeferredSetupInputLayerQueue(s SetupInputLayerFunc) {
	deferredSetupInputLayers = append(deferredSetupInputLayers, s)
}
