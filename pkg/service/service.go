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

	for _, ddns := range config.Cfg.DdnsList {
		recordType := ddns.Record.Type
		value, err := utils.GetIp(ddns.Record.Ipv6, true, ddns.Record.Interface)
		if err != nil {
			log.Error("get interface ip failed", zap.Reflect("ddns", ddns), zap.Error(err))
			continue
		}

		for _, rr := range ddns.RRs {
			record := api.FindRecordByRR(records, rr)
			if record == nil {
				log.Info("record not found, create new (not implement)")
				continue
			}

			recordId := *record.RecordId
			recordValue := *record.Value
			logFields := []zap.Field{
				zap.String("rr", rr),
				zap.String("recordId", recordId),
				zap.String("recordType", recordType),
				zap.String("value", value),
			}

			if recordValue == value {
				log.Info("record not change, skip", logFields...)
				continue
			}

			err = api.UpdateRecord(client, recordId, rr, recordType, value)
			if err != nil {
				log.Error("update record failed", append(logFields, zap.Error(err))...)
			}

			log.Info("update record success", logFields...)
		}
	}
}

func DdnsStart() error {
	ticker := time.Tick(time.Duration(config.Cfg.UpdateIntervalMin) * time.Minute)
	for {
		log.Info("start updateDns")
		UpdateDns()

		select {
		case <-ticker:
		}
	}
}
