package cfg

import (
	"database/sql"
	"os"
	// 引入数据库驱动注册及初始化
	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB
)

// InitDB 初始化数据库连接
func InitDB() (err error) {
	db, err = sql.Open("mysql", os.Getenv("DB_DSN"))
	if err != nil {
		return
	}
	err = db.Ping()
	if err != nil {
		return
	}
	return
}

// GetDB 获取数据库连接
func GetDB() *sql.DB {
	return db
}

// FreeDB 释放数据库连接
func FreeDB() {
	db.Close()
}
