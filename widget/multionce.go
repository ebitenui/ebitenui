package widget

import "sync"

type MultiOnce struct {
	once  sync.Once
	funcs []func()
}

func (m *MultiOnce) Append(f func()) {
	m.funcs = append(m.funcs, f)
}

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
