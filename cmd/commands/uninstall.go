package commands

import (
	"fmt"

	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal"
	"github.com/CCDirectLink/CCUpdaterCLI/local"
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

	for _, warning := range local.CheckLocal(context.Game(), nil) {
		stats.AddWarning(warning)
	}

	return stats, nil
}
