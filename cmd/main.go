package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/CCDirectLink/CCUpdaterCLI/cmd/api"
	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal"
	"github.com/CCDirectLink/CCUpdaterCLI/cmd/commands"
)

func assertContext() *internal.Context {
	context, err := internal.NewContext(nil)
	if err != nil {
		fmt.Printf("UNABLE TO FIND GAME in %s\n", err.Error())
		os.Exit(1)
	}
	return context
}
func assertOnlineContext() *internal.OnlineContext {
	onlineContext, err := assertContext().Upgrade()
	if err != nil {
		fmt.Printf("UNABLE TO GO ONLINE in %s\n", err.Error())
		os.Exit(1)
	}
	return onlineContext
}

func main() {
	flag.String("game", "", "if set it overrides the path of the game")

	port := flag.Int("port", 9392, "the port which the api server listens on")
	host := flag.String("host", "localhost", "the host which the api server listens on")

	flag.Bool("v", false, "makes certain commands report more verbose output")
	flag.Bool("all", false, "for list: indicates all kinds of packages should be shown")
	flag.Bool("force", false, "for commands that perform actions: ignores automatic dependency handling")

	flag.Parse()

	if len(os.Args) == 1 {
		printHelp()
		return
	}

	op := flag.Arg(0)
	args := flag.Args()[1:]

	switch op {
	case "install",
		"i":
		printStatsAndError(commands.Install(assertOnlineContext(), args))
	case "remove",
		"delete",
		"uninstall":
		printStatsAndError(commands.Uninstall(assertContext(), args))
	case "update":
		printStatsAndError(commands.Update(assertOnlineContext(), args))
	case "list":
		commands.List(assertOnlineContext())
	case "outdated":
		commands.Outdated(assertOnlineContext())
	case "api":
		api.StartAt(*host, *port)
	case "version":
		printVersion()
	case "help":
		printHelp()
	default:
		fmt.Printf("%s\n is not a command", op)
		printHelp()
		os.Exit(1)
	}
}

func printStatsAndError(stats *internal.Stats, err error) {
	if stats != nil && stats.Warnings != nil {
		for _, warning := range stats.Warnings {
			fmt.Printf("Warning in %s\n", warning)
		}
	}

	if err != nil {
		fmt.Printf("ERROR in %s\n", err.Error())
	}

	if stats != nil {
		fmt.Printf("Installed %d, updated %d, removed %d\n", stats.Installed, stats.Updated, stats.Removed)
	}
}
