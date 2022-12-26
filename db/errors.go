package db

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	NoRowsAffectedError             = Error("no rows affected")
	WrongPasswordHashError          = Error("wrong password hash")
	PhoneIncorrectError             = Error("phone number is incorrect")
	WrongConfirmationTypeError      = Error("wrong confirmation type")
	ConfirmationExpiredError        = Error("confirmation has expired")
	CantCreateBillingIDError        = Error("too many attemps to create user billing id")
	HaveNotFoundsToDoOperationError = Error("user have no founds to do this operation")
	AlreadyExistsError              = Error("object already exists")
	WrongUsernameError              = Error("wrong username")
)
