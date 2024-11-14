package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	DefaultAccountKey       = "C2y6yDjf5/R+ob0N8A7Cgv30VRDJIWEHLM+4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw/Jw=="
	EnvPrefix               = "COSMIUM_"
	ExplorerBaseUrlLocation = "/_explorer"
)

var Config = ServerConfig{}

func ParseFlags() {
	host := flag.String("Host", "localhost", "Hostname")
	port := flag.Int("Port", 8081, "Listen port")
	explorerPath := flag.String("ExplorerDir", "", "Path to cosmos-explorer files")
	tlsCertificatePath := flag.String("Cert", "", "Hostname")
	tlsCertificateKey := flag.String("CertKey", "", "Hostname")
	initialDataPath := flag.String("InitialData", "", "Path to JSON containing initial state")
	accountKey := flag.String("AccountKey", DefaultAccountKey, "Account key for authentication")
	disableAuthentication := flag.Bool("DisableAuth", false, "Disable authentication")
	disableTls := flag.Bool("DisableTls", false, "Disable TLS, serve over HTTP")
	persistDataPath := flag.String("Persist", "", "Saves data to given path on application exit")
	debug := flag.Bool("Debug", false, "Runs application in debug mode, this provides additional logging")

	flag.Parse()
	setFlagsFromEnvironment()

	Config.Host = *host
	Config.Port = *port
	Config.ExplorerPath = *explorerPath
	Config.TLS_CertificatePath = *tlsCertificatePath
	Config.TLS_CertificateKey = *tlsCertificateKey
	Config.InitialDataFilePath = *initialDataPath
	Config.PersistDataFilePath = *persistDataPath
	Config.DisableAuth = *disableAuthentication
	Config.DisableTls = *disableTls
	Config.Debug = *debug

	Config.DatabaseAccount = Config.Host
	Config.DatabaseDomain = Config.Host
	Config.DatabaseEndpoint = fmt.Sprintf("https://%s:%d/", Config.Host, Config.Port)
	Config.AccountKey = *accountKey
	Config.ExplorerBaseUrlLocation = ExplorerBaseUrlLocation
}

func setFlagsFromEnvironment() (err error) {
	flag.VisitAll(func(f *flag.Flag) {
		name := EnvPrefix + strings.ToUpper(strings.Replace(f.Name, "-", "_", -1))
		if value, ok := os.LookupEnv(name); ok {
			err2 := flag.Set(f.Name, value)
			if err2 != nil {
				err = fmt.Errorf("failed setting flag from environment: %w", err2)
			}
		}
	})

	return
}
