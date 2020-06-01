package widget

import "sync"

// MultiOnce works like sync.Once, but can execute any number of functions.
type MultiOnce struct {
	once  sync.Once
	funcs []func()
}

// Append adds f to the list of functions to be executed. If Do has been called already,
// calling Append will do nothing.
func (m *MultiOnce) Append(f func()) {
	m.funcs = append(m.funcs, f)
}

// Do executes all functions added using Append.
//
// Do executes the list of functions exactly once. Calling Do a second time will do nothing.
func (m *MultiOnce) Do() {
	m.once.Do(func() {
		defer func() {
			m.funcs = nil
		}()

		for _, f := range m.funcs {
			f()
		}
	})
}
