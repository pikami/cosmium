package config

import (
	"flag"
	"fmt"
)

const (
	DefaultAccountKey = "C2y6yDjf5/R+ob0N8A7Cgv30VRDJIWEHLM+4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw/Jw=="
)

var Config = ServerConfig{}

func ParseFlags() {
	host := flag.String("Host", "localhost", "Hostname")
	port := flag.Int("Port", 8081, "Listen port")
	explorerPath := flag.String("ExplorerDir", "/home/pk/pro/cosmos-explorer/dist", "Path to cosmos-explorer files")
	tlsCertificatePath := flag.String("Cert", "../example.crt", "Hostname")
	tlsCertificateKey := flag.String("CertKey", "../example.key", "Hostname")
	initialDataPath := flag.String("InitialData", "", "Path to JSON containing initial state")
	accountKey := flag.String("AccountKey", DefaultAccountKey, "Account key for authentication")
	disableAuthentication := flag.Bool("DisableAuth", false, "Disable authentication")

	flag.Parse()

	Config.Host = *host
	Config.Port = *port
	Config.ExplorerPath = *explorerPath
	Config.TLS_CertificatePath = *tlsCertificatePath
	Config.TLS_CertificateKey = *tlsCertificateKey
	Config.DataFilePath = *initialDataPath
	Config.DisableAuth = *disableAuthentication

	Config.DatabaseAccount = Config.Host
	Config.DatabaseDomain = Config.Host
	Config.DatabaseEndpoint = fmt.Sprintf("https://%s:%d/", Config.Host, Config.Port)
	Config.AccountKey = *accountKey
}
