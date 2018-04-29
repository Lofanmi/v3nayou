package routes

import (
	"crypto/subtle"
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/Lofanmi/v3nayou/cfg"
	"github.com/Lofanmi/v3nayou/spiders/gzhu"
	"github.com/Lofanmi/v3nayou/utils"

	"github.com/gin-gonic/gin"
	"github.com/parnurzeal/gorequest"
)

// Member 个人中心账号绑定
func Member(c *gin.Context) {
	var (
		name string
	)
	member := c.MustGet("member").(*cfg.Member)
	school := c.MustGet("school").(string)
	// 校验用户输入
	m, err := validate(c)
	if err != nil {
		output(c, err.Error(), -1, nil)
		return
	}
	// 校验短信验证码
	err = validateTel(c, m["tel"])
	if err != nil {
		output(c, err.Error(), -1, nil)
		return
	}
	err = validateCode(c, m["tcode"])
	if err != nil {
		output(c, err.Error(), -1, nil)
		return
	}
	err = validateExpire(c)
	if err != nil {
		output(c, err.Error(), -1, nil)
		return
	}
	// 登录验证
	name, err = tryLogin(m["sid"], m["psw"], "(主修)")
	if err != nil {
		output(c, err.Error(), -1, nil)
		return
	}
	if m["sid2"] != "" && m["psw2"] != "" {
		_, err = tryLogin(m["sid2"], m["psw2"], "(辅修)")
		if err != nil {
			output(c, err.Error(), -1, nil)
			return
		}
	}
	// 登录成功 入库
	err = storeToDB(m, member.OpenID, school, name)
	if err != nil {
		output(c, err.Error(), -1, nil)
		return
	}
	// 删除短信验证码的cookie
	utils.SetCookie(c, "t", "", -1)
	utils.SetCookie(c, "c", "", -1)
	utils.SetCookie(c, "e", "", -1)

	// 提交成功
	output(c, "提交成功", 0, nil)
	return
}

func validate(c *gin.Context) (m map[string]string, err error) {
	var (
		tel, tcode, email, sid, psw, sid2, psw2 string
	)
	m = make(map[string]string)
	// 手机号
	tel, err = checkEmpty(c.PostForm("tel"), "手机号")
	if err != nil {
		return nil, err
	}
	if err = checkTel(tel); err != nil {
		return nil, err
	}
	m["tel"] = tel
	// 短信验证码
	tcode, err = checkEmpty(c.PostForm("tcode"), "短信验证码")
	if err != nil {
		return nil, err
	}
	if err = checkTcode(tcode); err != nil {
		return nil, err
	}
	m["tcode"] = tcode
	// 邮箱
	email, err = checkEmpty(c.PostForm("email"), "邮箱")
	if err != nil {
		return nil, err
	}
	if err = checkEmail(email); err != nil {
		return nil, err
	}
	m["email"] = email
	// 学号(主修专业)
	sid, err = checkEmpty(c.PostForm("sid"), "学号")
	if err != nil {
		return nil, err
	}
	if err = checkSid(sid, "学号"); err != nil {
		return nil, err
	}
	m["sid"] = sid
	// 密码(主修专业)
	psw, err = checkEmpty(c.PostForm("psw"), "密码")
	if err != nil {
		return nil, err
	}
	m["psw"] = psw
	// 学号(辅修专业)
	sid2, _ = checkEmpty(c.PostForm("sid2"), "二专学号")
	if sid2 == "" {
		m["sid2"] = ""
		m["psw2"] = ""
	} else {
		if err = checkSid(sid2, "二专学号"); err != nil {
			return nil, err
		}
		m["sid2"] = sid2
		// 密码(辅修专业)
		psw2, err = checkEmpty(c.PostForm("psw2"), "二专密码")
		if err != nil {
			return nil, err
		}
		m["psw2"] = psw2
	}
	return m, nil
}

