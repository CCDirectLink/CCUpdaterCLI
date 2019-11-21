package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CCDirectLink/CCUpdaterCLI/cmd/commands"
	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal"
)

//UpdateRequest for incoming update requests
type UpdateRequest struct {
	Game  *string  `json:"game"`
	Names []string `json:"names"`
}

//UpdateResponse for update requests
type UpdateResponse struct {
	Success bool       `json:"success"`
	Message string     `json:"message,omitempty"`
	Stats   *internal.Stats `json:"stats,omitempty"`
}

//Update a mod via api request
func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", 405)
		return
	}

	setHeaders(w)

	decoder := json.NewDecoder(r.Body)
	stats, err := update(decoder)

	encoder := json.NewEncoder(w)
	if err == nil {
		encoder.Encode(&UpdateResponse{
			Success: true,
			Stats:   stats,
		})
	} else {
		encoder.Encode(&UpdateResponse{
			Success: false,
			Message: err.Error(),
			Stats:   stats,
		})
	}
}

func update(decoder *json.Decoder) (*internal.Stats, error) {
	var req UpdateRequest
	if err := decoder.Decode(&req); err != nil {
		return nil, fmt.Errorf("cmd/internal/api: Could not parse request body: %s", err.Error())
	}

	context, err := internal.NewOnlineContext(req.Game)
	if err != nil {
		return nil, fmt.Errorf("cmd/internal/api: Could not set game flag: %s", err.Error())
	}

	return commands.Update(context, req.Names)
}
