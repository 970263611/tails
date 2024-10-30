package othertool

import "regexp"

func IsValidIP(ip string) bool {
	// 正则表达式匹配 IP 地址
	ipRegex := `^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`
	re := regexp.MustCompile(ipRegex)
	return re.MatchString(ip)
}
