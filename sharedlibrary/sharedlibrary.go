package main

import "C"
import (
	"encoding/json"

	"github.com/pikami/cosmium/api"
	"github.com/pikami/cosmium/api/config"
	"github.com/pikami/cosmium/internal/repositories"
)

var currentServer *api.Server

//export Configure
func Configure(configurationJSON *C.char) bool {
	var configuration config.ServerConfig
	err := json.Unmarshal([]byte(C.GoString(configurationJSON)), &configuration)
	if err != nil {
		return false
	}
	config.Config = configuration
	return true
}

//export InitializeRepository
func InitializeRepository() {
	repositories.InitializeRepository()
}

//export StartAPI
func StartAPI() {
	currentServer = api.StartAPI()
}

//export StopAPI
func StopAPI() {
	if currentServer == nil {
		currentServer.StopServer <- true
		currentServer = nil
	}
}

//export GetState
func GetState() *C.char {
	stateJSON, err := repositories.GetState()
	if err != nil {
		return nil
	}
	return C.CString(stateJSON)
}

func main() {}
