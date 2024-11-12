package othertool

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
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

func CheckIp(ip string) bool {
	return RegularValidate(ip, REGEX_IP)
}

func CheckPortByString(port string) bool {
	p, err := strconv.Atoi(fmt.Sprintf("%v", port))
	if err != nil {
		return false
	}
	return CheckPort(p)
}

func CheckPort(port int) bool {
	return RangeValidate(port, 0, 65535)
}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func IsFileReadable(filename string) bool {
	file, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer file.Close()
	return true
}

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

func RemoveElement(arr []string, elementToRemove string) []string {
	var newArr []string
	for _, value := range arr {
		if value != elementToRemove {
			newArr = append(newArr, value)
		}
	}
	return newArr
}
