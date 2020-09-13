package server

type Config struct {
	VK
	Template Template       `json:"template" yaml:"template"`
	Listener ListenerConfig `json:"listeners" yaml:"listeners"`
}

type Template struct {
	File string `json:"file" yaml:"file"`
	Body string `json:"body" yaml:"body"`
}

type VK struct {
	UserAgent string   `json:"user_agent" yaml:"user_agent"`
	Tokens    []string `json:"tokens" yaml:"tokens"`
}

type ListenerConfig struct {
	Bind string     `json:"bind" yaml:"bind"`
	TLS  *TLSConfig `json:"tls" yaml:"tls"`
}

type TLSConfig struct {
	KeyFile  string `json:"key_file" yaml:"key_file"`
	CertFile string `json:"cert_file" yaml:"cert_file"`
}
