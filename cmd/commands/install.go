package commands

import (
	"fmt"

	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal"
	"github.com/CCDirectLink/CCUpdaterCLI"
	"github.com/CCDirectLink/CCUpdaterCLI/local"
)

var installed = 0

//Install a mod
func Install(context *internal.OnlineContext, args []string) (*internal.Stats, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("cmd: No mods installed since no mods were specified")
	}

	stats := &internal.Stats{}

	tx := make(ccmodupdater.PackageTX)
	installedMods := context.Game().Packages()
	for _, name := range args {
		if _, modExists := installedMods[name]; modExists {
			stats.AddWarning(fmt.Sprintf("cmd: Could not install '%s' because it was already installed", name))
			continue
		}
		tx[name] = ccmodupdater.PackageTXOperationInstall
	}
	
	err := context.Execute(tx, stats)
	for _, warning := range local.CheckLocal(context.Game(), context.RemotePackages()) {
		stats.AddWarning(warning)
	}
	return stats, err
}
