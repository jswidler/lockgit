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

package util

import (
	"os"

	"github.com/bmatcuk/doublestar"
)

type FileResult int

const (
	File      FileResult = 0
	Directory FileResult = 1
	Glob      FileResult = 2
)

func GetFiles(pattern string) (FileResult, []string, string, error) {
	stat, err := os.Lstat(pattern)
	if err == nil {
		if stat.IsDir() {
			// pattern is a specific directory
			newPattern := pattern + string(os.PathSeparator) + "**"
			matches, err := globNoDirectory(newPattern)
			return Glob, matches, newPattern, err
		} else {
			// pattern is a specific file
			return File, []string{pattern}, pattern, nil
		}
	} else if os.IsNotExist(err) {
		matches, err := globNoDirectory(pattern)
		return Glob, matches, pattern, err

	} else {
		return File, nil, pattern, err
	}
}

func globNoDirectory(pattern string) ([]string, error) {
	matches, err := doublestar.Glob(pattern)
	// filter out directories
	matches = Filter(matches, func(path string) bool {
		isDir, _ := ExistsDir(path)
		return !isDir
	})
	return matches, err
}

// Tests if a resource exists on the filesystem
func Exists(path string) (bool, error) {
	_, err := os.Lstat(path)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

// Tests if directory exists
func ExistsDir(path string) (bool, error) {
	file, err := os.Lstat(path)
	if err == nil {
		return file.Mode().IsDir(), nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

func Filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}
