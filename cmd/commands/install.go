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

	for _, name := range args {
		installedMods := context.Game().Packages()
		if _, modExists := installedMods[name]; modExists {
			stats.AddWarning(fmt.Sprintf("cmd: Could not install '%s' because it was already installed", name))
			continue
		}

		mod, hadMod := context.RemotePackages()[name]
		if !hadMod {
			return stats, fmt.Errorf("cmd: Could not find '%s' available for download", name)
		}


		if err := installOrUpdateMod(context, mod, stats); err != nil {
			return stats, err
		}
	}

	for _, warning := range local.CheckLocal(context.Game(), context.RemotePackages()) {
		stats.AddWarning(warning)
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
