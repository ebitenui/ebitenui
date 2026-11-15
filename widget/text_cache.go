package widget

import (
	"image"
	"image/color"
	"math"
	"runtime"

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
// instead of one per glyph. The images share 1024 by 1024 pages, textures
// of their own outside Ebitengine's atlases, so their draws batch and GPU
// memory grows a page at a time (4 MiB); a text taller than 256 pixels or
// wider than a page gets a texture of its own, and one over 2048 pixels on a
// side still draws glyph by glyph. It covers every Text, those inside
// buttons, labels and lists too. Set it before the interface is first drawn.
func SetTextCache(on bool) {
	textCacheOn = on
}

// textCache is a text drawn into an image, kept while what the text shows
// stays the same.
type textCache struct {
	image       *ebiten.Image   // the slot's part of its page, or own; nil while the text has no glyphs
	slot        *textSlot       // the text's room on a page
	slotCleanup runtime.Cleanup // frees slot when the collector drops the text
	own         *ebiten.Image   // a texture of its own, for a text too large for a page
	at          image.Point     // where the image's top left goes, from the text's origin
	key         textCacheKey
	valid       bool // key describes image
	tooLarge    bool // the text is drawn glyph by glyph
	renders     int  // times the text was drawn into the image
	op          ebiten.DrawImageOptions
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
	drainFreedTextSlots()
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
// the image at least the glyphs' box (see cacheTarget). It reports false
// when the box is larger than textCacheMaxSide on a side, giving its image
// up.
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
		t.dropCacheImage()
		return true
	}

	x0, y0 := math.Floor(b.minX), math.Floor(b.minY)
	sw, sh := int(math.Ceil(b.maxX)-x0), int(math.Ceil(b.maxY)-y0)
	if sw > textCacheMaxSide || sh > textCacheMaxSide {
		t.dropCacheImage()
		return false
	}
	dst := t.cacheTarget(sw, sh)
	c.at = image.Pt(int(x0), int(y0))
	// A slot's image keeps its page's coordinates: draw from its top left.
	ox, oy := float64(dst.Bounds().Min.X)-x0, float64(dst.Bounds().Min.Y)-y0
	for i := range lines {
		x, y := t.lineX(i, w)+ox, t.measurements.lineHeight*float64(i)+oy
		if hasShadow {
			t.drawLine(dst, i, x+shadow.OffsetX, y+shadow.OffsetY, shadow.Color)
		}
		t.drawLine(dst, i, x, y, nil)
	}
	c.renders++
	return true
}

// cacheTarget returns a cleared image at least sw by sh for the text to be
// drawn into: its slot while the box fits it and fills more than half of
// it, else a new slot, or a texture of its own when no slot can hold it.
func (t *Text) cacheTarget(sw, sh int) *ebiten.Image {
	c := t.cache
	if s := c.slot; s != nil && fitsWell(s.rect.Size(), sw, sh) {
		c.image.Clear()
		return c.image
	}
	t.dropSlot()
	if s := allocTextSlot(sw, sh); s != nil {
		img, ok := s.page.image.SubImage(s.rect).(*ebiten.Image)
		if ok {
			t.keepSlot(s)
			if c.own != nil {
				c.own.Deallocate()
				c.own = nil
			}
			c.image = img
			c.image.Clear()
			return c.image
		}
		freeTextSlot(s)
	}
	if c.own == nil || !fitsWell(c.own.Bounds().Size(), sw, sh) {
		if c.own != nil {
			c.own.Deallocate()
		}
		c.own = ebiten.NewImageWithOptions(image.Rect(0, 0, sw, sh), &ebiten.NewImageOptions{Unmanaged: true})
	} else {
		c.own.Clear()
	}
	c.image = c.own
	return c.image
}

// fitsWell reports whether an image of size holds a w by h box without
// being more than twice its width or height (a few pixels of slack aside).
func fitsWell(size image.Point, w, h int) bool {
	return size.X >= w && size.Y >= h && size.X <= 2*w+16 && size.Y <= 2*h+16
}

// dropCacheImage gives up the text's slot or texture.
func (t *Text) dropCacheImage() {
	c := t.cache
	t.dropSlot()
	if c.own != nil {
		c.own.Deallocate()
		c.own = nil
	}
	c.image = nil
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
