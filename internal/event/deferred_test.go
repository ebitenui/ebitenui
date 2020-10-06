package event

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

type deferredActionMock struct {
	mock.Mock
}

func TestAddDeferred(t *testing.T) {
	a1 := &deferredActionMock{}
	a2 := &deferredActionMock{}

	a1.On("Do").Run(func(args mock.Arguments) {
		AddDeferred(a2)
	})
	a2.On("Do")

	AddDeferred(a1)
	ExecuteDeferred()

	a1.AssertExpectations(t)
	a2.AssertExpectations(t)
}

func (d *deferredActionMock) Do() {
	d.Called()
}
