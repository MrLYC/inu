# Design: Interactive Pipe Command

## Overview
本设计文档描述如何实现 `inu interactive` 交互式命令，该命令将脱敏和还原操作结合在一个交互式流程中，支持用户在中间进行文本处理。

## Architecture

### Command Flow
```
交互式模式（单一模式）:
  User Input (原始文本)
      ↓
  Anonymize (LLM)
      ↓
  Output 脱敏文本 (stdout)
  Output 详细提示 (stderr)
      ↓
  Loop: Wait for stdin (用户输入处理后的文本)
      ↓
  当遇到分隔符或 EOF:
      ↓
  Restore (使用内存中的实体)
      ↓
  Output 还原文本 (stdout)
      ↓
  继续等待下一次输入（或用户 Ctrl+C 退出）
```

### Component Structure
```
cmd/inu/commands/
├── anonymize.go (现有)
├── restore.go (现有)
└── interactive.go (新增)
    ├── runInteractive() - 主执行逻辑
    ├── readUserInput() - 读取用户输入直到分隔符/EOF
    └── processAndRestore() - 处理并还原文本
```

## Key Design Decisions

### 1. 命令模式选择
**决策**: 单一交互式模式，无还原模式

**模式**:
- **唯一模式**: 脱敏 → 循环等待输入 → 还原 → 继续等待

**理由**:
- 用户场景明确：人工复制文本去处理，再粘贴回来
- 实体信息始终在内存中，不需要持久化
- 不需要支持管道或分步操作
- 简化命令设计，避免模式切换

**移除的功能**:
- ❌ 不需要 `--restore-mode`
- ❌ 不需要 `--entities` 参数
- ❌ 不需要输出实体信息
- ❌ 不需要管道集成

### 2. 实体信息管理
**决策**: 实体信息仅保存在内存中，不输出

**设计**:
- 脱敏后，实体信息存储在内存
- 用户看到脱敏文本，手动复制去处理
- 用户粘贴处理后的文本回终端
- 命令使用内存中的实体进行还原

**理由**:
- 用户不需要实体信息文件
- 整个流程在一个进程中完成
- 简化输出，只有文本内容

**移除的功能**:
- ❌ 不输出实体到 stdout/stderr
- ❌ 不需要 `--embed-entities` 标志
- ❌ 不需要实体分隔符（`<<<ENTITIES>>>`）
- ❌ 不需要实体序列化/反序列化

### 3. 多次输入循环模式
**决策**: 支持多次输入和还原，类似 heredoc

**流程**:
```go
// 1. 脱敏阶段
anonymizedText, entities := anonymizer.AnonymizeText(...)
fmt.Fprintln(os.Stdout, anonymizedText)

// 2. 输出详细提示
printDetailedPrompt(delimiter)

// 3. 循环等待输入
for {
    processedText := readUntilDelimiterOrEOF(delimiter)
    if processedText == "" {
        break  // EOF 且无内容，退出
    }

    // 4. 还原并输出
    restoredText := anonymizer.RestoreText(entities, processedText)
    fmt.Fprintln(os.Stdout, restoredText)

    // 5. 继续等待下一次输入
    fmt.Fprintln(os.Stderr, "\nReady for next input...")
}
```

**用户交互**:
1. 命令输出脱敏文本
2. 显示详细提示（如何输入、如何结束）
3. 用户粘贴处理后的文本
4. 用户输入分隔符（或 Ctrl+D）
5. 命令输出还原文本
6. 回到步骤 3，支持多次输入
7. 用户 Ctrl+C 或输入 EOF 退出

### 4. 输入分隔符设计
**决策**: 支持自定义分隔符，默认仅 EOF

**分隔符行为**:
- **默认**: 仅 EOF (Ctrl+D) 触发处理
- **自定义**: `--delimiter` 指定分隔符行
  - 遇到该行时触发处理
  - 分隔符行本身不包含在输入中
  - 处理后继续等待下一次输入

