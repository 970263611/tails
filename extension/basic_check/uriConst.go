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

func (req *request) getLogInfo() {
	defer req.deferMethod()
	tailsAdders := strings.Split(req.tailsAddr, ",")
	command := fmt.Sprintf("'tail -n 1000 %s | grep %s'", req.logFilePath, SUCCESS)
	for _, addr := range tailsAdders {
		commands := []string{
			"command_exec",
			cons.LOWER_C.GetCommandName(),
			command,
			"--f",
			addr,
		}
		logInfo := forward(commands, req.basicServer)
		parts := strings.Split(addr, ":")
		ip := parts[0]
		respEntry, ok := req.resultMap[ip]
		if ok {
			respEntry.LogInfo = logInfo
		}
	}

}

func (req *request) getPortInfo() {
	defer req.deferMethod()
	tailsAdders := strings.Split(req.tailsAddr, ",")
	for _, addr := range tailsAdders {
		parts := strings.Split(addr, ":")
		ip := parts[0]
		respEntry, ok := req.resultMap[ip]
		if !ok {
			continue
		}
		split := strings.Split(respEntry.Addr, ":")
		port := split[1]
		command := fmt.Sprintf("'lsof -i:%s'", port)
		commands := []string{
			"command_exec",
			cons.LOWER_C.GetCommandName(),
			command,
			"--f",
			addr,
		}
		portInfo := forward(commands, req.basicServer)
		lines := strings.Split(portInfo, "\n")
		// 去除可能的空行（由于字符串末尾可能有一个额外的换行符）
		var cleanedLines []string
		for _, line := range lines {
			if line != "" {
				cleanedLines = append(cleanedLines, line)
			}
		}
		respEntry.PortInfo = cleanedLines
	}
}

func (req *request) getInterfaceInfo() {
	defer req.deferMethod()
}

func forward(commands []string, basicServer *BasicServer) string {
	args, err := basicServer.context.LoadSystemParams(commands)
	if err != nil {
		return "加载参数失败"
	}
	resp, err := basicServer.context.Servlet(args, false)
	if err != nil {
		return err.Error()
	}
	return string(resp)
}

func (req *request) deferMethod() {
	req.wg.Done()
}
