package old_router

//
//import (
//	"Backend/db"
//	"github.com/gofiber/fiber/v2"
//)
//
//type Router struct {
//	App *fiber.App
//}
//
//func GetRouter() (r Router, err error) {
//	r.App = fiber.New()
//	return r, nil
//}
//func (r Router) ApplyHandlers(h Handlers, m Middleware) error {
//	auth := r.App.Group("/auth")
//	auth.Post("/reg", h.Reg)
//	auth.Post("/login", h.Login)
//	auth.Get("/logout", h.Logout)
//
//	r.App.Get("/is", h.Is)
//
//	self := r.App.Group("/self", m.GetAccessCheck([]int{db.ChildAccessLevel, db.GeneralAccessLevel, db.TeacherAccessLevel, db.SupporterAccessLevel, db.AdminAccessLevel}))
//	self.Get("/hi", SecureHi)
//
//	self.Get("/", h.Self)
//	self.Post("/change", h.SelfChange)
//
//	self.Get("/change/send", h.SelfChangeSendCode)
//	self.Get("/change/token", h.SelfChangeGetTokenFromCode)
//	self.Post("/change/secrets", h.SelfSecretsChange)
//
//	self.Get("/confirm/send", h.SelfConfirmSend)
//	self.Get("/confirm", h.SelfConfirm)
//
//	self.Get("/lessons", h.GetLessons)
//	self.Get("/chats", h.GetChats)
//	self.Get("/chats/:id", h.GetMessages)
//
//	// self.Get("/search/teacher/subject", h.GetTeachersBySubject)
//
//	child := r.App.Group("/child", m.GetAccessCheck([]int{db.ChildAccessLevel}))
//	child.Post("/lesson/request", h.RequestALesson)
//	child.Get("/lesson/check", h.GetRequestStatus)
//	child.Post("/lesson/mark", h.MarkLesson)
//
//	parent := r.App.Group("/parent", m.GetAccessCheck([]int{db.GeneralAccessLevel}))
//	parent.Post("/child/create", h.CreateChild)
//	parent.Get("/child/delete/send", h.DeleteChildSend)
//	parent.Post("/child/delete", h.DeleteChild)
//	parent.Get("/billing", h.GetBillingAccount)
//
//	teacher := r.App.Group("/teacher", m.GetAccessCheck([]int{db.TeacherAccessLevel}))
//	teacher.Get("/active", h.MarkActive)
//	teacher.Post("/subject/add", h.AddSubject)
//	teacher.Post("/subject/delete", h.DeleteSubject)
//	teacher.Get("/recommended", h.GetRecomendedSessions)
//	teacher.Post("/lesson/start", h.StartLesson)
//
//	admin := r.App.Group("/admin", m.GetAccessCheck([]int{db.AdminAccessLevel}))
//	admin.Post("/create/subject", h.CreateSubject)
//	admin.Post("/make/teacher", h.MakeUserTeacher)
//	return nil
//}
//func SecureHi(c *fiber.Ctx) error {
//	c.Status(200)
//	id := c.Locals("userid").(uint)
//	return c.JSON(struct {
//		UserID    uint `json:"user_id"`
//		PrivateOK bool `json:"private_ok"`
//	}{id, true})
//}
