package local

import (
	"fmt"
	"path/filepath"
	"os"
	"encoding/json"
	"github.com/Masterminds/semver"
	"github.com/CCDirectLink/CCUpdaterCLI"
)

type crossCodeChangeLog struct {
	Changelog []struct {
		Version string `json:"version"`
	} `json:"changelog"`
}

type crossCodePackage struct {
	version *semver.Version
}

func (cc crossCodePackage) Metadata() ccmodupdater.PackageMetadata {
	metadata := ccmodupdater.PackageMetadata{}
	metadata["name"] = "crosscode"
	metadata["ccmodType"] = "base"
	metadata["ccmodHumanName"] = "CrossCode"
	metadata["description"] = "CrossCode is the base game itself."
	metadata["version"] = cc.version.Original()
	return metadata
}
func (cc crossCodePackage) Remove() error {
	return fmt.Errorf("CrossCode cannot be removed")
}

type crossCodePackagePlugin struct {
	pkg crossCodePackage
}
func (ccp crossCodePackagePlugin) Packages() []ccmodupdater.LocalPackage {
	return []ccmodupdater.LocalPackage{
		ccp.pkg,
	}
}

// Attempts to get CrossCode as a package.
func NewCrossCodePackagePlugin(game *ccmodupdater.GameInstance) (ccmodupdater.LocalPackagePlugin, error) {
	// Firstly, find the changelog file. If it doesn't exist this isn't CrossCode.
	changelogPath := filepath.Join(game.Base(), "./assets/data/changelog.json")
	changelog, err := os.Open(changelogPath)
	if err != nil {
		return nil, fmt.Errorf("%s unopenable: %s", changelogPath, err.Error())
	}
	defer changelog.Close()
	
	var log *crossCodeChangeLog = &crossCodeChangeLog{}
	
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
	
	return crossCodePackagePlugin{
		crossCodePackage{
			version: version,
		},
	}, nil
}
