package command_exec

import (
	"basic"
	commandtool "basic/tool/command"
	log "github.com/sirupsen/logrus"
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
	var commandStr string
	length := len(commands)
	for i := 0; i < length; i++ {
		switch commands[i] {
		case "-c":
			i++
			if i < length {
				commandStr = commands[i]
			}
		default:
			continue
		}
	}
	if commandStr == "" {
		return []byte("parameter error")
	}
	var result string
	var err error
	if runtime.GOOS == "windows" {
		result, err = commandtool.ExecCmdByStr("cmd", "/C", commandStr)
	} else {
		result, err = commandtool.ExecCmdByStr(commandStr)
	}
	if err != nil {
		log.Error("exec command '" + commandStr + "' error, errmsg: " + err.Error())
		return []byte("command exec error : " + err.Error())
	}
	return []byte(result)
}
