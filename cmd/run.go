/*
Copyright © 2024 Steven A. Zaluk

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/stevezaluk/arcane-game-server/config"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stevezaluk/arcane-game-server/server"
)

var serv server.GameServer

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start the Game server",
	Long:  ``,
	PreRun: func(cmd *cobra.Command, args []string) {
		err := config.InitLogger()
		if err != nil {
			fmt.Println("Failed to initialize logger: ", err)
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[server - info] Starting game server...")
		serv.Start()
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("[server - info] Stopping game server...")
		serv.Stop()

		fileObject := viper.Get("log.fileObject").(*os.File)
		err := fileObject.Close()
		if err != nil {
			fmt.Println("Failed to close log file: ", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	err := config.ReadConfigFile(cfgFile)
	if err != nil {
		fmt.Println("Failed to read config file: ", err)
		os.Exit(1)
	}

	runCmd.Flags().IntP("port", "p", 8080, "Set the host port that the server should listen on")
	viper.BindPFlag("port", runCmd.Flags().Lookup("port"))

	runCmd.Flags().String("api.ip_address", "127.0.0.1", "The IP Address that the MTGJSON API is running on")
	viper.BindPFlag("api.ip_address", runCmd.Flags().Lookup("api.ip_address"))

	runCmd.Flags().Int("api.port", 8080, "The port that the MTGJSON API is running on")
	viper.BindPFlag("api.port", runCmd.Flags().Lookup("api.port"))

	runCmd.Flags().Bool("api.use_ssl", false, "Determine whether or not to use SSL when making API calls to MTGJSON API")
	viper.BindPFlag("api.use_ssl", runCmd.Flags().Lookup("api.use_ssl"))

	runCmd.Flags().String("api.email", "", "The email address that the game server should use for authenticating with the MTGJSON API")
	viper.BindPFlag("api.email", runCmd.Flags().Lookup("api.email"))

	runCmd.Flags().String("api.password", "", "The password that the game server should use for authenticating with the MTGJSON API")
	viper.BindPFlag("api.password", runCmd.Flags().Lookup("api.password"))

	runCmd.Flags().IntP("server.max_connections", "m", 4, "Set the max number of connections for the game server")
	viper.BindPFlag("server.max_connections", runCmd.Flags().Lookup("server.max_connections"))

	runCmd.Flags().String("log.path", "/var/log/arcane", "Set the directory where log files are saved")
	viper.BindPFlag("log.path", runCmd.Flags().Lookup("log.path"))
}
