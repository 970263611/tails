package basic_check

import "sync"

type request struct {
	basicServer *BasicServer
	addr        string
	tailsAddr   string
	logPath     string
	projectName string
	resultMap   map[string]*RespEntry
	resp        Resp
	wg          sync.WaitGroup
}

type Resp struct {
	ProjectName string       `json:"projectName"`
	Node        []*RespEntry `json:"node"`
}

type RespEntry struct {
	Addr          string           `json:"addr"`
	LogInfo       string           `json:"logInfo"`
	PortInfo      []string         `json:"portInfo"`
	InterFaceInfo []*InterfaceResp `json:"interFaceInfo"`
}

type InterfaceRequest struct {
	Header Header         `json:"header"`
	Body   map[string]any `json:"body"`
}

type Header struct {
	GlobalBusiTrackNo string `json:"globalBusiTrackNo"`
	ReqSysSriNo       string `json:"reqSysSriNo"`
	SendTime          string `json:"sendTime"`
}

type InterfaceResp struct {
	Url  string `json:"url"`
	Code string `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

type EdbListResp struct {
	Code string           `json:"code"`
	Data EdbListRespEntry `json:"data"`
	Msg  string           `json:"msg"`
}

type EdbListRespEntry struct {
	PageNum    int `json:"pageNum"`
	PageSize   int `json:"pageSize"`
	TotalCount int `json:"totalCount"`
	TotalPage  int `json:"totalPage"`
	List       []struct {
		MainTaskId      int    `json:"mainTaskId"`
		DataName        string `json:"dataName"`
		DataDesc        string `json:"dataDesc"`
		Category        bool   `json:"category"`
		SubTaskCount    int    `json:"subTaskCount"`
		SendSysCode     string `json:"sendSysCode"`
		RecvSysCode     string `json:"recvSysCode"`
		BusinessSysCode string `json:"businessSysCode"`
		MainTaskName    string `json:"mainTaskName"`
		CreateTime      string `json:"createTime"`
		TaskTypeName    string `json:"taskTypeName"`
		Source          string `json:"source"`
		TaskType        int    `json:"taskType"`
		FileDate        string `json:"fileDate"`
	} `json:"list"`
}
