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

// deleteKeyCmd represents the delete-key command
var deleteKeyCmd = &cobra.Command{
	Use:   "delete-key",
	Short: "Delete the key for the current vault",
	Long: `Delete the key for the current vault.

You will not be able to recover the key after running this command.  It requires --force to work.`,

	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := app.UnsetKey(cliFlags())
		log.FatalExit(err)
	},
}

func init() {
	rootCmd.AddCommand(deleteKeyCmd)

	addForceFlag(deleteKeyCmd, "force is required for this operation to succeed")
}
