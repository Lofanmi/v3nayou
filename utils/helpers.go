package utils

import (
	"net/http"
	"strings"

	"github.com/axgle/mahonia"
	"github.com/gin-gonic/gin"
)

// StrCut 截取字符串.
func StrCut(s, begin, end string) string {
	b := strings.Index(s, begin)
	if b == -1 {
		return ""
	}
	b += len(begin)
	e := strings.Index(s[b:], end)
	if e == -1 {
		return ""
	}
	e += b
	return s[b:e]
}

// Str2GBK UTF-8字符串转为GBK编码.
func Str2GBK(s string) string {
	return mahonia.NewEncoder("GBK").ConvertString(s)
}

// GBK2Str GBK字符串转为UTF-8编码.
func GBK2Str(s string) string {
	return mahonia.NewDecoder("GBK").ConvertString(s)
}

// StrMap2GBK 将map的所有字符串转为GBK编码.
// 由于正方教务系统还没有中文的key.
// 所以对value做转换就OK了.
func StrMap2GBK(m map[string]string) (result map[string]string) {
	result = map[string]string{}
	for k, v := range m {
		result[k] = Str2GBK(v)
	}
	return
}

// StrReplace 将字符串按照map提供的映射进行替换.
func StrReplace(s string, m map[string]string) string {
	for k, v := range m {
		s = strings.Replace(s, k, v, -1)
	}
	return s
}

// StrClean 剔除字符串的所有空白字符, 但会保留单个空格.
// 包含的字符有:
// "\r", "\n", "\t", 二个或三个空格.
func StrClean(s string) string {
	return StrReplace(s, map[string]string{
		"\r":  "",
		"\n":  "",
		"\t":  "",
		"  ":  "",
		"   ": "",
	})
}

// Encrypt 字符串加密
func Encrypt(s string) string {
	return s
}

// Decrypt 字符串解密
func Decrypt(s string) string {
	return s
}

// GetFullURL 获取
func GetFullURL(r *http.Request) string {
	scheme := "http://"
	if r.TLS != nil {
		scheme = "https://"
	}
	return strings.Join([]string{scheme, r.Host, r.RequestURI}, "")
}

// GinJSONData 返回统一的JSON响应数据
func GinJSONData(code int, data interface{}, message string) gin.H {
	return gin.H{
		"code":    code,
		"data":    data,
		"message": message,
	}
}

// SetCookie 设置Cookie
func SetCookie(c *gin.Context, name, value string, maxAge int) {
	c.SetCookie(name, value, maxAge, "/", "", false, true)
}
