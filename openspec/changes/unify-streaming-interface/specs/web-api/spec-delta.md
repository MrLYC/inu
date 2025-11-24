# Web API Spec Delta: Unify Streaming Interface

## Changes to Web API Spec

### 1. Restore Endpoint - Enhanced Response

**Location**: `specs/web-api/spec.md` - `/api/v1/restore` Endpoint

**Current Response**:
```json
{
  "restored_text": "string"
}
```

**New Response**:
```json
{
  "restored_text": "string",
  "unrestored_placeholders": [
    {
      "placeholder": "<个人信息[1].姓名.全名>",
      "reason": "not_found"
    },
    {
      "placeholder": "<组织机构[2].公司.名称>",
      "reason": "empty_values"
    }
  ]
}
```

**Field Descriptions**:
- `restored_text`: The restored text with successful replacements
- `unrestored_placeholders`: Optional array of failed restorations (omitted if empty)
  - `placeholder`: The normalized placeholder that couldn't be restored
  - `reason`: Failure reason
    - `"not_found"`: Placeholder not found in provided entities
    - `"empty_values"`: Entity exists but has no values

### 2. Restore Endpoint - Response Schema

**Location**: `specs/web-api/spec.md` - Response Schemas Section

**New Schema**:
```json
{
  "RestoreResponse": {
    "type": "object",
    "required": ["restored_text"],
    "properties": {
      "restored_text": {
        "type": "string",
        "description": "The restored text with placeholders replaced by original values"
      },
      "unrestored_placeholders": {
        "type": "array",
        "description": "List of placeholders that could not be restored (omitted if empty)",
        "items": {
          "$ref": "#/components/schemas/RestoreFailure"
        }
      }
    }
  },
  "RestoreFailure": {
    "type": "object",
    "required": ["placeholder", "reason"],
    "properties": {
      "placeholder": {
        "type": "string",
        "description": "The normalized placeholder that failed to restore",
        "example": "<个人信息[1].姓名.全名>"
      },
      "reason": {
        "type": "string",
        "enum": ["not_found", "empty_values"],
        "description": "Reason for restoration failure"
      }
    }
  }
}
```

## Backward Compatibility

✅ **Fully backward compatible**:
- Existing fields unchanged
- New field `unrestored_placeholders` is optional (omitted if empty)
- HTTP status codes unchanged
- Request format unchanged

Clients can:
1. **Ignore new field**: Continue using `restored_text` only
2. **Opt-in**: Check `unrestored_placeholders` for enhanced error handling

## Examples

### Example 1: Full Success
```json
POST /api/v1/restore
{
  "anonymized_text": "<个人信息[0].姓名.全名>在<地理位置[1].城市.名称>工作",
  "entities": [
    {
      "key": "<个人信息[0].姓名.全名>",
      "values": ["张三"]
    },
    {
      "key": "<地理位置[1].城市.名称>",
      "values": ["北京"]
    }
  ]
}

Response: 200 OK
{
  "restored_text": "张三在北京工作"
}
```
Note: `unrestored_placeholders` is omitted when empty.

### Example 2: Partial Failures
```json
POST /api/v1/restore
{
  "anonymized_text": "<个人信息[0].姓名.全名> <个人信息[999].电话.手机> <组织机构[5].公司.名称>",
  "entities": [
    {
      "key": "<个人信息[0].姓名.全名>",
      "values": ["张三"]
    },
    {
      "key": "<组织机构[5].公司.名称>",
      "values": []
    }
  ]
}

Response: 200 OK
{
  "restored_text": "张三 <个人信息[999].电话.手机> <组织机构[5].公司.名称>",
  "unrestored_placeholders": [
    {
      "placeholder": "<个人信息[999].电话.手机>",
      "reason": "not_found"
    },
    {
      "placeholder": "<组织机构[5].公司.名称>",
      "reason": "empty_values"
    }
  ]
}
```

### Example 3: Empty Entities
```json
POST /api/v1/restore
{
  "anonymized_text": "plain text without placeholders",
  "entities": []
}

Response: 200 OK
{
  "restored_text": "plain text without placeholders"
}
```

## Client Implementation Guide

### Handling Unrestored Placeholders

**JavaScript Example**:
```javascript
const response = await fetch('/api/v1/restore', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify(request)
});

const data = await response.json();

console.log('Restored:', data.restored_text);

if (data.unrestored_placeholders?.length > 0) {
  console.warn('Failed to restore:');
  data.unrestored_placeholders.forEach(f => {
    const reason = f.reason === 'not_found'
      ? 'not found in entities'
      : 'entity has no values';
    console.warn(`  - ${f.placeholder} (${reason})`);
  });
}
```

**Python Example**:
```python
response = requests.post('/api/v1/restore', json=request)
data = response.json()

print(f"Restored: {data['restored_text']}")

if failures := data.get('unrestored_placeholders'):
    print(f"\nWarning: {len(failures)} placeholder(s) could not be restored:")
    for f in failures:
        reason = ('not found in entities' if f['reason'] == 'not_found'
                  else 'entity has no values')
        print(f"  - {f['placeholder']} ({reason})")
```

## HTTP Status Codes

**No changes** to status code behavior:
- `200 OK`: Restoration completed (even with partial failures)
- `400 Bad Request`: Invalid request format
- `500 Internal Server Error`: Processing error

Partial restoration failures are **warnings**, not errors, so they return 200 OK.

## Updated Sections

Add to `specs/web-api/spec.md`:

1. **Restore Endpoint** - Update response schema
2. **Response Schemas** - Add `RestoreFailure` schema
3. **Examples** - Add partial failure examples
4. **Client Guide** - Add handling recommendations
