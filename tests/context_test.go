package tests

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/jswidler/lockgit/pkg/app"
)

func TestMultipleLockgits(t *testing.T) {
	baseOpts := opts("multitest")
	setupVault(t, baseOpts)
	basefile := createFile(baseOpts, "basefile")
	group1Opts := opts(filepath.Join("multitest", "group1"))
	setupVault(t, group1Opts)
	group1file := createFile(group1Opts, "group1file")
	group2Opts := opts(filepath.Join("multitest", "group2"))
	setupVault(t, group2Opts)
	group2file := createFile(group2Opts, "group2file")

	err := app.AddToVault(baseOpts, basefile)
	if err != nil {
		t.Errorf("expected to add basefile to base vault: %s", err)
	}

	err = app.AddToVault(baseOpts, group1file)
	if err == nil {
		t.Error("expected to fail to add group1file to base vault")
	}

	err = app.AddToVault(group1Opts, group2file)
	if err == nil {
		t.Error("expected to fail to add group2file to group1 vault")
	}
}

func TestAddLockgitFile(t *testing.T) {
	opts := opts("lockgitfiltest")
	setupVault(t, opts)

	keypath := filepath.Join(opts.Wd, ".lockgit", "key")

	err := app.AddToVault(opts, []string{keypath})
	if err == nil {
		t.Error("expected to fail to add key file to vault")
	}
}

func createFile(opts app.Options, name string) []string {
	path := filepath.Join(opts.Wd, name)
	_ = ioutil.WriteFile(path, []byte("some data"), 0644)
	return []string{path}
}
