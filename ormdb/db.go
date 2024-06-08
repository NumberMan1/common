package ormdb

import (
	"fmt"
	"github.com/NumberMan1/common/global/variable"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strconv"
)

// 注册数据库
func InitDb() *gorm.DB {
	fmt.Println("注册数据库")
	dsn := variable.Config.Mysql.User + ":" + variable.Config.Mysql.Password +
		"@tcp(" + variable.Config.Mysql.Host + ":" + strconv.FormatInt(int64(variable.Config.Mysql.Port), 10) + ")/" +
		variable.Config.Mysql.Database + "?" + variable.Config.Mysql.DbConfig
	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         191,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}
	config := &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true, NamingStrategy: schema.NamingStrategy{
		TablePrefix: variable.Config.Mysql.TablePrefix,
	}}
	if db, err := gorm.Open(mysql.New(mysqlConfig), config); err != nil {
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(0)
		sqlDB.SetMaxOpenConns(0)
		return db
	}
}
