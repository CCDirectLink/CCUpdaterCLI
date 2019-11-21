package local

import (
	"encoding/json"
	"os"
	"io/ioutil"
	"path/filepath"
	"github.com/Masterminds/semver"
	"github.com/CCDirectLink/CCUpdaterCLI"
)

type modPackage struct {
	base string
	loadedMetadata ccmodupdater.PackageMetadata
	dependencies map[string]string
}

func (mp modPackage) Metadata() ccmodupdater.PackageMetadata {
	return mp.loadedMetadata
}

func (mp modPackage) Dependencies() map[string]string {
	nMap := make(map[string]string)
	for id, version := range mp.dependencies {
		nMap[id] = version
	}
	return nMap
}

func (mp modPackage) Remove() error {
	return os.RemoveAll(mp.base)
}

// Ported from cmd/internal/local/modfinder.go
func getModPackage(base string) (ccmodupdater.LocalPackage, error) {
	file, err := os.Open(filepath.Join(base, "package.json"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data struct {
		Name              string             `json:"name"`
		Version           *string            `json:"version"`
		Description       *string            `json:"description"`
		Dependencies      *map[string]string `json:"dependencies"`
		CcmodDependencies *map[string]string `json:"ccmodDependencies"`
	}
	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		return nil, err
	}

	version := semver.MustParse("0.0.0")
	if data.Version != nil {
		ver, err := semver.NewVersion(*data.Version)
		if err == nil {
			version = ver
		}
	}

	var dependencies map[string]string
	if data.CcmodDependencies != nil {
		dependencies = *data.CcmodDependencies
	} else if data.Dependencies != nil {
		dependencies = *data.Dependencies
	}

	metadata := ccmodupdater.PackageMetadata{
		Name: data.Name,
		Type: ccmodupdater.PackageTypeMod,
		Description: "An installed mod.",
		Version: version,
	}
	if data.Description != nil {
		metadata.Description = *data.Description
	}
	if metadata.Name == "Simplify" {
		metadata.Type = ccmodupdater.PackageTypeBase
		metadata.Description = "Assistant to CCLoader."
	}
	
	return modPackage{
		base: base,
		loadedMetadata: metadata,
		dependencies: dependencies,
	}, nil
}

type modPackagePlugin struct {
	dir string
}

// NewModlikePackagePlugin creates a LocalPackagePlugin to scan a given `assets/mods`-like (that or `assets/tools`)
func NewModlikePackagePlugin(game *ccmodupdater.GameInstance, dir string) ccmodupdater.LocalPackagePlugin {
	return modPackagePlugin{
		dir: filepath.Join(game.Base(), dir),
	}
}

func (mpp modPackagePlugin) Packages() []ccmodupdater.LocalPackage {
	dirs, err := ioutil.ReadDir(mpp.dir)
	packages := []ccmodupdater.LocalPackage{}
	if err == nil {
		for _, dir := range dirs {
			if dir.IsDir() {
				mod, err := getModPackage(filepath.Join(mpp.dir, dir.Name()))
				if err == nil {
					packages = append(packages, mod)
				}
			}
		}
	}
	return packages
}
