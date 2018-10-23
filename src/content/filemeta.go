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
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"os"
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
	return fmt.Sprintf("%s\t%s", base64.RawURLEncoding.EncodeToString(f.Sha), f.RelPath)
}

//// 4 byte salt + 20 byte sha1 hash = 24 bytes
//func hashData(data []byte) []byte {
//	salt := make([]byte, 4, 24)
//	_, err := rand.Read(salt);
//	log.FatalPanic(err)
//
//	h := sha1.New()
//	h.Write(salt)
//	h.Write(data)
//
//	sha := h.Sum(nil)
//	hash := salt[0:24]
//	copy(hash[4:], sha)
//
//	return hash
//}
//
func (f Filemeta) CompareFileToHash() (bool, error) {
	h := sha1.New()
	h.Write(f.Sha[:4])

	err := readIn(f.AbsPath, h);
	if err != nil {
		return false, err
	}

	sha := h.Sum(nil)
	return bytes.Equal(sha, f.Sha[4:]), nil
}

func readIn(path string, dst io.Writer) (err error) {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := io.Copy(dst, f); err != nil {
		return err
	}
	return nil
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