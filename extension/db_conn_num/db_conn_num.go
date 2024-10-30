package db_conn_num

import (
	"basic"
	dbtool "basic/tool/db"
	"math"
	"strconv"
)

type DbConnNum struct{}

func GetInstance() basic.Component {
	return &DbConnNum{}
}

func (d DbConnNum) GetOrder() int {
	return math.MaxInt64
}

func (d DbConnNum) Register() *basic.ComponentMeta {
	meta := &basic.ComponentMeta{
		Key:       "db_conn_num",
		Describe:  "",
		Component: d,
	}
	return meta
}

func (d DbConnNum) Do(commands []string) (resp []byte) {
	//host := req["host"].(string)
	//port := req["port"].(int)
	//username := req["username"].(string)
	//password := req["password"].(string)
	//dbname := req["dbname"].(string)
	//dbBase, err := dbtool.CreateBaseDbByDbConfig(&dbtool.DbConfig{Host: host, Port: port, Username: username, Password: password, Dbname: dbname})
	dbBase, err := dbtool.CreateBaseDbByDbConfig(&dbtool.DbConfig{})
	if err != nil {
		return []byte("connect database fail, " + err.Error())
	}
	var count int64
	dbBase.Table("pg_stat_activity").Count(&count)
	return []byte("The current number of database connections is " + strconv.FormatInt(count, 10))
}
