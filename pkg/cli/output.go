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
func WriteOutput(content string, print bool, outputFile string) error {
	if print {
		fmt.Println(content)
	}

	if outputFile != "" {
		if err := os.WriteFile(outputFile, []byte(content), 0644); err != nil {
			return eris.Wrapf(err, "failed to write output to file: %s", outputFile)
		}
	}

	// If neither print nor output is specified, just do nothing (no error)
	return nil
}

// PrintEntitiesSimplified prints entities in simplified format: key: values.
func PrintEntitiesSimplified(entities []*anonymizer.Entity) {
	for _, entity := range entities {
		values := ""
		if len(entity.Values) > 0 {
			values = entity.Values[0]
			for i := 1; i < len(entity.Values); i++ {
				values += ", " + entity.Values[i]
			}
		}
		fmt.Printf("%s: %s\n", entity.Key, values)
	}
}
