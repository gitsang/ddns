package aliddns

import (
	"context"
	"log/slog"
	"time"

	"github.com/gitsang/ddns/pkg/ddns"
	"github.com/gitsang/ddns/pkg/logi"
	netx "github.com/gitsang/ddns/pkg/util/net"
	timex "github.com/gitsang/ddns/pkg/util/time"

	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v2/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	"github.com/alibabacloud-go/tea/tea"
	log "github.com/gitsang/golog"
	"go.uber.org/zap"
)

type Service struct {
	ctx    context.Context
	cancel context.CancelFunc
	client *alidns20150109.Client
	logh   slog.Handler

	interval    time.Duration
	ddnsConfigs []ddns.DdnsConfig
}

type AliDdnsServiceOptionFunc func(*Service) error

func WithLogHandler(logh slog.Handler) AliDdnsServiceOptionFunc {
	return func(s *Service) error {
		s.logh = logh
		return nil
	}
}

func WithAliClient(endpoint, accessKeyId, accessKeySecret string) AliDdnsServiceOptionFunc {
	return func(s *Service) error {
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

func WithInterval(interval string) AliDdnsServiceOptionFunc {
	return func(s *Service) error {
		interval, err := timex.ParseDuration(interval)
		if err != nil {
			return err
		}
		s.interval = interval
		return nil
	}
}

func WithDdnsConfigs(ddnsConfigs ...ddns.DdnsConfig) AliDdnsServiceOptionFunc {
	return func(s *Service) error {
		s.ddnsConfigs = ddnsConfigs
		return nil
	}
}

func defaultService() *Service {
	return &Service{
		client:      nil,
		logh:        logi.NopHandler,
		ddnsConfigs: []ddns.DdnsConfig{},
	}
}

func NewService(optfs ...AliDdnsServiceOptionFunc) (*Service, error) {
	s := defaultService()
	for _, optf := range optfs {
		if err := optf(s); err != nil {
			return nil, err
		}
	}
	return s, nil
}

func (s *Service) UpdateDns() {
	for _, ddns := range s.ddnsConfigs {
		if !ddns.Enable {
			return
		}
		logFields := []zap.Field{zap.Reflect("ddns", ddns)}

		// get ip
		ip, err := netx.GetIpWithPrefix(ddns.Interface, ddns.Prefix)
		if err != nil {
			log.Error("get interface ip failed", append(logFields, zap.Error(err))...)
			continue
		}
		logFields = append(logFields, zap.String("ip", ip))

		err = s.UpdateOrCreateRecord(ddns.Domain, ddns.RR, ddns.Type, ip)
		if err != nil {
			log.Error("update or create record failed", append(logFields, zap.Error(err))...)
			continue
		}
	}
}

func (s *Service) Start(ctx context.Context) error {
	s.ctx, s.cancel = context.WithCancel(ctx)
	defer s.cancel()

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	for ; true; <-ticker.C {
		select {
		case <-s.ctx.Done():
			return nil
		default:
		}

		s.UpdateDns()
	}

	return nil
}

func (s *Service) Stop() error {
	s.cancel()

	return nil
}
