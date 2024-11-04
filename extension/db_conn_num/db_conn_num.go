package db_conn_num

import (
	"basic"
	dbtool "basic/tool/db"
	othertool "basic/tool/other"
	"errors"
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

func (d DbConnNum) Register(globalContext *basic.Context) *basic.ComponentMeta {
	p1 := basic.Parameter{
		ParamType:    basic.STRING,
		CommandName:  "-h",
		StandardName: "host",
		Required:     true,
		CheckMethod: func(s string) error {
			if !othertool.CheckIp(s) {
				errors.New("ip is not valid")
			}
			return nil
		},
		Describe: "",
	}
	p2 := basic.Parameter{
		ParamType:    basic.INT,
		CommandName:  "-p",
		StandardName: "host",
		Required:     true,
		CheckMethod: func(s string) error {
			if !othertool.CheckPortByString(s) {
				errors.New("port is not valid")
			}
			return nil
		},
		Describe: "",
	}
	p3 := basic.Parameter{
		ParamType:    basic.STRING,
		CommandName:  "-u",
		StandardName: "username",
		Required:     true,
		Describe:     "",
	}
	p4 := basic.Parameter{
		ParamType:    basic.STRING,
		CommandName:  "-u",
		StandardName: "username",
		Required:     true,
		Describe:     "",
	}
	p5 := basic.Parameter{
		ParamType:    basic.STRING,
		CommandName:  "-w",
		StandardName: "password",
		Required:     true,
		Describe:     "",
	}
	p6 := basic.Parameter{
		ParamType:    basic.STRING,
		CommandName:  "-s",
		StandardName: "searchpath",
		Required:     true,
		Describe:     "",
	}
	return &basic.ComponentMeta{
		Key:       "db_conn_num",
		Describe:  "",
		Component: d,
		Params:    []basic.Parameter{p1, p2, p3, p4, p5, p6},
	}
}

func (d DbConnNum) Do(params map[string]any) (resp []byte) {
	config := &dbtool.DbConfig{
		Host:     params["host"].(string),
		Port:     params["port"].(int),
		Username: params["username"].(string),
		Password: params["password"].(string),
		Dbname:   params["dbname"].(string),
	}
	if params["searchpath"] != nil {
		config.SearchPath = params["searchpath"].(string)
	}
	dbBase, err := dbtool.CreateBaseDbByDbConfig(config)
	if err != nil {
		return []byte("connect database fail, " + err.Error())
	}
	var count int64
	dbBase.Table("pg_stat_activity").Count(&count)
	return []byte("The current number of database connections is " + strconv.FormatInt(count, 10))
}
