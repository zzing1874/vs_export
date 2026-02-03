### Requirement: 检测输入文件类型
系统 SHALL 根据输入文件的扩展名自动检测文件类型，支持 .sln 和 .vcxproj 两种格式。

#### Scenario: 识别 .sln 文件
- **WHEN** 用户提供的 `-s` 参数路径以 `.sln` 结尾
- **THEN** 系统使用解决方案文件处理流程

#### Scenario: 识别 .vcxproj 文件
- **WHEN** 用户提供的 `-s` 参数路径以 `.vcxproj` 结尾
- **THEN** 系统使用项目文件处理流程

#### Scenario: 不支持的文件扩展名
- **WHEN** 用户提供的文件扩展名既不是 `.sln` 也不是 `.vcxproj`
- **THEN** 系统返回错误信息，说明仅支持 .sln 和 .vcxproj 文件

### Requirement: 文件类型不区分大小写
系统 SHALL 对文件扩展名进行不区分大小写的匹配，`.SLN`、`.Sln`、`.sln` 应被视为相同。

#### Scenario: 大写扩展名
- **WHEN** 用户提供 `project.VCXPROJ` 作为输入
- **THEN** 系统正确识别为 .vcxproj 文件并处理

#### Scenario: 混合大小写扩展名
- **WHEN** 用户提供 `solution.Sln` 作为输入
- **THEN** 系统正确识别为 .sln 文件并处理

### Requirement: 处理流程分发
系统 SHALL 根据检测到的文件类型，将请求分发到相应的处理流程。

#### Scenario: .sln 文件使用现有流程
- **WHEN** 检测到输入为 .sln 文件
- **THEN** 系统调用 `sln.NewSln()` 解析解决方案文件，遍历其中的所有项目

#### Scenario: .vcxproj 文件使用单项目流程
- **WHEN** 检测到输入为 .vcxproj 文件
- **THEN** 系统直接调用 `project.NewProject()` 解析该项目文件，跳过解决方案解析

### Requirement: 错误消息清晰性
当文件类型不支持时，系统 SHALL 提供清晰的错误消息，明确说明支持的文件类型。

#### Scenario: 提示支持的文件类型
- **WHEN** 用户提供 `project.txt` 作为输入
- **THEN** 错误消息包含类似 "Only .sln and .vcxproj files are supported" 的说明
