package othertool

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	REGEX_IP             = "^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$"
	REGEX_DASE_AND_WORD  = "^-\\p{L}+"
	REGEX_2DASE_AND_WORD = "^--[a-zA-Z]+"
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
文件是否存在
*/
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

/*
*
文件是否可读
*/
func IsFileReadable(filename string) bool {
	file, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer file.Close()
	return true
}

/*
*
按照中文拆分字符串，入参：字符串本身，拆分后每个字符串的长度
*/
func SplitByChinese(s string, length int) []string {
	var result []string
	var current string
	for _, r := range s {
		if len(current)+len(string(r)) > length {
			result = append(result, current)
			current = ""
		}
		current += string(r)
	}
	if current != "" {
		result = append(result, current)
	}
	return result
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
