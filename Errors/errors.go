package Errors

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	// Auth
	BadSigningMethodError   = Error("bad token signing method")
	TokenIsIvalidError      = Error("token is invalid")
	TokenExpiredError       = Error("401 Unauthorized Token Is Expired")
	SessionIsIncorrectError = Error("Session is incorrect")
	TokenIsObsentError      = Error("token is obsent")
	TooWeakPassword         = Error("password is too weak")
	WrongCrendetails        = Error("wrong user auth crendetails")

	// DB
	NoRowsAffectedError             = Error("no rows affected")
	WrongPasswordHashError          = Error("wrong password hash")
	PhoneIncorrectError             = Error("phone number is incorrect")
	WrongConfirmationTypeError      = Error("wrong confirmation type")
	ConfirmationExpiredError        = Error("confirmation has expired")
	CantCreateBillingIDError        = Error("too many attemps to create user billing id")
	HaveNotFoundsToDoOperationError = Error("user have no founds to do this operation")
	AlreadyExistsError              = Error("object already exists")
	WrongUsernameError              = Error("wrong username")
	WrongEmailError                 = Error("wrong email")

	// Web
	RefreshIsObsentError = Error("cant find refresh token id db")
	CantParceBodyError   = Error("request body is invalid")
	BadRequest           = Error("wrong request data (general)")
	AlreadyLoggedIn      = Error("user is already logged in")

	// HTTP
	Unauthorizated401Error = Error("401 unauthorized")
	PageNotExistError      = Error("404 Not Found")
)
