package local

import (
	"encoding/json"
	"os"
	"io/ioutil"
	"path/filepath"
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

	metadata := make(ccmodupdater.PackageMetadata)
	err = json.NewDecoder(file).Decode(&metadata)
	if err != nil {
		return nil, err
	}
	
	if err = metadata.Verify(); err != nil {
		return nil, err
	}

	// Still have to mess with the metadata just a tad.
	// Specifically, we need to find mods which are connected into CCLoader, because these mods are all special
	if (metadata.Name() == "Simplify") || (metadata.Name() == "CCLoader display version") || (metadata.Name() == "OpenDevTools") {
		metadata["ccmodType"] = "base"
		metadata["description"] = "Assistant to CCLoader."
		// While this is false, it is also quite necessary because otherwise Simplify's CCLoader dep. messes with things (not good)
		delete(metadata, "dependencies")
		delete(metadata, "ccmodDependencies")
	}
	
	return modPackage{
		base: base,
		loadedMetadata: metadata,
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
