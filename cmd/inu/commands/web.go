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

package commands

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/mrlyc/inu/pkg/anonymizer"
	"github.com/mrlyc/inu/pkg/cli"
	"github.com/mrlyc/inu/pkg/web"
	"github.com/spf13/cobra"
)

var (
	webAddr        string
	webAdminUser   string
	webAdminToken  string
	webEntityTypes []string
)

// NewWebCmd creates the web command.
func NewWebCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "web",
		Short: "Start HTTP API server",
		Long: `Start a web server that provides RESTful API endpoints for text anonymization and restoration.

Authentication can be enabled by providing an admin token. If no token is provided, 
the server will run without authentication (not recommended for production).

Available endpoints:
  - GET  /              Web UI
  - GET  /health        Health check (no auth required)
  - GET  /api/v1/config Configuration
  - POST /api/v1/anonymize  Anonymize text
  - POST /api/v1/restore    Restore anonymized text`,
		RunE: runWeb,
	}

	cmd.Flags().StringVar(&webAddr, "addr", "127.0.0.1:8080", "Server address to listen on")
	cmd.Flags().StringVar(&webAdminUser, "admin-user", "admin", "Admin username for HTTP Basic Auth")
	cmd.Flags().StringVar(&webAdminToken, "admin-token", "", "Admin token/password for HTTP Basic Auth (leave empty to disable auth)")
	cmd.Flags().StringSliceVar(&webEntityTypes, "entity-types", anonymizer.DefaultEntityTypes, "Entity types to recognize")

	return cmd
}

func runWeb(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Check environment variables
	if err := cli.CheckRequiredEnvVars(); err != nil {
		return err
	}

	// Initialize LLM
	cli.ProgressMessage("Initializing LLM client...")
	llm, err := anonymizer.CreateOpenAIChatModel(ctx)
	if err != nil {
		return err
	}

	anon, err := anonymizer.NewHashHidePair(llm)
	if err != nil {
		return err
	}

	// Create web server
	config := &web.Config{
		Addr:       webAddr,
		AdminUser:  webAdminUser,
		AdminToken: webAdminToken,
	}

	server, err := web.NewServer(anon, config)
	if err != nil {
		return err
	}

	// Set entity types from command line flag
	server.SetEntityTypes(webEntityTypes)

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		if err := server.Start(); err != nil {
			errChan <- err
		}
	}()

	// Wait for shutdown signal or error
	select {
	case <-sigChan:
		// Graceful shutdown
		return server.Stop()
	case err := <-errChan:
		return err
	}
}
