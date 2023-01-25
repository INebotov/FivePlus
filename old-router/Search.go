package old_router

//
//import (
//	"github.com/gofiber/fiber/v2"
//	"strconv"
//)
//
//func (h Handlers) GetTeachersBySubject(c *fiber.Ctx) error {
//	subject := c.Query("subject")
//	if subject == "" {
//		return Drop400Error(c)
//	}
//	limit, err := strconv.Atoi(c.Query("limit", "none"))
//	if err != nil {
//		limit = 100
//	}
//
//	id, err := h.DB.GetSubjectIDFromName(subject)
//	if err != nil {
//		return Drop500Error(c, err)
//	}
//
//	teachers, err := h.DB.GetTeacherBySubject(id, limit)
//	if err != nil {
//		return Drop500Error(c, err)
//	}
//	c.Status(200)
//	return c.JSON(parceTeachersToTeachersReaponse(teachers))
//}
