package basic

import (
	log "github.com/sirupsen/logrus"
)

func Servlet(key string, commands []string) []byte {
	log.Info("call component '" + key + "'")
	c := components[Key(key)]
	if c == nil {
		result := "component '" + key + "' not found"
		log.Error(result)
		return []byte(result)
	}
	return c.Do(commands)
}
