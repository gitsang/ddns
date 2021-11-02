package api

import (
	"ddns/pkg/utils"

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

func DescribeDomainRecords(client *alidns20150109.Client, domainName string) (
	[]*alidns20150109.DescribeDomainRecordsResponseBodyDomainRecordsRecord, error) {

	describeDomainRecordsRequest := &alidns20150109.DescribeDomainRecordsRequest{
		DomainName: tea.String(domainName),
	}
	resp, err := client.DescribeDomainRecords(describeDomainRecordsRequest)
	if err != nil {
		return nil, err
	}
	records := resp.Body.DomainRecords.Record

	log.Debug(utils.FUNCTION()+"success", zap.Reflect("records", records))
	return records, nil
}

func FindRecordByRR(
	records []*alidns20150109.DescribeDomainRecordsResponseBodyDomainRecordsRecord,
	rr string) *alidns20150109.DescribeDomainRecordsResponseBodyDomainRecordsRecord {

	for _, rec := range records {
		if *rec.RR == rr {
			return rec
		}
	}
	return nil
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
