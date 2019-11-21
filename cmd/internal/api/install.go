package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal"
	"github.com/CCDirectLink/CCUpdaterCLI/cmd/commands"
)

//InstallRequest for incoming installation requests
type InstallRequest struct {
	Game  *string  `json:"game"`
	Names []string `json:"names"`
}

//InstallResponse for installation requests
type InstallResponse struct {
	Success bool       `json:"success"`
	Message string     `json:"message,omitempty"`
	Stats   *internal.Stats `json:"stats,omitempty"`
}

//Install a mod via api request
func Install(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", 405)
		return
	}

	setHeaders(w)

	decoder := json.NewDecoder(r.Body)
	stats, err := install(decoder)

	encoder := json.NewEncoder(w)
	if err == nil {
		encoder.Encode(&InstallResponse{
			Success: true,
			Stats:   stats,
		})
	} else {
		encoder.Encode(&InstallResponse{
			Success: false,
			Message: err.Error(),
			Stats:   stats,
		})
	}
}

func install(decoder *json.Decoder) (*internal.Stats, error) {
	var req InstallRequest
	if err := decoder.Decode(&req); err != nil {
		return nil, fmt.Errorf("cmd/internal/api: Could not parse request body: %s", err.Error())
	}

	context, err := internal.NewOnlineContext(req.Game)
	if err != nil {
		return nil, fmt.Errorf("cmd/internal/api: Could not set game flag: %s", err.Error())
	}

	return commands.Install(context, req.Names)
}
