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

package context

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jswidler/lockgit/src/log"
	"github.com/jswidler/lockgit/src/util"
	"github.com/pkg/errors"
)

type Context struct {
	WorkingPath  string // working dir
	ProjectPath  string // path of the parent for .lockgit
	LockgitPath  string // path to .lockgit
	DataPath     string // path to .lockgit/data
	ManifestPath string // path to .lockgit/manifest
	KeyPath      string // path to .lockgit/key
	Key          []byte // key bytes loaded from .lockgit/key if present
	KeyLoaded    bool   // true if the key was loaded
}

// Return a Context provided a base to begin traversal from.
// The context will be from the first .lockgit directory found
func FromPath(path string, keyRequired bool) (Context, error) {
	c := Context{}

	pathabs, err := filepath.Abs(path)
	if err != nil {
		log.FatalPanic(errors.Wrap(err, "could not make absolute path"))
	}

	lockgitPath, err := FindLockgit(pathabs)
	if err != nil {
		return c, err
	}
	c.WorkingPath = pathabs
	c.LockgitPath = lockgitPath
	c.ProjectPath = filepath.Dir(lockgitPath)
	c.DataPath = filepath.Join(lockgitPath, "data")
	c.ManifestPath = filepath.Join(lockgitPath, "manifest")
	c.KeyPath = filepath.Join(lockgitPath, "key")
	key, err := getKey(c.KeyPath)
	c.Key = key
	c.KeyLoaded = key != nil
	if keyRequired {
		log.FatalExit(err)
	}
	return c, nil
}

// Find the .lockgit directory given a path.  If there is no .lockgit directory
// in the provided path, each parent directory will be searched untill one is found.
// Returns the path or an error if none is found.
func FindLockgit(path string) (string, error) {
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

func getKey(path string) ([]byte, error) {
	key, err := ioutil.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("no key found at %s", path)
	}
	log.FatalPanic(err)
	if len(key) != 32 {
		log.FatalPanic(fmt.Errorf("key in %s is the wrong size", path))
	}
	return key, nil
}
