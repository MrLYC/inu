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
	"os"

	"github.com/rotisserie/eris"
	"github.com/spf13/viper"

	"github.com/mrlyc/inu/pkg/anonymizer"
)

// EntitiesConfig represents the YAML structure for entities file.
type EntitiesConfig struct {
	Entities []*anonymizer.Entity `yaml:"entities" mapstructure:"entities"`
}

// SaveEntitiesToYAML saves entities to a YAML file.
func SaveEntitiesToYAML(entities []*anonymizer.Entity, file string) error {
	v := viper.New()
	v.Set("entities", entities)

	if err := v.WriteConfigAs(file); err != nil {
		return eris.Wrapf(err, "failed to write entities to YAML file: %s", file)
	}

	return nil
}

// LoadEntitiesFromYAML loads entities from a YAML file using viper.
func LoadEntitiesFromYAML(file string) ([]*anonymizer.Entity, error) {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return nil, eris.Errorf("entities file does not exist: %s", file)
	}

	v := viper.New()
	v.SetConfigFile(file)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, eris.Wrapf(err, "failed to read entities file: %s", file)
	}

	var config EntitiesConfig
	if err := v.Unmarshal(&config); err != nil {
		return nil, eris.Wrapf(err, "failed to parse entities from YAML file: %s", file)
	}

	// Empty entities list is valid
	return config.Entities, nil
}
