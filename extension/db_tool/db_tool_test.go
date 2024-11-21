package db_tool

import (
	"basic/tool/db"
	"fmt"
	"testing"
)

func TestDoHandler(t *testing.T) {
	params := map[string]any{
		"host":       "127.0.0.1",
		"port":       5432,
		"username":   "postgres",
		"password":   "postgre",
		"dbname":     "postgres",
		"searchPath": "public",
		"sql":        "INSERT INTO public.test (id, name, age,address) VALUES (555, 'zhh', 21,'北京丰台'),(666, 'zhh', 21,'北京丰台');",
		"--force":    "",
	}
	instance := GetInstance(nil)
	resp := instance.Do(params)
	if resp != nil {
		fmt.Println(string(resp))
	}
}

func GetDb() *dbtool.BaseDb {
	dbConfig := dbtool.DbConfig{Host: "127.0.0.1", Port: 5432, Username: "postgres", Password: "postgre", Dbname: "postgres", SearchPath: "public"}
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable search_path=%s TimeZone=Asia/Shanghai",
		dbConfig.Host, dbConfig.Port, dbConfig.Username, dbConfig.Password, dbConfig.Dbname, dbConfig.SearchPath)
	conn, error := dbtool.CreateBaseDbByDsn(dsn)
	if error != nil {
		return nil
	}
	return conn
}

func TestExecSqlFile(t *testing.T) {
	err := ExecSqlFile("D:\\test.sql", GetDb())
	if err != nil {
		t.Errorf("ExecSql Fail %e", err)
	}
}
