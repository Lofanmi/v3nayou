package gzhu

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Lofanmi/v3nayou/utils"

	"github.com/parnurzeal/gorequest"
)

// Login 登录教务系统
func Login(sid, psw string, request *gorequest.SuperAgent) (bool, string, error) {
	var (
		url, body string
		resp      *http.Response
		errs      []error
	)

	url = "https://cas.gzhu.edu.cn/cas_server/login?service=http%3a%2f%2f202.192.18.182%2fLogin_gzdx.aspx"
	resp, body, errs = request.
		Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8").
		Set("Accept-Language", "zh-CN,zh;q=0.8,zh-TW;q=0.6").
		Set("Accept-Encoding", "gzip").
		Set("Cache-Control", "no-cache").
		Set("Connection", "keep-alive").
		Set("Pragma", "no-cache").
		Set("User-Agent", "Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/38.0.2125.101 Safari/537.36").
		Get(url).
		End()
	if errs != nil {
		return false, "", errors.New("Oh~教务系统挂了")
	}
	request.AddCookies(resp.Cookies())

	lt := utils.StrCut(body, `="lt" value="`, `"`)
	execution := utils.StrCut(body, `="execution" value="`, `"`)
	if lt == "" || execution == "" {
		return false, "", errors.New("连接教务系统失败")
	}

	_, body, _ = request.
		Type(gorequest.TypeForm).
		SendMap(map[string]string{
			"username":  sid,
			"password":  psw,
			"lt":        lt,
			"execution": execution,
			"captcha":   "",
			"warn":      "true",
			"_eventId":  "submit",
			"submit":    "登录",
		}).
		Post(url).
		End()
	body = utils.GBK2Str(body)

	success := strings.Contains(body, `id="xhxm">`+sid)
	name := utils.StrCut(utils.StrCut(body, `id="xhxm">`, "<"), "  ", "同学")

	return success, name, nil
}
