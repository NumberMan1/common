package variable

import (
	"github.com/coocood/freecache"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SysConfig struct {
	Mysql struct {
		Host        string `yaml:"host"`
		Port        int    `yaml:"port"`
		User        string `yaml:"user"`
		Password    string `yaml:"password"`
		Database    string `yaml:"database"`
		DbConfig    string `yaml:"db_config"`
		TablePrefix string `yaml:"table_prefix"`
	} `yaml:"mysql"`
	Zap struct {
		Director    string `yaml:"director"`
		LogLevel    string `yaml:"log_level"`
		EncodeLevel string `yaml:"encode_level"`
		LogPrefix   string `yaml:"log_prefix"`
	} `yaml:"zap"`
	BufferSize int `yaml:"buffer_size"`
}

//全局变量

var (
	GDb    *gorm.DB
	Config *SysConfig
	Cache  *freecache.Cache
	Log    *zap.Logger
)
