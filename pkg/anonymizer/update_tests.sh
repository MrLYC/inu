#!/bin/bash
# 更新测试文件以适配新接口

# 1. 替换 AnonymizeText 调用为 Anonymize + bytes.Buffer
sed -i '' 's/result, entities, err := anon\.AnonymizeText(ctx, \[\]string{\([^}]*\)}, \([^)]*\))/var buf bytes.Buffer\
entities, err := anon.Anonymize(ctx, []string{\1}, \2, \&buf)\
result := buf.String()/g' anonymizer_test.go

# 2. 替换只有错误检查的 AnonymizeText 调用
sed -i '' 's/_, _, err = anon\.AnonymizeText(ctx, \[\]string{\([^}]*\)}, \([^)]*\))/var buf bytes.Buffer\
_, err = anon.Anonymize(ctx, []string{\1}, \2, \&buf)/g' anonymizer_test.go
