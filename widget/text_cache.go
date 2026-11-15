package widget

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// textCacheOn is SetTextCache's setting.
var textCacheOn bool

// textCacheMaxSide is the longest side of a text kept as an image, in
// pixels; a larger text draws glyph by glyph.
const textCacheMaxSide = 2048

// SetTextCache sets whether texts draw from images of themselves. Each Text
// draws its glyphs into an image when what it shows changes (its label,
// face, color, shadow, padding, alignment, width when centered or aligned
// at the end, a hovered link) and draws that image every frame: one draw
// instead of one per glyph, for an image the size of the text in GPU memory.
// Each image is a texture of its own, outside Ebitengine's shared atlases:
// small images that last and are drawn into now and then would otherwise
// keep atlas pages from being reused, and the atlases grow. It covers every
// Text, those inside buttons, labels and lists too; a text over 2048 pixels
// on a side still draws glyph by glyph. Set it before the interface is
// first drawn.
func SetTextCache(on bool) {
	textCacheOn = on
}

// textCache is a text drawn into an image, kept while what the text shows
// stays the same.
type textCache struct {
	image    *ebiten.Image // nil while the text has no glyphs
	at       image.Point   // where the image's top left goes, from the text's origin
	key      textCacheKey
	valid    bool // key describes image
	tooLarge bool // the text is drawn glyph by glyph
	renders  int  // times the text was drawn into the image
	op       ebiten.DrawImageOptions
}

// textCacheKey is what a text's image depends on.
type textCacheKey struct {
	measured   int // the measurement it was drawn from
	width      int // the text's width, when lines are centered or end-aligned
	position   TextPosition
	padding    Insets
	color      [4]uint32
	shadow     bool
	shadowRGBA [4]uint32
	shadowAt   [2]float64
	linkIdle   [4]uint32
	linkHover  [4]uint32
	hovered    *bbCodeText // the hovered link's piece
}

// glyphBox is the box glyphs cover, in pixels.
type glyphBox struct {
	minX, minY, maxX, maxY float64
}

func rgba(c color.Color) (v [4]uint32) {
	if c != nil {
		v[0], v[1], v[2], v[3] = c.RGBA()
	}
	return v
}

// cacheKey returns what the text's image depends on now, the text being w
// pixels wide.
func (t *Text) cacheKey(w int) textCacheKey {
	p := t.computedParams
	k := textCacheKey{
		measured: t.measured,
		position: p.Position.HTextPosition,
		padding:  *p.Padding,
		color:    rgba(p.Color),
		hovered:  t.previousLink,
	}
	if k.position != TextPositionStart {
		k.width = w
	}
	if s := p.Shadow; s != nil && s.Color != nil {
		k.shadow, k.shadowRGBA, k.shadowAt = true, rgba(s.Color), [2]float64{s.OffsetX, s.OffsetY}
	}
	if l := p.LinkColor; l != nil {
		k.linkIdle, k.linkHover = rgba(l.Idle), rgba(l.Hover)
	}
	return k
}

// drawCached draws the text from its image, its first line's top left at
// origin, drawing the text into the image first when what it shows has
// changed. It reports false for a text too large to keep as an image, which
// draws glyph by glyph instead.
func (t *Text) drawCached(screen *ebiten.Image, origin image.Point) bool {
	if t.cache == nil {
		t.cache = &textCache{}
	}
	c := t.cache
	w := t.widget.Rect.Dx()
	if k := t.cacheKey(w); !c.valid || k != c.key {
		c.key, c.valid = k, true
		c.tooLarge = !t.renderCache(w)
	}
	if c.tooLarge {
		return false
	}
	if c.image != nil {
		c.op.GeoM.Reset()
		c.op.GeoM.Translate(float64(origin.X+c.at.X), float64(origin.Y+c.at.Y))
		screen.DrawImage(c.image, &c.op)
	}
	return true
}

// renderCache draws the text, w pixels wide, into its image as draw would,
// the image sized to the glyphs' box; an image large enough already is
// cleared and kept. It reports false when the box is larger than
// textCacheMaxSide on a side.
func (t *Text) renderCache(w int) bool {
	c := t.cache
	lines := len(t.measurements.processedLines)
	shadow := t.computedParams.Shadow
	hasShadow := shadow != nil && shadow.Color != nil
	b := glyphBox{math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)}
	for i := range lines {
		if !t.measurements.processedLinesLaid[i] {
			t.layOutLine(i)
		}
		x, y := t.lineX(i, w), t.measurements.lineHeight*float64(i)
		t.addLineBox(&b, i, x, y)
		if hasShadow {
			t.addLineBox(&b, i, x+shadow.OffsetX, y+shadow.OffsetY)
		}
	}
	if b.minX > b.maxX { // no glyphs
		c.image = nil
		return true
	}

	x0, y0 := math.Floor(b.minX), math.Floor(b.minY)
	sw, sh := int(math.Ceil(b.maxX)-x0), int(math.Ceil(b.maxY)-y0)
	if sw > textCacheMaxSide || sh > textCacheMaxSide {
		return false
	}
	if c.image != nil && c.image.Bounds().Dx() >= sw && c.image.Bounds().Dy() >= sh {
		c.image.Clear()
	} else {
		if c.image != nil {
			c.image.Deallocate()
		}
		c.image = ebiten.NewImageWithOptions(image.Rect(0, 0, sw, sh), &ebiten.NewImageOptions{Unmanaged: true})
	}
	c.at = image.Pt(int(x0), int(y0))
	for i := range lines {
		x, y := t.lineX(i, w)-x0, t.measurements.lineHeight*float64(i)-y0
		if hasShadow {
			t.drawLine(c.image, i, x+shadow.OffsetX, y+shadow.OffsetY, shadow.Color)
		}
		t.drawLine(c.image, i, x, y, nil)
	}
	c.renders++
	return true
}

// addLineBox grows b to cover line i's laid out glyphs, the line starting
// at (x, y).
func (t *Text) addLineBox(b *glyphBox, i int, x, y float64) {
	if !t.ProcessBBCode {
		b.add(t.measurements.processedLineGlyphs[i], x, y)
		return
	}
	for _, piece := range t.measurements.processedLines[i] {
		b.add(piece.glyphs, x, y)
		x += piece.width
	}
}

// add grows b to cover glyphs drawn from (x, y).
func (b *glyphBox) add(glyphs []text.Glyph, x, y float64) {
	for _, g := range glyphs {
		s := g.Image.Bounds().Size()
		gx, gy := x+g.X, y+g.Y
		b.minX, b.minY = min(b.minX, gx), min(b.minY, gy)
		b.maxX, b.maxY = max(b.maxX, gx+float64(s.X)), max(b.maxY, gy+float64(s.Y))
	}
}
