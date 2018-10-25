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
	"encoding/base64"
	"fmt"
)

// Meta data about a secret
type Filemeta struct {
	AbsPath string
	RelPath string
	Sha     []byte
}

func NewFilemeta(absPath string, datafile Datafile) Filemeta {
	return Filemeta{
		AbsPath: absPath,
		RelPath: datafile.Path,
		Sha:     datafile.Hash(),
	}
}

func (f Filemeta) ShaString() string {
	return base64.RawURLEncoding.EncodeToString(f.Sha)
}

func (f Filemeta) String() string {
	return fmt.Sprintf("%s\t%s", f.ShaString(), f.RelPath)
}

/* if you are worried about hash collisions for some reason...
func hashFileSecure(path string) ([]byte, error) {
	h := sha256.New()
	if err := readIn(path,h); err != nil {
		return nil, err
	}
	hash := h.Sum(nil)
	return bcrypt.GenerateFromPassword(hash, -1)
}

func CompareToHashSecure(path string, hash []byte) (bool, error) {
	h := sha256.New()
	if err := readIn(path,h); err != nil {
		return false, err
	}
	sha := h.Sum(nil)
	err := bcrypt.CompareHashAndPassword(hash, sha)
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	}
	return err == nil, err
}
*/
