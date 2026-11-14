package widget

import (
	"testing"

	"github.com/ebitenui/ebitenui/input"
	"github.com/matryer/is"
)

// Setting up the input layers of children that don't elevate them
// allocates nothing.
func TestContainer_SetupInputLayerAllocatesNothing(t *testing.T) {
	is := is.New(t)

	c := NewContainer()
	for range 10 {
		c.AddChild(newSimpleWidget(10, 10, nil))
	}
	c.Validate()
	def := func(input.SetupInputLayerFunc) {}

	is.Equal(testing.AllocsPerRun(100, func() { c.SetupInputLayer(def) }), 0.0)
}
