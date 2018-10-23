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
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	c "github.com/jswidler/lockgit/src/content"
	"github.com/jswidler/lockgit/src/context"
	"github.com/jswidler/lockgit/src/gitignore"
	"github.com/jswidler/lockgit/src/log"
	u "github.com/jswidler/lockgit/src/util"
	"github.com/pkg/errors"
)

func openFromVault(ctx context.Context, filemeta c.Filemeta, params Options) error {
	if !params.Force {
		matches, err := filemeta.CompareFileToHash()
		if matches {
			log.Verbose(fmt.Sprintf("skipping %s - file exists and matches hash", filemeta.RelPath))
			return nil
		} else if err == nil {
			log.Info(fmt.Sprintf("skipping %s - file exists but has changed.  To discard live version enable --force", filemeta.RelPath))
			return nil
		} else if err != nil && !os.IsNotExist(err) {
			log.LogError(err)
			return err
		}
	}

	// The file does not exist, or force is enabled:
	datafile, err := c.ReadDatafile(ctx, filemeta)
	if err != nil {
		return err
	}
	data, err := datafile.DecodeData()
	if err != nil {
		return err
	}
	_ = os.Mkdir(filepath.Dir(filemeta.AbsPath), 0755)
	err = ioutil.WriteFile(filemeta.AbsPath, data, os.FileMode(datafile.Perm))
	if err != nil {
		return err
	}
	log.Verbose(fmt.Sprintf("saved secret to %s", filemeta.AbsPath))
	return nil
}

func deletePlaintextFile(ctx context.Context, filemeta c.Filemeta, params Options) error {
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

		datafile.MatchesHash(filemeta.Sha)
		if !datafile.MatchesHash(filemeta.Sha) {
			return fmt.Errorf("%s has changed.  To delete anyway enable --force\n", filemeta.RelPath)
		}
	}
	err = os.Remove(filemeta.AbsPath)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("could not delete %s", filemeta.AbsPath))
	}
	return nil
}

func addFile(ctx context.Context, manifest *c.Manifest, absPath string, params Options) error {
	// the location of the file relative to project path
	relPath, err := filepath.Rel(ctx.ProjectPath, absPath)
	if err != nil {
		return err
	}
	info, err := os.Stat(absPath)
	if err != nil {
		return err
	} else if !info.Mode().IsRegular() {
		return fmt.Errorf("%s cannot be added because it is not a regular file", absPath)
	}
	relRoot := strings.Split(relPath, string(os.PathSeparator))[0]
	if relRoot == ".." {
		return fmt.Errorf("%s cannot be added because it is not in the project directory %s", absPath, ctx.ProjectPath)
	} else if relRoot == ".lockgit" {
		return fmt.Errorf("%s cannot be added because it is in the .lockgit directory", absPath)
	}

	mindx := manifest.Find(relPath)
	if !params.Force && mindx >= 0 {
		return fmt.Errorf("%s is already in the vault - enable --force or use commit to update instead", absPath)
	}

	filedata, err := ioutil.ReadFile(absPath)
	log.FatalPanic(err)
	datafile := c.Datafile{
		Ver:  1,
		Data: base64.RawStdEncoding.EncodeToString(filedata),
		Path: relPath,
		Perm: int(info.Mode().Perm()),
	}
	filemeta := c.NewFilemeta(absPath, datafile)

	datafile.Write(ctx, filemeta)

	if mindx >= 0 {
		oldDatafile := c.MakeDatafilePath(ctx, manifest.Files[mindx])
		_ = os.Remove(oldDatafile)
		manifest.Files[mindx] = filemeta
	} else {
		manifest.Add(filemeta)
	}

	if !params.NoUpdateGitignore {
		gitignore.Add(ctx.ProjectPath, filemeta.RelPath)
	}

	return nil
}

func deleteFileFromVault(ctx context.Context, manifest *c.Manifest, absPath string) error {
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

func ensureSameContext(ctx context.Context, files []string) error {
	for _, filename := range files {
		lockgit, err := context.FindLockgit(filename)
		if err != nil || ctx.LockgitPath != lockgit {
			// One of the files is in a different vault from the original
			fileRelWd := relativeIfPossible(ctx.WorkingPath, filename)
			altLockgitRelWd := relativeIfPossible(ctx.WorkingPath, lockgit)
			lockgitRelWd := relativeIfPossible(ctx.WorkingPath, ctx.LockgitPath)
			return fmt.Errorf("%s is in vault %s and not in the active vault %s", fileRelWd, altLockgitRelWd, lockgitRelWd)
		}
	}
	return nil
}

func relativeIfPossible(base, target string) string {
	path, err := filepath.Rel(base, target)
	if err != nil {
		return target
	}
	return path
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
		return nil, err
	}
	if len(k) != 32 {
		return nil, errors.New("key is not the correct length")
	}
	return k, nil
}
