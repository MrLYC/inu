# CLI Spec Delta: 统一流式接口

## MODIFIED: Requirement: 脱敏命令 (Line 114)

### Change Summary
- 接口方法重命名: `AnonymizeTextStream()` → `Anonymize()`
- 行为保持不变(流式输出)

### Updated Scenarios
无需更新 scenarios,因为命令行为完全一致,只是内部接口调用变化。

---

## MODIFIED: Requirement: 还原命令 (Line 204)

### Change Summary
- 添加详细的失败占位符警告输出
- 新增返回值: 未还原的占位符列表
- 警告输出到 stderr,不影响 stdout 的主输出

### New Scenario: 部分占位符无法还原时显示警告

#### Scenario: 显示无法还原的占位符
- **WHEN** 用户执行 `inu restore --entities entities.yaml --content "{{PERSON_0}} met {{PERSON_1}}"`
- **AND** 实体文件只包含 `PERSON_0` 的定义
- **THEN** 系统应该：
  1. 将还原后的文本输出到 stdout: `Alice met {{PERSON_1}}`
  2. 在 stderr 显示警告:
     ```
     Warning: 1 placeholder(s) could not be restored:
       - PERSON_1
     ```
  3. 返回退出码 0 (部分失败不视为错误)

#### Scenario: 全部占位符成功还原
- **WHEN** 用户执行 `inu restore --entities entities.yaml --content "{{PERSON_0}} met {{PERSON_1}}"`
- **AND** 实体文件包含所有占位符的定义
- **THEN** 系统应该：
  1. 将完全还原的文本输出到 stdout
  2. 不显示任何警告
  3. 返回退出码 0

#### Scenario: 无占位符的文本
- **WHEN** 用户执行 `inu restore --entities entities.yaml --content "plain text"`
- **THEN** 系统应该：
  1. 直接输出原文本到 stdout
  2. 不显示任何警告
  3. 返回退出码 0

#### Scenario: 多个占位符失败时列出所有失败项
- **WHEN** 用户执行还原命令,且有 3 个占位符无法还原
- **THEN** stderr 应该显示:
  ```
  Warning: 3 placeholder(s) could not be restored:
    - PERSON_1
    - PERSON_2
    - ORG_5
  ```

#### Scenario: 警告不影响管道输出
- **WHEN** 用户执行 `inu restore --entities entities.yaml --content "..." | grep "Alice"`
- **AND** 有部分占位符无法还原
- **THEN** 系统应该：
  1. stdout 传递给管道(只包含还原后的文本)
  2. 警告输出到 stderr(不影响管道)
  3. 用户可以使用 `2>/dev/null` 抑制警告

#### Scenario: 输出到文件时仍显示警告
- **WHEN** 用户执行 `inu restore --entities entities.yaml --content "..." --output restored.txt`
- **AND** 有部分占位符无法还原
- **THEN** 系统应该：
  1. 将还原后的文本写入文件和 stdout
  2. 在 stderr 显示警告(无论输出目标是什么)

### Design Notes
- **为什么返回退出码 0?**: 部分还原失败是正常业务场景,不应视为命令执行失败
- **为什么输出到 stderr?**: 保持 stdout 纯净,便于管道和重定向
- **警告格式**: 清晰列出每个失败的占位符,便于用户检查实体文件

---

## MODIFIED: Requirement: 交互式命令 (Line 297)

### Change Summary
- 脱敏阶段: 接口方法重命名 `AnonymizeTextStream()` → `Anonymize()`
- 还原阶段: 添加详细的失败占位符警告输出
- 行为与 restore 命令保持一致

### New Scenario: 还原阶段显示失败占位符警告

#### Scenario: 还原时部分占位符失败
- **WHEN** 用户在交互式模式输入处理后的文本
- **AND** 输入包含无法还原的占位符
- **THEN** 系统应该：
  1. 将还原后的文本输出到 stdout
  2. 在 stderr 显示警告:
     ```
     Warning: 2 placeholder(s) could not be restored:
       - PERSON_1
       - ORG_3
     ```
  3. 显示 "Ready for next input..." 提示
  4. 继续等待下一次输入

