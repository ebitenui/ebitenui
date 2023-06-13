//go:build js
// +build js

package jsUtil

import (
	"regexp"
	"syscall/js"
)

var MOBILE_BROWSER_REGEX = regexp.MustCompile("(?i)Android|webOS|iPhone|iPad|iPod|BlackBerry|Windows Phone")

func IsMobileBrowser() bool {
	navigator := js.Global().Get("navigator")
	userAgent := navigator.Get("userAgent")
	return MOBILE_BROWSER_REGEX.Match([]byte(userAgent.String()))
}

func Prompt(title string, value string) (string, bool) {
	prompt := js.Global().Get("prompt")
	result := prompt.Invoke(title, value)
	if !result.IsNull() && !result.IsUndefined() {
		return result.String(), true
	} else {
		return "", false
	}

}
