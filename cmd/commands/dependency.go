package commands

import (
	"fmt"

	"github.com/CCDirectLink/CCUpdaterCLI"
	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal"
	"github.com/Masterminds/semver"
)

func installDependencies(context *internal.OnlineContext, mod ccmodupdater.LocalPackage, stats *internal.Stats) error {
	for name, version := range mod.Dependencies() {
		if err := installDependency(context, name, version, stats); err != nil {
			return err
		}
	}
	return nil
}

func installDependency(context *internal.OnlineContext, name, version string, stats *internal.Stats) error {	
	ver, err := semver.NewConstraint(version)
	if err != nil {
		stats.AddWarning(fmt.Sprintf("cmd: Dependency on mod '%s' had an invalid version number '%s'", name, version))
		return nil
	}

	localMod, localHasMod := context.Game().Packages()[name]
	if localHasMod {
		if ver.Check(localMod.Metadata().Version) {
			// We don't need to do anything.
			return nil
		}
	}
	
	// --- From here on, we can just not bother to care about localMod ---

	newest, ok := context.RemotePackages()[name]
	if !ok {
		stats.AddWarning(fmt.Sprintf("cmd: Mod '%s' not available: %s", name, err))
		return err
	}

	if !ver.Check(newest.Metadata().Version) {
		stats.AddWarning(fmt.Sprintf("cmd: Could not update mod '%s' to %s because the newest version is %s", name, version, newest.Metadata().Version))
		return nil
	}

	return installOrUpdateMod(context, newest, stats)
}
