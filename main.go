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
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/rotisserie/eris"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

type Entity struct {
	Key        string   `json:"key"`
	EntityType string   `json:"type"`
	ID         string   `json:"id"`
	Category   string   `json:"category"`
	Detail     string   `json:"detail"`
	Values     []string `json:"values"`
}

// Has is the entity anonymization example struct.
type Has struct {
	anonymizeTemplate *prompt.DefaultChatTemplate
	llm               model.BaseChatModel
}

// createAnonymizeMessages creates messages for anonymization.
func (h *Has) createAnonymizeMessages(ctx context.Context, types []string, text string) ([]*schema.Message, error) {
	encodedTypes, err := json.Marshal(types)
	if err != nil {
		return nil, eris.Wrap(err, "failed to marshal types")
	}

	messages, err := h.anonymizeTemplate.Format(ctx, map[string]any{
		"types": string(encodedTypes),
		"text":  text,
	})

	if err != nil {
		return nil, eris.Wrap(err, "failed to format message")
	}

	return messages, nil
}

// AnonymizeText anonymizes the given text based on the specified entity types.
func (h *Has) AnonymizeText(ctx context.Context, types []string, text string) (string, []*Entity, error) {
	messages, err := h.createAnonymizeMessages(ctx, types, text)
	if err != nil {
		return "", nil, eris.Wrap(err, "failed to create anonymize messages")
	}

	response, err := h.llm.Generate(ctx, messages)
	if err != nil {
		return "", nil, eris.Wrap(err, "failed to generate response")
	}

	splited := strings.SplitN(response.Content, "<<<PAIR>>>", 2)
	if len(splited) != 2 {
		return "", nil, fmt.Errorf("invalid response format, expected 2 parts but got %d", len(splited))
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
func (h *Has) RestoreText(ctx context.Context, entities []*Entity, text string) (string, error) {
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

// NewHas creates a new Has instance.
func NewHas(chatModel model.BaseChatModel) (*Has, error) {
	anonymizeTemplate := prompt.FromMessages(schema.FString,
		schema.UserMessage(`Anonymize the text with the given entity types, then output the tag-to-original mapping; if nothing is found, reply "None".
Specified types: {types}
<text>{text}</text>`),
	)

	return &Has{
		anonymizeTemplate: anonymizeTemplate,
		llm:               chatModel,
	}, nil
}

// createOpenAIChatModel creates an OpenAI chat model instance.
func createOpenAIChatModel(ctx context.Context) (model.BaseChatModel, error) {
	key := os.Getenv("OPENAI_API_KEY")
	modelName := os.Getenv("OPENAI_MODEL_NAME")
	baseURL := os.Getenv("OPENAI_BASE_URL")
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL: baseURL,
		Model:   modelName,
		APIKey:  key,
	})
	return chatModel, eris.Wrap(err, "failed to create openai chat model")
}

func main() {
	ctx := context.Background()
	llm, err := createOpenAIChatModel(ctx)
	if err != nil {
		log.Fatalf("create chat model failed, err=%v", err)
	}

	has, err := NewHas(llm)
	if err != nil {
		log.Fatalf("create has failed, err=%v", err)
	}

	text := "张三的身份证号是 110101199001011234，他的电话号码是 13800138000。"
	types := []string{"人名", "联系方式", "职务", "密码", "组织", "地址", "文件", "账号", "网址", "IP"}

	result, entities, err := has.AnonymizeText(ctx, types, text)
	if err != nil {
		log.Fatalf("anonymize text failed, err=%v", err)
	}

	log.Printf("anonymize result: %s", result)
	log.Printf("anonymize mapping: %+s", entities)

	restoredText, err := has.RestoreText(ctx, entities, result)
	if err != nil {
		log.Fatalf("restore text failed, err=%v", err)
	}

	log.Printf("restored text: %s", restoredText)
}
