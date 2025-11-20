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

	"github.com/mrlyc/inu/pkg/anonymizer"
	"github.com/mrlyc/inu/pkg/cli"
	"github.com/spf13/cobra"
)

var (
	anonymizeFile           string
	anonymizeContent        string
	anonymizeEntityTypes    []string
	anonymizePrint          bool
	anonymizePrintEntities  bool
	anonymizeOutput         string
	anonymizeOutputEntities string
)

var defaultEntityTypes = []string{
	"个人信息", "业务信息", "资产信息", "账户信息",
	"位置数据", "文档名称", "组织机构", "岗位称谓",
}

// NewAnonymizeCmd creates the anonymize command.
func NewAnonymizeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "anonymize",
		Short: "Anonymize sensitive information in text",
		Long: `Anonymize text by detecting and replacing sensitive entities with placeholders.
The anonymized entities can be saved to a YAML file for later restoration.`,
		RunE: runAnonymize,
	}

	cmd.Flags().StringVarP(&anonymizeFile, "file", "f", "", "Read input from file")
	cmd.Flags().StringVarP(&anonymizeContent, "content", "c", "", "Input content as string")
	cmd.Flags().StringSliceVarP(&anonymizeEntityTypes, "entity-types", "t", nil, "Entity types to detect (comma-separated)")
	cmd.Flags().BoolVarP(&anonymizePrint, "print", "p", false, "Print anonymized text to stdout")
	cmd.Flags().BoolVar(&anonymizePrintEntities, "print-entities", false, "Print entities in simplified format to stdout")
	cmd.Flags().StringVarP(&anonymizeOutput, "output", "o", "", "Write anonymized text to file")
	cmd.Flags().StringVarP(&anonymizeOutputEntities, "output-entities", "e", "", "Write entities to YAML file")

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
	if len(entityTypes) == 0 {
		entityTypes = defaultEntityTypes
	}

	// Initialize LLM
	cli.ProgressMessage("Initializing LLM client...")
	llm, err := anonymizer.CreateOpenAIChatModel(ctx)
	if err != nil {
		return err
	}

	anon, err := anonymizer.New(llm)
	if err != nil {
		return err
	}

	// Anonymize text
	cli.ProgressMessage("Anonymizing text...")
	result, entities, err := anon.AnonymizeText(ctx, entityTypes, input)
	if err != nil {
		return err
	}

	// Output anonymized text
	if err := cli.WriteOutput(result, anonymizePrint, anonymizeOutput); err != nil {
		return err
	}

	// Output entities
	if anonymizePrintEntities {
		cli.PrintEntitiesSimplified(entities)
	}

	if anonymizeOutputEntities != "" {
		if err := cli.SaveEntitiesToYAML(entities, anonymizeOutputEntities); err != nil {
			return err
		}
		cli.ProgressMessage("Entities saved to: %s", anonymizeOutputEntities)
	}

	cli.ProgressMessage("Anonymization complete")
	return nil
}
