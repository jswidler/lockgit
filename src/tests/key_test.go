package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jswidler/lockgit/src/app"
)

const testKey = "X57J6W76UA6HXQ7VZ5KLLMMVYGVHTTGNQS3EJPR6AOQRP2WBZAUQ"

func TestSetKey(t *testing.T) {
	opts := opts("unlocktest")
	setupVault(t, opts)

	keyPath := filepath.Join(opts.Wd, ".lockgit", "key")
	_ = os.Remove(keyPath)

	err := app.SetKey(opts, testKey)
	if err != nil {
		t.Errorf("set key failed %s", err)
	}
	key := app.GetKey(opts)
	if key != testKey {
		t.Error("key was not successfully recalled")
	}
}

func TestSetKeyFailsIfKeyIsSet(t *testing.T) {
	opts := opts("unlocktest")
	setupVault(t, opts)

	err := app.SetKey(opts, testKey)
	if err == nil {
		t.Error("expected set key to fail and not overwrite key without force")
	}
}

func TestSetKeyOverwritesWithForce(t *testing.T) {
	opts := opts("unlocktest")
	setupVault(t, opts)

	// Check get key does not panic
	_ = app.GetKey(opts)

	opts.Force = true
	err := app.SetKey(opts, testKey)
	if err != nil {
		t.Errorf("set key failed %s", err)
	}
	key := app.GetKey(opts)
	if key != testKey {
		t.Error("key was not successfully recalled")
	}
}
