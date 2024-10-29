package command_exec

import (
	commandtool "basic/tool/command"
	"runtime"
)

const name = "command_exec"

func Register() (key string, f func(req map[string]any) (resp []byte), metaData any) {
	return name, doHandler, nil
}

func doHandler(req map[string]any) (resp []byte) {
	commandStr := req["command"].(string)
	var result string
	var err error
	if runtime.GOOS == "windows" {
		result, err = commandtool.ExecCmdByStr("cmd", "/C", commandStr)
	} else {
		result, err = commandtool.ExecCmdByStr(commandStr)
	}
	if err != nil {
		return []byte("command exec error : " + err.Error())
	}
	return []byte(result)
}
