package config_update

type QueryConfigRespEntry struct {
	TotalCount     int
	PageNumber     int
	PagesAvailable int
	PageItems      []Items
}

type Items struct {
	Id               string      `json:"id"`
	DataId           string      `json:"dataId"`
	Group            string      `json:"group"`
	Content          string      `json:"content"`
	Md5              interface{} `json:"md5"`
	EncryptedDataKey string      `json:"encryptedDataKey"`
	Tenant           string      `json:"tenant"`
	AppName          string      `json:"appName"`
	Type             string      `json:"type"`
}

type UpdateConfigRespEntry struct {
	Code    int
	Message string
	Data    bool
}

type LoginRespEntry struct {
	token string
}
