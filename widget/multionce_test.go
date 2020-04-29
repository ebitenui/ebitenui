package widget

import (
	"testing"

	"github.com/matryer/is"
)

func TestMultiOnce_Do(t *testing.T) {
	is := is.New(t)

	count := 0

	m := MultiOnce{}
	m.Append(func() {
		count++
	})
	m.Append(func() {
		count++
	})

	m.Do()
	m.Do()
	is.Equal(count, 2)
}
