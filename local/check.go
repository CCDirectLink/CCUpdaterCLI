package local

import (
	"fmt"
	"github.com/CCDirectLink/CCUpdaterCLI"
)

// CheckLocal returns a series of warnings. remote is optional, and allows additional information to be provided.
func CheckLocal(game *ccmodupdater.GameInstance, remote map[string]ccmodupdater.RemotePackage) []string {
	warnings := []string{}
	local := game.Packages()
	
	for pkgName, pkg := range local {
		for depName, constraint := range pkg.Metadata().Dependencies() {
			if depName == "Simplify" {
				// Dependencies on Simplify are dependencies on CCLoader
				continue
			}
			dep, hasDep := local[depName]
			if !hasDep {
				warnings = append(warnings, fmt.Sprintf("cmd: '%s' was left without it's dependency '%s'", pkgName, depName))
			} else {
				depVersion := dep.Metadata().Version()
				if !constraint.Check(depVersion) {
					warnings = append(warnings, fmt.Sprintf("cmd: '%s' has a dependency on '%s', but the version of that dependency (%s) doesn't match what's expected (%s)", pkgName, depName, depVersion, constraint))
				}
			}
		}
	}
	
	// CCLoader maintenance
	var remoteCCLoader ccmodupdater.RemotePackage
	remoteCCLoader, hasRemoteCCLoader := nil, false
	if remote != nil {
		remoteCCLoader, hasRemoteCCLoader = remote["ccloader"]
	}
	localCCLoader, hasCCLoader := local["ccloader"]
	if !hasCCLoader {
		warnings = append(warnings, "CCLoader wasn't installed; if mods are being installed, this may be required.")
	} else if hasRemoteCCLoader {
		if remoteCCLoader.Metadata().Version().GreaterThan(localCCLoader.Metadata().Version()) {
			warnings = append(warnings, "CCLoader is out of date; the mod may malfunction if it uses newer features.")
		}
	}
	
	return warnings
}
