// Package html compiles a deliberately restricted HTML subset into ebitenui widgets.
// It is not a browser: unsupported elements and attributes are errors.
package html

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/ebitenui/ebitenui/widget"
	nethtml "golang.org/x/net/html"
)

// Options supplies the theme inherited by all generated widgets. If nil, attach
// a theme to Root before validating or rendering, as with hand-built widgets.
type Options struct{ Theme *widget.Theme }

// Document owns the generated widget tree. Root is an implicit vertical container.
// IDs refer to the concrete objects (for example *widget.Button), not their base Widget.
type Document struct {
	Root *widget.Container
	ids  map[string]widget.PreferredSizeLocateableWidget
}

// GetByID returns the concrete widget with id, or nil if none exists.
func (d *Document) GetByID(id string) widget.PreferredSizeLocateableWidget { return d.ids[id] }

// Parse loads an HTML string.
func Parse(source string, options Options) (*Document, error) {
	return Load(strings.NewReader(source), options)
}

// LoadFile loads an HTML file, closing it before returning.
func LoadFile(path string, options Options) (*Document, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Load(f, options)
}

type element struct {
	tag      string
	attrs    map[string]string
	children []*element
	text     strings.Builder
}

// Load reads HTML and generates a new tree on each call. Tags must be explicitly
// closed, except the void input and br elements. No scripts, CSS, or browser
// error recovery are supported. Errors never return a partial Document.
func Load(r io.Reader, options Options) (*Document, error) {
	root := &element{tag: "div", attrs: map[string]string{}}
	stack := []*element{root}
	z := nethtml.NewTokenizer(r)
	for {
		kind := z.Next()
		switch kind {
		case nethtml.ErrorToken:
			if err := z.Err(); err != io.EOF {
				return nil, err
			}
			if len(stack) != 1 {
				return nil, fmt.Errorf("html: unclosed <%s>", stack[len(stack)-1].tag)
			}
			d := &Document{ids: make(map[string]widget.PreferredSizeLocateableWidget)}
			w, err := d.compile(root)
			if err != nil {
				return nil, err
			}
			d.Root = w.(*widget.Container)
			d.Root.GetWidget().SetTheme(options.Theme)
			return d, nil
		case nethtml.CommentToken:
		case nethtml.DoctypeToken:
			if !strings.EqualFold(strings.TrimSpace(z.Token().Data), "html") {
				return nil, fmt.Errorf("html: only <!doctype html> is supported")
			}
		case nethtml.TextToken:
			parent := stack[len(stack)-1]
			parent.text.WriteString(string(z.Text()))
		case nethtml.StartTagToken, nethtml.SelfClosingTagToken:
			token := z.Token()
			e := &element{tag: token.Data, attrs: make(map[string]string)}
			for _, a := range token.Attr {
				if _, exists := e.attrs[a.Key]; exists {
					return nil, fmt.Errorf("html: duplicate attribute %q on <%s>", a.Key, e.tag)
				}
				e.attrs[a.Key] = a.Val
			}
			parent := stack[len(stack)-1]
			if e.tag == "br" {
				if !textElement(parent.tag) || len(e.attrs) != 0 {
					return nil, fmt.Errorf("html: <br> is only supported inside text or buttons, without attributes")
				}
				parent.text.WriteByte('\n')
			} else {
				parent.children = append(parent.children, e)
			}
			if kind == nethtml.StartTagToken && e.tag != "input" && e.tag != "br" {
				stack = append(stack, e)
			}
		case nethtml.EndTagToken:
			tag := z.Token().Data
			if len(stack) == 1 || stack[len(stack)-1].tag != tag {
				return nil, fmt.Errorf("html: unexpected closing </%s>", tag)
			}
			stack = stack[:len(stack)-1]
		}
	}
}

func textElement(tag string) bool {
	return tag == "p" || tag == "span" || tag == "label" || tag == "button"
}
func containerElement(tag string) bool { return tag == "div" || tag == "html" || tag == "body" }

