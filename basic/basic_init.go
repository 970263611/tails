package basic

import (
	iface "basic/interfaces"
	"extension/batch_retry"
	"extension/command_exec"
	"extension/db_conn_num"
	"extension/db_tool"
	"extension/edb_retry"
	"extension/high_availability"
	"extension/redis_tool"
	"extension/tcp_num"
	"extension/uri_switch"
)

/*
*
这里设置配置文件默认值
*/
var defaultParams = map[string]any{}

func InitComponent() {
	addInitComponent(server_check.GetInstance(globalContext))
	addInitComponent(command_exec.GetInstance(globalContext))
	addInitComponent(db_conn_num.GetInstance(globalContext))
	addInitComponent(db_tool.GetInstance(globalContext))
	addInitComponent(redis_tool.GetInstance(globalContext))
	addInitComponent(web_server.GetInstance(globalContext))
	addInitComponent(monitor_check.GetInstance(globalContext))
	addInitComponent(tcp_num.GetInstance(globalContext))
	addInitComponent(uri_switch.GetInstance(globalContext))
	addInitComponent(server_on_off.GetInstance(globalContext))
	addInitComponent(batch_retry.GetInstance(globalContext))
	addInitComponent(edb_retry.GetInstance(globalContext))
	addInitComponent(high_availability.GetInstance(globalContext))
}

func InitGlobalContext(ctx iface.Context) {
	if globalContext == nil {
		globalContext = ctx.(*Context)
		globalContext.cache = make(map[string]map[string]any)
		globalContext.components = make(map[string]*componentMeta)
	}
}
