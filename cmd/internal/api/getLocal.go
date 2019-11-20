package api

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal/local"
	"github.com/CCDirectLink/CCUpdaterCLI/public"
)

//LocalModsRequest for incoming installed mod list requests
type LocalModsRequest struct {
	Game *string `json:"game"`
}

//LocalModsResponse contains a list of installed mods
type LocalModsResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Mods    []public.PackageMetadata `json:"mods"`
}

//GetLocalMods returns all installed mods
func GetLocalMods(w http.ResponseWriter, r *http.Request) {
	var decoder *json.Decoder
	if r.Method == "POST" {
		decoder = json.NewDecoder(r.Body)
	}

	setHeaders(w)

	mods, err := getLocalMods(decoder)

	encoder := json.NewEncoder(w)
	if err == nil {
		encoder.Encode(&LocalModsResponse{
			Success: true,
			Mods:    mods,
		})
	} else {
		encoder.Encode(&LocalModsResponse{
			Success: false,
			Message: err.Error(),
		})
	}
}

func getLocalMods(decoder *json.Decoder) ([]public.PackageMetadata, error) {
	if decoder != nil {
		var req LocalModsRequest
		if err := decoder.Decode(&req); err != nil {
			return nil, fmt.Errorf("cmd/internal/api: Could not parse request body: %s", err.Error())
		}

		if req.Game != nil {
			if err := flag.Set("game", *req.Game); err != nil {
				return nil, fmt.Errorf("cmd/internal/api: Could set game flag: %s", err.Error())
			}
		}
	}

	game, err := local.GetGame()
	if err != nil {
		return nil, err
	}

	total := []public.PackageMetadata{}
	
	for _, v := range game.Packages() {
		total = append(total, v.Metadata())
	}
	
	return total, nil
}
