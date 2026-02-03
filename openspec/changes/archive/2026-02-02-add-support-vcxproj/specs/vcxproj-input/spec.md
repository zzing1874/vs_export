## ADDED Requirements

### Requirement: 支持 .vcxproj 文件作为输入
系统 SHALL 接受 .vcxproj 文件路径作为 `-s` 参数的输入，并生成对应的 compile_commands.json 文件。

#### Scenario: 使用 .vcxproj 文件生成编译命令
- **WHEN** 用户执行 `vs_export -s project.vcxproj -c "Debug|Win32"`
- **THEN** 系统解析该 .vcxproj 文件并生成包含所有编译命令的 compile_commands.json

#### Scenario: .vcxproj 文件路径不存在
- **WHEN** 用户提供的 .vcxproj 文件路径不存在
- **THEN** 系统返回错误信息，提示文件未找到

#### Scenario: .vcxproj 文件格式错误
- **WHEN** 提供的 .vcxproj 文件不是有效的 XML 格式
- **THEN** 系统返回错误信息，提示文件格式错误

### Requirement: 解析项目配置
系统 SHALL 从 .vcxproj 文件中提取指定配置的编译信息，包括源文件列表、include 路径和预处理器定义。

#### Scenario: 提取有效配置
- **WHEN** 用户指定的配置（如 "Debug|Win32"）在 .vcxproj 文件中存在
- **THEN** 系统提取该配置下的 AdditionalIncludeDirectories 和 PreprocessorDefinitions

#### Scenario: 指定的配置不存在
- **WHEN** 用户指定的配置在 .vcxproj 文件中不存在
- **THEN** 系统返回错误信息，列出可用的配置列表

### Requirement: 处理项目变量
系统 SHALL 正确展开 .vcxproj 文件中的项目变量，包括 `$(ProjectDir)`、`$(Configuration)` 和 `$(Platform)`。

#### Scenario: 展开项目目录变量
- **WHEN** .vcxproj 文件的 include 路径包含 `$(ProjectDir)`
- **THEN** 系统将其替换为 .vcxproj 文件所在的绝对路径

#### Scenario: 展开配置变量
- **WHEN** include 路径包含 `$(Configuration)` 或 `$(Platform)`
- **THEN** 系统将其替换为用户通过 `-c` 参数指定的对应值

### Requirement: 处理解决方案变量
当输入为单个 .vcxproj 文件时，系统 SHALL 将 `$(SolutionDir)` 变量设置为项目文件所在目录。

#### Scenario: 展开 SolutionDir 变量
- **WHEN** .vcxproj 文件中的路径包含 `$(SolutionDir)`
- **THEN** 系统将其替换为 .vcxproj 文件所在目录的绝对路径

### Requirement: 生成编译命令
系统 SHALL 为 .vcxproj 中的每个源文件生成一条编译命令，格式与从 .sln 文件生成的一致。

#### Scenario: 生成完整的编译命令
- **WHEN** 系统处理一个包含 3 个 .cpp 源文件的 .vcxproj
- **THEN** compile_commands.json 包含 3 条编译命令，每条包含 directory、file 和 command 字段

#### Scenario: 编译命令格式一致性
- **WHEN** 同一项目分别通过 .sln 和 .vcxproj 文件处理
- **THEN** 生成的编译命令的 include 路径和预处理器定义应该相同
