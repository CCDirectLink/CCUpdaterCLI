package commands

import (
	"fmt"

	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal"
)

//Uninstall removes a mod from a directory
func Uninstall(context *internal.Context, args []string) (*internal.Stats, error) {
	stats := &internal.Stats{}
	for _, name := range args {
		mod, modExists := context.Game().Packages()[name]
		if !modExists {
			stats.AddWarning(fmt.Sprintf("cmd: Could not remove mod '%s' because it couldn't be found", name))
		} else {
			err := mod.Remove()
			if err != nil {
				stats.AddWarning(fmt.Sprintf("cmd: Could not remove mod '%s' because of an error in %s", name, err.Error()))
			}
		}

		stats.Removed++
	}
	for _, pkg := range context.Game().Packages() {
		for _, name := range args {
			_, hasDep := pkg.Metadata().Dependencies()[name]
			if hasDep {
				stats.AddWarning(fmt.Sprintf("cmd: '%s' was left without it's dependency '%s'", pkg.Metadata().Name(), name))
			}
		}
	}

	return stats, nil
}
