package monitor_check

// LoginResp login接口返回体结构
type LoginResp struct {
	Code    int
	Data    LoginRespEntry
	Message string
}

// LoginRespEntry login接口返回对象
type LoginRespEntry struct {
	Token string
}

type A1Resp struct {
	Code    int
	Data    []A1RespEntry
	Message string
}

type A1RespEntry struct {
	//总成功率
	TotalSucRate int
	//系统吞吐量
	TradeTps string
}

type A5Resp struct {
	Code    int
	Data    []A5RespEntry
	Message string
}

type A5RespEntry struct {
	//中心标志
	CenterFlag string
	//服务名称
	SubsystemCode string
	//服务中文名称
	SubsystemName string
	//最小节点数
	TotalNum int
	//当前运行节点数
	CurrentNum int
	//当前异常IP
	ErrorIp string
	//当前异常节点数
	ErrorNum int
	//当前服务数低于最小服务数
	AlarmStatus bool
}

type A6Resp struct {
	Code    int
	Data    []A6RespEntry
	Message string
}

type A6RespEntry struct {
	//中心标志
	CenterFlag string
	//主机端口号
	HostPort int
	//启用标志0-停用 1-启用
	ActiveFlag string
	//端口状态
	PortState bool
}

type A10Resp struct {
	Code    int
	Data    []A10RespEntry
	Message string
}

type A10RespEntry struct {
	//Job名字
	JobName string
	//中心标识集合
	CenterList []int
	//告警级别
	Level int
	//是否暂停告警
	Pause bool
	//开始时间
	BeginTime string
	//持续时间（分钟数）
	AlarmSustainTime string
}
