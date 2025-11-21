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
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/mrlyc/inu/pkg/anonymizer"
	"github.com/mrlyc/inu/pkg/cli"
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"
)

var (
	interactiveFile        string
	interactiveContent     string
	interactiveEntityTypes []string
	interactiveNoPrompt    bool
)

// NewInteractiveCmd creates the interactive command.
func NewInteractiveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "interactive",
		Short: "Interactive anonymization and restoration workflow",
		Long: `Interactively anonymize text, wait for user to process it externally,
then restore it using entities kept in memory. Supports multiple restoration cycles.

Typical workflow:
1. Command anonymizes your text and displays it
2. Copy the anonymized text
3. Process it externally (e.g., paste to ChatGPT for summarization)
4. Paste the processed text back
5. Press Ctrl+D (Unix) or Ctrl+Z (Windows) to restore
6. Repeat steps 3-5 as needed

The entities are kept in memory throughout the session, so you can process
the anonymized text multiple times with different external tools.`,
		RunE: runInteractive,
	}

	flags := cmd.Flags()
	flags.StringVarP(&interactiveFile, "file", "f", "", "Read input from file")
	flags.StringVarP(&interactiveContent, "content", "c", "", "Input content as string")
	flags.StringSliceVarP(&interactiveEntityTypes, "entity-types", "t", anonymizer.DefaultEntityTypes, "Entity types to detect (comma-separated)")
	flags.BoolVar(&interactiveNoPrompt, "no-prompt", false, "Disable detailed prompts (show minimal messages only)")

	return cmd
}

func runInteractive(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Check environment variables
	if err := cli.CheckRequiredEnvVars(); err != nil {
		return err
	}

	// Read original input
	var stdin *os.File
	if interactiveFile == "" && interactiveContent == "" {
		stdin = os.Stdin
	}

	input, err := cli.ReadInput(interactiveFile, interactiveContent, stdin)
	if err != nil {
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

	// Anonymize text with streaming output
	cli.ProgressMessage("Anonymizing text...")
	fmt.Fprintln(os.Stderr, "\n"+strings.Repeat("=", 60))
	fmt.Fprintln(os.Stderr, "ANONYMIZED TEXT:")
	fmt.Fprintln(os.Stderr, strings.Repeat("=", 60))
	entities, err := anon.AnonymizeTextStream(ctx, interactiveEntityTypes, input, os.Stdout)
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, strings.Repeat("=", 60))

	// Print prompt to stderr
	printPrompt(interactiveNoPrompt)

	// Interactive restoration loop
	for {
		processedText, err := readUntilEOF()
		if err != nil {
			return err
		}

		// Empty input means EOF without content, exit gracefully
		if processedText == "" {
			break
		}

		// Restore text using in-memory entities
		cli.ProgressMessage("Restoring text...")
		restoredText, err := anon.RestoreText(ctx, entities, processedText)
		if err != nil {
			// Best-effort restoration - output what we can
			fmt.Fprintln(os.Stderr, "Warning: Some placeholders could not be restored")
			restoredText = processedText
		}

		// Output restored text to stdout with clear separation
		fmt.Fprintln(os.Stderr, "\n"+strings.Repeat("=", 60))
		fmt.Fprintln(os.Stderr, "RESTORED TEXT:")
		fmt.Fprintln(os.Stderr, strings.Repeat("=", 60))
		fmt.Println(restoredText)
		fmt.Fprintln(os.Stderr, strings.Repeat("=", 60))

		// Show ready message for next input
		if interactiveNoPrompt {
			fmt.Fprintln(os.Stderr, "\nWaiting for input...")
		} else {
			fmt.Fprintln(os.Stderr, "\n📝 Ready for next input (Ctrl+D to restore, Ctrl+C to exit)")
		}
	}

	cli.ProgressMessage("Exiting")
	return nil
}

// printPrompt prints usage instructions to stderr
func printPrompt(noPrompt bool) {
	if noPrompt {
		fmt.Fprintln(os.Stderr, "\nWaiting for input...")
		return
	}

	fmt.Fprintln(os.Stderr, "\n"+strings.Repeat("-", 60))
	fmt.Fprintln(os.Stderr, "✅ Anonymization Complete")
	fmt.Fprintln(os.Stderr, strings.Repeat("-", 60))
	fmt.Fprintln(os.Stderr, "Next steps:")
	fmt.Fprintln(os.Stderr, "  1. Copy the anonymized text above")
	fmt.Fprintln(os.Stderr, "  2. Process it externally (e.g., paste to ChatGPT)")
	fmt.Fprintln(os.Stderr, "  3. Paste the processed text below")
	fmt.Fprintln(os.Stderr, "  4. Press Ctrl+D (Unix) or Ctrl+Z (Windows) to restore")
	fmt.Fprintln(os.Stderr, strings.Repeat("-", 60))
	fmt.Fprintln(os.Stderr, "\n📝 Paste your processed text here:")
}

// readUntilEOF reads from stdin until EOF (Ctrl+D) is encountered
func readUntilEOF() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return "", eris.Wrap(err, "failed to read input")
	}

	// EOF reached
	return strings.Join(lines, "\n"), nil
}
