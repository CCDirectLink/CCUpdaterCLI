package public

// Moved from cmd/internal/global/ccmoddb.go

import (
	"encoding/json"
	"net/http"
	"github.com/Masterminds/semver"
)

const link = "https://raw.githubusercontent.com/CCDirectLink/CCModDB/master/mods.json"

// ccModDB contains data about mods
type ccModDB struct {
	Mods map[string]ccModDBMod `json:"mods"`
}

// ccModDBMod defines the CCModDb mod structure
type ccModDBMod struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	License     *string `json:"license"`
	Page        []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	ArchiveLink string `json:"archive_link"`
	Hash        struct {
		Sha256 string `json:"sha256"`
	} `json:"hash"`
	Version string `json:"version"`
	// Not sure how this is supposed to be used. - 20kdc
	Dir     *struct {
		Any string `json:"any"`
	} `json:"dir"`
}

var ccModDBData *ccModDB

// fetchModDataFromCCModDB
func fetchModDataFromCCModDB() (*ccModDB, error) {
	if ccModDBData != nil {
		return ccModDBData, nil
	}

	res, err := http.Get(link)
	if err != nil {
		return nil, err
	}

	ccModDBData = &ccModDB{}
	err = json.NewDecoder(res.Body).Decode(ccModDBData)
	return ccModDBData, err
}

// GetRemotePackages retrieves all the remote packages that can be found right now.
func GetRemotePackages() (map[string]RemotePackage, error) {
	ccmoddb, err := fetchModDataFromCCModDB()
	if err != nil {
		return nil, err
	}
	// Start with CCLoader already in there
	packages := map[string]RemotePackage{
		"ccloader": ccLoaderRemotePackage{},
	}
	for _, mod := range ccmoddb.Mods {
		version, err := semver.NewVersion(mod.Version)
		if err != nil {
			continue
		}
		pkg := modRemotePackage{
			data: mod,
			version: version,
		}
		packages[pkg.Metadata().Name] = pkg
	}
	return packages, nil
}
