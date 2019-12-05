package commands

import (
	"fmt"

	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal"
	"github.com/CCDirectLink/CCUpdaterCLI/local"
	"github.com/CCDirectLink/CCUpdaterCLI"
)

//Update a mod
func Update(context *internal.OnlineContext, args []string) (*internal.Stats, error) {
	if len(args) == 0 {
		args = updateOutdated(context)
	}

	localPackages := context.Game().Packages()
	stats := &internal.Stats{}

	tx := make(ccmodupdater.PackageTX)
	for _, name := range args {
		if localPackages[name] == nil {
			stats.AddWarning(fmt.Sprintf("cmd: Couldn't update mod '%s' because it isn't installed.", name))
		} else {
			tx[name] = ccmodupdater.PackageTXOperationInstall
		}
	}
	
	err := context.Execute(tx, stats)
	for _, warning := range local.CheckLocal(context.Game(), context.RemotePackages()) {
		stats.AddWarning(warning)
	}
	return stats, err
}

func updateOutdated(context *internal.OnlineContext) []string {
	args := []string{}
	remotePackages := context.RemotePackages()
	for modName, mod := range context.Game().Packages() {
		remotePkg, hasRemote := remotePackages[modName]
		if !hasRemote {
			continue
		}

		// remoteVer.Compare(localVer) is > 0 for outdated mods.
		if remotePkg.Metadata().Version().Compare(mod.Metadata().Version()) <= 0 {
			continue
		}

		args = append(args, modName)
	}

	return args
}

