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
	"log"

	"github.com/mrlyc/inu/pkg/anonymizer"
)

func main() {
	ctx := context.Background()
	llm, err := anonymizer.CreateOpenAIChatModel(ctx)
	if err != nil {
		log.Fatalf("create chat model failed, err=%v", err)
	}

	anon, err := anonymizer.New(llm)
	if err != nil {
		log.Fatalf("create anonymizer failed, err=%v", err)
	}

	text := "张三的身份证号是 110101199001011234，他的电话号码是 13800138000。"
	types := []string{"人名", "联系方式", "职务", "密码", "组织", "地址", "文件", "账号", "网址", "IP"}

	result, entities, err := anon.AnonymizeText(ctx, types, text)
	if err != nil {
		log.Fatalf("anonymize text failed, err=%v", err)
	}

	log.Printf("anonymize result: %s", result)
	log.Printf("anonymize mapping: %+v", entities)

	restoredText, err := anon.RestoreText(ctx, entities, result)
	if err != nil {
		log.Fatalf("restore text failed, err=%v", err)
	}

	log.Printf("restored text: %s", restoredText)
}
