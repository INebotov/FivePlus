package web

import "github.com/gin-gonic/gin"

func (r *Router) CreateSubject(c *gin.Context) {
	c.String(200, "Not implemented %s", c.FullPath())
}
func (r *Router) CreateGrade(c *gin.Context) {
	c.String(200, "Not implemented %s", c.FullPath())
}
func (r *Router) GetTeacherRequests(c *gin.Context) {
	c.String(200, "Not implemented %s", c.FullPath())
}
func (r *Router) AcceptTeacherRequest(c *gin.Context) {
	c.String(200, "Not implemented %s", c.FullPath())
}
