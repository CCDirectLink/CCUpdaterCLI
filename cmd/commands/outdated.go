package commands

import (
	"fmt"

	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal"
)

//Outdated displays old mods and their new version
func Outdated(context *internal.OnlineContext) {
	outdated := false
	remote := context.RemotePackages()
	for modName, mod := range context.Game().Packages() {
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
