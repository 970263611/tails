package suspend_server

type SelectGetWayIDRespEntry struct {
	Id   string
	Code string
}

type SelectRuleRespEntry struct {
	Name     string
	RuleYaml string
	RuleJson string
}

type ApiRespEntry struct {
	Code    string
	Message string
}
