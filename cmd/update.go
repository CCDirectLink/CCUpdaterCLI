package cmd

import (
	"fmt"

	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal/local"
	"github.com/CCDirectLink/CCUpdaterCLI/public"
)

//Update a mod
func Update(args []string) (*Stats, error) {
	game, err := local.GetGame()
	if err != nil {
		return nil, fmt.Errorf("cmd: Could not find game folder")
	}

	remote, err := public.GetRemotePackages()
	if err != nil {
		return nil, fmt.Errorf("cmd: Could not download mod data because an error occured in %s", err.Error())
	}

	if len(args) == 0 {
		return updateOutdated(game, remote)
	}

	stats := &Stats{}
	for _, name := range args {
		remoteMod, remoteModExists := remote[name]
		if !remoteModExists {
			stats.AddWarning(fmt.Sprintf("cmd: Couldn't update mod '%s' because no remote version exists.", name))
		} else {
			if err := installOrUpdateMod(game, remote, remoteMod, stats); err != nil {
				return stats, err
			}
		}
	}

	return stats, nil
}

func updateOutdated(game *public.GameInstance, remote map[string]public.RemotePackage) (*Stats, error) {
	stats := &Stats{}
	for modName, mod := range game.Packages() {
		remotePkg, hasRemote := remote[modName]
		if !hasRemote {
			continue
		}

		// remoteVer.Compare(localVer) is > 0 for outdated mods.
		if remotePkg.Metadata().Version.Compare(mod.Metadata().Version) <= 0 {
			continue
		}

		if err := installOrUpdateMod(game, remote, remotePkg, stats); err != nil {
			return stats, err
		}
	}

	return stats, nil
}

