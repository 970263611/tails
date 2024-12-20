package basic_check

type request struct {
	c           chan string
	addr        string
	tailsAddr   string
	groupName   string
	projectName string
	logFilePath string
	resultMap   map[string]*RespEntry
	resp        Resp
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
	PortInfo    string
}
