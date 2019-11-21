package local

import (
	"github.com/CCDirectLink/CCUpdaterCLI"
	"fmt"
)

// Makes a set of all the LocalPackagePlugins.
func AllLocalPackagePlugins(game *ccmodupdater.GameInstance) ([]ccmodupdater.LocalPackagePlugin, error) {
	// No CrossCode makes the whole thing invalid for saner debugging.
	ccp, err := NewCrossCodePackagePlugin(game)
	if err != nil {
		return nil, fmt.Errorf("Not CrossCode: %s", err.Error())
	}
	return []ccmodupdater.LocalPackagePlugin{
		NewCCLoaderPackagePlugin(game),
		ccp,
		NewModlikePackagePlugin(game, "assets/mods"),
	}, nil
}
