package monitor_server

// StatisticResp statistic接口返回体结构
type StatisticResp struct {
	Code    int
	Data    []StatisticRespEntry
	Message string
}

// StatisticRespEntry statistic接口返回对象
type StatisticRespEntry struct {
	CenterFlg     string
	TotalSucCount int
}
