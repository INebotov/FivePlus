package router

import (
	"BackendSimple/db"
	"crypto/sha512"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type SelfChange struct {
	Changes []struct {
		FiledName string `json:"filed_name"`
		NewValue  string `json:"new_value"`
	} `json:"changes"`
}
type ResSelfChange struct {
	Code          int      `json:"code"`
	Message       string   `json:"message"`
	ValuesUpdated []string `json:"values_updated"`
}
type ResSelf struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`

	EmailConfemed  bool `json:"email_confemed"`
	PhoneConfermed bool `json:"phone_confermed"`

	Devices []string `json:"devices"`

	UserType string

	Additional []any `json:"additional"`
}

type UserConfidentional struct {
	UserName string `json:"user_name"`
	Name     string `json:"name"`
	Email    string `json:"email"`

	Role  string `json:"role"`
	Photo string `json:"photo"`
}
type TeacherConfidentionsl struct {
	ID uint `json:"id"`

	Comments []db.Comment `json:"comments"`

	Mark int `json:"mark"`

	Subjects []db.Subject `json:"subjects"`
}
type LessonsConfidentional struct {
	ID      uint       `json:"id"`
	Name    string     `json:"name"`
	Subject db.Subject `json:"subject"`

	Student UserConfidentional    `json:"student"`
	Teacher TeacherConfidentionsl `json:"teacher"`

	TimeStarted int64 `json:"time_started"`
	TimeEnded   int64 `json:"time_ended"`
}

func parceTeachersToTeachersReaponse(s []db.Teacher) []TeacherConfidentionsl {
	res := make([]TeacherConfidentionsl, len(s))

	for i, e := range s {
		res[i] = TeacherConfidentionsl{
			ID:   e.ID,
			Mark: e.ResultMark,
		}
	}
	return res
}
func parseUserToUserConfidentioal(u []db.User) []UserConfidentional {
	res := make([]UserConfidentional, len(u))
	for i, el := range u {
		res[i] = UserConfidentional{
			Name:  el.Name,
			Email: el.Email,

			Role:     el.LevelString,
			Photo:    el.Photo,
			UserName: el.UserName,
		}
	}
	return res
}
func (h Handlers) Self(c *fiber.Ctx) error {
	id := c.Locals("userid").(uint)
	if id == 0 {
		return Drop400Error(c)
	}

	level := c.Locals("level").(int)

	user := db.User{Model: gorm.Model{ID: id}}
	err := h.DB.GetUser(&user)
	if err != nil {
		return Drop500Error(c, err)
	}

	dev, err := h.DB.GetAllUsersDevices(user.ID)
	if err != nil {
		return Drop500Error(c, err)
	}

	var Additional []any

	if level == db.TeacherAccessLevel {
		t := db.Teacher{
			UserID: id,
		}

		err = h.DB.GetTeacher(&t)
		if err != nil {
			return Drop500Error(c, err)
		}

		Additional = append(Additional, TeacherConfidentionsl{
			ID:       t.ID,
			Comments: t.Comments,
			Mark:     t.ResultMark,
			Subjects: t.Subjects,
		})
	} else if level == db.GeneralAccessLevel {
		dep, err := h.DB.GetUsersChild(id)
		if err == nil {
			Additional = append(Additional, parseUserToUserConfidentioal(dep))
		}
	} else if level == db.ChildAccessLevel {
		dep, err := h.DB.GetUsersParents(id)
		if err == nil {
			Additional = append(Additional, parseUserToUserConfidentioal(dep))
		}
	}

	c.Status(200)
	return c.JSON(ResSelf{
		Name:          user.Name,
		Email:         user.Email,
		UserType:      user.LevelString,
		EmailConfemed: user.EmailConfermed,
		Devices:       dev,
		Additional:    Additional,
	})
}

func parceSelfChangeToUser(c SelfChange) (db.User, []string) {
	var Res db.User
	Fileds := make([]string, 0, len(c.Changes))
	for _, el := range c.Changes {
		switch el.FiledName {
		case "name":
			Res.Name = el.NewValue
			Fileds = append(Fileds, "name")
		case "photo":
			Res.Photo = el.NewValue
			Fileds = append(Fileds, "photo")
		}
	}
	return Res, Fileds
}
func (h Handlers) SelfChange(c *fiber.Ctx) error {
	id := c.Locals("userid").(uint)
	if id == 0 {
		return Drop400Error(c)
	}

	var u SelfChange
	if err := c.BodyParser(&u); err != nil {
		return Drop500Error(c, err)
	}

	user, fileds := parceSelfChangeToUser(u)
	user.ID = id

	err := h.DB.UpdateUser(user)
	if err != nil {
		return Drop500Error(c, err)
	}

	c.Status(200)
	return c.JSON(ResSelfChange{
		Code:          200,
		Message:       "Successfully updated",
		ValuesUpdated: fileds,
	})
}

func (h Handlers) SelfChangeSendCode(c *fiber.Ctx) error {
	id := c.Locals("userid").(uint)
	if id == 0 {
		return Drop400Error(c)
	}

	confirmThru, err := strconv.Atoi(c.Query("thru", "none"))
	if err != nil {
		return Drop400Error(c)
	}

	var userMy db.User
	userMy.ID = id

	err = h.DB.GetUser(&userMy)
	if err != nil {
		return Drop500Error(c, err)
	}

	var conf db.Confirmation
	var value string
	var ex time.Duration
	if confirmThru == db.EmailConfirmationType {
		conf = db.Confirmation{
			UserID: id,
			Type:   db.EmailConfirmationType,
			Action: db.ProfileSecretsChangeConfirmationAction,
		}
		value = userMy.Email
		ex = h.EmailConfExpired
	} else {
		return Drop400Error(c)
	}

	err = h.DB.CreateConfirmation(&conf, value, ex)
	if err != nil {
		return Drop500Error(c, err)
	}
	c.Status(200)
	return c.JSON(StandartResponse{
		Code:    200,
		Message: "Please confirm operation thurue email or phone",
	})
}
func (h Handlers) SelfChangeGetTokenFromCode(c *fiber.Ctx) error {
	id := c.Locals("userid").(uint)
	if id == 0 {
		return Drop400Error(c)
	}
	code, err := strconv.Atoi(c.Query("code", "none"))
	if err != nil {
		return Drop400Error(c)
	}
	err = h.DB.ConfirmOperation(id, code, db.ProfileSecretsChangeConfirmationAction)
	if err != nil {
		return Drop401Error(c)
	}

	key, err := h.Auth.GenerateChangeToken(id)
	if err != nil {
		return Drop500Error(c, err)
	}
	c.Status(200)
	c.Cookie(&fiber.Cookie{
		Name:     "changekey",
		Value:    key.Key,
		Expires:  key.Expires,
		Path:     "/self/change",
		HTTPOnly: true,
	})
	return c.JSON(StandartResponse{
		Code:    200,
		Message: "Successfully created tocken",
	})
}
func parceSecretsToUser(c SelfChange) (db.User, []string) {
	var user db.User
	Fileds := make([]string, 0, len(c.Changes))

	for _, el := range c.Changes {
		switch el.FiledName {
		case "email":
			user.Email = el.NewValue
			Fileds = append(Fileds, "email")
		case "password":
			passH := sha512.Sum512([]byte(el.NewValue))
			pass := passH[:]
			user.Password = pass
			Fileds = append(Fileds, "password")
		case "username":
			user.UserName = el.NewValue
			Fileds = append(Fileds, "username")
		}
	}
	return user, Fileds
}
func (h Handlers) SelfSecretsChange(c *fiber.Ctx) error {
	id := c.Locals("userid").(uint)
	if id == 0 {
		return Drop400Error(c)
	}

	//level := c.Locals("level").(uint)

	var u SelfChange
	if err := c.BodyParser(&u); err != nil {
		return Drop500Error(c, err)
	}
	key := c.Cookies("changekey")
	cl, err := h.Auth.GetTokenClaims(key)
	sub, ok := cl["sub"]
	if err != nil || !ok || uint(sub.(float64)) != id {
		return Drop400Error(c)
	}

	user, fileds := parceSecretsToUser(u)
	//if level == db.TeacherAccessLevel {
	//	// Teacher Change
	//}

	user.ID = id

	if contains(fileds, "email") {
		err = h.DB.ChangeUserBool(id, "email_confermed", false)
		if err != nil {
			return err
		}
	}

	err = h.DB.UpdateUser(user)
	if err != nil {
		return Drop500Error(c, err)
	}

	c.Status(200)
	return c.JSON(ResSelfChange{
		Code:          200,
		Message:       "Successfully updated",
		ValuesUpdated: fileds,
	})
}

func (h Handlers) SelfConfirmSend(c *fiber.Ctx) error {
	id := c.Locals("userid").(uint)
	if id == 0 {
		return Drop400Error(c)
	}
	what := c.Query("what", "none")

	user, err := h.DB.GetUserById(id)
	if err != nil {
		return Drop500Error(c, err)
	}

	if what == "email" {
		Emailconf := db.Confirmation{
			UserID: user.ID,
			Type:   db.EmailConfirmationType,
			Action: db.RegistrationConfirmationAction,
		}
		err = h.DB.CreateConfirmation(&Emailconf, user.Email, h.EmailConfExpired)
		if err != nil {
			return Drop500Error(c, err)
		}
	} else {
		return Drop400Error(c)
	}

	c.Status(200)
	return c.JSON(StandartResponse{
		Code:    200,
		Message: "Use code to confirm operation",
	})
}
func (h Handlers) SelfConfirm(c *fiber.Ctx) error {
	id := c.Locals("userid").(uint)
	if id == 0 {
		return Drop400Error(c)
	}

	code, err := strconv.Atoi(c.Query("code", "none"))
	if err != nil {
		return Drop400Error(c)
	}
	what := c.Query("what", "none")

	var userMy db.User
	userMy.ID = id
	if what == "email" {
		userMy.EmailConfermed = true
	} else {
		return Drop400Error(c)
	}

	err = h.DB.ConfirmOperation(id, code, db.RegistrationConfirmationAction)
	if err != nil {
		return Drop401Error(c)
	}

	err = h.DB.UpdateUser(userMy)
	if err != nil {
		return Drop500Error(c, err)
	}

	c.Status(200)
	return c.JSON(ResSelfChange{
		Code:          200,
		Message:       "Successfully confermed",
		ValuesUpdated: []string{what},
	})
}

func (h Handlers) fillLessonsConfidention(l []db.Lesson, level int) ([]LessonsConfidentional, error) {
	res := make([]LessonsConfidentional, len(l))
	var sub db.Subject
	var student db.User
	var teacher db.Teacher
	for i, e := range l {
		sub = db.Subject{
			ID: e.ID,
		}
		err := h.DB.GetSubject(&sub)
		if err != nil {
			return nil, err
		}

		if level == db.ChildAccessLevel {
			teacher.ID = e.TeacherID
			err = h.DB.GetTeacher(&teacher)
			if err != nil {
				return nil, err
			}
			res[i].Teacher = TeacherConfidentionsl{
				ID:       teacher.ID,
				Comments: teacher.Comments,
				Mark:     teacher.ResultMark,
				Subjects: teacher.Subjects,
			}
		} else {
			student.ID = e.StudentID
			err = h.DB.GetUser(&student)
			if err != nil {
				return nil, err
			}
			res[i].Student = UserConfidentional{
				UserName: student.UserName,
				Name:     student.Name,
				Email:    student.Email,
				Role:     student.LevelString,
				Photo:    student.Photo,
			}
		}

		res[i] = LessonsConfidentional{
			Name:        e.Name,
			Subject:     sub,
			TimeStarted: e.TimeStarted,
			TimeEnded:   e.TimeEnded,
		}
	}
	return res, nil
}
func (h Handlers) GetLessons(c *fiber.Ctx) error {
	id := c.Locals("userid").(uint)
	if id == 0 {
		return Drop400Error(c)
	}
	level := c.Locals("level").(int)

	limit, err := strconv.Atoi(c.Query("limit", "none"))
	if err != nil || limit > 100 {
		limit = 25
	}

	startFrom, err := strconv.Atoi(c.Query("startfrom", "none"))
	if err != nil {
		startFrom = 0
	}

	if level == db.TeacherAccessLevel {
		id, err = h.DB.GetTeacherIDByUserId(id)
		if err != nil {
			return Drop500Error(c, err)
		}
	}

	lessons, err := h.DB.GetLessons(id, limit, startFrom)
	if err != nil {
		return Drop500Error(c, err)
	}

	res, err := h.fillLessonsConfidention(lessons, level)
	if err != nil {
		return Drop500Error(c, err)
	}

	c.Status(200)
	return c.JSON(res)
}
