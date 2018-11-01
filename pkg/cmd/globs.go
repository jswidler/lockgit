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

package cmd

import (
	"fmt"

	"github.com/jswidler/lockgit/pkg/app"
	"github.com/spf13/cobra"
)

// globsCmd represents the globs command
var globsCmd = &cobra.Command{
	Use:   "globs",
	Short: "List all saved glob patterns",
	Long:  "List all saved glob patterns.\n\n" + globHelp,

	Aliases: []string{"glob"},

	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		files := app.LsGlobs(cliFlags())
		for _, f := range files {
			fmt.Println(f)
		}
	},
}

func init() {
	rootCmd.AddCommand(globsCmd)
}

const globHelp = `Glob patterns supported are like those found in a .gitignore file.  The following special terms are supported:

  *           matches any sequence of non-path-separators
  **          matches any sequence of characters, including path separator
  ?           matches any single non-path-separator character
  [class]     matches any single non-path-separator character against a class of characters (see below) 
  {alt1,...}  matches a sequence of characters if one of the comma-separated alternatives matches

Character classes support the following:

  [abc]       matches any single character within the set
  [a-z]       matches any single character in the range
  [^class]    matches any single character which does not match the class

Note: It is a good idea to surround glob patterns with quotes to prevent shell wildcard expansion.`
