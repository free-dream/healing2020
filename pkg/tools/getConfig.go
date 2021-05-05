package tools

import (
	//"fmt"
	"gopkg.in/ini.v1"
)

func GetConfig(section string, keyName string) string {
	cfg, err := ini.Load("config/app.ini")
	if err != nil {
		//fmt.Printf("err:%v",err)
		return ""
	}
	return cfg.Section(section).Key(keyName).String()
}

func IsDebug() bool {
	if GetConfig("debug", "debug") == "false" {
		return false
	}
	return true
}
