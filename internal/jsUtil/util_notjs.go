//go:build !js
// +build !js

package jsUtil

func IsMobileBrowser() bool {
	return false
}

func Prompt(title string, value string, cursorPos int, callback InsertCallBack) (string, bool) {
	return "", false
}
func SetCursorPosition(posX int) {
}
func GetCursorPosition() int {
	return 0
}
