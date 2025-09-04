package widget

import (
	"image"
	"testing"

	"github.com/ebitenui/ebitenui/input"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/matryer/is"
	"github.com/stretchr/testify/mock"
)

type controlMock struct {
	mock.Mock
}

func (s *controlMock) Validate() {
}

func TestContainer_Render(t *testing.T) {
	w := NewWidget()
	m := controlMock{}
	m.On("GetWidget").Maybe().Return(w)
	m.On("PreferredSize").Maybe().Return(50, 50)
	m.On("SetLocation", mock.Anything).Maybe()
	m.On("Render", mock.Anything, mock.Anything)

	c := newContainer(t,
		ContainerOpts.Layout(newRowLayout(t)))
	c.AddChild(&m)

	render(c, t)

	m.AssertExpectations(t)
}

func TestContainer_Render_AutoDisableChildren(t *testing.T) {
	is := is.New(t)

	w := NewWidget()
	m := controlMock{}
	m.On("GetWidget").Maybe().Return(w)
	m.On("PreferredSize").Maybe().Return(50, 50)
	m.On("SetLocation", mock.Anything).Maybe()
	m.On("Render", mock.Anything, mock.Anything).Maybe()

	c := newContainer(t,
		ContainerOpts.AutoDisableChildren(),
		ContainerOpts.Layout(newRowLayout(t)))
	c.AddChild(&m)

	c.widget.Disabled = true
	render(c, t)

	is.True(w.Disabled)
}

func TestContainer_SetupInputLayer(t *testing.T) {
	def := func(s input.SetupInputLayerFunc) {
		// nothing to do
	}

	w := NewWidget()
	m := controlMock{}
	m.On("GetWidget").Maybe().Return(w)
	m.On("SetupInputLayer", mock.AnythingOfType("input.DeferredSetupInputLayerFunc"))

	c := newContainer(t,
		ContainerOpts.Layout(newRowLayout(t)))
	c.AddChild(&m)

	c.SetupInputLayer(def)

	m.AssertExpectations(t)
}

func TestContainer_Replace(t *testing.T) {
	c := newContainer(t)
	m1 := controlMock{}
	m2 := controlMock{}
	m3 := controlMock{}

	m1.On("GetWidget").Maybe().Return(NewWidget())
	m2.On("GetWidget").Maybe().Return(NewWidget())
	m3.On("GetWidget").Maybe().Return(NewWidget())

	c.AddChild(&m1, &m2, &m3)

	if len(c.children) != 3 {
		t.Fatalf("expected 3 children, got %d", len(c.children))
	}

	replaced := controlMock{}
	replaced.On("GetWidget").Maybe().Return(NewWidget())

	c.ReplaceChild(&m2, &replaced)

	if len(c.children) != 3 {
		t.Fatalf("expected 3 children, got %d", len(c.children))
	}

	if c.children[0] != &m1 {
		t.Fatalf("expected original child at index 0")
	}

	if c.children[1] != &replaced {
		t.Fatalf("expected replaced child at index 1")
	}

	if c.children[2] != &m3 {
		t.Fatalf("expected original child at index 2")
	}
}

func (c *controlMock) GetWidget() *Widget {
	args := c.Called()
	if arg, ok := args.Get(0).(*Widget); ok {
		return arg
	}
	return nil
}

func (c *controlMock) PreferredSize() (int, int) {
	c.Called()
	return -1, -1
}

func (c *controlMock) SetLocation(rect image.Rectangle) {
	c.Called(rect)
}

func (c *controlMock) Render(screen *ebiten.Image) {
	c.Called(screen)
}

func (c *controlMock) SetupInputLayer(def input.DeferredSetupInputLayerFunc) {
	c.Called(def)
}

func newContainer(t *testing.T, opts ...ContainerOpt) *Container {
	t.Helper()
	return NewContainer(opts...)
}
