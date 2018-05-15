package gzhu

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Lofanmi/v3nayou/utils"

	"github.com/parnurzeal/gorequest"
)

// Schedule 获取课表
func Schedule(sid string, request *gorequest.SuperAgent) (map[string]string, error) {
	var (
		url, body, data string
		r               *regexp.Regexp
	)

	url = "http://202.192.18.182/tjkbcx.aspx?xh=" + sid

	_, body, _ = request.Get(url).End()
	body = utils.GBK2Str(body)
	if body == "" {
		return nil, fmt.Errorf("查不到课表~请检查官方教务是否已更新~")
	}

	r, _ = regexp.Compile(`selected="selected"[^>]+>(.*?)<`)
	list := r.FindAllStringSubmatch(body, -1)
	if len(list) < 6 {
		return nil, fmt.Errorf("获取课表失败")
	}
	m := map[string]string{
		"xn":      string(list[0][1]),
		"xq":      string(list[1][1]),
		"grade":   string(list[2][1]),
		"academy": string(list[3][1]),
		"major":   string(list[4][1]),
		"clazz":   string(list[5][1]),
	}

	data = `<table id="Table6"` +
		utils.StrClean(utils.StrCut(body, `<table id="Table6"`, `</table>`)) +
		"</table>"
	data = strings.Replace(data, ` id="Table6"`, ` class="schedule-table"`, -1)
	data = utils.StrReplace(data, map[string]string{"星期": ""})

	r, _ = regexp.Compile(` align="Center"`)
	data = r.ReplaceAllString(data, "")
	r, _ = regexp.Compile(`<tr><td colspan="2">早晨</td>.+?</tr>`)
	data = r.ReplaceAllString(data, "")
	r, _ = regexp.Compile(` width="\d+%"`)
	data = r.ReplaceAllString(data, "")

	data = strings.Replace(data, "#<br>", "", -1)

	m["data"] = r.ReplaceAllString(data, "")

	return m, nil
}
