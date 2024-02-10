package main

import (
	"fmt"

	"github.com/pikami/cosmium/api"
)

func main() {
	fmt.Println("Hello world")

	router := api.CreateRouter()
	router.RunTLS(":8081", "../example.crt", "../example.key")
}
