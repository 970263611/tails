package command_exec

import (
	"basic"
	commandtool "basic/tool/command"
	log "github.com/sirupsen/logrus"
	"runtime"
	"strings"
)

type CommandExec struct{}

func GetInstance() basic.Component {
	return &CommandExec{}
}

func (c CommandExec) GetName() string {
	return "command_exec"
}

func (c CommandExec) GetDescribe() string {
	return "执行系统命令或者调用系统脚本"
}

func (c CommandExec) Register(globalContext *basic.Context) *basic.ComponentMeta {
	p1 := basic.Parameter{
		ParamType:    basic.STRING,
		CommandName:  "-c",
		StandardName: "command",
		Required:     true,
		Describe:     "",
	}
	return &basic.ComponentMeta{
		ComponentType: basic.EXECUTE,
		Component:     c,
		Params:        []basic.Parameter{p1},
	}
}

func (c CommandExec) Start(globalContext *basic.Context) error {
	return nil
}

func (c CommandExec) Do(params map[string]any) (resp []byte) {
	commandStr := params["command"].(string)
	var result string
	var err error
	if runtime.GOOS == "windows" {
		commandArr := strings.Split(commandStr, " ")
		result, err = commandtool.ExecCmdByStr(append([]string{"cmd", "/C"}, commandArr...)...)
	} else {
		commandArr := make([][]string, 0)
		for _, command := range strings.Split(commandStr, "|") {
			commandArr = append(commandArr, strings.Split(strings.TrimSpace(command), " "))
		}
		result, err = commandtool.ExecCmdByPipe(commandArr...)
	}
	if err != nil {
		log.Error("exec command '" + commandStr + "' error, errmsg: " + err.Error())
		return []byte("命令行执行失败 : " + err.Error())
	}
	return []byte(result)
}
