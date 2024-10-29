package dbtool

import (
	"fmt"
	"testing"
)

func TestCreateBaseDbByDbConfig(t *testing.T) {
	dbConfig := &DbConfig{"127.0.0.1", 15431, "postgres", "postgres", "postgres", "public"}
	conn, err := CreateBaseDbByDbConfig(dbConfig)
	if err != nil {
		t.Logf("connet error %v", err.Error())
	}
	t.Logf("connet success %v", conn)
}

func TestCreateBaseDbByDsn(t *testing.T) {
	dbConfig := DbConfig{"127.0.0.1", 15431, "postgres", "postgres", "postgres", "public"}
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable search_path=%s TimeZone=Asia/Shanghai",
		dbConfig.Host, dbConfig.Port, dbConfig.Username, dbConfig.Password, dbConfig.Dbname, dbConfig.SearchPath)
	conn, _ := CreateBaseDbByDsn(dsn)
	t.Logf("connet success %v", conn)
}
