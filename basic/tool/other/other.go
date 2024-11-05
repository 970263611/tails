package othertool

import (
	"fmt"
	"regexp"
	"strconv"
)

const (
	REGEX_IP            = "^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$"
	REGEX_DASE_AND_WORD = "^-\\p{L}+"
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
