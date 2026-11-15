package ebitenui

import (
	"image"

	"github.com/ebitenui/ebitenui/widget"
)

// nearestFocuser returns the candidate nearest to from in direction d, or
// nil when none lies that way. For the four sides it follows the rules of
// Android's focus search: a candidate must lie beyond from's edge in that
// direction; one overlapping from's row (or column, going up or down) wins
// over one beside it; among those, the nearest wins, by 13 times the square
// of the distance along the direction plus the square of the offset across
// it, so a candidate straight ahead beats a nearer one off to the side. For
// the diagonals it takes the nearest candidate whose center lies in that
// quarter. Ties go to the earlier candidate.
func nearestFocuser(from widget.Focuser, d widget.FocusDirection, candidates []widget.Focuser) widget.Focuser {
	src := from.GetWidget().Rect
	var best widget.Focuser
	bestInBeam, bestScore := false, 0
	for _, c := range candidates {
		if c == from {
			continue
		}
		dst := c.GetWidget().Rect
		ok, inBeam, score := focusScore(src, dst, d)
		if !ok {
			continue
		}
		if best == nil || (inBeam && !bestInBeam) || (inBeam == bestInBeam && score < bestScore) {
			best, bestInBeam, bestScore = c, inBeam, score
		}
	}
	return best
}

// focusScore reports whether dst lies in direction d from src, whether it
// overlaps src's row or column across d, and its distance score, lower
// being nearer.
func focusScore(src, dst image.Rectangle, d widget.FocusDirection) (ok, inBeam bool, score int) {
	sc, dc := center(src), center(dst)
	dx, dy := dc.X-sc.X, dc.Y-sc.Y
	rows := dst.Min.Y < src.Max.Y && dst.Max.Y > src.Min.Y
	cols := dst.Min.X < src.Max.X && dst.Max.X > src.Min.X
	var major, minor int
	switch d {
	case widget.FOCUS_EAST:
		ok = (src.Min.X < dst.Min.X || src.Max.X <= dst.Min.X) && src.Max.X < dst.Max.X
		major, minor, inBeam = max(0, dst.Min.X-src.Max.X), dy, rows
	case widget.FOCUS_WEST:
		ok = (src.Max.X > dst.Max.X || src.Min.X >= dst.Max.X) && src.Min.X > dst.Min.X
		major, minor, inBeam = max(0, src.Min.X-dst.Max.X), dy, rows
	case widget.FOCUS_SOUTH:
		ok = (src.Min.Y < dst.Min.Y || src.Max.Y <= dst.Min.Y) && src.Max.Y < dst.Max.Y
		major, minor, inBeam = max(0, dst.Min.Y-src.Max.Y), dx, cols
	case widget.FOCUS_NORTH:
		ok = (src.Max.Y > dst.Max.Y || src.Min.Y >= dst.Max.Y) && src.Min.Y > dst.Min.Y
		major, minor, inBeam = max(0, src.Min.Y-dst.Max.Y), dx, cols
	case widget.FOCUS_NORTHEAST:
		return dx > 0 && dy < 0, false, dx*dx + dy*dy
	case widget.FOCUS_SOUTHEAST:
		return dx > 0 && dy > 0, false, dx*dx + dy*dy
	case widget.FOCUS_SOUTHWEST:
		return dx < 0 && dy > 0, false, dx*dx + dy*dy
	case widget.FOCUS_NORTHWEST:
		return dx < 0 && dy < 0, false, dx*dx + dy*dy
	case widget.FOCUS_NEXT, widget.FOCUS_PREVIOUS:
		return false, false, 0 // the tab order, not a place
	}
	return ok, inBeam, 13*major*major + minor*minor
}

func center(r image.Rectangle) image.Point {
	return image.Pt((r.Min.X+r.Max.X)/2, (r.Min.Y+r.Max.Y)/2)
}
