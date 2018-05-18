package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Lofanmi/v3nayou/cfg"
	"github.com/Lofanmi/v3nayou/utils"

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
			map[string]string{"name": "四六级", "link": "http://m.tool.finded.net/index.php?m=Cet&school=gzhu", "icon": "/static/icons/cet.png"},
			map[string]string{"name": "实时公交", "link": "http://m.tool.finded.net/index.php?m=Bus&school=gzhu", "icon": "/static/icons/bus.png"},
			map[string]string{"name": "教学评价", "link": "[empty]", "icon": "/static/icons/teacher.png"},
			map[string]string{"name": "图书", "link": "http://m.book.finded.net/?code=1-3xwYVV6uCCyWxeNU-1xBtPHA-2x-2x", "icon": "/static/icons/books.png"},
		},
		"gzhu.icons.others": []map[string]string{
			map[string]string{"name": "常用电话", "link": "http://m.tool.finded.net/index.php?m=Baike&c=Contacts&a=index", "icon": "/static/icons/tel.png"},
			map[string]string{"name": "错误反馈", "link": "bug", "icon": "/static/icons/bug.png"},
		},
		"gzhu.ads": []map[string]string{
			map[string]string{"name": "参与开发", "link": "dev", "img": "/static/join.png"},
		},
		"gzhu.links": []map[string]string{
			map[string]string{"name": "哪有", "link": "http://mp.weixin.qq.com/mp/getmasssendmsg?__biz=MzA5NDQ1MTkyNA==#wechat_webview_type=1&wechat_redirect"},
			map[string]string{"name": "有独", "link": "http://mp.weixin.qq.com/mp/getmasssendmsg?__biz=MzA5NTg5NzExMg==#wechat_webview_type=1&wechat_redirect"},
			map[string]string{"name": "种草时间", "link": "http://mp.weixin.qq.com/mp/getmasssendmsg?__biz=MzI3NTA4MDQzMw==#wechat_webview_type=1&wechat_redirect"},
			map[string]string{"name": "关于我们", "link": "http://m.tool.finded.net/index.php?m=Home&c=Article&a=index&id=42"},
		},
		// 中大
		"sysu.name":       "中山大学",
		"sysu.name.short": "中大",
		"sysu.icons.mainly": []map[string]string{
			map[string]string{"name": "成绩", "link": "http://wjw.sysu.edu.cn/", "icon": "/static/icons/score.png"},
			map[string]string{"name": "课表", "link": "http://wjw.sysu.edu.cn/", "icon": "/static/icons/schedule.png"},
			map[string]string{"name": "校历", "link": "https://mp.weixin.qq.com/s/BekgSO8SpFJvjv0LuDZoJQ", "icon": "/static/icons/calendar.png"},
			map[string]string{"name": "个人中心", "link": "[empty]", "icon": "/static/icons/ucenter.png"},
			map[string]string{"name": "四六级", "link": "http://m.tool.finded.net/index.php?m=Cet&school=sysu", "icon": "/static/icons/cet.png"},
			map[string]string{"name": "实时公交", "link": "http://m.tool.finded.net/index.php?m=Bus&school=sysu", "icon": "/static/icons/bus.png"},
			map[string]string{"name": "教学评价", "link": "[empty]", "icon": "/static/icons/teacher.png"},
			map[string]string{"name": "图书", "link": "http://m.book.finded.net/?code=n1-1xbIFlI7rZNcsKYqhqw5g-2x-2x", "icon": "/static/icons/books.png"},
		},
		"sysu.icons.others": []map[string]string{
			map[string]string{"name": "公选排行榜", "link": "http://www.courstack.com/course/sysu", "icon": "/static/icons/chart.png"},
			map[string]string{"name": "错误反馈", "link": "bug", "icon": "/static/icons/bug.png"},
		},
		"sysu.ads": []map[string]string{
			map[string]string{"name": "参与开发", "link": "dev", "img": "/static/join.png"},
		},
		"sysu.links": []map[string]string{
			map[string]string{"name": "哪有", "link": "http://mp.weixin.qq.com/mp/getmasssendmsg?__biz=MzA5NDQ1MTkyNA==#wechat_webview_type=1&wechat_redirect"},
			map[string]string{"name": "有独", "link": "http://mp.weixin.qq.com/mp/getmasssendmsg?__biz=MzA5NTg5NzExMg==#wechat_webview_type=1&wechat_redirect"},
			map[string]string{"name": "种草时间", "link": "http://mp.weixin.qq.com/mp/getmasssendmsg?__biz=MzI3NTA4MDQzMw==#wechat_webview_type=1&wechat_redirect"},
			map[string]string{"name": "关于我们", "link": "http://m.tool.finded.net/index.php?m=Home&c=Article&a=index&id=42"},
		},
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
	// 文章列表
	articles, _ := getArticlesFromDB(school)
	m["articles"] = articles

	// 友情链接
	m["links"] = config(school + ".links")

	// 广大的学生未绑定教务系统, 链接到个人中心绑定教务.
	// 中大的学生无需绑定, 直接跳转至官方微教务查询.
	if school == "gzhu" {
		// 取出微信中间件的用户信息, 判断是否绑定了手机号和教务系统.
		member := c.MustGet("member").(*cfg.Member)
		sid, psw := member.EduMajor(school)
		if sid == "" || psw == "" {
			// 成绩
			m["mainly_icons"].([]interface{})[0].(map[string]interface{})["link"] = "ucenter"
			// 课表
			m["mainly_icons"].([]interface{})[1].(map[string]interface{})["link"] = "ucenter"
		}
	}

	c.JSON(http.StatusOK, utils.GinJSONData(http.StatusOK, m, "成功拉取数据"))

	return
}

func config(key string) interface{} {
	if _, ok := cfgs[key]; !ok {
		return nil
	}

	copy := make(map[string]interface{})
	bytes, _ := json.Marshal(cfgs)
	json.Unmarshal(bytes, &copy)

	v, _ := copy[key]

	return v
}

func getArticlesFromDB(school string) (result []map[string]interface{}, err error) {
	var (
		rows                               *sql.Rows
		title, abstract, cover, link, date string
		id, atype                          int
	)

	if school == "gzhu" {
		school = "1"
	} else if school == "sysu" {
		school = "7"
	}

	db := cfg.GetDB()

	rows, err = db.Query("SELECT id,title,abstract,cover,link,date,atype FROM `articles` WHERE `school` = ? AND `status` = 1 ORDER BY `sort`", school)
	defer rows.Close()
	if err != nil {
		return
	}

	for rows.Next() {
		err = rows.Scan(&id, &title, &abstract, &cover, &link, &date, &atype)
		if err != nil {
			return
		}
		if atype == 1 {
			link = fmt.Sprintf("http://m.tool.finded.net/index.php?c=Article&id=%d", id)
		}
		m := map[string]interface{}{
			"id":       id,
			"title":    title,
			"abstract": abstract,
			"cover":    cover,
			"link":     link,
			"date":     date,
		}
		result = append(result, m)
	}

	return
}
