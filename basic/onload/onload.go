package onload

import (
	"basic"
	"extension/batch_server"
	"extension/check_server"
	"extension/command_exec"
	"extension/db_conn_num"
	"extension/edb_server"
	"extension/monitor_server"
	"extension/nacos_server"
	"extension/redis_server"
	"extension/sql_server"
	"extension/suspend_server"
	"extension/tcp_server"
	"extension/web_server"
)

func init() {
	basic.AddInitComponent(check_server.GetInstance())
	basic.AddInitComponent(command_exec.GetInstance())
	basic.AddInitComponent(db_conn_num.GetInstance())
	basic.AddInitComponent(sql_server.GetInstance())
	basic.AddInitComponent(redis_server.GetInstance())
	basic.AddInitComponent(web_server.GetInstance())
	basic.AddInitComponent(monitor_server.GetInstance())
	basic.AddInitComponent(tcp_server.GetInstance())
	basic.AddInitComponent(suspend_server.GetInstance())
	basic.AddInitComponent(nacos_server.GetInstance())
	basic.AddInitComponent(batch_server.GetInstance())
	basic.AddInitComponent(edb_server.GetInstance())

}
