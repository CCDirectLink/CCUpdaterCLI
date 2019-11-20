package cmd

import (
	"fmt"

	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal/local"
	"github.com/CCDirectLink/CCUpdaterCLI/public"
)

var installed = 0

//Install a mod
func Install(args []string) (*Stats, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("cmd: No mods installed since no mods were specified")
	}

	game, err := local.GetGame()
	if err != nil {
		return nil, fmt.Errorf("cmd: Could not find game folder")
	}

	remote, err := public.GetRemotePackages()
	if err != nil {
		return nil, fmt.Errorf("cmd: Could not download mod data because an error occured in %s", err.Error())
	}

	stats := &Stats{}

	for _, name := range args {
		installedMods := game.Packages()
		if _, modExists := installedMods[name]; modExists {
			stats.AddWarning(fmt.Sprintf("cmd: Could not install '%s' because it was already installed", name))
			continue
		}

		mod, hadMod := remote[name]
		if !hadMod {
			return stats, fmt.Errorf("cmd: Could not find '%s' available for download: %s", name, err.Error())
		}


		if err := installOrUpdateMod(game, remote, mod, stats); err != nil {
			return stats, err
		}
	}

	_, hasCCLoader := game.Packages()["ccloader"]
	if !hasCCLoader {
		stats.AddWarning("cmd: CCLoader wasn't installed (`ccmu install ccloader`); if mods are being installed, this may be required.")
	}
	
	return stats, nil
}

func installOrUpdateMod(game *public.GameInstance, remote map[string]public.RemotePackage, mod public.RemotePackage, stats *Stats) error {
	modName := mod.Metadata().Name

	localMod, hasMod := game.Packages()[modName]
	if hasMod {
		if err := localMod.Remove(); err != nil {
			return fmt.Errorf("cmd: Could not remove '%s' for update because an error occured in %s", modName, err.Error())
		}
	}
	
	if err := mod.Install(game); err != nil {
		return fmt.Errorf("cmd: Could not install '%s' because an error occured in %s", modName, err.Error())
	}

	localMod, hasMod = game.Packages()[modName]
	if !hasMod {
		stats.AddWarning(fmt.Sprintf("cmd: Installed '%s' but it seems to be an invalid mod", modName))
		return nil
	}

	stats.Installed++
	return installDependencies(game, remote, localMod, stats)
}
