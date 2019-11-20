package cmd

import (
	"fmt"

	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal/local"
	"github.com/CCDirectLink/CCUpdaterCLI/public"
)

//Outdated displays old mods and their new version
func Outdated() {
	game, err := local.GetGame()
	if err != nil {
		fmt.Printf("Could not find game folder. Make sure you executed the command inside the game folder.\n")
		return
	}

	remote, err := public.GetRemotePackages()
	if err != nil {
		fmt.Printf("Could not download mod data because an error occured in %s.\n", err.Error())
		return
	}

	outdated := false
	for modName, mod := range game.Packages() {
		remoteMod, remoteHasMod := remote[modName]
		thisModIsOutdated := false
		if remoteHasMod {
			thisModIsOutdated = remoteMod.Metadata().Version.Compare(mod.Metadata().Version) > 0
		}
		if thisModIsOutdated {
			if !outdated {
				outdated = true
				fmt.Println("New     Current Name")
			}

			fmt.Printf("%s   %s   %s\n", remoteMod.Metadata().Version, mod.Metadata().Version, modName)
		}
	}
}
