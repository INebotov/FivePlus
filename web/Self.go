package web

import "github.com/gin-gonic/gin"

func (r *Router) Self(c *gin.Context) {
	c.String(200, "Not implemented %s", c.FullPath())
}

func (r *Router) Confirm(c *gin.Context) {
	c.String(200, "Not implemented %s", c.FullPath())
}

func (r *Router) Lessons(c *gin.Context) {
	c.String(200, "Not implemented %s", c.FullPath())
}

func (r *Router) Rooms(c *gin.Context) {
	c.String(200, "Not implemented %s", c.FullPath())
}

func (r *Router) InvLink(c *gin.Context) {
	c.String(200, "Not implemented %s", c.FullPath())
}

func (r *Router) Accept(c *gin.Context) {
	c.String(200, "Not implemented %s", c.FullPath())
}

func (r *Router) ToBeTeacher(c *gin.Context) {
	c.String(200, "Not implemented %s", c.FullPath())
}

func (r *Router) Billing(c *gin.Context) {
	c.String(200, "Not implemented %s", c.FullPath())
}
