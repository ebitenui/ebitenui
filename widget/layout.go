package widget

import (
	"image"
)

type Layouter interface {
	PreferredSize(widgets []PreferredSizeLocateableWidget) (int, int)
	Layout(widgets []PreferredSizeLocateableWidget, rect image.Rectangle)
}

type Relayoutable interface {
	RequestRelayout()
}

type Locateable interface {
	SetLocation(rect image.Rectangle)
}

type Locater interface {
	WidgetAt(x int, y int) HasWidget
}

type Insets struct {
	Top    int
	Left   int
	Right  int
	Bottom int
}

type Direction int

const (
	DirectionHorizontal = Direction(iota)
	DirectionVertical
)

func NewInsetsSimple(widthHeight int) Insets {
	return Insets{
		Top:    widthHeight,
		Left:   widthHeight,
		Right:  widthHeight,
		Bottom: widthHeight,
	}
}

func (i Insets) Apply(rect image.Rectangle) image.Rectangle {
	rect.Min = rect.Min.Add(image.Point{i.Left, i.Top})
	rect.Max = rect.Max.Sub(image.Point{i.Right, i.Bottom})
	return rect
}

func (i Insets) Dx() int {
	return i.Left + i.Right
}

func (i Insets) Dy() int {
	return i.Top + i.Bottom
}
