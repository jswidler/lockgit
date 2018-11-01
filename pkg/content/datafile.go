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
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jswidler/lockgit/pkg/log"
	"github.com/pkg/errors"
)

type Datafile struct {
	Ver  int
	Data string
	Path string
	Perm int
}

func NewDatafile(ctx Context, absPath string) (Datafile, error) {
	d := Datafile{}
	relPath := ctx.ProjRelPath(absPath)
	info, err := os.Lstat(absPath)
	if err != nil {
		return d, err
	} else if !info.Mode().IsRegular() {
		return d, errors.Errorf("%s is not a regular file", ctx.ProjRelPath(absPath))
	}
	filedata, err := ioutil.ReadFile(absPath)
	if err != nil {
		return d, errors.Wrap(err, "unable to read")
	}
	datafile := Datafile{
		Ver:  1,
		Data: base64.RawStdEncoding.EncodeToString(filedata),
		Path: relPath,
		Perm: int(info.Mode().Perm()),
	}
	return datafile, nil
}

func (d Datafile) Serialize() []byte {
	jsondata, err := json.Marshal(d)
	log.FatalPanic(err)
	return jsondata
}

func (d Datafile) Write(ctx Context, filemeta Filemeta) {
	path := MakeDatafilePath(ctx, filemeta)
	ciphertext := encrypt(ctx.Key, compress(d.Serialize()))
	err := ioutil.WriteFile(path, ciphertext, 0644)
	log.FatalPanic(err)
}

func MakeDatafilePath(ctx Context, filemeta Filemeta) string {
	return filepath.Join(ctx.DataPath, base64.RawURLEncoding.EncodeToString(filemeta.Sha))
}

func (d Datafile) Hash() []byte {
	salt := make([]byte, 4, 24)
	_, err := rand.Read(salt)
	log.FatalPanic(err)

	h := sha1.New()
	h.Write(salt)
	h.Write(d.Serialize())

	sha := h.Sum(nil)
	hash := salt[0:24]
	copy(hash[4:], sha)
	return hash
}

func (d Datafile) MatchesHash(hash []byte) bool {
	h := sha1.New()
	h.Write(hash[:4])
	h.Write(d.Serialize())
	sha := h.Sum(nil)
	return bytes.Equal(sha, hash[4:])
}

func ReadDatafile(ctx Context, filemeta Filemeta) (Datafile, error) {
	data := Datafile{}
	ciphertext, err := ioutil.ReadFile(MakeDatafilePath(ctx, filemeta))
	if err != nil {
		return data, err
	}
	plaintext := decompress(decrypt(ctx.Key, ciphertext))
	err = json.Unmarshal(plaintext, &data)
	return data, err
}

func (d Datafile) DecodeData() ([]byte, error) {
	return base64.RawStdEncoding.DecodeString(d.Data)
}

func compress(data []byte) []byte {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(data)
	w.Close()
	return b.Bytes()
}

func decompress(data []byte) []byte {
	var out bytes.Buffer
	b := bytes.NewBuffer(data)
	r, _ := zlib.NewReader(b)
	io.Copy(&out, r)
	r.Close()
	return out.Bytes()
}

func encrypt(key, plaintext []byte) []byte {
	block, err := aes.NewCipher(key)
	log.FatalPanic(err)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	_, err = io.ReadFull(rand.Reader, iv)
	log.FatalPanic(err)
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	return ciphertext
}

func decrypt(key, ciphertext []byte) []byte {
	block, err := aes.NewCipher(key)
	log.FatalPanic(err)
	if len(ciphertext) < aes.BlockSize {
		panic("blocksize is too small")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)
	return ciphertext
}
