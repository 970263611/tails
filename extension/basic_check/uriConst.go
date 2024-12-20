package basic_check

import (
	cons "basic/constants"
	"errors"
	"fmt"
	"strings"
)

func (req *request) initResp() {
	m := make(map[string]*RespEntry)
	resp := Resp{
		ProjectName: req.projectName,
		Node:        []*RespEntry{},
	}
	adders := strings.Split(req.addr, ",")
	for i := range adders {
		addr := adders[i]
		parts := strings.Split(addr, ":")
		ip := parts[0]
		m[ip] = &RespEntry{
			Location:    req.groupName,
			ProjectName: req.projectName,
			Addr:        adders[i],
		}
		resp.Node = append(resp.Node, m[ip])
	}
	req.resultMap = m
	req.resp = resp
}

func (req *request) getLogPath() error {
	var logFilePath string
	groupName := req.groupName
	projectName := req.projectName
	switch groupName {
	case CPYY:
		logFilePath = fmt.Sprintf("%s/%s/%s", cpyyLogPath, projectName, logFileName)
	case GSCF:
		logFilePath = fmt.Sprintf("%s/%s/%s", gscfLogPath, projectName, logFileName)
	case GRCF:
		logFilePath = fmt.Sprintf("%s/%s/%s", grcfLogPath, projectName, logFileName)
	default:
		return errors.New("请输出正确的groupName 应该为:CPYY GSCF GRCF其中之一")
	}
	req.logFilePath = logFilePath
	return nil
}

func (req *request) getLog(basicServer *BasicServer) {
	tailsAdders := strings.Split(req.tailsAddr, ",")
	command := fmt.Sprintf("'tail -n 1000 %s | grep %s'", req.logFilePath, SUCCESS)
	for _, addr := range tailsAdders {
		//检查日志
		commands := []string{
			"command_exec",
			cons.LOWER_C.GetCommandName(),
			command,
			"--f",
			addr,
		}
		logInfo := getLogInfo(commands, basicServer)
		parts := strings.Split(addr, ":")
		port := parts[0]
		respEntry, ok := req.resultMap[port]
		if ok {
			respEntry.LogInfo = logInfo
		}
	}
}

func getLogInfo(commands []string, basicServer *BasicServer) string {
	args, err := basicServer.context.LoadSystemParams(commands)
	if err != nil {
		return "log 转发失败 加载参数失败"
	}
	resp, err := basicServer.context.Servlet(args, false)
	if err != nil {
		return err.Error()
	}
	return string(resp)
}
