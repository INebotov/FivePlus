package old_router

//
//import (
//	"BackendSimple/db"
//	"github.com/gofiber/fiber/v2"
//)
//
//type CreateSubjectReq struct {
//	Name        string `json:"name"`
//	Photo       string `json:"photo"`
//	Description string `json:"description"`
//
//	StaticID string `json:"staticID"`
//}
//type MakeUserTeacherReq struct {
//	UserName string `json:"user_name"`
//	ID       uint   `json:"id"`
//}
//
//func (h Handlers) CreateSubject(c *fiber.Ctx) error {
//	var req CreateSubjectReq
//	if err := c.BodyParser(&req); err != nil {
//		return Drop500Error(c, err)
//	}
//	err := h.DB.CreateSubject(&db.Subject{
//		Name:        req.Name,
//		Photo:       req.Photo,
//		Description: req.Description,
//		StaticID:    req.StaticID,
//	})
//	if err != nil {
//		return Drop500Error(c, err)
//	}
//
//	c.Status(200)
//	return c.JSON(StandartResponse{
//		Code:    200,
//		Message: "Successfully created subject",
//	})
//}
//func (h Handlers) MakeUserTeacher(c *fiber.Ctx) error {
//	var req MakeUserTeacherReq
//	if err := c.BodyParser(&req); err != nil || (req.ID == 0 && req.UserName == "") {
//		return Drop400Error(c)
//	}
//	var err error
//	if req.ID == 0 {
//		req.ID, err = h.DB.GetIDByUsername(req.UserName)
//		if err != nil {
//			return Drop500Error(c, err)
//		}
//	}
//
//	err = h.DB.MakeUserTeacher(req.ID)
//	if err != nil {
//		return Drop500Error(c, err)
//	}
//
//	c.Status(200)
//	return c.JSON(StandartResponse{
//		Code:    200,
//		Message: "Successfully added teacher!",
//	})
//}
