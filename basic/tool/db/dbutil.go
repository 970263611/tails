package dbtool

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type BaseDb struct {
	*gorm.DB
}

type DbConfig struct {
	Host       string
	Port       int
	Username   string
	Password   string
	Dbname     string
	SearchPath string `default:public`
}

/*
*
通过配置信息创建数据库连接
*/
func CreateBaseDbByDbConfig(dbConfig *DbConfig) (*BaseDb, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable search_path=%s TimeZone=Asia/Shanghai",
		dbConfig.Host, dbConfig.Port, dbConfig.Username, dbConfig.Password, dbConfig.Dbname, dbConfig.SearchPath)
	return CreateBaseDbByDsn(dsn)
}

/*
*
通过字符串类型创建数据库连接
*/
func CreateBaseDbByDsn(dsn string) (*BaseDb, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	} else {
		return &BaseDb{db}, nil
	}

}
