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
	is.True(c.effectBuffer == buffer)                          // the same off-screen image
	is.Equal(child.drawnOn.Bounds(), img.Rect(30, 20, 80, 50)) // at the new place

	c.SetOpacity(1)
	is.Equal(c.effectBuffer, nil) // released
}

// The off-screen image covers the container and a margin around it, within
// the screen: a short slide keeps it, a move past the margin makes another.
func TestContainer_EffectBufferFitsTheContainer(t *testing.T) {
	is := is.New(t)

	screen := ebiten.NewImage(1000, 600)
	c, _ := newRecordingContainer(t)
	c.SetLocation(img.Rect(300, 200, 500, 300))
	c.SetOpacity(0.5)
	c.Render(screen)
	is.Equal(c.effectBuffer.Bounds(), img.Rect(236, 136, 564, 364)) // 64 px around it, not the screen

	buffer := c.effectBuffer
	c.SetLocation(img.Rect(340, 200, 540, 300)) // 40 px along, as a popup easing in
	c.Render(screen)
	is.True(c.effectBuffer == buffer) // kept

	c.SetLocation(img.Rect(700, 200, 900, 300))
	c.Render(screen)
	is.True(c.effectBuffer != buffer) // made again around it
	is.Equal(c.effectBuffer.Bounds(), img.Rect(636, 136, 964, 364))

	c.SetLocation(img.Rect(950, 550, 1050, 650)) // partly off the screen
	c.Render(screen)
	is.Equal(c.effectBuffer.Bounds(), img.Rect(886, 486, 1000, 600)) // up to its edge
}

// A fading container keeps its off-screen image; an invisible one, such as
// a window faded out as it closes, keeps none and makes none.
func TestContainer_InvisibleKeepsNoBuffer(t *testing.T) {
	is := is.New(t)

	screen := ebiten.NewImage(100, 100)
	c, _ := newRecordingContainer(t)
	c.SetOpacity(0.5)
	c.Render(screen)
	buffer := c.effectBuffer
	is.True(buffer != nil)
	c.SetOpacity(0.25) // fading
	c.Render(screen)
	is.True(c.effectBuffer == buffer) // kept through the fade

	c.SetOpacity(0)
	is.Equal(c.effectBuffer, nil) // released
	c.Render(screen)
	is.Equal(c.effectBuffer, nil) // and none made
}

func TestWindow_SetOpacity(t *testing.T) {
	is := is.New(t)

	w := NewWindow(WindowOpts.Contents(NewContainer()))
	w.SetOpacity(0.25)

	c, ok := w.GetContainer().(*Container)
	is.True(ok)
	is.Equal(c.Opacity(), float32(0.25))
}
