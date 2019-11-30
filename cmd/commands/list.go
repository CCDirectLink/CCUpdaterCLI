package commands

import (
	"fmt"
	"flag"

	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal"
	"github.com/CCDirectLink/CCUpdaterCLI"
)

//List prints a list of all available mods
func List(context *internal.OnlineContext) {
	verbose := flag.Lookup("v").Value.String() != "false"
	all := flag.Lookup("all").Value.String() != "false"
	
	// Collect packages
	localPackages := context.Game().Packages()
	remotePackages := context.RemotePackages()
	
	// Collate packages
	packagesLatest := map[string]ccmodupdater.Package{}
	for k, v := range localPackages {
		packagesLatest[k] = v
	}
	for k, v := range remotePackages {
		oldPackage, oldExists := packagesLatest[k]
		if oldExists {
			oldPackageMeta := oldPackage.Metadata()
			if oldPackageMeta.Version().Compare(v.Metadata().Version()) < 0 {
				packagesLatest[k] = v
			}
		} else {
			packagesLatest[k] = v
		}
	}
	
	// Output
	prefix := ""
	if all && verbose {
		prefix = "\t"
	}
	for ptype := ccmodupdater.PackageTypeFirst; ptype < ccmodupdater.PackageTypeAfterLast; ptype++ {
		if all {
			if verbose {
				fmt.Printf("%s\n", ptype.StringPlural())
			}
		} else {
			if ptype != ccmodupdater.PackageTypeMod {
				continue
			}
		}
		for k, v := range packagesLatest {
			latestMeta := v.Metadata()
			local, localExists := localPackages[k]
			remote, remoteExists := remotePackages[k]
			localMeta := ccmodupdater.PackageMetadata{}
			if localExists {
				localMeta = local.Metadata()
			}
			remoteMeta := ccmodupdater.PackageMetadata{}
			if remoteExists {
				remoteMeta = remote.Metadata()
			}
			
			if latestMeta.Type() == ptype {
				if verbose {
					fmt.Printf("%s%s %s:\n%s\t%s\n", prefix, k, latestMeta.Version(), prefix, latestMeta.Description())
					status := "Not Installed"
					if localExists {
						status = "Installed"
					}
					if localExists && remoteExists {
						upToDateness := localMeta.Version().Compare(remoteMeta.Version())
						if upToDateness == 0 {
							status = "Up to date"
						} else if upToDateness == -1 {
							status = "Outdated (" + localMeta.Version().Original() + " installed)"
						} else if upToDateness == 1 {
							status = "Development Build"
						}
					}
					fmt.Printf("%s\t%s\n", prefix, status)
				} else {
					fmt.Printf("%s%s %s\n", prefix, latestMeta.Version().Original(), k)
				}
			}
		}
	}
}
