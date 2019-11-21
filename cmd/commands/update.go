package commands

import (
	"fmt"

	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal"
)

//Update a mod
func Update(context *internal.OnlineContext, args []string) (*internal.Stats, error) {
	if len(args) == 0 {
		return updateOutdated(context)
	}

	remotePackages := context.RemotePackages()
	stats := &internal.Stats{}
	for _, name := range args {
		remoteMod, remoteModExists := remotePackages[name]
		if !remoteModExists {
			stats.AddWarning(fmt.Sprintf("cmd: Couldn't update mod '%s' because no remote version exists.", name))
		} else {
			if err := installOrUpdateMod(context, remoteMod, stats); err != nil {
				return stats, err
			}
		}
	}

	return stats, nil
}

func updateOutdated(context *internal.OnlineContext) (*internal.Stats, error) {
	stats := &internal.Stats{}
	remotePackages := context.RemotePackages()
	for modName, mod := range context.Game().Packages() {
		remotePkg, hasRemote := remotePackages[modName]
		if !hasRemote {
			continue
		}

		// remoteVer.Compare(localVer) is > 0 for outdated mods.
		if remotePkg.Metadata().Version.Compare(mod.Metadata().Version) <= 0 {
			continue
		}

		if err := installOrUpdateMod(context, remotePkg, stats); err != nil {
			return stats, err
		}
	}

	return stats, nil
}

