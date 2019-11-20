package public
import (
	"io/ioutil"
	"path/filepath"
)

// GameInstance represents an active CrossCode instance.
type GameInstance struct {
	// Read-only: The base directory of this CrossCode game instance.
	base string
	// Read-only: The CrossCode package instance
	game LocalPackage
}

// Base returns the base directory of the GameInstance.
func (gi GameInstance) Base() string {
	return gi.base
}

// NewGameInstance creates a new GameInstance with the given base directory.
func NewGameInstance(base string) (*GameInstance, error) {
	gi := &GameInstance{
		base: base,
	}
	game, err := gi.getCrossCodePackage()
	gi.game = game
	if err != nil {
		return nil, err
	}
	return gi, nil
}

// Packages returns a map of the LocalPackages that are installed, where the keys are the .Metadata().Name values for those packages.
func (gi *GameInstance) Packages() map[string]LocalPackage {
	// Start with CrossCode itself already in there
	packages := map[string]LocalPackage{
		"crosscode": gi.game,
	}
	lp, err := gi.getCCLoaderPackage()
	if err == nil {
		packages["ccloader"] = lp
	}
	// Mods (adapted from modfinder.go)
	modsDir := filepath.Join(gi.base, "assets/mods")
	dirs, err := ioutil.ReadDir(modsDir)
	if err == nil {
		for _, dir := range dirs {
			if dir.IsDir() {
				mod, err := getModPackage(filepath.Join(modsDir, dir.Name()))
				if err == nil {
					packages[mod.Metadata().Name] = mod
				}
			}
		}
	}
	return packages
}
