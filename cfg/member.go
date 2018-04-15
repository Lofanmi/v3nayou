package cfg

import (
	"encoding/json"
	"time"

	"github.com/Lofanmi/v3nayou/utils"
)

// Member 会员信息
type Member struct {
	ID        int    `json:"id"`
	OpenID    string `json:"openid"`
	Name      string `json:"name"`
	Tel       string `json:"tel"`
	Email     string `json:"email"`
	Edu       string `json:"edu"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`

	m *map[string]map[string]string
}

// EduMajor 获取主修专业账号
func (member *Member) EduMajor(school string) (sid, psw string) {
	if member.m == nil {
		member.m = new(map[string]map[string]string)
		json.Unmarshal([]byte(member.Edu), member.m)
	}
	sid, _ = (*member.m)[school]["sid"]
	psw, _ = (*member.m)[school]["psw"]

	psw = utils.Decrypt(psw)

	return
}

// EduSecond 获取辅修专业账号
func (member *Member) EduSecond(school string) (sid2, psw2 string) {
	if member.m == nil {
		member.m = new(map[string]map[string]string)
		json.Unmarshal([]byte(member.Edu), member.m)
	}
	sid2, _ = (*member.m)[school]["sid2"]
	psw2, _ = (*member.m)[school]["psw2"]

	psw2 = utils.Decrypt(psw2)

	return
}

// CreatedTs 获取创建时间戳
func (member *Member) CreatedTs() int64 {
	t, _ := time.Parse(TimeLayout, member.CreatedAt)
	return t.Unix()
}

// UpdatedTs 获取更新时间戳
func (member *Member) UpdatedTs() int64 {
	t, _ := time.Parse(TimeLayout, member.UpdatedAt)
	return t.Unix()
}
