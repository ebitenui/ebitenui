package widget

import (
	"image"
	"math"

	"github.com/ebitenui/ebitenui/event"
)

// WindowTransition makes a window ease into place when it is added to a
// UI, and out when it closes, before the UI removes it: in by its In
// motion, out by its Out motion.
type WindowTransition struct {
	// Duration is how many updates each way takes.
	Duration int
	// In is how the window comes in, Out how it goes out.
	In, Out WindowMotion
	// Ease maps the transition's progress, from 0 to 1, to how far the
	// window has come; nil eases out cubically. Closing plays it backward.
	Ease func(t float64) float64
}

// WindowMotion is one way of a window's transition: from or to an offset
// from its place, fading or not. The zero motion keeps the window in place,
// opaque: it shows at once, or goes at once.
type WindowMotion struct {
	// Offset is where the window comes from, or goes to, relative to its
	// place, in pixels: (0, 40) from below, (-40, 0) from the left.
	Offset image.Point
	// Fade fades the window in, or out, through its opacity.
	Fade bool
}

// Transition makes the window ease in when it is added to a UI and out when
// it closes: Close, or the function AddWindow returns, starts the way out,
// and the UI removes the window, firing ClosedEvent, once it is over. A
// window turned back on its way takes the other motion from where it is.
// A transition of no duration does nothing.
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
func (w *Window) Reopen() {
	if w.closing {
		w.closing = false
		w.applyTransition()
	}
}

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

// applyTransition shows the window where its transition has it, by the
// motion of the way it goes: off its place by the part of the motion's
// offset still to go, and faded as far if the motion fades.
func (w *Window) applyTransition() {
	t := w.transition
	m := t.In
	if w.closing {
		m = t.Out
	}
	ease := t.Ease
	if ease == nil {
		ease = easeOutCubic
	}
	v := ease(float64(w.progress) / float64(t.Duration))
	off := image.Pt(int(math.Round((1-v)*float64(m.Offset.X))), int(math.Round((1-v)*float64(m.Offset.Y))))
	if off != w.shownAt {
		place := w.placeRect()
		w.shownAt = off
		w.container.SetLocation(place.Add(off))
	}
	opacity := float32(1)
	if m.Fade {
		opacity = float32(v)
	}
	if c, ok := w.container.(*Container); ok && c.Opacity() != opacity {
		c.SetOpacity(opacity)
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
