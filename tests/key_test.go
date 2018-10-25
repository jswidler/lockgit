package tests

import (
	"testing"

	"github.com/jswidler/lockgit/pkg/app"
)

const testKey = "X57J6W76UA6HXQ7VZ5KLLMMVYGVHTTGNQS3EJPR6AOQRP2WBZAUQ"

func TestSetKey(t *testing.T) {
	opts := opts("setkey")
	setupVault(t, opts)

	// Delete the original key
	reloadConfig(opts)
	opts.Force = true
	err := app.UnsetKey(opts)
	if err != nil {
		t.Errorf("unset key failed %s", err)
	}

	// Now try and set the new key
	opts.Force = false
	reloadConfig(opts)
	err = app.SetKey(opts, testKey)
	if err != nil {
		t.Errorf("set key failed %s", err)
	}

	reloadConfig(opts)
	key := app.GetKey(opts)
	if key != testKey {
		t.Error("key was not successfully recalled")
	}
}

func TestSetKeyFailsIfKeyIsSet(t *testing.T) {
	opts := opts("setkeyfail")
	setupVault(t, opts)

	err := app.SetKey(opts, testKey)
	if err == nil {
		t.Error("expected set key to fail and not overwrite key without force")
	}
}

func TestSetKeyOverwritesWithForce(t *testing.T) {
	opts := opts("setkeyforce")
	setupVault(t, opts)

	// Check get key does not panic
	_ = app.GetKey(opts)

	opts.Force = true
	err := app.SetKey(opts, testKey)
	if err != nil {
		t.Errorf("set key failed %s", err)
	}

	reloadConfig(opts)
	key := app.GetKey(opts)
	if key != testKey {
		t.Error("key was not successfully recalled")
	}
}
