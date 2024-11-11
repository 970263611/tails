package onload

import (
	"basic"
	"extension/check_server"
	"extension/command_exec"
	"extension/db_conn_num"
	"extension/redis_server"
	"extension/sql_server"
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
	basic.AddInitComponent(tcp_server.GetInstance())
}
