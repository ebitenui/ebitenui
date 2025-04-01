package mobile

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
