package api

import (
	"encoding/json"
	"net/http"

	"github.com/CCDirectLink/CCUpdaterCLI"
	"github.com/CCDirectLink/CCUpdaterCLI/remote"
)

//GlobalModsResponse contains a list of available mods
type GlobalModsResponse struct {
	Success bool                  `json:"success"`
	Message string                `json:"message,omitempty"`
	Mods    map[string]ccmodupdater.PackageMetadata `json:"mods"`
}

//GetGlobalMods returns all available mods
func GetGlobalMods(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)

	mods, err := getGlobalMods()

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

func getGlobalMods() (map[string]ccmodupdater.PackageMetadata, error) {
	res, err := remote.GetRemotePackages()
	if err != nil {
		return nil, err
	}

	metadata := map[string]ccmodupdater.PackageMetadata{}
	for k, v := range res {
		metadata[k] = v.Metadata()
	}
	return metadata, nil
}
