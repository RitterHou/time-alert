package main

import (
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
)

var (
	root string
	conf = make(map[string]string)
)

const (
	fileName = "TimeAlert.ini"
	content  = `; Generated by https://github.com/RitterHou/time-alert
; 10/20/30/60/..., alert will be touched off if current_minute % this_number == 0
alert_time_point=30
; disable alert on this hours
; disabled_hours=22,23`
)

func init() {
	root, _ = os.Getwd()
	confFile := path.Join(root, fileName)
	if _, err := os.Stat(confFile); os.IsNotExist(err) {
		writeFile(confFile, content)
	}
	confContent := readFile(confFile)

	var re = regexp.MustCompile("(;[\\d\\D]*?\n)") // 移除注释内容
	confContent = re.ReplaceAllString(confContent, "")

	for _, line := range strings.Split(confContent, "\n") {
		if strings.Contains(line, "=") {
			value := strings.Split(line, "=")
			conf[value[0]] = value[1]
		}
	}
}

func getConf() map[string]string {
	return conf
}

func readFile(filePath string) string {
	data, _ := ioutil.ReadFile(filePath)
	return string(data)
}

func writeFile(filePath string, data string) {
	ioutil.WriteFile(filePath, []byte(data), 0644)
}