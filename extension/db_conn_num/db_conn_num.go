package db_conn_num

import (
	dbtool "basic/tool/db"
	"strconv"
)

const name = "db_conn_num"

var params_map = map[string]string{
	"host":     "string",
	"port":     "int",
	"username": "string",
	"password": "string",
	"dbname":   "string",
}

func Register() (key string, f func(req map[string]any) (resp []byte), metaData any) {
	return name, doHandler, nil
}

func doHandler(req map[string]any) (resp []byte) {
	host := req["host"].(string)
	port := req["port"].(int)
	username := req["username"].(string)
	password := req["password"].(string)
	dbname := req["dbname"].(string)
	dbBase, err := dbtool.CreateBaseDbByDbConfig(&dbtool.DbConfig{Host: host, Port: port, Username: username, Password: password, Dbname: dbname})
	if err != nil {
		return []byte("connect database fail, " + err.Error())
	}
	var count int64
	dbBase.Table("pg_stat_activity").Count(&count)
	return []byte("The current number of database connections is " + strconv.FormatInt(count, 10))
}
