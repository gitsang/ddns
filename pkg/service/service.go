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

	records, err := api.DescribeDomainRecords(client, config.Cfg.Domain)
	if err != nil {
		log.Error("list record failed", zap.Error(err))
		return
	}

	for _, ddns := range config.Cfg.DDNSs {
		if !ddns.Enable {
			continue
		}
		logFields := []zap.Field{zap.Reflect("ddns", ddns)}

		ip, err := utils.GetIpWithPrefix(ddns.Interface, ddns.Prefix)
		if err != nil {
			log.Error("get interface ip failed", append(logFields, zap.Error(err))...)
			continue
		}
		logFields = append(logFields, zap.String("ip", ip))

		record := api.FindRecordByRR(records, ddns.RR)
		if record == nil { // create
			err = api.CreateRecord(client, config.Cfg.Domain, ddns.RR, ddns.Type, ip)
			if err != nil {
				log.Error("create record failed", append(logFields, zap.Error(err))...)
				continue
			}

		} else { // update
			recordId := *record.RecordId
			recordValue := *record.Value
			logFields = append(logFields, zap.String("recordId", recordId), zap.String("recordValue", recordValue))
			if recordValue == ip {
				log.Debug("record not change, skip", logFields...)
				continue
			}

			err = api.UpdateRecord(client, recordId, ddns.RR, ddns.Type, ip)
			if err != nil {
				log.Error("update record failed", append(logFields, zap.Error(err))...)
				continue
			}
		}

		log.Info("update record success", logFields...)
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
