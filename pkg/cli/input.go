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

package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/rotisserie/eris"
)

// ReadInput reads input from file, content string, or stdin with priority: file > content > stdin.
func ReadInput(file, content string, stdin io.Reader) (string, error) {
	// Priority 1: file
	if file != "" {
		data, err := os.ReadFile(file)
		if err != nil {
			return "", eris.Wrapf(err, "failed to read file: %s", file)
		}
		return string(data), nil
	}

	// Priority 2: content
	if content != "" {
		return content, nil
	}

	// Priority 3: stdin
	if stdin == nil {
		return "", eris.New("no input provided: use --file, --content, or pipe to stdin")
	}

	data, err := io.ReadAll(stdin)
	if err != nil {
		return "", eris.Wrap(err, "failed to read from stdin")
	}

	if len(data) == 0 {
		return "", eris.New("no input provided: use --file, --content, or pipe to stdin")
	}

	return string(data), nil
}

// CheckRequiredEnvVars checks if required environment variables are set and returns a friendly error.
func CheckRequiredEnvVars() error {
	apiKey := os.Getenv("OPENAI_API_KEY")
	modelName := os.Getenv("OPENAI_MODEL_NAME")

	var missing []string
	if apiKey == "" {
		missing = append(missing, "OPENAI_API_KEY")
	}
	if modelName == "" {
		missing = append(missing, "OPENAI_MODEL_NAME")
	}

	if len(missing) > 0 {
		return eris.Errorf(`Required environment variables are not set: %v

Please configure them:
  export OPENAI_API_KEY="your-api-key"
  export OPENAI_MODEL_NAME="gpt-4"
  export OPENAI_BASE_URL="https://api.openai.com/v1"  # optional

For more information, see: https://github.com/MrLYC/inu#configuration`, missing)
	}

	return nil
}

// ProgressMessage prints a progress message to stderr.
func ProgressMessage(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}
