package db_conn_num

import (
	"basic"
	dbtool "basic/tool/db"
	othertool "basic/tool/other"
	"errors"
	"strconv"
)

type DbConnNum struct{}

func GetInstance() basic.Component {
	return &DbConnNum{}
}

func (c DbConnNum) GetName() string {
	return "db_conn_num"
}

func (c DbConnNum) GetDescribe() string {
	return "数据库连接数查询,当前仅支持pg数据库"
}

func (d DbConnNum) Register(globalContext *basic.Context) *basic.ComponentMeta {
	p1 := basic.Parameter{
		ParamType:    basic.STRING,
		CommandName:  "-h",
		StandardName: "host",
		Required:     true,
		CheckMethod: func(s string) error {
			if !othertool.CheckIp(s) {
				return errors.New("IP不合法")
			}
			return nil
		},
		Describe: "IP地址,支持ipv4,例:192.168.0.1",
	}
	p2 := basic.Parameter{
		ParamType:    basic.INT,
		CommandName:  "-p",
		StandardName: "port",
		Required:     true,
		CheckMethod: func(s string) error {
			if !othertool.CheckPortByString(s) {
				return errors.New("端口不合法")
			}
			return nil
		},
		Describe: "端口号,0 ~ 65535之间",
	}
	p3 := basic.Parameter{
		ParamType:    basic.STRING,
		CommandName:  "-u",
		StandardName: "username",
		Required:     true,
		Describe:     "数据库用户名",
	}
	p4 := basic.Parameter{
		ParamType:    basic.STRING,
		CommandName:  "-P",
		StandardName: "password",
		Required:     true,
		Describe:     "数据库密码",
	}
	p5 := basic.Parameter{
		ParamType:    basic.STRING,
		CommandName:  "-d",
		StandardName: "dbname",
		Required:     true,
		Describe:     "数据库名称",
	}
	p6 := basic.Parameter{
		ParamType:    basic.STRING,
		CommandName:  "-s",
		StandardName: "searchpath",
		Required:     true,
		Describe:     "数据库schema",
	}
	return &basic.ComponentMeta{
		ComponentType: basic.EXECUTE,
		Component:     d,
		Params:        []basic.Parameter{p1, p2, p3, p4, p5, p6},
	}
}

func (d DbConnNum) Start(globalContext *basic.Context) error {
	return nil
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
