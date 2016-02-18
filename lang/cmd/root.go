// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/kildevaeld/lang"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var service *lang.Service

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "lang",
	Short: "",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:

}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	configDir, err := lang.ConfigDir()

	if err != nil {
		printError(err)
	}

	if !filepath.IsAbs(configDir) {
		if configDir, err = filepath.Abs(configDir); err != nil {
			printError(err)
		}
	}

	service = lang.New(lang.Config{
		Root: configDir,
	})

	bs, _ := ioutil.ReadFile("../generator/manifest.json")

	var def map[string]lang.Definition
	json.Unmarshal(bs, &def)

	for _, v := range def {
		service.AddDefinition(v)
	}

	if err := RootCmd.Execute(); err != nil {
		printError(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.lang.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".lang") // name of config file (without extension)
	viper.AddConfigPath("$HOME") // adding home directory as first search path
	viper.AutomaticEnv()         // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
