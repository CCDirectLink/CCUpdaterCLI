package remote

import (
	"path/filepath"
	"fmt"
	"os"
	"github.com/CCDirectLink/CCUpdaterCLI"
)

type ccLoaderRemotePackage struct {
}

func (cc ccLoaderRemotePackage) Metadata() ccmodupdater.PackageMetadata {
	metadata := ccmodupdater.PackageMetadata{}
	metadata["name"] = "ccloader"
	metadata["ccmodHumanName"] = "CCLoader"
	metadata["description"] = "CCLoader is a mod loader."
	// This is 0.0.1 above local.ccloader's version to ensure updates are always considered available for now
	metadata["version"] = "1.0.1"
	return metadata
}

func (cc ccLoaderRemotePackage) Install(game *ccmodupdater.GameInstance) error {
	err := os.MkdirAll("installing", os.ModePerm)
	if err != nil {
		return err
	}
	defer os.RemoveAll("installing")

	downloadFile, err := download("https://github.com/CCDirectLink/CCLoader/archive/master.zip")
	if err != nil {
		return fmt.Errorf("In ccloader download: %s", err.Error())
	}
	// target dir is CCLoader-master
	extractDir, err := extract(downloadFile)
	if err != nil {
		return fmt.Errorf("In ccloader extract: %s", err.Error())
	}
	err = copyDir(game.Base(), filepath.Join(extractDir, "CCLoader-master"))
	if err != nil {
		return fmt.Errorf("In ccloader copy from %s: %s", extractDir, err.Error())
	}
	return nil
}
