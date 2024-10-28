package commandtool

import (
	"basic/tool/charset"
	"os/exec"
	"runtime"
	"strings"
)

var charsetType charsettool.Charset

func init() {
	switch runtime.GOOS {
	case "windows":
		charset, _ := getWindowsCharset()
		if charset == "GBK" {
			charsetType = charsettool.GBK
		} else {
			charsetType = charsettool.UTF8
		}
	default:
		charsetType = charsettool.UTF8
	}
}

func ExecCmdByStr(name string, arg ...string) (string, error) {
	command := exec.Command(name, arg...)
	data, err := command.CombinedOutput()
	if err != nil {
		if "windows" == runtime.GOOS {
			return "", err
		}
	}
	if "windows" == runtime.GOOS {
		data, _ = charsettool.Convert(data, charsetType, charsettool.UTF8)
	}
	return string(data), err
}

func getWindowsCharset() (string, error) {
	cmd := exec.Command("cmd", "/C", "chcp")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	codePageStr := string(output)
	if strings.Contains(codePageStr, "936") {
		return "GBK", nil
	} else {
		return "UTF-8", nil
	}
}
