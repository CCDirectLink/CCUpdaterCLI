package api

import (
	"fmt"
	"net/http"

	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal/api"
)

//Start api server
func Start() {
	StartAt("localhost", 9392)
}

//StartAt host and port
func StartAt(host string, port int) {
	url := fmt.Sprintf("%s:%d", host, port)
	fmt.Printf("API server listening on %s\n", url)

	http.HandleFunc("/api/v1/install", api.Install)
	http.HandleFunc("/api/v1/uninstall", api.Uninstall)
	http.HandleFunc("/api/v1/update", api.Update)
	http.HandleFunc("/api/v1/get/local", api.GetLocalMods)
	http.HandleFunc("/api/v1/get/global", api.GetGlobalMods)

	http.ListenAndServe(url, nil)
}
