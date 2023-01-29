package config

import (
	"time"

	log "github.com/gitsang/golog"
	"github.com/jinzhu/configor"
	"go.uber.org/zap"
)

var Version = "unknown"

type Config struct {
	Log struct {
		Level string `default:"info"`
		File  string `default:"/opt/ddns/log/ddns.log"`
	}

	AccessKeyId       string
	AccessKeySecret   string
	UpdateIntervalMin int `default:"60"`

	DDNSs []struct {
		Enable    bool
		Domain    string
		Type      string
		RR        string
		Interface string
		Prefix    string
	}
}

var Cfg Config

func LoadConfig(file string) error {
	err := configor.New(&configor.Config{
		ENVPrefix:          "DDNS",
		AutoReload:         true,
		AutoReloadInterval: time.Minute,
		AutoReloadCallback: func(config interface{}) {
			log.Info("config auto reload", zap.Reflect("Cfg", Cfg))
		},
	}).Load(&Cfg, file)
	if err != nil {
		return err
	}

	log.Debug("load config success", zap.Reflect("Cfg", Cfg))
	return nil
}
