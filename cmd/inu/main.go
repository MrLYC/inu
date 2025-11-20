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

	types := []string{"个人信息", "业务信息", "资产信息", "账户信息", "位置数据", "文档名称", "组织机构", "岗位称谓"}
	text := `
	张三的身份证号是 110101199001011234，他的电话号码是 13800138000。 
	他住在北京市朝阳区。张三的银行账户是 6222021001123456789，电子邮箱是 zhangsan@example.com。
	他的信用卡号是 4111111111111111，有效期到 12/25，CVV 码是 123。
	张三在公司 ABC Tech 工作，职位是 软件工程师，员工编号是 E12345。
	他最近购买了一辆车，车牌号是 京A12345，车型是 特斯拉 Model 3。
	张三的上级是 李四，职位是 技术经理，张三一般叫他 老李。
	`

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
