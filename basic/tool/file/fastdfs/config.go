package fdfs_client

import (
	"bufio"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
)

type config struct {
	trackerAddr []string
	maxConns    int
}

func newConfig(configName string) (*config, error) {
	config := &config{}
	f, err := os.Open(configName)
	if err != nil {
		return nil, err
	}
	splitFlag := "\n"
	if runtime.GOOS == "windows" {
		splitFlag = "\r\n"
	}
	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadString('\n')
		line = strings.TrimSuffix(line, splitFlag)
		str := strings.SplitN(line, "=", 2)
		switch str[0] {
		case "tracker_server":
			config.trackerAddr = append(config.trackerAddr, strings.ReplaceAll(str[1], "\n", ""))
		case "maxConns":
			config.maxConns, err = strconv.Atoi(strings.ReplaceAll(str[1], "\n", ""))
			if err != nil {
				return nil, err
			}
		}
		if err != nil {
			if err == io.EOF {
				return config, nil
			}
			return nil, err
		}
	}
	return config, nil
}
