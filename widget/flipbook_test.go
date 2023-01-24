package widget

import (
	"testing"

	"github.com/ebitenui/ebitenui/event"
	"github.com/matryer/is"
	"github.com/stretchr/testify/mock"
)

func TestFlipBook_SetPage_AlwaysRenderSingleWidget(t *testing.T) {
	is := is.New(t)

	f := newFlipBook(t)

	pages := []*controlMock{}
	order := []*controlMock{}

	for i := 0; i < 3; i++ {
		w := NewWidget()
		p := controlMock{}
		p.On("GetWidget").Maybe().Return(w)
		p.On("PreferredSize").Maybe().Return(50, 50)
		p.On("SetLocation", mock.Anything).Maybe()
		p.On("Render", mock.Anything, mock.Anything).Maybe().Run(func(_ mock.Arguments) {
			order = append(order, &p)
		})
		pages = append(pages, &p)
	}

	expectedOrder := []*controlMock{pages[0], pages[1], pages[2], pages[0]}

	for _, p := range expectedOrder {
		f.SetPage(p)
		render(f, t)
	}

	pages[0].AssertNumberOfCalls(t, "Render", 2)
	pages[1].AssertNumberOfCalls(t, "Render", 1)
	pages[2].AssertNumberOfCalls(t, "Render", 1)

	is.Equal(order, expectedOrder)
}

func newFlipBook(t *testing.T, opts ...FlipBookOpt) *FlipBook {
	t.Helper()

	f := NewFlipBook(opts...)
	event.ExecuteDeferred()
	render(f, t)
	return f
}
