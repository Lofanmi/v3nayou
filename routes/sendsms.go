package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/Lofanmi/v3nayou/utils"

	"github.com/gin-gonic/gin"
	"github.com/parnurzeal/gorequest"
)

// SendSMS 发送短信验证码 5分钟内有效
func SendSMS(c *gin.Context) {
	// 手机号
	tel, err := checkEmpty(c.Param("tel"), "手机号")
	if err != nil {
		output(c, err.Error(), -1, nil)
		return
	}
	if err = checkTel(tel); err != nil {
		output(c, err.Error(), -1, nil)
		return
	}

	code := makeCode()

	if err = sendsms(tel, code); err != nil {
		output(c, err.Error(), -1, nil)
		return
	}

	// 验证信息加密并进cookie中
	tel, code, expire := toCookie(tel, code)

	maxAge := 5 * 60
	utils.SetCookie(c, "t", tel, maxAge)
	utils.SetCookie(c, "c", code, maxAge)
	utils.SetCookie(c, "e", expire, maxAge)

	output(c, "发送成功", 0, nil)
	return
}

// [0, 9999)
func makeCode() string {
	rand.Seed(time.Now().UnixNano())
	x := rand.Intn(10000)
	return fmt.Sprintf("%04d", x)
}

// tel, code, expire
func toCookie(tel, code string) (t, c, e string) {
	t = utils.Encrypt(tel)
	c = utils.Encrypt(code)
	e = utils.Encrypt(
		fmt.Sprintf("%d", time.Now().Add(5*time.Minute).Unix()),
	)
	return
}

func sendsms(tel, code string) error {
	if os.Getenv("GIN_MODE") != "release" {
		return nil
	}

	m := map[string]string{
		"apikey": os.Getenv("SMS_APIKEY"),
		"mobile": tel,
		"text": fmt.Sprintf(
			"【哪有校园服务】本次操作的验证码是：%s，验证码5分钟内有效，如非本人操作，可不必理会。", code,
		),
	}

	url := "https://sms.yunpian.com/v2/sms/single_send.json"

	r := gorequest.New()
	_, body, _ := r.Post(url).
		Type(gorequest.TypeForm).
		SendMap(m).
		EndBytes()

	result := make(map[string]interface{})
	json.Unmarshal(body, &result)

	switch result["code"].(int) {
	case 0:
		return nil
	case 3:
		return errors.New("短信平台余额不足")
	case 8:
		return errors.New("发送过于频繁, 30秒后重试")
	case 9:
		return errors.New("发送过于频繁, 5分钟后重试")
	case 17:
		return errors.New("发送过于频繁, 为保护用户不被骚扰, 请明天重试")
	case 22:
		return errors.New("发送过于频繁, 1小时后重试")
	case 33:
		return errors.New("发送过于频繁, 操作非法")
	default:
		return errors.New(result["msg"].(string))
	}
}
