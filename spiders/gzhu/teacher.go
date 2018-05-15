package gzhu

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/Lofanmi/v3nayou/utils"

	"github.com/parnurzeal/gorequest"
)

// Status 获取教学评价状态
func Status(sid string, request *gorequest.SuperAgent) ([]map[string]string, error) {
	var (
		url, body, data string
	)

	url = "http://202.192.18.182/xsjxpj2.aspx?xh=" + sid

	_, body, _ = request.Get(url).End()
	body = utils.GBK2Str(body)
	if body == "" {
		return nil, fmt.Errorf("评教系统无法连接")
	}
	if strings.Contains(body, "您已经评价过") {
		return nil, errors.New("已经评价好咯")
	}
	if strings.Contains(body, "评教系统已关闭") {
		return nil, errors.New("评教系统已关闭")
	}

	// 获取课程名称, 教师列表, 课程链接列表
	data = utils.StrCut(body, `<table class="datelist"`, `</table>`)
	data = strings.Replace(data, "【课堂教学】", "", -1)
	tr := strings.Split(data, "</tr>")

	status := []map[string]string{}

	for _, s := range tr {
		m := make(map[string]string)
		// 课程名称
		m["name"] = utils.StrCut(s, "<td>", "</td>")
		// 课程链接
		scriptName := strings.TrimLeft(utils.StrCut(s, "window.open('", "'"), "/")
		m["url"] = "http://202.192.18.182/" + scriptName
		// 授课教师
		m["teacher"] = utils.StrCut(s, "&nbsp;&nbsp;", "(")
		status = append(status, m)
	}

	return status, nil
}

// Comment 教学评价
func Comment(sid string, status []map[string]string, request *gorequest.SuperAgent) (bool, []map[string]string, error) {
	var (
		url, body, avg, s1, s2, s3, s4, comment, message string
		datelisthead                                     bool
	)

	if len(status) == 0 {
		return false, nil, fmt.Errorf("教师列表为空")
	}

	result := []map[string]string{}

	for _, m := range status {
		_, body, _ = request.Get(m["url"]).End()
		body = utils.GBK2Str(body)
		avg, s1, s2, s3, s4, comment = randComment()
		form := map[string]string{
			"__VIEWSTATE":          subViewState(body),
			"__VIEWSTATEGENERATOR": "C8894877",
			"Button1":              "保  存",
			"TextBox1":             "",
			// 随机生成的 4 个成绩
			"DataGrid1:_ctl2:txt_pf=": s1,
			"DataGrid1:_ctl3:txt_pf=": s2,
			"DataGrid1:_ctl4:txt_pf=": s3,
			"DataGrid1:_ctl5:txt_pf=": s4,
			// 随机挑选的评语
			"txt_pjxx": comment,
		}
		if strings.Contains(body, "__EVENTVALIDATION") {
			form["__EVENTVALIDATION"] = subEventValidation(body)
		}
		_, body, _ = request.
			Type(gorequest.TypeForm).
			SendMap(utils.StrMap2GBK(form)).
			Post(url).
			End()
		body = ""
		m := map[string]string{
			"avg":     avg,
			"s1":      s1,
			"s2":      s2,
			"s3":      s3,
			"s4":      s4,
			"comment": comment,
		}
		result = append(result, m)
	}

	// 验证是否评价成功
	url = "http://202.192.18.182/xsjxpj2.aspx?xh=" + sid
	_, body, _ = request.Get(url).End()
	body = utils.GBK2Str(body)
	form := map[string]string{
		"__VIEWSTATE":          subViewState(body),
		"__VIEWSTATEGENERATOR": "0008A89B",
		"__EVENTTARGET":        "",
		"__EVENTARGUMENT":      "",
		"btn_tj":               " 提 交 ",
	}
	if strings.Contains(body, "__EVENTVALIDATION") {
		form["__EVENTVALIDATION"] = subEventValidation(body)
	}
	_, body, _ = request.
		Type(gorequest.TypeForm).
		SendMap(utils.StrMap2GBK(form)).
		Post(url).
		End()

	message = utils.StrCut(body, "alert('", "'")
	datelisthead = strings.Contains(body, "datelisthead")

	return message == "完成评价！" && datelisthead, result, nil
}

var (
	comments = []string{
		"老师的方式非常适合我，根据本课程知识结构的特点，层次分明。理论和实际相结合，使知识更条理化。",
		"老师治学严谨，要求严格，循循善诱，平易近人；注意启发和调动学生的积极性，课堂气氛较为活跃。",
		"老师对待教学认真负责，语言生动，条理清晰，举例充分恰当。",
		"老师课堂内容充实，简单明了，使学生能够轻轻松松掌握知识，教学内容丰富有效。",
		"老师对待学生严格要求，能够鼓励学生踊跃发言，使课堂气氛比较积极热烈。",
		"老师课堂内容充实，简单明了，使学生能够轻松掌握知识。教学内容丰富，过程中尊重学生，很受同学欢迎。",
		"老师教学认真，课堂效率高，授课内容详细，我们大部分都能跟着老师思路，整节课学下来有收获、欣喜。",
		"老师常常理论联系实际，课上穿插实际问题，使同学们对自己所学专业有初步了解，为今后学习打下基础。",
		"老师治学严谨，对学生严格要求。课堂中，他循循善诱，强调独立思考，引导学生进行启发是思维。",
		"上课时，老师能够从学生实际出发，适当缓和课堂气氛，充分调动学生学习的积极性。",
		"老师上课诙谐有趣，他善于用凝练的语言将复杂难于理解的过程公式清晰、明确的表达出来。",
	}
)

// 随机生成分数和评语
// Intn(4)=>[0, 4) 生成96-99的随机数
func randComment() (avg, s1, s2, s3, s4, comment string) {
	// 初始化随机种子
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// 随机生成4个分数
	i1 := 96 + r.Intn(3)
	i2 := 96 + r.Intn(3)
	i3 := 96 + r.Intn(3)
	i4 := 96 + r.Intn(3)
	// 平均分
	avg = strconv.FormatFloat(float64(i1+i2+i3+i4)/4.0, 'f', 2, 64)
	if strings.Contains(avg, ".") {
		avg = strings.TrimRight(avg, "0")
		avg = strings.TrimRight(avg, ".")
	}
	s1 = fmt.Sprintf("%d", i1)
	s2 = fmt.Sprintf("%d", i2)
	s3 = fmt.Sprintf("%d", i3)
	s4 = fmt.Sprintf("%d", i4)
	// 随机评语
	comment = comments[rand.Intn(len(comments))]
	// 返回
	return
}
