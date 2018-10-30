package tests

import (
	"testing"

	"github.com/jswidler/lockgit/pkg/app"
	"github.com/jswidler/lockgit/pkg/content"
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
		t.Errorf("set key failed :%s", err)
	}

	reloadConfig(opts)
	key := app.GetKey(opts)
	if key != testKey {
		t.Error("key was not successfully recalled")
	}
}

func TestUnsetKey(t *testing.T) {
	opts := opts("unsetkey")
	setupVault(t, opts)

	err := app.UnsetKey(opts)
	if err == nil {
		t.Error("expected unset-key to fail without force")
	}

	opts.Force = true
	err = app.UnsetKey(opts)
	if err != nil {
		t.Errorf("expected unset-key to work with force: %s", err)
	}

	err = app.UnsetKey(opts)
	if err == nil {
		t.Error("expected unset-key to fail when the key is unset")
	}

	opts.Force = false
	_, err = content.FromPath(opts.Wd)
	if !content.IsKeyLoadError(err) {
		t.Error("expected a KeyLoadError after unsetting the key")
	}

}
