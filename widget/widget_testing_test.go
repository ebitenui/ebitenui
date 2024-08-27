package widget

import (
	img "image"
	"log"
	"os"
	"sync"
	"testing"

	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type simpleWidget struct {
	widget          *Widget
	preferredWidth  int
	preferredHeight int
}

var loadFontOnce sync.Once
var fontFace2 text.Face

func newSimpleWidget(preferredWidth int, preferredHeight int, ld interface{}) *simpleWidget {
	return &simpleWidget{
		widget:          NewWidget(WidgetOpts.LayoutData(ld)),
		preferredWidth:  preferredWidth,
		preferredHeight: preferredHeight,
	}
}

func (s *simpleWidget) GetWidget() *Widget {
	return s.widget
}

func (s *simpleWidget) PreferredSize() (int, int) {
	return s.preferredWidth, s.preferredHeight
}

func (s *simpleWidget) SetLocation(rect img.Rectangle) {
	s.widget.Rect = rect
}

func loadFont(t *testing.T) text.Face {
	t.Helper()

	loadFontOnce.Do(func() {

		data, err := os.Open("testdata/fonts/NotoSans-Regular.ttf")
		if err != nil {
			panic(err)
		}
		s, err := text.NewGoTextFaceSource(data)
		if err != nil {
			log.Fatal(err)
		}

		fontFace2 = &text.GoTextFace{
			Source: s,
			Size:   20,
		}
	})

	return fontFace2
}

func newImageEmpty(t *testing.T) *ebiten.Image {
	t.Helper()
	return newImageEmptySize(1, 1, t)
}

func newImageEmptySize(width int, height int, t *testing.T) *ebiten.Image {
	t.Helper()
	return ebiten.NewImage(width, height)
}

func newNineSliceEmpty(t *testing.T) *image.NineSlice {
	t.Helper()
	return image.NewNineSliceSimple(newImageEmpty(t), 0, 0)
}

func leftMouseButtonClick(w HasWidget, t *testing.T) {
	t.Helper()
	leftMouseButtonPress(w, t)
	leftMouseButtonRelease(w, t)
}

func leftMouseButtonPress(w HasWidget, t *testing.T) {
	t.Helper()

	w.GetWidget().MouseButtonPressedEvent.Fire(&WidgetMouseButtonPressedEventArgs{
		Widget:  w.GetWidget(),
		Button:  ebiten.MouseButtonLeft,
		OffsetX: 0,
		OffsetY: 0,
	})

	event.ExecuteDeferred()
}

func leftMouseButtonRelease(w HasWidget, t *testing.T) {
	t.Helper()

	w.GetWidget().MouseButtonReleasedEvent.Fire(&WidgetMouseButtonReleasedEventArgs{
		Widget:  w.GetWidget(),
		Button:  ebiten.MouseButtonLeft,
		OffsetX: 0,
		OffsetY: 0,
		Inside:  true,
	})

	event.ExecuteDeferred()
}

func render(r Renderer, t *testing.T) {
	t.Helper()

	screen := ebiten.NewImage(1, 1)
	r.Render(screen)
	event.ExecuteDeferred()
}
