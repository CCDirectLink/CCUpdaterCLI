package local

import (
	"flag"
	"os"
	"github.com/CCDirectLink/CCUpdaterCLI/public"
)

//GetGame using the current working directory or flags
func GetGame() (*public.GameInstance, error) {
	dir, err := getDir()
	if err != nil {
		return nil, err
	}
	return public.NewGameInstance(dir)
}

func getDir() (string, error) {
	game := flag.Lookup("game")
	if game != nil {
		return game.Value.String(), nil
	}

	return os.Getwd()
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if !os.IsNotExist(err) {
		return true, err
	}
	return false, nil
}
