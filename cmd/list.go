package cmd

import (
	"fmt"
	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal/local"
	"github.com/CCDirectLink/CCUpdaterCLI/public"
	"flag"
)

//List prints a list of all available mods
func List() {
	verbose := flag.Lookup("v").Value.String() != "false"
	all := flag.Lookup("all").Value.String() != "false"
	
	// Collect packages
	
	localPackages := map[string]public.LocalPackage{}
	
	game, err := local.GetGame()
	if err == nil {
		localPackages = game.Packages()
	}

	remotePackages, err := public.GetRemotePackages()
	if err != nil {
		fmt.Printf("cmd: Unable to continue, unable to get remote packages : %s\n", err.Error())
		return
	}
	
	// Collate packages
	packagesLatest := map[string]public.Package{}
	for k, v := range localPackages {
		packagesLatest[k] = v
	}
	for k, v := range remotePackages {
		oldPackage, oldExists := packagesLatest[k]
		if oldExists {
			oldPackageMeta := oldPackage.Metadata()
			if oldPackageMeta.Version.Compare(v.Metadata().Version) < 0 {
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
	for ptype := public.PackageTypeFirst; ptype < public.PackageTypeAfterLast; ptype++ {
		if all {
			if verbose {
				fmt.Printf("%s\n", ptype.StringPlural())
			}
		} else {
			if ptype != public.PackageTypeMod {
				continue
			}
		}
		for k, v := range packagesLatest {
			latestMeta := v.Metadata()
			local, localExists := localPackages[k]
			remote, remoteExists := remotePackages[k]
			localMeta := public.PackageMetadata{}
			if localExists {
				localMeta = local.Metadata()
			}
			remoteMeta := public.PackageMetadata{}
			if remoteExists {
				remoteMeta = remote.Metadata()
			}
			
			if latestMeta.Type == ptype {
				if verbose {
					fmt.Printf("%s%s %s:\n%s\t%s\n", prefix, k, latestMeta.Version, prefix, latestMeta.Description)
					status := "Not Installed"
					if localExists {
						status = "Installed"
					}
					if localExists && remoteExists {
						upToDateness := localMeta.Version.Compare(remoteMeta.Version)
						if upToDateness == 0 {
							status = "Up to date"
						} else if upToDateness == -1 {
							status = "Outdated (" + localMeta.Version.Original() + " installed)"
						} else if upToDateness == 1 {
							status = "Development Build"
						}
					}
					fmt.Printf("%s\t%s\n", prefix, status)
				} else {
					fmt.Printf("%s%s %s\n", prefix, latestMeta.Version.Original(), k)
				}
			}
		}
	}
}
