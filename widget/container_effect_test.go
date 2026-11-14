package widget

import (
	img "image"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/matryer/is"
)

// recordingWidget remembers the image it was last drawn onto.
type recordingWidget struct {
	widget  *Widget
	drawnOn *ebiten.Image
}

func (r *recordingWidget) GetWidget() *Widget             { return r.widget }
func (r *recordingWidget) PreferredSize() (int, int)      { return 10, 10 }
func (r *recordingWidget) SetLocation(rect img.Rectangle) { r.widget.Rect = rect }
func (r *recordingWidget) Validate()                      {}
func (r *recordingWidget) Render(screen *ebiten.Image)    { r.drawnOn = screen }

func newRecordingContainer(t *testing.T) (*Container, *recordingWidget) {
	t.Helper()
	child := &recordingWidget{widget: NewWidget()}
	c := NewContainer(ContainerOpts.Layout(NewAnchorLayout()))
	c.AddChild(child)
	c.Validate()
	c.SetLocation(img.Rect(10, 10, 60, 40))
	return c, child
}

func TestContainer_DrawsDirectlyWhenOpaque(t *testing.T) {
	is := is.New(t)

	screen := ebiten.NewImage(100, 100)
	c, child := newRecordingContainer(t)
	c.Render(screen)

	is.Equal(c.Opacity(), float32(1))
	is.Equal(child.drawnOn, screen) // drawn onto the screen
	is.Equal(c.effectBuffer, nil)   // no off-screen image
}

func TestContainer_OpacityDrawsOffScreen(t *testing.T) {
	is := is.New(t)

	screen := ebiten.NewImage(100, 100)
	c, child := newRecordingContainer(t)
	c.SetOpacity(0.5)
	c.Render(screen)

	is.True(child.drawnOn != screen)                           // drawn off screen
	is.Equal(child.drawnOn.Bounds(), img.Rect(10, 10, 60, 40)) // covering the container

	buffer := c.effectBuffer
	c.SetLocation(img.Rect(30, 20, 80, 50)) // sliding, as a window opening
	c.Render(screen)
	is.Equal(c.effectBuffer, buffer)                           // the same off-screen image
	is.Equal(child.drawnOn.Bounds(), img.Rect(30, 20, 80, 50)) // at the new place

	c.SetOpacity(1)
	is.Equal(c.effectBuffer, nil) // released
}

func TestWindow_SetOpacity(t *testing.T) {
	is := is.New(t)

	w := NewWindow(WindowOpts.Contents(NewContainer()))
	w.SetOpacity(0.25)

	c, ok := w.GetContainer().(*Container)
	is.True(ok)
	is.Equal(c.Opacity(), float32(0.25))
}
