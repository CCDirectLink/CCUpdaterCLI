package ccmodupdater
import (
	"fmt"
	"github.com/Masterminds/semver"
)

// PackageType represents a specific kind of package.
type PackageType int

// PackageTypeFirst is the first PackageType.
const PackageTypeFirst PackageType = 0

const (
	// PackageTypeBase represents a hidden base package that should not be visible from UI (but may be visible in CLI).
	PackageTypeBase PackageType = iota
	// PackageTypeMod represents a mod.
	PackageTypeMod
	// PackageTypeTool represents a tool.
	PackageTypeTool
	// PackageTypeAfterLast is above the last PackageType, to allow iteration in for loops. If it appears, it means "unknown".
	PackageTypeAfterLast
)

func (pt PackageType) String() string {
	switch pt {
		case PackageTypeBase:
			return "Base"
		case PackageTypeMod:
			return "Mod"
		case PackageTypeTool:
			return "Tool"
	}
	return string(pt)
}

// StringPlural is similar to String, but the resulting string is a plural
func (pt PackageType) StringPlural() string {
	if pt == PackageTypeBase {
		return "Base Packages"
	}
	return pt.String() + "s"
}

// PackageMetadata represents the metadata for a package. For further details, please see the inputLocationsFormat.md file.
type PackageMetadata map[string]interface{}

// Verify verifies that the PackageMetadata is valid and won't panic.
func (pm PackageMetadata) Verify() error {
	name := pm["name"]
	switch name.(type) {
		case string:
		default:
			return fmt.Errorf("'name' not string")
	}
	version := pm["version"]
	switch ver := version.(type) {
		case nil:
		case string:
			_, err := semver.NewVersion(ver)
			if err != nil {
				return fmt.Errorf("'version' invalid: %s", err.Error())
			}
		default:
			return fmt.Errorf("'version' invalid: not string")
	}
	tp := pm["ccmodType"]
	switch pmType := tp.(type) {
		case nil:
			// Not present, not a problem
		case string:
			if pm.Type() == PackageTypeAfterLast {
				return fmt.Errorf("'ccmodType' invalid: %s not recognized type", pmType)
			}
		default:
			return fmt.Errorf("'ccmodType' invalid: not string")
	}
	dependencies := pm["ccmodDependencies"]
	if dependencies == nil {
		dependencies = pm["dependencies"]
	}
	switch depMap := dependencies.(type) {
		case nil:
			// Not present, not a problem
		case map[string]interface{}:
			for k, v := range depMap {
				switch ver := v.(type) {
					case string:
						_, err := semver.NewConstraint(ver)
						if err != nil {
							return fmt.Errorf("dependency %s is invalid: %s", k, err.Error())
						}
					default:
						return fmt.Errorf("dependency %s is invalid: not string", k)
				}
			}
		default:
			return fmt.Errorf("'ccmodDependencies' is invalid: not object")
	}
	// Un-functional stuff
	de := pm["description"]
	switch de.(type) {
		case nil:
			// Not present, not a problem
		case string:
		default:
			return fmt.Errorf("'description' invalid: not string")
	}
	return nil
}

// Name gets the name of the package.
func (pm PackageMetadata) Name() string {
	return pm["name"].(string)
}

// Type gets the type of the package.
func (pm PackageMetadata) Type() PackageType {
	tp := pm["ccmodType"]
	if tp == nil {
		return PackageTypeMod
	}
	tpString := tp.(string)
	if tpString == "mod" {
		return PackageTypeMod
	} else if tpString == "base" {
		return PackageTypeBase
	} else if tpString == "tool" {
		return PackageTypeTool
	}
	// Used as an "unknown" value. Shouldn't get here if called from general code.
	// Used by Verify() to make sure it doesn't.
	return PackageTypeAfterLast
}

// Version gets the version of the package.
func (pm PackageMetadata) Version() *semver.Version {
	ver := pm["version"]
	if ver == nil {
		return semver.MustParse("0.0.0");
	}
	return semver.MustParse(pm["version"].(string))
}

// Dependencies gets the dependencies of the package.
func (pm PackageMetadata) Dependencies() map[string]*semver.Constraints {
	deps := make(map[string]*semver.Constraints)
	depsBaseUnk := pm["ccmodDependencies"]
	if depsBaseUnk == nil {
		depsBaseUnk := pm["dependencies"]
		if depsBaseUnk == nil {
			return deps
		}
	}
	depsBase := depsBaseUnk.(map[string]interface{})
	for k, v := range depsBase {
		cn, err := semver.NewConstraint(v.(string))
		if err != nil {
			panic(err)
		}
		deps[k] = cn
	}
	return deps
}

// Description gets the description of the package.
func (pm PackageMetadata) Description() string {
	description := pm["description"]
	if description == nil {
		return ""
	}
	return description.(string)
}


// Package represents a package, local or remote.
type Package interface {
	// Returns the metadata for this package.
	Metadata() PackageMetadata
}

// RemotePackage represents a package on a remote server.
type RemotePackage interface {
	Package
	// Installs the package. This does not check dependencies. This may affect other packages.
	Install(target *GameInstance) error
}

// LocalPackage represents a package installed in a GameInstance. It is invalidated when the GameInstance is modified.
type LocalPackage interface {
	Package
	// Attempts to remove the package. This does not check dependencies. This may affect other packages.
	Remove() error
}

// LocalPackagePlugin represents a plugin to the system that adds a scanner for local packages. It is assumed that the *GameInstance is supplied upon construction. The system is built this way because local packages are 
type LocalPackagePlugin interface {
	Packages() []LocalPackage
}
