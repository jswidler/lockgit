package tests

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/jswidler/lockgit/pkg/app"
)

func TestInitVaultSuccess(t *testing.T) {
	// Initialize the vault
	opts := opts("inittest")
	setupVault(t, opts)

	// Test there is a vault and a key
	testLockgitAndKeyExists(t, opts.Wd)

	// test the key has been added to gitignore
	gitignore, _ := os.Open(filepath.Join(opts.Wd, ".gitignore"))
	reader := bufio.NewReader(gitignore)
	for {
		readString, err := reader.ReadString('\n')
		if readString == ".lockgit/key" {
			// found it
			break
		} else if err != nil {
			t.Error(".lockgit/key not found in .gitignore")
		}
	}
}

func TestInitVaultFailsIfVaultExists(t *testing.T) {
	// Initialize the vault
	opts := opts("inittest")
	setupVault(t, opts)

	// Re-initialize the vault
	err := app.InitVault(opts)
	if err == nil {
		t.Error("expected vault init to fail, but it didn't")
	}
}

func TestInitVaultWithoutGitignore(t *testing.T) {
	// Initialize the vault
	opts := opts("inittest")
	opts.NoUpdateGitignore = true
	setupVault(t, opts)

	testLockgitAndKeyExists(t, opts.Wd)

	// gitignore should not exist
	_, err := os.Stat(filepath.Join(opts.Wd, ".gitignore"))
	if !os.IsNotExist(err) {
		t.Error("expected .gitignore to not exist, but it does")
	}
}

func testLockgitAndKeyExists(t *testing.T, testProjectPath string) {
	// Test there is a .lockgit folder created
	lockgitPath := filepath.Join(testProjectPath, ".lockgit")
	info, _ := os.Stat(lockgitPath)
	if info == nil {
		t.Fatal(".lockgit directory not found")
	} else if !info.IsDir() {
		t.Fatal(".lockgit exists but is not a directory")
	} else if info.Mode().Perm() != 0755 {
		t.Error(".lockgit exists but has wrong permissons")
	}

	// Test there is a keyfile
	keyPath := filepath.Join(lockgitPath, "key")
	info, _ = os.Stat(keyPath)
	if info == nil {
		t.Error(".lockgit/key not found")
	} else {
		if !info.Mode().IsRegular() {
			t.Error(".lockgit/key exists but is not a regular file")
		}
		if info.Mode().Perm() != 0644 {
			t.Error(".lockgit/key exists but has wrong permissons")
		}
		bytes, _ := ioutil.ReadFile(keyPath)
		if len(bytes) != 32 {
			t.Error(".lockgit/key is not 32 bytes")
		}
	}
}