**示例**:
```bash
# 默认模式 - 仅 EOF
$ inu interactive -f input.txt
<脱敏文本>
[提示信息]
# 用户粘贴文本，按 Ctrl+D
<还原文本>
[继续等待...]

# 自定义分隔符
$ inu interactive -f input.txt --delimiter "END"
<脱敏文本>
[提示信息: 输入 END 结束]
# 用户粘贴文本
# 用户输入 END
<还原文本>
[继续等待...]
```

**理由**:
- 类似 heredoc 的使用体验
- 支持多次输入和还原
- 用户可选择方便的结束方式

### 5. 错误处理策略
**场景考虑**:
1. 用户输入的处理后文本破坏了占位符格式
2. 用户输入空文本
3. LLM 调用失败

**处理策略**:
- **占位符破坏**: 尽力还原，保留无法匹配的部分，不报错
- **空输入**: 跳过本次，继续等待下一次输入
- **LLM 失败**: 标准错误处理，退出并提示

## Command Specification

### 命令格式
```bash
inu interactive [flags]
```

### Flags

#### 输入相关
- `-f, --file <path>`: 从文件读取原始文本
- `-c, --content <text>`: 从命令行参数读取原始文本
- （默认）：从 stdin 读取原始文本

#### 实体类型
- `-t, --entity-types <types>`: 指定要检测的实体类型（逗号分隔）
- （默认）：使用所有默认实体类型

#### 交互控制
- `--delimiter <text>`: 自定义输入分隔符（默认仅 EOF）
  - 示例: `--delimiter "END"` 表示输入 END 时触发处理
- `--no-prompt`: 禁用交互提示信息（仅保留核心提示）

### Usage Examples

#### Example 1: 基本交互式使用（EOF 模式）
```bash
$ inu interactive -f sensitive.txt

# 输出脱敏文本：
<个人信息[0].姓名.全名> works at <组织机构[0].公司.ABC Tech>

# stderr 提示：
=== Anonymization Complete ===
The text above has been anonymized.
You can now:
1. Copy the anonymized text
2. Process it externally (e.g., paste to ChatGPT)
3. Paste the processed text back here
4. Press Ctrl+D (Unix) or Ctrl+Z (Windows) to restore

Waiting for your input...

# 用户粘贴处理后的文本：
<个人信息[0].姓名.全名> is a great employee at <组织机构[0].公司.ABC Tech>
^D

# 输出还原文本：
张三 is a great employee at ABC Tech

# 继续等待下一次输入：
Ready for next input (Ctrl+D to restore, Ctrl+C to exit)...
```

#### Example 2: 自定义分隔符模式
```bash
$ inu interactive -c "张三的电话是 13800138000" --delimiter "END"

# 输出：
<个人信息[0].姓名.全名>的电话是<个人信息[1].电话.号码>

# 提示：
=== Anonymization Complete ===
Paste processed text below, then type 'END' on a new line to restore.
Press Ctrl+C to exit.

# 用户输入：
这是 <个人信息[0].姓名.全名> 的联系方式
END

# 输出：
这是 张三 的联系方式

# 继续等待：
Ready for next input...
```

#### Example 3: 多次处理
```bash
$ inu interactive -f report.txt --delimiter "==END=="

# 第一次：总结
<脱敏文本>
[粘贴 ChatGPT 总结]
==END==
<还原的总结>

# 第二次：翻译
[粘贴 ChatGPT 翻译]
==END==
<还原的翻译>

# 第三次：...
^C  # 退出
```

## Implementation Details

### Entity Serialization Format
实体信息使用 JSON 数组格式：
```json
[
  {
    "key": "<个人信息[0].姓名.全名>",
    "type": "个人信息",
    "id": "0",
    "category": "姓名",
    "detail": "张三",
    "values": ["张三"]
  }
]
```

### Interactive Input Handling
```go
func readProcessedText() (string, error) {
    fmt.Fprintln(os.Stderr, "Enter processed text (press Ctrl+D when done):")

    var lines []string
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        return "", err
    }

    return strings.Join(lines, "\n"), nil
}
```

### Restore Mode Detection
还原模式下，需要解析实体信息：

**从文件读取**:
```go
if restoreMode && entitiesFile != "" {
    entities, err := cli.LoadEntitiesFromYAML(entitiesFile)
    // ...
}
```

