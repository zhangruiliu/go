package g

import (
	"g/configure"
	"io"

	"github.com/jinzhu/gorm"
)

const SUCCESS = `{"result":"success"}`
const FAILURE = `{"result":"failure"}`

const SUCCESS_INFO = "success"
const FAILURE_INFO = "failure"

// 全局变量
var (
	// Writer 输出
	Writer io.WriteCloser
	// Configure 所有配置
	Configure configure.Configure
	DB0       *gorm.DB
	Sqls      map[string]string //sql 模板
	DBcheck   *gorm.DB
	Dbid      *gorm.DB
)
