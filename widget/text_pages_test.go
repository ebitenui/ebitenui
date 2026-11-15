package widget

import (
	img "image"
	"testing"

	"github.com/matryer/is"
)

// emptyTextPages gives a test a pool of text pages of its own.
func emptyTextPages(t *testing.T) {
	t.Helper()
	saved := textPages.pages
	textPages.pages = nil
	t.Cleanup(func() { textPages.pages = saved })
}

// Slots of one height share a shelf left to right, a pixel apart; a taller
// one opens a shelf below, and a full shelf sends a slot to a new one.
func TestTextPage_Shelves(t *testing.T) {
	is := is.New(t)

	p := &textPage{}
	a := p.alloc(300, 14)
	b := p.alloc(300, 14)
	tall := p.alloc(300, 30)
	is.Equal(a.rect, img.Rect(0, 0, 300, 14))
	is.Equal(b.rect, img.Rect(301, 0, 601, 14))   // beside a, past the gap
	is.Equal(tall.rect, img.Rect(0, 16, 300, 46)) // below the 16 px shelf

	c := p.alloc(300, 14)
	d := p.alloc(300, 14)                       // the 14 px shelf holds three
	is.Equal(c.rect, img.Rect(602, 0, 902, 14)) // third on the shelf
	is.Equal(d.rect, img.Rect(0, 48, 300, 62))  // a new shelf below the 32 px one
	is.Equal(p.used, 5)
}

// Of the shelves that fit a slot, the least tall takes it.
func TestTextPage_TightestShelf(t *testing.T) {
	is := is.New(t)

	p := &textPage{}
	p.alloc(500, 38)                          // a 40 px shelf
	is.Equal(p.alloc(600, 30).rect.Min.Y, 40) // too wide for it: a 32 px shelf below
	is.Equal(p.alloc(300, 30).rect.Min.Y, 40) // both fit; the 32 px one wins
	is.Equal(p.alloc(300, 30).rect.Min.Y, 0)  // the 32 px one is full
}

// A full page still gives a short slot room on a taller shelf.
func TestTextPage_AnyShelfWhenFull(t *testing.T) {
	is := is.New(t)

	p := &textPage{}
	first := p.alloc(textPageSize-1, 250)
	for range 3 {
		p.alloc(textPageSize-1, 250) // four 256 px shelves fill the page
	}
	is.True(p.alloc(100, 14) == nil)
	p.free(first)
	s := p.alloc(100, 14)
	is.True(s != nil)
	is.Equal(s.rect.Min.Y, 0) // on the emptied top shelf
}

// A freed slot's room goes to the next slot that fits, joined with free room
// beside it; empty shelves at the bottom are dropped, so their rows serve
// any height.
func TestTextPage_Free(t *testing.T) {
	is := is.New(t)

	p := &textPage{}
	p.alloc(300, 14)
	b := p.alloc(300, 14)
	c := p.alloc(300, 14)
	p.free(b)
	again := p.alloc(300, 14)
	is.Equal(again.rect, b.rect) // b's room again

	p.free(again)
	p.free(c)                                                   // joins the room on both sides
	is.Equal(p.alloc(722, 14).rect, img.Rect(301, 0, 1023, 14)) // b's, c's and the rest
	is.Equal(p.used, 2)

	low := p.alloc(300, 40)
	is.Equal(low.rect.Min.Y, 16)
	p.free(low)
	is.Equal(len(p.shelves), 1)               // the empty bottom shelf went
	is.Equal(p.alloc(300, 60).rect.Min.Y, 16) // its rows serve a taller slot
}

// Slots too tall or too wide for a page aren't given; a full page sends
// slots to a new page, and a page left empty beyond the first is released.
func TestTextPages_Pool(t *testing.T) {
	is := is.New(t)
	emptyTextPages(t)

	is.True(allocTextSlot(100, textPageMaxSlot+1) == nil) // too tall
	is.True(allocTextSlot(textPageSize, 14) == nil)       // too wide with its gap

	perPage := textPageSize / (textPageMaxSlot + 8) // 264 px shelves
	slots := make([]*textSlot, 0, perPage)
	for range perPage {
		slots = append(slots, allocTextSlot(textPageSize-1, textPageMaxSlot))
	}
	is.Equal(len(textPages.pages), 1)
	extra := allocTextSlot(textPageSize-1, textPageMaxSlot)
	is.Equal(len(textPages.pages), 2) // the first page is full
	is.True(extra.page != slots[0].page)

	freeTextSlot(extra)
	is.Equal(len(textPages.pages), 1) // released
	for _, s := range slots {
		freeTextSlot(s)
	}
	is.Equal(len(textPages.pages), 1) // the first page stays
}
