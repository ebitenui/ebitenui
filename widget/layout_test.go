package widget

import (
	"image"
	"testing"

	"github.com/matryer/is"
)

func TestNewInsetsSimple(t *testing.T) {
	is := is.New(t)

	i := NewInsetsSimple(15)

	is.Equal(i.Left, 15)
	is.Equal(i.Right, 15)
	is.Equal(i.Top, 15)
	is.Equal(i.Bottom, 15)
}

func TestInsets_Apply(t *testing.T) {
	is := is.New(t)

	i := Insets{
		Left:   10,
		Right:  20,
		Top:    30,
		Bottom: 40,
	}
	r := image.Rect(25, 35, 145, 155)

	is.Equal(i.Apply(r), image.Rect(35, 65, 125, 115))
}

func TestInsets_Dx(t *testing.T) {
	is := is.New(t)

	i := Insets{
		Left:  10,
		Right: 20,
	}

	is.Equal(i.Dx(), 30)
}

func TestInsets_Dy(t *testing.T) {
	is := is.New(t)

	i := Insets{
		Top:    30,
		Bottom: 40,
	}

	is.Equal(i.Dy(), 70)
}
