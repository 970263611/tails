package dbtool

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DbConfig struct {
	Host       string
	Port       int
	Username   string
	Password   string
	Dbname     string
	SearchPath string
}

func ConnectionByDbConfig(dbConfig DbConfig) BaseDb {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable search_path=%s TimeZone=Asia/Shanghai",
		dbConfig.Host, dbConfig.Port, dbConfig.Username, dbConfig.Password, dbConfig.Dbname, dbConfig.SearchPath)
	return ConnectionByDsn(dsn)
}

func ConnectionByDsn(dsn string) BaseDb {
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return BaseDb{db}
}
