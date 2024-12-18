package config

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/pikami/cosmium/internal/logger"
)

const (
	DefaultAccountKey       = "C2y6yDjf5/R+ob0N8A7Cgv30VRDJIWEHLM+4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw/Jw=="
	EnvPrefix               = "COSMIUM_"
	ExplorerBaseUrlLocation = "/_explorer"
)

func ParseFlags() ServerConfig {
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

	config := ServerConfig{}
	config.Host = *host
	config.Port = *port
	config.ExplorerPath = *explorerPath
	config.TLS_CertificatePath = *tlsCertificatePath
	config.TLS_CertificateKey = *tlsCertificateKey
	config.InitialDataFilePath = *initialDataPath
	config.PersistDataFilePath = *persistDataPath
	config.DisableAuth = *disableAuthentication
	config.DisableTls = *disableTls
	config.Debug = *debug
	config.AccountKey = *accountKey

	config.PopulateCalculatedFields()

	return config
}

func (c *ServerConfig) PopulateCalculatedFields() {
	c.DatabaseAccount = c.Host
	c.DatabaseDomain = c.Host
	c.DatabaseEndpoint = fmt.Sprintf("https://%s:%d/", c.Host, c.Port)
	c.ExplorerBaseUrlLocation = ExplorerBaseUrlLocation
	logger.EnableDebugOutput = c.Debug
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
