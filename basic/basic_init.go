package basic

import (
	iface "basic/interfaces"
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

func InitGlobalContext(ctx iface.Context) {
	if globalContext == nil {
		globalContext = ctx.(*Context)
		globalContext.cache = make(map[string]map[string]any)
		globalContext.components = make(map[string]*componentMeta)
	}
}

func init() {
	addInitComponent(check_server.GetInstance(globalContext))
	addInitComponent(command_exec.GetInstance(globalContext))
	addInitComponent(db_conn_num.GetInstance(globalContext))
	addInitComponent(sql_server.GetInstance(globalContext))
	addInitComponent(redis_server.GetInstance(globalContext))
	addInitComponent(web_server.GetInstance(globalContext))
	addInitComponent(monitor_server.GetInstance(globalContext))
	addInitComponent(tcp_server.GetInstance(globalContext))
	addInitComponent(suspend_server.GetInstance(globalContext))
	addInitComponent(nacos_server.GetInstance(globalContext))
	addInitComponent(batch_server.GetInstance(globalContext))
	addInitComponent(edb_server.GetInstance(globalContext))
}
