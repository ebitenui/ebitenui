package widget

import (
	"slices"
	"testing"

	"github.com/matryer/is"
)

// ownContainer is a caller's container type, embedding a Container.
type ownContainer struct{ *Container }

// The focus search goes into every container: a Panel, or a type of the
// caller's that embeds a Container, as into the types it knows by name; and
// into a tab book's open tab.
func TestContainer_GetFocusersSearchesEveryContainer(t *testing.T) {
	is := is.New(t)

	panel := NewPanel(ContainerOpts.BackgroundImage(newNineSliceEmpty(t)))
	own := ownContainer{NewContainer()}
	for _, c := range []struct {
		inner PreferredSizeLocateableWidget
		add   func(...PreferredSizeLocateableWidget) RemoveChildFunc
	}{{panel, panel.AddChild}, {own, own.AddChild}} {
		b := newButton(t)
		c.add(b)
		outer := NewContainer()
		outer.AddChild(c.inner)
		found := outer.GetFocusers()
		is.Equal(len(found), 1)
		is.True(found[0] == Focuser(b))
	}

	b := newButton(t)
	tab := NewTabBookTab(TabBookTabOpts.Label("Tab"))
	tab.AddChild(b)
	outer := NewContainer()
	outer.AddChild(newTabBook(t, TabBookOpts.Tabs(tab)))
	is.True(slices.Contains(outer.GetFocusers(), Focuser(b))) // in the open tab
}
