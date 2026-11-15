package widget

import (
	"image"
	"runtime"
	"slices"
	"sync"
	"sync/atomic"

	"github.com/hajimehoshi/ebiten/v2"
)

// Cached texts live on shared pages: textures of their own, outside
// Ebitengine's atlases, each holding many texts on shelves. A texture costs
// the GPU at least 64 KiB however small (Direct3D), so a texture per text
// would cost far more than its pixels; one page holds a screenful of texts,
// and their draws share one source.

const (
	textPageSize    = 1024 // a page's side, in pixels
	textPageMaxSlot = 256  // the tallest slot a page gives; taller texts get a texture of their own
	textShelfStep   = 8    // shelf heights are multiples of this
)

// textSlot is a text's room on a page: rect, and a pixel's gap right and
// below it.
type textSlot struct {
	page  *textPage
	shelf int // the shelf's index on the page
	rect  image.Rectangle
}

// textPage is a page's texture and its shelves, top to bottom.
type textPage struct {
	image   *ebiten.Image
	shelves []textShelf
	used    int // slots given and not freed
}

// textShelf is a row of slots up to h pixels tall, gaps included, with the
// spans of x still free, in order and joined.
type textShelf struct {
	y, h int
	free []textSpan
}

type textSpan struct{ x0, x1 int }

// textPages is the pool of pages. Slots of texts the collector drops wait in
// freed, since cleanups run on another goroutine, until the next draw from a
// page or allocation frees them.
var textPages struct {
	pages   []*textPage
	mu      sync.Mutex
	freed   []*textSlot
	pending atomic.Bool
}

// allocTextSlot returns room for a w by h text on a page, making a page when
// none has room, or nil when the text is too tall or too wide for a page.
func allocTextSlot(w, h int) *textSlot {
	drainFreedTextSlots()
	if w+1 > textPageSize || h > textPageMaxSlot {
		return nil
	}
	for _, p := range textPages.pages {
		if s := p.alloc(w, h); s != nil {
			return s
		}
	}
	p := &textPage{image: ebiten.NewImageWithOptions(image.Rect(0, 0, textPageSize, textPageSize),
		&ebiten.NewImageOptions{Unmanaged: true})}
	textPages.pages = append(textPages.pages, p)
	return p.alloc(w, h)
}

// freeTextSlot gives s's room back, and releases its page when the page is
// left empty and another one remains.
func freeTextSlot(s *textSlot) {
	p := s.page
	p.free(s)
	if p.used == 0 && len(textPages.pages) > 1 {
		textPages.pages = slices.DeleteFunc(textPages.pages, func(q *textPage) bool { return q == p })
		p.image.Deallocate()
	}
}

// freeLater queues s to be freed on the drawing goroutine; it is the cleanup
// of a text the collector drops.
func freeLater(s *textSlot) {
	textPages.mu.Lock()
	textPages.freed = append(textPages.freed, s)
	textPages.mu.Unlock()
	textPages.pending.Store(true)
}

// drainFreedTextSlots frees the slots of texts the collector dropped.
func drainFreedTextSlots() {
	if !textPages.pending.Load() {
		return
	}
	textPages.mu.Lock()
	freed := textPages.freed
	textPages.freed = nil
	textPages.pending.Store(false)
	textPages.mu.Unlock()
	for _, s := range freed {
		freeTextSlot(s)
	}
}

// alloc gives room for a w by h text: on the tightest shelf at most half
// again as tall, else on a new shelf below the others, else on any shelf
// tall enough; nil when the page is full.
func (p *textPage) alloc(w, h int) *textSlot {
	gw, gh := w+1, h+1
	best := -1
	for i := range p.shelves {
		s := &p.shelves[i]
		if s.h >= gh && s.h <= gh+gh/2 && s.fits(gw) && (best < 0 || s.h < p.shelves[best].h) {
			best = i
		}
	}
	if best < 0 {
		y, sh := p.bottom(), (gh+textShelfStep-1)/textShelfStep*textShelfStep
		if y+sh <= textPageSize {
			p.shelves = append(p.shelves, textShelf{y: y, h: sh, free: []textSpan{{0, textPageSize}}})
			best = len(p.shelves) - 1
		}
	}
	if best < 0 {
		for i := range p.shelves {
			if p.shelves[i].h >= gh && p.shelves[i].fits(gw) {
				best = i
				break
			}
		}
	}
	if best < 0 {
		return nil
	}
	s := &p.shelves[best]
	x := s.take(gw)
	p.used++
	return &textSlot{page: p, shelf: best, rect: image.Rect(x, s.y, x+w, s.y+h)}
}

// free gives s's room back to its shelf, and drops empty shelves at the
// bottom so their rows serve any height.
func (p *textPage) free(s *textSlot) {
	p.shelves[s.shelf].give(textSpan{s.rect.Min.X, s.rect.Max.X + 1})
	p.used--
	for n := len(p.shelves); n > 0 && p.shelves[n-1].empty(); n-- {
		p.shelves = p.shelves[:n-1]
	}
}

// bottom is where a new shelf starts.
func (p *textPage) bottom() int {
	if len(p.shelves) == 0 {
		return 0
	}
	last := p.shelves[len(p.shelves)-1]
	return last.y + last.h
}

func (s *textShelf) fits(w int) bool {
	return slices.ContainsFunc(s.free, func(f textSpan) bool { return f.x1-f.x0 >= w })
}

// take cuts w pixels from the first free span long enough and returns
// where they start; the shelf must fit w.
func (s *textShelf) take(w int) int {
	i := slices.IndexFunc(s.free, func(f textSpan) bool { return f.x1-f.x0 >= w })
	x := s.free[i].x0
	if s.free[i].x0 += w; s.free[i].x0 == s.free[i].x1 {
		s.free = slices.Delete(s.free, i, i+1)
	}
	return x
}

// give returns span f to the shelf, joined with free spans beside it.
func (s *textShelf) give(f textSpan) {
	i, _ := slices.BinarySearchFunc(s.free, f.x0, func(g textSpan, x int) int { return g.x0 - x })
	s.free = slices.Insert(s.free, i, f)
	if i+1 < len(s.free) && s.free[i].x1 == s.free[i+1].x0 {
		s.free[i].x1 = s.free[i+1].x1
		s.free = slices.Delete(s.free, i+1, i+2)
	}
	if i > 0 && s.free[i-1].x1 == s.free[i].x0 {
		s.free[i-1].x1 = s.free[i].x1
		s.free = slices.Delete(s.free, i, i+1)
	}
}

func (s *textShelf) empty() bool {
	return len(s.free) == 1 && s.free[0] == textSpan{0, textPageSize}
}

// keepSlot gives the text a slot that frees itself when the collector drops
// the text.
func (t *Text) keepSlot(s *textSlot) {
	t.cache.slot = s
	t.cache.slotCleanup = runtime.AddCleanup(t, freeLater, s)
}

// dropSlot frees the text's slot now, if it has one.
func (t *Text) dropSlot() {
	c := t.cache
	if c.slot == nil {
		return
	}
	c.slotCleanup.Stop()
	freeTextSlot(c.slot)
	c.slot = nil
}
