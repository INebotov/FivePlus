package web

import (
	"Backend/Errors"
	"Backend/Kafka"
	"Backend/Metrics"
	"Backend/auth"
	"Backend/db"
	"fmt"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

type Router struct {
	app *gin.Engine

	DataBase *db.DB
	Auth     *auth.Auth
	Metrics  *Metrics.Metrics
	Log      *zap.Logger

	Kafka struct {
		Email *Kafka.Kafka
		SMS   *Kafka.Kafka
		Media *Kafka.Kafka
	}

	EmailConfExpired time.Duration
	PhoneConfExpired time.Duration

	CloseFunc func() error

	Host string
	Port int

	Error Errors.HttpErrors

	pathAccess map[string][]int
}

func (r *Router) Init() error {
	r.app = gin.New()

	r.pathAccess = make(map[string][]int)

	r.app.Use(ginzap.Ginzap(r.Log, time.RFC3339, false))
	r.app.Use(r.RequestCounter)
	r.app.Use(r.UnseriousIDScanning)
	r.app.Use(r.RequestLog) // Only debug!
	r.app.Use(r.Refresh)
	r.app.Use(r.AccessCheck)

	r.Error = Errors.Init(Errors.HttpErrors{
		Metrics: r.Metrics,
		Log:     r.Log,
	})

	r.CloseFunc = func() error { return nil }
	return nil
}
func (r *Router) Start() error {
	return r.app.Run(fmt.Sprintf("%s:%d", r.Host, r.Port))
}

func (r *Router) ApplyHandlers() error {
	// Access: 1-Child 2-General 3-Teacher 4-Supporter 5-Admin
	r.applyMetrics()

	//2)/self + /self/confirm
	//3)/child, /children
	//4)/child/{id}/statistic
	//5)/lesson
	//6)Чат (/messages/{chat_id}) + другие пути для чата
	//7)/rooms, /lessons, /rooms/{id}, /lessons/{id}
	//8)Функции преподов (/teacher/tobe, /teacher/lesson, /teacher/lesson/{action})
	//9)Биллинг
	//10)/subjects, /grades, /teachers
	//11)/userlogo/{by}/value}, /person/{by}/{value}
	//11)Приглашения (/invlink, /accept, /invite)
	//11)Работа с проблемами (/problems, /support)
	//12)Админские функции (/admin/teacher, /admin/grade, /admin/subject
	//13)/mark/{what}

	r.applyAuth()

	return nil
}
func (r *Router) applyMetrics() {
	r.app.GET("/metrics", gin.WrapH(r.Metrics.GetHandler()))
}
func (r *Router) applyAuth() {
	// Auth
	r.pathAccess["/auth/reg"] = []int{0}
	r.app.POST("/auth/reg", r.Registration)

	r.pathAccess["/auth/login"] = []int{0}
	r.app.POST("/auth/login", r.Login)

	r.pathAccess["/auth/logout"] = []int{1, 2, 3, 4, 5}
	r.app.GET("/auth/logout", r.Logout)
}
func (r *Router) applySelf() {
	// Auth
	r.pathAccess["/auth/reg"] = []int{0}
	r.app.POST("/auth/reg", r.Registration)

	r.pathAccess["/auth/login"] = []int{0}
	r.app.POST("/auth/login", r.Login)

	r.pathAccess["/auth/logout"] = []int{1, 2, 3, 4, 5}
	r.app.GET("/auth/logout", r.Logout)
}
