package uri_switch

import (
	"fmt"
	"testing"
)

func TestDoHandler(t *testing.T) {
	params := map[string]any{
		"host":     "localhost",
		"port":     9999,
		"username": "yanhaixin",
		"password": "B10BF0502F1E0573E391A8AF50D9FB0DA22FE6AD862D3276",
		"systemId": "885e492d-bec8-4268-a2fc-ce0162c066fc",
		"code":     "GATEWAY-GSCT-TEST",
		"name":     "cfzt-edb",
	}
	instance := GetInstance(nil)
	resp := instance.Do(params)
	if resp != nil {
		fmt.Println(string(resp))
	}
}

// 测试apiName 不存在
func TestCreateApi(t *testing.T) {
	params := map[string]any{
		"host":     "localhost",
		"port":     9999,
		"username": "yanhaixin",
		"password": "B10BF0502F1E0573E391A8AF50D9FB0DA22FE6AD862D3276",
		"systemId": "885e492d-bec8-4268-a2fc-ce0162c066fc",
		"code":     "GATEWAY-GSCT-TEST",
		"name":     "cfzt-edb",
		"enabled":  "false",
		"uri":      "/cfzt-edb/query/2",
		"apiName":  "cfzt-edb-ddd",
	}
	instance := GetInstance(nil)
	resp := instance.Do(params)
	if resp != nil {
		fmt.Println(string(resp))
	}
}

// 测试apiName 存在
func TestUpdateApi(t *testing.T) {
	params := map[string]any{
		"host":     "localhost",
		"port":     9999,
		"username": "yanhaixin",
		"password": "B10BF0502F1E0573E391A8AF50D9FB0DA22FE6AD862D3276",
		"systemId": "885e492d-bec8-4268-a2fc-ce0162c066fc",
		"code":     "GATEWAY-GSCT-TEST",
		"name":     "cfzt-edb",
		"enabled":  "false",
		"uri":      "/cfzt-edb/query/2",
		"apiName":  "cfzt-edb",
	}
	instance := GetInstance(nil)
	resp := instance.Do(params)
	if resp != nil {
		fmt.Println(string(resp))
	}
}

// 测试更新,封停, uri不存在
func TestUpdateApi02(t *testing.T) {
	params := map[string]any{
		"host":     "localhost",
		"port":     9999,
		"username": "yanhaixin",
		"password": "B10BF0502F1E0573E391A8AF50D9FB0DA22FE6AD862D3276",
		"systemId": "885e492d-bec8-4268-a2fc-ce0162c066fc",
		"code":     "GATEWAY-GSCT-TEST",
		"name":     "cfzt-edb",
		"enabled":  "false",
		"uri":      "/cfzt-edb/query/2",
		"apiName":  "cfzt-edb",
	}
	instance := GetInstance(nil)
	resp := instance.Do(params)
	if resp != nil {
		fmt.Println(string(resp))
	}
}

// 测试更新,封停, uri存在
func TestUpdateApi03(t *testing.T) {
	params := map[string]any{
		"host":     "localhost",
		"port":     9999,
		"username": "yanhaixin",
		"password": "B10BF0502F1E0573E391A8AF50D9FB0DA22FE6AD862D3276",
		"systemId": "885e492d-bec8-4268-a2fc-ce0162c066fc",
		"code":     "GATEWAY-GSCT-TEST",
		"name":     "cfzt-edb",
		"enabled":  "false",
		"uri":      "/cfzt-edb/query",
		"apiName":  "cfzt-edb",
	}
	instance := GetInstance(nil)
	resp := instance.Do(params)
	if resp != nil {
		fmt.Println(string(resp))
	}
}

// 测试更新,解停, uri不存在
func TestUpdateApi04(t *testing.T) {
	params := map[string]any{
		"host":     "localhost",
		"port":     9999,
		"username": "yanhaixin",
		"password": "B10BF0502F1E0573E391A8AF50D9FB0DA22FE6AD862D3276",
		"systemId": "885e492d-bec8-4268-a2fc-ce0162c066fc",
		"code":     "GATEWAY-GSCT-TEST",
		"name":     "cfzt-edb",
		"enabled":  "true",
		"uri":      "/cfzt-edb/query/2",
		"apiName":  "cfzt-edb",
	}
	instance := GetInstance(nil)
	resp := instance.Do(params)
	if resp != nil {
		fmt.Println(string(resp))
	}
}

// 测试更新,解停, uri存在
func TestUpdateApi05(t *testing.T) {
	params := map[string]any{
		"host":     "localhost",
		"port":     9999,
		"username": "yanhaixin",
		"password": "B10BF0502F1E0573E391A8AF50D9FB0DA22FE6AD862D3276",
		"systemId": "885e492d-bec8-4268-a2fc-ce0162c066fc",
		"code":     "GATEWAY-GSCT-TEST",
		"name":     "cfzt-edb",
		"enabled":  "true",
		"uri":      "/cfzt-edb/query",
		"apiName":  "cfzt-edb",
	}
	instance := GetInstance(nil)
	resp := instance.Do(params)
	if resp != nil {
		fmt.Println(string(resp))
	}
}
