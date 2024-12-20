package commandtool

import (
	"basic/tool/charset"
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

var charsetType charsettool.Charset

/*
*
系统字符集初始化
*/
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

/*
*
执行系统指令
*/
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

/*
*
部分命令在执行过滤后如果没有匹配到内容，会返回一个非0的状态码，
该状态码不应被认为是命令执行错误，例如grep，如果没有匹配到内容会返回1
此时1代表未找到匹配内容，而并非命令执行失败
因而需特殊处理此类命令
*/
var specialCommandsMap = map[string]string{
	"grep": "1",
}

const SPECIAL_COMMANDS_RESULT string = "未找到匹配内容"

/*
*
管道流模式执行系统指令
*/
func ExecCmdByPipe(command ...[]string) (string, error) {
	var input, output, errorOutPut *bytes.Buffer
	first := true
	for _, cmdStr := range command {
		cmd := exec.Command(cmdStr[0], cmdStr[1:]...)
		if !first {
			input = output
			cmd.Stdin = input
		}
		// 创建一个缓冲区来捕获原生错误
		errorOutPut = &bytes.Buffer{}
		cmd.Stderr = errorOutPut

		output = &bytes.Buffer{}
		cmd.Stdout = output
		if err := cmd.Start(); err != nil {
			return "", fmt.Errorf(" 错误信息为: %s", errorOutPut.String())
		}
		if err := cmd.Wait(); err != nil {
			if _, ok1 := specialCommandsMap[cmdStr[0]]; ok1 {
				if exitError, ok2 := err.(*exec.ExitError); ok2 {
					if exitError.ExitCode() == 1 {
						return SPECIAL_COMMANDS_RESULT, nil
					}
				}
			}
			return "", fmt.Errorf(" 错误信息为: %s", errorOutPut.String())
		}
		first = false
	}
	return output.String(), nil
}

/*
*
获取window下字符集
*/
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
