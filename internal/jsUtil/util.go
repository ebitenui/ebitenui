package jsUtil

type InsertCallBack func(t string) string

type SelectAllCallback func()

type MobileInputMode string

const (
	TEXT      = MobileInputMode("text")
	DECIMAL   = MobileInputMode("decimal")
	NUMERIC   = MobileInputMode("numeric")
	TELEPHONE = MobileInputMode("tel")
	SEARCH    = MobileInputMode("search")
	EMAIL     = MobileInputMode("email")
	URL       = MobileInputMode("url")
)

const WASM = "wasm"
const JS = "js"
