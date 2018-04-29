package routes

import (
	"net/http"

	"github.com/Lofanmi/v3nayou/cfg"
	"github.com/Lofanmi/v3nayou/utils"

	"github.com/gin-gonic/gin"
)

// UCenter 个人中心路由
func UCenter(c *gin.Context) {
	m := make(map[string]interface{})

	school := c.MustGet("school").(string)
	member := c.MustGet("member").(*cfg.Member)

	// 手机号
	m["tel"] = member.Tel
	// 邮箱
	m["email"] = member.Email
	// OpenID
	m["openid"] = member.OpenID
	// 主修
	m["sid"], m["psw"] = member.EduMajor(school)
	// 辅修
	m["sid2"], m["psw2"] = member.EduSecond(school)

	// 微信昵称
	m["nickname"] = member.WechatAuthObj().Nickname
	// 头像
	m["headimgurl"] = member.WechatAuthObj().HeadImgURL

	c.JSON(http.StatusOK, utils.GinJSONData(http.StatusOK, m, cfg.MessageSuccess))

	return
}
