package cmd

import (
	"fmt"
	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal/local"
)

//Uninstall removes a mod from a directory
func Uninstall(args []string) (*Stats, error) {
	game, err := local.GetGame()
	if err != nil {
		return nil, fmt.Errorf("cmd: Could not find game folder")
	}
	
	stats := &Stats{}
	for _, name := range args {
		mod, modExists := game.Packages()[name]
		if !modExists {
			stats.AddWarning(fmt.Sprintf("cmd: Could not remove mod '%s' because it couldn't be found", name))
			return stats, err
		}

		err = mod.Remove()
		if err != nil {
			stats.AddWarning(fmt.Sprintf("cmd: Could not remove mod '%s' because of an error in %s", name, err.Error()))
		}

		stats.Removed++
	}

	return stats, nil
}
