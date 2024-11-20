package cons

type ComponentType int

const (
	EXECUTE ComponentType = iota
	SYSTEM
)

type ParamType int8

const (
	NO_VALUE ParamType = iota
	INT
	STRING
)

const (
	COMPONENT_KEY string = "componentKey"
	WEB_KEY       string = "web_server"
	NEEDHELP      string = "--help"
	GID           string = "globalID"
	SALT          string = "--salt"
	CONFIG_SALT   string = "jasypt.salt"
)
