package widget

import (
	img "image"
	"image/color"
	"strconv"
	"testing"

	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/image"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/matryer/is"
)

func TestList_SelectedEntry_Initial(t *testing.T) {
	is := is.New(t)

	entries := []interface{}{"first", "second", "third"}

	list := newList(t,
		ListOpts.Entries(entries),

		ListOpts.EntryLabelFunc(func(e interface{}) string {
			result, _ := e.(string)
			return result
		}),

		ListOpts.EntrySelectedHandler(func(_ *ListEntrySelectedEventArgs) {
			is.Fail() // event fired without previous action
		}))

	is.Equal(list.SelectedEntry(), nil)
}

func TestList_NoSliderOpts(t *testing.T) {
	entries := []interface{}{"first", "second", "third"}
	_ = NewList(
		ListOpts.Entries(entries),

		ListOpts.EntryLabelFunc(func(e interface{}) string {
			result, _ := e.(string)
			return result
		}),

		ListOpts.EntrySelectedHandler(func(_ *ListEntrySelectedEventArgs) {
		}),
		ListOpts.ScrollContainerImage(&ScrollContainerImage{
			Idle:     newNineSliceEmpty(t),
			Disabled: newNineSliceEmpty(t),
			Mask:     newNineSliceEmpty(t),
		}),

		ListOpts.HideHorizontalSlider(),
		ListOpts.HideVerticalSlider(),

		ListOpts.EntryFontFace(loadFont(t)),

		ListOpts.EntryColor(&ListEntryColor{
			Unselected:                 color.Transparent,
			Selected:                   color.Transparent,
			DisabledUnselected:         color.Transparent,
			DisabledSelected:           color.Transparent,
			SelectedBackground:         color.Transparent,
			DisabledSelectedBackground: color.Transparent,
		}),
	)
}

func TestList_SetSelectedEntry(t *testing.T) {
	is := is.New(t)

	entries := []interface{}{"first", "second", "third"}

	var eventArgs *ListEntrySelectedEventArgs
	numEvents := 0

	list := newList(t,
		ListOpts.Entries(entries),

		ListOpts.EntryLabelFunc(func(e interface{}) string {
			result, _ := e.(string)
			return result
		}),

		ListOpts.EntrySelectedHandler(func(args *ListEntrySelectedEventArgs) {
			eventArgs = args
			numEvents++
		}))

	list.SetSelectedEntry(entries[1])
	event.ExecuteDeferred()

	is.Equal(eventArgs.Entry, entries[1])
	is.Equal(list.SelectedEntry(), entries[1])

	list.SetSelectedEntry(entries[1])
	event.ExecuteDeferred()

	is.Equal(numEvents, 1)
}

func TestList_UpdateEntry(t *testing.T) {
	is := is.New(t)

	type listEntry struct {
		label string
	}

	entries := []interface{}{
		&listEntry{label: "first"},
		&listEntry{label: "second"},
		&listEntry{label: "third"},
	}

	list := newList(t,
		ListOpts.Entries(entries),

		ListOpts.EntryLabelFunc(func(e interface{}) string {
			return e.(*listEntry).label
		}),
	)

	oldButton := list.buttons[1]
	entry := entries[1].(*listEntry)

	is.Equal(list.buttons[1].Text().Label, "second")

	entry.label = "updated"

	is.Equal(list.buttons[1].Text().Label, "second")

	list.UpdateEntry(entry)

	is.True(list.buttons[1] != oldButton)
	is.Equal(list.buttons[1].Text().Label, "updated")
}

func TestList_EntrySelectedEvent_IncludesList(t *testing.T) {
	is := is.New(t)

	entries := []interface{}{"first", "second", "third"}

	var eventArgs *ListEntrySelectedEventArgs
	var selectedEntry interface{}

	list := newList(t,
		ListOpts.Entries(entries),

		ListOpts.EntryLabelFunc(func(e interface{}) string {
			result, _ := e.(string)
			return result
		}),

		ListOpts.EntrySelectedHandler(func(args *ListEntrySelectedEventArgs) {
			eventArgs = args
			selectedEntry = args.List.SelectedEntry()
		}))

	list.SetSelectedEntry(entries[1])
	event.ExecuteDeferred()

	is.Equal(eventArgs.List, list)
	is.Equal(selectedEntry, entries[1])
}

