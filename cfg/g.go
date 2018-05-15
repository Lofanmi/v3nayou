package cfg

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	// TimeLayout 时间格式化
	TimeLayout = "2006-01-02 15:04:05"
	// MessageSuccess 成功消息
	MessageSuccess = "success"
)

// LoadEnv 加载环境变量文件, 并检查必填的环境变量
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	must := []string{
		"DB_DSN",
		"TEST_SID",
		"TEST_PSW",
		"HOME_TPL",
		"ADMIN_USER",
		"ADMIN_PASS",
		"WECHAT_APPID",
		"WECHAT_APPSECRET",
		"SMS_APIKEY",
		"APP_KEY",
	}

	for _, key := range must {
		if os.Getenv(key) == "" {
			panic(fmt.Sprintf("Invalid env variable: %s", key))
		}
	}
}
