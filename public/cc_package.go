package public

import (
	"fmt"
	"path/filepath"
	"os"
	"encoding/json"
	"github.com/Masterminds/semver"
)

// VerySecretVerion has to be exported for some reason, but is private.
type VerySecretVerion struct {
	Changelog []struct {
		Version string `json:"version"`
	} `json:"changelog"`
}

type crossCodePackage struct {
	version *semver.Version
}

func (cc crossCodePackage) Metadata() PackageMetadata {
	return PackageMetadata{
		Name: "crosscode",
		Type: PackageTypeBase,
		Description: "CrossCode is the base game itself.",
		Version: cc.version,
	}
}
func (cc crossCodePackage) Remove() error {
	return fmt.Errorf("CrossCode cannot be removed")
}
func (cc crossCodePackage) Dependencies() map[string]string {
	return map[string]string{}
}

// Attempts to get CrossCode as a package.
func (gi *GameInstance) getCrossCodePackage() (LocalPackage, error) {
	// Firstly, find the changelog file. If it doesn't exist this isn't CrossCode.
	changelog, err := os.Open(filepath.Join(gi.base, "./assets/data/changelog.json"))
	if err != nil {
		return nil, fmt.Errorf("CrossCode not found")
	}
	defer changelog.Close()
	
	var log *VerySecretVerion = &VerySecretVerion{}
	
	if err = json.NewDecoder(changelog).Decode(log) ; err != nil {
		return nil, err
	}
	
	if len(log.Changelog) == 0 {
		return nil, fmt.Errorf("Changelog had no entries")
	}
	
	version, err := semver.NewVersion(log.Changelog[0].Version)
	if err != nil {
		return nil, err
	}
	
	return crossCodePackage{
		version: version,
	}, nil
}
