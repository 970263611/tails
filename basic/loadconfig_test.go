package basic

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
	"testing"
)

/**
示例配置
app:
  server:
    host:
    -- 192.168.27.1
    -- 192.168.27.1
    username: haha
log:
  level: info
*/

func loadConfig() *Context {
	yml := "app:\n  server:\n    host:\n    -- 192.168.27.1\n    -- 192.168.27.1\n    username: haha\nlog:\n  level: info"
	reader := strings.NewReader(yml)
	v := viper.New()
	v.SetConfigType("yaml")
	for key, value := range defaultParams {
		v.SetDefault(key, value)
	}
	v.ReadConfig(reader)
	return &Context{Config: v}
}

// 第一种用法
type Config struct {
	App struct {
		Server struct {
			Host     []string `mapstructure:"host"`
			Username string   `mapstructure:"username"`
		} `mapstructure:"server"`
	} `mapstructure:"app"`
	Log struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"username"`
}

func TestLoadConfig(t *testing.T) {
	var config Config
	err := Unmarshal(&config)
	if err != nil {
		fmt.Println(err)
	}
	t.Log(config)
}

// 第二种用法
func TestLoadConfig2(t *testing.T) {
	//context := loadConfig()
	//host := context.Config.GetStringSlice("app.server.host")
	//username := context.Config.GetString("app.server.username");
	//...
}
