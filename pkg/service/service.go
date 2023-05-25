package service

import (
	"ddns/pkg/api"
	"ddns/pkg/config"
	"ddns/pkg/utils"
	"time"

	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v2/client"
	log "github.com/gitsang/golog"
	"go.uber.org/zap"
)

var client *alidns20150109.Client

func UpdateDns() {
	var err error
	client, err = api.CreateClient(config.Cfg.AccessKeyId, config.Cfg.AccessKeySecret)
	if err != nil {
		log.Error("create client failed", zap.Error(err))
		return
	}

	for _, ddns := range config.Cfg.DDNSs {
		if !ddns.Enable {
			return
		}
		logFields := []zap.Field{zap.Reflect("ddns", ddns)}

		// get ip
		ip, err := utils.GetIpWithPrefix(ddns.Interface, ddns.Prefix)
		if err != nil {
			log.Error("get interface ip failed", append(logFields, zap.Error(err))...)
			continue
		}
		logFields = append(logFields, zap.String("ip", ip))

		err = api.UpdateOrCreateRecord(client, ddns.Domain, ddns.RR, ddns.Type, ip)
		if err != nil {
			log.Error("update or create record failed", append(logFields, zap.Error(err))...)
			continue
		}
	}
}

func DdnsStart() error {
	ticker := time.Tick(time.Duration(config.Cfg.UpdateIntervalMin) * time.Minute)
	for {
		UpdateDns()

		select {
		case <-ticker:
		}
	}
}
