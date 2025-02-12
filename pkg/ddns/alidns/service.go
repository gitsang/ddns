package alidns

import (
	"log/slog"

	"github.com/gitsang/ddns/pkg/logi"

	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v2/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	"github.com/alibabacloud-go/tea/tea"
)

type DnsProvider struct {
	logh   slog.Handler
	client *alidns20150109.Client
}

type DnsProviderOptionFunc func(*DnsProvider) error

func WithLogHandler(logh slog.Handler) DnsProviderOptionFunc {
	return func(s *DnsProvider) error {
		s.logh = logh
		return nil
	}
}

func WithAliClient(endpoint, accessKeyId, accessKeySecret string) DnsProviderOptionFunc {
	return func(s *DnsProvider) error {
		config := &openapi.Config{
			AccessKeyId:     tea.String(accessKeyId),
			AccessKeySecret: tea.String(accessKeySecret),
			Endpoint:        tea.String(endpoint),
		}
		client, err := alidns20150109.NewClient(config)
		if err != nil {
			slog.New(s.logh).Error("create client failed", slog.Any("err", err))
			return err
		}
		s.client = client
		return nil
	}
}

func defaultDnsProvider() *DnsProvider {
	return &DnsProvider{
		client: nil,
		logh:   logi.NopHandler,
	}
}

func NewDnsProvider(optfs ...DnsProviderOptionFunc) (*DnsProvider, error) {
	s := defaultDnsProvider()
	for _, optf := range optfs {
		if err := optf(s); err != nil {
			return nil, err
		}
	}
	return s, nil
}
