package routes

import (
	"fmt"
	"strconv"

	"github.com/Lofanmi/v3nayou/cfg"
	"github.com/Lofanmi/v3nayou/spiders/gzhu"

	"github.com/gin-gonic/gin"
	"github.com/parnurzeal/gorequest"
)

// Score 成绩
func Score(c *gin.Context) {
	var (
		sid, psw, xn, xq string
		success          bool
		err              error
		m                []map[string]string
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

	xnRaw := c.Param("xn")
	xqRaw := c.Param("xq")

	xn = xnRaw
	xq = xqRaw

	if xq == "all" {
		xq = ""
	}
	if xn == "all" {
		xn = ""
		xq = ""
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

	// 获取成绩
	m, err = gzhu.Score(sid, xn, xq, r)
	if err != nil {
		output(c, err.Error(), -1, nil)
		return
	}

	result := buildResult(sid, xnRaw, xqRaw, &m)

	output(c, "获取成绩成功", 0, result)
	return
}

func buildResult(sid, xn, xq string, mlist *[]map[string]string) (result map[string]interface{}) {
	result = make(map[string]interface{})

	result["type"] = "standardGP"

	result["xn"] = xn
	result["xq"] = xq

	xns := []string{}
	year, _ := strconv.Atoi("20" + sid[0:2])
	for i := 0; i < 5; i++ {
		xns = append(xns, fmt.Sprintf("%d-%d", year, year+1))
		year++
	}
	xns = append(xns, "all")

	result["xns"] = xns
	result["xqs"] = []string{"1", "2", "all"}

	result["list"] = *mlist

	return
}
