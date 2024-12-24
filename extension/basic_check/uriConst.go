package basic_check

import (
	cons "basic/constants"
	"basic/tool/net"
	"basic/tool/utils"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	reqSysSriNo string = "11697105108"

	SUCCESS string = "应用启动成功"

	access string = "cfzt-access"

	accessUrl string = "/access/login"

	edb string = "cfzt-edb"

	edbListUrl string = "/api/mainTask/list"

	edbReconsumeUrl string = "/api/mainTask/reconsumeByName"

	edbSubtaskUrl string = "/subtask/list"

	edbUplod string = "cfzt-edb-upload"

	edbUploadUrl string = "/edb/upload/queryTaskList"

	power string = "cfzt-power"

	powerUserUrl string = "/api/user/info"

	powerTokenUrl string = "/api/uaas/token/check"

	powerCodeUrl string = "/api/uaas/login/code"

	workFlow string = "cfzt-workflow"

	workFlowUrl string = "/api/repository/list"

	batch string = "cfzt-batch"

	batchUrl string = "/queryByName"
)

func (req *request) check() error {
	projectName := req.projectName
	if projectName != access && projectName != edb && projectName != edbUplod && projectName != power && projectName != workFlow && projectName != batch {
		return errors.New("请输出正确的projectName 应该为:cfzt-access,cfzt-edb,cfzt-edb-upload,cfzt-power,cfzt-worlflow,cfzt-batch 其中之一")
	}
	return nil
}

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
			Addr:          adders[i],
			InterFaceInfo: []*InterfaceResp{},
		}
		resp.Node = append(resp.Node, m[ip])
	}
	req.resultMap = m
	req.resp = resp
}

func (req *request) getLogInfo() {
	defer req.deferMethod()
	//转发
	tailsAdders := strings.Split(req.tailsAddr, ",")
	command := fmt.Sprintf("'tail -n 5000 %s | grep %s'", req.logPath, SUCCESS)
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
	switch req.projectName {
	case access:
		req.getAccessInfo()
	case edb:
		req.getEdbInfo()
	case edbUplod:
		req.getEdbUploadInfo()
	case power:
		req.getPowerInfo()
	case workFlow:
		req.getWorkflowInfo()
	case batch:
		req.getBatchInfo()
	default:

	}
}

func (req *request) getAccessInfo() {
	var interfaceRequest InterfaceRequest
	m := map[string]any{
		"username": "admin",
		"password": "6ccfe46554c1774d360da4ea4b787e06f243edc7922362392b618f3fa87a844e30929927a7038bdf212e0a23db2cdf538922f2d852418fc64222c8d28416c6ef96e9848f6aff59049364666f64e4c26efb722d036782e88d205828c8093c3d8dfa1895af96b553721908ba6a81901f4e",
	}
	interfaceRequest.Body = m
	req.commonRequest(accessUrl, &interfaceRequest)
}

func (req *request) getEdbInfo() {
	var interfaceRequest InterfaceRequest
	var m map[string]any
	var dataName = "JGQL"

	m = map[string]any{
		"pageNum":     "1",
		"pageSize":    "1",
		"sendSysCode": "",
	}
	interfaceRequest.Body = m
	req.commonRequest(edbListUrl, &interfaceRequest)

	m = map[string]any{
		"dataName": dataName,
	}
	interfaceRequest.Body = m
	req.commonRequest(edbReconsumeUrl, &interfaceRequest)

	//重试完，查询详情接口
	m = map[string]any{
		"pageNum":     "1",
		"pageSize":    "10",
		"sendSysCode": "",
	}
	interfaceRequest.Body = m
	uuid, _ := utils.GenerateUUID()
	h := Header{
		GlobalBusiTrackNo: uuid,
		SendTime:          strconv.FormatInt(time.Now().Unix(), 10),
		ReqSysSriNo:       reqSysSriNo,
	}
	interfaceRequest.Header = h
	header := http.Header{}
	header.Set("Content-Type", "application/json")
	//获取ip,循环发送
	address := strings.Split(req.addr, ",")
	for _, addr := range address {
		parts := strings.Split(addr, ":")
		ip := parts[0]
		_, ok := req.resultMap[ip]
		if !ok {
			continue
		}
		urlPrefix := fmt.Sprintf("%s://%s", "http", addr)
		var resp EdbListResp
		err := net.PostRespStruct(urlPrefix+edbListUrl, interfaceRequest, header, &resp)
		if err != nil {
			log.Error("发错接口失败：", err)
			continue
		}
		listResp := resp.Data
		var mainTaskId int
		for _, entry := range listResp.List {
			if entry.DataName == dataName {
				mainTaskId = entry.MainTaskId
				break
			}
		}
		m = map[string]any{
			"mainTaskId": mainTaskId,
		}
		interfaceRequest.Body = m
		req.commonRequest(edbSubtaskUrl, &interfaceRequest)
	}

}

