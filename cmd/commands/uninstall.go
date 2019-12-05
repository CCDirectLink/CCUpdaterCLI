package commands

import (
	"fmt"

	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal"
	"github.com/CCDirectLink/CCUpdaterCLI/local"
	"github.com/CCDirectLink/CCUpdaterCLI"
)

//Uninstall removes a mod from a directory
func Uninstall(context *internal.Context, args []string) (*internal.Stats, error) {
	stats := &internal.Stats{}

	tx := make(ccmodupdater.PackageTX)
	for _, name := range args {
		pkg := context.Game().Packages()[name]
		if pkg == nil {
			stats.AddWarning(fmt.Sprintf("cmd: Couldn't remove '%s' because it couldn't be found", name))
		} else {
			tx[name] = ccmodupdater.PackageTXOperationRemove
		}
	}
	
	err := context.Execute(tx, stats)
	for _, warning := range local.CheckLocal(context.Game(), nil) {
		stats.AddWarning(warning)
	}
	return stats, err
}
