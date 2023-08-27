//go:build !js
// +build !js

package jsUtil

func IsMobileBrowser() bool {
	return false
}

func Prompt(mode MobileInputMode, title string, value string, cursorPos int, yPos int, callback InsertCallBack) {

}
func SetCursorPosition(posX int) {
}
func GetCursorPosition() int {
	return 0
}
