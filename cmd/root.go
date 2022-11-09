// MIT License
//
// Copyright (c) 2022 TJ Johnson <tj10200>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tj10200/docker-compose-orchestrator/pkg/orchestrator"
	"os"
)

var ComposeFiles []string
var DependencyConfig string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "orchestrator",
	Short: "run docker compose setup.",
	Long:  "Run your docker compose services with proper dependency control",

	Run: func(cmd *cobra.Command, args []string) {
		orch, err := orchestrator.NewOrchestrator()
		if err != nil {
			log.Fatalf("Failed to create orchestrator instance: %s", err)
		}

		err = orch.Run()
		if err != nil {
			log.Fatalf("Orchestrator failed to run: %s", err)
		}
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
	rootCmd.Flags().StringSliceVar(&ComposeFiles, "docker-compose-files", []string{}, "compose files to use for config settings")
	rootCmd.Flags().StringVar(&DependencyConfig, "config", "config.toml", "configuration file listing dependency graph for compose services")
}