func (d *Document) compile(e *element) (widget.PreferredSizeLocateableWidget, error) {
	fail := func(format string, args ...any) (widget.PreferredSizeLocateableWidget, error) {
		return nil, fmt.Errorf("html: <%s>: %s", e.tag, fmt.Sprintf(format, args...))
	}
	allowed := map[string]bool{"id": true, "disabled": true, "width": true, "height": true}
	switch {
	case containerElement(e.tag):
		for _, key := range []string{"data-layout", "data-spacing", "data-padding", "data-columns"} {
			allowed[key] = true
		}
	case textElement(e.tag):
	case e.tag == "input":
		for _, key := range []string{"type", "value", "placeholder", "checked"} {
			allowed[key] = true
		}
	default:
		return fail("unsupported element")
	}
	for key := range e.attrs {
		if !allowed[key] {
			return fail("unsupported attribute %q", key)
		}
	}
	for _, key := range []string{"disabled", "checked"} {
		if value, ok := e.attrs[key]; ok && value != "" && !strings.EqualFold(value, key) {
			return fail("%s is a boolean HTML attribute; omit it for false", key)
		}
	}
	number := func(key string, fallback int) (int, error) {
		value, ok := e.attrs[key]
		if !ok {
			return fallback, nil
		}
		n, err := strconv.Atoi(value)
		if err != nil || n < 0 {
			return 0, fmt.Errorf("html: <%s>: %s must be a non-negative integer", e.tag, key)
		}
		return n, nil
	}
	width, err := number("width", 0)
	if err != nil {
		return nil, err
	}
	height, err := number("height", 0)
	if err != nil {
		return nil, err
	}
	var w widget.PreferredSizeLocateableWidget
	if containerElement(e.tag) {
		if strings.TrimSpace(e.text.String()) != "" {
			return fail("put text inside <p>, <span>, or <label>")
		}
		spacing, err := number("data-spacing", 0)
		if err != nil {
			return nil, err
		}
		padding, err := number("data-padding", 0)
		if err != nil {
			return nil, err
		}
		insets := &widget.Insets{Top: padding, Bottom: padding, Left: padding, Right: padding}
		var layout widget.Layouter
		mode := e.attrs["data-layout"]
		switch mode {
		case "", "column", "row":
			if _, ok := e.attrs["data-columns"]; ok {
				return fail("data-columns requires data-layout=grid")
			}
			direction := widget.DirectionVertical
			if mode == "row" {
				direction = widget.DirectionHorizontal
			}
			layout = widget.NewRowLayout(widget.RowLayoutOpts.Direction(direction), widget.RowLayoutOpts.Spacing(spacing), widget.RowLayoutOpts.Padding(insets))
		case "grid":
			columns, err := number("data-columns", 1)
			if err != nil {
				return nil, err
			}
			if columns == 0 {
				return fail("data-columns must be positive")
			}
			layout = widget.NewGridLayout(widget.GridLayoutOpts.Columns(columns), widget.GridLayoutOpts.Spacing(spacing, spacing), widget.GridLayoutOpts.Padding(insets))
		default:
			return fail("unsupported data-layout %q", mode)
		}
		c := widget.NewContainer(widget.ContainerOpts.Layout(layout))
		for _, child := range e.children {
			cw, err := d.compile(child)
			if err != nil {
				return nil, err
			}
			c.AddChild(cw)
		}
		w = c
	} else {
		if len(e.children) != 0 {
			return fail("nested elements are not supported (except <br> in text)")
		}
		label := strings.TrimSpace(e.text.String())
		switch e.tag {
		case "p", "span", "label":
			w = widget.NewText(widget.TextOpts.TextLabel(label))
		case "button":
			w = widget.NewButton(widget.ButtonOpts.TextLabel(label))
		case "input":
			if label != "" {
				return fail("input cannot contain text")
			}
			typ := e.attrs["type"]
			_, checked := e.attrs["checked"]
			switch typ {
			case "", "text", "password":
				if checked {
					return fail("checked requires type=checkbox")
				}
				t := widget.NewTextInput(widget.TextInputOpts.Placeholder(e.attrs["placeholder"]), widget.TextInputOpts.Secure(typ == "password"))
				t.SetText(e.attrs["value"])
				w = t
			case "checkbox":
				if _, ok := e.attrs["placeholder"]; ok {
					return fail("placeholder is not supported on checkbox")
				}
				if _, ok := e.attrs["value"]; ok {
					return fail("value is not supported on checkbox")
				}
				state := widget.WidgetUnchecked
				if checked {
					state = widget.WidgetChecked
				}
				w = widget.NewCheckbox(widget.CheckboxOpts.InitialState(state))
			default:
				return fail("unsupported input type %q", typ)
			}
		}
	}
	base := w.GetWidget()
	base.MinWidth, base.MinHeight = width, height
	_, base.Disabled = e.attrs["disabled"]
	if id, ok := e.attrs["id"]; ok {
		if id == "" {
			return fail("id must not be empty")
		}
		if _, exists := d.ids[id]; exists {
			return fail("duplicate id %q", id)
		}
		d.ids[id] = w
	}
	return w, nil
}
