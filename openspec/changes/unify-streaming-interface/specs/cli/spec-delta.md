# CLI Spec Delta: Unify Streaming Interface

## Changes to CLI Spec

### 1. Restore Command - Enhanced Output

**Location**: `specs/cli/spec.md` - Restore Command Section

**Current Behavior**:
- Restores text silently
- No indication of failures

**New Behavior**:
- Displays warnings for unrestored placeholders
- Distinguishes between two failure types:
  1. **Not Found**: Placeholder not found in entities file
  2. **Empty Values**: Entity exists but has no values

**Output Format**:
```
Warning: 2 placeholder(s) could not be restored:
  - <个人信息[1].姓名.全名> (not found in entities file)
  - <组织机构[2].公司.名称> (entity has no values)
```

**Exit Code**:
- Returns `0` even with partial failures (warnings, not errors)

### 2. Interactive Command - Enhanced Feedback

**Location**: `specs/cli/spec.md` - Interactive Command Section

**Current Behavior**:
- Restores text without feedback on failures

**New Behavior**:
- Shows same warning format as restore command
- Warnings appear after restored text output

### 3. Anonymize Command - Interface Rename

**Location**: `specs/cli/spec.md` - Anonymize Command Section

**Change**:
- Internal interface method renamed from `AnonymizeTextStream` → `Anonymize`
- No user-facing changes

## Technical Implementation

### Streaming Output
All restore operations now use streaming output via `io.Writer`:
- Reduces memory usage for large texts
- Enables real-time output display
- Supports multiple output targets (stdout, file, MultiWriter)

### Failure Tracking
```go
type RestoreFailure struct {
    Placeholder string // e.g., "<个人信息[1].姓名.全名>"
    Reason      string // "not_found" | "empty_values"
}
```

### CLI Warning Display
Warnings are written to `stderr` to separate them from main output, allowing:
```bash
inu restore -i anonymized.txt -e entities.json > restored.txt 2> warnings.txt
```

## Backward Compatibility

✅ **Fully backward compatible**:
- Exit codes unchanged (0 for success)
- Output format unchanged (restored text)
- Additional warnings are opt-out (redirect stderr)
- Entity file format unchanged

## Examples

### Restore with Failures
```bash
$ inu restore -i input.txt -e entities.json

这是一段文本 <个人信息[999].姓名.全名> <组织机构[5].公司.名称>

Warning: 2 placeholder(s) could not be restored:
  - <个人信息[999].姓名.全名> (not found in entities file)
  - <组织机构[5].公司.名称> (entity has no values)
```

### Interactive Mode
```bash
$ inu interactive

Enter text to anonymize (Ctrl+D to finish):
张三在北京工作
^D

Anonymized text:
<个人信息[0].姓名.全名>在<地理位置[1].城市.名称>工作

Enter anonymized text to restore (Ctrl+D to finish):
<个人信息[0].姓名.全名>在<地理位置[1].城市.名称>工作
^D

Restored text:
张三在北京工作
```

## Updated Sections

Add to `specs/cli/spec.md`:

1. **Restore Command** - Add "Failure Warnings" subsection
2. **Interactive Command** - Add "Restoration Feedback" subsection
3. **Error Handling** - Document warning vs error distinction
