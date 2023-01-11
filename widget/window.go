package widget

import (
	"image"

	"github.com/ebitenui/ebitenui/input"
	"github.com/hajimehoshi/ebiten/v2"
)

type RemoveWindowFunc func()

type Window struct {
	Modal        bool
	closeOnClick bool
	closeFunc    RemoveWindowFunc
	Contents     *Container
}

type WindowOpt func(w *Window)

type WindowOptions struct {
}

var WindowOpts WindowOptions

func NewWindow(opts ...WindowOpt) *Window {
	w := &Window{}

	for _, o := range opts {
		o(w)
	}
	if w.closeOnClick {
		w.Contents.GetWidget().MouseButtonReleasedEvent.AddHandler(func(args interface{}) {
			a := args.(*WidgetMouseButtonReleasedEventArgs)
			if !a.Inside && w.closeFunc != nil {
				w.closeFunc()
			}
		})
	}
	return w
}

func (o WindowOptions) Contents(c *Container) WindowOpt {
	return func(w *Window) {
		w.Contents = c
	}
}

func (o WindowOptions) Modal() WindowOpt {
	return func(w *Window) {
		w.Modal = true
	}
}

func (o WindowOptions) CloseOnClickOut() WindowOpt {
	return func(w *Window) {
		w.closeOnClick = true
	}
}

func (w *Window) SetCloseFunction(close RemoveWindowFunc) {
	w.closeFunc = close
}

func (w *Window) Close() {
	if w.closeFunc != nil {
		w.closeFunc()
	}
}

func (w *Window) SetLocation(rect image.Rectangle) {
	w.Contents.SetLocation(rect)
}

func (w *Window) RequestRelayout() {
	w.Contents.RequestRelayout()
}

func (w *Window) SetupInputLayer(def input.DeferredSetupInputLayerFunc) {
	if w.Modal {
		w.Contents.GetWidget().ElevateToNewInputLayer(&input.Layer{
			DebugLabel: "modal window",
			EventTypes: input.LayerEventTypeAll,
			BlockLower: true,
			FullScreen: true,
		})
	}
}

func (w *Window) Render(screen *ebiten.Image, def DeferredRenderFunc) {
	w.Contents.Render(screen, def)
}
