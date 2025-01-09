package tlsprovider

import (
	"crypto/tls"

	"github.com/pikami/cosmium/internal/logger"
)

func GetDefaultTlsConfig() *tls.Config {
	cert, err := tls.X509KeyPair([]byte(certificate), []byte(certificateKey))
	if err != nil {
		logger.ErrorLn("Failed to parse certificate and key:", err)
		return &tls.Config{}
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
}
