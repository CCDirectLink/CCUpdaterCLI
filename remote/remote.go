package remote

import (
	"runtime"
	"encoding/json"
	"net/http"
	"github.com/CCDirectLink/CCUpdaterCLI"
)

const link = "https://raw.githubusercontent.com/20kdc/CCModDB/master/npDatabase.json"

// ccModDB contains data about mods
type ccModDB map[string]ccModDBMod

type ccModDBHash struct {
	SHA256 string `json:"sha256"`
}

// Returns runtime.GOOS, substituting known differences into Node.js platform values.
func whatPlatformAreWe() string {
	sysPlatform := runtime.GOOS
	if sysPlatform == "windows" {
		return "win32"
	}
	return sysPlatform
}

type ccModDBInstallationMethod struct {
	Type string `json:"type"`
	Platform *string `json:"platform"`

	URL string `json:"url"`
	Hash ccModDBHash `json:"hash"`
	Source *string `json:"source"`
}
// ccModDBMod defines the CCModDb;NP mod structure
type ccModDBMod struct {
	Metadata ccmodupdater.PackageMetadata `json:"metadata"`
	Installation []ccModDBInstallationMethod `json:"installation"`
}

var ccModDBData ccModDB

// fetchModDataFromCCModDB
func fetchModDataFromCCModDB() (ccModDB, error) {
	if ccModDBData != nil {
		return ccModDBData, nil
	}

	res, err := http.Get(link)
	if err != nil {
		return nil, err
	}

	ccModDBData = make(ccModDB)
	err = json.NewDecoder(res.Body).Decode(&ccModDBData)
	return ccModDBData, err
}

// GetRemotePackages retrieves all the remote packages that can be found right now.
func GetRemotePackages() (map[string]ccmodupdater.RemotePackage, error) {
	ccmoddb, err := fetchModDataFromCCModDB()
	if err != nil {
		return nil, err
	}
	// Start with CCLoader already in there
	packages := map[string]ccmodupdater.RemotePackage{
		"ccloader": ccLoaderRemotePackage{},
	}
	for _, mod := range ccmoddb {
		if err := mod.Metadata.Verify(); err != nil {
			// Uhoh... Should we warning here????
			continue
		}
		pkg := modRemotePackage{data: mod}
		packages[pkg.Metadata().Name()] = pkg
	}
	return packages, nil
}
