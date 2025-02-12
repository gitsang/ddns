package ddns

type DnsProvider interface {
	UpdateOrCreateRecord(record Record, ip string) error
}
