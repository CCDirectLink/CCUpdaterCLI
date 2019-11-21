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
}

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
	}, nil
}
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
	packages, err := remote.GetRemotePackages()
	if err != nil {
		return nil, err
	}
	rwc := &OnlineContext{
		Context: *ctx,
		remote: packages,
	}
	return rwc, nil
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
