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

const (
	DataStoreJson   = "json"
	DataStoreBadger = "badger"
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
	logLevel := NewEnumValue("info", []string{"debug", "info", "error", "silent"})
	flag.Var(logLevel, "LogLevel", fmt.Sprintf("Sets the logging level %s", logLevel.AllowedValuesList()))
	dataStore := NewEnumValue("json", []string{DataStoreJson, DataStoreBadger})
	flag.Var(dataStore, "DataStore", fmt.Sprintf("Sets the data store %s", dataStore.AllowedValuesList()))

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
	config.AccountKey = *accountKey
	config.LogLevel = logLevel.value
	config.DataStore = dataStore.value

	config.PopulateCalculatedFields()

	return config
}

func (c *ServerConfig) PopulateCalculatedFields() {
	c.DatabaseAccount = c.Host
	c.DatabaseDomain = c.Host
	c.DatabaseEndpoint = fmt.Sprintf("https://%s:%d/", c.Host, c.Port)
	c.ExplorerBaseUrlLocation = ExplorerBaseUrlLocation

	switch c.LogLevel {
	case "debug":
		logger.SetLogLevel(logger.LogLevelDebug)
	case "info":
		logger.SetLogLevel(logger.LogLevelInfo)
	case "error":
		logger.SetLogLevel(logger.LogLevelError)
	case "silent":
		logger.SetLogLevel(logger.LogLevelSilent)
	default:
		logger.SetLogLevel(logger.LogLevelInfo)
	}

	fileInfo, err := os.Stat(c.PersistDataFilePath)
	if c.PersistDataFilePath != "" && !os.IsNotExist(err) {
		if err != nil {
			logger.ErrorLn("Failed to get file info for persist path:", err)
			os.Exit(1)
		}

		if c.DataStore == DataStoreJson && fileInfo.IsDir() {
			logger.ErrorLn("--Persist cannot be a directory when using json data store")
			os.Exit(1)
		}

		if c.DataStore == DataStoreBadger && !fileInfo.IsDir() {
			logger.ErrorLn("--Persist must be a directory when using Badger data store")
			os.Exit(1)
		}
	}

	if c.DataStore == DataStoreBadger && c.InitialDataFilePath != "" {
		logger.ErrorLn("InitialData option is currently not supported with Badger data store")
		os.Exit(1)
	}
}

func (c *ServerConfig) ApplyDefaultsToEmptyFields() {
	if c.Host == "" {
		c.Host = "localhost"
	}
	if c.Port == 0 {
		c.Port = 8081
	}
	if c.AccountKey == "" {
		c.AccountKey = DefaultAccountKey
	}
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
