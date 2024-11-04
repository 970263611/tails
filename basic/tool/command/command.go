package commandtool

import (
	"basic/tool/charset"
	"bytes"
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

func ExecCmdByStr(commandArr ...string) (string, error) {
	command := exec.Command(commandArr[0], commandArr[1:]...)
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

func ExecCmdByPipe(command ...[]string) (string, error) {
	var input, output *bytes.Buffer
	first := true
	for _, cmdStr := range command {
		cmd := exec.Command(cmdStr[0], cmdStr[1:]...)
		if !first {
			input = output
			cmd.Stdin = input
		}
		output = &bytes.Buffer{}
		cmd.Stdout = output
		if err := cmd.Start(); err != nil {
			return "", err
		}
		if err := cmd.Wait(); err != nil {
			return "", err
		}
		first = false
	}
	return output.String(), nil
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
