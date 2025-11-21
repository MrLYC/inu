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
	"io"
	"os"

	"github.com/mrlyc/inu/pkg/anonymizer"
	"github.com/mrlyc/inu/pkg/cli"
	"github.com/spf13/cobra"
)

var (
	anonymizeFile           string
	anonymizeContent        string
	anonymizeEntityTypes    []string
	anonymizeNoPrint        bool
	anonymizeOutput         string
	anonymizeOutputEntities string
)

// NewAnonymizeCmd creates the anonymize command.
func NewAnonymizeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "anonymize",
		Short: "Anonymize sensitive information in text",
		Long: `Anonymize text by detecting and replacing sensitive entities with placeholders.
The anonymized entities can be saved to a YAML file for later restoration.`,
		RunE: runAnonymize,
	}

	flags := cmd.Flags()
	flags.StringVarP(&anonymizeFile, "file", "f", "", "Read input from file")
	flags.StringVarP(&anonymizeContent, "content", "c", "", "Input content as string")
	flags.StringSliceVarP(&anonymizeEntityTypes, "entity-types", "t", anonymizer.DefaultEntityTypes, "Entity types to detect (comma-separated)")
	flags.BoolVar(&anonymizeNoPrint, "no-print", false, "Do not print output to stdout (default: print to stdout)")
	flags.StringVarP(&anonymizeOutput, "output", "o", "", "Write anonymized text to file")
	flags.StringVarP(&anonymizeOutputEntities, "output-entities", "e", "", "Write entities to YAML file")

	return cmd
}

func runAnonymize(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Check environment variables
	if err := cli.CheckRequiredEnvVars(); err != nil {
		return err
	}

	// Read input
	var stdin *os.File
	if anonymizeFile == "" && anonymizeContent == "" {
		stdin = os.Stdin
	}

	input, err := cli.ReadInput(anonymizeFile, anonymizeContent, stdin)
	if err != nil {
		return err
	}

	// Determine entity types
	entityTypes := anonymizeEntityTypes

	// Initialize LLM
	cli.ProgressMessage("=== Initializing LLM client... ===")
	llm, err := anonymizer.CreateOpenAIChatModel(ctx)
	if err != nil {
		return err
	}

	anon, err := anonymizer.NewHashHidePair(llm)
	if err != nil {
		return err
	}

	// Determine output writer
	var writer io.Writer
	var fileWriter *os.File

	if !anonymizeNoPrint && anonymizeOutput != "" {
		// Output to both stdout and file
		fileWriter, err = os.Create(anonymizeOutput)
		if err != nil {
			return err
		}
		defer fileWriter.Close()
		writer = io.MultiWriter(os.Stdout, fileWriter)
	} else if !anonymizeNoPrint {
		// Output to stdout only
		writer = os.Stdout
	} else if anonymizeOutput != "" {
		// Output to file only
		fileWriter, err = os.Create(anonymizeOutput)
		if err != nil {
			return err
		}
		defer fileWriter.Close()
		writer = fileWriter
	} else {
		// No output (no-print and no output file)
		writer = io.Discard
	}

	// Anonymize text with streaming
	cli.ProgressMessage("=== Anonymizing text... ===")
	entities, err := anon.AnonymizeTextStream(ctx, entityTypes, input, writer)
	if err != nil {
		return err
	}

	cli.ProgressMessage("=== Anonymization complete ===")
	// Output entities to stderr
	cli.WriteEntitiesToStderr(entities, anonymizeNoPrint)

	// Output entities to file
	if anonymizeOutputEntities != "" {
		if err := cli.SaveEntitiesToYAML(entities, anonymizeOutputEntities); err != nil {
			return err
		}
		cli.ProgressMessage("Entities saved to: %s", anonymizeOutputEntities)
	}

	cli.ProgressMessage("All done")
	return nil
}
