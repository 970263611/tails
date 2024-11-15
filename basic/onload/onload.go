package onload

import (
	"basic"
	"extension/check_server"
	"extension/command_exec"
	"extension/db_conn_num"
	"extension/monitor_server"
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

}
