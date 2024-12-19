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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stevezaluk/arcane-game-server/server"
)

var serv server.GameServer

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start the Game server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[server - info] Starting game server...")
		serv.Start()
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("[server - info] Stopping game server...")
		serv.Stop()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().IntP("port", "p", 8080, "Set the host port that the server should listen on")
	viper.BindPFlag("port", runCmd.Flags().Lookup("port"))

	runCmd.Flags().String("log-path", "/var/log/arcane", "Set the directory where log files are saved")
	viper.BindPFlag("log.path", runCmd.Flags().Lookup("log-path"))
}
