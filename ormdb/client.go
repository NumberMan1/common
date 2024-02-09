package ormdb

import (
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func ConnectToDB(dataBase, dsn string) (db *gorm.DB, err error) {
	switch dataBase {
	case "mysql":
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	case "postgres":
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	case "sqlserver":
		db, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	case "tidb":
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	case "clickhouse":
		db, err = gorm.Open(clickhouse.Open(dsn), &gorm.Config{})
	}
	return
}
