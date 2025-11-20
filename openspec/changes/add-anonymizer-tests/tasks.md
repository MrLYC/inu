# Implementation Tasks

## 1. 设计 Mock 策略
- [x] 1.1 分析 `model.BaseChatModel` 接口结构
- [x] 1.2 确定 mock 实现方式（接口实现而非第三方框架）
- [x] 1.3 设计可配置的响应数据结构

## 2. 实现 LLM Mock
- [ ] 2.1 创建 `pkg/anonymizer/mock_llm_test.go`
- [ ] 2.2 实现 `mockChatModel` 结构体
  - [ ] 实现 `Generate(ctx, messages)` 方法
  - [ ] 支持可配置的响应内容
  - [ ] 支持可配置的错误返回
- [ ] 2.3 实现响应构造辅助函数
  - [ ] `newMockResponse(anonymizedText, mapping)` - 构造标准响应
  - [ ] `newMockErrorResponse(err)` - 构造错误响应

## 3. AnonymizeText 测试
- [ ] 3.1 创建 `pkg/anonymizer/anonymizer_test.go`
- [ ] 3.2 测试正常场景
  - [ ] 单个实体匿名化（如：姓名）
  - [ ] 多个实体匿名化（如：姓名+电话）
  - [ ] 多种实体类型混合
  - [ ] 同类型多个实体（如：多个姓名）
- [ ] 3.3 测试边界情况
  - [ ] 空文本输入
  - [ ] 无匹配实体（LLM 返回 "None"）
  - [ ] 空实体类型列表
- [ ] 3.4 测试错误处理
  - [ ] LLM API 调用失败
  - [ ] 响应格式错误（缺少 `<<<PAIR>>>`）
  - [ ] JSON 解析失败（mapping 格式错误）
  - [ ] 实体 key 格式错误（无法匹配正则）

## 4. RestoreText 测试
- [ ] 4.1 测试正常场景
  - [ ] 单个实体还原
  - [ ] 多个实体还原
  - [ ] 顺序不敏感（实体数组顺序任意）
- [ ] 4.2 测试边界情况
  - [ ] 空文本输入
  - [ ] 空实体列表
  - [ ] 实体无 Values（应跳过）
  - [ ] 文本中不包含任何实体 key（原样返回）
- [ ] 4.3 测试还原准确性
  - [ ] 完整往返测试：原文 -> 匿名化 -> 还原 == 原文（使用 mock）

## 5. New 构造函数测试
- [ ] 5.1 测试成功创建实例
  - [ ] 验证返回的 Anonymizer 不为 nil
  - [ ] 验证内部字段正确初始化
- [ ] 5.2 测试传入 nil chatModel（如果需要边界检查）

## 6. 辅助函数测试
- [ ] 6.1 测试 `createAnonymizeMessages`
  - [ ] 验证生成的消息格式
  - [ ] 验证类型和文本正确嵌入
  - [ ] 测试 JSON 编码错误处理

## 7. 集成和验证
- [ ] 7.1 运行所有测试：`go test ./pkg/anonymizer/...`
- [ ] 7.2 检查测试覆盖率：`go test -cover ./pkg/anonymizer/...`
- [ ] 7.3 确保覆盖率达到 80% 以上（核心逻辑）
- [ ] 7.4 验证测试不调用真实 LLM API（无需环境变量）
- [ ] 7.5 在 CI 中验证测试可以稳定通过

## 8. 文档更新
- [ ] 8.1 在测试文件中添加注释说明 mock 策略
- [ ] 8.2 更新 `README.md` 中的测试部分（如有必要）
- [ ] 8.3 添加测试运行示例到文档
