/*
 * Copyright 2024 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/mrlyc/inu/cmd/inu/commands"
)

var (
	// Version information (injected at build time)
	Version   = "dev"
	Commit    = "unknown"
	BuildTime = "unknown"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "inu",
		Short:   "Text anonymization and restoration tool",
		Long:    `Inu is a CLI tool for anonymizing sensitive information in text using LLM and restoring it back.`,
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", Version, Commit, BuildTime),
	}

	// Add subcommands
	rootCmd.AddCommand(commands.NewAnonymizeCmd())
	rootCmd.AddCommand(commands.NewRestoreCmd())
	rootCmd.AddCommand(commands.NewInteractiveCmd())
	rootCmd.AddCommand(commands.NewWebCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
