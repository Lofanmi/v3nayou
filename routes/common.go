package routes

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/Lofanmi/v3nayou/utils"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func output(c *gin.Context, message string, code int, data interface{}) {
	if code == 0 {
		code = 200
	} else {
		code = 422
	}
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.IndentedJSON(
		http.StatusOK,
		utils.GinJSONData(code, data, message),
	)
}

func checkEmpty(value, fieldName string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", errors.Errorf("%s不能为空", fieldName)
	}
	return value, nil
}

func checkTel(value string) error {
	if !regexp.MustCompile("^1[0-9]{10}$").MatchString(value) {
		return errors.Errorf("手机号格式不正确")
	}
	return nil
}

func checkEmail(value string) error {
	// https://github.com/go-playground/validator/blob/v9/regexes.go
	emailRegexString := "^(?:(?:(?:(?:[a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:\\(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22)))@(?:(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
	if !regexp.MustCompile(emailRegexString).MatchString(value) {
		return errors.Errorf("邮箱格式不正确")
	}
	return nil
}

func checkTcode(value string) error {
	if !regexp.MustCompile("^[0-9]{4}$").MatchString(value) {
		return errors.Errorf("短信验证码格式不正确")
	}
	return nil
}

func checkSid(value, fieldName string) error {
	if len(value) < 10 {
		return errors.Errorf("%s格式不正确", fieldName)
	}
	if !regexp.MustCompile("^[a-zA-Z0-9]+$").MatchString(value) {
		return errors.Errorf("%s格式不正确", fieldName)
	}
	return nil
}
