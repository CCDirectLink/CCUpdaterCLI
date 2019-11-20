package public
import (
	"github.com/Masterminds/semver"
	"fmt"
	"path/filepath"
	"io/ioutil"
	"os"
)

type modRemotePackage struct {
	data ccModDBMod
	version *semver.Version
}

func (mrp modRemotePackage) Metadata() PackageMetadata {
	return PackageMetadata{
		Name: mrp.data.Name,
		Type: PackageTypeMod,
		Description: mrp.data.Description,
		Version: mrp.version,
	}
}

func (mrp modRemotePackage) Install(game *GameInstance) error {
	err := os.MkdirAll("installing", os.ModePerm)
	if err != nil {
		return err
	}
	defer os.RemoveAll("installing")

	file, err := download(mrp.data.ArchiveLink)
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

	dir, err := extract(file)
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	pkgDir, found, err := findPackage(dir)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("cmd/internal: Could not find package metadata of mod '%s'", mrp.data.Name)
	}

	modDir := filepath.Join(game.Base(), "assets/mods", mrp.data.Name)
	err = copyDir(modDir, pkgDir)
	if err != nil {
		return err
	}

	return nil
}

func findPackage(dir string) (string, bool, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return dir, false, err
	}

	for _, file := range files {
		if !file.IsDir() && file.Name() == "package.json" {
			return dir, true, nil
		}
	}

	for _, file := range files {
		if file.IsDir() {
			res, found, err := findPackage(filepath.Join(dir, file.Name()))
			if err != nil {
				return res, found, err
			}

			if found {
				return res, true, nil
			}
		}
	}

	return dir, false, nil
}
