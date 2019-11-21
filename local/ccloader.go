package local

import (
	"fmt"
	"path/filepath"
	"os"
	"github.com/Masterminds/semver"
	"github.com/CCDirectLink/CCUpdaterCLI"
)

type ccLoaderPackage struct {
}

func (cc ccLoaderPackage) Metadata() ccmodupdater.PackageMetadata {
	return ccmodupdater.PackageMetadata{
		Name: "ccloader",
		Type: ccmodupdater.PackageTypeBase,
		Description: "CCLoader is a mod loader.",
		// Please see ccLoaderRemotePackage
		Version: semver.MustParse("1.0.0"),
	}
}
func (cc ccLoaderPackage) Remove() error {
	return fmt.Errorf("CCLoader cannot be automatically removed right now")
}
func (cc ccLoaderPackage) Dependencies() map[string]string {
	return map[string]string{}
}

type ccloaderPackagePlugin struct {
	dir string
}

// NewCCLoaderPackagePlugin creates a LocalPackagePlugin given the game base.
func NewCCLoaderPackagePlugin(game *ccmodupdater.GameInstance) ccmodupdater.LocalPackagePlugin {
	return ccloaderPackagePlugin{
		dir: game.Base(),
	}
}

func (ccl ccloaderPackagePlugin) Packages() []ccmodupdater.LocalPackage {
	proof, err := os.Open(filepath.Join(ccl.dir, "ccloader/index.html"))
	if err != nil {
		return []ccmodupdater.LocalPackage{}
	}
	proof.Close()
	return []ccmodupdater.LocalPackage{
		ccLoaderPackage{},
	}
}
