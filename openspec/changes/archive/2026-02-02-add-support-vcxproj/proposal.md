## Why

当前工具只支持以 .sln 文件作为输入，但在某些场景下用户只有单个 .vcxproj 项目文件而没有完整的解决方案文件。增加对 .vcxproj 文件的直接支持，可以让用户在这些场景下也能生成 compile_commands.json。

## What Changes

- 新增对 .vcxproj 文件格式的识别和直接处理能力
- 修改命令行参数解析逻辑，支持 `-s` 参数接受 .vcxproj 文件路径
- 当输入为 .vcxproj 文件时，跳过 .sln 文件解析流程，直接解析项目文件
- 保持现有 .sln 文件处理逻辑不变，向后兼容

## Capabilities

### New Capabilities
- `vcxproj-input`: 支持直接以 .vcxproj 文件作为输入，解析项目配置并生成编译命令

### Modified Capabilities
- `file-detection`: 增强文件类型检测逻辑，区分 .sln 和 .vcxproj 文件并选择相应的处理流程

## Impact

**代码影响：**
- `main.go`: 修改命令行参数处理和文件类型检测逻辑
- `sln/` 包: 可能需要重构以支持单个项目文件的处理流程

**用户影响：**
- 用户现在可以使用 `vs_export -s project.vcxproj -c "Debug|x64"` 直接处理单个项目文件
- 保持现有 .sln 文件使用方式不变，完全向后兼容

**依赖：**
- 无新增外部依赖
