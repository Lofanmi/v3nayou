package routes

import (
	"encoding/json"

	"github.com/Lofanmi/v3nayou/cfg"
	"github.com/Lofanmi/v3nayou/spiders/gzhu"
	"github.com/Lofanmi/v3nayou/utils"

	"github.com/gin-gonic/gin"
	"github.com/parnurzeal/gorequest"
)

// TeacherStatus 获取教学评价状态
func TeacherStatus(c *gin.Context) {
	var (
		sid, psw string
		success  bool
		err      error
		m        []map[string]string
		data     []byte
	)

	school := c.GetString("school")
	member := c.MustGet("member").(*cfg.Member)
	sid, psw = member.EduMajor(school)
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

	// 获取教学评价状态
	m, err = gzhu.Status(sid, r)
	if err != nil {
		output(c, err.Error(), -1, nil)
		return
	}

	data, err = json.Marshal(m)
	if err != nil {
		output(c, "无法编码教学评价状态", -1, nil)
		return
	}

	output(c, "成功获取教学评价状态", 0, map[string]string{
		"payload": utils.Encrypt(string(data)),
	})

	return
}

// TeacherComment 教学评价
func TeacherComment(c *gin.Context) {
	var (
		sid, psw, s string
		success     bool
		err         error
		m           []map[string]string
	)

	school := c.GetString("school")
	member := c.MustGet("member").(*cfg.Member)
	sid, psw = member.EduMajor(school)
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

	s = utils.Decrypt(c.PostForm("payload"))
	if s == "" {
		output(c, "评价数据无效", -1, nil)
	}
	err = json.Unmarshal([]byte(s), m)
	if err != nil {
		output(c, "无法读取评价数据", -1, nil)
		return
	}
	if len(m) == 0 {
		output(c, "教师列表为空", -1, nil)
		return
	}

	// 教学评价
	success, m, err = gzhu.Comment(sid, m, r)
	if err != nil {
		output(c, err.Error(), -1, nil)
		return
	}
	if !success {
		output(c, "评价失败", -1, m)
	}

	return
}
