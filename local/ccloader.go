package local

import (
	"fmt"
	"path/filepath"
	"os"
	"io/ioutil"
	"github.com/CCDirectLink/CCUpdaterCLI"
	"encoding/json"
)

type ccLoaderPackage struct {
	dir string
	metadata ccmodupdater.PackageMetadata
}

func (cc ccLoaderPackage) Metadata() ccmodupdater.PackageMetadata {
	return cc.metadata
}
func (cc ccLoaderPackage) Remove() error {
	if err := ioutil.WriteFile(filepath.Join(cc.dir, "package.json"), []byte("{\"name\": \"CrossCode\", \"version\" : \"1.2.3\", \"main\": \"assets/node-webkit.html\", \"chromium-args\" : \"--ignore-gpu-blacklist\", \"window\" : { \"toolbar\" : false, \"icon\" : \"favicon.png\", \"width\" : 1136, \"height\": 640, \"fullscreen\" : false }}"), os.ModePerm); err != nil {
		return fmt.Errorf("Couldn't replace package.json (Installation may be broken now!!!):", err.Error())
	}
	if err := os.RemoveAll(filepath.Join(cc.dir, "ccloader")); err != nil {
		return fmt.Errorf("Couldn't remove CCLoader: %s", err.Error())
	}
	// If someone messed with Simplify, this might not work, so put it after the CCLoader removal.
	// This'll ensure the goal is achieved even if the details are broken.
	if err := os.RemoveAll(filepath.Join(cc.dir, "assets/mods/simplify")); err != nil {
		return fmt.Errorf("Couldn't remove Simplify: %s", err.Error())
	}
	return nil
}

type ccloaderPackagePlugin struct {
	dir string
}

// NewCCLoaderPackagePlugin creates a LocalPackagePlugin given the game base.
func NewCCLoaderPackagePlugin(game *ccmodupdater.GameInstance) ccmodupdater.LocalPackagePlugin {
	return ccloaderPackagePlugin{
		dir: game.Base(),
	}
}

func (ccl ccloaderPackagePlugin) Packages() []ccmodupdater.LocalPackage {
	proof, err := os.Open(filepath.Join(ccl.dir, "ccloader/package.json"))
	if err != nil {
		return []ccmodupdater.LocalPackage{}
	}
	defer proof.Close()

	metadata := make(ccmodupdater.PackageMetadata)
	err = json.NewDecoder(proof).Decode(&metadata)
	if err != nil {
		return []ccmodupdater.LocalPackage{}
	}
	
	if err = metadata.Verify(); err != nil {
		return []ccmodupdater.LocalPackage{}
	}
	return []ccmodupdater.LocalPackage{
		ccLoaderPackage{ccl.dir, metadata},
	}
}