func TestList_EntrySelectedEvent_User(t *testing.T) {
	is := is.New(t)

	entries := []interface{}{"first", "second", "third"}

	var eventArgs *ListEntrySelectedEventArgs
	numEvents := 0

	list := newList(t,
		ListOpts.Entries(entries),

		ListOpts.EntryLabelFunc(func(e interface{}) string {
			result, _ := e.(string)
			return result
		}),

		ListOpts.EntrySelectedHandler(func(args *ListEntrySelectedEventArgs) {
			eventArgs = args
			numEvents++
		}))

	leftMouseButtonClick(list.buttons[1], t)

	is.Equal(eventArgs.Entry, entries[1])
	is.Equal(list.SelectedEntry(), entries[1])

	leftMouseButtonClick(list.buttons[1], t)

	is.Equal(numEvents, 1)
}

func TestList_EntrySelectedEvent_User_AllowReselect(t *testing.T) {
	is := is.New(t)

	entries := []interface{}{"first", "second", "third"}

	var eventArgs *ListEntrySelectedEventArgs
	numEvents := 0

	list := newList(t,
		ListOpts.Entries(entries),

		ListOpts.EntryLabelFunc(func(e interface{}) string {
			result, _ := e.(string)
			return result
		}),

		ListOpts.AllowReselect(),

		ListOpts.EntrySelectedHandler(func(args *ListEntrySelectedEventArgs) {
			eventArgs = args
			numEvents++
		}))

	leftMouseButtonClick(list.buttons[1], t)

	is.Equal(eventArgs.Entry, entries[1])
	is.Equal(list.SelectedEntry(), entries[1])

	leftMouseButtonClick(list.buttons[1], t)

	is.Equal(eventArgs.Entry, entries[1])
	is.Equal(eventArgs.PreviousEntry, entries[1])
	is.Equal(list.SelectedEntry(), entries[1])

	is.Equal(numEvents, 2)
}

func TestList_Entries_SortedOnCreate(t *testing.T) {
	is := is.New(t)

	list := newList(t,
		ListOpts.Entries([]interface{}{"third", "first", "second"}),
		ListOpts.EntrySortFunc(func(a, b any) int {
			left := a.(string)
			right := b.(string)
			switch {
			case left < right:
				return -1
			case left > right:
				return 1
			default:
				return 0
			}
		}),
		ListOpts.EntryLabelFunc(func(e interface{}) string {
			return e.(string)
		}),
	)

	is.Equal(list.Entries()[0], "first")
	is.Equal(list.Entries()[1], "second")
	is.Equal(list.Entries()[2], "third")
	is.Equal(list.buttons[0].Text().Label, "first")
	is.Equal(list.buttons[1].Text().Label, "second")
	is.Equal(list.buttons[2].Text().Label, "third")
}

func TestList_AddEntry_ResortsAndReusesButtons(t *testing.T) {
	is := is.New(t)

	list := newList(t,
		ListOpts.Entries([]interface{}{2, 4}),
		ListOpts.EntrySortFunc(func(a, b any) int {
			return a.(int) - b.(int)
		}),
		ListOpts.EntryLabelFunc(func(e interface{}) string {
			return strconv.Itoa(e.(int))
		}),
	)

	button2 := list.buttons[0]
	button4 := list.buttons[1]

	list.AddEntry(3)

	is.Equal(list.Entries()[0], 2)
	is.Equal(list.Entries()[1], 3)
	is.Equal(list.Entries()[2], 4)
	is.Equal(list.buttons[0], button2)
	is.Equal(list.buttons[2], button4)
	is.Equal(list.buttons[1].Text().Label, "3")

	children := list.listContent.Children()
	is.Equal(children[0], button2)
	is.Equal(children[1], list.buttons[1])
	is.Equal(children[2], button4)
}

