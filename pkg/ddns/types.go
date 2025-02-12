package ddns

type Record struct {
	Type   string `json:"type" yaml:"type" default:"A"`
	Domain string `json:"domain" yaml:"domain" default:"example.com"`
	RR     string `json:"rr" yaml:"rr" default:"example.com"`
}

type InterfaceProvider struct {
	Interface string `json:"interface" yaml:"interface" default:"eth0"`
	Prefix    string `json:"prefix" yaml:"prefix" default:"192.168"`
}

type Provider struct {
	Type      string            `json:"type" yaml:"type" default:"interface"`
	Interface InterfaceProvider `json:"interface" yaml:"interface"`
	Command   string            `json:"command" yaml:"command" default:"ifconfig eth0 | grep inet | awk '{print $2}'"`
}

type DdnsConfig struct {
	Enable   bool     `json:"enable" yaml:"enable" default:"false"`
	Record   Record   `json:"record" yaml:"record"`
	Provider Provider `json:"provider" yaml:"provider"`
}
