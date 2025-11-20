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
	"os"

	"github.com/mrlyc/inu/pkg/anonymizer"
	"github.com/rotisserie/eris"
)

// WriteOutput writes content to stdout and/or file based on flags.
// By default (noPrint=false), content is written to stdout.
// If outputFile is specified, content is also written to the file.
// If noPrint=true, stdout output is suppressed.
func WriteOutput(content string, noPrint bool, outputFile string) error {
	// Write to stdout by default, unless noPrint is true
	if !noPrint {
		fmt.Println(content)
	}

	// Write to file if specified
	if outputFile != "" {
		if err := os.WriteFile(outputFile, []byte(content), 0644); err != nil {
			return eris.Wrapf(err, "failed to write output to file: %s", outputFile)
		}
	}

	return nil
}

// WriteEntitiesToStderr writes entity information to stderr in simplified format.
// By default (noPrint=false), entities are written to stderr.
// If noPrint=true, stderr output is suppressed.
func WriteEntitiesToStderr(entities []*anonymizer.Entity, noPrint bool) {
	if noPrint || len(entities) == 0 {
		return
	}

	for _, entity := range entities {
		values := ""
		if len(entity.Values) > 0 {
			values = entity.Values[0]
			for i := 1; i < len(entity.Values); i++ {
				values += ", " + entity.Values[i]
			}
		}
		fmt.Fprintf(os.Stderr, "%s: %s\n", entity.Key, values)
	}
}

// PrintEntitiesSimplified is deprecated. Use WriteEntitiesToStderr instead.
// Kept for backward compatibility during transition.
func PrintEntitiesSimplified(entities []*anonymizer.Entity) {
	WriteEntitiesToStderr(entities, false)
}
