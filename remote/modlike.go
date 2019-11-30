package remote
import (
	"fmt"
	"path/filepath"
	"github.com/CCDirectLink/CCUpdaterCLI"
)

type modRemotePackage struct {
	data ccModDBMod
}

func (mrp modRemotePackage) Metadata() ccmodupdater.PackageMetadata {
	return mrp.data.Metadata
}

func (mrp modRemotePackage) Install(game *ccmodupdater.GameInstance) error {
	errors := []error{};
	target := filepath.Join(game.Base(), "assets/mods", mrp.data.Metadata.Name())
	for _, method := range mrp.data.Installation {
		err := tryExecuteInstallationMethod(method, target)
		if err != nil {
			errors = append(errors, err)
		} else {
			return nil
		}
	}

	if len(errors) == 1 {
		return errors[0]
	}
	
	return fmt.Errorf("All installation methods failed.")
}
