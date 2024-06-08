package global

import (
	"encoding/json"
	"fmt"
	"github.com/NumberMan1/common/global/variable"
	"github.com/NumberMan1/common/logger"
	"github.com/NumberMan1/common/ormdb"
	"github.com/coocood/freecache"
	"os"
)

func InitConf(configPath string) {
	fmt.Println("加载配置文件dev_config.json")
	file, _ := os.Open(configPath)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logger.SLoggerConsole.Error("Error closing")
		}
	}(file)
	decoder := json.NewDecoder(file)
	config := variable.SysConfig{}
	err := decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error:", err)
		panic("加载配置出错")
	}

	variable.Config = &config
}

func Init(configPath string) {
	InitConf(configPath)
	variable.Log = logger.Zap()
	variable.GDb = ormdb.InitDb()
	variable.Cache = freecache.NewCache(1024 * 1024 * 1024)
}
