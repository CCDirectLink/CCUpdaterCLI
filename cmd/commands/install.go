package commands

import (
	"fmt"

	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal"
	"github.com/CCDirectLink/CCUpdaterCLI"
	"github.com/CCDirectLink/CCUpdaterCLI/remote"
)

var installed = 0

//Install a mod
func Install(context *internal.OnlineContext, args []string) (*internal.Stats, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("cmd: No mods installed since no mods were specified")
	}

	remote, err := remote.GetRemotePackages()
	if err != nil {
		return nil, fmt.Errorf("cmd: Could not download mod data because an error occured in %s", err.Error())
	}

	stats := &internal.Stats{}

	for _, name := range args {
		installedMods := context.Game().Packages()
		if _, modExists := installedMods[name]; modExists {
			stats.AddWarning(fmt.Sprintf("cmd: Could not install '%s' because it was already installed", name))
			continue
		}

		mod, hadMod := remote[name]
		if !hadMod {
			return stats, fmt.Errorf("cmd: Could not find '%s' available for download: %s", name, err.Error())
		}


		if err := installOrUpdateMod(context, mod, stats); err != nil {
			return stats, err
		}
	}

	_, hasCCLoader := context.Game().Packages()["ccloader"]
	if !hasCCLoader {
		stats.AddWarning("cmd: CCLoader wasn't installed (`ccmu install ccloader`); if mods are being installed, this may be required.")
	}
	
	return stats, nil
}

func installOrUpdateMod(context *internal.OnlineContext, mod ccmodupdater.RemotePackage, stats *internal.Stats) error {
	modName := mod.Metadata().Name()
	if err := installDependencies(context, mod, stats); err != nil {
		return fmt.Errorf("cmd: Could not install '%s' dependencies: %s", modName, err.Error())
	}

	localMod, hasMod := context.Game().Packages()[modName]
	if hasMod {
		if err := localMod.Remove(); err != nil {
			return fmt.Errorf("cmd: Could not remove '%s' for update because an error occured in %s", modName, err.Error())
		}
	}
	
	if err := mod.Install(context.Game()); err != nil {
		return fmt.Errorf("cmd: Could not install '%s' because an error occured in %s", modName, err.Error())
	}

	localMod, hasMod = context.Game().Packages()[modName]
	if !hasMod {
		stats.AddWarning(fmt.Sprintf("cmd: Installed '%s' but it seems to be an invalid mod", modName))
		return nil
	}

	stats.Installed++
	return nil
}
