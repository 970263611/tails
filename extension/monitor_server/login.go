package monitor_server

// LoginReq login接口发送数据的结构体
type LoginReq struct {
	userId   string
	password string
}

// LoginResp login接口返回体结构
type LoginResp struct {
	Code    int
	Data    []LoginRespEntry
	Message string
}

// LoginRespEntry login接口返回对象
type LoginRespEntry struct {
	Token string
}
