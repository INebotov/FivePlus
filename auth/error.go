package auth

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	BadSigningMethodError = Error("bad token signing method")
	TokenIsIvalidError    = Error("token is invalid")
)
