package router

import "github.com/gofiber/fiber/v2"

type ClassicError string

func (e ClassicError) Error() string {
	return string(e)
}

const (
	RefreshIsObsentError = ClassicError("cant find refresh token id db")
)

type HttpError struct {
	Code int
	Body Error
}
type Error struct {
	Code          string `json:"code"`
	Message       string `json:"message"`
	PossibleFixes string `json:"possible_fixes"`
}

func Uncategorized500Error() HttpError {
	return HttpError{
		Code: 500,
		Body: Error{
			Code:          "U500E",
			Message:       "In our server something went wrong :(",
			PossibleFixes: "Thus error has occupied on our server. We do all to make this never happened again",
		},
	}
}
func BadRequestUncatigorised() HttpError {
	return HttpError{
		Code: 400,
		Body: Error{
			Code:          "BR400U",
			Message:       "We think you made mistake in your request.",
			PossibleFixes: "Check your data again maybe it will help.",
		},
	}
}
func Unauthorizated401() HttpError {
	return HttpError{
		Code: 401,
		Body: Error{
			Code:          "UA400",
			Message:       "You have no access or yours credentials is incorrect",
			PossibleFixes: "Check again you tocen or password",
		},
	}
}

func Drop500Error(c *fiber.Ctx, err error) error {
	httpErr := Uncategorized500Error()
	Nerr := c.Status(httpErr.Code).JSON(httpErr.Body)
	if Nerr != nil {
		return Nerr
	}
	return err
}

func Drop400Error(c *fiber.Ctx) error {
	httpErr := BadRequestUncatigorised()
	err := c.Status(httpErr.Code).JSON(httpErr.Body)
	if err != nil {
		return err
	}
	return nil
}

func Drop401Error(c *fiber.Ctx) error {
	httpErr := Unauthorizated401()
	err := c.Status(httpErr.Code).JSON(httpErr.Body)
	if err != nil {
		return err
	}
	return nil
}
