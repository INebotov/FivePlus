package web

import (
	"Backend/Errors"
	"Backend/db"
	"crypto/sha256"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type ResSelf struct {
	db.User
	Devices []db.Device `json:"devices"`
	Balance float64     `json:"balance"`
}

func (r *Router) GetSelf(c *gin.Context) {
	userID := c.GetUint("ID")
	userLevel := c.GetInt("LEVEL")

	user := db.User{ID: userID}
	err := r.DataBase.GetUser(&user)
	if err != nil {
		r.Error.DropError(c, err)
		return
	}

	dev := r.DataBase.DetAllUserDevices(user.ID)
	var balance float64

	if userLevel != db.ChildAccessLevel {
		balance, err = r.DataBase.GetBalance(user.BillingID)
		if err != nil {
			r.Error.DropError(c, err)
			return
		}
	}

	c.JSON(200, ResSelf{
		User:    user,
		Devices: dev,
		Balance: balance,
	})
}

type ChangeSelf struct {
	FiledName string `json:"filed_name"`
	NewValue  string `json:"new_value"`
	Operation string `json:"operation"`
}

func (r *Router) parceChanges(changes []ChangeSelf, userLevel, userID uint) (res db.User, isSecured bool, operations []string) {
	isSecured = false
	operations = make([]string, len(changes))
	for i, e := range changes {
		operations[i] = e.FiledName
		switch e.FiledName {
		case "name":
			if e.Operation == "change" {
				res.Name = e.NewValue
			}
		case "email":
			if e.Operation == "change" {
				res.Email = e.NewValue
				err := r.DataBase.ChangeUserBool(userID, "email_confermed", false)
				if err != nil {
					continue
				}
				isSecured = true
			}
		case "phone":
			if e.Operation == "change" {
				res.Phone = e.NewValue
				err := r.DataBase.ChangeUserBool(userID, "phone_confirmed", false)
				if err != nil {
					continue
				}
				isSecured = true
			}
		case "password":
			if e.Operation == "change" {
				hasher := sha256.New()
				hasher.Write([]byte(e.NewValue))

				res.Password = fmt.Sprintf("%x", hasher.Sum(nil))
				isSecured = true
			}
		case "photo":
			if e.Operation == "change" {
				res.Name = e.NewValue
			}
		case "username":
			if e.Operation == "change" {
				res.Name = e.NewValue
				isSecured = true
			}
		case "subject":
			sub, err := strconv.ParseUint(e.NewValue, 10, 32)
			if err != nil || userLevel != db.TeacherAccessLevel {
				continue
			}
			if e.Operation == "add" {
				if r.DataBase.AddSubject(r.DataBase.GetTeacherIDByUserId(userID), uint(sub)) != nil {
					continue
				}
			} else if e.Operation == "delete" {
				if r.DataBase.DeleteSubject(r.DataBase.GetTeacherIDByUserId(userID), uint(sub)) != nil {
					continue
				}
			}
		case "about":
			if e.Operation == "change" {
				res.About = e.NewValue
			}
		case "grade":
			gr, err := strconv.ParseUint(e.NewValue, 10, 32)
			if err != nil || userLevel != db.TeacherAccessLevel {
				continue
			}
			if e.Operation == "add" {
				if r.DataBase.AddGrade(r.DataBase.GetTeacherIDByUserId(userID), uint(gr)) != nil {
					continue
				}
			} else if e.Operation == "delete" {
				if r.DataBase.DeleteGrade(r.DataBase.GetTeacherIDByUserId(userID), uint(gr)) != nil {
					continue
				}
			}
		}
	}

	return
}

type ChangeResp struct {
	StandartResponse
	FieldsUpdated []string `json:"fields_updated"`
}

func (r *Router) ChangeSelf(c *gin.Context) {
	var body []ChangeSelf
	if err := c.BindJSON(&body); err != nil {
		r.Error.DropError(c, Errors.CantParceBodyError)
		return
	}

	userID := c.GetUint("ID")
	levelID := c.GetUint("LEVEL")

	user, secured, ops := r.parceChanges(body, userID, levelID)
	user.ID = userID

	if secured {
		if !c.GetBool("ACTION_PRESENT") || c.GetBool("ACTION_EXPIRED") || c.GetBool("ACTION_INVALID") {
			r.Error.DropError(c, Errors.Unauthorizated401Error)
			return
		}

		sli := c.GetStringSlice("ACTION_OPERATION")
		for _, el := range ops {
			for _, a := range sli {
				if a != el {
					r.Error.DropError(c, Errors.Unauthorizated401Error)
					return
				}
			}
		}
	}

	err := r.DataBase.UpdateUser(user)
	if err != nil {
		r.Error.DropError(c, err)
		return
	}

	c.JSON(200, ChangeResp{
		StandartResponse: StandartResponse{200, "Successfully updated user"},
		FieldsUpdated:    ops,
	})
}

func (r *Router) ConfirmSend(c *gin.Context) {
	what := c.Param("what")
	thru := c.Query("thru")

	userID := c.GetUint("ID")

	conf := db.Confirmation{
		UserID: userID,
	}
	var thruInt int
	var exp time.Duration

	if thru == "email" {
		thruInt = db.EmailConfirmationType
		exp = r.EmailConfExpired
	} else if thru == "phone" {
		exp = r.PhoneConfExpired
		thruInt = db.PhoneConfirmationType
	} else {
		r.Error.DropError(c, Errors.BadRequest)
		return
	}

	conf.ActionType = what
	switch what {
	case "email_confirm":
		conf.Type = db.EmailConfirmationType
	case "phone_confirm":
		conf.Type = db.PhoneConfirmationType
	case "profile_change":
		conf.Type = thruInt
	case "child_delete":
		conf.Type = thruInt
	case "universal":
		conf.Type = thruInt
	default:
		r.Error.DropError(c, Errors.BadRequest)
		return
	}

	err := r.DataBase.CreateConfirmation(&conf, exp)
	if err != nil {
		r.Error.DropError(c, err)
		return
	}

	c.JSON(200, StandartResponse{
		Code:    200,
		Message: fmt.Sprintf("Confirmation has been sent. Please check your %s.", thru),
	})

}
func (r *Router) GetActonToken(c *gin.Context) {
	var body struct {
		Code int `json:"code"`
	}
	if err := c.BindJSON(&body); err != nil {
		r.Error.DropError(c, Errors.CantParceBodyError)
		return
	}

	what := c.Param("what")
	userID := c.GetUint("ID")

	ok := false
	for _, e := range []string{} {
		if e == what {
			ok = true
		}
	}

	if !ok {
		r.Error.DropError(c, Errors.BadRequest)
		return
	}

	err := r.DataBase.ConfirmOperation(userID, body.Code, what)
	if err != nil {
		r.Error.DropError(c, err)
		return
	}

	token, err := r.Auth.GenerateChangeToken(userID, what)
	if err != nil {
		r.Error.DropError(c, err)
		return
	}

	c.SetCookie("action", token.Key, int(token.Expires.Unix()), "/", "*", false, true)
	c.JSON(200, StandartResponse{200, "Successfully created action token"})
}

func (r *Router) Lessons(c *gin.Context) {
	c.String(200, "Not implemented %s", c.FullPath())
}
func (r *Router) Rooms(c *gin.Context) {
	c.String(200, "Not implemented %s", c.FullPath())
}

func (r *Router) GenInvLink(c *gin.Context) {
	c.String(200, "Not implemented %s", c.FullPath())
}
func (r *Router) DisableInvLink(c *gin.Context) {
	c.String(200, "Not implemented %s", c.FullPath())
}
func (r *Router) GetMyInvLinks(c *gin.Context) {
	c.String(200, "Not implemented %s", c.FullPath())
}

func (r *Router) AcceptInvite(c *gin.Context) {
	c.String(200, "Not implemented %s", c.FullPath())
}
func (r *Router) GetMyInvites(c *gin.Context) {
	c.String(200, "Not implemented %s", c.FullPath())
}

func (r *Router) SendToBeTeacher(c *gin.Context) {
	c.String(200, "Not implemented %s", c.FullPath())
}
func (r *Router) DeleteToBeTeacher(c *gin.Context) {
	c.String(200, "Not implemented %s", c.FullPath())
}
func (r *Router) GetToBeTeacher(c *gin.Context) {
	c.String(200, "Not implemented %s", c.FullPath())
}

func (r *Router) Billing(c *gin.Context) {
	c.String(200, "Not implemented %s", c.FullPath())
}
