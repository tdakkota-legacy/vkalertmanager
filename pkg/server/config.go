package server

type ListenerConfig struct {
	Bind string     `json:"bind" yaml:"bind"`
	TLS  *TLSConfig `json:"tls" yaml:"tls"`
}

type TLSConfig struct {
	KeyFile  string `json:"key_file" yaml:"key_file"`
	CertFile string `json:"cert_file" yaml:"cert_file"`
}
