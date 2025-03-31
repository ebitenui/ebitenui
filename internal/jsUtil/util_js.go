//go:build js
// +build js

package jsUtil

import (
	"regexp"
	"strconv"
	"syscall/js"
)

var MOBILE_BROWSER_REGEX = regexp.MustCompile("(?i)Android|webOS|iPhone|iPad|iPod|BlackBerry|Windows Phone")

var document js.Value

var insertCB InsertCallBack

var selectAllCB SelectAllCallback

var started bool

var offsetTop int

func init() {
	document = js.Global().Get("document")

	//Create a hidden html input element that will capture keystrokes
	p := document.Call("createElement", "input")
	p.Set("id", "tempInput")
	p.Set("style", "height:0px; width:1px; margin:0px; position: fixed; overflow:hidden; top:-10px; border:0px; padding:0px")
	document.Get("body").Call("appendChild", p)

	//Add a listener on the hidden html input for keystrokes
	p.Call("addEventListener", "input", js.FuncOf(handleInput), false)
	p.Call("addEventListener", "select", js.FuncOf(handleSelection), false)

	//Get the canvas and attach an event listener for screen touches
	requestAnimationFrame := js.Global().Get("requestAnimationFrame")
	requestAnimationFrame.Invoke(js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		canvas := document.Get("body").Call("getElementsByTagName", "canvas").Index(0)
		offsetTop = canvas.Get("offsetTop").Int()
		canvas.Call("addEventListener", "touchstart", js.FuncOf(handleClick), false)
		canvas.Call("addEventListener", "touchend", js.FuncOf(handleClick), false)
		return nil
	}))
}

func IsMobileBrowser() bool {
	navigator := js.Global().Get("navigator")
	userAgent := navigator.Get("userAgent")
	return MOBILE_BROWSER_REGEX.Match([]byte(userAgent.String()))
}

func Prompt(mode MobileInputMode, title string, value string, cursorPos int, yPos int, cb InsertCallBack, sa SelectAllCallback) {
	insertCB = cb
	selectAllCB = sa
	p := document.Call("getElementById", "tempInput")

	//Configure the hidden html input element based on what our library has for the input
	p.Call("setAttribute", "inputmode", string(mode))
	p.Set("value", value)
	p.Call("setSelectionRange", cursorPos, cursorPos)
	p.Get("style").Call("setProperty", "top", strconv.Itoa(offsetTop+yPos)+"px")

	//Indicate we've started capturing input
	started = true
}

// TODO: fix cursor position from being called every loop
func SetCursorPosition(cursorPos int, cursorPos2 int) {
	p := document.Call("getElementById", "tempInput")
	p.Call("setSelectionRange", cursorPos, cursorPos2)
}

func GetCursorPosition() int {
	p := document.Call("getElementById", "tempInput")
	return p.Get("selectionStart").Int()
}

func handleClick(this js.Value, args []js.Value) any {
	//If we have clicked on one of the inputs, shift focus to the input to open the keyboard
	if started {
		p := document.Call("getElementById", "tempInput")
		p.Call("focus")
		started = false
	}
	return nil
}

var previousValue = ""
var previousPosition = 0

// Process changes on the hidden html text input
func handleInput(this js.Value, args []js.Value) any {
	newTextString := args[0].Get("target").Get("value").String()
	if insertCB != nil {
		result := insertCB(newTextString)
		if result != newTextString {
			p := document.Call("getElementById", "tempInput")
			p.Set("value", result)
		}
		if result == previousValue {
			SetCursorPosition(previousPosition, previousPosition)
		} else {
			previousPosition = GetCursorPosition()
			previousValue = result
		}
	}
	return nil
}

func handleSelection(this js.Value, args []js.Value) any {
	if selectAllCB != nil {
		start := args[0].Get("target").Get("selectionStart").Int()
		end := args[0].Get("target").Get("selectionEnd").Int()
		if start == 0 && end == len([]rune(args[0].Get("target").Get("value").String())) {
			selectAllCB()
		}
	}
	return nil
}
