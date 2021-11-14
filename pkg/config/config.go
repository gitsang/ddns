package config

import (
	"time"

	log "github.com/gitsang/golog"
	"github.com/jinzhu/configor"
	"go.uber.org/zap"
)

type DDNS struct {
	RRs    []string
	Record struct {
		Type      string
		Interface string
		Ipv6      bool
	}
}

type Config struct {
	AccessKeyId       string
	AccessKeySecret   string
	Domain            string
	UpdateIntervalMin int
	DdnsList          []DDNS
}

var Cfg Config

func LoadConfig(file string) error {
	err := configor.New(&configor.Config{
		ENVPrefix:          "DDNS",
		AutoReload:         true,
		AutoReloadInterval: time.Minute,
		AutoReloadCallback: func(config interface{}) {
			log.Debug("config auto reload", zap.Reflect("Cfg", Cfg))
		},
	}).Load(&Cfg, file)
	if err != nil {
		return err
	}

	log.Info("load config success", zap.Reflect("Cfg", Cfg))
	return nil
}
