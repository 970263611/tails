package config_update

import (
	cons "basic/constants"
	iface "basic/interfaces"
	"basic/tool/net"
	"basic/tool/utils"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type transfer struct {
	host        string
	port        int
	username    string
	password    string
	namespaceId string
	group       string
	dataId      string
	outPut      string
	upload      string
	Items       Items
}

const (
	GetConfigUrl    string = "/nacos/v1/cs/configs"
	UpdateConfigUrl string = "/nacos/v2/cs/config"
	SUCCESS         int    = 0
	JSON            string = "json"
	TEXT            string = "text"
	XML             string = "xml"
	YAML            string = "yaml"
	HTML            string = "html"
	Properties      string = "properties"
)

type ConfigServer struct{}

func GetInstance(globalContext iface.Context) iface.Component {
	return &ConfigServer{}
}

func (c *ConfigServer) GetName() string {
	return "config_update"
}

func (c *ConfigServer) GetDescribe() string {
	return "Nacos更新配置文件 " +
		"\n例: 获取配置 config_update -h 127.0.0.1 -p 8848 -u nacos -w nacos -n test1112 -g group -d sit.yml -o /appcpyy/ " +
		"\n例: 上传配置 config_update -h 127.0.0.1 -p 8848 -u nacos -w nacos -n test1112 -g group -d sit.yml -U /appcpyy/sit.yml"
}

func (c *ConfigServer) Register(cm iface.ComponentMeta) {
	cm.AddParameters(cons.STRING, cons.LOWER_H, "ip", "host", true, func(s string) error {
		if !utils.CheckIp(s) {
			return errors.New("nacos服务ip不合法")
		}
		return nil
	}, "nacos服务的主机地址")
	cm.AddParameters(cons.INT, cons.LOWER_P, "port", "port", true, func(s string) error {
		if !utils.CheckPortByString(s) {
			return errors.New("nacos服务port不合法")
		}
		return nil
	}, "nacos服务的端口")
	cm.AddParameters(cons.STRING, cons.LOWER_U, "username", "username", true, nil, "nacos登录用户名")
	cm.AddParameters(cons.STRING, cons.LOWER_W, "password", "password", true, nil, "nacos登录密码")
	cm.AddParameters(cons.STRING, cons.LOWER_N, "namespaceId", "namespaceId", true, nil, "命名空间")
	cm.AddParameters(cons.STRING, cons.LOWER_G, "group", "group", true, nil, "组名")
	cm.AddParameters(cons.STRING, cons.LOWER_D, "dataId", "dataId", true, nil, "配置名称")
	cm.AddParameters(cons.STRING, cons.LOWER_O, "outPut", "outPut", false, nil, "将配置文件导出到指定目录,默认当前目录")
	cm.AddParameters(cons.STRING, cons.UPPER_U, "upload", "upload", false, nil, "文件路径,将修改后的配置文件进行上传更新")
}

func (c *ConfigServer) Do(params map[string]any) (resp []byte) {
	res := &transfer{
		host:        params["host"].(string),
		port:        params["port"].(int),
		username:    params["username"].(string),
		password:    params["password"].(string),
		namespaceId: params["namespaceId"].(string),
		group:       params["group"].(string),
		dataId:      params["dataId"].(string),
	}
	if params["outPut"] != nil {
		res.outPut = params["outPut"].(string)
	}
	// 获取配置,根据是否上传 判断是否是获取配置还是更新配置
	err := res.getConfig()
	if err != nil {
		return []byte(err.Error())
	}
	if params["upload"] != nil {
		res.upload = params["upload"].(string)
		err := res.updateConfig()
		if err != nil {
			return []byte(err.Error())
		}
		return []byte("更新配置文件成功!")
	} else {
		err = res.downConfigFile()
		if err != nil {
			return []byte(err.Error())
		}
		return []byte("配置文件写入本地成功")
	}
}

func (t *transfer) getConfig() error {
	urlPrefix := fmt.Sprintf("%s://%s:%d", "http", t.host, t.port)
	queryParams := url.Values{}
	queryParams.Add("username", t.username)
	queryParams.Add("password", t.password)
	queryParams.Add("namespaceId", t.namespaceId)
	queryParams.Add("tenant", t.namespaceId)
	queryParams.Add("group", t.group)
	queryParams.Add("dataId", t.dataId)

	queryParams.Add("search", "blur")
	queryParams.Add("pageNo", "1")
	queryParams.Add("pageSize", "10")

	var configResp QueryConfigRespEntry
	err := net.GetRespStruct(urlPrefix+GetConfigUrl, queryParams, nil, &configResp)
	if err != nil {
		return fmt.Errorf("Nacos获取配置接口失败%v", err)
	}
	if configResp.TotalCount != 1 {
		return fmt.Errorf("获取Nacos配置个数为: %v,请确认唯一", configResp.TotalCount)
	}
	items := configResp.PageItems[0]
	t.Items = items
	log.Info(fmt.Sprintf("获取配置成功,配置为:%s", items.Content))
	return nil
}

func (t *transfer) updateConfig() error {
	data, err := t.checkConfigFile()
	if err != nil {
		return err
	}

	// 更新配置
	urlPrefix := fmt.Sprintf("%s://%s:%d", "http", t.host, t.port)

	queryParams := url.Values{}
	queryParams.Add("username", t.username)
	queryParams.Add("password", t.password)
	queryParams.Add("namespaceId", t.namespaceId)
	queryParams.Add("group", t.group)
	queryParams.Add("dataId", t.dataId)
	queryParams.Add("content", data)
	queryParams.Add("type", t.Items.Type)
	var fullURL = urlPrefix + UpdateConfigUrl
	if queryParams != nil {
		u, err := url.Parse(fullURL)
		if err != nil {
			return fmt.Errorf("解析 URL 失败: %v", err)
		}
		u.RawQuery = queryParams.Encode()
		fullURL = u.String()
	}
	postData := map[string]string{
		"namespaceId": t.namespaceId,
		"dataId":      t.dataId,
	}
	header := http.Header{}
	header.Add("Content-Type", "application/x-www-form-urlencoded")
	var updateConfigRespEntry UpdateConfigRespEntry
	err = net.PostRespStruct(fullURL, postData, header, &updateConfigRespEntry)
	if err != nil {
		return fmt.Errorf("Nacos发布配置接口失败: " + err.Error())
	}
	if updateConfigRespEntry.Code != SUCCESS {
		return fmt.Errorf("Nacos发布配置接口失败: " + updateConfigRespEntry.Message)
	}
	if !updateConfigRespEntry.Data {
		return fmt.Errorf("Nacos发布配置接口失败: " + updateConfigRespEntry.Message)
	}
	log.Info(fmt.Sprintf("更新配置文件成功:%s", t.dataId))
	return nil
}

func (t *transfer) downConfigFile() error {
	var dir string
	var err error
	if t.outPut != "" {
		dir = t.outPut
		if err = os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	} else {
		dir, err = os.Getwd()
		if err != nil {
			return err
		}
	}
	filePath := filepath.Join(dir, t.dataId)
	err = os.WriteFile(filePath, []byte(t.Items.Content), 0644)
	if err != nil {
		return err
	}
	log.Info(fmt.Sprintf("配置文件写入本地成功,位置为:%s", filePath))
	return nil
}

func (t *transfer) checkConfigFile() (string, error) {
	filePath := t.upload
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	switch t.Items.Type {
	case JSON:
		if !isValidJSON(data) {
			return "", errors.New("不是合法的 JSON 文件。")
		}
	case XML:
		if !isValidXML(data) {
			return "", errors.New("不是合法的 XML 文件。")
		}
	case YAML:
		if !isValidYAML(data) {
			return "", errors.New("不是合法的 YAML 文件。")
		}
	case TEXT:
		if !isValidTEXT(data) {
			return "", errors.New("不是合法的 TEXT 文件。")
		}
	case HTML:
		if !isValidHTML(data) {
			return "", errors.New("不是合法的 HTML 文件。")
		}
	case Properties:
		if !isValidProperties(data) {
			return "", errors.New("不是合法的 Properties  文件。")
		}
	}
	return string(data), nil
}
