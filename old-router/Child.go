package old_router

//
//import (
//	"BackendSimple/db"
//	"github.com/gofiber/fiber/v2"
//	"strconv"
//)
//
//type LessonMark struct {
//	LessonID uint `json:"lesson_id"`
//
//	Mark    int    `json:"mark"`
//	Comment string `json:"comment"`
//}
//
//type RequestLesson struct {
//	SubjectID uint   `json:"subject_id"`
//	Question  string `json:"question"`
//}
//type ResponseLessonRequest struct {
//	Code      int    `json:"code"`
//	Message   string `json:"message"`
//	RequestID uint   `json:"request_id"`
//
//	Satisfied bool `json:"satisfied"`
//
//	ChatToken string `json:"chat_token"`
//	ChatID    string `json:"chat_id"`
//	LessonID  uint   `json:"lesson_id"`
//}
//
//func (h Handlers) RequestALesson(c *fiber.Ctx) error {
//	id := c.Locals("userid").(uint)
//	if id == 0 {
//		return Drop400Error(c)
//	}
//
//	var u RequestLesson
//	if err := c.BodyParser(&u); err != nil {
//		return Drop400Error(c)
//	}
//
//	less := db.LessonRequest{
//		UserID:    id,
//		SubjectID: u.SubjectID,
//		Question:  u.Question,
//	}
//	err := h.DB.CreateLessonRequest(&less)
//	if err != nil {
//		return Drop500Error(c, err)
//	}
//	c.Status(200)
//	return c.JSON(ResponseLessonRequest{
//		Code:      200,
//		RequestID: less.ID,
//		Message:   "Successfully created lesson request",
//		Satisfied: less.Satisfied,
//	})
//}
//func (h Handlers) GetRequestStatus(c *fiber.Ctx) error {
//	id := c.Locals("userid").(uint)
//	if id == 0 {
//		return Drop400Error(c)
//	}
//	reqID, err := strconv.ParseUint(c.Query("id", "none"), 10, 32)
//	if err != nil {
//		return Drop400Error(c)
//	}
//
//	req := db.LessonRequest{ID: uint(reqID)}
//	err = h.DB.GetLessonRequest(&req)
//	if err != nil {
//		return Drop400Error(c)
//	}
//
//	if !req.Satisfied {
//		c.Status(200)
//		return c.JSON(ResponseLessonRequest{
//			Code:      200,
//			RequestID: req.ID,
//			Satisfied: false,
//		})
//	}
//
//	l := db.Lesson{
//		ID: req.LessonID,
//	}
//	err = h.DB.GetLesson(&l)
//	if err != nil {
//		return Drop500Error(c, err)
//	}
//
//	token, err := h.Auth.GenChatToken(id, l.ChatID)
//	if err != nil {
//		return Drop500Error(c, err)
//	}
//
//	c.Status(200)
//	return c.JSON(ResponseLessonRequest{
//		Code:      200,
//		Message:   "Lesson has been accepted!",
//		RequestID: req.ID,
//		Satisfied: true,
//		ChatToken: token,
//		ChatID:    l.ChatID,
//		LessonID:  l.ID,
//	})
//}
//
//func (h Handlers) MarkLesson(c *fiber.Ctx) error {
//	level := c.Locals("level").(int)
//	if level == 0 {
//		return Drop400Error(c)
//	}
//
//	id := c.Locals("userid").(uint)
//	if id == 0 {
//		return Drop400Error(c)
//	}
//
//	if level != db.ChildAccessLevel {
//		return Drop401Error(c)
//	}
//	var req LessonMark
//	if err := c.BodyParser(&req); err != nil {
//		return Drop500Error(c, err)
//	}
//
//	if req.Mark > 5 || req.Mark < 0 {
//		return Drop400Error(c)
//	}
//
//	l := db.Lesson{
//		ID:        req.LessonID,
//		StudentID: id,
//	}
//	err := h.DB.GetLesson(&l)
//	if err != nil {
//		return Drop500Error(c, err)
//	}
//	comment := db.Comment{
//		UserID:  id,
//		Mark:    req.Mark,
//		Message: req.Comment,
//	}
//	err = h.DB.MarkLesson(req.LessonID, comment)
//	if err != nil {
//		return Drop500Error(c, err)
//	}
//	c.Status(200)
//	return c.JSON(StandartResponse{
//		Code:    200,
//		Message: "Successfully marked lesson!",
//	})
//}
