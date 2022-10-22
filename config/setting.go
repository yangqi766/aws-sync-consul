package config

import (
	"log"

	"github.com/go-ini/ini"
)

type AWS struct {
	AccessKey string
	SecretKey string
}

type ALI struct {
	AccessKey string
	SecretKey string
}

type Consul struct {
	Address string
	Port    int
}

var (
	AwsSetting    = &AWS{}
	AliSetting    = &ALI{}
	ConsulSetting = &Consul{}
	cfg           *ini.File
)

// Setup initialize the configuration instance
func Setup() {
	var err error
	cfg, err = ini.Load("config/app.ini")
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'config/app.ini': %v", err)
	}

	mapTo("AWS", AwsSetting)
	mapTo("ALI", AliSetting)
	mapTo("Consul", ConsulSetting)
}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