**从 stdin 解析**（嵌入模式）:
```go
func parseEmbeddedEntities(input string) (text string, entities []*Entity, error) {
    // 查找分隔符
    parts := strings.Split(input, "<<<ENTITIES>>>")
    if len(parts) != 2 {
        return input, nil, nil // 无嵌入实体
    }

    text = strings.TrimSpace(parts[0])

    // 解析实体 JSON
    endParts := strings.Split(parts[1], "<<<END_ENTITIES>>>")
    entitiesJSON := strings.TrimSpace(endParts[0])

    var entities []*Entity
    err := json.Unmarshal([]byte(entitiesJSON), &entities)
    return text, entities, err
}
```

## Testing Strategy

### Unit Tests
- `TestInteractiveCommand_Anonymize`: 测试脱敏阶段
- `TestInteractiveCommand_Restore`: 测试还原阶段
- `TestInteractiveCommand_EntitySerialization`: 测试实体序列化
- `TestInteractiveCommand_ParseEmbeddedEntities`: 测试嵌入实体解析
- `TestInteractiveCommand_RestoreMode`: 测试还原模式

### Integration Tests
使用 `io.Pipe` 模拟 stdin/stdout 交互：
```go
func TestInteractiveCommand_Interactive(t *testing.T) {
    // 1. 创建管道模拟 stdin
    r, w := io.Pipe()

    // 2. 启动 goroutine 写入处理后的文本
    go func() {
        w.Write([]byte("processed text\n"))
        w.Close()
    }()

    // 3. 执行命令（使用 r 作为 stdin）
    // 4. 验证输出
}
```

### Manual Testing
```bash
# 测试 1: 基本功能
inu interactive -c "张三的电话是 13800138000"

# 测试 2: 文件输入
echo "敏感信息" > test.txt
inu interactive -f test.txt

# 测试 3: 嵌入模式
inu interactive -c "张三" --embed-entities

# 测试 4: 还原模式
inu anonymize -c "张三" -e /tmp/entities.yaml -o /tmp/anon.txt
cat /tmp/anon.txt | inu interactive --restore-mode --entities /tmp/entities.yaml
```

## User Experience Considerations

### 提示信息
- **脱敏完成**: `Anonymized text (entities saved internally):`
- **等待输入**: `Waiting for processed text (press Ctrl+D when done)...`
- **还原完成**: `Restored text:`
- **错误提示**: 清晰说明问题和解决方法

### 进度指示
所有进度信息输出到 stderr，不干扰数据流：
```
[stderr] Initializing LLM client...
[stderr] Anonymizing text...
[stdout] <anonymized text>
[stderr] <<<ENTITIES>>>
[stderr] {...}
[stderr] <<<END_ENTITIES>>>
[stderr] Waiting for input...
```

### 错误处理
常见错误及提示：
- 缺少输入：`Error: No input provided. Use --file, --content, or stdin.`
- 还原模式缺少实体：`Error: --entities required in --restore-mode.`
- 无效实体格式：`Error: Failed to parse entities: invalid JSON format.`

## Future Enhancements

### 1. 流式处理
支持逐行处理大文件：
```bash
inu interactive --stream -f large.txt
```

### 2. 双向管道支持
探索支持完整管道的可能性（技术难度较高）

### 3. 实体编辑
允许用户在中间阶段修改实体映射：
```bash
inu interactive --interactive-entities
```

### 4. 批处理模式
处理多个文本文件：
```bash
inu interactive --batch files/*.txt
```

## Alternatives Considered

### Alternative 1: 修改 anonymize 命令支持 --wait 标志
```bash
inu anonymize -f input.txt --wait
```
**缺点**:
- 混淆 `anonymize` 命令的单一职责
- 不够直观

### Alternative 2: 创建 workflow 子命令
```bash
inu workflow interactive
```
**缺点**:
- 命令嵌套过深
- 名称不够描述性

### Alternative 3: 使用 TUI（Terminal UI）
创建全屏交互式界面。
**缺点**:
- 实现复杂度高
- 不适合脚本和自动化
- 引入额外依赖

**选择 `pipe` 命令的理由**:
- 名称直观（管道处理）
- 保持命令简单性
- 易于集成和自动化
- 符合 Unix 哲学
