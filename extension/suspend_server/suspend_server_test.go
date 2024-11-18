package suspend_server

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
		"enabled":  "false",
		"uri":      "/cfzt-edb/query",
		"apiName":  "cfzt-edb",
	}
	instance := GetInstance()
	resp := instance.Do(params)
	if resp != nil {
		fmt.Println(string(resp))
	}
}
