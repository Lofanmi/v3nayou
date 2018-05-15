package routes

import (
	"os"

	"github.com/gin-gonic/gin"
)

// DB DB
var DB = make(map[string]string)

// SetupRouter 路由注册
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Home 首页
	r.GET("/school/:id", WechatMiddleware, Home)
	// Start 应用启动
	r.GET("/api/start", WechatMiddleware, SchoolMiddleware, Start)
	// UCenter 个人中心
	r.GET("/api/ucenter", WechatMiddleware, SchoolMiddleware, UCenter)
	// SendSMS 发送短信验证码
	r.GET("/api/sendsms/:tel", WechatMiddleware, SchoolMiddleware, SendSMS)
	// Member 绑定教务手机号
	r.POST("/api/member", WechatMiddleware, SchoolMiddleware, Member)
	// Schedule 课表
	r.GET("/api/schedule/:type", WechatMiddleware, SchoolMiddleware, Schedule)
	// Score 成绩
	r.GET("/api/score/:type/:xn/:xq", WechatMiddleware, SchoolMiddleware, Score)
	// TeacherStatus 教学评价状态
	r.GET("/api/teacher/status", WechatMiddleware, SchoolMiddleware, TeacherStatus)
	// TeacherComment 教学评价
	r.GET("/api/teacher/comment", WechatMiddleware, SchoolMiddleware, TeacherComment)

	authorized := r.Group("/admin/", gin.BasicAuth(gin.Accounts{
		os.Getenv("ADMIN_USER"): os.Getenv("ADMIN_PASS"),
	}))
	authorized.POST("login", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)
		c.JSON(200, gin.H{"message": user})
	})

	return r
}
