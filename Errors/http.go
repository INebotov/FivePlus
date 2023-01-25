package Errors

import (
	"Backend/Metrics"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
)

type StandartErrResponse struct {
	Code    int    `json:"code"`
	Error   string `json:"error"`
	Message string `json:"message"`

	HandleFunc func(c *gin.Context) `json:"-"`
}

type HttpErrors struct {
	Metrics *Metrics.Metrics
	Log     *zap.Logger

	mapa map[error]StandartErrResponse

	Standart500Error StandartErrResponse
}

func Init(e HttpErrors) HttpErrors {
	res := make(map[error]StandartErrResponse)

	res[BadSigningMethodError] = StandartErrResponse{
		Code:    401,
		Error:   "A-BTSM",
		Message: "You token is invalid please try to re-login",
		HandleFunc: func(c *gin.Context) {
			c.SetCookie("access", "", 1, "/", "*", false, true)
			c.SetCookie("refresh", "", 1, "/", "*", false, true)
		},
	}
	res[TokenIsIvalidError] = StandartErrResponse{
		Code:    401,
		Error:   "A-TIIE",
		Message: "You token is invalid please try to re-login",
		HandleFunc: func(c *gin.Context) {
			c.SetCookie("access", "", 1, "/", "*", false, true)
			c.SetCookie("refresh", "", 1, "/", "*", false, true)
		},
	}
	res[NoRowsAffectedError] = StandartErrResponse{
		Code:       404,
		Error:      "D-NRAE",
		Message:    "Record you requested not found. Please check you data and try again.",
		HandleFunc: func(c *gin.Context) {},
	}
	res[WrongPasswordHashError] = StandartErrResponse{
		Code:       400,
		Error:      "V-WPHE",
		Message:    "For some reason your password hash is not correct! Please try again if that dont help unfortunately you may need to re-register.",
		HandleFunc: func(c *gin.Context) {},
	}
	res[PhoneIncorrectError] = StandartErrResponse{
		Code:       400,
		Error:      "V-PIIE",
		Message:    "You phone is incorrect! Please check your data.",
		HandleFunc: func(c *gin.Context) {},
	}
	res[WrongConfirmationTypeError] = StandartErrResponse{
		Code:       400,
		Error:      "V-WCTE",
		Message:    "Confirmation type you provided is not supportable. Please check your data again",
		HandleFunc: func(c *gin.Context) {},
	}
	res[ConfirmationExpiredError] = StandartErrResponse{
		Code:       400,
		Error:      "A-CHEE",
		Message:    "Confirmation has expired!",
		HandleFunc: func(c *gin.Context) {},
	}
	res[CantCreateBillingIDError] = StandartErrResponse{
		Code:       500,
		Error:      "D-CCBI",
		Message:    "For some reason we cant create billing id for your account. Please try again.",
		HandleFunc: func(c *gin.Context) {},
	}
	res[HaveNotFoundsToDoOperationError] = StandartErrResponse{
		Code:       403,
		Error:      "B-HNFE",
		Message:    "You have no enough founds to do this action.",
		HandleFunc: func(c *gin.Context) {},
	}
	res[AlreadyExistsError] = StandartErrResponse{
		Code:       400,
		Error:      "D-RAEE",
		Message:    "Record you trying to create already exists.",
		HandleFunc: func(c *gin.Context) {},
	}
	res[WrongUsernameError] = StandartErrResponse{
		Code:       400,
		Error:      "D-WUNE",
		Message:    "You provided wrong username.",
		HandleFunc: func(c *gin.Context) {},
	}
	res[RefreshIsObsentError] = StandartErrResponse{
		Code:       401,
		Error:      "A-RTAE",
		Message:    "Your refresh token is absent. Please re-login.",
		HandleFunc: func(c *gin.Context) {},
	}
	res[RefreshIsObsentError] = StandartErrResponse{
		Code:       401,
		Error:      "A-RTAE",
		Message:    "Your refresh token is absent. Please re-login.",
		HandleFunc: func(c *gin.Context) {},
	}

	e.mapa = res
	return e
}

func (e HttpErrors) Http_MetricsError(err error) {
	e.Log.Error("Error in metrics while doing request", zap.Field{
		Key:    "Error",
		Type:   zapcore.ErrorType,
		String: err.Error(),
	})
}
func actionFromError(err string) string {
	categoryError := strings.Split(err, "-")
	if len(categoryError) != 2 {
		return ""
	}
	switch categoryError[0] {
	case "A":
		return "auth_failure"
	case "B":
		return ""
	case "V":
		return "bad_request"
	case "D":
		return "error"
	}
	return ""
}
func (e HttpErrors) DropError(c *gin.Context, err error) {
	if err == nil {
		return
	}
	e.Log.Error("Error while handling request", zap.Field{
		Key:    "Error",
		Type:   zapcore.ErrorType,
		String: err.Error(),
	})
	errTwo := e.Metrics.Action("error")
	if errTwo != nil {
		e.Http_MetricsError(errTwo)
	}

	val, ok := e.mapa[err]
	if !ok {
		c.JSON(500, e.Standart500Error)
		return
	}

	act := actionFromError(val.Error)
	if act != "" {
		errTwo := e.Metrics.Action(act)
		if errTwo != nil {
			e.Http_MetricsError(errTwo)
		}
	}

	val.HandleFunc(c)
	c.JSON(val.Code, StandartErrResponse{
		Code:    val.Code,
		Error:   val.Error,
		Message: val.Message,
	})
}
