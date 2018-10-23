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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/jswidler/lockgit/src/content"
	"github.com/jswidler/lockgit/src/context"
	"github.com/jswidler/lockgit/src/log"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
)

type Params struct {
	NoUpdateGitignore		bool
	Force					bool
}

func InitVault(params Params) {
	projectPath, err := os.Getwd()
	log.FatalPanic(err)

	err = initVault(projectPath, params)
	log.FatalExit(err)
}


func SetKey(keystr string, force bool) {
	ctx, _ := loadcm(loadcmopts{ctxOnly:true})

	if !force && ctx.KeyLoaded {
		log.FatalExit(fmt.Errorf("key already exists, use --force to overwrite"))
	}
	key, err := keyToBytes(keystr)
	log.FatalExit(err)

	err = ioutil.WriteFile(ctx.KeyPath, key, 0644)
	log.FatalPanic(err)

	log.Info("key saved")
}

func RevealKey() {
	ctx, _ := loadcm(loadcmopts{ctxOnly:true, keyRequired:true})
	fmt.Println(keyToString(ctx.Key))
}

func Ls() {
	_, manifest := loadcm(loadcmopts{})
	for _, filemeta := range manifest.Files {
		fmt.Println(filemeta.RelPath)
	}
}

func AddToVault(files []string, params Params) {
	ctx, manifest := loadcm(loadcmopts{keyRequired:true, allowEmpty:true})

	// map inputs to absolute paths
	pathsToAbs(&files)

	err := ensureSameContext(ctx, files)
	log.FatalExit(errors.Wrap(err, "failed to add"))

	for _, filename := range files {
		err := addFile(ctx, &manifest, filename, params)
		if err != nil {
			log.LogError(err)
		} else {
			log.Info(fmt.Sprintf("added %s to vault", filename))
		}
	}
	manifest.Export()
}

func RemoveFromVault(files []string) {
	ctx, manifest := loadcm(loadcmopts{keyRequired:true})

	// map inputs to absolute paths
	pathsToAbs(&files)

	for _, filename := range files {
		err := deleteFileFromVault(ctx, &manifest, filename)
		if err != nil {
			log.LogError(err)
		} else {
			log.Info(fmt.Sprintf("removed %s from vault", filename))
		}
	}
	manifest.Export()
}

func Commit() {
	ctx, manifest := loadcm(loadcmopts{keyRequired:true})
	changes := false
	params := Params{Force:true}
	for _, filemeta := range manifest.Files {
		datafile, err := content.NewDatafile(ctx, filemeta.AbsPath)
		if err != nil && !os.IsNotExist(err) {
			log.LogError(err)
			continue
		}
		if !datafile.MatchesHash(filemeta.Sha) {
			err := addFile(ctx, &manifest, filemeta.AbsPath, params)
			if err != nil {
				log.LogError(err)
			} else {
				log.Info(fmt.Sprintf("%s updated", filemeta.AbsPath))
				changes = true
			}
		}
	}
	if !changes {
		log.Info("no changes")
	} else {
		manifest.Export()
	}
}

func Status() {
	ctx, manifest := loadcm(loadcmopts{filesRequired:true})

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"file", "updated", "hash"})
	table.SetBorder(false)

	for _, filemeta := range manifest.Files {
		var updated string
		datafile, err := content.NewDatafile(ctx, filemeta.AbsPath)
		if err != nil {
			if !os.IsNotExist(err) {
				log.LogError(err)
			}
			updated = "unavailable"
		} else {
			updated = strconv.FormatBool(!datafile.MatchesHash(filemeta.Sha))
		}

		table.Append([]string{filemeta.RelPath, updated, filemeta.ShaString()})
	}

	table.Render()
}

func OpenVault(params Params) {
	ctx, manifest := loadcm(loadcmopts{keyRequired:true, filesRequired:true})
	for _, filemeta := range manifest.Files {
		if err := openFromVault(ctx, filemeta, params); err != nil {
			log.LogError(errors.Wrapf(err, "error extracting secret"))
		}
	}
}

func CloseVault(params Params) {
	ctx, manifest := loadcm(loadcmopts{keyRequired:true, filesRequired:true})
	for _, filemeta := range manifest.Files {
		if err := deletePlaintextFile(ctx, filemeta, params); err != nil {
			log.LogError(err)
		}
	}
}

type loadcmopts struct {
	ctxOnly			bool
	keyRequired     bool
	filesRequired   bool
	allowEmpty      bool
}
func loadcm(opts loadcmopts) (context.Context, content.Manifest) {
	ctx, err := context.FromWd(opts.keyRequired)
	if err != nil {
		log.FatalExit(errors.Wrap(err, "unable to load lockgit"))
	} else if opts.ctxOnly {
		return ctx, content.Manifest{}
	}
	_ = os.Mkdir(ctx.DataPath, 0755)
	manifest := content.ImportManifest(ctx)
	if !opts.allowEmpty && len(manifest.Files) == 0 {
		fmt.Println("vault is empty")
		os.Exit(0)
	}
	return ctx, manifest
}

func pathsToAbs(files *[]string) {
	for i, file := range *files {
		path, err := filepath.Abs(file)
		if err != nil {
			log.FatalPanic(errors.Wrapf(err, "cannot make absolute path from %s", file))
		}
		(*files)[i] = path
	}
}