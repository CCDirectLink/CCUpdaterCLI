package remote
import (
	"fmt"
	"os"
	"path/filepath"
)

// Attempts to install the given installation method to the given directory.
func tryExecuteInstallationMethod(method ccModDBInstallationMethod, newDir string) error {
	if method.Type != "modZip" {
		return fmt.Errorf("Unable to interpret installation method of type %s", method.Type)
	}
	if method.Platform != nil {
		platform := *method.Platform
		if platform != whatPlatformAreWe() {
			return fmt.Errorf("Installation method requires platform %s", platform)
		}
	}
	
	err := os.MkdirAll("installing", os.ModePerm)
	if err != nil {
		return fmt.Errorf("Unable to make temp directory: %s", err.Error())
	}
	defer os.RemoveAll("installing")

	file, err := download(method.URL)
	if err != nil {
		return fmt.Errorf("Unable to download: %s", err.Error())
	}
	defer os.Remove(file.Name())

	dir, err := extract(file)
	if err != nil {
		return fmt.Errorf("Unable to extract: %s", err.Error())
	}
	defer os.RemoveAll(dir)
	
	dirSrc := dir
	if method.Source != nil {
		dirSrc = filepath.Join(dirSrc, *method.Source)
	}
	
	err = copyDir(newDir, dirSrc)
	if err != nil {
		return err
	}
	return nil
}
