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
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/jswidler/lockgit/pkg/content"
	"github.com/jswidler/lockgit/pkg/log"
	"github.com/jswidler/lockgit/pkg/util"
)

// Returns (headers, rows) for a table
func Status(opts Options) ([]string, [][]string) {
	ctx, manifest := loadcm(opts.Wd, loadcmopts{})
	headers := []string{"file", "updated", "pattern", "hash"}

	// Collect all the files which are tracked by patterns
	patternMatched := make([]string, 0, 64)
	if len(ctx.Config.Patterns) > 0 {
		for _, pattern := range ctx.Config.Patterns {
			absPattern := filepath.Join(ctx.ProjectPath, pattern)

			_, files, _, err := util.GetFiles(absPattern)
			log.FatalPanic(err)
			patternMatched = append(patternMatched, files...)
		}
	}

	if len(manifest.Files) == 0 && len(patternMatched) == 0 {
		log.Info("vault is empty")
		return nil, nil
	}

	table := make(statusTable, 0, 32)

	// iterate through files in the manifest
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

		patternMatched = util.Filter(patternMatched, func(path string) bool {
			return path != filemeta.AbsPath
		})

		table = append(table, []string{
			filemeta.RelPath,
			updated,
			firstMatchedPattern(datafile.Path, ctx.Config.Patterns),
			filemeta.ShaString(),
		})
	}

	// iterate through any files matched, but that were not seen in the manifest
	for _, notCommited := range patternMatched {
		table = append(table, []string{
			ctx.ProjRelPath(notCommited),
			"new file",
			firstMatchedPattern(ctx.ProjRelPath(notCommited), ctx.Config.Patterns),
			"",
		})
	}

	// sort the table and return it
	sort.Sort(table)
	return headers, table
}

type statusTable [][]string

func (s statusTable) Len() int {
	return len(s)
}

func (s statusTable) Less(i, j int) bool {
	return s[i][0] < s[j][0]
}

func (s statusTable) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
