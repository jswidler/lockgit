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
	"os"

	"github.com/jswidler/lockgit/pkg/app"
	"github.com/jswidler/lockgit/pkg/log"
	"github.com/spf13/cobra"
)

// setKeyCmd represents the unlock command
var setKeyCmd = &cobra.Command{
	Use:   "set-key <KEY>",
	Short: "Set the key for the current vault",
	Long: `Set the key for the current vault.

The key for the vault can be displayed using the reveal-key command.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Fprint(os.Stderr, "invalid key\n")
			os.Exit(1)
		}
		err := app.SetKey(app.Options{Wd: wd, Force: force}, args[0])
		log.FatalExit(err)
	},
}

func init() {
	rootCmd.AddCommand(setKeyCmd)
	addForceFlag(setKeyCmd, "allow overwriting of an existing key")
}
