package router

//
//import (
//	"BackendSimple/db"
//	"crypto/sha512"
//	"github.com/gofiber/fiber/v2"
//	"gorm.io/gorm"
//	"strconv"
//	"time"
//)
//
//type CreateChild struct {
//	UserName string `json:"user_name"`
//	Name     string `json:"name"`
//
//	SameEmail bool   `json:"same_email"`
//	Email     string `json:"email"`
//
//	Password string `json:"password"`
//}
//type BillingAccountResponse struct {
//	Code         int              `json:"code"`
//	Balance      int64            `json:"balance"`
//	Transactions []db.Transaction `json:"transactions"`
//}
//type DeleteChild struct {
//	Code     int    `json:"code"`
//	UserName string `json:"user_name"`
//	ID       uint   `json:"id"`
//}
//
//func (h Handlers) CreateChild(c *fiber.Ctx) error {
//	var u CreateChild
//	if err := c.BodyParser(&u); err != nil {
//		return Drop500Error(c, err)
//	}
//	id := c.Locals("userid").(uint)
//	if id == 0 {
//		return Drop400Error(c)
//	}
//
//	parent, err := h.DB.GetUserById(id)
//	if err != nil {
//		return Drop500Error(c, err)
//	}
//
//	var child db.User
//	child.UserName = u.UserName
//	child.Name = u.Name
//	passH := sha512.Sum512([]byte(u.Password))
//	child.Password = passH[:]
//
//	if u.SameEmail {
//		child.Email = parent.Email
//		child.EmailConfermed = parent.EmailConfermed
//	} else {
//		child.Email = u.Email
//		child.EmailConfermed = false
//	}
//	child.AccessLevel = db.ChildAccessLevel
//	child.BillingID = parent.BillingID
//
//	err = h.DB.CreateChild(id, &child)
//	if err != nil {
//		return Drop500Error(c, err)
//	}
//	c.Status(200)
//	return c.JSON(StandartResponse{
//		Code:    200,
//		Message: "Successfully created child!",
//	})
//}
//func (h Handlers) DeleteChildSend(c *fiber.Ctx) error {
//	id := c.Locals("userid").(uint)
//	if id == 0 {
//		return Drop400Error(c)
//	}
//
//	confirmThru, err := strconv.Atoi(c.Query("thru", "none"))
//	if err != nil {
//		return Drop400Error(c)
//	}
//
//	var userMy db.User
//	userMy.ID = id
//
//	err = h.DB.GetUser(&userMy)
//	if err != nil {
//		return Drop500Error(c, err)
//	}
//
//	var conf db.Confirmation
//	var ex time.Duration
//	if confirmThru == db.EmailConfirmationType {
//		conf = db.Confirmation{
//			UserID: id,
//			Type:   db.EmailConfirmationType,
//			Action: db.DeteteChildAction,
//		}
//		ex = h.EmailConfExpired
//	} else {
//		return Drop400Error(c)
//	}
//
//	err = h.DB.CreateConfirmation(&conf, userMy, ex)
//	if err != nil {
//		return Drop500Error(c, err)
//	}
//	c.Status(200)
//	return c.JSON(StandartResponse{
//		Code:    200,
//		Message: "Please confirm operation thurue email or phone",
//	})
//}
//func (h Handlers) DeleteChild(c *fiber.Ctx) error {
//	id := c.Locals("userid").(uint)
//	if id == 0 {
//		return Drop400Error(c)
//	}
//	var req DeleteChild
//	if err := c.BodyParser(&req); err != nil {
//		return Drop500Error(c, err)
//	}
//	err := h.DB.ConfirmOperation(id, req.Code, db.DeteteChildAction)
//	if err != nil {
//		return Drop401Error(c)
//	}
//
//	err = h.DB.DeleteChild(id, &db.User{Model: gorm.Model{ID: req.ID}, UserName: req.UserName})
//	if err != nil {
//		return Drop500Error(c, err)
//	}
//	c.Status(200)
//	return c.JSON(StandartResponse{
//		Code:    200,
//		Message: "Successfully deleted child!",
//	})
//}
//
//func (h Handlers) GetBillingAccount(c *fiber.Ctx) error {
//	id := c.Locals("userid").(uint)
//	if id == 0 {
//		return Drop400Error(c)
//	}
//
//	user, err := h.DB.GetUserById(id)
//	if err != nil {
//		return Drop500Error(c, err)
//	}
//	var account db.BillingAccount
//	account.UserBillingID = user.BillingID
//
//	err = h.DB.GetBillingAccount(&account)
//	if err != nil {
//		return Drop500Error(c, err)
//	}
//
//	tranes, err := h.DB.GetAllTransactions(user.BillingID)
//	if err != nil {
//		return Drop500Error(c, err)
//	}
//
//	c.Status(200)
//	return c.JSON(BillingAccountResponse{
//		Code:         200,
//		Balance:      account.Balance,
//		Transactions: tranes,
//	})
//}
