package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/CCDirectLink/CCUpdaterCLI"
	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal"
)

//GlobalModsRequest for incoming list available mods requests
type GlobalModsRequest struct {
	Game *string `json:"game"`
}

//GlobalModsResponse contains a list of available mods
type GlobalModsResponse struct {
	Success bool                  `json:"success"`
	Message string                `json:"message,omitempty"`
	Mods    map[string]ccmodupdater.PackageMetadata `json:"mods"`
}

//GetGlobalMods returns all available mods
func GetGlobalMods(w http.ResponseWriter, r *http.Request) {
	//var decoder *json.Decoder
	if r.Method == "POST" {
		//decoder = json.NewDecoder(r.Body)
	}

	setHeaders(w)

	mods, err := getGlobalMods(decoder)

	encoder := json.NewEncoder(w)
	if err == nil {
		encoder.Encode(&GlobalModsResponse{
			Success: true,
			Mods:    mods,
		})
	} else {
		encoder.Encode(&GlobalModsResponse{
			Success: false,
			Message: err.Error(),
		})
	}
}

func getGlobalMods(decoder *json.Decoder) (map[string]ccmodupdater.PackageMetadata, error) {
	var dir *string = nil
	if decoder != nil {
		var req GlobalModsRequest
		if err := decoder.Decode(&req); err != nil {
			return nil, fmt.Errorf("cmd/internal/api: Could not parse request body: %s", err.Error())
		}
		dir = req.Game
	}

	res, err := internal.NewOnlineContext(dir)
	if err != nil {
		return nil, err
	}

	metadata := map[string]ccmodupdater.PackageMetadata{}
	for k, v := range res.RemotePackages() {
		metadata[k] = v.Metadata()
	}
	return metadata, nil
}
