package commandtool

import (
	"testing"
)

func TestExecCmdByStr(t *testing.T) {
	str, err := ExecCmdByStr("cmd", "/C", "dir", "C:\\Program Files")
	if err != nil {
		t.Errorf("err:%v", err)
	} else {
		t.Logf("str:%v", str)
	}
}
