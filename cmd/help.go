package main

import "fmt"

func printHelp() {
	fmt.Println("Usage: ccmu [options] [command] <args...>")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  --game <path>         Sets the game folder used for operations")
	fmt.Println("")
	fmt.Println("Some commands have command-specific options.")
	fmt.Println("These must still be placed in the [options] area.")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  install <name>        Installs one or more packages")
	fmt.Println("  uninstall <name>      Uninstall one or more packages")
	fmt.Println("  update [name...]      Updates one or more packages")
	fmt.Println("  list                  Lists local and remote mods.")
	fmt.Println("    --all : Shows all packages, not just mods.")
	fmt.Println("    -v : Shows more information about the mods/packages.")
	fmt.Println("  outdated              Show the names and versions of outdated packages")
	fmt.Println("  version               Display the version of this tool")
	fmt.Println("  help                  Display this message")
}
