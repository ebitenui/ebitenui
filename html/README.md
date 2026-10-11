# Declarative HTML layouts

`github.com/ebitenui/ebitenui/html` loads a small HTML subset into actual ebitenui objects. It is not a browser and does not implement CSS, JavaScript, links, or browser layout/error recovery. Unsupported tags and attributes return errors.

```html
<div id="menu" data-layout="column" data-spacing="12" data-padding="16">
  <p id="title">Welcome &amp; enjoy!</p>
  <div data-layout="row" data-spacing="8">
    <input id="name" type="text" placeholder="Your name" width="200">
    <button id="start">Start</button>
  </div>
  <input id="sound" type="checkbox" checked>
</div>
```

```go
import (
    "github.com/ebitenui/ebitenui"
    uihtml "github.com/ebitenui/ebitenui/html"
    "github.com/ebitenui/ebitenui/themes"
    "github.com/ebitenui/ebitenui/widget"
)

// In your initialization function:
doc, err := uihtml.LoadFile("menu.html", uihtml.Options{
    Theme: themes.GetBasicDarkTheme(),
})
if err != nil {
    return err
}
ui := &ebitenui.UI{Container: doc.Root}
start := doc.GetByID("start").(*widget.Button)
start.ClickedEvent.AddHandler(func(args interface{}) {
    name := doc.GetByID("name").(*widget.TextInput).GetText()
    _ = name // Start your game here.
})
_ = ui // Use ui.Update() and ui.Draw(screen) in your game.
```

`Load(io.Reader, Options)` and `Parse(string, Options)` are also available. Each load produces independent objects. `GetByID` returns the concrete object in the tree, or nil when absent. IDs are case-sensitive and must be nonempty and unique. Keep the document to access its widgets later; handlers are attached in Go.

## Supported subset

| HTML | ebitenui object | Attributes beyond the common attributes |
| --- | --- | --- |
| `div`, `html`, `body` | `*widget.Container` | `data-layout`, `data-spacing`, `data-padding`, `data-columns` |
| `p`, `span`, `label` | `*widget.Text` | None |
| `button` | `*widget.Button` | None |
| `input` (default or `type="text"`) | `*widget.TextInput` | `type`, `value`, `placeholder` |
| `input type="password"` | Secure `*widget.TextInput` | `type`, `value`, `placeholder` |
| `input type="checkbox"` | `*widget.Checkbox` | `type`, `checked` |

Common attributes: `id`, `width`, `height`, `disabled`. Width and height set **minimum** size in integer pixels, not fixed browser dimensions. `disabled` and `checked` are HTML boolean attributes: use a bare attribute, an empty value, or the attribute's name; omit the attribute for false. The value `"false"` is rejected to avoid ambiguity.

Container layout defaults to `column`. `data-layout="row"` uses horizontal RowLayout; `data-layout="grid"` uses GridLayout with `data-columns` (default 1, must be positive). `data-spacing` separates children (both axes for grids), and `data-padding` applies uniform container padding. Both default to zero and accept non-negative integer pixels. `data-columns` is only valid for grids.

An implicit vertical container is always returned as `Document.Root`, so multiple top-level widgets are allowed. Optional `html` and `body` tags are ordinary containers. `<!doctype html>` and comments are ignored. Tags/attribute names are case-insensitive, as in HTML. HTML character entities are decoded.

Text must be inside `p`, `span`, `label`, or `button`. Surrounding whitespace is trimmed; internal whitespace is preserved. `<br>` inserts a newline within those elements. Other nested elements inside text/buttons are rejected. A `label` is just displayed text: browser `for` associations are not supported.

Use explicit closing tags or self-closing tags; `input` and `br` are void elements. Browser-style omitted closing tags are not supported. Use a complete theme (including the themes for widgets you use), just as for handwritten ebitenui code. With `Options{}`, attach a theme to `doc.Root` before validation/rendering. Loading builds objects but does not validate their visual resources.
