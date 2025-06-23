package ddns

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/gitsang/ddns/pkg/util/execx"
	netx "github.com/gitsang/ddns/pkg/util/net"
	timex "github.com/gitsang/ddns/pkg/util/time"
	"github.com/gitsang/logi"
)

type Service struct {
	logh slog.Handler

	interval    time.Duration
	dnsProvider DnsProvider
	ddnsConfigs []DdnsConfig
}

type ServiceOptionFunc func(*Service) error

func WithLogHandler(logh slog.Handler) ServiceOptionFunc {
	return func(s *Service) error {
		s.logh = logh
		return nil
	}
}

func WithInterval(interval string) ServiceOptionFunc {
	return func(s *Service) error {
		interval, err := timex.ParseDuration(interval)
		if err != nil {
			return err
		}
		s.interval = interval
		return nil
	}
}

func WithDnsProvider(provider DnsProvider) ServiceOptionFunc {
	return func(s *Service) error {
		s.dnsProvider = provider
		return nil
	}
}

func WithDdnsConfigs(ddnsConfigs ...DdnsConfig) ServiceOptionFunc {
	return func(s *Service) error {
		s.ddnsConfigs = ddnsConfigs
		return nil
	}
}

func defaultService() *Service {
	return &Service{
		logh:        logi.NopHandler,
		interval:    30 * time.Minute,
		dnsProvider: nil,
		ddnsConfigs: []DdnsConfig{},
	}
}

func NewService(optfs ...ServiceOptionFunc) (*Service, error) {
	s := defaultService()
	for _, optf := range optfs {
		if err := optf(s); err != nil {
			return nil, err
		}
	}
	return s, nil
}

func (s *Service) UpdateDns() {
	for _, config := range s.ddnsConfigs {
		if !config.Enable {
			continue
		}
		logger := slog.New(s.logh).With(slog.Any("config", config))

		// get ip
		var ips []string
		switch config.Provider.Type {
		case "interface":
			var err error
			ips, err = netx.GetIpsWithPrefix(
				config.Provider.Interface.Interface,
				config.Provider.Interface.Prefix)
			if err != nil {
				logger.Error("get ip failed", slog.Any("err", err))
				continue
			}
		case "command":
			out, err := execx.RunBash(config.Provider.Command)
			if err != nil {
				logger.Error("get ip failed", slog.String("out", out), slog.Any("err", err))
				continue
			}
			ips = strings.Split(out, "\n")
		}
		logger = logger.With(slog.Any("ip", ips[0]))

		// update dns
		err := s.dnsProvider.UpdateOrCreateRecord(config.Record, ips[0])
		if err != nil {
			logger.Error("update dns failed", slog.Any("err", err))
			continue
		}
	}
}

func (s *Service) Start(ctx context.Context) error {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	for {
		s.UpdateDns()

		select {
		case <-ctx.Done():
			slog.New(s.logh).Info("service is stopping...")
			return nil
		case <-ticker.C:
		}
	}
}
