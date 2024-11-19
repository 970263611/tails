package command_exec

import (
	cons "basic/constants"
	iface "basic/interfaces"
	commandtool "basic/tool/command"
	log "github.com/sirupsen/logrus"
	"runtime"
	"strings"
)

type CommandExec struct{}

func GetInstance(globalContext iface.Context) iface.Component {
	return &CommandExec{}
}

func (c CommandExec) GetName() string {
	return "command_exec"
}

func (c CommandExec) GetDescribe() string {
	return "执行系统命令行,例:command_exec -c 'ps -ef|grep java'"
}

func (c CommandExec) Register(cm iface.ComponentMeta) {
	cm.AddParameters(cons.STRING, "-c", "", "command", true, nil, "命令行,用单引号或双引号包裹;例如:'ll -al' 或 'ps -ef|grep java'")
}

func (c CommandExec) Start() error {
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
