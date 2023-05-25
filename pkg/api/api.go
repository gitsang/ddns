package api

import (
	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v2/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	"github.com/alibabacloud-go/tea/tea"
	log "github.com/gitsang/golog"
	"go.uber.org/zap"
)

const RequestEndpoint = "dns.aliyuncs.com"

func CreateClient(accessKeyId string, accessKeySecret string) (_result *alidns20150109.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String(RequestEndpoint),
	}

	_result = &alidns20150109.Client{}
	_result, _err = alidns20150109.NewClient(config)
	return _result, _err
}

func FindRecord(client *alidns20150109.Client, domainName, rr, typ string) (
	*alidns20150109.DescribeDomainRecordsResponseBodyDomainRecordsRecord, error) {

	describeDomainRecordsRequest := &alidns20150109.DescribeDomainRecordsRequest{
		DomainName:  tea.String(domainName),
		RRKeyWord:   tea.String(rr),
		TypeKeyWord: tea.String(typ),
		SearchMode:  tea.String("EXACT"),
	}
	resp, err := client.DescribeDomainRecords(describeDomainRecordsRequest)
	if err != nil {
		return nil, err
	}
	records := resp.Body.DomainRecords.Record
	log.Debug("describe domain records success", zap.Reflect("records", records))

	for _, rec := range records {
		if *rec.RR == rr {
			return rec, nil
		}
	}

	return nil, nil
}

func UpdateRecord(client *alidns20150109.Client, id, rr, typ, value string) error {
	updateDomainRecordRequest := &alidns20150109.UpdateDomainRecordRequest{
		RecordId: tea.String(id),
		RR:       tea.String(rr),
		Type:     tea.String(typ),
		Value:    tea.String(value),
	}
	_, err := client.UpdateDomainRecord(updateDomainRecordRequest)
	if err != nil {
		return err
	}

	return nil
}

func CreateRecord(client *alidns20150109.Client, domain, rr, typ, value string) error {
	addDomainRecordRequest := &alidns20150109.AddDomainRecordRequest{
		DomainName: tea.String(domain),
		RR:         tea.String(rr),
		Type:       tea.String(typ),
		Value:      tea.String(value),
	}
	_, err := client.AddDomainRecord(addDomainRecordRequest)
	if err != nil {
		return err
	}

	return nil
}

func UpdateOrCreateRecord(client *alidns20150109.Client, domain, rr, typ, rec string) error {
	logFields := []zap.Field{
		zap.String("domain", domain),
		zap.String("rr", rr),
		zap.String("typ", typ),
		zap.String("rec", rec),
	}
	defer func() {
		log.Info("UpdateOrCreateRecord end", logFields...)
	}()

	// find record
	record, err := FindRecord(client, domain, rr, typ)
	if err != nil {
		return err
	}

	// create or update
	if record == nil { // create
		err = CreateRecord(client, domain, rr, typ, rec)
		if err != nil {
			logFields = append(logFields, zap.Error(err))
			return err
		}
		logFields = append(logFields, zap.String("message", "create record success"))
	} else { // update
		recordId := *record.RecordId
		recordValue := *record.Value
		logFields = append(logFields, zap.String("recordId", recordId), zap.String("recordValue", recordValue))

		if recordValue == rec {
			logFields = append(logFields, zap.String("message", "record not change, skip"))
			return nil
		}

		err = UpdateRecord(client, recordId, rr, typ, rec)
		if err != nil {
			logFields = append(logFields, zap.Error(err))
			return err
		}
		logFields = append(logFields, zap.String("message", "update record success"))
	}

	return nil
}
