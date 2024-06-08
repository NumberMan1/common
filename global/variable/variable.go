package variable

import (
	"github.com/coocood/freecache"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SysConfig struct {
	Mysql struct {
		Host        string `json:"host"`
		Port        int    `json:"port"`
		User        string `json:"user"`
		Password    string `json:"password"`
		Database    string `json:"database"`
		DbConfig    string `json:"db_config"`
		TablePrefix string `json:"table_prefix"`
	} `json:"mysql"`
	Zap struct {
		Director    string `json:"director"`
		LogLevel    string `json:"log_level"`
		EncodeLevel string `json:"encode_level"`
		LogPrefix   string `json:"log_prefix"`
	} `json:"zap"`
	BufferSize int `json:"buffer_size"`
}

//全局变量

var (
	GDb    *gorm.DB
	Config *SysConfig
	Cache  *freecache.Cache
	Log    *zap.Logger
)
