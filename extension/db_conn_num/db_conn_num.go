package db_conn_num

import (
	cons "basic/constants"
	iface "basic/interfaces"
	dbtool "basic/tool/db"
	"basic/tool/utils"
	"errors"
	"strconv"
)

type DbConnNum struct{}

func GetInstance(globalContext iface.Context) iface.Component {
	return &DbConnNum{}
}

func (c DbConnNum) GetName() string {
	return "db_conn_num"
}

func (c DbConnNum) GetDescribe() string {
	return "数据库连接数查询,当前仅支持pg数据库"
}

func (d DbConnNum) Register(cm iface.ComponentMeta) {
	cm.AddParameters(cons.STRING, cons.LOWER_H, "db_conn_num.host", "host", true,
		func(s string) error {
			if !utils.CheckIp(s) {
				return errors.New("IP不合法")
			}
			return nil
		}, "IP地址,支持ipv4,例:192.168.0.1")
	cm.AddParameters(cons.INT, cons.LOWER_P, "db_conn_num.port", "port", true,
		func(s string) error {
			if !utils.CheckPortByString(s) {
				return errors.New("端口不合法")
			}
			return nil
		}, "端口号,0 ~ 65535之间")
	cm.AddParameters(cons.STRING, cons.LOWER_U, "db_conn_num.username", "username", true, nil, "数据库用户名")
	cm.AddParameters(cons.STRING, cons.LOWER_W, "db_conn_num.password", "password", true, nil, "数据库密码")
	cm.AddParameters(cons.STRING, cons.LOWER_D, "db_conn_num.dbname", "dbname", true, nil, "数据库名称")
	cm.AddParameters(cons.STRING, cons.LOWER_S, "db_conn_num.searchpath", "schema", true, nil, "数据库schema")
}

func (d DbConnNum) Start() error {
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
		return []byte("连接数据库失败: " + err.Error())
	}
	var count int64
	dbBase.Table("pg_stat_activity").Count(&count)
	return []byte("当前数据库连接数:" + strconv.FormatInt(count, 10))
}
