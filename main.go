package main

import (
	"fmt"

	"github.com/pikami/cosmium/api"
	"github.com/pikami/cosmium/api/config"
	"github.com/pikami/cosmium/internal/repositories"
)

func main() {
	config.ParseFlags()

	if config.Config.DataFilePath != "" {
		repositories.LoadStateFS(config.Config.DataFilePath)
	}

	router := api.CreateRouter()
	router.RunTLS(
		fmt.Sprintf(":%d", config.Config.Port),
		config.Config.TLS_CertificatePath,
		config.Config.TLS_CertificateKey)
}
