package public

import (
	"fmt"
	"path/filepath"
	"os"
	"github.com/Masterminds/semver"
)

type ccLoaderPackage struct {
}

func (cc ccLoaderPackage) Metadata() PackageMetadata {
	return PackageMetadata{
		Name: "ccloader",
		Type: PackageTypeBase,
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

// Attempts to get the CCLoader package.
func (gi *GameInstance) getCCLoaderPackage() (LocalPackage, error) {
	proof, err := os.Open(filepath.Join(gi.base, "./ccloader/index.html"))
	if err != nil {
		return nil, fmt.Errorf("CCLoader not found")
	}
	proof.Close()
	return ccLoaderPackage{
	}, nil
}
