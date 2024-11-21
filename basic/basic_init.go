package basic

import (
	iface "basic/interfaces"
	"extension/batch_retry"
	"extension/command_exec"
	"extension/db_conn_num"
	"extension/db_tool"
	"extension/edb_retry"
	"extension/high_availability"
	"extension/oneclick_inspection"
	"extension/redis_tool"
	"extension/svc_check"
	"extension/svc_switch"
	"extension/tcp_num"
	"extension/uri_switch"
	"extension/web_svc"
)

func InitGlobalContext(ctx iface.Context) {
	if globalContext == nil {
		globalContext = ctx.(*Context)
		globalContext.cache = make(map[string]map[string]any)
		globalContext.components = make(map[string]*componentMeta)
	}
}

func InitComponent() {
	addInitComponent(svc_check.GetInstance(globalContext))
	addInitComponent(command_exec.GetInstance(globalContext))
	addInitComponent(db_conn_num.GetInstance(globalContext))
	addInitComponent(db_tool.GetInstance(globalContext))
	addInitComponent(redis_tool.GetInstance(globalContext))
	addInitComponent(web_svc.GetInstance(globalContext))
	addInitComponent(oneclick_inspection.GetInstance(globalContext))
	addInitComponent(tcp_num.GetInstance(globalContext))
	addInitComponent(uri_switch.GetInstance(globalContext))
	addInitComponent(svc_switch.GetInstance(globalContext))
	addInitComponent(batch_retry.GetInstance(globalContext))
	addInitComponent(edb_retry.GetInstance(globalContext))
	addInitComponent(high_availability.GetInstance(globalContext))
}
