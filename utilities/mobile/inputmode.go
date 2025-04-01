package mobile

type InputMode string

const (
	TEXT      = InputMode("text")
	DECIMAL   = InputMode("decimal")
	NUMERIC   = InputMode("numeric")
	TELEPHONE = InputMode("tel")
	SEARCH    = InputMode("search")
	EMAIL     = InputMode("email")
	URL       = InputMode("url")
)
