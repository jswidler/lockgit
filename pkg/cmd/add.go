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
	"github.com/jswidler/lockgit/pkg/app"
	"github.com/jswidler/lockgit/pkg/log"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add <file|glob> ...",
	Short: "Add files and glob patterns to the vault",
	Long: `Add files to the vault, either individually or with a glob pattern.

For each input, if it matches a file exactly, only that single will be added to the vault.  If the input matches a
directory exactly, all files of that directory and subdirectories will be added to the vault, and the glob pattern 
"<input>/**" will be added to the vault.  Otherwise, the input will be interpreted as a glob pattern.

Files and glob patterns added to the vault can be removed with the rm command.

` + globHelp,

	Example: `  Add a single file called credentials in the current directory:
  lockgit add credentials

  Add all .key files in the src directory and its subfolders:
  lockgit add 'src/**/*.key'`,

	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := app.AddToVault(cliFlags(), args)
		log.FatalExit(err)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addForceFlag(addCmd, "allow overwriting of existing files in the vault")
}
