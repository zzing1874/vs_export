## Context

当前 `vs_export` 工具的架构：
- `main.go` 接受 `-s` 参数指定 .sln 文件路径
- `sln.NewSln()` 解析 .sln 文件，提取其中引用的所有 .vcxproj 文件路径
- `project.NewProject()` 解析每个 .vcxproj 文件，提取编译配置（include 路径、预处理器定义等）
- `sln.CompileCommandsJson()` 遍历所有项目，生成 compile_commands.json

实际上，.vcxproj 文件已经包含了生成 compile_commands.json 所需的所有信息（源文件列表、include 路径、预处理器定义等）。.sln 文件的主要作用是：
1. 聚合多个项目
2. 提供 `$(SolutionDir)` 变量

对于单个 .vcxproj 文件的使用场景，不需要 .sln 文件也能生成完整的编译命令。

## Goals / Non-Goals

**Goals:**
- 支持 `-s` 参数直接接受 .vcxproj 文件路径
- 保持现有 .sln 文件处理逻辑和行为不变
- 最小化代码变动，复用现有的 `Project` 结构和解析逻辑

**Non-Goals:**
- 不修改 .vcxproj 文件的解析逻辑（`project.go` 的 XML 解析保持不变）
- 不改变 compile_commands.json 的输出格式
- 不引入新的命令行参数

## Decisions

### 决策 1：文件类型检测位置
**决策：** 在 `main.go` 中根据文件扩展名检测输入类型

**理由：**
- 简单直接，使用 `filepath.Ext()` 或字符串检查即可
- 保持 `sln` 包的职责单一，不需要修改包接口
- 用户输入验证应该在入口处完成

**备选方案：** 在 `sln.NewSln()` 内部自动检测
- 缺点：混淆了 `Sln` 结构的语义（不再只表示解决方案）

### 决策 2：处理 .vcxproj 文件的方式
**决策：** 创建一个只包含单个项目的 `Sln` 实例

**理由：**
- 复用现有的 `sln.CompileCommandsJson()` 逻辑，无需重复编写
- `Sln` 结构已经支持项目列表，单个项目只是特殊情况
- 最小化代码变动，降低引入 bug 的风险

**备选方案 1：** 直接调用 `project.NewProject()` 并单独实现导出逻辑
- 缺点：需要复制 `CompileCommandsJson()` 中的逻辑（环境变量替换、格式转换等）

**备选方案 2：** 重构 `sln` 包，提取项目处理逻辑为独立函数
- 缺点：重构范围较大，增加变更风险

### 决策 3：`$(SolutionDir)` 变量的处理
**决策：** 当输入为 .vcxproj 时，将 `$(SolutionDir)` 设置为项目文件所在目录

**理由：**
- 对于单项目场景，项目目录和解决方案目录是同一概念
- 保证包含 `$(SolutionDir)` 的路径引用仍能正确解析
- 避免引入未定义变量错误

**备选方案：** 设置为空或不处理
- 缺点：如果 .vcxproj 文件中使用了 `$(SolutionDir)` 变量会导致路径错误

### 决策 4：代码组织方式
**决策：** 在 `main.go` 中添加文件类型判断分支

**实现方式：**
```go
if filepath.Ext(*path) == ".vcxproj" {
    // 直接解析项目文件
    pro, err := sln.NewProject(*path)
    // 手动构造 Sln 实例
    solution := sln.Sln{
        SolutionDir: filepath.Dir(absPath),
        ProjectList: []sln.Project{pro},
    }
} else {
    // 现有的 .sln 解析逻辑
    solution, err := sln.NewSln(*path)
}
```

**理由：**
- 清晰的分支逻辑，易于理解和维护
- 不影响 `sln` 包的公共接口
- 后续处理流程完全一致

## Risks / Trade-offs

### [风险] 包访问限制
`sln.Sln` 结构的字段可能是私有的，无法在 `main.go` 中直接构造

**缓解措施：** 
- 检查 `sln.Sln` 的字段可见性（当前是公开的，`SolutionDir` 和 `ProjectList` 都是大写开头）
- 如果需要，可以在 `sln` 包中添加一个辅助构造函数 `NewSlnFromProject()`

### [风险] `$(SolutionDir)` 语义不一致
对于某些复杂项目结构，项目目录和解决方案目录可能不同

**缓解措施：**
- 文档中说明当使用 .vcxproj 作为输入时，`$(SolutionDir)` 将被设置为项目所在目录
- 如果用户需要自定义解决方案目录，仍然应该使用 .sln 文件

### [权衡] 不添加专用命令行参数
不引入 `-p` 参数专门用于 .vcxproj 文件

**理由：**
- 保持接口简洁，用户不需要记忆多个参数
- 自动检测文件类型对用户更友好
- 权衡：错误消息需要更明确，提示支持的文件类型
