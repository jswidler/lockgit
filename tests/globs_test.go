package tests

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/jswidler/lockgit/pkg/app"
)

func TestAddGlob(t *testing.T) {
	opts := opts("globtest")
	setupVault(t, opts)
	createFilesC(opts.Wd)

	addedGlobs := app.LsGlobs(opts)
	if len(addedGlobs) != 0 {
		t.Errorf("expected 0 globs in the vault")
	}

	app.AddToVault(opts, []string{"**/fileb*"})
	addedFiles := app.Ls(opts)
	if len(addedFiles) != 7 {
		t.Errorf("expected 7 files in the vault")
	}
	addedGlobs = app.LsGlobs(opts)
	if len(addedGlobs) != 1 {
		t.Errorf("expected 1 globs in the vault")
	}

	// rerunning the same command should do nothing
	app.AddToVault(opts, []string{"**/fileb*"})
	addedFiles = app.Ls(opts)
	if len(addedFiles) != 7 {
		t.Errorf("expected 7 files in the vault")
	}
	addedGlobs = app.LsGlobs(opts)
	if len(addedGlobs) != 1 {
		t.Errorf("expected 1 globs in the vault")
	}

	// But a different command should add more files
	app.AddToVault(opts, []string{"**/filea*"})
	addedFiles = app.Ls(opts)
	if len(addedFiles) != 14 {
		t.Errorf("expected 14 files in the vault")
	}
	addedGlobs = app.LsGlobs(opts)
	if len(addedGlobs) != 2 {
		t.Errorf("expected 2 globs in the vault")
	}
}

func TestRemoveGlob(t *testing.T) {
	opts := opts("globtest2")
	setupVault(t, opts)
	createFilesC(opts.Wd)
	app.AddToVault(opts, []string{"**/filea*"})
	filecount := app.Ls(opts)
	if len(filecount) != 7 {
		t.Errorf("expected 7 files in the vault")
	}

	app.RemoveFromVault(opts, []string{"**/filea*"})
	filecount = app.Ls(opts)
	if len(filecount) != 0 {
		t.Errorf("expected 0 files in the vault")
	}
	globcount := app.LsGlobs(opts)
	if len(globcount) != 0 {
		t.Errorf("expected 0 globs in the vault")
	}
}

func TestDoNotAddGlobWithoutMatches(t *testing.T) {
	opts := opts("globtest3")
	setupVault(t, opts)
	createFilesC(opts.Wd)

	err := app.AddToVault(opts, []string{"**/filec*", "**/filea*"})
	if err == nil {
		t.Errorf("expected add to partially fail")
	}

	globs := app.LsGlobs(opts)
	if len(globs) != 1 || globs[0] != "**/filea*" {
		t.Errorf("expected only one glob to be saved")
	}
	files := app.Ls(opts)
	if len(files) != 7 {
		t.Errorf("expected to add 7 files")
	}
}

func createFilesC(projectdir string) {
	d1 := filepath.Join(projectdir, "dir1")
	d2 := filepath.Join(projectdir, "dir2")
	d11 := filepath.Join(d1, "dir11")
	d12 := filepath.Join(d1, "dir12")
	d21 := filepath.Join(d2, "dir21")
	d22 := filepath.Join(d2, "dir22")
	_ = os.Mkdir(d1, 0755)
	_ = os.Mkdir(d2, 0755)
	_ = os.Mkdir(d11, 0755)
	_ = os.Mkdir(d12, 0755)
	_ = os.Mkdir(d21, 0755)
	_ = os.Mkdir(d22, 0755)

	files := []string{
		filepath.Join(projectdir, "filea"),
		filepath.Join(d1, "filea1"),
		filepath.Join(d2, "filea2"),
		filepath.Join(d11, "filea11"),
		filepath.Join(d12, "filea12"),
		filepath.Join(d21, "filea21"),
		filepath.Join(d22, "filea22"),
		filepath.Join(projectdir, "fileb"),
		filepath.Join(d1, "fileb1"),
		filepath.Join(d2, "fileb2"),
		filepath.Join(d11, "fileb11"),
		filepath.Join(d12, "fileb12"),
		filepath.Join(d21, "fileb21"),
		filepath.Join(d22, "fileb22"),
	}

	for _, file := range files {
		_ = ioutil.WriteFile(file, []byte(data1), 0644)
	}
}
