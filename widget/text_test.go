package widget

import (
	"bufio"
	"image"
	"image/color"
	"strings"
	"testing"

	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/input"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/matryer/is"
)

func TestText_Lines(t *testing.T) {
	is := is.New(t)

	for _, label := range []string{"", "one", "one\ntwo", "one\r\ntwo", "one\n", "one\n\ntwo", "\n"} {
		tw := newText(t, label)

		var want []string
		s := bufio.NewScanner(strings.NewReader(label))
		for s.Scan() {
			want = append(want, s.Text())
		}
		is.Equal(tw.measurements.processedLineTexts, want) // lines of label
	}
}

func TestText_WrappedLinesKeepEveryCharacter(t *testing.T) {
	is := is.New(t)

	label := "a wrapped label with enough words in it to run over several lines"
	tw := newText(t, label, TextOpts.MaxWidth(120))

	is.True(len(tw.measurements.processedLineTexts) > 1)
	is.Equal(strings.Join(tw.measurements.processedLineTexts, ""), label)
}

func TestText_LinkHover(t *testing.T) {
	is := is.New(t)

	var entered []string
	tw := newText(t, "go [link=home]there[/link] now",
		TextOpts.ProcessBBCode(true),
		TextOpts.LinkCursorEnteredHandler(func(args *LinkEventArgs) {
			entered = append(entered, args.Id)
		}))
	tw.SetLocation(image.Rect(0, 0, 400, 100))

	// Move the link right; it is hovered where it is now drawn.
	tw.Label = "go further [link=home]there[/link] now"
	render(tw, t)

	face := *loadFont(t)
	before, _ := text.Measure("go further ", face, 0)
	link, _ := text.Measure("there", face, 0)
	cursor := &fixedCursor{x: int(before + link/2), y: int(tw.measurements.lineHeight / 2)}
	input.SetCursorUpdater(cursor)
	t.Cleanup(func() { input.SetCursorUpdater(nil) })

	tw.Update(&UpdateObject{})
	event.ExecuteDeferred()

	is.Equal(entered, []string{"home"}) // link entered
}

func newText(t *testing.T, label string, opts ...TextOpt) *Text {
	t.Helper()

	tw := NewText(append(opts, TextOpts.Text(label, loadFont(t), color.White))...)
	render(tw, t)
	return tw
}

// fixedCursor is a cursor that stays at one position, with no buttons pressed.
type fixedCursor struct{ x, y int }

func (c *fixedCursor) Update()                                         {}
func (c *fixedCursor) AfterUpdate()                                    {}
func (c *fixedCursor) Draw(*ebiten.Image)                              {}
func (c *fixedCursor) AfterDraw(*ebiten.Image)                         {}
func (c *fixedCursor) MouseButtonPressed(ebiten.MouseButton) bool      { return false }
func (c *fixedCursor) MouseButtonJustPressed(ebiten.MouseButton) bool  { return false }
func (c *fixedCursor) MouseButtonJustReleased(ebiten.MouseButton) bool { return false }
func (c *fixedCursor) CursorPosition() (int, int)                      { return c.x, c.y }
func (c *fixedCursor) GetCursorImage(string) *ebiten.Image             { return nil }
func (c *fixedCursor) GetCursorOffset(string) image.Point              { return image.Point{} }
