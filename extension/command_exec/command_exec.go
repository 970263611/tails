package command_exec

import (
	"basic"
	commandtool "basic/tool/command"
	log "github.com/sirupsen/logrus"
	"math"
	"runtime"
	"strings"
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
	pipeFlag := false
	length := len(commands)
	for i := 0; i < length; i++ {
		switch commands[i] {
		case "-c":
			i++
			if i < length {
				commandStr = commands[i]
			}
		case "-p":
			pipeFlag = true
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
		commandArr := strings.Split(commandStr, " ")
		result, err = commandtool.ExecCmdByStr(append([]string{"cmd", "/C"}, commandArr...)...)
	} else if pipeFlag {
		commandArr := make([][]string, 0)
		for _, command := range strings.Split(commandStr, "|") {
			commandArr = append(commandArr, strings.Split(strings.TrimSpace(command), " "))
		}
		result, err = commandtool.ExecCmdByPipe(commandArr...)
	} else {
		commandArr := strings.Split(commandStr, " ")
		result, err = commandtool.ExecCmdByStr(commandArr...)
	}
	if err != nil {
		log.Error("exec command '" + commandStr + "' error, errmsg: " + err.Error())
		return []byte("command exec error : " + err.Error())
	}
	return []byte(result)
}
