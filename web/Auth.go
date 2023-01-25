package web

import (
	"Backend/Errors"
	"Backend/auth"
	"Backend/db"
	"crypto/sha256"
	"fmt"
	"github.com/gin-gonic/gin"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

type Registration struct {
	UserName string `json:"user_name"`
	Name     string `json:"name"`
	Email    string `json:"email"`

	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func (r *Router) Registration(c *gin.Context) {
	var body Registration
	if err := c.BindJSON(&body); err != nil {
		r.Error.DropError(c, Errors.CantParceBodyError)
		return
	}

	accessPresent := c.GetBool("ACCESS_PRESENT")
	refreshPresent := c.GetBool("REFRESH_PRESENT")
	if accessPresent || refreshPresent {
		r.Error.DropError(c, Errors.AlreadyLoggedIn)
		return
	}

	if passwordvalidator.Validate(body.Password, 60) != nil {
		r.Error.DropError(c, Errors.TooWeakPassword)
		return
	}

	hasher := sha256.New()
	hasher.Write([]byte(body.Password))

	user := db.User{
		UserName:       body.UserName,
		Name:           body.Name,
		Email:          body.Email,
		EmailConfermed: false,
		Password:       fmt.Sprintf("%x", hasher.Sum(nil)),
		AccessLevel:    db.GeneralAccessLevel,
	}

	err := r.DataBase.CreateUser(&user)
	if err != nil {
		r.Error.DropError(c, err)
		return
	}

	token, err := r.Auth.GenerateToken(user.ID, user.AccessLevel)
	if err != nil {
		r.Error.DropError(c, err)
		return
	}

	err = r.DataBase.CreateSession(db.Session{
		Refresh:   token.Refresh.Key,
		ExpiresIn: token.Refresh.Expires.Unix(),
		UserLevel: user.AccessLevel,
		UserID:    user.ID,
		Device:    c.GetHeader("User-Agent"),
	})
	if err != nil {
		r.Error.DropError(c, err)
		return
	}

	c.SetCookie("access", token.Access.Key, int(token.Access.Expires.Unix()), "/", "*", false, true)
	c.SetCookie("refresh", token.Refresh.Key, int(token.Refresh.Expires.Unix()), "/", "*", false, true)

	c.JSON(200, StandartResponse{
		Code:    200,
		Message: "Successfully created account! Please confirm Email and Phone.",
	})
}

type Login struct {
	LoginBy  string `json:"login_by"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (r *Router) Login(c *gin.Context) {
	var body Login
	if err := c.BindJSON(&body); err != nil {
		r.Error.DropError(c, Errors.CantParceBodyError)
		return
	}

	accessPresent := c.GetBool("ACCESS_PRESENT")
	refreshPresent := c.GetBool("REFRESH_PRESENT")
	if accessPresent || refreshPresent {
		r.Error.DropError(c, Errors.AlreadyLoggedIn)
		return
	}

	user := db.User{}

	switch body.LoginBy {
	case "email":
		user.Email = body.Login
	case "phone":
		user.Phone = body.Login
	case "username":
		user.UserName = body.Login
	}

	hasher := sha256.New()
	hasher.Write([]byte(body.Password))

	user.Password = fmt.Sprintf("%x", hasher.Sum(nil))

	if r.DataBase.GetUser(&user) != nil || user.ID == 0 {
		r.Error.DropError(c, Errors.WrongCrendetails)
		return
	}

	token, err := r.Auth.GenerateToken(user.ID, user.AccessLevel)
	if err != nil {
		r.Error.DropError(c, err)
		return
	}

	err = r.DataBase.CreateSession(db.Session{
		Refresh:   token.Refresh.Key,
		ExpiresIn: token.Refresh.Expires.Unix(),
		UserLevel: user.AccessLevel,
		UserID:    user.ID,
		Device:    c.GetHeader("User-Agent"),
	})
	if err != nil {
		r.Error.DropError(c, err)
		return
	}

	c.SetCookie("access", token.Access.Key, int(token.Access.Expires.Unix()), "/", "*", false, true)
	c.SetCookie("refresh", token.Refresh.Key, int(token.Refresh.Expires.Unix()), "/", "*", false, true)

	c.JSON(200, StandartResponse{
		Code:    200,
		Message: "Successful login!",
	})
}
func (r *Router) Logout(c *gin.Context) {
	accessPresent := c.GetBool("ACCESS_PRESENT")
	refreshPresent := c.GetBool("REFRESH_PRESENT")

	if !accessPresent && !refreshPresent {
		r.Error.DropError(c, Errors.BadRequest)
		return
	}

	refresh := c.GetString("REFRESH")

	err := r.DataBase.DeleteSession(auth.Token{Key: refresh})
	if err != nil {
		r.Error.DropError(c, err)
		return
	}

	c.SetCookie("access", "", -1, "/", "*", false, true)
	c.SetCookie("refresh", "", -1, "/", "*", false, true)

	c.JSON(200, StandartResponse{
		Code:    200,
		Message: "Successful logout!",
	})
}
