package routes

import (
	"net/http"

	"github.com/Lofanmi/v3nayou/cfg"
	"github.com/Lofanmi/v3nayou/utils"

	"github.com/gin-gonic/gin"
)

// Home 首页路由(公众号自定义菜单链接入口)
func Home(c *gin.Context) {
	var school string
	switch c.Param("id") {
	case "1":
		school = "gzhu"
	case "7":
		school = "sysu"
	default:
		school = ""
	}
	if school == "" {
		// 非法请求
		return
	}
	utils.SetCookie(c, "school", school, 3600*24*365*5)
	c.Data(http.StatusOK, "text/html; charset=UTF-8", cfg.HomeTmpl())
}
