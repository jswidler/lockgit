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

// File is named zzz so init() is run last, and the command list is full and ready to be sorted

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jswidler/lockgit/pkg/app"
	"github.com/jswidler/lockgit/pkg/log"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var noUpdateGitignore bool
var wd string
var force bool

func cliFlags() app.Options {
	return app.Options{
		Wd:                wd,
		NoUpdateGitignore: noUpdateGitignore,
		Force:             force,
	}
}

var rootCmd = &cobra.Command{
	Use:   "lockgit",
	Short: "A secret vault for git repos",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		wd, err = os.Getwd()
		log.FatalPanic(err)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// This is the order the commands will be sorted in the help output
var cmdOrder = []string{
	"init",
	"set-key", "reveal-key", "delete-key",
	"add", "mv", "rm",
	"status", "commit",
	"open", "close",
	"ls", "globs",
}

func init() {
	cobra.OnInitialize(initConfig)

	cmds := cmdList(rootCmd.Commands())
	sort.Sort(cmds)

	cobra.EnableCommandSorting = false
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ~/.lockgit.yml)")
	rootCmd.PersistentFlags().BoolVarP(&noUpdateGitignore, "no-update-gitignore", "", false, "disable updating .gitignore file")

	viper.BindPFlag("no-update-gitignore", rootCmd.PersistentFlags().Lookup("no-update-gitignore"))
}

func initConfig() {
	InitConfig(cfgFile)
}

// Reads in config file
func InitConfig(file string) {
	if file != "" {
		// Use specified file
		viper.SetConfigFile(file)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		log.FatalExit(err)
		// Use the .lockgit.yml config file in the home dir
		viper.SetConfigFile(filepath.Join(home, ".lockgit.yml"))
	}

	// viper.AutomaticEnv() // read in environment variables that match
	viper.AllSettings()
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {

		//fmt.Println("Using config file:", viper.ConfigFileUsed())

		noUpdateGitignore = viper.GetBool("no-update-gitignore")
	}
}

// Return a validator for named positional arguments from the command-line
func cobraNamedPositionalArgs(argNames ...string) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) > len(argNames) {
			return errors.Errorf("too many arguments: got %d, expected %d", len(args), len(argNames))
		} else if len(args) < len(argNames) {
			var s string
			if len(args)+1 != len(argNames) {
				s = "s"
			}
			return errors.Errorf("missing required argument%s: %s", s, strings.Join(argNames[len(args):], ", "))
		} else {
			return nil
		}
	}
}

// Map each element to a negative number such the the first element in the list gets the number furthest from 0 and the
// last element of the list gets -1.
func mapit(arr []string) map[string]int {
	out := make(map[string]int)
	l := len(arr)
	for i, e := range arr {
		out[e] = i - l
	}
	return out
}

var orderMap = mapit(cmdOrder)

type cmdList []*cobra.Command

func (c cmdList) Len() int {
	return len(c)
}

func (c cmdList) Less(i, j int) bool {
	score := orderMap[c[i].Name()] - orderMap[c[j].Name()]
	if score == 0 {
		return c[i].Name() < c[j].Name()
	} else {
		return score < 0
	}
}

func (c cmdList) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func addForceFlag(cmd *cobra.Command, msg string) {
	cmd.Flags().BoolVarP(&force, "force", "f", false, msg)
}
