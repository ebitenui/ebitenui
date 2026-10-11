package html_test

import (
	"errors"
	"image"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	uihtml "github.com/ebitenui/ebitenui/html"
	"github.com/ebitenui/ebitenui/themes"
	"github.com/ebitenui/ebitenui/widget"
)

func TestLoad(t *testing.T) {
	theme := themes.GetBasicDarkTheme()
	d, err := uihtml.Parse(`<!DOCTYPE html><!-- ignored -->
 <html><body><div id="panel" data-layout="grid" data-columns="2" data-spacing="8" data-padding="4">
 <p id="title">Hello &amp; welcome<br>Next line</p>
 <button id="save" width="120" height="30" disabled>Save</button>
 <input id="name" value="Ada &lt;3" placeholder="Name">
 <input id="password" type="password">
 <input id="check" type="checkbox" checked="checked">
 <div id="row" data-layout="row"><span>Left</span><label>Right</label></div>
 </div></body></html>`, uihtml.Options{Theme: theme})
	if err != nil {
		t.Fatal(err)
	}
	if d.Root.GetWidget().GetTheme() != theme {
		t.Fatal("theme not installed")
	}
	panel := d.GetByID("panel").(*widget.Container)
	if len(panel.Children()) != 6 {
		t.Fatalf("children: %d", len(panel.Children()))
	}
	if d.GetByID("title").(*widget.Text).Label != "Hello & welcome\nNext line" {
		t.Fatal("text/entities/br")
	}
	button := d.GetByID("save").(*widget.Button)
	if !button.GetWidget().Disabled || button.GetWidget().MinWidth != 120 || button.GetWidget().MinHeight != 30 {
		t.Fatal("widget attributes")
	}
	if d.GetByID("name").(*widget.TextInput).GetText() != "Ada <3" {
		t.Fatal("input value")
	}
	if d.GetByID("check").(*widget.Checkbox).State() != widget.WidgetChecked {
		t.Fatal("checkbox state")
	}
	if d.GetByID("missing") != nil {
		t.Fatal("missing id")
	}
	if panel.Children()[1] != button {
		t.Fatal("lookup must return the actual tree object")
	}
	d.Root.Validate()
	d.Root.SetLocation(image.Rect(0, 0, 800, 600))
	d.Root.PreferredSize()
}

func TestRejectUnsupportedAndInvalidHTML(t *testing.T) {
	for _, source := range []string{
		`<script>alert(1)</script>`, `<img src="x">`, `<div style="display:flex"></div>`,
		`<a href="x">link</a>`, `<input type="radio">`, `<p><span>nested</span></p>`,
		`<div>bare text</div>`, `<button onclick="save()">Save</button>`,
		`<div id="x"></div><button id="x"></button>`, `<div id=""></div>`,
		`<div id="x" id="y"></div>`, `<div width="-1"></div>`, `<div height="10px"></div>`,
		`<div data-layout="flex"></div>`, `<div data-columns="2"></div>`,
		`<div data-layout="grid" data-columns="0"></div>`, `<div data-spacing="oops"></div>`,
		`<div data-padding="-2"></div>`, `<input checked>`, `<input type="checkbox" value="x">`,
		`<input type="checkbox" placeholder="x">`, `<input type="checkbox" checked="false">`,
		`<button disabled="false"></button>`, `<div><button></div>`, `<div>`, `</div>`,
		`<br>`, `<p>a<br id="x">b</p>`, `<!doctype svg>`,
	} {
		t.Run(source, func(t *testing.T) {
			d, err := uihtml.Parse(source, uihtml.Options{})
			if err == nil || d != nil {
				t.Fatalf("expected error and no document: %v", err)
			}
		})
	}
}

func TestFragmentsAndIndependentLoads(t *testing.T) {
	source := `<BUTTON ID="b">One &quot;two&quot;</BUTTON><input id="i"/><div/>`
	a, err := uihtml.Parse(source, uihtml.Options{})
	if err != nil {
		t.Fatal(err)
	}
	b, err := uihtml.Parse(source, uihtml.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if len(a.Root.Children()) != 3 || a.GetByID("b") == b.GetByID("b") {
		t.Fatal("independent fragments")
	}
	empty, err := uihtml.Parse(" \n<!-- empty -->", uihtml.Options{})
	if err != nil || len(empty.Root.Children()) != 0 {
		t.Fatal("empty document", err)
	}
}

type brokenReader struct{}

func (brokenReader) Read([]byte) (int, error) { return 0, io.ErrUnexpectedEOF }

func TestIO(t *testing.T) {
	if _, err := uihtml.Load(brokenReader{}, uihtml.Options{}); !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "layout.html")
	if err := os.WriteFile(path, []byte(`<div id="root"><p>File</p></div>`), 0600); err != nil {
		t.Fatal(err)
	}
	d, err := uihtml.LoadFile(path, uihtml.Options{})
	if err != nil || d.GetByID("root") == nil {
		t.Fatal(err)
	}
	if _, err := uihtml.LoadFile(path+".missing", uihtml.Options{}); err == nil {
		t.Fatal("missing file")
	}
	if _, err := uihtml.Load(strings.NewReader(`<p>Reader</p>`), uihtml.Options{}); err != nil {
		t.Fatal(err)
	}
}
