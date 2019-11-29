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

package app

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	c "github.com/jswidler/lockgit/pkg/content"
	"github.com/jswidler/lockgit/pkg/log"
	u "github.com/jswidler/lockgit/pkg/util"
	"github.com/pkg/errors"
)

func openFromVault(ctx c.Context, filemeta c.Filemeta, params Options) error {
	if !params.Force {
		datafile, err := c.NewDatafile(ctx, filemeta.AbsPath)
		if err == nil {
			// Able to read the file
			if datafile.MatchesCurrent(filemeta) {
				log.Verbose(fmt.Sprintf("skipping %s - file exists and is unchanged",
					ctx.RelPath(filemeta.AbsPath)))
				return nil
			} else {
				log.Info(fmt.Sprintf("skipping %s - file exists but has changed.  To discard live version enable --force",
					ctx.RelPath(filemeta.AbsPath)))
				return nil
			}
		} else if !os.IsNotExist(err) {
			// Not really sure what happened
			return err
		}
	}
	// if here, the file does not exist, or force is enabled
	datafile, err := c.ReadDatafile(ctx, filemeta)
	if err != nil {
		return err
	}
	data, err := datafile.DecodeData()
	if err != nil {
		return err
	}
	absPath := filepath.Join(ctx.ProjectPath, datafile.Path())
	_ = os.MkdirAll(filepath.Dir(absPath), 0755)
	err = ioutil.WriteFile(absPath, data, os.FileMode(datafile.Perm()))
	if err != nil {
		return err
	}
	log.Verbose(fmt.Sprintf("saved secret to %s", ctx.RelPath(absPath)))
	return nil
}

func deletePlaintextFile(ctx c.Context, filemeta c.Filemeta, params Options) error {
	exists, err := u.Exists(filemeta.AbsPath)
	if err != nil {
		return err
	} else if !exists {
		return nil
	}

	if !params.Force {
		datafile, err := c.NewDatafile(ctx, filemeta.AbsPath)
		if err != nil {
			return err
		}

		if !datafile.MatchesCurrent(filemeta) {
			return fmt.Errorf("%s has changed.  To delete anyway enable --force\n", ctx.RelPath(filemeta.AbsPath))
		}
	}
	err = os.Remove(filemeta.AbsPath)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("could not delete %s", ctx.RelPath(filemeta.AbsPath)))
	}
	return nil
}

func addFile(ctx c.Context, manifest *c.Manifest, absPath string, opts Options) error {
	relPath := ctx.ProjRelPath(absPath)
	relRoot := strings.Split(relPath, string(os.PathSeparator))[0]
	if relRoot == ".." {
		return fmt.Errorf("%s cannot be added because it is not in the project directory %s", ctx.RelPath(absPath), ctx.RelPath(ctx.ProjectPath))
	} else if relRoot == ".lockgit" {
		return fmt.Errorf("%s cannot be added because it is in the .lockgit directory", ctx.RelPath(absPath))
	}

	mindx := manifest.Find(relPath)
	if !opts.Force && mindx >= 0 {
		return fmt.Errorf("%s is already in the vault - enable --force or use commit to update instead", ctx.RelPath(absPath))
	}

	datafile, err := c.NewDatafile(ctx, absPath)
	if err != nil {
		return err
	}
	filemeta := c.NewFilemeta(absPath, datafile)

	datafile.Write(filemeta)

	if mindx >= 0 {
		oldDatafile := c.MakeDatafilePath(ctx, manifest.Files[mindx])
		_ = os.Remove(oldDatafile)
		manifest.Files[mindx] = filemeta
	} else {
		manifest.Add(filemeta)
	}

	return nil
}

func deleteFileFromVault(ctx c.Context, manifest *c.Manifest, absPath string) error {
	relPath, err := filepath.Rel(ctx.ProjectPath, absPath)
	if err != nil {
		return err
	}
	mindx := manifest.Find(relPath)
	if mindx < 0 {
		return fmt.Errorf("not found in manifest %s", relPath)
	}

	datafilepath := c.MakeDatafilePath(ctx, manifest.Files[mindx])
	_ = os.Remove(datafilepath)
	manifest.Files = append(manifest.Files[:mindx], manifest.Files[mindx+1:]...)
	return nil
}

func ensureSameContext(ctx c.Context, files []string) error {
	for _, filename := range files {
		fileCtx, _ := c.FromPath(filename)
		if ctx.LockgitPath != fileCtx.LockgitPath {
			// One of the files is in a different vault from the original
			if fileCtx.LockgitPath == "" {
				return fmt.Errorf("%s is not in the active vault %s",
					ctx.RelPath(filename), ctx.RelPath(ctx.LockgitPath))
			} else {
				return fmt.Errorf("%s is in vault %s and not in the active vault %s",
					ctx.RelPath(filename), ctx.RelPath(fileCtx.LockgitPath), ctx.RelPath(ctx.LockgitPath))
			}
		}
	}
	return nil
}

func genKey() []byte {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		panic(err)
	}
	return key
}

func keyToString(key []byte) string {
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(key)
}

func keyToBytes(key string) ([]byte, error) {
	k, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(key)
	if err != nil {
		return nil, errors.Wrap(err, "invalid key")
	}
	if len(k) != 32 {
		return nil, errors.New("invalid key: key is not the correct length")
	}
	return k, nil
}
