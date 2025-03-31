//go:build !js
// +build !js

package jsUtil

func IsMobileBrowser() bool {
	return false
}

func Prompt(mode MobileInputMode, title string, value string, cursorPos int, yPos int, callback InsertCallBack, selectAll SelectAllCallback) {

}

func SetCursorPosition(cursorPos int, cursorPos2 int) {
}

func GetCursorPosition() int {
	return 0
}
