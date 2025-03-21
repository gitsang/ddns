package alidns

import (
	"log/slog"
	"time"

	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/gitsang/ddns/pkg/ddns"
)

func (s *DnsProvider) FindRecord(domain, rr, typ string) (*alidns20150109.DescribeDomainRecordsResponseBodyDomainRecordsRecord, error) {
	var (
		entryTime = time.Now()
		logger    = slog.New(s.logh).With(
			slog.String("domain", domain),
			slog.String("rr", rr),
			slog.String("typ", typ),
		)
	)

	defer func() {
		logger = logger.With(slog.String("cost", time.Since(entryTime).String()))
		logger.Debug("end")
	}()

	describeDomainRecordsRequest := &alidns20150109.DescribeDomainRecordsRequest{
		DomainName: tea.String(domain),
		RRKeyWord:  tea.String(rr),
		Type:       tea.String(typ),
	}
	resp, err := s.client.DescribeDomainRecords(describeDomainRecordsRequest)
	if err != nil {
		logger = logger.With(slog.Any("err", err))
		return nil, err
	}

	records := resp.Body.DomainRecords.Record
	logger = logger.With(slog.Any("records", records))

	for _, rec := range records {
		if *rec.RR == rr {
			logger = logger.With(slog.Any("rec", rec))
			return rec, nil
		}
	}

	return nil, nil
}

func (s *DnsProvider) UpdateRecord(id, rr, typ, value string) error {
	var (
		entryTime = time.Now()
		logger    = slog.New(s.logh).With(
			slog.String("id", id),
			slog.String("rr", rr),
			slog.String("typ", typ),
			slog.String("value", value),
		)
	)

	defer func() {
		logger = logger.With(slog.String("cost", time.Since(entryTime).String()))
		logger.Debug("end")
	}()

	updateDomainRecordRequest := &alidns20150109.UpdateDomainRecordRequest{
		RecordId: tea.String(id),
		RR:       tea.String(rr),
		Type:     tea.String(typ),
		Value:    tea.String(value),
	}
	_, err := s.client.UpdateDomainRecord(updateDomainRecordRequest)
	if err != nil {
		logger = logger.With(slog.Any("err", err))
		return err
	}

	return nil
}

func (s *DnsProvider) CreateRecord(domain, rr, typ, value string) error {
	var (
		entryTime = time.Now()
		logger    = slog.New(s.logh).With(
			slog.String("domain", domain),
			slog.String("rr", rr),
			slog.String("typ", typ),
			slog.String("value", value),
		)
	)

	defer func() {
		logger = logger.With(slog.String("cost", time.Since(entryTime).String()))
		logger.Debug("end")
	}()

	addDomainRecordRequest := &alidns20150109.AddDomainRecordRequest{
		DomainName: tea.String(domain),
		RR:         tea.String(rr),
		Type:       tea.String(typ),
		Value:      tea.String(value),
	}
	_, err := s.client.AddDomainRecord(addDomainRecordRequest)
	if err != nil {
		logger = logger.With(slog.Any("err", err))
		return err
	}

	return nil
}

func (s *DnsProvider) UpdateOrCreateRecord(record ddns.Record, value string) error {
	var (
		entryTime = time.Now()
		logger    = slog.New(s.logh).With(
			slog.Any("record", record),
			slog.String("value", value),
		)
	)

	defer func() {
		logger = logger.With(slog.String("cost", time.Since(entryTime).String()))
		logger.Debug("end")
	}()

	// find record
	aliRecord, err := s.FindRecord(record.Domain, record.RR, record.Type)
	if err != nil {
		return err
	}

	// create if not exist
	if aliRecord == nil {
		err = s.CreateRecord(record.Domain, record.RR, record.Type, value)
		if err != nil {
			logger = logger.With(slog.Any("err", err))
			return err
		}
		logger.Info("create record success")
		return nil
	}

	// update if exist
	recordId := *aliRecord.RecordId
	recordValue := *aliRecord.Value
	logger = logger.With(slog.String("recordId", recordId), slog.String("recordValue", recordValue))
	if recordValue == value {
		logger.Info("record not change, skip")
		return nil
	}

	err = s.UpdateRecord(recordId, record.RR, record.Type, value)
	if err != nil {
		logger = logger.With(slog.Any("err", err))
		return err
	}
	logger.Info("update record success")

	return nil
}
