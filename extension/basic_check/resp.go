package basic_check

import "sync"

type request struct {
	basicServer *BasicServer
	addr        string
	tailsAddr   string
	groupName   string
	projectName string
	logFilePath string
	resultMap   map[string]*RespEntry
	resp        Resp
	wg          sync.WaitGroup
}

type Resp struct {
	ProjectName string
	Node        []*RespEntry
}

type RespEntry struct {
	Location    string
	ProjectName string
	Addr        string
	LogInfo     string
	PortInfo    []string
}
