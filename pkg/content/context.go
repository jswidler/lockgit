// Copyright © 2018 Jesse Swidler <jswidler@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package content

import (
	"encoding/base32"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jswidler/lockgit/pkg/log"
	"github.com/jswidler/lockgit/pkg/util"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Context struct {
	WorkingPath string   // working dir
	ProjectPath string   // path of the parent for .lockgit
	LockgitPath string   // path to .lockgit
	DataPath    string   // path to .lockgit/data
	ConfigPath  string   // path to .lockgit/data/lgconfig
	Config      LgConfig // Config data

	Key []byte // key bytes loaded from lgconfig (if key is present)
}

// Return a Context provided a base to begin traversal from.
// The context will be from the first .lockgit directory found
func FromPath(path string) (Context, error) {
	c := Context{}

	pathabs, err := filepath.Abs(path)
	if err != nil {
		log.FatalPanic(errors.Wrap(err, "could not make absolute path"))
	}

	lockgitPath, err := findLockgit(pathabs)
	if err != nil {
		return c, err
	}
	c.WorkingPath = pathabs
	c.LockgitPath = lockgitPath
	c.ProjectPath = filepath.Dir(lockgitPath)
	c.DataPath = filepath.Join(lockgitPath, "data")
	c.ConfigPath = filepath.Join(c.LockgitPath, "lgconfig")

	c.Config, err = ReadConfig(c)
	if os.IsNotExist(err) {
		// TODO: Handle differently in V1 - error if file is missing?
		// v0.5 -> v0.6+:  For now, assume this file is missing because it was created with an old version of lockgit.
		// 				  Update the lockgit vault by creating a the config file and moving the key if it exists
		c.Config = NewLgConfig()
		c.Config.Write(c.ConfigPath)
		key, err := readKeyOldV05(c)
		if err == nil {
			keyStr := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(key)
			viper.Set("vaults."+c.Config.Id+".key", keyStr)
			viper.Set("vaults."+c.Config.Id+".path", c.ProjectPath)
			viper.WriteConfig()
		}
	} else if err != nil {
		log.FatalPanic(errors.Wrap(err, "could not read .lockgit/lgconfig"))
	}

	// Update path in config file if it has changed
	pathKey := "vaults." + c.Config.Id + ".path"
	if c.ProjectPath != viper.GetString(pathKey) {
		viper.Set(pathKey, c.ProjectPath)
		viper.WriteConfig()
	}

	c.Key, err = readKey(c)
	return c, err
}

func readKey(c Context) ([]byte, error) {
	keyStr := viper.GetString("vaults." + c.Config.Id + ".key")

	if keyStr == "" {
		return nil, &KeyLoadError{fmt.Sprintf("no key for %s found in %s", c.ProjectPath, viper.ConfigFileUsed())}
	}

	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(keyStr)
	if err != nil {
		return nil, &KeyLoadError{fmt.Sprintf("error attempting to read key for %s in %s: %s", c.ProjectPath, viper.ConfigFileUsed(), err.Error())}
	} else if len(key) != 32 {
		return key, &KeyLoadError{fmt.Sprintf("key for %s in %s is the wrong size", c.ProjectPath, viper.ConfigFileUsed())}
	}
	return key, nil
}

func readKeyOldV05(c Context) ([]byte, error) {
	keyPath := filepath.Join(c.LockgitPath, "key")
	key, err := ioutil.ReadFile(keyPath)
	if os.IsNotExist(err) {
		return key, &KeyLoadError{fmt.Sprintf("no key found at %s", keyPath)}
	} else if err != nil {
		return key, &KeyLoadError{fmt.Sprintf("error attempting to read key at %s: %s", keyPath, err.Error())}
	} else if len(key) != 32 {
		return key, &KeyLoadError{fmt.Sprintf("key in %s is the wrong size", keyPath)}
	}
	return key, nil
}

// Find the .lockgit directory given a path.  If there is no .lockgit directory
// in the provided path, each parent directory will be searched untill one is found.
// Returns the path or an error if none is found.
func findLockgit(path string) (string, error) {
	for {
		lockgitPath := filepath.Join(path, ".lockgit")
		if exist, _ := util.ExistsDir(lockgitPath); exist {
			return lockgitPath, nil
		}
		path = filepath.Dir(path)
		if path == "/" {
			return "", errors.New("no lockgit vault found")
		}
	}
}

func (c Context) ImportManifest() (Manifest, error) {
	return ImportManifest(c)
}

func (c Context) RelPath(absPath string) string {
	path, err := filepath.Rel(c.WorkingPath, absPath)
	if err != nil {
		return absPath
	}
	return path
}

func (c Context) ProjRelPath(absPath string) string {
	relPath, err := filepath.Rel(c.ProjectPath, absPath)
	log.FatalPanic(err)
	return relPath
}