#### Scenario: 还原全部成功
- **WHEN** 用户输入的处理后文本所有占位符都能还原
- **THEN** 系统应该：
  1. 输出完全还原的文本到 stdout
  2. 不显示任何警告
  3. 显示 "Ready for next input..." 提示

#### Scenario: 多次输入处理中的警告
- **WHEN** 用户进行多次输入处理
- **AND** 第一次输入有失败占位符,第二次输入全部成功
- **THEN** 系统应该：
  1. 第一次处理后显示警告
  2. 第二次处理后不显示警告
  3. 每次处理独立,警告不累积

### Updated Scenario: 详细提示信息 (Line 333)

#### Scenario: 详细提示信息
- **WHEN** 命令启动并完成脱敏
- **THEN** stderr 应该输出类似：
  ```
  === Anonymization Complete ===
  The text above has been anonymized.

  Now enter your processed text below.
  When finished, press Ctrl+D (EOF) to restore.
  ```
- **AND** 用户输入处理后的文本并触发还原
- **AND** 如果有占位符无法还原,显示警告:
  ```
  Warning: 1 placeholder(s) could not be restored:
    - PERSON_2

  Ready for next input...
  ```

### Design Notes
- **交互式体验**: 警告在还原输出后立即显示,紧接 "Ready for next input..." 提示
- **一致性**: 警告格式与 restore 命令完全一致
- **不中断流程**: 警告不会停止交互循环,用户可以继续下一次输入

---

## Implementation Impact

### Affected Commands
1. **anonymize**: 内部接口调用变化,无用户可见变化
2. **restore**: 新增警告输出功能,改善用户体验
3. **interactive**: 新增警告输出功能,保持与 restore 一致

### Backward Compatibility
- ✅ **命令行参数**: 无变化
- ✅ **标准输出格式**: 无变化(主输出保持纯净)
- ✅ **退出码**: 无变化(部分失败仍返回 0)
- ✅ **管道兼容**: 警告输出到 stderr,不影响管道
- ⚠️ **标准错误输出**: 新增警告信息(可能影响解析 stderr 的脚本)

### Migration Guide for Users

**场景 1: 脚本中使用 restore 命令**
```bash
# 旧行为
inu restore --entities entities.yaml --content "{{PERSON_0}}" > output.txt
# stdout: 还原后的文本
# stderr: (空)

# 新行为
inu restore --entities entities.yaml --content "{{PERSON_0}} and {{PERSON_1}}" > output.txt
# stdout: 还原后的文本(不变)
# stderr: Warning: 1 placeholder(s) could not be restored: ...

# 如果需要抑制警告
inu restore --entities entities.yaml --content "..." 2>/dev/null > output.txt
```

**场景 2: 管道中使用 restore**
```bash
# 旧行为和新行为完全兼容
echo "{{PERSON_0}}" | inu restore --entities entities.yaml | grep "Alice"
# 警告自动输出到 stderr,不影响管道
```

**场景 3: 检查是否有失败的占位符**
```bash
# 新功能: 捕获 stderr 检查警告
inu restore --entities entities.yaml --content "..." 2> errors.txt
if [ -s errors.txt ]; then
    echo "有占位符无法还原,请检查:"
    cat errors.txt
fi
```

### Testing Requirements

#### Unit Tests
- `pkg/anonymizer/anonymizer_test.go`:
  - TestRestoreText_WithUnrestoredPlaceholders
  - TestRestoreText_AllRestored
  - TestRestoreText_NoPlaceholders

#### Integration Tests
- `cmd/inu/commands/restore_test.go`:
  - 验证警告输出格式
  - 验证 stdout 和 stderr 分离
  - 验证退出码保持为 0

- `cmd/inu/commands/interactive_test.go`:
  - 验证还原阶段的警告输出
  - 验证多次输入的独立处理

---

## References

- **Original Spec**: openspec/specs/cli/spec.md
- **Proposal**: openspec/changes/unify-streaming-interface/proposal.md
- **Design**: openspec/changes/unify-streaming-interface/design.md
- **Tasks**: openspec/changes/unify-streaming-interface/tasks.md
