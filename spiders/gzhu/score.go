package gzhu

import (
	"fmt"
	"strings"

	"github.com/Lofanmi/v3nayou/utils"

	"github.com/parnurzeal/gorequest"
)

var scorekey = map[string]int{
	// 学年
	"xn": 0,
	// 学期
	"xq": 1,
	// 课程名称
	"name": 3,
	// 课程性质
	"cate": 4,
	// 学分
	"credit": 6,
	// 绩点
	"point": 7,
	// 成绩
	"score": 8,
	// 补考成绩
	"retake": 10,
	// 重修
	"restudy": 11,
}

// Score 获取成绩
func Score(sid, xn, xq string, request *gorequest.SuperAgent) ([]map[string]string, error) {
	var (
		url, body string
	)

	url = `http://202.192.18.182/xscj_gc.aspx?xh=` + sid

	_, body, _ = request.Get(url).End()
	body = utils.GBK2Str(body)

	form := map[string]string{
		"__VIEWSTATE":          subViewState(body),
		"__VIEWSTATEGENERATOR": subViewStateGenerator(body),
		"__EVENTVALIDATION":    subEventValidation(body),
		"Button1":              "按学期查询",
		"ddlXN":                xn,
		"ddlXQ":                xq,
	}

	_, body, _ = request.
		Type(gorequest.TypeForm).
		SendMap(utils.StrMap2GBK(form)).
		Post(url).
		End()

	body = utils.GBK2Str(body)
	if body == "" {
		return nil, fmt.Errorf(`尚未查询到成绩, 换个学期试试看~`)
	}

	data := strings.Split(utils.StrCut(body, `<td>重修标记</td>`, `</table>`), `</tr>`)

	detail := []map[string]string{}

	for i := 0; i < len(data)-2; i++ {
		m := make(map[string]string)

		td := strings.Split(data[i+1], `</td>`)

		for key, position := range scorekey {
			s := utils.StrReplace(
				td[position],
				map[string]string{
					"<td>":                       "",
					"&nbsp;":                     "",
					"<tr>\r\n\t\t":               "",
					"<tr class=\"alt\">\r\n\t\t": "",
				},
			)
			m[key] = utils.StrClean(s)
		}

		detail = append(detail, m)
	}

	if len(detail) <= 0 {
		return nil, fmt.Errorf(`尚未查询到成绩, 再查一次试试看?`)
	}

	return detail, nil
}
