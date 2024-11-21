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
	ReleaseStatus  int
	GatewayId      string
}

type LiuKongQueryRespEntry struct {
	Id string
	//来源模式 0是路由 1是API
	ResourceMode int
	//路由名称/API组名称
	Resource        string
	ControlBehavior int
	ParseStrategy   int
	FieldName       string
	Prop            bool
	Grade           int
	//线程个数
	Count       float32
	IntervalSec int
	//秒
	Region               string
	Burst                int
	MatchStrategy        int
	RouteName            string
	Type                 string
	MatchType            string
	Index                int
	MaxQueueingTimeoutMs int
	GatewayId            string
}
