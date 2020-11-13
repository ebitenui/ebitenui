package event

import "testing"

type nopAction struct {
}

func BenchmarkExecuteDeferred(b *testing.B) {
	a := &nopAction{}
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			AddDeferred(a)
		}
		ExecuteDeferred()
	}
}

func (a *nopAction) Do() {
	// empty
}
