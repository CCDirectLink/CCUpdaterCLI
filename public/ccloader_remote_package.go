package public

import (
	"path/filepath"
	"github.com/Masterminds/semver"
	"fmt"
	"os"
)

type ccLoaderRemotePackage struct {
}

func (cc ccLoaderRemotePackage) Metadata() PackageMetadata {
	return PackageMetadata{
		Name: "ccloader",
		Type: PackageTypeBase,
		Description: "CCLoader is a mod loader.",
		// Please see ccLoaderPackage (note this is higher than that)
		Version: semver.MustParse("1.0.1"),
	}
}
func (cc ccLoaderRemotePackage) Dependencies() map[string]string {
	return map[string]string{}
}
func (cc ccLoaderRemotePackage) Install(game *GameInstance) error {
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
