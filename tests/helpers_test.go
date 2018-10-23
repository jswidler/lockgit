package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jswidler/lockgit/src/app"
)

func setupVault(t *testing.T, params app.Options) {
	cleanDir(params.Wd)
	// Initialize the vault
	err := app.InitVault(params)
	if err != nil {
		t.Fatal("InitVault returned an error", err)
	}
}

func opts(dir string) app.Options {
	wd, _ := os.Getwd()
	return app.Options{Wd: filepath.Join(wd, "..", "build", "test", "vaults", dir)}
}

func cleanDir(path string) {
	_ = os.RemoveAll(path)
	_ = os.MkdirAll(path, 0755)
}
