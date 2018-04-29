package routes

import (
	"github.com/Lofanmi/v3nayou/cfg"
	"github.com/Lofanmi/v3nayou/spiders/gzhu"

	"github.com/gin-gonic/gin"
	"github.com/parnurzeal/gorequest"
)

// Schedule 课表
func Schedule(c *gin.Context) {
	var (
		sid, psw string
		success  bool
		err      error
		m        map[string]string
	)
	school := c.GetString("school")
	member := c.MustGet("member").(*cfg.Member)

	if c.Param("type") == "major" {
		sid, psw = member.EduMajor(school)
	} else {
		sid, psw = member.EduSecond(school)
	}

	if sid == "" {
		output(c, "账号尚未绑定", -1, nil)
		return
	}

	r := gorequest.New().SetDoNotClearSuperAgent(true)

	// 登录教务系统
	success, _, err = gzhu.Login(sid, psw, r)
	if err != nil {
		output(c, err.Error(), -1, nil)
		return
	}
	if !success {
		output(c, err.Error(), -1, nil)
		return
	}

	// 获取课表
	m, err = gzhu.Schedule(sid, r)
	if err != nil {
		output(c, err.Error(), -1, nil)
		return
	}

	output(c, "获取课表成功", 0, m)
	return
}
