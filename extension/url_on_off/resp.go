package uri_switch

type SelectGetWayIDRespEntry struct {
	Id   string
	Code string
}

type SelectRuleRespEntry struct {
	Name     string
	RuleYaml string
	RuleJson string
}

type RespEntry struct {
	Code    string
	Message string
}

type ApiQueryRespEntry struct {
	Id             string
	ApiName        string
	PredicateItems string
	//0是精确 1 是前缀 2是正则
	ReleaseStatus int
	GatewayId     string
}

type LiuKongQueryRespEntry struct {
	Id string `json:"id"`
	//来源模式 0是路由 1是API
	ResourceMode int `json:"resourceMode"`
	//路由名称/API组名称
	Resource        string `json:"resource"`
	ControlBehavior int    `json:"controlBehavior"`
	ParseStrategy   int    `json:"parseStrategy"`
	FieldName       string `json:"fieldName"`
	Prop            bool   `json:"prop"`
	//0是线程数 1是QPS
	Grade int `json:"grade"`
	//线程个数
	Count       float32 `json:"count"`
	IntervalSec int     `json:"intervalSec"`
	//秒
	Region               string `json:"region"`
	Burst                int    `json:"burst"`
	MatchStrategy        int    `json:"matchStrategy"`
	RouteName            string `json:"routeName"`
	Type                 string `json:"type"`
	MatchType            string `json:"matchType"`
	Index                int    `json:"index"`
	MaxQueueingTimeoutMs int    `json:"maxQueueingTimeoutMs"`
	GatewayId            string `json:"gatewayId"`
}
