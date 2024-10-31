package command_exec

import (
	"basic"
	commandtool "basic/tool/command"
	"math"
	"runtime"
)

type CommandExec struct{}

func GetInstance() basic.Component {
	return &CommandExec{}
}

func (c CommandExec) GetOrder() int {
	return math.MaxInt64
}

func (c CommandExec) Register(globalContext *basic.Context) *basic.ComponentMeta {
	meta := &basic.ComponentMeta{
		Key:       "command_exec",
		Describe:  "",
		Component: c,
	}
	return meta
}

func (c CommandExec) Do(commands []string) (resp []byte) {
	//commandStr := req["command"].(string)
	commandStr := ""
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
