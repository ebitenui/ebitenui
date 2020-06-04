package widget

import (
	img "image"

	"github.com/blizzy78/ebitenui/input"
	"github.com/hajimehoshi/ebiten"
)

// A FlipBook is a container that always renders exactly one child widget: the current page.
// The current page will be embedded in a FillLayout.
type FlipBook struct {
	containerOpts  []ContainerOpt
	fillLayoutOpts []FillLayoutOpt

	init          *MultiOnce
	container     *Container
	removeCurrent RemoveChildFunc
}

// FlipBookOpt is a function that configures f.
type FlipBookOpt func(f *FlipBook)

// FlipBookOpts contains functions that configure a FlipBook.
const FlipBookOpts = flipBookOpts(true)

type flipBookOpts bool

// NewFlipBook constructs a new FlipBook configured with opts.
func NewFlipBook(opts ...FlipBookOpt) *FlipBook {
	f := &FlipBook{
		init: &MultiOnce{},
	}

	f.init.Append(f.createWidget)

	for _, o := range opts {
		o(f)
	}

	return f
}

// WithContainerOpts configures a FlipBook with opts.
func (o flipBookOpts) WithContainerOpts(opts ...ContainerOpt) FlipBookOpt {
	return func(f *FlipBook) {
		f.containerOpts = append(f.containerOpts, opts...)
	}
}

// WithPadding configures a FlipBook with padding i.
func (o flipBookOpts) WithPadding(i Insets) FlipBookOpt {
	return func(f *FlipBook) {
		f.fillLayoutOpts = append(f.fillLayoutOpts, FillLayoutOpts.WithPadding(i))
	}
}

// GetWidget implements HasWidget.
func (f *FlipBook) GetWidget() *Widget {
	f.init.Do()
	return f.container.GetWidget()
}

// PreferredSize implements PreferredSizer.
func (f *FlipBook) PreferredSize() (int, int) {
	f.init.Do()
	return f.container.PreferredSize()
}

// SetLocation implements Locateable.
func (f *FlipBook) SetLocation(rect img.Rectangle) {
	f.init.Do()
	f.container.SetLocation(rect)
}

// RequestRelayout implements Relayoutable.
func (f *FlipBook) RequestRelayout() {
	f.init.Do()
	f.container.RequestRelayout()
}

// SetupInputLayer implements InputLayerer.
func (f *FlipBook) SetupInputLayer(def input.DeferredSetupInputLayerFunc) {
	f.init.Do()
	f.container.SetupInputLayer(def)
}

// Render implements Renderer.
func (f *FlipBook) Render(screen *ebiten.Image, def DeferredRenderFunc) {
	f.init.Do()
	f.container.Render(screen, def)
}

func (f *FlipBook) createWidget() {
	f.container = NewContainer(append(f.containerOpts, []ContainerOpt{
		ContainerOpts.WithLayout(NewFillLayout(f.fillLayoutOpts...)),
	}...)...)
	f.containerOpts = nil
	f.fillLayoutOpts = nil
}

// SetPage sets the current page to be rendered to page. The previous page will no longer be rendered.
//
// Note that when switching to a new page, it may be necessary to re-layout parent containers if the pages
// are of different sizes.
func (f *FlipBook) SetPage(page HasWidget) {
	f.init.Do()

	if f.removeCurrent != nil {
		f.removeCurrent()
	}

	f.removeCurrent = f.container.AddChild(page)
}
