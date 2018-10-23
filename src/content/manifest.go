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
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jswidler/lockgit/src/context"
	"github.com/jswidler/lockgit/src/log"
)

type Manifest struct {
	Files []Filemeta
	path  string
}

func (m Manifest) Export() {
	err := ioutil.WriteFile(m.path, m.serialize(), 0644)
	log.FatalPanic(err)
}

func ImportManifest(ctx context.Context) Manifest {
	m := Manifest{
		Files: make([]Filemeta, 0, 32),
		path:  ctx.ManifestPath,
	}
	file, err := os.Open(ctx.ManifestPath)
	if os.IsNotExist(err) {
		return m
	}
	log.FatalExit(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tokens := strings.SplitN(scanner.Text(), "\t", 2)
		if len(tokens) != 2 {
			log.FatalExit(fmt.Errorf("%s does not have the expected format", ctx.ManifestPath))
		}

		sha, err := base64.RawURLEncoding.DecodeString(tokens[0])
		if err != nil {
			log.FatalExit(fmt.Errorf("%s does not have the expected format", ctx.ManifestPath))
		}
		m.Add(Filemeta{Sha: sha, RelPath: tokens[1], AbsPath: filepath.Join(ctx.ProjectPath, tokens[1])})
	}

	err = scanner.Err()
	log.FatalPanic(err)

	return m
}

func (m Manifest) Find(projRelPath string) int {
	i := sort.Search(len(m.Files), func(i int) bool { return m.Files[i].RelPath >= projRelPath })
	if i < len(m.Files) && m.Files[i].RelPath == projRelPath {
		return i
	}
	return -1
}

func (m *Manifest) Add(filemeta Filemeta) {
	m.Files = append(m.Files, filemeta)
	m.sort()
}

func (m Manifest) serialize() []byte {
	m.sort()
	var buffer bytes.Buffer
	for _, v := range m.Files {
		buffer.WriteString(v.String())
		buffer.WriteString("\n")
	}
	return buffer.Bytes()
}

func (m Manifest) sort() {
	sort.Slice(m.Files, func(i, j int) bool {
		return m.Files[i].RelPath < m.Files[j].RelPath
	})
}
