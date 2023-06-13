//go:build !js
// +build !js

package jsUtil

func IsMobileBrowser() bool {
	return false
}

func Prompt(title string, value string) (string, bool) {
	return "", false
}
