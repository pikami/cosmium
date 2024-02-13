package config

type ServerConfig struct {
	DatabaseAccount  string
	DatabaseDomain   string
	DatabaseEndpoint string

	ExplorerPath        string
	Port                int
	Host                string
	TLS_CertificatePath string
	TLS_CertificateKey  string
}
