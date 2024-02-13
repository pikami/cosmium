package main

import (
	"fmt"

	"github.com/pikami/cosmium/api"
	"github.com/pikami/cosmium/api/config"
)

func main() {
	config.ParseFlags()

	router := api.CreateRouter()
	router.RunTLS(
		fmt.Sprintf(":%d", config.Config.Port),
		config.Config.TLS_CertificatePath,
		config.Config.TLS_CertificateKey)
}
