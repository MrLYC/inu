# cli Specification Delta

## MODIFIED Requirements

### Requirement: 还原命令
系统 SHALL 提供 `restore` 子命令来还原脱敏的文本，支持占位符格式变化的容错匹配。

#### Scenario: 从标准输入读取匿名文本并还原
- **WHEN** 用户执行 `echo "<个人信息[0].姓名.全名>" | inu restore --entities entities.yaml`
- **THEN** 系统应该读取标准输入和实体文件，还原文本并输出到标准输出（默认行为）

#### Scenario: 从文件读取匿名文本
- **WHEN** 用户执行 `inu restore --file anonymized.txt --entities entities.yaml`
- **THEN** 系统应该读取文件内容进行还原并输出到标准输出

#### Scenario: 从命令行参数读取匿名文本
- **WHEN** 用户执行 `inu restore --content "<个人信息[0].姓名.全名>" --entities entities.yaml`
- **THEN** 系统应该还原提供的内容字符串并输出到标准输出

#### Scenario: 输入优先级
- **WHEN** 用户同时指定 `--file`、`--content` 和标准输入
- **THEN** 系统应该按优先级使用：`--file` > `--content` > 标准输入

#### Scenario: 输出还原文本到文件
- **WHEN** 用户执行 `inu restore --file anonymized.txt --entities entities.yaml --output restored.txt`
- **THEN** 系统应该将还原后的文本写入文件
- **AND** 同时输出到标准输出（默认行为）

#### Scenario: 只输出到文件不显示
- **WHEN** 用户执行 `inu restore --content "text" --entities entities.yaml --no-print --output restored.txt`
- **THEN** 系统应该只将还原后的文本写入文件，不输出到标准输出

#### Scenario: 缺少实体文件时报错
- **WHEN** 用户执行 `inu restore --content "text"` 但未指定 `--entities`
- **THEN** 系统应该退出并显示错误：需要提供实体文件

#### Scenario: 实体文件格式错误时报错
- **WHEN** 用户提供的 `--entities` 文件不是有效的 YAML 格式
- **THEN** 系统应该退出并显示清晰的解析错误信息

#### Scenario: 还原包含额外空格的占位符
- **WHEN** 用户执行 `echo "< 个人信息 [0]. 姓名. 张三 >" | inu restore --entities entities.yaml`
- **AND** 实体文件包含标准格式的键 `<个人信息[0].姓名.全名>`
- **THEN** 系统应该成功匹配并还原为原始值 "张三"
- **AND** 输出应该是 "张三"（不是占位符）

#### Scenario: 还原包含中文标点的占位符
- **WHEN** 用户执行 `echo "<业务信息[2]。系统。名称>" | inu restore --entities entities.yaml`
- **AND** 实体文件包含标准格式的键 `<业务信息[2].系统.名称>`
- **THEN** 系统应该成功匹配并还原（中文标点 `。` 被归一化为 `.`）
- **AND** 输出应该是原始值（不是占位符）

#### Scenario: 还原包含全角字符的占位符
- **WHEN** 用户执行 `echo "<账户信息[　０　].银行账户.６２２２０２１００１１２３４５６７８９>" | inu restore --entities entities.yaml`
- **AND** 实体文件包含标准格式的键 `<账户信息[0].银行账户.6222021001123456789>`（半角数字）
- **THEN** 系统应该成功匹配并还原（全角字符被归一化为半角）
- **AND** 输出应该是原始值（不是占位符）

#### Scenario: 还原混合格式变化的占位符
- **WHEN** 用户执行 `echo "< 业务信息 [2]。 系统 。 名称 >" | inu restore --entities entities.yaml`
- **AND** 实体文件包含标准格式的键 `<业务信息[2].系统.名称>`
- **THEN** 系统应该成功匹配并还原（同时处理空格、中文标点）
- **AND** 输出应该是原始值（不是占位符）

#### Scenario: 部分占位符无法匹配时的行为
- **WHEN** 用户执行 `echo "<个人信息[0].姓名.全名> and < unknown >" | inu restore --entities entities.yaml`
- **AND** 实体文件只包含 `<个人信息[0].姓名.全名>`
- **THEN** 系统应该还原已知占位符为 "张三"
- **AND** 未知占位符保留原样 "< unknown >"
- **AND** 输出应该是 "张三 and < unknown >"

#### Scenario: 归一化不影响标准格式占位符
- **WHEN** 用户执行 `echo "<个人信息[0].姓名.全名>" | inu restore --entities entities.yaml`
- **AND** 实体文件包含标准格式的键 `<个人信息[0].姓名.全名>`
- **THEN** 系统应该正常还原（归一化对标准格式无影响）
- **AND** 输出应该是 "张三"
- **AND** 行为与之前版本完全一致（向后兼容）
