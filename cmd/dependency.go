package cmd

import (
	"fmt"

	"github.com/CCDirectLink/CCUpdaterCLI/public"
	"github.com/Masterminds/semver"
)

func installDependencies(game *public.GameInstance, remote map[string]public.RemotePackage, mod public.LocalPackage, stats *Stats) error {
	for name, version := range mod.Dependencies() {
		if err := installDependency(game, remote, name, version, stats); err != nil {
			return err
		}
	}
	return nil
}

func installDependency(game *public.GameInstance, remote map[string]public.RemotePackage, name, version string, stats *Stats) error {	
	ver, err := semver.NewConstraint(version)
	if err != nil {
		stats.AddWarning(fmt.Sprintf("cmd: Dependency on mod '%s' had an invalid version number '%s'", name, version))
		return nil
	}

	localMod, localHasMod := game.Packages()[name]
	if localHasMod {
		if ver.Check(localMod.Metadata().Version) {
			// We don't need to do anything.
			return nil
		}
	}
	
	// --- From here on, we can just not bother to care about localMod ---

	newest, ok := remote[name]
	if !ok {
		stats.AddWarning(fmt.Sprintf("cmd: Mod '%s' not available: %s", name, err))
		return err
	}

	if !ver.Check(newest.Metadata().Version) {
		stats.AddWarning(fmt.Sprintf("cmd: Could not update mod '%s' to %s because the newest version is %s", name, version, newest.Metadata().Version))
		return nil
	}

	return installOrUpdateMod(game, remote, newest, stats)
}
