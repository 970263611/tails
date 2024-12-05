package cons

const (
	COMPONENT_KEY string = "componentKey"
	CONFIG_SALT   string = "jasypt.salt"
	WEB_SVC       string = "web_svc"
	SYSTEM_PARAMS string = "systemparams"
	DIFFUSE       string = "diffuse"
	ENC_JAR              = "jasypt-1.9.3.jar"
)

const (
	RESULT_CODE    = "code"
	RESULT_DATA    = "data"
	RESULT_SUCCESS = "00"
	RESULT_ERROR   = "99"
)

const (
	SYSPARAM_HELP        string = "--help,--h"
	SYSPARAM_CONFIG      string = "--config,--c"
	SYSPARAM_KEY         string = "--key"
	SYSPARAM_SALT        string = "--salt"
	SYSPARAM_FORWORD     string = "--f,--forward"
	SYSPARAM_GID         string = "--gid"
	SYSPARAM_LOG_LEVEL   string = "--loglevel"
	SYSPARAM_LOG_OUTTYPE string = "--logouttype"
)

// ParamType ================start=============
type ParamTypeIface interface {
	GetParamType() int
}

type ParamType int

func (c ParamType) GetParamType() int {
	return int(c)
}

const (
	NO_VALUE ParamType = iota
	INT
	STRING
	KEY_NAME
)

// ParamType ================end=============

// ComponentType ================start=============
type ComponentTypeIface interface {
	GetComponentType() string
}

type ComponentType int

func (c ComponentType) GetComponentType() int {
	return int(c)
}

const (
	EXECUTE ComponentType = iota
	SYSTEM
)

// ComponentType ================end=============

// ComponentName ================start=============
type CommandNameIface interface {
	GetCommandName() string
}

type CommandName string

func (c CommandName) GetCommandName() string {
	return string(c)
}

const (
	LOWER_A CommandName = "-a" //单/多地址 单ip:port 多ip:port,ip1:port1 使用,分割
	LOWER_C CommandName = "-c" //字符串输入，命令行输入
	LOWER_D CommandName = "-d" //数据库名
	UPPER_D CommandName = "-D" //查看当前数据库key的数量
	LOWER_E CommandName = "-e" //enable
	LOWER_F CommandName = "-f" //强制执行
	LOWER_G CommandName = "-g" //组名，组别
	LOWER_H CommandName = "-h" //ip地址
	UPPER_H CommandName = "-H" //ip2地址
	LOWER_I CommandName = "-i" //输入路径
	LOWER_K CommandName = "-k" //key
	UPPER_L CommandName = "-L" //redis list
	LOWER_N CommandName = "-n" //名称，发送方名称
	UPPER_N CommandName = "-N" //名称，接收方名称
	LOWER_M CommandName = "-m" //redis模式
	UPPER_M CommandName = "-M" //redis哨兵模式主节点名称
	LOWER_O CommandName = "-o" //输出路径，输出
	LOWER_P CommandName = "-p" //第一端口
	UPPER_P CommandName = "-P" //第二端口
	LOWER_R CommandName = "-r" //重试，重复执行，启动高可用模块
	LOWER_S CommandName = "-s" //状态展示status
	UPPER_S CommandName = "-S" //redis set
	LOWER_T CommandName = "-t" //节点信息交互
	LOWER_U CommandName = "-u" //用户名
	UPPER_U CommandName = "-U" //uri
	LOWER_V CommandName = "-v" //查看节点信息
	UPPER_V CommandName = "-V" //占用其它类型，无意义
	LOWER_W CommandName = "-w" //第一密码
	UPPER_W CommandName = "-W" //第二密码
	LOWER_X CommandName = "-x" //占用其它类型，无意义
	LOWER_Y CommandName = "-y" //占用其它类型，无意义
	LOWER_Z CommandName = "-z" //占用其它类型，无意义
	UPPER_X CommandName = "-X" //占用其它类型，无意义
	UPPER_Y CommandName = "-Y" //占用其它类型，无意义
	UPPER_Z CommandName = "-Z" //占用其它类型，无意义
)

const (
	SINGLE   string = "single"
	SENTINEL string = "sentinel"
	CLUSTER  string = "cluster"
)

// ComponentName ================end=============
