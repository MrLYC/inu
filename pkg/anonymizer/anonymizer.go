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

package anonymizer

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/rotisserie/eris"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

// Anonymizer is the entity anonymization handler.
type Anonymizer struct {
	anonymizeTemplate *prompt.DefaultChatTemplate
	llm               model.BaseChatModel
}

// createAnonymizeMessages creates messages for anonymization.
func (a *Anonymizer) createAnonymizeMessages(ctx context.Context, types []string, text string) ([]*schema.Message, error) {
	encodedTypes, err := json.Marshal(types)
	if err != nil {
		return nil, eris.Wrap(err, "failed to marshal types")
	}

	messages, err := a.anonymizeTemplate.Format(ctx, map[string]any{
		"types": string(encodedTypes),
		"text":  text,
	})

	if err != nil {
		return nil, eris.Wrap(err, "failed to format message")
	}

	return messages, nil
}

// AnonymizeText anonymizes the given text based on the specified entity types.
func (a *Anonymizer) AnonymizeText(ctx context.Context, types []string, text string) (string, []*Entity, error) {
	messages, err := a.createAnonymizeMessages(ctx, types, text)
	if err != nil {
		return "", nil, eris.Wrap(err, "failed to create anonymize messages")
	}

	response, err := a.llm.Generate(ctx, messages)
	if err != nil {
		return "", nil, eris.Wrap(err, "failed to generate response")
	}

	splited := strings.SplitN(response.Content, "<<<PAIR>>>", 2)
	if len(splited) != 2 {
		return "", nil, fmt.Errorf("invalid response format, expected 2 parts but got %d, %s", len(splited), response.Content)
	}

	anonymizedText := strings.TrimSpace(splited[0])
	mappingStr := strings.TrimSpace(splited[1])
	var mapping map[string][]string
	err = json.Unmarshal([]byte(mappingStr), &mapping)
	if err != nil {
		return "", nil, eris.Wrap(err, "failed to unmarshal mapping")
	}

	// key format: <EntityType[ID].Category.Detail>
	keyParseRe := regexp.MustCompile(`<(.+?)\[(.+?)\]\.(.+?)\.(.+?)>`)
	entities := make([]*Entity, 0, len(mapping))
	for key, values := range mapping {
		matches := keyParseRe.FindStringSubmatch(key)
		if len(matches) != 5 {
			return "", nil, fmt.Errorf("invalid key format: %s", key)
		}

		entityType := matches[1]
		id := matches[2]
		category := matches[3]
		detail := matches[4]

		entities = append(entities, &Entity{
			Key:        key,
			EntityType: entityType,
			ID:         id,
			Category:   category,
			Detail:     detail,
			Values:     values,
		})
	}

	return anonymizedText, entities, nil
}

// RestoreText restores the original text from the anonymized text using the provided entities.
func (a *Anonymizer) RestoreText(ctx context.Context, entities []*Entity, text string) (string, error) {
	var replaceMapping []string
	for _, entity := range entities {
		if len(entity.Values) == 0 {
			continue
		}
		replaceMapping = append(replaceMapping, entity.Key, entity.Values[0])
	}

	replacer := strings.NewReplacer(replaceMapping...)
	return replacer.Replace(text), nil
}

// New creates a new Anonymizer instance.
func New(chatModel model.BaseChatModel) (*Anonymizer, error) {
	anonymizeTemplate := prompt.FromMessages(schema.FString,
		schema.UserMessage(`Anonymize the text with the given entity types, then output the tag-to-original mapping; if nothing is found, reply "None".
Specified types: {types}
<text>{text}</text>`),
	)

	return &Anonymizer{
		anonymizeTemplate: anonymizeTemplate,
		llm:               chatModel,
	}, nil
}
