package onload

import (
	"basic"
	"extension/check_server"
	"extension/command_exec"
	"extension/db_conn_num"
	"extension/sql_server"
)

func init() {
	var list = make([]basic.Component, 0)
	list = append(list, check_server.GetInstance())
	list = append(list, command_exec.GetInstance())
	list = append(list, db_conn_num.GetInstance())
	list = append(list, sql_server.GetInstance())
	n := len(list)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if list[j].GetOrder() > list[j+1].GetOrder() {
				list[j], list[j+1] = list[j+1], list[j]
			}
		}
	}
	basic.Assemble(list)
}
