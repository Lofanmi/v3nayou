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
	r.GET("/start", WechatMiddleware, SchoolMiddleware, Start)
	// UCenter 个人中心
	r.GET("/ucenter", WechatMiddleware, SchoolMiddleware, UCenter)
	// SendSMS 发送短信验证码
	r.GET("/sendsms/:tel", WechatMiddleware, SchoolMiddleware, SendSMS)
	// Member 绑定教务手机号
	r.POST("/member", WechatMiddleware, SchoolMiddleware, Member)
	// Schedule 课表
	r.GET("/schedule/:type", WechatMiddleware, SchoolMiddleware, Schedule)
	// Score 成绩
	r.GET("/score/:type/:xn/:xq", WechatMiddleware, SchoolMiddleware, Score)
	// TeacherStatus 教学评价状态
	r.GET("/teacher/status", WechatMiddleware, SchoolMiddleware, TeacherStatus)
	// TeacherComment 教学评价
	r.GET("/teacher/comment", WechatMiddleware, SchoolMiddleware, TeacherComment)

	authorized := r.Group("/admin/", gin.BasicAuth(gin.Accounts{
		os.Getenv("ADMIN_USER"): os.Getenv("ADMIN_PASS"),
	}))
	authorized.POST("login", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)
		c.JSON(200, gin.H{"message": user})
	})

	return r
}
