package event

import "testing"

type nopAction struct {
}

func BenchmarkExecuteDeferredActions(b *testing.B) {
	a := &nopAction{}
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			AddDeferred(a)
		}
		ExecuteDeferredActions()
	}
}

func (a *nopAction) Do() {
}