func newList(t *testing.T, opts ...ListOpt) *List {
	t.Helper()

	l := NewList(append(opts, []ListOpt{
		ListOpts.ScrollContainerImage(&ScrollContainerImage{
			Idle:     newNineSliceEmpty(t),
			Disabled: newNineSliceEmpty(t),
			Mask:     newNineSliceEmpty(t),
		}),

		ListOpts.SliderParams(&SliderParams{
			TrackImage: &SliderTrackImage{},
			HandleImage: &ButtonImage{
				Idle:    newNineSliceEmpty(t),
				Pressed: newNineSliceEmpty(t),
			}}),

		ListOpts.EntryFontFace(loadFont(t)),

		ListOpts.EntryColor(&ListEntryColor{
			Unselected:                 color.Transparent,
			Selected:                   color.Transparent,
			DisabledUnselected:         color.Transparent,
			DisabledSelected:           color.Transparent,
			SelectedBackground:         color.Transparent,
			DisabledSelectedBackground: color.Transparent,
		}),
	}...)...)

	event.ExecuteDeferred()
	render(l, t)
	return l
}

func listEntryButtons(l *List) []*Button {
	return l.buttons
}

// Entries can be widgets, built per entry, with neither a label function nor
// a font: each shows inside its entry's button, clicking the button selects
// the entry, and UpdateEntry builds the widget again.
func TestList_EntryWidgets(t *testing.T) {
	is := is.New(t)

	built := map[any]int{}
	widgets := map[any]*simpleWidget{}
	var selected any
	l := NewList(
		ListOpts.Entries([]any{"a", "b"}),
		ListOpts.EntryWidgetFunc(func(e any) PreferredSizeLocateableWidget {
			built[e]++
			widgets[e] = newSimpleWidget(30, 15, nil)
			return widgets[e]
		}),
		ListOpts.EntrySelectedHandler(func(args *ListEntrySelectedEventArgs) { selected = args.Entry }),
		ListOpts.ScrollContainerImage(&ScrollContainerImage{Idle: newNineSliceEmpty(t), Mask: newNineSliceEmpty(t)}),
		ListOpts.SliderParams(&SliderParams{TrackImage: &SliderTrackImage{}, HandleImage: &ButtonImage{Idle: newNineSliceEmpty(t), Pressed: newNineSliceEmpty(t)}}),
		ListOpts.EntryColor(&ListEntryColor{
			Unselected: color.Transparent, Selected: color.Transparent, DisabledUnselected: color.Transparent,
			DisabledSelected: color.Transparent, SelectedBackground: color.Transparent, DisabledSelectedBackground: color.Transparent,
		}),
	)
	l.SetLocation(img.Rect(0, 0, 100, 100))
	event.ExecuteDeferred()
	render(l, t)

	buttons := listEntryButtons(l)
	is.Equal(len(buttons), 2)
	is.True(!widgets["b"].GetWidget().Rect.Empty() && widgets["b"].GetWidget().Rect.In(buttons[1].GetWidget().Rect)) // inside its button

	leftMouseButtonClick(buttons[1], t)
	is.Equal(selected, "b") // selected by clicking

	l.UpdateEntry("b")
	is.Equal(built["b"], 2) // built again
}

// newEntryImages returns entry images, each its own.
func newEntryImages(t *testing.T) *ListEntryImage {
	t.Helper()
	return &ListEntryImage{
		Unselected: newNineSliceEmpty(t), DisabledUnselected: newNineSliceEmpty(t),
		Focused: newNineSliceEmpty(t), Selecting: newNineSliceEmpty(t),
		Selected: newNineSliceEmpty(t), DisabledSelected: newNineSliceEmpty(t),
		SelectedFocused: newNineSliceEmpty(t), SelectingFocused: newNineSliceEmpty(t),
	}
}

// isEntryImages checks that a list draws its entries with images, each state
// its own.
func isEntryImages(is *is.I, l *List, images *ListEntryImage) {
	is.Helper()
	unselected, selected := l.computedParams.entryUnselectedColor, l.computedParams.entrySelectedColor
	is.True(unselected.Idle == images.Unselected)
	is.True(unselected.Disabled == images.DisabledUnselected)
	is.True(unselected.Hover == images.Focused)
	is.True(unselected.Pressed == images.Selecting)
	is.True(selected.Idle == images.Selected)
	is.True(selected.Disabled == images.DisabledSelected)
	is.True(selected.Hover == images.SelectedFocused)
	is.True(selected.Pressed == images.SelectingFocused)
}

