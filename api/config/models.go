package config

type ServerConfig struct {
	DatabaseAccount  string `json:"databaseAccount"`
	DatabaseDomain   string `json:"databaseDomain"`
	DatabaseEndpoint string `json:"databaseEndpoint"`
	AccountKey       string `json:"accountKey"`

	ExplorerPath            string `json:"explorerPath"`
	Port                    int    `json:"port"`
	Host                    string `json:"host"`
	TLS_CertificatePath     string `json:"tlsCertificatePath"`
	TLS_CertificateKey      string `json:"tlsCertificateKey"`
	InitialDataFilePath     string `json:"initialDataFilePath"`
	PersistDataFilePath     string `json:"persistDataFilePath"`
	DisableAuth             bool   `json:"disableAuth"`
	DisableTls              bool   `json:"disableTls"`
	LogLevel                string `json:"logLevel"`
	ExplorerBaseUrlLocation string `json:"explorerBaseUrlLocation"`
}
