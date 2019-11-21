package ccmodupdater

// GameInstance represents an active instance.
type GameInstance struct {
	// Read-only: The base directory.
	base string
	// Read-write: Local plugins
	LocalPlugins []LocalPackagePlugin
}

// Base returns the base directory of the GameInstance.
func (gi GameInstance) Base() string {
	return gi.base
}

// NewGameInstance creates a new GameInstance with the given base directory. This has no LocalPlugins, so isn't fully usable by default.
func NewGameInstance(base string) *GameInstance {
	gi := &GameInstance{
		base: base,
		LocalPlugins: []LocalPackagePlugin{},
	}
	return gi
}

// Packages returns a map of the LocalPackages that are installed, where the keys are the .Metadata().Name values for those packages.
func (gi *GameInstance) Packages() map[string]LocalPackage {
	packages := map[string]LocalPackage{}
	for _, v := range gi.LocalPlugins {
		for _, pkg := range v.Packages() {
			packages[pkg.Metadata().Name] = pkg
		}
	}
	return packages
}
