package internal

import (
	"fmt"
	"os"
	"flag"
	
	"github.com/CCDirectLink/CCUpdaterCLI"
	"github.com/CCDirectLink/CCUpdaterCLI/local"
	"github.com/CCDirectLink/CCUpdaterCLI/remote"
)

//Context contains the context details for this command.
type Context struct {
	game *ccmodupdater.GameInstance
	upgraded *OnlineContext
}

//NewContext creates a new local context.
func NewContext(dir *string) (*Context, error) {
	if dir == nil {
		game := flag.Lookup("game")
		var dirVal string
		var err error
		if game != nil {
			dirVal = game.Value.String()
		} else {
			dirVal, err = os.Getwd()
			if err != nil {
				return nil, fmt.Errorf("Unable to get working directory (as game directory): %s", err)
			}
		}
		dir = &dirVal
	}
	game := ccmodupdater.NewGameInstance(*dir)
	plugins, err := local.AllLocalPackagePlugins(game)
	if err != nil {
		return nil, fmt.Errorf("Unable to prepare for checking local packages: %s", err)
	}
	game.LocalPlugins = plugins
	return &Context{
		game,
		nil,
	}, nil
}

//NewOnlineContext creates a new online context.
func NewOnlineContext(dir *string) (*OnlineContext, error) {
	ctx, err := NewContext(dir)
	if err != nil {
		return nil, err
	}
	rwc, err := ctx.Upgrade()
	if err != nil {
		return nil, err
	}
	return rwc, nil
}

func (ctx *Context) Game() *ccmodupdater.GameInstance {
	return ctx.game
}

//Upgrade upgrades the Context to an OnlineContext.
func (ctx *Context) Upgrade() (*OnlineContext, error) {
	if ctx.upgraded != nil {
		return ctx.upgraded, nil
	}
	packages, err := remote.GetRemotePackages()
	if err != nil {
		return nil, err
	}
	rwc := &OnlineContext{
		Context: *ctx,
		remote: packages,
	}
	ctx.upgraded = rwc
	rwc.upgraded = rwc
	return rwc, nil
}

// Execute executes a package transaction.
func (ctx *Context) Execute(tx ccmodupdater.PackageTX, stats *Stats) error {
	upgraded := ctx.upgraded
	var tc ccmodupdater.PackageTXContext
	if upgraded != nil {
		tc = ccmodupdater.PackageTXContext{
			LocalPackages: ctx.game.Packages(),
			RemotePackages: upgraded.remote,
		}
	} else {
		tc = ccmodupdater.PackageTXContext{
			LocalPackages: ctx.game.Packages(),
			RemotePackages: make(map[string]ccmodupdater.RemotePackage),
		}
	}
	forceFlag := flag.Lookup("force")
	force := false
	if forceFlag != nil {
		force = forceFlag.Value.String() != "false"
	}
	if !force {
		solutions, err := tc.Solve(tx)
		if err != nil {
			return err
		}
		if len(solutions) > 1 {
			return fmt.Errorf("Dependency issue; can solve this in multiple ways. (This shouldn't happen in the current system.) %v", solutions)
		}
		if len(solutions) == 0 {
			return fmt.Errorf("Internal error caused no solutions to be returned yet no error was returned.")
		}
		tx = solutions[0]
	}
	return tc.Perform(ctx.game, tx, func (pkg string, pre bool, remove bool, install bool) {
		if install && remove {
			if pre {
				fmt.Fprintf(os.Stderr, "updating %s\n", pkg)
			} else {
				stats.Updated++
			}
		} else if install {
			if pre {
				fmt.Fprintf(os.Stderr, "installing %s\n", pkg)
			} else {
				stats.Installed++
			}
		} else if remove {
			if pre {
				fmt.Fprintf(os.Stderr, "removing %s\n", pkg)
			} else {
				stats.Removed++
			}
		}
	})
}

//OnlineContext contains the details for an online context.
type OnlineContext struct {
	Context
	remote map[string]ccmodupdater.RemotePackage
}

//RemotePackages returns all the remote packages.
func (rwc *OnlineContext) RemotePackages() map[string]ccmodupdater.RemotePackage {
	target := map[string]ccmodupdater.RemotePackage{}
	for k, v := range rwc.remote {
		target[k] = v
	}
	return target
}

//Stats contains the statistics about the installed mods
type Stats struct {
	Installed int `json:"installed"`
	Updated   int `json:"updated"`
	Removed   int `json:"removed"`

	Warnings []string `json:"warnings,omitempty"`
}

//AddWarning to the statistics
func (stats *Stats) AddWarning(warning string) {
	stats.Warnings = append(stats.Warnings, warning)
}
