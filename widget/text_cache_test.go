package widget

import (
	"bytes"
	img "image"
	"image/color"
	"strings"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/matryer/is"
)

// cachedText returns a text drawn from an image of itself, 200 by 100 px,
// drawn once.
func cachedText(t *testing.T, label string, opts ...TextOpt) *Text {
	t.Helper()
	SetTextCache(true)
	t.Cleanup(func() { SetTextCache(false) })
	txt := NewText(append([]TextOpt{TextOpts.Text(label, loadFont(t), color.White)}, opts...)...)
	txt.SetLocation(img.Rect(0, 0, 200, 100))
	txt.Validate()
	txt.Render(ebiten.NewImage(200, 100))
	return txt
}

// A cached text draws its glyphs into its image once, and again only when
// what it shows changes; moving it moves the image.
func TestText_CacheDrawsOnceUntilChanged(t *testing.T) {
	is := is.New(t)

	txt := cachedText(t, "Hello\nworld", TextOpts.Position(TextPositionCenter, TextPositionStart))
	screen := ebiten.NewImage(200, 100)
	is.Equal(txt.cache.renders, 1)
	txt.Render(screen)
	txt.Render(screen)
	is.Equal(txt.cache.renders, 1) // drawn from the image

	for i, change := range []func(){
		func() { txt.Label = "Hello" },
		func() { txt.SetColor(color.Black) },
		func() { txt.SetShadow(&TextShadow{Color: color.Black, OffsetX: 1, OffsetY: 1}) },
		func() { txt.SetLocation(img.Rect(0, 0, 150, 100)) }, // centered lines move
		func() { txt.SetPadding(&Insets{Left: 4}) },
		func() { txt.SetFace(loadFont(t)) },
		func() { txt.SetPosition(&TextPositioning{HTextPosition: TextPositionEnd}) },
	} {
		change()
		txt.Render(screen)
		is.Equal(txt.cache.renders, i+2)
	}

	txt.SetLocation(img.Rect(10, 10, 160, 110)) // moved, same width
	txt.Render(screen)
	is.Equal(txt.cache.renders, 8) // the image moves with it
}

// A text aligned at its start doesn't depend on its width; one showing a
// hovered link draws again in the link's hover color.
func TestText_CacheWidthAndLinks(t *testing.T) {
	is := is.New(t)

	txt := cachedText(t, "A [link=x]link[/link]", TextOpts.ProcessBBCode(true),
		TextOpts.LinkColor(&TextLinkColor{Idle: color.White, Hover: color.Black}))
	screen := ebiten.NewImage(200, 100)
	txt.SetLocation(img.Rect(0, 0, 150, 100))
	txt.Render(screen)
	is.Equal(txt.cache.renders, 1) // same lines at any width

	txt.previousLink = txt.measurements.processedLines[0][1] // hovered, as handleLinkEvents leaves it
	txt.Render(screen)
	is.Equal(txt.cache.renders, 2)

	txt.computedParams.LinkColor = &TextLinkColor{Idle: color.Black, Hover: color.White} // as a new theme gives
	txt.Render(screen)
	is.Equal(txt.cache.renders, 3)
}

// A shorter text draws into the image it has; a larger one gets a new image.
func TestText_CacheKeepsItsImage(t *testing.T) {
	is := is.New(t)

	txt := cachedText(t, "Hello")
	screen := ebiten.NewImage(200, 100)
	image := txt.cache.image
	is.True(image != nil)

	txt.Label = "Hell"
	txt.Render(screen)
	is.True(txt.cache.image == image) // reused

	txt.Label = "Hello, world"
	txt.Render(screen)
	is.True(txt.cache.image != image) // too small, replaced
	is.True(txt.cache.image.Bounds().Dx() > image.Bounds().Dx())

	txt.Label = ""
	txt.Render(screen)
	is.True(txt.cache.image == nil) // nothing left to draw
}

// A text taller or wider than textCacheMaxSide draws glyph by glyph, without
// trying again each frame; an empty one has nothing to keep.
func TestText_CacheLimits(t *testing.T) {
	is := is.New(t)

	tall := cachedText(t, strings.Repeat("line\n", 200))
	tall.Render(ebiten.NewImage(200, 100))
	is.True(tall.cache.image == nil)
	is.True(tall.cache.tooLarge)

	empty := cachedText(t, "")
	is.True(empty.cache.image == nil)
}

// Without SetTextCache, texts keep no image.
func TestText_CacheOff(t *testing.T) {
	is := is.New(t)

	txt := NewText(TextOpts.Text("Hello", loadFont(t), color.White))
	txt.SetLocation(img.Rect(0, 0, 200, 100))
	render(txt, t)
	is.True(txt.cache == nil)
}

// Drawn from its image, a text allocates no more than one draw does (before
// the game runs, Ebitengine allocates as it queues each draw).
func TestText_CacheDrawsWithoutAllocating(t *testing.T) {
	is := is.New(t)

	txt := cachedText(t, "Hello\nworld", TextOpts.Shadow(color.Black, 1, 1))
	screen := ebiten.NewImage(200, 100)
	op := &ebiten.DrawImageOptions{}
	oneDraw := testing.AllocsPerRun(10, func() { screen.DrawImage(txt.cache.image, op) })
	is.Equal(testing.AllocsPerRun(10, func() { txt.Render(screen) }), oneDraw)
}

func BenchmarkText_Draw(b *testing.B) {
	src, err := text.NewGoTextFaceSource(bytes.NewReader(data))
	if err != nil {
		b.Fatal(err)
	}
	var face text.Face = &text.GoTextFace{Source: src, Size: 20}
	label := strings.Repeat(strings.Repeat("Lorem ipsum dolor sit amet. ", 2)+"\n", 20) // 1 100 glyphs
	for _, cache := range []bool{false, true} {
		name := "glyphs"
		if cache {
			name = "image"
		}
		b.Run(name, func(b *testing.B) {
			SetTextCache(cache)
			defer SetTextCache(false)
			txt := NewText(TextOpts.Text(label, &face, color.White))
			txt.SetLocation(img.Rect(0, 0, 800, 600))
			txt.Validate()
			screen := ebiten.NewImage(800, 600)
			txt.Render(screen)
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				txt.Render(screen)
			}
		})
	}
}
