package commands

import (
	"fmt"

	"github.com/CCDirectLink/CCUpdaterCLI"
	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal"
	"github.com/Masterminds/semver"
)

func installDependencies(context *internal.OnlineContext, mod ccmodupdater.Package, stats *internal.Stats) error {
	for name, version := range mod.Metadata().Dependencies() {
		if err := installDependency(context, name, version, stats); err != nil {
			return err
		}
	}
	return nil
}

func installDependency(context *internal.OnlineContext, name string, ver *semver.Constraints, stats *internal.Stats) error {	
	localMod, localHasMod := context.Game().Packages()[name]
	if localHasMod {
		if ver.Check(localMod.Metadata().Version()) {
			// We don't need to do anything.
			return nil
		}
	}
	
	// --- From here on, we can just not bother to care about localMod ---

	newest, ok := context.RemotePackages()[name]
	if !ok {
		return fmt.Errorf("cmd: Dependency '%s' not available", name)
	}

	if !ver.Check(newest.Metadata().Version()) {
		stats.AddWarning(fmt.Sprintf("cmd: Could not update mod '%s' to '%s' because the available version is %s", name, ver, newest.Metadata().Version()))
		return nil
	}

	return installOrUpdateMod(context, newest, stats)
}
