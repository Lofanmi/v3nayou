package routes

import (
	"encoding/json"
	"net/http"

	"github.com/Lofanmi/v3nayou/cfg"
	"github.com/Lofanmi/v3nayou/utils"
	"github.com/parnurzeal/gorequest"

	"github.com/gin-gonic/gin"
)

var (
	cfgs = map[string]interface{}{
		// 广大
		"gzhu.name":       "广州大学",
		"gzhu.name.short": "广大",
		"gzhu.icons.mainly": []map[string]string{
			map[string]string{"name": "成绩", "link": "score", "icon": "/static/icons/score.png"},
			map[string]string{"name": "课表", "link": "schedule", "icon": "/static/icons/schedule.png"},
			map[string]string{"name": "校历", "link": "http://mp.weixin.qq.com/s/WIn3rG1BjxRx3EdOn9zvrw", "icon": "/static/icons/calendar.png"},
			map[string]string{"name": "个人中心", "link": "ucenter", "icon": "/static/icons/ucenter.png"},
			map[string]string{"name": "四六级", "link": "[empty]", "icon": "/static/icons/cet.png"},
			map[string]string{"name": "时间胶囊", "link": "[empty]", "icon": "/static/icons/capsule.png"},
			map[string]string{"name": "教学评价", "link": "[empty]", "icon": "/static/icons/teacher.png"},
			map[string]string{"name": "图书", "link": "http://m.book.finded.net/?code=1-3xwYVV6uCCyWxeNU-1xBtPHA-2x-2x", "icon": "/static/icons/books.png"},
		},
		"gzhu.icons.others": []map[string]string{
			map[string]string{"name": "公选排行榜", "link": "[empty]", "icon": "/static/icons/charts.png"},
			map[string]string{"name": "常用电话", "link": "http://m.tool.finded.net/index.php?m=Baike&c=Contacts&a=index", "icon": "/static/icons/strategy.png"},
		},
		"gzhu.ads": []map[string]string{
			map[string]string{"name": "如何正确使用哪有?", "link": "#", "img": "http://m.nayou.finded.net/uploads/ads/529de50d95a28442dd5718532ba432835c144588.jpg"},
		},
		"gzhu.links": []map[string]string{
			map[string]string{"name": "哪有", "link": "	http://mp.weixin.qq.com/mp/getmasssendmsg?__biz=MzA5NDQ1MTkyNA==#wechat_webview_type=1&wechat_redirect"},
			map[string]string{"name": "有独", "link": "http://mp.weixin.qq.com/mp/getmasssendmsg?__biz=MzA5NTg5NzExMg==#wechat_webview_type=1&wechat_redirect"},
			map[string]string{"name": "种草时间", "link": "http://mp.weixin.qq.com/mp/getmasssendmsg?__biz=MzI3NTA4MDQzMw==#wechat_webview_type=1&wechat_redirect"},
			map[string]string{"name": "关于我们", "link": "#"},
		},
		// 中大
		// ...
	}
)

// Start 应用启动路由
func Start(c *gin.Context) {
	school := c.GetString("school")

	m := make(map[string]interface{})
	// 学校名称
	m["school"] = config(school + ".name.short").(string)
	// 图标
	m["mainly_icons"] = config(school + ".icons.mainly")
	m["others_icons"] = config(school + ".icons.others")
	// 广告位
	m["ads"] = config(school + ".ads")
	// 文章
	// TODO: 后期有时间再迁移到同一张数据表, 现在先调API吧
	// m["articles"] = []map[string]string{}
	_, body, _ := gorequest.New().Get("http://m.nayou.finded.net/api/articles/1").EndBytes()
	articles := []map[string]interface{}{}
	json.Unmarshal(body, &articles)
	m["articles"] = articles

	// 友情链接
	m["links"] = config(school + ".links")

	// 广大的学生未绑定教务系统, 链接到个人中心绑定教务.
	// 中大的学生无需绑定, 直接跳转至官方微教务查询.
	if school == "gzhu" {
		// 取出微信中间件的用户信息, 判断是否绑定了教务系统.
		member := c.MustGet("member").(*cfg.Member)
		sid, psw := member.EduMajor(school)
		if sid == "" || psw == "" {
			// 成绩
			m["mainly_icons"].([]map[string]string)[0]["link"] = "ucenter"
			// 课表
			m["mainly_icons"].([]map[string]string)[1]["link"] = "ucenter"
		}
	}

	c.JSON(http.StatusOK, utils.GinJSONData(http.StatusOK, m, cfg.MessageSuccess))

	return
}

func config(key string) interface{} {
	v, ok := cfgs[key]
	if !ok {
		return nil
	}
	return v
}
