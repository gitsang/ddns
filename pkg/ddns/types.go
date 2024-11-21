package ddns

type DdnsConfig struct {
	Enable    bool   `json:"enable" yaml:"enable" default:"false"`
	Type      string `json:"type" yaml:"type" default:"A"`
	Domain    string `json:"domain" yaml:"domain" default:"example.com"`
	RR        string `json:"rr" yaml:"rr" default:"example.com"`
	Interface string `json:"interface" yaml:"interface" default:"eth0"`
	Prefix    string `json:"prefix" yaml:"prefix" default:"192.168"`
}
