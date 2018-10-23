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

	"github.com/jswidler/lockgit/src/log"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var noUpdateGitignore bool
var wd string

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

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.lockgit.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&noUpdateGitignore, "no-update-gitignore", "", false, "disable updating .gitignore file")

	viper.BindPFlag("noUpdateGitignore", rootCmd.PersistentFlags().Lookup("no-update-gitignore"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// Search config in home directory with name ".lockgit" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".lockgit")
	}

	viper.AutomaticEnv() // read in environment variables that match
	viper.AllSettings()
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {

		//fmt.Println("Using config file:", viper.ConfigFileUsed())

		// Not sure why the binding is not working...
		noUpdateGitignore = viper.GetBool("noUpdateGitignore")
	}
}

var force bool

func addForceFlag(cmd *cobra.Command, msg string) {
	cmd.Flags().BoolVarP(&force, "force", "f", false, msg)
}
