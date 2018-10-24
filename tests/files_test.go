package tests

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/jswidler/lockgit/pkg/app"
)

const data1 = "this is some data"
const data2 = "this is different data"

func TestAddFile(t *testing.T) {
	opts := opts("addfiletest")
	setupVault(t, opts)

	files := createFilesA(opts.Wd)

	err := app.AddToVault(opts, files)
	if err != nil {
		t.Fatalf("failed to add files %s", err)
	}

	expected := []string{"filea", "foo/fileb"}
	ls := app.Ls(opts)
	if !reflect.DeepEqual(expected, ls) {
		t.Fatalf("ls returned %s instead of %s", ls, expected)
	}
}

func TestRemoveFile(t *testing.T) {
	opts := opts("removefiletest")
	setupVault(t, opts)

	files := createFilesA(opts.Wd)

	err := app.AddToVault(opts, files)
	if err != nil {
		t.Fatalf("failed to add files %s", err)
	}

	ls := app.Ls(opts)
	if len(ls) != 2 {
		t.Fatalf("expected ls to have 2 files, but it has %d", len(ls))
	}

	app.RemoveFromVault(opts, files[:1])
	ls = app.Ls(opts)
	if len(ls) != 1 {
		t.Fatalf("expected ls to have 1 files, but it has %d", len(ls))
	}
}

func TestRestoreFile(t *testing.T) {
	opts := opts("restoretest")
	setupVault(t, opts)

	filesA := createFilesA(opts.Wd)
	filesB := createFilesB(opts.Wd)
	err := app.AddToVault(opts, filesA)
	if err != nil {
		t.Fatalf("failed to add files %s", err)
	}

	app.CloseVault(opts)

	for _, f := range filesA {
		_, err := os.Stat(f)
		if !os.IsNotExist(err) {
			t.Fatalf("failed to delete %s", f)
		}
	}
	for _, f := range filesB {
		_, err := os.Stat(f)
		if err != nil {
			t.Errorf("expected to find %s", f)
		}
	}

	app.OpenVault(opts)

	bytes, err := ioutil.ReadFile(filesA[0])
	if err != nil {
		t.Errorf("unable to read %s", filesA[0])
	} else if string(bytes) != data1 {
		t.Errorf("%s has the wrong data", filesA[0])
	}
	info, _ := os.Stat(filesA[1])
	if info == nil || info.Mode().Perm() != 0600 {
		t.Errorf("%s has the wrong permissions", filesA[1])
	}
}

func TestUpdateFile(t *testing.T) {
	opts := opts("update")
	setupVault(t, opts)

	file := filepath.Join(opts.Wd, "filea")
	_ = ioutil.WriteFile(file, []byte(data1), 0644)
	err := app.AddToVault(opts, []string{file})
	if err != nil {
		t.Fatalf("failed to add test file to vault %s", err)
	}

	_ = ioutil.WriteFile(file, []byte(data2), 0644)

	err = app.AddToVault(opts, []string{file})
	if err == nil {
		t.Fatal("should have failed to add changed test file to vault")
	}

	err = app.Commit(opts)
	if err != nil {
		t.Fatalf("failed to commit change to vault %s", err)
	}

	app.CloseVault(opts)
	_, err = os.Stat(file)
	if !os.IsNotExist(err) {
		t.Errorf("failed to delete %s", file)
	}

	app.OpenVault(opts)

	bytes, _ := ioutil.ReadFile(file)
	if string(bytes) != data2 {
		t.Errorf("the file was not updated with the expected data")
	}
}

func createFilesA(projectdir string) []string {
	foodir := filepath.Join(projectdir, "foo")
	filea := filepath.Join(projectdir, "filea")
	fileb := filepath.Join(foodir, "fileb")

	_ = os.Mkdir(foodir, 0755)
	_ = ioutil.WriteFile(filea, []byte(data1), 0644)
	_ = ioutil.WriteFile(fileb, []byte(data2), 0600)

	return []string{filea, fileb}
}

func createFilesB(projectdir string) []string {
	bardir := filepath.Join(projectdir, "bar")
	filea := filepath.Join(projectdir, "filec")
	fileb := filepath.Join(bardir, "filed")

	_ = os.Mkdir(bardir, 0755)
	_ = ioutil.WriteFile(filea, []byte(data1), 0644)
	_ = ioutil.WriteFile(fileb, []byte(data2), 0600)

	return []string{filea, fileb}
}
