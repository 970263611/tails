package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

const (
	REGEX_IP             = "^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$"
	REGEX_DASE_AND_WORD  = "^-\\p{L}+"
	REGEX_2DASE_AND_WORD = "^--[a-zA-Z]+"
	REGEX_DOMAIN_NAME    = "^([a-zA-Z0-9]([a-zA-Z0-9\\-]{0,61}[a-zA-Z0-9])?\\.)+[a-zA-Z]{2,}$"
)

/*
*
正则校验
*/
func RegularValidate(val, regex string) bool {
	reg := regexp.MustCompile(regex)
	return reg.MatchString(val)
}

/*
*
数字范围校验
*/
func RangeValidate(val, min, max int) bool {
	if val < min || val > max {
		return false
	}
	return true
}

/*
*
验证IP是否规范
*/
func CheckIp(ip string) bool {
	return RegularValidate(ip, REGEX_IP)
}

/*
*
验证端口号是否规范
*/
func CheckPortByString(port string) bool {
	p, err := strconv.Atoi(fmt.Sprintf("%v", port))
	if err != nil {
		return false
	}
	return CheckPort(p)
}

/*
*
验证端口号是否规范
*/
func CheckPort(port int) bool {
	return RangeValidate(port, 0, 65535)
}

/*
*
校验地址是否合法，必须是合法域名或者ip:port,可以有多组ip:port,必须是,分割。
*/
func CheckAddr(multipleAddr string) bool {
	addrs := strings.Split(multipleAddr, ",")
	for _, addr := range addrs {
		if RegularValidate(addr, REGEX_DOMAIN_NAME) {
			return true
		}
		arr := strings.Split(addr, ":")
		if len(arr) != 2 {
			return false
		}
		if !CheckIp(arr[0]) || !CheckPortByString(arr[1]) {
			return false
		}
	}
	return true
}

/*
*
字符串按照空格拆分为字符串数组，但是被单引号或者双引号包裹的部分不做拆分
*/
func SplitString(s string) []string {
	re := regexp.MustCompile(`'(?:[^']|'')*'|"(?:[^"]|\"\")*"|\S+`)
	return re.FindAllString(s, -1)
}

/*
*
检查字符串是否是boolean类型
*/
func CheckIsBooleanByString(s string) bool {
	// 忽略大小写
	s = strings.ToLower(s)
	return s == "true" || s == "false"
}

/*
*
获取连续数量的某个字符串
*/
func GetBlankByNum(num int, s string) string {
	if num < 1 {
		return ""
	} else {
		var sb strings.Builder
		for i := 0; i < num; i++ {
			sb.WriteString(s)
		}
		return sb.String()
	}
}

/*
*
生成uuid
*/
func GenerateUUID() (string, error) {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		return "", err
	}

	// 设置UUID版本为4（随机生成）和变体信息
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	return hex.EncodeToString(uuid), nil
}

/*
*
获取当前线程唯一标识
*/
func GetGoroutineID() string {
	var buf [1 << 20]byte
	n := runtime.Stack(buf[:], false)
	stackTrace := string(buf[:n])
	lines := strings.Split(stackTrace, "\n")
	for _, line := range lines {
		if strings.Contains(line, "goroutine") {
			parts := strings.Fields(line)
			for i, part := range parts {
				if part == "goroutine" && i+1 < len(parts) {
					return parts[i+1]
				}
			}
		}
	}
	return ""
}

/*
*
ENC解密
*/
func JasyptDec(input, password string) (string, error) {
	dMsg, err := Decrypt(input, password)
	if err != nil {
		msg := fmt.Sprintf("jasypt解密失败: %v", err)
		log.Error(msg)
		return input, errors.New(msg)
	}
	return dMsg, nil
}

/*
*
成对去掉字符串两侧单引号或双引号
*/
func RemoveQuotes(s string) string {
	if len(s) >= 2 && ((s[0] == '\'' && s[len(s)-1] == '\'') || (s[0] == '"' && s[len(s)-1] == '"')) {
		return s[1 : len(s)-1]
	}
	return s
}

/*
*
将\或\\替换为/
*/
func ReplaceSlash(str string) string {
	if runtime.GOOS == "windows" {
		// 先替换双反斜杠
		str = strings.Replace(str, "\\\\", "/", -1)
		// 再替换单反斜杠
		str = strings.Replace(str, "\\", "/", -1)
	}
	return str
}

/*
*
获取tails服务所在路径,例如:/home/dimple/
*/
func GetRootPath() string {
	rootpath, _ := os.Executable()
	return ReplaceSlash(PathDir(rootpath)) + "/"
}

/*
*
获取path中的路径部分
*/
func PathDir(filepath string) string {
	return path.Dir(ReplaceSlash(filepath))
}

/*
*
判断路径是否是绝对路径
*/
func IsAbsolutePath(pathStr string) bool {
	if len(pathStr) == 0 {
		return false
	}
	if strings.HasPrefix(pathStr, "/") {
		return true
	}
	if len(pathStr) >= 2 && strings.Contains("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz", pathStr[0:1]) && strings.Contains(":/", pathStr[1:2]) {
		return true
	}
	return false
}

/*
*
获取给定路径的绝对路径
*/
func GetAbsolutePath(path string) string {
	path = ReplaceSlash(path)
	if IsAbsolutePath(path) {
		return path
	} else {
		return GetRootPath() + path
	}
}
