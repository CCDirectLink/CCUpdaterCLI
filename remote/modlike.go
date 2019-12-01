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
	typ := mrp.data.Metadata.Type()
	
	pkgName := mrp.data.Metadata.Name()

	// -- Work out installation details --
	
	var target string
	if typ == ccmodupdater.PackageTypeMod {
		target = filepath.Join(game.Base(), "assets/mods", pkgName)
	} else if typ == ccmodupdater.PackageTypeBase {
		if pkgName == "ccloader" {
			target = game.Base()
		} else {
			return fmt.Errorf("Unable to handle special behavior.")
		}
	} else {
		return fmt.Errorf("Unable to handle package type %s", mrp.data.Metadata.Type())
	}
	
	// -- It begins! --
	
	errors := []error{};
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
