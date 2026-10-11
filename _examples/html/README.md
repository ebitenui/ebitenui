# HTML layout example

Run from the repository root:

```sh
go run ./_examples/html
```

Enter a name, choose whether the greeting should be enthusiastic, then click **Say hello**. **Reset** restores the initial state.

- [`layout.html`](layout.html) declares the containers, row/grid layouts, text, buttons, and inputs.
- [`main.go`](main.go) loads the embedded HTML with a theme, gets concrete widgets by their HTML IDs, and attaches ordinary ebitenui event handlers in Go.

The HTML is embedded so no runtime file path is required. To load an editable file from disk instead, use `uihtml.LoadFile("path/to/layout.html", options)`.

See the [HTML package documentation](../../html/README.md) for the supported subset. This is not a browser: CSS and JavaScript are not supported.
