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
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jswidler/lockgit/pkg/log"
	"github.com/jswidler/lockgit/pkg/util"
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
}

type KeyLoadError struct {
	err error
}

func (err KeyLoadError) Error() string {
	return err.Error()
}

func IsKeyLoadError(err error) bool {
	if err != nil {
		switch err.(type) {
		case KeyLoadError:
			return true
		default:
			return false
		}
	}
	return false
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
	c.ManifestPath = filepath.Join(lockgitPath, "manifest")
	c.KeyPath = filepath.Join(lockgitPath, "key")
	key, err := ioutil.ReadFile(c.KeyPath)
	if os.IsNotExist(err) {
		return c, KeyLoadError{errors.Errorf("no key found at %s", c.KeyPath)}
	} else if err != nil {
		return c, KeyLoadError{errors.Wrapf(err, "error attempting to read key at %s", c.KeyPath)}
	} else if len(key) != 32 {
		return c, KeyLoadError{errors.Errorf("key in %s is the wrong size", c.KeyPath)}
	} else {
		c.Key = key
	}
	return c, nil
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

func (c Context) RelPath(absPath string) string {
	path, err := filepath.Rel(c.WorkingPath, absPath)
	if err != nil {
		return absPath
	}
	return path
}
