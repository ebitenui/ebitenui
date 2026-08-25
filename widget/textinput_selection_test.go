package widget

import "testing"

func TestSelectedText_NonASCII_RuneIndexing(t *testing.T) {
	const want = "Ünïcödé Ship Name"

	tin := NewTextInput()
	tin.SetText(want)
	tin.SelectAll()

	if got := tin.SelectedText(); got != want {
		t.Fatalf("SelectedText() after SelectAll() on %q = %q, want %q", want, got, want)
	}
}

func TestDeleteSelectedText_NonASCII_RuneIndexing(t *testing.T) {
	const text = "Ünïcödé Ship Name"

	tin := NewTextInput()
	tin.SetText(text)
	tin.SelectAll()
	tin.DeleteSelectedText()

	if got := tin.GetText(); got != "" {
		t.Fatalf("GetText() after SelectAll()+DeleteSelectedText() = %q, want empty", got)
	}
}

func TestDeleteSelectedText_PositionalNotSubstringSearch(t *testing.T) {
	tin := NewTextInput()
	tin.SetText("aXa")

	tin.dragStartIndex = 2
	tin.cursorPosition = 3

	if got := tin.SelectedText(); got != "a" {
		t.Fatalf("SelectedText() for [2:3] of \"aXa\" = %q, want %q", got, "a")
	}

	tin.DeleteSelectedText()

	const want = "aX"
	if got := tin.GetText(); got != want {
		t.Fatalf("GetText() after deleting [2:3] of \"aXa\" = %q, want %q", got, want)
	}
}

func TestDeleteSelectedText_PositionalDeleteAlsoCorrectsCursor(t *testing.T) {
	tin := NewTextInput()
	tin.SetText("aXa")
	tin.dragStartIndex = 2
	tin.cursorPosition = 3

	tin.DeleteSelectedText()

	if tin.cursorPosition != 2 {
		t.Fatalf("cursorPosition after deleting [2:3] of \"aXa\" = %d, want 2", tin.cursorPosition)
	}
	tin.Insert("Y")
	if got := tin.GetText(); got != "aXY" {
		t.Fatalf("GetText() after positional delete + Insert(\"Y\") = %q, want %q", got, "aXY")
	}
}
