package widget

import (
	"bufio"
	"image"
	"image/color"
	"slices"
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

func TestText_GlyphsFollowText(t *testing.T) {
	is := is.New(t)

	small := *loadFont(t)
	goFace, ok := small.(*text.GoTextFace)
	is.True(ok)
	var big text.Face = &text.GoTextFace{Source: goFace.Source, Size: 30}

	screen := ebiten.NewImage(400, 200)
	draw := func(tw *Text) {
		tw.Validate()
		tw.Render(screen)
	}

	plain := newText(t, "first line\nsecond line")
	bbcode := newText(t, "a [color=#ff0000]red[/color] word", TextOpts.ProcessBBCode(true))
	for _, change := range []func(tw *Text){
		func(*Text) {},
		func(tw *Text) { tw.Label = "another label\nwith [color=#00ff00]two[/color] lines" },
		func(tw *Text) { tw.SetFace(&big) },
	} {
		change(plain)
		draw(plain)
		total := 0
		for i, line := range plain.measurements.processedLineTexts {
			want := visibleGlyphs(line, *plain.computedParams.Face)
			is.True(slices.Equal(plain.measurements.processedLineGlyphs[i], want)) // line glyphs
			total += len(want)
		}
		is.Equal(len(plain.glyphs), total) // the buffer holds only the current glyphs

		change(bbcode)
		draw(bbcode)
		for _, pieces := range bbcode.measurements.processedLines {
			for _, piece := range pieces {
				is.True(slices.Equal(piece.glyphs, visibleGlyphs(piece.text, *bbcode.computedParams.Face))) // piece glyphs
			}
		}
	}
}

func TestText_LaysOutOnlyLinesDrawn(t *testing.T) {
	is := is.New(t)

	tw := newText(t, "first line\nsecond line") // drawn on a 1x1 screen

	is.Equal(tw.measurements.processedLinesLaid, []bool{true, false})
	is.Equal(tw.measurements.processedLineGlyphs[1], nil) // the line below the screen
}

func TestText_FaceChangedInPlace(t *testing.T) {
	is := is.New(t)

	goFace, ok := (*loadFont(t)).(*text.GoTextFace)
	is.True(ok)
	face := *goFace
	var tf text.Face = &face
	tw := NewText(TextOpts.Text("Hello", &tf, color.White))
	render(tw, t)
	before := tw.measurements.boundingBoxWidth

	face.Size *= 2 // a new size, without SetFace
	render(tw, t)

	is.True(tw.measurements.boundingBoxWidth > 1.5*before)                                    // measured again
	is.True(slices.Equal(tw.measurements.processedLineGlyphs[0], visibleGlyphs("Hello", tf))) // laid out again
}

func TestText_SetFaceLaysOutAgain(t *testing.T) {
	is := is.New(t)

	limited := text.NewLimitedFace(*loadFont(t))
	limited.AddUnicodeRange('a', 'b') // 'c' falls back to the missing glyph
	var tf text.Face = limited
	tw := NewText(TextOpts.Text("abc", &tf, color.White))
	render(tw, t)
	before := slices.Clone(tw.measurements.processedLineGlyphs[0])

	limited.AddUnicodeRange('c', 'c') // the same metrics
	tw.SetFace(&tf)
	render(tw, t)

	got := tw.measurements.processedLineGlyphs[0]
	is.True(slices.Equal(got, visibleGlyphs("abc", tf))) // laid out again
	is.True(!slices.Equal(got, before))                  // with the real 'c'
}

// visibleGlyphs returns the glyphs text.AppendGlyphs lays out for s, without
// those that have no image.
func visibleGlyphs(s string, face text.Face) []text.Glyph {
	var glyphs []text.Glyph
	for _, g := range text.AppendGlyphs(nil, s, face, nil) {
		if g.Image != nil {
			glyphs = append(glyphs, g)
		}
	}
	return glyphs
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