func (req *request) getEdbUploadInfo() {
	var interfaceRequest InterfaceRequest
	m := map[string]any{
		"dataName":        "",
		"sendSysCode":     "",
		"recvSysCode":     "",
		"startTime":       "",
		"endTime":         "",
		"businessSysCode": "",
		"edbStatus":       "",
		"pageNum":         "1",
		"pageSize":        "1",
	}
	interfaceRequest.Body = m
	req.commonRequest(edbUploadUrl, &interfaceRequest)
}

func (req *request) getPowerInfo() {
	var interfaceRequest InterfaceRequest
	m := map[string]any{
		"tellerNo":    "20200715750",
		"accessToken": "",
		"clientId":    "972438975409164290",
	}
	interfaceRequest.Body = m
	req.commonRequest(powerUserUrl, &interfaceRequest)

	m = map[string]any{
		"userToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyTmFtZSI6IuiCluafkOWFsCIsImV4cCI6MTczNDY1ODQzMCwidXNlcklkIjoiMjAxMjAzMzgyOTAiLCJpYXQiOjE3MzQ2NTY1MTAsImp0aSI6IjMzOTJhOTA2LWU3ZjYtNGUwMS04NGI2LTVhYzdhZmYyYzIxNSJ9.PYRVtL_FW37En20oRW3_afjBmPR2lE746v4rUkIuDtg",
	}
	interfaceRequest.Body = m
	req.commonRequest(powerTokenUrl, &interfaceRequest)

	m = map[string]any{
		"userLoginName": "20200813510",
		"userLoginPwd":  "92284ac2475d0d677ac41626c8e41d421c881bcc11cd4591542bd9d14d5d6890f364f7f4031a49c9d5b945bd2ad9dd6387f16b8e1d3af4e2bbaf82675ca1cbba1e3dee1b1e23210b386e5baa0f32a885b915da598495d5cd86e4f0d8ab1087dbf19ef0d438d8080fa53fd966aa1e9068",
	}
	interfaceRequest.Body = m
	req.commonRequest(powerCodeUrl, &interfaceRequest)
}

func (req *request) getWorkflowInfo() {
	var interfaceRequest InterfaceRequest
	m := map[string]any{
		"systemId": "99711240106",
		"defName":  "",
		"status":   "false",
		"pageNum":  "1",
		"pageSize": "1",
	}
	interfaceRequest.Body = m
	req.commonRequest(workFlowUrl, &interfaceRequest)
}

func (req *request) getBatchInfo() {
	var interfaceRequest InterfaceRequest
	m := map[string]any{
		"name": "queryLegalFeedbackJob",
	}
	interfaceRequest.Body = m
	req.commonRequest(batchUrl, &interfaceRequest)
}

func (req *request) commonRequest(urlSuffix string, interFaceRequest *InterfaceRequest) {
	//请求参数
	uuid, _ := utils.GenerateUUID()
	h := Header{
		GlobalBusiTrackNo: uuid,
		SendTime:          strconv.FormatInt(time.Now().Unix(), 10),
		ReqSysSriNo:       reqSysSriNo,
	}
	interFaceRequest.Header = h
	header := http.Header{}
	header.Set("Content-Type", "application/json")
	//获取ip,循环发送
	address := strings.Split(req.addr, ",")
	for _, addr := range address {
		parts := strings.Split(addr, ":")
		ip := parts[0]
		respEntry, ok := req.resultMap[ip]
		if !ok {
			continue
		}
		urlPrefix := fmt.Sprintf("%s://%s", "http", addr)
		var resp InterfaceResp
		err := net.PostRespStruct(urlPrefix+urlSuffix, interFaceRequest, header, &resp)
		if err != nil {
			log.Error("发错查询任务接口失败：", err)
			resp.Msg = err.Error()
		}
		resp.Url = urlPrefix + urlSuffix
		respEntry.InterFaceInfo = append(respEntry.InterFaceInfo, &resp)
	}
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
