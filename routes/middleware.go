package routes

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Lofanmi/v3nayou/cfg"
	"github.com/Lofanmi/v3nayou/utils"

	"github.com/gin-gonic/gin"
	"github.com/parnurzeal/gorequest"
)

// SchoolMiddleware 学校中间件
func SchoolMiddleware(c *gin.Context) {
	school, _ := c.Cookie("school")
	if school == "" || (school != "gzhu" && school != "sysu") {
		// 错误页: 访问出错, 请从微信公众号进入
		c.Set("school", "[empty]")
		c.Next()
		return
	}
	c.Set("school", school)
	c.Next()
}

func loadFromCookie(c *gin.Context) (*cfg.Member, error) {
	o, _ := c.Cookie("o")
	o = utils.Decrypt(o)
	return getMemberByOpenID(o)
}

// WechatMiddleware 微信登录中间件
func WechatMiddleware(c *gin.Context) {
	var (
		member       *cfg.Member
		openid, auth string
		err          error
	)

	member, _ = loadFromCookie(c)

	if member == nil || member.WechatAuthObj() == nil {
		if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
			c.JSON(http.StatusForbidden, nil)
			return
		}
		code := c.Query("code")
		if code == "" {
			c.Redirect(http.StatusFound, getRedirectURL(c))
			return
		}
		openid, auth, err = getWechatAuth(code)
		if err != nil {
			c.Redirect(http.StatusFound, getRedirectURL(c))
			return
		}
		member, err = storeMember(openid, auth)
		if err != nil {
			c.Redirect(http.StatusFound, getRedirectURL(c))
			return
		}
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	if openid != "" {
		utils.SetCookie(c, "o", openid, 3600*24*365*5)
	}
	c.Set("member", member)
	c.Next()

	return
}

func getRedirectURL(c *gin.Context) string {
	return fmt.Sprintf(
		"https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_userinfo&state=1#wechat_redirect",
		os.Getenv("WECHAT_APPID"),
		utils.GetFullURL(c.Request),
	)
}

// https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421140842
func getWechatAuth(code string) (openid, auth string, err error) {
	var (
		url  string
		body []byte
		m    map[string]interface{}
	)

	url = fmt.Sprintf(
		"https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		os.Getenv("WECHAT_APPID"),
		os.Getenv("WECHAT_APPSECRET"),
		code,
	)

	r := gorequest.New()

	m = make(map[string]interface{})
	_, body, _ = r.Get(url).EndStruct(&m)

	if _, ok := m["errcode"]; ok {
		err = errors.New(m["errmsg"].(string))
		return
	}

	openid = m["openid"].(string)

	url = fmt.Sprintf(
		"https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN",
		m["access_token"],
		openid,
	)
	m = make(map[string]interface{})
	_, body, _ = r.Get(url).EndStruct(&m)
	auth = string(body)

	return
}

func getMemberByOpenID(openid string) (member *cfg.Member, err error) {
	var (
		id         int
		name       string
		tel        string
		email      string
		edu        string
		wechatauth string
		createdAt  string
		updatedAt  string
	)

	if openid == "" {
		return nil, nil
	}

	db := cfg.GetDB()
	row := db.QueryRow(
		"SELECT id,name,tel,email,edu,wechatauth,created_at,updated_at FROM `members` WHERE `openid` = ?",
		openid,
	)

	err = row.Scan(
		&id,
		&name,
		&tel,
		&email,
		&edu,
		&wechatauth,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return
	}

	member = &cfg.Member{
		ID:         id,
		OpenID:     openid,
		Name:       name,
		Tel:        tel,
		Email:      email,
		Edu:        edu,
		WechatAuth: wechatauth,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}

	return
}

func storeMember(openid, auth string) (member *cfg.Member, err error) {
	var (
		stmt *sql.Stmt
	)

	member, err = getMemberByOpenID(openid)

	t := time.Now().Format(cfg.TimeLayout)
	db := cfg.GetDB()

	if member == nil {
		stmt, err = db.Prepare("INSERT INTO `members` (openid,edu,wechatauth,created_at,updated_at) values(?,?,?,?,?)")
		if err != nil {
			return
		}

		_, err = stmt.Exec(openid, "", auth, t, t)
		if err != nil {
			return
		}

		member, err = getMemberByOpenID(openid)
		return
	}

	stmt, err = db.Prepare("UPDATE `members` SET wechatauth=?, updated_at=? WHERE openid=?")
	if err != nil {
		return
	}

	_, err = stmt.Exec(auth, t, openid)
	if err != nil {
		return
	}

	member, err = getMemberByOpenID(openid)
	return
}
