package api

import (
	"encoding/json"
	//"flag"
	"fmt"
	"net/http"
)

//GlobalModsRequest for incoming list available mods requests
type GlobalModsRequest struct {
	Game *string `json:"game"`
}

//GlobalModsResponse contains a list of available mods
type GlobalModsResponse struct {
	Success bool                  `json:"success"`
	Message string                `json:"message,omitempty"`
	//Mods    map[string]global.Mod `json:"mods"`
}

//GetGlobalMods returns all available mods
func GetGlobalMods(w http.ResponseWriter, r *http.Request) {
	//var decoder *json.Decoder
	if r.Method == "POST" {
		//decoder = json.NewDecoder(r.Body)
	}

	setHeaders(w)

	//mods, err := getGlobalMods(decoder)
	err := fmt.Errorf("NYI")

	encoder := json.NewEncoder(w)
	if err == nil {
		encoder.Encode(&GlobalModsResponse{
			Success: true,
			//Mods:    mods,
		})
	} else {
		encoder.Encode(&GlobalModsResponse{
			Success: false,
			Message: err.Error(),
		})
	}
}
/*
func getGlobalMods(decoder *json.Decoder) (map[string]global.Mod, error) {
	if decoder != nil {
		var req GlobalModsRequest
		if err := decoder.Decode(&req); err != nil {
			return nil, fmt.Errorf("cmd/internal/api: Could not parse request body: %s", err.Error())
		}

		if req.Game != nil {
			if err := flag.Set("game", *req.Game); err != nil {
				return nil, fmt.Errorf("cmd/internal/api: Could set game flag: %s", err.Error())
			}
		}
	}

	res, err := global.FetchModData()
	if err != nil {
		return nil, err
	}

	return res.Mods, nil
}
*/
