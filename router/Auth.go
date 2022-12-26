package router

import (
	"BackendSimple/auth"
	"BackendSimple/db"
	"crypto/sha512"
	"github.com/gofiber/fiber/v2"
	"gopkg.in/dealancer/validate.v2"
)

type Reg struct {
	UserName string `json:"user_name"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type Login struct {
	LoginBy string `json:"login_by"`
	Login   string `json:"login"`

	Password string `json:"password"`
}
type ItIs struct {
	Value string `json:"value"`
	Is    bool
}

func sendKeys(c *fiber.Ctx, tokens auth.TokenPair) {
	c.Cookie(&fiber.Cookie{
		Name:     "access",
		Value:    tokens.Access,
		Expires:  tokens.AccessExpires,
		HTTPOnly: true,
	})
	c.Cookie(&fiber.Cookie{
		Name:     "refresh",
		Value:    tokens.Refresh,
		Expires:  tokens.RefreshExpires,
		HTTPOnly: true,
	})
	c.Status(200)
}
func (h Handlers) Reg(c *fiber.Ctx) error {
	var u Reg
	if err := c.BodyParser(&u); err != nil {
		return Drop500Error(c, err)
	}
	passH := sha512.Sum512([]byte(u.Password))
	pass := passH[:]

	user := db.User{
		UserName:       u.UserName,
		Name:           u.Name,
		Email:          u.Email,
		Password:       pass,
		AccessLevel:    db.GeneralAccessLevel,
		EmailConfermed: false,
	}
	err := validate.Validate(&user)
	if err != nil || h.DB.CheckUserPersistance(user) {
		return Drop400Error(c)
	}

	err = h.DB.CreateUser(&user)
	if err != nil {
		return Drop500Error(c, err)
	}

	err = h.DB.CreateBillingAccount(user.BillingID)
	if err != nil {
		return Drop500Error(c, err)
	}

	tokens, err := h.Auth.GenerateToken(user.ID, db.GeneralAccessLevel, c.GetReqHeaders()["User-Agent"])
	if err != nil {
		return Drop500Error(c, err)
	}

	sendKeys(c, tokens)

	return c.JSON(StandartResponse{
		Code:    200,
		Message: "Successfully registered!",
	})
}
func (h Handlers) Login(c *fiber.Ctx) error {
	var u Login
	if err := c.BodyParser(&u); err != nil {
		return Drop400Error(c)
	}
	passH := sha512.Sum512([]byte(u.Password))
	pass := passH[:]

	user := db.User{Password: pass}

	if u.LoginBy == "username" {
		user.UserName = u.Login
	} else if u.LoginBy == "email" {
		user.Email = u.Login
	} else {
		return Drop400Error(c)
	}

	err := h.DB.GetUser(&user)
	if err != nil {
		return Drop401Error(c)
	}

	tokens, err := h.Auth.GenerateToken(user.ID, user.AccessLevel, c.GetReqHeaders()["User-Agent"])
	if err != nil {
		return Drop500Error(c, err)
	}
	sendKeys(c, tokens)
	return c.JSON(StandartResponse{
		Code:    200,
		Message: "Successfully authorized!",
	})
}
func (h Handlers) Logout(c *fiber.Ctx) error {
	refresh := c.Cookies("refresh")
	access := c.Cookies("access")

	if refresh == "" || access == "" {
		return Drop400Error(c)
	}

	keyToDel := db.RefreshToken{
		Refresh: refresh,
	}

	err := h.DB.DelToken(&keyToDel)
	if err != nil {
		return Drop500Error(c, err)
	}
	c.Status(200)
	c.ClearCookie("refresh", "access")
	return c.JSON(StandartResponse{Code: 200, Message: "Successfully logged out!"})
}

func (h Handlers) isUserName(c *fiber.Ctx, username string) error {
	u := db.User{UserName: username}
	err := h.DB.GetUser(&u)
	c.Status(200)
	if err == nil {
		return c.JSON(ItIs{
			Value: username,
			Is:    true,
		})
	}
	return c.JSON(ItIs{
		Value: username,
		Is:    false,
	})
}
func (h Handlers) isEmail(c *fiber.Ctx, email string) error {
	u := db.User{Email: email}
	err := h.DB.GetUser(&u)
	if err == nil {
		c.Status(200)
		return c.JSON(ItIs{
			Value: email,
			Is:    true,
		})
	}
	c.Status(200)
	return c.JSON(ItIs{
		Value: email,
		Is:    false,
	})
}
func (h Handlers) Is(c *fiber.Ctx) error {
	email := c.Query("email")
	username := c.Query("username")
	if email == "" && username == "" {
		return Drop400Error(c)
	}
	if email != "" {
		return h.isEmail(c, email)
	} else {
		return h.isUserName(c, username)
	}
}
