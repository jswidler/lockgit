package tests

import (
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

	// TODO: Test there is a key in the config file

}
