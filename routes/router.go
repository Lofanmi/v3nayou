package routes

import (
	"github.com/gin-gonic/gin"
)

// SetupRouter 路由注册
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Home 首页
	r.GET("/school/:id", WechatMiddleware(), Home)

	api := r.Group("/api")
	api.
		Use(WechatMiddleware()).
		Use(SchoolMiddleware())
	{
		// Start 应用启动
		api.GET("/start", Start)
		// UCenter 个人中心
		api.GET("/ucenter", UCenter)
		// SendSMS 发送短信验证码
		api.GET("/sendsms/:tel", SendSMS)
		// Member 绑定教务手机号
		api.POST("/member", Member)
		// Schedule 课表
		api.GET("/schedule/:type", Schedule)
		// Score 成绩
		api.GET("/score/:type/:xn/:xq", Score)
		// TeacherStatus 教学评价状态
		api.GET("/teacher/status", TeacherStatus)
		// TeacherComment 教学评价
		api.GET("/teacher/comment", TeacherComment)
	}

	return r
}
