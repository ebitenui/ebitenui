//go:build !js
// +build !js

package jsUtil

import "github.com/ebitenui/ebitenui/utilities/mobile"

func IsMobileBrowser() bool {
	return false
}

func Prompt(mode mobile.InputMode, title string, value string, cursorPos int, yPos int, callback InsertCallBack, selectAll SelectAllCallback) {

}

func SetCursorPosition(cursorPos int, cursorPos2 int) {
}

func GetCursorPosition() int {
	return 0
}
