package tests

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/jswidler/lockgit/pkg/app"
)

func TestStatus(t *testing.T) {
	opts := opts("statustest")
	setupVault(t, opts)
	createFilesC(opts.Wd)

	app.AddToVault(opts, []string{"dir1/*"})
	addedFiles := app.Ls(opts)
	if len(addedFiles) != 2 {
		t.Errorf("expected 2 files in the vault")
	}

	_, table := app.Status(opts)
	if len(table) != 2 {
		t.Errorf("expected 2 files in the vault")
	}

	_ = ioutil.WriteFile(filepath.Join(opts.Wd, "dir1", "filea1"), []byte(data2), 0644)
	_ = ioutil.WriteFile(filepath.Join(opts.Wd, "dir1", "filea15"), []byte(data2), 0644)

	headers, table := app.Status(opts)
	if headers[1] != "updated" {
		t.Errorf("expected second column to be the updated column")
	}
	if table[0][0] != "dir1/filea1" {
		t.Errorf("expected first file to be filea1")
	}
	if table[0][1] != "true" {
		t.Errorf("expected first file to be updated")
	}
	if table[1][0] != "dir1/filea15" {
		t.Errorf("expected second file to be filea15")
	}
	if table[1][1] != "new file" {
		t.Errorf("expected second file to be new")
	}
	if table[2][0] != "dir1/fileb1" {
		t.Errorf("expected third file to be fileb1")
	}
	if table[2][1] != "false" {
		t.Errorf("expected first file to the same")
	}

	if table[0][3] == "" {
		t.Errorf("expected first file to have an id")
	}
	if table[1][3] != "" {
		t.Errorf("expected second file not to have an id")
	}

	app.AddToVault(opts, []string{"dir2"})
	_, table = app.Status(opts)
	if len(table) != 9 {
		t.Errorf("expected 9 files in the vault")
	}

	if table[0][2] != "dir1/*" {
		t.Errorf("expected first file to have pattern dir1/*")
	}
	if table[3][2] != "dir2/**" {
		t.Errorf("expected fourth file not to have pattern dir2/**")
	}
}