func validateTel(c *gin.Context, input string) error {
	var tel string

	tel, _ = c.Cookie("t")
	tel = utils.Decrypt(tel)

	if tel == "" {
		return errors.New("请获取手机验证码")
	}

	if subtle.ConstantTimeCompare([]byte(input), []byte(tel)) != 1 {
		return errors.New("请填写获取验证码的手机号")
	}

	return nil
}

func validateCode(c *gin.Context, input string) error {
	var code string

	code, _ = c.Cookie("c")
	code = utils.Decrypt(code)

	if code == "" {
		return errors.New("请获取短信验证码")
	}

	if subtle.ConstantTimeCompare([]byte(input), []byte(code)) != 1 {
		return errors.New("短信验证码不正确")
	}

	return nil
}

func validateExpire(c *gin.Context) error {
	var expire string

	expire, _ = c.Cookie("t")
	expire = utils.Decrypt(expire)

	if expire == "" {
		return errors.New("请重新获取验证码")
	}

	ts, err := strconv.ParseInt(expire, 10, 64)
	if err != nil {
		return errors.New("验证失败, 请重新获取短信")
	}

	if ts < time.Now().Unix() {
		return errors.New("验证码已过期, 请重新获取")
	}

	return nil
}

func tryLogin(sid, psw, extra string) (name string, err error) {
	var (
		success bool
	)
	r := gorequest.New().SetDoNotClearSuperAgent(true)
	success, name, err = gzhu.Login(sid, psw, r)
	if err != nil {
		return "", err
	}
	if !success {
		return "", errors.New("教务系统连接超时" + extra)
	}
	return name, nil
}

func storeToDB(m map[string]string, openid, school, name string) error {
	var (
		id  int
		edu string
	)

	db := cfg.GetDB()

	row := db.QueryRow("SELECT id,edu FROM `members` WHERE `openid` = ?", openid)
	err := row.Scan(&id, &edu)
	if err != nil && err.Error() != sql.ErrNoRows.Error() {
		return errors.New("[500]保存用户信息出错")
	}

	if err != nil && err.Error() == sql.ErrNoRows.Error() {
		return insertToDB(m, openid, school, name)
	}

	return updateToDB(m, openid, school, name, edu)
}

func insertToDB(m map[string]string, openid, school, name string) error {
	db := cfg.GetDB()

	stmt, err := db.Prepare("INSERT INTO `members` (openid, name, tel, email, edu) values (?, ?, ?, ?, ?)")
	if err != nil {
		return errors.New("[501]保存用户信息出错")
	}

	bytes, _ := json.Marshal(map[string]map[string]string{
		school: map[string]string{
			"sid":  m["sid"],
			"psw":  utils.Encrypt(m["psw"]),
			"sid2": m["sid2"],
			"psw2": utils.Encrypt(m["psw2"]),
		},
	})
	edu := string(bytes)

	_, err = stmt.Exec(openid, name, m["tel"], m["email"], edu)
	if err != nil {
		return errors.New("[502]保存用户信息出错")
	}

	return nil
}

func updateToDB(m map[string]string, openid, school, name, edu string) error {
	db := cfg.GetDB()

	stmt, err := db.Prepare("UPDATE `members` SET name=?, tel=?, email=?, edu=? WHERE openid=?")
	if err != nil {
		return errors.New("[503]保存用户信息出错")
	}

	e := make(map[string]map[string]string)
	json.Unmarshal([]byte(edu), &e)

	e[school] = map[string]string{
		"sid":  m["sid"],
		"psw":  utils.Encrypt(m["psw"]),
		"sid2": m["sid2"],
		"psw2": utils.Encrypt(m["psw2"]),
	}
	bytes, _ := json.Marshal(e)
	edu = string(bytes)

	_, err = stmt.Exec(name, m["tel"], m["email"], edu, openid)
	if err != nil {
		return errors.New("[504]保存用户信息出错")
	}

	return nil
}
