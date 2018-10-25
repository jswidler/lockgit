package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jswidler/lockgit/pkg/app"
	"github.com/jswidler/lockgit/pkg/cmd"
	"github.com/spf13/viper"
)

func setupVault(t *testing.T, opts app.Options) {
	cleanDir(opts.Wd)

	reloadConfig(opts)

	// Initialize the vault
	err := app.InitVault(opts)
	if err != nil {
		t.Fatal("InitVault returned an error", err)
	}

	// Reload the config to get the key
	reloadConfig(opts)
}

func reloadConfig(opts app.Options) {
	viper.Reset()
	cmd.InitConfig(filepath.Join(opts.Wd, "config.yml"))
}

func opts(dir string) app.Options {
	wd, _ := os.Getwd()
	return app.Options{Wd: filepath.Join(wd, "..", "build", "test", "vaults", dir)}
}

func cleanDir(path string) {
	_ = os.RemoveAll(path)
	_ = os.MkdirAll(path, 0755)
}
