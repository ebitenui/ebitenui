package widget

import (
	"image/color"
	"testing"

	"github.com/ebitenui/ebitenui/event"
	"github.com/matryer/is"
)

func TestTabBook_Tab_Initial(t *testing.T) {
	is := is.New(t)

	tab1 := NewTabBookTab("Tab 1", newTabContainerOpts()...)
	tab2 := NewTabBookTab("Tab 2", newTabContainerOpts()...)

	tb := newTabBook(t,
		TabBookOpts.Tabs(tab1, tab2),
		TabBookOpts.TabSelectedHandler(func(_ *TabBookTabSelectedEventArgs) {
			is.Fail() // event fired without previous action
		}))

	is.Equal(tb.Tab(), tab1)
}

func TestTabBook_SetTab(t *testing.T) {
	is := is.New(t)

	var eventArgs *TabBookTabSelectedEventArgs
	numEvents := 0

	tab1 := NewTabBookTab("Tab 1", newTabContainerOpts()...)
	tab2 := NewTabBookTab("Tab 2", newTabContainerOpts()...)

	tb := newTabBook(t,
		TabBookOpts.Tabs(tab1, tab2),
		TabBookOpts.TabSelectedHandler(func(args *TabBookTabSelectedEventArgs) {
			eventArgs = args
			numEvents++
		}))

	tb.SetTab(tab2)
	event.ExecuteDeferred()

	is.Equal(tb.Tab(), tab2)
	is.Equal(eventArgs.Tab, tab2)
	is.Equal(eventArgs.PreviousTab, tab1)

	tb.SetTab(tab2)
	event.ExecuteDeferred()
	is.Equal(numEvents, 1)
}

func TestTabBook_TabSelectedEvent_User(t *testing.T) {
	is := is.New(t)

	var eventArgs *TabBookTabSelectedEventArgs
	numEvents := 0

	tab1 := NewTabBookTab("Tab 1", newTabContainerOpts()...)
	tab2 := NewTabBookTab("Tab 2", newTabContainerOpts()...)

	tb := newTabBook(t,
		TabBookOpts.Tabs(tab1, tab2),
		TabBookOpts.TabSelectedHandler(func(args *TabBookTabSelectedEventArgs) {
			eventArgs = args
			numEvents++
		}))

	leftMouseButtonClick(tabBookButtons(tb)[1], t)

	is.Equal(tb.Tab(), tab2)
	is.Equal(eventArgs.Tab, tab2)
	is.Equal(eventArgs.PreviousTab, tab1)

	leftMouseButtonClick(tabBookButtons(tb)[1], t)
	is.Equal(numEvents, 1)
}

func newTabBook(t *testing.T, opts ...TabBookOpt) *TabBook {
	t.Helper()

	tb := NewTabBook(append(opts, []TabBookOpt{
		TabBookOpts.TabButtonImage(&ButtonImage{
			Idle: newNineSliceEmpty(t),
		}),
		TabBookOpts.TabButtonText(loadFont(t), &ButtonTextColor{
			Idle:     color.Transparent,
			Disabled: color.Transparent,
		}),
	}...)...)

	event.ExecuteDeferred()
	render(tb, t)
	return tb
}

func tabBookButtons(t *TabBook) []*Button {
	buttons := []*Button{}
	for _, tab := range t.tabs {
		buttons = append(buttons, t.tabToButton[tab])
	}
	return buttons
}

func newTabContainerOpts() []ContainerOpt {
	result := []ContainerOpt{}

	return result
}
