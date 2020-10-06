package widget

import (
	"testing"

	"github.com/matryer/is"
)

func TestMultiOnce_Do(t *testing.T) {
	is := is.New(t)

	m := MultiOnce{}

	count := 0
	f := func() {
		count++
	}

	m.Append(f)
	m.Append(f)

	m.Do()
	is.Equal(count, 2)

	m.Do()
	is.Equal(count, 2)
}
