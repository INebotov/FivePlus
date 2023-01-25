package old_router

//
//import (
//	"BackendSimple/db"
//	"fmt"
//	"github.com/gofiber/fiber/v2"
//	"strconv"
//)
//
//type AddSubjectReq struct {
//	ID   uint   `json:"id"`
//	Name string `json:"name"`
//}
//
//type RequestsToSend struct {
//	Code     int          `json:"code"`
//	Requests []RequestRes `json:"requests"`
//}
//type RequestRes struct {
//	ID          uint               `json:"id"`
//	TimeCreated int64              `json:"time_created"`
//	Student     UserConfidentional `json:"student"`
//
//	Subject  db.Subject `json:"subject"`
//	Question string     `json:"question"`
//}
//
//func (h Handlers) MarkActive(c *fiber.Ctx) error {
//	id := c.Locals("userid").(uint)
//	if id == 0 {
//		return Drop400Error(c)
//	}
//
//	tid, err := h.DB.GetTeacherIDByUserId(id)
//	if err != nil {
//		return Drop500Error(c, err)
//	}
//
//	active := c.Query("active") == "true"
//
//	err = h.DB.ChangeActive(tid, active)
//	if err != nil {
//		return Drop500Error(c, err)
//	}
//	c.Status(200)
//	return c.JSON(StandartResponse{
//		Code:    200,
//		Message: fmt.Sprintf("Successfully marked as Active=%v", active),
//	})
//}
//
//func (h Handlers) AddSubject(c *fiber.Ctx) error {
//	id := c.Locals("userid").(uint)
//	if id == 0 {
//		return Drop400Error(c)
//	}
//
//	var u AddSubjectReq
//	if err := c.BodyParser(&u); err != nil {
//		return Drop400Error(c)
//	}
//
//	techID, err := h.DB.GetTeacherIDByUserId(id)
//	if err != nil {
//		return Drop500Error(c, err)
//	}
//
//	if u.ID == 0 {
//		u.ID, err = h.DB.GetSubjectIDFromName(u.Name)
//		if err != nil {
//			return Drop500Error(c, err)
//		}
//	}
//
//	err = h.DB.TeacherAddSubject(techID, u.ID)
//	if err != nil {
//		return Drop500Error(c, err)
//	}
//	c.Status(200)
//
//	return c.JSON(StandartResponse{
//		Code:    200,
//		Message: "Successfully added!",
//	})
//}
//func (h Handlers) DeleteSubject(c *fiber.Ctx) error {
//	id := c.Locals("userid").(uint)
//	if id == 0 {
//		return Drop400Error(c)
//	}
//
//	var u AddSubjectReq
//	if err := c.BodyParser(&u); err != nil {
//		return Drop400Error(c)
//	}
//
//	techID, err := h.DB.GetTeacherIDByUserId(id)
//	if err != nil {
//		return Drop500Error(c, err)
//	}
//
//	if u.ID == 0 {
//		u.ID, err = h.DB.GetSubjectIDFromName(u.Name)
//		if err != nil {
//			return Drop500Error(c, err)
//		}
//	}
//	err = h.DB.TeacherDeleteSubject(techID, u.ID)
//	if err != nil {
//		return Drop500Error(c, err)
//	}
//	c.Status(200)
//
//	return c.JSON(StandartResponse{
//		Code:    200,
//		Message: "Successfully deleted!",
//	})
//}
//
//func (h Handlers) fillRequests(r []db.LessonRequest) ([]RequestRes, error) {
//	res := make([]RequestRes, len(r))
//	for i, e := range r {
//		var student db.User
//		student.ID = e.UserID
//		err := h.DB.GetUser(&student)
//		if err != nil {
//			return nil, err
//		}
//
//		stud := UserConfidentional{
//			Name:  student.Name,
//			Email: student.Email,
//
//			Role:     student.LevelString,
//			Photo:    student.Photo,
//			UserName: student.UserName,
//		}
//
//		var subject db.Subject
//		subject.ID = e.SubjectID
//		err = h.DB.GetSubject(&subject)
//		if err != nil {
//			return nil, err
//		}
//
//		res[i] = RequestRes{
//			ID:          e.ID,
//			TimeCreated: e.CreatedAt.Unix(),
//			Student:     stud,
//			Subject:     subject,
//			Question:    e.Question,
//		}
//	}
//	return res, nil
//}
//func (h Handlers) GetRecomendedSessions(c *fiber.Ctx) error {
//	id := c.Locals("userid").(uint)
//	if id == 0 {
//		return Drop400Error(c)
//	}
//
//	limit, err := strconv.Atoi(c.Query("limit", "none"))
//	if err != nil || limit > 100 {
//		limit = 25
//	}
//
//	teachID, err := h.DB.GetTeacherIDByUserId(id)
//	if id == 0 {
//		return Drop400Error(c)
//	}
//
//	teacher := db.Teacher{ID: teachID}
//	err = h.DB.GetTeacher(&teacher)
//	if err != nil {
//		return Drop500Error(c, err)
//	}
//	subjectIds := make([]uint, len(teacher.Subjects))
//	for i, el := range teacher.Subjects {
//		subjectIds[i] = el.ID
//	}
//
//	requests, err := h.DB.GetAllPendingLessonRequests(subjectIds, limit)
//	if err != nil {
//		return Drop500Error(c, err)
//	}
//	res, err := h.fillRequests(requests)
//	if err != nil {
//		return Drop500Error(c, err)
//	}
//	c.Status(200)
//	return c.JSON(RequestsToSend{
//		Code:     200,
//		Requests: res,
//	})
//}
//
//type StartLessonReq struct {
//	LessonReqID uint `json:"lesson_req_id"`
//}
//type StartLessonResponse struct {
//	Code    int    `json:"code"`
//	Message string `json:"message"`
//
//	ChatToken string `json:"chat_token"`
//	ChatID    string `json:"chat_id"`
//	LessonID  uint   `json:"lesson_id"`
//}
//
//func (h Handlers) StartLesson(c *fiber.Ctx) error {
//	id := c.Locals("userid").(uint)
//	if id == 0 {
//		return Drop400Error(c)
//	}
//
//	teachID, err := h.DB.GetTeacherIDByUserId(id)
//	if err != nil {
//		return Drop400Error(c)
//	}
//
//	var u StartLessonReq
//	if err := c.BodyParser(&u); err != nil || u.LessonReqID == 0 {
//		return Drop400Error(c)
//	}
//
//	r := db.LessonRequest{
//		ID: u.LessonReqID,
//	}
//	err = h.DB.GetLessonRequest(&r)
//	if err != nil {
//		return Drop500Error(c, err)
//	}
//
//	if r.Satisfied {
//		return Drop400Error(c)
//	}
//
//	chat := h.Chat.CreateRoom(r.Question, id, r.UserID)
//
//	l := db.Lesson{
//		Name:      r.Question,
//		SubjectID: r.SubjectID,
//		TeacherID: teachID,
//		StudentID: r.UserID,
//		ChatID:    chat.ID,
//	}
//	err = h.DB.CreateLesson(&l, &r)
//	if err != nil {
//		return Drop500Error(c, err)
//	}
//
//	room := db.ChatRoom{
//		ID: chat.ID,
//	}
//	err = h.DB.CreateRoom(&room)
//	if err != nil {
//		return Drop500Error(c, err)
//	}
//	token, err := h.Auth.GenChatToken(id, chat.ID)
//	if err != nil {
//		return Drop500Error(c, err)
//	}
//	c.Status(200)
//	return c.JSON(StartLessonResponse{
//		Code:      200,
//		Message:   "Successfully started lesson!",
//		ChatToken: token,
//		ChatID:    room.ID,
//		LessonID:  l.ID,
//	})
//}
