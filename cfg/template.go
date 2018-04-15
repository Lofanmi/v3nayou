package cfg

import (
	"io/ioutil"
	"os"
)

var (
	tmpl []byte
)

// InitTmpl 初始化模板, 从硬盘中加载HTML.
func InitTmpl() {
	filename := os.Getenv("HOME_TPL")
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	tmpl = data
}

// HomeTmpl 获取首页模板.
func HomeTmpl() []byte {
	return tmpl
}