// Entry images draw the entries' backgrounds, rounded or bordered, in place
// of their colors, each state its own: the selected entry's button takes the
// selected ones. A state without an image keeps its color, and an
// unselected entry without one stays clear.
func TestList_EntryImages(t *testing.T) {
	is := is.New(t)

	images := newEntryImages(t)
	l := newList(t, ListOpts.EntryImage(images), ListOpts.Entries([]any{"a", "b"}),
		ListOpts.EntryLabelFunc(func(e any) string { return e.(string) })) //nolint:forcetypeassert // strings
	isEntryImages(is, l, images)
	l.SetSelectedEntry("b")
	is.True(listEntryButtons(l)[1].computedParams.Image.Idle == images.Selected)   // the selected entry
	is.True(listEntryButtons(l)[0].computedParams.Image.Idle == images.Unselected) // the others

	focused := color.NRGBA{1, 2, 3, 255}
	partial := newList(t, ListOpts.EntryImage(&ListEntryImage{Selected: images.Selected}), ListOpts.Entries([]any{"a"}),
		ListOpts.EntryLabelFunc(func(e any) string { return e.(string) })) //nolint:forcetypeassert // strings
	colors := *partial.definedParams.EntryColor
	colors.SelectedFocusedBackground = focused
	partial.definedParams.EntryColor = &colors
	partial.Validate()
	is.True(partial.computedParams.entrySelectedColor.Idle == images.Selected)
	is.True(partial.computedParams.entrySelectedColor.Hover == image.NewNineSliceColor(focused)) // its color
	is.True(partial.computedParams.entryUnselectedColor.Idle == image.NewNineSliceColor(color.Transparent))
}

// A theme's entry images reach a list, and its combo button entry images a
// list combo button's list, unless images are given.
func TestList_EntryImagesFromTheTheme(t *testing.T) {
	is := is.New(t)

	themed, comboThemed := newEntryImages(t), newEntryImages(t)
	theme := &Theme{ListTheme: &ListParams{EntryImage: themed}, ListComboButtonTheme: &ListComboButtonParams{List: &ListParams{EntryImage: comboThemed}}}
	label := func(e any) string { return e.(string) } //nolint:forcetypeassert // strings

	l := newList(t, ListOpts.Entries([]any{"a"}), ListOpts.EntryLabelFunc(label))
	l.GetWidget().SetTheme(theme)
	l.Validate()
	isEntryImages(is, l, themed)

	given := newEntryImages(t)
	l.definedParams.EntryImage = given
	l.Validate()
	isEntryImages(is, l, given)

	c := newListComboButton(t, ListComboButtonOpts.Entries([]any{"a"}), ListComboButtonOpts.EntryLabelFunc(label, label))
	c.GetWidget().SetTheme(theme)
	c.Validate()
	isEntryImages(is, c.list, comboThemed) // a combo's list takes the combo's

	c.definedParams.List.EntryImage = given
	c.Validate()
	isEntryImages(is, c.list, given)
}

// With SlidersWhenNeeded, the vertical slider shows only while the entries
// overflow the list, and keeps its room while hidden.
func TestList_SlidersWhenNeeded(t *testing.T) {
	is := is.New(t)

	l := newList(t, ListOpts.SlidersWhenNeeded(), ListOpts.HideHorizontalSlider(),
		ListOpts.Entries([]any{"first"}), ListOpts.EntryLabelFunc(func(e any) string { return e.(string) })) //nolint:forcetypeassert // strings
	l.SetLocation(img.Rect(0, 0, 200, 120))
	l.RequestRelayout() // as a parent laying it out does
	screen := ebiten.NewImage(200, 120)
	l.Render(screen) // laid out
	l.Render(screen)
	is.Equal(l.vSlider.GetWidget().GetVisibility(), Visibility_Hide_Blocking) // one entry fits

	for i := range 20 {
		l.AddEntry("entry " + strconv.Itoa(i))
	}
	l.Render(screen)
	is.Equal(l.vSlider.GetWidget().GetVisibility(), Visibility_Show) // twenty-one do not
}
