package widget

import (
	"image"
	"math"

	"github.com/ebitenui/ebitenui/event"
)

// WindowTransition makes a window ease into place when it is added to a
// UI, and back out when it closes, before the UI removes it.
type WindowTransition struct {
	// Duration is how many updates the transition takes.
	Duration int
	// Offset is where the window comes from and goes back to, relative to
	// its place, in pixels.
	Offset image.Point
	// Fade fades the window in and out as it moves, through its opacity.
	Fade bool
	// Ease maps the transition's progress, from 0 to 1, to how far the
	// window has come; nil eases out cubically. Closing plays it backward.
	Ease func(t float64) float64
}

// Transition makes the window ease in when it is added to a UI and out when
// it closes: Close, or the function AddWindow returns, starts the way out,
// and the UI removes the window, firing ClosedEvent, once it is over. A
// transition of no duration does nothing.
func (o WindowOptions) Transition(t WindowTransition) WindowOpt {
	return func(w *Window) {
		if t.Duration <= 0 {
			return
		}
		w.transition = &t
		w.transitionEnded = &event.Event{}
		w.transitionEnded.AddHandler(func(any) {
			if w.closing && w.progress == 0 && w.remove != nil {
				remove := w.remove
				w.remove = nil
				remove()
			}
		})
	}
}

// Closing reports whether the window is on its way out: its transition
// runs backward, and the UI removes it at the end.
func (w *Window) Closing() bool { return w.closing }

// Reopen turns a window on its way out back toward its place. AddWindow
// does it for a window still closing.
func (w *Window) Reopen() { w.closing = false }

// startTransition puts the window at the start of its way in, as it is
// added to a UI that removes it with remove.
func (w *Window) startTransition(remove func()) {
	w.progress, w.closing, w.remove = 0, false, remove
	w.applyTransition()
}

// updateTransition moves the window a step along its transition. Once its
// way out is over, it asks to be removed when the update's events run, so
// that the UI's loop over its windows is not disturbed.
func (w *Window) updateTransition() {
	switch {
	case w.closing && w.progress > 0:
		w.progress--
	case !w.closing && w.progress < w.transition.Duration:
		w.progress++
	default:
		return
	}
	w.applyTransition()
	if w.closing && w.progress == 0 {
		w.transitionEnded.Fire(nil)
	}
}

// applyTransition shows the window where its transition has it: off its
// place by the part of the offset still to go, and faded as far.
func (w *Window) applyTransition() {
	t := w.transition
	ease := t.Ease
	if ease == nil {
		ease = easeOutCubic
	}
	v := ease(float64(w.progress) / float64(t.Duration))
	off := image.Pt(int(math.Round((1-v)*float64(t.Offset.X))), int(math.Round((1-v)*float64(t.Offset.Y))))
	if off != w.shownAt {
		place := w.placeRect()
		w.shownAt = off
		w.container.SetLocation(place.Add(off))
	}
	if t.Fade {
		w.SetOpacity(float32(v))
	}
}

// placeRect is where the window is placed: where it shows, less the offset
// its transition shows it at.
func (w *Window) placeRect() image.Rectangle {
	return w.container.GetWidget().Rect.Sub(w.shownAt)
}

func easeOutCubic(t float64) float64 {
	u := 1 - t
	return 1 - u*u*u
}
