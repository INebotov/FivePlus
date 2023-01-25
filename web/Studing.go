package web

import "github.com/gin-gonic/gin"

func (r *Router) RequestLesson(c *gin.Context) {
	c.String(200, "Not implemented %s", c.FullPath())
}
func (r *Router) CheckLessonRequest(c *gin.Context) {
	c.String(200, "Not implemented %s", c.FullPath())
}
func (r *Router) GetRecommendedLessons(c *gin.Context) {
	c.String(200, "Not implemented %s", c.FullPath())
}
func (r *Router) ActWithLesson(c *gin.Context) {
	c.String(200, "Not implemented %s", c.FullPath())
}
