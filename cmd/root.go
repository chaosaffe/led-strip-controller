// Copyright © 2016 Jason Murray <jason@chaosaffe.io>
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
	"fmt"
	"log"
	"net/http"

	"goji.io"

	"github.com/chaosaffe/led-strip-controller/api"
	"github.com/chaosaffe/led-strip-controller/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var stripFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "led-strip-controller",
	Short: "Does stuff",
	Long:  `Does stuff`,
	Run: func(cmd *cobra.Command, args []string) {

		ss := config.BuildStrips(stripFile)
		defer func() {
			ss.AllOff()
		}()

		mux := goji.NewMux()
		api.Register(mux, ss)
		http.ListenAndServe(":8000", mux)

	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.led-strip-controller.yaml)")
	RootCmd.PersistentFlags().StringVar(&stripFile, "strips-config", "", "LED strip definition file path")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".led-strip-controller") // name of config file (without extension)
	viper.AddConfigPath("$HOME")                 // adding home directory as first search path
	viper.AutomaticEnv()                         // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
