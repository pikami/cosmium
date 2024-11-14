package config

type ServerConfig struct {
	DatabaseAccount  string
	DatabaseDomain   string
	DatabaseEndpoint string
	AccountKey       string

	ExplorerPath            string
	Port                    int
	Host                    string
	TLS_CertificatePath     string
	TLS_CertificateKey      string
	InitialDataFilePath     string
	PersistDataFilePath     string
	DisableAuth             bool
	DisableTls              bool
	Debug                   bool
	ExplorerBaseUrlLocation string
}
